package version

import (
	"bytes"
	"errors"
	"io"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aziontech/azion-cli/utils"
	"github.com/tidwall/gjson"
)

// ErrBuildFailed is what BuildDone wraps when a build settles without a
// usable version.
var ErrBuildFailed = errors.New("build did not produce a usable version")

// StatusError carries the HTTP status of a failed v6 call next to the error
// the CLI already shows for it, so a wait can tell a call worth repeating from
// one that never will succeed. Unwrap returns the shown error, so errors.Is
// against the utils sentinels keeps working.
type StatusError struct {
	StatusCode int
	// RetryAfter is the delay the API asked for with Retry-After, or zero.
	RetryAfter time.Duration
	Err        error
}

func (e *StatusError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("unexpected response from the API (status %d)", e.StatusCode)
	}
	return e.Err.Error()
}
func (e *StatusError) Unwrap() error { return e.Err }

// NewStatusError wraps shown with the status and Retry-After of resp. Without
// a response (the request never got one) it returns shown unchanged. It never
// returns nil: a call that failed stays failed even when nothing explains it,
// so its caller never goes on with an empty result.
func NewStatusError(resp *http.Response, shown error) error {
	if resp == nil {
		if shown == nil {
			return errors.New("the API call failed without a response")
		}
		return shown
	}
	if shown == nil {
		shown = fmt.Errorf("unexpected response from the API (status %d)", resp.StatusCode)
	}
	return &StatusError{
		StatusCode: resp.StatusCode,
		RetryAfter: retryAfter(resp.Header.Get("Retry-After"), time.Now()),
		Err:        shown,
	}
}

// retryAfter parses Retry-After, which is either a number of seconds or an
// HTTP date. A missing, malformed or past value is zero.
func retryAfter(v string, now time.Time) time.Duration {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(v); err == nil {
		if secs < 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	if at, err := http.ParseTime(v); err == nil && at.After(now) {
		return at.Sub(now)
	}
	return 0
}

// Retryable is the retry rule for waits on the v6 API, in the shape
// poll.Options.Retryable takes. It reads the errors the clients already
// return (utils.ErrorPerStatusCodeV4), plus the status the polled reads
// attach with NewStatusError:
//
//   - 429 is retried after Retry-After, or on the next poll when it is absent;
//   - 5xx, a timeout and a request that got no response are retried on the
//     next poll (ErrorPerStatusCodeV4 maps all of them to ErrorInternalServerError
//     or ErrorTimeoutAPICall);
//   - any other error ends the wait.
func Retryable(err error) (bool, time.Duration) {
	var se *StatusError
	if errors.As(err, &se) {
		switch {
		case se.StatusCode == http.StatusTooManyRequests:
			return true, se.RetryAfter
		case se.StatusCode >= 500:
			return true, 0
		}
		return false, 0
	}
	if errors.Is(err, utils.ErrorInternalServerError) || errors.Is(err, utils.ErrorTimeoutAPICall) {
		return true, 0
	}
	return false, 0
}

// BuildError is one build that did not reach ready. A step that waits on
// several builds joins one BuildError per failure with errors.Join, so the
// user sees everything that went wrong at once.
type BuildError struct {
	Kind      string // "application", "connector", ...
	Name      string
	ID        int64
	VersionID string
	// State is the last state seen: the failed one, or the one the build was
	// stuck in when the wait gave up.
	State   State
	TraceID string // when the API returns one
	Err     error
}

func (e *BuildError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %q (id %d, version %s)", e.Kind, e.Name, e.ID, e.VersionID)
	if e.State != "" {
		fmt.Fprintf(&b, " ended in %s", e.State)
	}
	if e.Err != nil {
		fmt.Fprintf(&b, ": %v", e.Err)
	}
	if e.TraceID != "" {
		fmt.Fprintf(&b, " (trace id %s)", e.TraceID)
	}
	return b.String()
}

func (e *BuildError) Unwrap() error { return e.Err }

// ErrorFor is the error a v6 client returns for a failed call: the v4
// handling (utils.ErrorPerStatusCodeV4), except for a 409. v4 reads every 409
// as a name already in use, which is what its create retries rely on, but in
// v6 a 409 also means a state conflict (a release whose resources can no
// longer change, say). So a 409 is a name in use only when the API's message
// says so; otherwise the error is the API's message.
//
// errBody is what utils.LogAndRewindBodyV4 returns, a summary of the body, not
// the body itself; the API's JSON is read from httpResp, which that call
// rewinds.
func ErrorFor(errBody string, httpResp *http.Response, err error) error {
	if httpResp == nil || httpResp.StatusCode != http.StatusConflict {
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}
	detail := apiMessage(rawBody(httpResp))
	if detail == "" || nameTaken(detail) {
		return utils.ErrorNameInUse
	}
	return fmt.Errorf("%w: %s", ErrConflict, detail)
}

// ErrConflict is a 409 that is not a name in use.
var ErrConflict = errors.New("the API refused the change")

// rawBody reads the response body and rewinds it for whoever reads it next.
func rawBody(resp *http.Response) string {
	if resp == nil || resp.Body == nil {
		return ""
	}
	data, err := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewReader(data))
	if err != nil {
		return ""
	}
	return string(data)
}

// apiMessage reads the messages of a v6 error body: each error's detail, and
// its meta.message when that says more.
func apiMessage(body string) string {
	if body == "" {
		return ""
	}
	var out []string
	for _, e := range gjson.Get(body, "errors").Array() {
		detail := e.Get("detail").String()
		if m := e.Get("meta.message").String(); m != "" && !strings.Contains(detail, m) {
			detail = strings.TrimSpace(detail + " " + m)
		}
		if detail != "" {
			out = append(out, detail)
		}
	}
	if len(out) == 0 {
		if d := gjson.Get(body, "detail").String(); d != "" {
			out = append(out, d)
		}
	}
	return strings.Join(out, "; ")
}

func nameTaken(msg string) bool {
	lower := strings.ToLower(msg)
	for _, phrase := range utils.NameTaken {
		if strings.Contains(lower, strings.ToLower(phrase)) {
			return true
		}
	}
	return strings.Contains(lower, "name") && strings.Contains(lower, "in use")
}
