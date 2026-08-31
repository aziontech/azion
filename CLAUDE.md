# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## ⚡ Important: v4 vs v3 APIs

**v4 is the PRIMARY and LONG-TERM SUPPORTED API.** v3 is legacy/deprecated.

| Aspect | v4 (Primary) | v3 (Legacy) |
|--------|--------------|-------------|
| API Type | REST | GraphQL |
| SDK | `azionapi-v4-go-sdk-dev` | `azionapi-go-sdk` |
| Commands | `pkg/cmd/` | `pkg/v3commands/` |
| API Clients | `pkg/api/` | `pkg/v3api/` |
| Manifest | `pkg/manifest/` | `pkg/v3manifest/` |
| Status | **Active development** | Maintenance only |

**When adding new features:** Always use v4 paths (`pkg/cmd/`, `pkg/api/`, `pkg/manifest/`).

**v3 commands are loaded only when:** `HasBlockAPIV4Flag()` returns true (user has v3-only flag set).

## Build, Test, and Development Commands

```bash
# Build (production environment)
make build

# Development with hot-reload
make dev

# Run all tests with coverage
make test

# Run specific test file
go test -v -run TestNewCmd ./pkg/cmd/create/applications/

# Run specific test
go test -v -run TestName ./path/to/package/

# Lint code
make lint

# Security scanning
make sec          # GoSec
make govulncheck  # Vulnerability check

# Generate documentation
make docs

# Cross-compile for all platforms
make cross-build
```

## Project Overview

Azion CLI is a Go-based command-line interface for managing applications on the Azion Edge Platform. It supports initializing, building, and deploying applications using various frameworks (Next.js, Vue, Angular, Astro, React, Vite, etc.).

**Go Version:** 1.25.9 (strictly enforced - follows "Release - 1" policy)

## Folder Structure

```
azion/
├── cmd/
│   ├── azion/main.go          # Main CLI entry point
│   └── gen_docs/main.go       # Documentation generator
├── pkg/
│   ├── api/                   # v4 REST API clients (22 resources)
│   ├── cmd/                   # v4 command implementations
│   ├── manifest/              # v4 manifest interpreter
│   ├── v3api/                 # v3 GraphQL API clients (11 resources) [LEGACY]
│   ├── v3commands/            # v3 command implementations [LEGACY]
│   ├── v3manifest/            # v3 manifest interpreter [LEGACY]
│   ├── cmdutil/               # Factory, flags, interfaces
│   ├── config/                # Configuration management
│   ├── contracts/             # Data structures for CLI
│   ├── iostreams/             # Console I/O
│   ├── logger/                # Zap-based logging
│   ├── output/                # Output formatting utilities
│   ├── token/                 # Authentication and token management
│   ├── testutils/             # Test helpers
│   ├── httpmock/              # HTTP mocking for tests
│   ├── github/                # GitHub API integration
│   ├── metric/                # Usage metrics collection
│   ├── schedule/              # Scheduled command execution
│   ├── dry_run/               # Dry-run mode
│   └── utils/                 # Shared utilities
├── messages/                  # User-facing message strings (i18n)
│   ├── build/
│   ├── create/
│   ├── delete/
│   ├── deploy/
│   └── ...                    # One directory per command
├── env/
│   ├── prod                   # Production endpoints
│   ├── stage                  # Staging endpoints
│   └── local                  # Local development
└── scripts/
    ├── install.sh             # Official installer
    ├── e2e.sh                 # End-to-end tests
    └── completions.sh         # Shell completions
```

## Architecture

### Entry Points
- `cmd/azion/main.go` - Main CLI entry point
- `cmd/gen_docs/main.go` - Documentation generator

### Command Structure (Cobra Framework)

The CLI uses the Cobra framework. Commands are organized in `pkg/cmd/`:

```
azion (root)
├── init, build, deploy, dev        # Core workflow
├── create, describe, list          # CRUD operations
├── delete, update                  # CRUD operations
├── config, profiles, login, logout # Configuration
├── link, unlink, clone, sync       # Project management
├── purge, warmup, rollback         # Operations
└── logs, whoami, version           # Utilities
```

### Core Components

**Factory Pattern (`pkg/cmdutil/factory.go`):**
```go
type Factory struct {
    HttpClient *http.Client
    IOStreams  *iostreams.IOStreams
    Config     config.Config
    Flags
}
```

All commands receive a `*cmdutil.Factory` for dependency injection.

**Configuration (`pkg/config/`):**
- Default path: `~/.azion/`
- Files: `settings.toml`, `profiles.json`, `metrics.json`, `schedule.json`
- Supports multiple profiles via `profiles.json`

**Manifest System (`pkg/manifest/`):**
- Defines applications as code via `.edge/manifest.json`
- Resource types: Applications, Functions, CacheSettings, RulesEngine, Connectors, Workloads, Firewalls, Purge
- Resource creation order matters (functions → applications → cache settings → rules)

**Pre-Command Flow (`pkg/cmd/root/pre_command.go`):**
1. Ensure `profiles.json` exists (creates default if missing)
2. Set HTTP timeout from `--timeout` flag (default 50s)
3. Handle `--token` flag (validate and save)
4. Check metrics authorization
5. Check for CLI updates (every 24h)

**Messages (`messages/`):**
- User-facing strings organized by command
- Pattern: `msg "github.com/aziontech/azion-cli/messages/<command>"`
- Keeps UI text separate from logic for i18n

### Test Pattern

Tests use `pkg/testutils/testutils.go`:

```go
func TestExample(t *testing.T) {
    tests := []struct {
        name   string
        args   []string
        mock   func() *httpmock.Registry
        output string
        err    error
    }{
        {
            name: "test case name",
            args: []string{"--flag", "value"},
            mock: func() *httpmock.Registry {
                mock := httpmock.Registry{}
                mock.Register(
                    httpmock.REST("POST", "endpoint"),
                    httpmock.JSONFromFile("./fixtures/response.json"),
                )
                return &mock
            },
            output: "expected output",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f, out, _ := testutils.NewFactory(tt.mock())
            cmd := NewCmd(f)
            cmd.SetArgs(tt.args)
            _, err := cmd.ExecuteC()
            assert.Equal(t, tt.output, out.String())
        })
    }
}
```

Test fixtures go in `./fixtures/` subdirectories next to test files.

## Key Patterns

### Creating a New Command (v4)

1. Create `pkg/cmd/<category>/<name>.go`
2. Define `NewCmd(f *cmdutil.Factory) *cobra.Command`
3. Add flags with `cmd.Flags()`
4. Call API client from `pkg/api/<resource>/`
5. Add messages in `messages/<category>/`
6. Format output with `pkg/output/`
7. Create test file `<name>_test.go` with mock

### API Client Pattern

Each resource has a client in `pkg/api/<resource>/client.go`:

```go
func NewClient(c *http.Client, url string, token string) *Client {
    conf := sdk.NewConfiguration()
    conf.HTTPClient = c
    conf.AddDefaultHeader("Authorization", "token "+token)
    conf.UserAgent = "Azion_CLI/" + version.BinVersion
    return &Client{apiClient: sdk.NewAPIClient(conf)}
}
```

### Output Formatting

Use `pkg/output/` for consistent formatting:

```go
logger.FInfoFlags(f.IOStreams.Out, msg.Success, f.Format, f.Out)
```

Supported formats: text (default), json, yaml.

### Error Handling

Use `utils/errors.go` for standardized error messages mapped to HTTP status codes.

## Environment Variables

The CLI reads from `env/prod`, `env/stage`, `env/local`:

```
STORAGE_URL=
AUTH_URL=
API_URL=
API_V4_URL=
CONSOLE=
TEMPLATE_BRANCH=
TEMPLATE_MAJOR=
```

## Commit Convention

Semantic commits are enforced via git hooks:
- `{type}({scope}): {subject}`
- Types: feat, fix, docs, style, refactor, test, chore
- Install hooks: `git config core.hooksPath hooks`

## Go Code Style

### For Range Loops

Use value iteration (`for _, item := range slice`) as the standard pattern:

```go
// Preferred: access value directly
for _, item := range items {
    if item.Name == target {
        found = &item
        break
    }
}

// Avoid: index-based access (used only when index is needed)
for i := range items {
    items[i].Name = "modified" // OK when modifying slice elements
}
```

**Rationale:** Value iteration is more readable and matches existing codebase patterns (see `pkg/cmd/init/init.go` for examples).