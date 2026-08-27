package github

const workflowContent = `name: Deploy Application using Azion CLI

on:
  workflow_dispatch:
    inputs:
      branch:
        required: true
        type: choice
        default: main
        options:
          - main

jobs:
  deploy:
    name: Deploy
    runs-on: ubuntu-latest

    permissions:
      contents: write

    steps:
      - name: Checkout
        uses: actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1
        with:
          fetch-depth: 0

      - name: Use Node.js 22.x
        uses: actions/setup-node@820762786026740c76f36085b0efc47a31fe5020 # v7.0.0
        with:
          node-version: 22

      - name: Detect package manager
        id: pm
        run: |
          if [ -f pnpm-lock.yaml ]; then
            echo "name=pnpm" >> "$GITHUB_OUTPUT"
          elif [ -f yarn.lock ]; then
            echo "name=yarn" >> "$GITHUB_OUTPUT"
          else
            echo "name=npm" >> "$GITHUB_OUTPUT"
          fi

      - name: Enable Corepack
        if: steps.pm.outputs.name != 'npm'
        run: corepack enable

      - name: Install dependencies (npm)
        if: steps.pm.outputs.name == 'npm'
        run: |
          if [ -f package-lock.json ]; then npm ci; else npm install; fi

      - name: Install dependencies (pnpm)
        if: steps.pm.outputs.name == 'pnpm'
        run: pnpm i --ignore-scripts --frozen-lockfile && pnpm approve-builds --all && pnpm i --frozen-lockfile

      - name: Install dependencies (yarn)
        if: steps.pm.outputs.name == 'yarn'
        run: yarn install --immutable

      - name: Install Azion CLI
        run: |
          curl -o azionlinux https://downloads.azion.com/linux/x86_64/azion
          sudo mv azionlinux /usr/bin/azion
          sudo chmod u+x /usr/bin/azion

      - name: CLI version
        run: azion --version

      # Configure a personal token in your github secrets
      # You may create a personal token by running 'azion create personal-token'
      - name: Configure token
        run: |
          azion -t ${{ secrets.AZION_PERSONAL_TOKEN }}  
          azion whoami

      - name: Azion Build
        run: |
          azion build

      # You may add the --sync flag to sync local and remote resources
      - name: Azion Deploy
        run: |
          azion deploy --local --debug

      - name: Commit Azion files
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
          git add azion/azion.json azion.config.* || true
          # Commit only if there are staged changes
          if git diff --cached --quiet; then
            echo "No Azion changes to commit. Skipping push."
          else
            git commit -m "chore: update azion files"
            # Rebase in case remote has new commits, then push
            git pull --rebase origin ${{ inputs.branch }}
            git push
          fi
`
