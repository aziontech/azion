---
name: azion
description: |
  Azion CLI command reference and workflow patterns. Use this skill whenever the user:
  - Asks how to use azion CLI commands (build, deploy, init, create, etc.)
  - Wants to know CLI flags, options, and syntax
  - Needs help setting up Azion projects or applications
  - Wants to configure CI/CD with Azion CLI
  - Has questions about CLI troubleshooting
  - Wants to generate azion.json or args.json configurations
  - Mentions "azion" in the context of commands, deployment, or configuration
  - Asks about edge computing, CDN deployment, or static site deployment on Azion

  Trigger even if the user doesn't explicitly say "CLI" - they might ask "how do I deploy to Azion" or "how do I create an edge application."
---

# Azion CLI Skill

Comprehensive reference for Azion CLI v4.19.1 commands, workflows, and patterns.

## Quick Start

```bash
# Install
curl -fsSL https://cli.azion.app/install.sh | bash

# Login
azion login

# Deploy a new project
azion init --name "my-project"
azion build && azion deploy --auto
```

## Decision Trees

### New Project vs Existing Project

```
User has...
├── No code yet → azion init --name "project-name"
│                  └── Use --template flag for specific templates
│
└── Existing code → azion link --name "project-name" --preset <framework>
                     └── Then: azion deploy --auto
```

### When to Use --auto

| Scenario | Use --auto? | Why |
|----------|-------------|-----|
| CI/CD pipelines | Yes | No interactive prompts |
| Quick prototypes | Yes | Faster iteration |
| First-time setup | No | Learn the options |
| Custom configuration | No | Need to specify options |

### Choosing a Preset

| Framework | Preset Value | Notes |
|-----------|--------------|-------|
| Next.js | `next` | Most common for React SSR |
| React SPA | `react` | Client-side only |
| Vue.js | `vue` | Vue 3 recommended |
| Angular | `angular` | Requires specific build |
| Astro | `astro` | Static/hybrid |
| Hexo | `hexo` | Static blog |
| Pure HTML | `static` | No build step |
| Vanilla JS | `javascript` | No framework |
| TypeScript | `typescript` | No framework |

## Command Quick Reference

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

## Common Workflows

### New Project Setup
```bash
azion login
azion init --name "my-project"
azion build
azion deploy
```

### Deploy Existing Project
```bash
azion link --name "my-app" --preset nextjs
azion deploy
```

### Configuration as Code
```bash
azion config init --config-dir ./config
# Edit azion.json as needed
azion config apply --config-dir ./config
```

### Debugging Applications
```bash
azion logs cells --tail --function-id 1234
azion logs http --tail --limit 50
azion describe application --id 1234
```

### Cache Management
```bash
azion purge --urls "example.com/path1,example.com/path2"
azion purge --wildcard "example.com/*"
azion warmup --url "https://example.com" --max-urls 500
```

## Common Pitfalls

### 1. v3 vs v4 API Confusion

**Problem:** Some users have legacy v3 API access which uses different commands.

**Solution:**
- v4 is the **primary and long-term supported API**
- v3 commands are in `pkg/v3commands/` (legacy, deprecated)
- If you see v3-style commands failing, ensure your account has v4 access
- Use `pkg/cmd/`, `pkg/api/`, `pkg/manifest/` paths for new features

### 2. Token Expiration

**Problem:** Personal tokens expire and commands suddenly fail with 401 errors.

**Solution:**
```bash
# Check current auth
azion whoami

# Create long-lived token (up to 9 months)
azion create personal-token --name "ci-token" --expiration "9m"

# Use token directly
azion --token "your-token" deploy --auto
```

### 3. Build Failures

**Problem:** Build fails due to framework-specific issues.

**Solution:**
```bash
# Skip framework build if issues
azion build --skip-framework-build

# Use custom entrypoint
azion build --entry ./src/index.js

# Enable debug logging
azion build --debug
```

### 4. Deployment Timeout

**Problem:** Deploy takes too long and times out.

**Solution:**
```bash
# Increase timeout
azion deploy --timeout 120

# Use dry-run to test
azion deploy --dry-run

# Check logs for issues
azion logs cells --tail
```

### 5. Missing Template Flag

**Problem:** Want to use a specific starter template.

**Solution:** Use the `--template` flag:
```bash
# List available templates (future feature)
azion init --help

# Use specific template
azion init --template nextjs-starter --name "my-project"
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

## Reference Files

For complete command documentation, see:
- `references/commands.md` - Full command reference with all flags
- `references/workflows.md` - Step-by-step workflow guides
- `scripts/generate_config.py` - Configuration file generator

## When to Read Reference Files

- **Need all flags for a command?** → Read `references/commands.md`
- **Setting up CI/CD?** → Read `references/workflows.md` → CI/CD Integration section
- **Creating config files programmatically?** → Run `scripts/generate_config.py`
