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
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Use Node.js 20.x
        uses: actions/setup-node@v4
        with:
          node-version: 20

	  # If you want to change package manager, change the command below and add the necessary configurations
      - name: Install dependencies
        run: npm install

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
          azion deploy --local --skip-build

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
