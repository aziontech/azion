# Runbook — azion (CLI)

## Service Identity

| Field | Value |
|-------|-------|
| Name | azion (CLI) |
| Type | CLI tool |
| Language | Go 1.25 |
| Framework | Cobra + Viper |
| Repository | github.com/aziontech/azion |
| Distribution | Homebrew, apt, rpm, apk, Chocolatey, WinGet, S3 |

## Building

```bash
# Build for current platform
make build

# Cross-compile all platforms
make cross-build

# Development with hot reload
make dev

# Generate CLI reference docs
make docs
```

Build injects version, API URLs, and metrics key via ldflags.

## Testing

```bash
# Unit tests with coverage
make test
# Reports: cover/azioncoverage.html

# Linting
make lint

# Security scan
make sec

# Vulnerability check
make govulncheck
```

## CI/CD

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| test_and_build.yml | PR | Tests, gosec, GoReleaser validation |
| deploy_prod.yml | Push to main | Tag, build all platforms, release (S3 + GoReleaser + Chocolatey + WinGet) |
| deploy_stage.yml | Manual | Stage environment release |
| ci-compliance.yml | PR, weekly | Azion compliance checks |
| ci-security.yml | PR, weekly | Security scanning |
| e2e_test.yml | Manual | End-to-end integration tests |
| generate_docs.yml | Manual | Auto-generate CLI docs |

## Release Process

1. Merge PR to `main`
2. `deploy_prod.yml` auto-triggers:
   - Bumps git tag
   - GoReleaser builds binaries for 11 platform/arch combinations
   - Uploads to `azion-downloads` S3 bucket
   - Creates GitHub release
   - Publishes Chocolatey package
   - Updates WinGet manifest

## Common Issues

### 1. Private Module Access

**Symptoms**: `go mod download` fails for `aziontech/*` packages.

**Resolution**: Set `GOPRIVATE=github.com/aziontech/*` and configure git with a valid GitHub token: `git config --global url."https://${TOKEN}@github.com/".insteadOf "https://github.com/"`.

### 2. GoReleaser Validation Fails

**Symptoms**: CI fails at GoReleaser check step.

**Resolution**: Run `goreleaser check --config .goreleaser.yaml` locally. Common causes: invalid YAML syntax, missing ldflags variables, or architecture matrix issues.

### 3. API Version Conflicts

**Symptoms**: Commands fail with unexpected responses.

**Resolution**: The CLI supports both v3 and v4 APIs. If a resource only exists in v3, ensure the correct API flag or command variant is used. V3 commands are under `pkg/v3commands/`.

### 4. Template Init Fails

**Symptoms**: `azion init` fails to fetch templates or clone samples.

**Resolution**: Check network access to `api.azion.com/v4/utils/project_samples` and `github.com/aziontech/azion-samples.git`. Ensure the authentication token is valid (`azion login`).

### 5. Deploy S3 Upload Fails

**Symptoms**: Deploy fails at file upload step.

**Resolution**: S3 credentials are auto-created on first deploy and cached in `~/.azion/credentials.toml` with 1-year expiry. If expired, delete credentials file and redeploy.

## Monitoring

- **Telemetry**: Anonymous usage metrics via Segment.io (user opt-in)
- **Logging**: Structured JSON via Zap, configurable with `--log-level` flag
- **Timing**: Deploy command prints per-phase timing summary

## Escalation

| Level | Contact | When |
|-------|---------|------|
| L1 | Team Dev Tools & Integrations | CLI bugs, command issues |
| L2 | Team API / Platform | API integration issues, SDK bugs |
