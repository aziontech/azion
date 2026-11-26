package version

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aziontech/azion-cli/utils"
)

// A StatusError never panics when printed, and NewStatusError never wraps nil.
func TestStatusErrorWithoutAnError(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Header: http.Header{}}
	err := NewStatusError(resp, nil)
	if err == nil || err.Error() == "" {
		t.Fatalf("NewStatusError(200, nil) = %v, want an error", err)
	}
	if got := (&StatusError{StatusCode: 502}).Error(); got == "" {
		t.Error("an empty StatusError printed nothing")
	}
	if NewStatusError(nil, nil) == nil {
		t.Error("NewStatusError(nil, nil) = nil")
	}
}

// conflict is a 409 as a client sees it: the body rewound on the response,
// and errBody the summary utils.LogAndRewindBodyV4 returns, not the JSON.
func conflict(body string) (*http.Response, string) {
	resp := &http.Response{StatusCode: http.StatusConflict, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body))}
	summary, err := utils.LogAndRewindBodyV4(resp)
	if err != nil {
		summary = ""
	}
	return resp, summary
}

// A 409 is a name in use only when the API says so; any other conflict keeps
// the API's message instead of being reported as a name in use.
func TestConflictKeepsTheAPIMessage(t *testing.T) {
	resp, body := conflict(`{"errors":[{"code":"43009","title":"Conflict","status":"409","detail":"Validation failed on field 'resources'.","meta":{"field":"resources","message":"The field 'resources' is immutable after creation."}}]}`)
	err := ErrorFor(body, resp, nil)
	if errors.Is(err, utils.ErrorNameInUse) {
		t.Fatalf("a resources conflict was reported as a name in use: %v", err)
	}
	if !errors.Is(err, ErrConflict) || !strings.Contains(err.Error(), "immutable after creation") {
		t.Errorf("ErrorFor(409 resources) = %v, want the API's message", err)
	}

	resp, body = conflict(`{"errors":[{"status":"409","detail":"An application with the name 'site' already exists."}]}`)
	if err := ErrorFor(body, resp, nil); !errors.Is(err, utils.ErrorNameInUse) {
		t.Errorf("ErrorFor(409 name) = %v, want ErrorNameInUse", err)
	}

	resp, body = conflict("")
	if err := ErrorFor(body, resp, nil); !errors.Is(err, utils.ErrorNameInUse) {
		t.Errorf("ErrorFor(409 without a body) = %v, want ErrorNameInUse", err)
	}
}
