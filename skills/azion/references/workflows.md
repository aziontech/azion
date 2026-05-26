# Azion CLI Workflows Guide

Step-by-step guides for common Azion CLI workflows.

## Table of Contents

- [Initial Setup](#initial-setup)
- [Deploying Applications](#deploying-applications)
- [Managing Edge Functions](#managing-edge-functions)
- [Configuration as Code](#configuration-as-code)
- [Cache Management](#cache-management)
- [Object Storage](#object-storage)
- [CI/CD Integration](#cicd-integration)
- [Debugging & Troubleshooting](#debugging--troubleshooting)
- [Multi-Account Management](#multi-account-management)

---

## Initial Setup

### Install Azion CLI

```bash
# Recommended: Remote Script
curl -fsSL https://cli.azion.app/install.sh | bash

# Homebrew (macOS)
brew install azion

# Chocolatey (Windows)
choco install azion

# Winget (Windows)
winget install aziontech.azion
```

### First-Time Login

```bash
# Interactive login
azion login

# Non-interactive (CI/CD)
azion login --username your@email.com --password "your-password"

# Or use token directly
azion --token "your-personal-token" whoami
```

### Verify Setup

```bash
# Check version
azion version

# Verify authentication
azion whoami

# List your applications
azion list application
```

---

## Deploying Applications

### New Project from Template

```bash
# 1. Initialize project (interactive mode)
azion init --name "my-project"

# 2. Follow interactive prompts to select:
#    - Framework (Next.js, React, Vue, etc.)
#    - Package manager
#    - Deploy options

# 3. Or use fully automated mode with preset
azion init --name "my-project" --preset nextjs --auto

# 4. Or use a specific starter template
azion init --name "my-project" --template nextjs-starter

# 5. Build and deploy
azion build
azion deploy --auto
```

### Deploy Existing Project

```bash
# 1. Navigate to your project
cd /path/to/your/project

# 2. Link to Azion
azion link --name "my-app" --preset nextjs

# 3. Deploy
azion deploy --auto
```

### Deploy Static Files

```bash
# Deploy from specific directory
azion deploy --path ./dist

# Skip build if files are ready
azion deploy --path ./dist --skip-build
```

### Deploy with Custom Config

```bash
# Use custom config directory
azion deploy --config-dir ./custom-config

# Sync with remote before deploy
azion deploy --sync

# Local build and deploy
azion deploy --local
```

### Framework-Specific Examples

**Next.js:**
```bash
azion init --name "nextjs-app" --preset next --auto
cd nextjs-app
azion deploy --auto
```

**React:**
```bash
azion init --name "react-app" --preset react --auto
cd react-app
azion deploy --auto
```

**Vue:**
```bash
azion init --name "vue-app" --preset vue --auto
cd vue-app
azion deploy --auto
```

**Astro:**
```bash
azion init --name "astro-app" --preset astro --auto
cd astro-app
azion deploy --auto
```

**Static HTML:**
```bash
azion init --name "static-site" --preset static --auto
cd static-site
azion deploy --path ./public --auto
```

---

## Managing Edge Functions

### Create Edge Function

```bash
# Create function from code file
azion create function \
  --name "my-function" \
  --code ./functions/main.js \
  --active true

# Or use JSON file
azion create function --file ./function-config.json
```

### Function JSON Configuration

```json
{
  "name": "my-function",
  "active": true,
  "code": "async function handleRequest(request) {\n  return new Response('Hello!', { status: 200 });\n}\naddEventListener('fetch', event => {\n  event.respondWith(handleRequest(event.request));\n})"
}
```

### Create Function Instance

```bash
# Link function to application
azion create function-instance \
  --application-id 1234 \
  --function-id 5678 \
  --name "my-instance"
```

### View Function Logs

```bash
# Real-time logs
azion logs cells --tail

# Filter by function
azion logs cells --function-id 5678 --tail

# With prettified output
azion logs cells --tail --pretty --limit 50
```

### Update Function

```bash
azion update function \
  --id 5678 \
  --code ./functions/updated.js \
  --active true
```

---

## Configuration as Code

### Initialize Configuration

```bash
# Create azion.json
azion config init

# In specific directory
azion config init --config-dir ./config --force
```

### Apply Configuration

```bash
# Apply to platform
azion config apply --config-dir ./config
```

### Sync with Remote

```bash
# Pull remote config to local
azion sync

# Generate azion.config file
azion sync --iac --extension ts
```

### Sample azion.json Structure

```json
{
  "id": 1234567890,
  "name": "my-application",
  "application": {
    "active": true,
    "delivery_protocol": "https",
    "http3": true,
    "origins": [
      {
        "name": "origin-1",
        "origin_type": "single_origin",
        "address": "origin.example.com"
      }
    ]
  },
  "functions": [
    {
      "name": "edge-function-1",
      "path": "./functions/main.js"
    }
  ],
  "cache_settings": {
    "browser_cache_settings": "honor",
    "cdn_cache_settings": "override",
    "cdn_cache_settings_maximum_ttl": 3600
  }
}
```

---

## Cache Management

### Purge Cache

```bash
# Purge specific URLs
azion purge --urls "example.com/page1,example.com/page2"

# Purge by wildcard pattern
azion purge --wildcard "example.com/blog/*"

# Purge by cache key
azion purge --cachekey "example.com/@@cookie_name=cookie_value"

# Purge from tiered cache
azion purge --layer tiered_cache --urls "example.com/*"
```

### Warmup Cache

```bash
# Basic warmup
azion warmup --url "https://example.com"

# With custom settings
azion warmup \
  --url "https://example.com/products" \
  --max-urls 500 \
  --max-concurrent 5 \
  --timeout 10000
```

### Rollback Deployment

```bash
# Revert to previous deployment
azion rollback --connector-id aaaa-bbbb-cccc-dddd
```

---

## Object Storage

### Create Bucket

```bash
# Create bucket with read-only edge access
azion create storage bucket \
  --name "my-bucket" \
  --edge-access read_only

# With read-write access
azion create storage bucket \
  --name "uploads-bucket" \
  --edge-access read_write
```

### Upload Object

```bash
# Upload file to bucket
azion create storage object \
  --bucket-name "my-bucket" \
  --object-key "images/logo.png" \
  --source ./local/logo.png
```

### List Buckets

```bash
azion list storage
```

### Describe Bucket

```bash
azion describe storage --bucket-name "my-bucket"
```

---

## CI/CD Integration

### GitHub Actions

```yaml
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
        run: |
          curl -fsSL https://cli.azion.app/install.sh | bash
          echo "$HOME/.azion/bin" >> $GITHUB_PATH

      - name: Build Application
        run: azion build
        env:
          AZION_TOKEN: ${{ secrets.AZION_TOKEN }}

      - name: Deploy to Edge
        run: azion deploy --auto --silent
        env:
          AZION_TOKEN: ${{ secrets.AZION_TOKEN }}
```

### GitLab CI

```yaml
deploy:
  stage: deploy
  image: node:18
  before_script:
    - curl -fsSL https://cli.azion.app/install.sh | bash
    - export PATH="$HOME/.azion/bin:$PATH"
  script:
    - azion build
    - azion deploy --auto --silent
  variables:
    AZION_TOKEN: $AZION_TOKEN
  only:
    - main
```

### Jenkins Pipeline

```groovy
pipeline {
  agent any
  environment {
    AZION_TOKEN = credentials('azion-token')
  }
  stages {
    stage('Install CLI') {
      steps {
        sh 'curl -fsSL https://cli.azion.app/install.sh | bash'
      }
    }
    stage('Build') {
      steps {
        sh '~/.azion/bin/azion build'
      }
    }
    stage('Deploy') {
      steps {
        sh '~/.azion/bin/azion deploy --auto --silent'
      }
    }
  }
}
```

### Environment Variables Pattern

```bash
# Set token via environment
export AZION_TOKEN="your-personal-token"

# Commands will use token automatically
azion deploy --auto
```

---

## Debugging & Troubleshooting

### Enable Debug Logging

```bash
# Debug level
azion deploy --debug

# Or set log level
azion deploy --log-level debug
```

### View Logs

```bash
# Function logs (real-time)
azion logs cells --tail --pretty

# HTTP logs (real-time)
azion logs http --tail --pretty

# Filter by function ID
azion logs cells --function-id 1234 --limit 100
```

### Inspect Resources

```bash
# Application details
azion describe application --id 1234

# Function details
azion describe function --id 5678

# List all applications
azion list application --format json
```

### Common Issues

**Authentication Failed:**
```bash
# Re-login
azion logout
azion login

# Or use token
azion --token "your-token" whoami
```

**Build Fails:**
```bash
# Skip framework build
azion build --skip-framework-build

# Use custom entry
azion build --entry ./src/custom-entry.js

# Debug output
azion build --debug
```

**Deploy Timeout:**
```bash
# Increase timeout
azion deploy --timeout 120

# Or check logs
azion logs cells --tail
```

**Token Expired:**
```bash
# Create new token
azion create personal-token --name "new-token" --expiration "9m"

# Use new token
azion --token "new-token" deploy
```

---

## Multi-Account Management

### Create Profiles

```bash
# Create profile for account 1
azion create profile --name "work-account"

# Create profile for account 2
azion create profile --name "personal-account"
```

### List Profiles

```bash
azion profiles
```

### Switch Profiles

```bash
# Use specific profile
azion --config ./work-config deploy

# Or set environment
export AZION_CONFIG=./work-config
azion deploy
```

### Delete Profile

```bash
azion delete profile --name "old-profile"
```