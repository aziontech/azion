package manifest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	logger.New(zapcore.DebugLevel)
	os.Exit(m.Run())
}

// ---------------------------------------------------------------------------
// Manifest interpreter
// ---------------------------------------------------------------------------

func TestManifestInterpreter_ManifestPath(t *testing.T) {
	chdir(t, t.TempDir())

	// Read the directory back rather than reusing t.TempDir(): on macOS the
	// temp dir is reached through a symlink that Getwd resolves.
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	path, err := NewManifestInterpreter().ManifestPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := workingDir + manifestFilePath
	if path != expected {
		t.Errorf("ManifestPath() = %q, want %q", path, expected)
	}
}

func TestManifestInterpreter_ReadManifest(t *testing.T) {
	manifestPath := filepath.Join(t.TempDir(), "manifest.json")
	content := `{
		"functions": [{"name": "myfunc", "path": "worker.js"}],
		"applications": [{
			"name": "myapp",
			"rules": [],
			"cache_settings": [],
			"functions_instances": [{"name": "myinst", "function": "myfunc", "active": true}]
		}]
	}`
	if err := os.WriteFile(manifestPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	f, _, _ := testutils.NewFactory(nil)
	msgs := []string{}

	manifest, err := NewManifestInterpreter().ReadManifest(manifestPath, f, &msgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifest.Functions) != 1 || manifest.Functions[0].Name != "myfunc" {
		t.Errorf("functions = %v, want one entry named myfunc", manifest.Functions)
	}
	if len(manifest.Applications) != 1 {
		t.Fatalf("applications = %v, want one entry", manifest.Applications)
	}
	instances := manifest.Applications[0].FunctionsInstances
	if len(instances) != 1 || instances[0].Function.Name != "myfunc" {
		t.Errorf("function instances = %v, want one referencing myfunc", instances)
	}
}

func TestManifestInterpreter_ReadManifestMissingFile(t *testing.T) {
	f, _, _ := testutils.NewFactory(nil)
	msgs := []string{}

	if _, err := NewManifestInterpreter().ReadManifest("does-not-exist.json", f, &msgs); err == nil {
		t.Error("expected an error for a missing manifest, got nil")
	}
}

// ---------------------------------------------------------------------------
// Args file resolution
// ---------------------------------------------------------------------------

func TestArgsFilePath_HonorsConfigDir(t *testing.T) {
	tests := []struct {
		projectConf string
		expected    string
	}{
		{projectConf: "azion", expected: filepath.Join("azion", "args.json")},
		{projectConf: "myconf", expected: filepath.Join("myconf", "args.json")},
		{projectConf: "nested/conf", expected: filepath.Join("nested", "conf", "args.json")},
	}

	for _, tt := range tests {
		rc := &ResourceContext{ProjectConf: tt.projectConf}
		if got := rc.argsFilePath(); got != tt.expected {
			t.Errorf("argsFilePath() with config-dir %q = %q, want %q", tt.projectConf, got, tt.expected)
		}
	}
}

func TestUnmarshalJsonArgs_ReadsArgsFromConfigDir(t *testing.T) {
	writeArgsFile(t, "myconf", `{"key":"value"}`)

	// The path recorded in azion.json is relative to the project root
	args, found, err := unmarshalJsonArgs(filepath.Join("myconf", "args.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Error("expected the args file to be reported as found")
	}
	if !reflect.DeepEqual(args, map[string]interface{}{"key": "value"}) {
		t.Errorf("got %v, want map[key:value]", args)
	}
}

func TestUnmarshalJsonArgs_AcceptsAbsolutePath(t *testing.T) {
	argsPath := filepath.Join(t.TempDir(), "args.json")
	if err := os.WriteFile(argsPath, []byte(`{"key":"value"}`), 0644); err != nil {
		t.Fatal(err)
	}
	chdir(t, t.TempDir())

	args, found, err := unmarshalJsonArgs(argsPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Error("expected the args file to be reported as found")
	}
	if !reflect.DeepEqual(args, map[string]interface{}{"key": "value"}) {
		t.Errorf("got %v, want map[key:value]", args)
	}
}

func TestUnmarshalJsonArgs_MissingFileIsEmpty(t *testing.T) {
	chdir(t, t.TempDir())

	args, found, err := unmarshalJsonArgs(filepath.Join("azion", "args.json"))
	if err != nil {
		t.Fatalf("missing args file should not error, got: %v", err)
	}
	if found {
		t.Error("a missing args file must not be reported as found")
	}
	if len(args) != 0 {
		t.Errorf("expected empty args, got %v", args)
	}
}

func TestUnmarshalJsonArgs_InvalidJsonErrors(t *testing.T) {
	workingDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(workingDir, "args.json"), []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}
	chdir(t, workingDir)

	if _, _, err := unmarshalJsonArgs("args.json"); err == nil {
		t.Error("expected an error for malformed args.json, got nil")
	}
}

// ---------------------------------------------------------------------------
// Function instances
// ---------------------------------------------------------------------------

func TestApplyFunctionInstances_ArgsFromConfigDir(t *testing.T) {
	argsPath := writeArgsFile(t, "myconf", `{"token":"abc"}`)
	expected := map[string]interface{}{"token": "abc"}

	tests := []struct {
		name      string
		reference contracts.FunctionReference
		confArgs  string
	}{
		{
			name:      "function referenced by name",
			reference: contracts.FunctionReference{Name: "myfunc"},
			confArgs:  argsPath,
		},
		{
			name:      "function referenced by a known id",
			reference: contracts.FunctionReference{ID: 42},
			confArgs:  argsPath,
		},
		{
			// Not described in azion.json, so it falls back to the project args file
			name:      "function referenced by an unknown id",
			reference: contracts.FunctionReference{ID: 12345},
			confArgs:  argsPath,
		},
		{
			// Entry written before args paths were recorded, or by hand
			name:      "config entry without an args path",
			reference: contracts.FunctionReference{Name: "myfunc"},
			confArgs:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := applyInstance(t, confWithFunction(tt.confArgs), "myconf",
				contracts.FunctionInstance{Name: "myinst", Function: tt.reference})

			if !reflect.DeepEqual(payload["args"], expected) {
				t.Errorf("args sent = %v, want %v", payload["args"], expected)
			}
		})
	}
}

func TestApplyFunctionInstances_ManifestArgsWin(t *testing.T) {
	argsPath := writeArgsFile(t, "myconf", `{"token":"from-file"}`)

	payload, _ := applyInstance(t, confWithFunction(argsPath), "myconf", contracts.FunctionInstance{
		Name:     "myinst",
		Function: contracts.FunctionReference{Name: "myfunc"},
		Args:     map[string]interface{}{"token": "from-manifest"},
	})

	expected := map[string]interface{}{"token": "from-manifest"}
	if !reflect.DeepEqual(payload["args"], expected) {
		t.Errorf("args sent = %v, want %v", payload["args"], expected)
	}
}

// Instances are updated with PATCH, so a project with no args file must omit
// the field rather than send an empty object that clears existing arguments.
func TestApplyFunctionInstances_NoArgsFileOmitsField(t *testing.T) {
	chdir(t, t.TempDir())

	payload, _ := applyInstance(t, confWithFunction(""), "myconf", contracts.FunctionInstance{
		Name: "myinst", Function: contracts.FunctionReference{Name: "myfunc"},
	})

	if _, present := payload["args"]; present {
		t.Errorf("args must be omitted when no args file exists, got %v", payload["args"])
	}
}

func TestApplyFunctionInstances_RecordsArgsPathInConfig(t *testing.T) {
	writeArgsFile(t, "myconf", `{"token":"abc"}`)

	_, written := applyInstance(t, confWithFunction(""), "myconf", contracts.FunctionInstance{
		Name: "myinst", Function: contracts.FunctionReference{ID: 12345},
	})

	if written == nil {
		t.Fatal("expected azion.json to be written")
	}
	expected := filepath.Join("myconf", "args.json")
	for _, fn := range written.Function {
		if fn.Args != expected {
			t.Errorf("function %q recorded args path %q, want %q", fn.Name, fn.Args, expected)
		}
	}
}

func TestApplyFunctionInstances_RequiresApplicationID(t *testing.T) {
	f, _, _ := testutils.NewFactory(&httpmock.Registry{})
	msgs := []string{}
	rc := NewResourceContext(f, &contracts.AzionApplicationOptions{}, &contracts.ManifestV4{},
		"azion", &msgs, func(*contracts.AzionApplicationOptions, string) error { return nil })

	err := rc.ApplyFunctionInstances([]contracts.FunctionInstance{{Name: "myinst"}})
	if err == nil {
		t.Error("expected an error when the application id is missing, got nil")
	}
}

// ---------------------------------------------------------------------------
// Resource deletion
// ---------------------------------------------------------------------------

func TestDeleteResources_SkipDeletionAbsent(t *testing.T) {
	// Empty global maps so no deletions are attempted
	CacheIds = map[string]int64{}
	RuleIds = map[string]contracts.RuleIdsStruct{}

	f, _, _ := testutils.NewFactory(nil)
	msgs := []string{}

	conf := &contracts.AzionApplicationOptions{
		Application: contracts.AzionJsonDataApplication{ID: 123},
		// SkipDeletion is intentionally left as nil to simulate absence in JSON
	}

	if err := deleteResources(context.Background(), f, conf, &msgs); err != nil {
		t.Fatalf("deleteResources failed with SkipDeletion absent: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// chdir switches to dir for the duration of the test. Config paths recorded in
// azion.json are relative to the project root, so tests need a known root.
func chdir(t *testing.T, dir string) {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previous)
	})
}

// writeArgsFile creates confDir/args.json inside a fresh working directory,
// switches to it, and returns the path as azion.json would record it.
func writeArgsFile(t *testing.T, confDir, content string) string {
	t.Helper()
	workingDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(workingDir, confDir), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workingDir, confDir, "args.json"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	chdir(t, workingDir)
	return filepath.Join(confDir, "args.json")
}

func confWithFunction(argsPath string) *contracts.AzionApplicationOptions {
	return &contracts.AzionApplicationOptions{
		Application: contracts.AzionJsonDataApplication{ID: 99},
		Function: []contracts.AzionJsonDataFunction{
			{ID: 42, Name: "myfunc", Args: argsPath},
		},
	}
}

const functionInstanceResponse = `{"state":"executed","data":{"id":777,"name":"myinst","function":42,"active":true,"args":{}}}`

// applyInstance runs ApplyFunctionInstances against a mocked API and returns
// both the payload sent to the API and the config that would be persisted.
func applyInstance(
	t *testing.T,
	conf *contracts.AzionApplicationOptions,
	projectConf string,
	inst contracts.FunctionInstance,
) (map[string]interface{}, *contracts.AzionApplicationOptions) {
	t.Helper()

	var payload map[string]interface{}
	mock := &httpmock.Registry{}
	mock.Register(
		httpmock.REST("POST", "workspace/applications/99/functions"),
		func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(body, &payload); err != nil {
				t.Fatal(err)
			}
			return httpmock.JSONFromString(functionInstanceResponse)(req)
		},
	)

	f, _, _ := testutils.NewFactory(mock)
	msgs := []string{}
	var written *contracts.AzionApplicationOptions
	rc := NewResourceContext(f, conf, &contracts.ManifestV4{}, projectConf, &msgs,
		func(c *contracts.AzionApplicationOptions, _ string) error {
			written = c
			return nil
		})

	if err := rc.ApplyFunctionInstances([]contracts.FunctionInstance{inst}); err != nil {
		t.Fatalf("ApplyFunctionInstances: %v", err)
	}
	return payload, written
}
