// Package v6test is a fake API v6 for tests. Unlike pkg/httpmock, whose stubs
// answer once each, a route here answers with a sequence of responses and
// repeats the last one — which is what a version that goes queued, building,
// ready looks like to a client that polls it.
package v6test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// Response is one answer of a route.
type Response struct {
	Status int
	Body   any // marshalled to JSON; a string is sent as is
	Header map[string]string
}

// JSON is a 200 with body.
func JSON(body any) Response { return Response{Status: http.StatusOK, Body: body} }

// Status is an answer with only a status and an optional body.
func Status(code int, body ...any) Response {
	r := Response{Status: code}
	if len(body) > 0 {
		r.Body = body[0]
	}
	return r
}

// Call is one request the fake received.
type Call struct {
	Method string
	Path   string
	Body   map[string]any
}

// API is the fake. Its zero value answers 404 to everything.
type API struct {
	mu     sync.Mutex
	routes map[string]*route
	Calls  []Call
}

type route struct {
	responses []Response
	served    int
}

// On sets the answers for method and path (without the /v6 prefix or host;
// for example "GET /workspace/applications/1/versions/V1").
func (a *API) On(method, path string, responses ...Response) *API {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.routes == nil {
		a.routes = map[string]*route{}
	}
	a.routes[method+" "+path] = &route{responses: responses}
	return a
}

// Client is an *http.Client whose requests the fake answers.
func (a *API) Client() *http.Client { return &http.Client{Transport: a} }

// URL is the base URL to give the clients.
const URL = "http://api.v6.test"

func (a *API) RoundTrip(req *http.Request) (*http.Response, error) {
	var body map[string]any
	if req.Body != nil {
		data, _ := io.ReadAll(req.Body)
		_ = json.Unmarshal(data, &body)
	}

	a.mu.Lock()
	key := req.Method + " " + req.URL.Path
	a.Calls = append(a.Calls, Call{Method: req.Method, Path: req.URL.Path, Body: body})
	rt, ok := a.routes[key]
	var res Response
	if ok && len(rt.responses) > 0 {
		res = rt.responses[min(rt.served, len(rt.responses)-1)]
		rt.served++
	} else {
		res = Status(http.StatusNotFound, map[string]any{"detail": "no route for " + key})
	}
	a.mu.Unlock()

	var data []byte
	switch b := res.Body.(type) {
	case nil:
	case string:
		data = []byte(b)
	default:
		data, _ = json.Marshal(b)
	}
	status := res.Status
	if status == 0 {
		status = http.StatusOK
	}
	h := http.Header{"Content-Type": []string{"application/json"}}
	for k, v := range res.Header {
		h.Set(k, v)
	}
	return &http.Response{StatusCode: status, Header: h, Body: io.NopCloser(bytes.NewReader(data)), Request: req}, nil
}

// Count is how many times method and path were called.
func (a *API) Count(method, path string) int {
	a.mu.Lock()
	defer a.mu.Unlock()
	n := 0
	for _, c := range a.Calls {
		if c.Method == method && c.Path == path {
			n++
		}
	}
	return n
}

// Bodies returns the request bodies sent to method and path, in order.
func (a *API) Bodies(method, path string) []map[string]any {
	a.mu.Lock()
	defer a.mu.Unlock()
	var out []map[string]any
	for _, c := range a.Calls {
		if c.Method == method && c.Path == path {
			out = append(out, c.Body)
		}
	}
	return out
}

// Summary lists the calls, for failure messages.
func (a *API) Summary() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	var b strings.Builder
	for _, c := range a.Calls {
		fmt.Fprintf(&b, "%s %s\n", c.Method, c.Path)
	}
	return b.String()
}

// Page wraps results the way the v6 list endpoints do.
func Page(results ...any) map[string]any {
	if results == nil {
		results = []any{}
	}
	return map[string]any{"count": len(results), "results": results}
}

// Data wraps an object the way the v6 create and retrieve endpoints do.
func Data(obj any) map[string]any { return map[string]any{"data": obj} }

// Version is a versioned object with the fields the deploy reads.
func Version(id int64, name, versionID, state string) map[string]any {
	return map[string]any{"id": id, "name": name, "version_id": versionID, "version_state": state}
}
