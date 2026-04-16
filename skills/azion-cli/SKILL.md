---
name: azion-cli
description: |
  Azion CLI command reference and workflow patterns. Use this skill whenever the user:
  - Asks how to use azion CLI commands (build, deploy, init, create, etc.)
  - Wants to know CLI flags, options, and syntax
  - Needs help setting up Azion projects or applications
  - Wants to configure CI/CD with Azion CLI
  - Has questions about CLI troubleshooting
  - Wants to generate azion.json or args.json configurations
  - Mentions "azion" in the context of commands, deployment, or configuration

  Trigger even if the user doesn't explicitly say "CLI" - they might ask "how do I deploy to Azion" or "how do I create an edge application."
---

# Azion CLI Skill

Comprehensive reference for Azion CLI v4.19.1 commands, workflows, and patterns.

## Overview

Azion CLI is a command-line interface for managing Azion edge applications, functions, and resources. Install with:

```bash
curl -fsSL https://cli.azion.app/install.sh | bash
# Or via package managers:
brew install azion
choco install azion
winget install aziontech.azion
```

## Quick Command Reference

### Project Lifecycle
| Command | Purpose |
|---------|---------|
| `azion init` | Initialize new project from template |
| `azion build` | Build application locally |
| `azion deploy` | Deploy to Azion edge |
| `azion link` | Link existing project to Azion |
| `azion unlink` | Unlink project from Azion |
| `azion dev` | Start local development server |

### Resource Management
| Command | Purpose |
|---------|---------|
| `azion create <resource>` | Create new resources |
| `azion update <resource>` | Update existing resources |
| `azion delete <resource>` | Delete resources |
| `azion describe <resource>` | View resource details |
| `azion list <resource>` | List all resources |

### Configuration
| Command | Purpose |
|---------|---------|
| `azion config init` | Initialize azion.json config file |
| `azion config apply` | Apply configuration to platform |
| `azion config delete` | Delete configuration |
| `azion sync` | Sync local config with remote |

### Authentication & Profiles
| Command | Purpose |
|---------|---------|
| `azion login` | Login to Azion account |
| `azion logout` | Logout from account |
| `azion whoami` | Display current user |
| `azion profiles` | Manage multiple profiles |

### Operations
| Command | Purpose |
|---------|---------|
| `azion logs cells` | View Edge Functions logs |
| `azion logs http` | View HTTP event logs |
| `azion purge` | Clear cache entries |
| `azion warmup` | Preload URLs into edge cache |
| `azion rollback` | Revert to previous deployment |
| `azion clone` | Clone applications |

## Common Workflows

### New Project Setup
```bash
# Login first
azion login

# Initialize from template
azion init --name "my-project"

# Build and deploy
azion build
azion deploy
```

### Deploy Existing Project
```bash
# Link existing codebase
azion link --name "my-app" --preset nextjs

# Deploy
azion deploy
```

### Configuration as Code
```bash
# Create config file
azion config init --config-dir ./config

# Edit azion.json as needed, then apply
azion config apply --config-dir ./config
```

### Debugging Applications
```bash
# View function logs in real-time
azion logs cells --tail --function-id 1234

# View HTTP logs
azion logs http --tail --limit 50

# Check application details
azion describe application --id 1234
```

### Cache Management
```bash
# Purge by URL
azion purge --urls "example.com/path1,example.com/path2"

# Purge by wildcard
azion purge --wildcard "example.com/*"

# Warmup cache
azion warmup --url "https://example.com" --max-urls 500
```

## Global Flags

Available on ALL commands:

| Flag | Description |
|------|-------------|
| `-c, --config <path>` | Config folder for current command |
| `-d, --debug` | Debug level logging |
| `--format <type>` | Output format (json) |
| `-h, --help` | Show help |
| `-l, --log-level <level>` | Log level: debug/info/error |
| `--no-color` | Disable colors |
| `--out <file>` | Export output to file |
| `-s, --silent` | Silence all logs |
| `--timeout <sec>` | HTTP timeout (default 50) |
| `-t, --token <token>` | Personal token for auth |
| `-y, --yes` | Auto-answer yes to prompts |

## Framework Presets

Supported frameworks for `init` and `link`:

```
next           Next.js
react          React
vue            Vue.js
angular        Angular
astro          Astro
hexo           Hexo
static         Static HTML
javascript     Vanilla JavaScript
typescript     Vanilla TypeScript
```

## Resource Types for Create/Update/Delete

```
application         Edge Application
cache-setting       Cache Settings
connector           Serverless Connector
firewall            Edge Firewall
firewall-instance   Firewall Function Instance
firewall-rule       Firewall Rule
function            Edge Function
function-instance   Function Instance
network-list        Network List
origin              Origin
personal-token      Personal Token
profile             CLI Profile
rules-engine        Rules Engine Rule
storage             Object Storage
variables           Environment Variables
waf                 Web Application Firewall
waf-exceptions      WAF Exceptions
workload            Edge Workload
workload-deployment Workload Deployment
```

## CI/CD Integration Pattern

```yaml
# .github/workflows/deploy.yml
name: Deploy to Azion

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Azion CLI
        run: curl -fsSL https://cli.azion.app/install.sh | bash

      - name: Build
        run: azion build
        env:
          AZION_TOKEN: ${{ secrets.AZION_TOKEN }}

      - name: Deploy
        run: azion deploy --auto
        env:
          AZION_TOKEN: ${{ secrets.AZION_TOKEN }}
```

## Troubleshooting

### Authentication Issues
```bash
# Re-login if token expired
azion logout
azion login --username your@email.com --password "password"

# Or use token directly
azion --token "your-personal-token" deploy
```

### Build Failures
```bash
# Skip framework build if issues
azion build --skip-framework-build

# Use custom entrypoint
azion build --entry ./src/index.js

# Enable debug logging
azion build --debug
```

### Deployment Issues
```bash
# Dry run to test without deploying
azion deploy --dry-run

# Sync local config with remote
azion deploy --sync

# Check logs
azion logs cells --tail
```

## Reference Files

For complete command documentation, see:
- `references/commands.md` - Full command reference with all flags
- `references/workflows.md` - Step-by-step workflow guides