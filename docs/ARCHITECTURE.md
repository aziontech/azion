# Architecture — azion (CLI)

## Overview

Command-line interface for the Azion Edge Platform. Built with Go and Cobra, it enables developers to initialize, build, deploy, and manage edge applications and all Azion products from the terminal. Supports multi-framework templates (Next.js, Vue, Angular, Astro, React, Vite), dual API versions (v3 legacy + v4), and cross-platform distribution.

## Command Structure

```
azion
├── init          Initialize edge app from templates (fetched via API)
├── build         Build application locally
├── dev           Local development server
├── deploy        Deploy to Azion Edge (local or remote mode)
├── deploy-remote Remote browser-based deployment
├── clone         Clone existing edge application
├── link / unlink Link/unlink project to Azion
├── login/logout  Authentication (browser OAuth or token)
├── create        Create resources (19 subcommands: applications, domains, functions, etc.)
├── list          List resources (18 subcommands)
├── describe      Describe resource details (15 subcommands)
├── update        Update resources (20 subcommands)
├── delete        Delete resources (14 subcommands)
├── purge         Purge edge caches
├── logs          View application logs
├── sync          Sync local and remote state
├── rollback      Rollback to previous versions
├── profiles      Manage authentication profiles
├── whoami        Show authenticated user
└── version       Show CLI version
```

## Project Structure

```
cmd/azion/main.go              Entry point — factory setup, root command
pkg/
├── cmd/                        Command implementations (29 packages)
│   ├── root/root.go            Root command with v3/v4 API routing
│   ├── init/                   Template-based project scaffolding
│   ├── deploy/                 Deployment orchestration (S3 + script runner)
│   ├── create/list/update/     CRUD for all Azion products
│   └── ...
├── api/                        API client wrappers (24 product packages)
├── cmdutil/factory.go          Dependency injection factory
├── token/                      Token validation, profile management
├── config/                     Configuration (~/.azion/ settings)
├── manifest/                   Project manifest (azion.json) handling
├── output/                     Output formatting (JSON, table, text)
├── metric/                     Usage telemetry (Segment.io)
├── v3api/                      Legacy v3 API support (13 packages)
├── v3commands/                 v3 command variants (23 packages)
└── constants/                  API URLs, defaults
messages/                       i18n message strings (40+ packages)
scripts/                        Build/release scripts
```

## API Integration

```
CLI command
    │
    ├── v4 SDK (azionapi-v4-go-sdk-dev) ──→ api.azion.com/v4
    │     Applications, Workloads, Functions, Firewall, WAF,
    │     Domains, Storage, Variables, Network Lists
    │
    ├── v3 SDK (azionapi-go-sdk) ──→ api.azionapi.net
    │     Legacy product endpoints
    │
    ├── AWS SDK v2 ──→ S3 (storage uploads, bucket management)
    │
    └── GraphQL (machinebox/graphql) ──→ Analytics queries
```

Authentication: token-based via `Authorization: token [token]` header. Tokens validated against `/user/me`. Multi-profile support stored in `~/.azion/settings.toml`.

## Deploy Flow

```
azion deploy
    ├── 1. Build project (if needed)
    ├── 2. Create S3 bucket + credentials (first deploy)
    ├── 3. Upload static files to S3
    ├── 4. Call script-runner API (remote execution)
    ├── 5. Poll execution logs (2s intervals)
    └── 6. Return deployment URL + timing summary
```

Supports dry-run mode (`--dry-run`) for simulation without upload.

## Dependencies

| Dependency | Purpose |
|-----------|---------|
| spf13/cobra + viper | CLI framework + config |
| azionapi-v4-go-sdk-dev | v4 API SDK |
| azionapi-go-sdk | v3 legacy API SDK |
| go-git/v5 | Git operations (clone templates) |
| aws-sdk-go-v2 | S3 storage operations |
| AlecAivazis/survey/v2 | Interactive prompts |
| go.uber.org/zap | Structured logging |
| fatih/color | Terminal colors |
| segmentio/analytics-go/v3 | Usage telemetry |
| pelletier/go-toml/v2 | Config file parsing |

## Distribution

Cross-compiled via GoReleaser for linux/darwin/freebsd/windows across amd64/arm64/386/arm/ppc64. Distributed through Homebrew, apt, rpm, apk, Chocolatey, WinGet, and direct S3 downloads.

## Configuration

- **User config**: `~/.azion/` (settings.toml, profiles.toml, credentials.toml)
- **Project config**: `azion/` directory (azion.json manifest)
- **Environment**: `AZIONCLI_*` prefix, auto-loaded by Viper
