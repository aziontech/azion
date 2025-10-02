# Azion CLI Debugging Guide

Use this guide to collect the right information and quickly diagnose issues when using the Azion CLI (`azion`). It focuses on actionable steps and the most common pitfalls observed in this repository.

## Quick Checklist

- **Enable debug logs**
  - Append `--debug` to any command to print HTTP logs, request/response diagnostics, and internal traces.
  - Example:
    ```bash
    azion list application --debug
    ```

- **Verify you are on the latest version**
  - Print your CLI version:
    ```bash
    azion version
    ```
  - If outdated, update using your package manager or reinstall from the latest release.

- **Confirm your auth context**
  - Check who is logged in:
    ```bash
    azion whoami
    ```
  - If you see an error like `Not logged in`, run:
    ```bash
    azion login
    ```

    or

    ```bash
    azion -t <PERSONAL_TOKEN>
    ```

## Deploying Projects

- **Match your account version with the correct config format**
  - V3 accounts should use the legacy `azion.json` (v3) format.
  - V4 accounts should use the new `azion.json` (v4) format.
  - Mixing formats will cause parsing or ID conversion errors (e.g. messages like `ErrorConvertId`, `ErrorConvertIdFunction`).
  - The same goes for your azion.config file. You should respect the [V3](https://github.com/aziontech/lib/blob/v1.20.6/packages/config/src/configProcessor/helpers/azion.config.example.ts)/[V4](https://github.com/aziontech/lib/blob/v2.1.2/packages/config/src/configProcessor/helpers/azion.config.example.ts) format.

- **If link/init + deploy fails, try deploy directly**
  - Some link/init paths automate multiple steps and may hide the root cause. If it fails, run `deploy` directly to isolate the issue:
    ```bash
    azion deploy --debug
    ```

    or

    ```bash
    azion deploy --local --debug
    ```
  - If direct deploy succeeds but link/init fails, open an [issue](https://github.com/aziontech/azion/issues) with logs:
    - Command used, full output with `--debug`, your OS/arch, CLI version, and a sanitized `azion.json` snippet.

- **Prefer local mode for debugging**
  - The `--local` flag for `deploy` runs a local deploy and is easier to debug issues in `manifest.json`/rules/application mapping.
  - Example:
    ```bash
    azion deploy --local --debug
    ```

### Handling "name already in use" errors

- Check existing resources via CLI
  - Use list commands to confirm whether the resource already exists. If it does, future deploys should update the existing resource instead of creating a new one:
    ```bash
    azion list application
    azion list workload
    azion list function
    azion list cache-setting
    ```
  - If you find existing resources, you can manually edit your `azion.json` to include the corresponding IDs.

- You may also verify in the [Azion Console](https://console.azion.com)
  - Visit [Azion Console](https://console.azion.com) and check whether the resource you are trying to create already exists (Application, Workload, Function, Cache Setting, etc.).
  - This is helpful for a visual confirmation and for reviewing additional details (owners, timestamps, references) that might explain conflicts.

- If resources were created but IDs were not recorded in `azion.json`
  - If you don't mind redoing the deploy from scratch, you can remove the resources referenced by your `azion.json` in one go:
    ```bash
    # Be at your project's root (where azion/azion.json lives)
    azion delete application --cascade
    ```
  - The cascade delete will attempt to delete the application and related resources declared in your `azion.json`, allowing a clean redeploy.

- Use unlink for a guided cleanup
  - The `unlink` command also offers to cascade delete and clean local state:
    ```bash
    azion unlink
    ```
  - This is helpful when your local project has drifted from the remote state and you want to re-link or re-init cleanly.
  - If the unlink option is the one chose, remember to run `azion link` again to re-link your project.

## Common Issues and How to Resolve

- **Parsing input or prompt errors**
  - Many commands use interactive prompts. If you see `Failed to parse your response` or `utils.ErrorParseResponse`, try supplying all required flags explicitly to avoid prompts, or re-run with `--debug` to see where parsing failed.

- **Config discovery and azion.json reads**
  - The CLI uses the working directory to locate config. Use `pwd` to ensure you’re at your project root.
  - If `azion.json` is missing, commands may return errors like `ErrorMissingAzionJson`. Create or link your project before running deploy-related commands.

- **HTTP 4xx/5xx troubleshooting**
  - The CLI maps status codes to friendly errors (see `utils/errors.go`). For 400, the body is scanned for keys like `detail`, `invalid_order_field`, `minimum_tls_version`, and name-in-use hints.
  - For 500s, errors are normalized (timeouts vs. generic 500). Provide full `--debug` logs when reporting issues. More often that not, 5xx errors are API problems. Please check [azion status page](https://status.azion.com/) for more information.

## Collecting Diagnostics for a Bug Report

- **Minimum data to include**
  - Command executed with full args (mask tokens) and `--debug` output.
  - CLI version from `azion version` and OS (e.g., `uname -a`).
  - Whether your account is V3 or V4 and which `azion.json` format you are using.
  - Sanitized `azion.json` (remove secrets and proprietary values).

- **Repro steps**
  - List exact steps starting from a clean directory or a minimal sample repo, including file tree and contents relevant to the issue (e.g., `manifest.json`, `azion.json`, code files referenced by `--code`).

## Tips for Contributors (Developers of this repo)

- **Run vet and linters locally**
  ```bash
  go vet ./...
  golangci-lint run ./...
  ```

- **Avoid brittle tests**
  - Match mock endpoints with actual clients.
  - Use errors from `messages/...` and wrap with `%w`.
  - When asserting errors, compare with `errors.Is` or string equality against constants when needed.

- **CI nuances**
  - The workflow may skip v3 tests while v3 is being sunset. Run all tests locally if you’re changing shared code:
    ```bash
    make test
    ```

## Still stuck?

- Search existing issues and discussions in the repo.
- Open a new issue with the diagnostics above.
- If it’s a blocking production case, contact Azion support with your `--debug` logs and `azion version` output.
- Check our [help channels](https://www.azion.com/en/documentation/products/get-help/) for more information.
