# CLI Metrics Collection and Analytics

This document explains how the Azion CLI collects, stores, and sends usage metrics to Segment for analytics.

## Overview

The Azion CLI collects anonymous usage metrics to better understand user needs and enhance the application. Metrics collection is **opt-in** - users are asked for consent the first time they use the CLI or create a new profile.

## User Consent Flow

### First-Time Consent

When a user runs a command for the first time (without having previously answered the metrics consent question), the CLI prompts:

```
To better understand user needs and enhance our application, we gather anonymous data. Do you agree to participate? (Y/n)
```

This consent flow is handled in [`pkg/cmd/root/pre_command.go`](pkg/cmd/root/pre_command.go:239):

```go
// 0 = authorization was not asked yet, 1 = accepted, 2 = denied
func checkAuthorizeMetricsCollection(cmd *cobra.Command, globalFlagAll bool, settings *token.Settings, activeProfile string) error {
    if settings.AuthorizeMetricsCollection > 0 || cmd.Name() == "completion" {
        return nil
    }

    authorize := confirmFn(globalFlagAll, msg.AskCollectMetrics, true)
    if authorize {
        settings.AuthorizeMetricsCollection = 1
    } else {
        settings.AuthorizeMetricsCollection = 2
    }

    if err := token.WriteSettings(*settings, activeProfile); err != nil {
        return err
    }

    return nil
}
```

### Consent Values

The `AuthorizeMetricsCollection` field in [`token.Settings`](pkg/token/models.go:37) stores the user's consent status:

| Value | Meaning |
|-------|---------|
| `0` | Not asked yet - user hasn't been prompted |
| `1` | **Accepted** - user agreed to share metrics |
| `2` | **Denied** - user declined to share metrics |

### When Consent is Prompted

- First time using the CLI with a new profile
- Creating a new profile (see [`pkg/cmd/create/profile/profile.go`](pkg/cmd/create/profile/profile.go:64))
- Consent is **not** prompted for the `completion` command (internal shell completion)

## Metrics Collection

### Collection Trigger

Metrics are collected after every command execution in [`pkg/cmd/root/root.go`](pkg/cmd/root/root.go:237):

```go
func Execute(f *factoryRoot) {
    logger.New(zapcore.InfoLevel)

    cmd := f.CmdRoot()
    err := cmd.Execute()
    executionTime := time.Since(f.startTime).Seconds()

    // 1 = authorize; anything different than 1 means that the user did not authorize
    if f.globalSettings != nil {
        if f.globalSettings.AuthorizeMetricsCollection == 1 {
            activeProfile := f.factory.GetActiveProfile()
            errMetrics := metric.TotalCommandsCount(cmd, f.commandName, executionTime, err, activeProfile)
            if errMetrics != nil {
                logger.Debug("Error while saving metrics", zap.Error(err))
            }
        }
    }
    // ...
}
```

### Collection Process

1. **Command Name Extraction**: The command path is dynamically extracted and rewritten (e.g., `azion list application` becomes `list-application`)
2. **Execution Tracking**: Execution time is measured from command start
3. **Success/Failure Status**: Determined by whether the command returned an error
4. **Local Storage**: Metrics are stored locally in `metrics.json` per profile

### Metrics Data Structure

The [`command`](pkg/metric/count.go:17) struct defines the collected metrics:

```go
type command struct {
    TotalSuccess   int     // Number of successful executions
    TotalFailed    int     // Number of failed executions
    ExecutionTime  float64 // Last execution time in seconds
    CLIVersion     string  // CLI version used
    VulcanVersion  string  // Deprecated: kept for legacy metrics compatibility
    BundlerVersion string  // Bundler version
    Shell          string  // User's shell (bash, zsh, etc.)
}
```

### Ignored Commands

The following commands are **excluded** from metrics collection (see [`pkg/metric/count.go`](pkg/metric/count.go:26)):

```go
var ignoredCommands = map[string]bool{
    "__complete": true,  // Shell completion internal command
    "completion": true,  // Shell completion generation command
}
```

### Local Storage

Metrics are stored locally in the profile's configuration directory:

```
~/.config/azion/<profile>/metrics.json
```

The file contains a map of command names to their metrics:

```json
{
  "list-application": {
    "TotalSuccess": 5,
    "TotalFailed": 1,
    "ExecutionTime": 1.234,
    "CLIVersion": "2.0.0",
    "VulcanVersion": "1.5.0",
    "BundlerVersion": "1.5.0",
    "Shell": "/bin/zsh"
  },
  "deploy": {
    "TotalSuccess": 10,
    "TotalFailed": 2,
    "ExecutionTime": 45.678,
    "CLIVersion": "2.0.0",
    "VulcanVersion": "1.5.0",
    "BundlerVersion": "7.2.0",
    "Shell": "/bin/zsh"
  }
}
```

## Sending Metrics to Segment

### Timing: Every 24 Hours

Metrics are sent to Segment **once every 24 hours** at the beginning of the first command execution after the 24-hour window has passed. This is handled in [`pkg/cmd/root/pre_command.go`](pkg/cmd/root/pre_command.go:141):

```go
func checkForUpdateAndMetrics(cVersion string, f *cmdutil.Factory, settings *token.Settings) error {
    logger.Debug("Verifying if an update is required")
    activeProfile := f.GetActiveProfile()
    // checks if 24 hours have passed since the last check
    if time.Since(settings.LastCheck) < 24*time.Hour && !settings.LastCheck.IsZero() {
        return nil
    }

    // checks if user is Logged in before sending metrics
    if verifyUserInfo(settings) {
        metric.Send(settings, activeProfile)
    }
    // ... update check continues
}
```

### Prerequisites for Sending

Before sending metrics, the CLI verifies:

1. **24 hours have passed** since the last check (`settings.LastCheck`)
2. **User is logged in** - has both `ClientId` and `Email` (see [`verifyUserInfo()`](pkg/cmd/root/pre_command.go:258))
3. **User has consented** to metrics collection (`AuthorizeMetricsCollection == 1`)

### Segment Integration

Metrics are sent using the Segment analytics-go SDK (see [`pkg/metric/segment.go`](pkg/metric/segment.go)):

```go

func Send(settings *token.Settings, profile string) {
    client := analytics.New(SEGMENT_KEY)
    defer client.Close()

    metrics := readLocalMetrics(profile)

    os := runtime.GOOS
    arch := runtime.GOARCH

    for event, cmd := range metrics {
        err := client.Enqueue(analytics.Track{
            UserId: settings.ClientId,
            Event:  fmt.Sprintf("cli_%s", event),
            Properties: analytics.NewProperties().
                Set("email", settings.Email).
                Set("cli version", cmd.CLIVersion).
                Set("vulcan version", cmd.VulcanVersion). // Deprecated: kept for legacy metrics compatibility
                Set("bundler version", cmd.BundlerVersion).
                Set("total successful", cmd.TotalSuccess).
                Set("total failed", cmd.TotalFailed).
                Set("total", cmd.TotalSuccess+cmd.TotalFailed).
                Set("shell", cmd.Shell).
                Set("execution time", cmd.ExecutionTime).
                Set("operating system", os).
                Set("architecture", arch).
                Set("client id", settings.ClientId),
        })
        // ...
    }

    clean(profile) // Clear metrics after sending
}
```

### Properties Sent to Segment

Each command event includes the following properties:

| Property | Description | Source |
|----------|-------------|--------|
| `email` | User's email address | User settings |
| `cli version` | CLI version used | Build version |
| `vulcan version` | Vulcan bundler version (deprecated) | GitHub API |
| `bundler version` | Bundler version | GitHub API |
| `total successful` | Successful execution count | Local metrics |
| `total failed` | Failed execution count | Local metrics |
| `total` | Total executions (success + failed) | Calculated |
| `shell` | User's shell path | System detection |
| `execution time` | Last execution time (seconds) | Timer measurement |
| `operating system` | OS (darwin, linux, windows) | runtime.GOOS |
| `architecture` | CPU architecture | runtime.GOARCH |
| `client id` | User's client ID | User settings |

### Event Naming

Events are prefixed with `cli_` followed by the command name:

- `cli_list-application`
- `cli_deploy`
- `cli_create-application`
- etc.

### Post-Send Cleanup

After successfully sending metrics to Segment, the local metrics file is cleared:

```go
func clean(profile string) {
    err := os.WriteFile(location(profile), []byte{}, 0666)
    if err != nil {
        return
    }
}
```

This ensures metrics are not sent twice and starts a fresh collection period.

## Data Flow Summary

The metrics collection and sending process follows these steps:

1. **User runs a command**
   - The CLI entry point is [`Execute()`](pkg/cmd/root/root.go:237)

2. **Pre-Command Checks** (in [`doPreCommandCheck()`](pkg/cmd/root/pre_command.go:33))
   - Check if user consent was given (`AuthorizeMetricsCollection`)
   - If not asked yet, prompt user for consent
   - Check if 24 hours have passed since last metrics send (`settings.LastCheck`)
   - Verify user is logged in (`ClientId && Email`)

3. **Send Metrics to Segment** (if conditions met)
   - Read local `metrics.json` file
   - Enqueue each command's metrics to Segment API
   - Clear local metrics file after successful send
   - Update `LastCheck` timestamp to current time

4. **Command Executes**
   - The actual CLI command runs

5. **Post-Command Metrics Collection**
   - Extract command name dynamically
   - Measure execution time
   - Determine success/failure status
   - Save to local `metrics.json` (if user consented)

## File Locations

| File | Path | Purpose |
|------|------|---------|
| Settings | `~/.azion/<profile>/settings.toml` | User settings including consent |
| Metrics | `~/.azion/<profile>/metrics.json` | Local metrics storage |
| Profiles | `~/.azion/profiles.json` | Profile currently in use |

## Key Implementation Files

| File | Description |
|------|-------------|
| [`pkg/cmd/root/root.go`](pkg/cmd/root/root.go) | Metrics collection trigger after command execution |
| [`pkg/cmd/root/pre_command.go`](pkg/cmd/root/pre_command.go) | Consent handling and metrics sending logic |
| [`pkg/metric/count.go`](pkg/metric/count.go) | Metrics collection and local storage |
| [`pkg/metric/segment.go`](pkg/metric/segment.go) | Segment API integration |
| [`pkg/token/models.go`](pkg/token/models.go) | Settings data structure |
| [`messages/root/messages.go`](messages/root/messages.go) | User-facing messages |

## Privacy Considerations

1. **Opt-in by default**: Users must explicitly consent before any metrics are collected
2. **Transparent**: Users are informed about data collection purpose
3. **Minimal data**: Only essential usage data is collected
4. **User identification**: Uses Client ID and email from user's authenticated session
5. **No sensitive data**: Command arguments and flags are not recorded
6. **Local-first**: Metrics are stored locally and sent in batches every 24 hours
