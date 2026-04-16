# Azion CLI Commands Reference

Complete reference for all Azion CLI v4.19.1 commands.

## Table of Contents

- [Project Commands](#project-commands)
  - [init](#init)
  - [build](#build)
  - [deploy](#deploy)
  - [dev](#dev)
  - [link](#link)
  - [unlink](#unlink)
- [Resource Commands](#resource-commands)
  - [create](#create)
  - [update](#update)
  - [delete](#delete)
  - [describe](#describe)
  - [list](#list)
  - [clone](#clone)
- [Configuration Commands](#configuration-commands)
  - [config init](#config-init)
  - [config apply](#config-apply)
  - [config delete](#config-delete)
  - [sync](#sync)
- [Authentication Commands](#authentication-commands)
  - [login](#login)
  - [logout](#logout)
  - [whoami](#whoami)
  - [profiles](#profiles)
- [Operations Commands](#operations-commands)
  - [logs cells](#logs-cells)
  - [logs http](#logs-http)
  - [purge](#purge)
  - [warmup](#warmup)
  - [rollback](#rollback)
- [Global Flags](#global-flags)

---

## Project Commands

### init

Initialize a new project from a starter template.

```bash
azion init [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--auto` | Run without interruptions |
| `--local` | Run build and deploy locally |
| `--name <string>` | Application name |
| `--package-manager <string>` | Package manager (npm/yarn/pnpm) |
| `--skip-framework-build` | Skip framework build phase |
| `--sync` | Sync local azion.json with remote |
| `--template <string>` | Use a specific starter template |
| `-h, --help` | Show help |

**Examples:**
```bash
azion init
azion init --name "my-project"
azion init --name "my-project" --preset nextjs --auto
azion init --name "my-project" --template nextjs-starter
```

**Template Examples:**
```bash
# Use a Next.js starter template
azion init --template nextjs-starter --name "my-nextjs-app"

# Use a React starter template
azion init --template react-starter --name "my-react-app"

# Combine with auto mode for CI/CD
azion init --template vue-starter --name "my-vue-app" --auto
```

---

### build

Build application locally.

```bash
azion build [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--config-dir <string>` | Custom config directory (default "azion") |
| `--entry <string>` | Code entrypoint (default "./main.js") |
| `--preset <string>` | Application preset |
| `--skip-framework-build` | Skip framework build phase |
| `--use-node-polyfills <string>` | Use node polyfills in build |
| `--use-own-worker <string>` | Use custom worker expression |
| `-h, --help` | Show help |

**Examples:**
```bash
azion build
azion build --entry ./src/index.js
azion build --preset nextjs --skip-framework-build
```

---

### deploy

Deploy application to Azion edge.

```bash
azion deploy [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--auto` | Run without interruptions |
| `--config-dir <string>` | Custom config directory (default "azion") |
| `--dry-run` | Simulate deploy without actions |
| `--env <string>` | Custom .env file path (default ".edge/.env") |
| `--local` | Build and deploy locally |
| `--no-prompt` | Return errors instead of prompts |
| `--path <string>` | Path to static files |
| `--skip-build` | Skip build step |
| `--skip-framework-build` | Skip framework build phase |
| `--sync` | Sync local azion.json with remote |
| `--writable-bucket` | Create bucket with read-write access |
| `-h, --help` | Show help |

**Examples:**
```bash
azion deploy
azion deploy --auto
azion deploy --path ./dist --skip-build
azion deploy --dry-run  # Test without deploying
```

---

### dev

Start local development server.

```bash
azion dev [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--port <int>` | Localhost port |
| `--skip-framework-build` | Skip framework build phase |
| `-h, --help` | Show help |

**Examples:**
```bash
azion dev
azion dev --port 3000
```

---

### link

Link existing project to Azion application.

```bash
azion link [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--auto` | Run without interruptions |
| `--config-dir <string>` | Custom config directory (default "azion") |
| `--local` | Run build and deploy locally |
| `--name <string>` | Application name |
| `--package-manager <string>` | Package manager (npm/yarn/pnpm) |
| `--preset <string>` | Application template |
| `--remote <string>` | Clone remote repository |
| `--skip-framework-build` | Skip framework build phase |
| `--sync` | Sync local azion.json with remote |
| `-h, --help` | Show help |

**Examples:**
```bash
azion link
azion link --name "my-app" --preset react
azion link --remote https://github.com/user/repo
```

---

### unlink

Unlink project from Azion.

```bash
azion unlink [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `-h, --help` | Show help |

---

## Resource Commands

### create

Create a new resource.

```bash
azion create <resource> [flags]
```

**Available Resources:**
| Resource | Description |
|----------|-------------|
| `application` | Create Edge Application |
| `cache-setting` | Create Cache Settings |
| `connector` | Create Connector |
| `firewall` | Create Firewall |
| `firewall-instance` | Create Firewall Function Instance |
| `firewall-rule` | Create Firewall Rule |
| `function` | Create Edge Function |
| `function-instance` | Create Function Instance |
| `network-list` | Create Network List |
| `origin` | Create Origin |
| `personal-token` | Create Personal Token |
| `profile` | Create CLI Profile |
| `rules-engine` | Create Rules Engine Rule |
| `storage` | Create Storage bucket/object |
| `variables` | Create Environment Variable |
| `waf` | Create WAF |
| `waf-exceptions` | Create WAF Exception |
| `workload` | Create Workload |
| `workload-deployment` | Create Workload Deployment |

**Examples:**
```bash
azion create application --name "my-app"
azion create function --name "my-function" --code ./function.js --active true
azion create variables --key "API_KEY" --value "secret" --secret true
azion create storage bucket --name "my-bucket" --edge-access read_only
```

---

### update

Update an existing resource.

```bash
azion update <resource> [flags]
```

**Available Resources:** Same as create (except `profile`, `workload-deployment`).

**Examples:**
```bash
azion update application --id 1234 --active true
azion update function --id 5678 --code ./updated.js
```

---

### delete

Delete a resource.

```bash
azion delete <resource> [flags]
```

**Available Resources:** Same as create.

**Examples:**
```bash
azion delete application --id 1234
azion delete function --id 5678
```

---

### describe

Display resource details.

```bash
azion describe <resource> [flags]
```

**Available Resources:** application, cache-setting, connector, firewall, firewall-instance, firewall-rule, function, function-instance, network-list, origin, personal-token, rules-engine, storage, variables, waf, waf-exceptions, workload.

**Examples:**
```bash
azion describe application --id 1234
azion describe function --id 5678
```

---

### list

List all resources of a type.

```bash
azion list <resource> [flags]
```

**Available Resources:** application, cache-setting, connector, firewall, firewall-instance, firewall-rule, function, function-instance, network-list, origin, personal-token, rules-engine, storage, variables, waf, waf-exceptions, workload, workload-deployment.

**Examples:**
```bash
azion list application
azion list function --format json
azion list workload
```

---

### clone

Clone an application.

```bash
azion clone application [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--application-id <int>` | ID of application to clone |
| `--name <string>` | Name for new application |
| `-h, --help` | Show help |

**Example:**
```bash
azion clone application --application-id 1234 --name "my-app-copy"
```

---

## Configuration Commands

### config init

Create azion.json configuration file.

```bash
azion config init [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--config-dir <string>` | Directory for config file (default ".") |
| `--force` | Overwrite existing config |
| `-h, --help` | Show help |

**Examples:**
```bash
azion config init
azion config init --config-dir ./my-project --force
```

---

### config apply

Apply configuration to Azion Platform.

```bash
azion config apply [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--config-dir <string>` | Directory containing config (default ".") |
| `-h, --help` | Show help |

**Examples:**
```bash
azion config apply
azion config apply --config-dir ./config
```

---

### config delete

Delete all resources from azion.json.

```bash
azion config delete [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--force` | Force deletion without confirmation |
| `-h, --help` | Show help |

---

### sync

Synchronize local config with remote resources.

```bash
azion sync [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--config-dir <string>` | Custom config directory (default "azion") |
| `--env <string>` | Custom .env file path (default ".edge/.env") |
| `--extension <string>` | Config file extension (mjs/cjs/ts/js, default "mjs") |
| `--iac` | Generate azion.config file |
| `-h, --help` | Show help |

**Examples:**
```bash
azion sync
azion sync --iac --extension ts
```

---

## Authentication Commands

### login

Login to Azion account.

```bash
azion login [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--password <string>` | Password |
| `--username <string>` | Email address |
| `-h, --help` | Show help |

**Examples:**
```bash
azion login
azion login --username user@email.com --password "password"
```

---

### logout

Logout from Azion account.

```bash
azion logout [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `-h, --help` | Show help |

---

### whoami

Display current logged-in user.

```bash
azion whoami [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `-h, --help` | Show help |

---

### profiles

Manage CLI profiles.

```bash
azion profiles [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `-h, --help` | Show help |

**Related commands:**
```bash
azion create profile --name "profile-name"
azion delete profile --name "profile-name"
```

---

## Operations Commands

### logs cells

View Edge Functions console logs.

```bash
azion logs cells [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--function-id <string>` | Filter by function ID |
| `--limit <string>` | Logs per request (default "100") |
| `--pretty` | Prettified output |
| `--tail` | Continuous log stream |
| `-h, --help` | Show help |

**Examples:**
```bash
azion logs cells
azion logs cells --tail
azion logs cells --function-id 1234 --limit 50 --pretty
```

---

### logs http

View HTTP event logs.

```bash
azion logs http [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--limit <string>` | Logs per request (default "100") |
| `--pretty` | Prettified output |
| `--tail` | Continuous log stream |
| `-h, --help` | Show help |

**Examples:**
```bash
azion logs http
azion logs http --tail --limit 200 --pretty
```

---

### purge

Delete objects from cache.

```bash
azion purge [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--cachekey <string>` | URLs with cache keys to purge |
| `--layer <string>` | Cache layer: "cache" or "tiered_cache" (default "cache") |
| `--urls <string>` | Comma-separated URLs to purge |
| `--wildcard <string>` | Wildcard URL pattern |
| `-h, --help` | Show help |

**Examples:**
```bash
azion purge --urls "example.com/page1,example.com/page2"
azion purge --wildcard "example.com/*"
azion purge --cachekey "example.com/@@cookie=value"
azion purge --layer tiered_cache --urls "example.com/*"
```

---

### warmup

Preload URLs into edge cache.

```bash
azion warmup [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--max-concurrent <int>` | Max concurrent requests (default 2) |
| `--max-urls <int>` | Max URLs to process (default 1500) |
| `--timeout <int>` | Request timeout in ms (default 8000) |
| `--url <string>` | Base URL to warm up |
| `-h, --help` | Show help |

**Examples:**
```bash
azion warmup --url "https://example.com"
azion warmup --url "https://example.com/products" --max-urls 500 --max-concurrent 5
```

---

### rollback

Revert to previous deployment.

```bash
azion rollback [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--config-dir <string>` | Custom config directory (default "azion") |
| `--connector-id <int>` | Storage connector ID |
| `-h, --help` | Show help |

**Example:**
```bash
azion rollback --connector-id aaaa-bbbb-cccc-dddd
```

---

## Global Flags

Available on ALL commands:

| Flag | Description |
|------|-------------|
| `-c, --config <string>` | Config folder for current command |
| `-d, --debug` | Debug level logging |
| `--format <string>` | Output format (json) |
| `-h, --help` | Show help |
| `-l, --log-level <string>` | Log level: debug/info/error (default "info") |
| `--no-color` | Disable colored output |
| `--out <string>` | Export output to file |
| `-s, --silent` | Silence all logs |
| `--timeout <int>` | HTTP timeout in seconds (default 50) |
| `-t, --token <string>` | Personal token for auth |
| `-v, --version` | Show CLI version |
| `-y, --yes` | Auto-answer yes to prompts |

---

## Version

Display CLI version:

```bash
azion version
azion -v
azion --version
```