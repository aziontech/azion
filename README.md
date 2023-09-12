# Azion CLI
[![MIT License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![CLI Reference](https://img.shields.io/badge/cli-reference-green.svg)](https://github.com/aziontech/azion-cli/wiki/azion)
[![Go Report Card](https://goreportcard.com/badge/github.com/aziontech/azion-cli)](https://goreportcard.com/report/github.com/aziontech/azion-cli)

**Azion CLI** is a user-friendly way to work with the Azion Edge Platform, allowing you to create and manage applications through simple commands. It makes possible the initialization, build, and deployment of applications, from simple static pages to different frameworks, such as:

- Next.js 
- Vue
- Angular
- Astro
- Hexo
- Vite

Through it, you can manage all Azion products, create automation using CI/CD scripts or pipelines, provision multiple services that make up your application with a few commands, and also manage your Azion configurations as code.

The developer-friendly way to interact with Azion!

## Quick links

- [Downloading](#downloading)
- [Building](#building)
- [Setup Autocomplete](https://github.com/aziontech/azion-cli/wiki/Azion-CLI-autocompletion)
- [How to Use](#How-to-Use)
- [Commands Reference](https://github.com/aziontech/azion-cli/wiki/azion)
- [Contributing](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [License](#License)


## Downloading

>**Attention**: if you've downloaded `azioncli` in an older version than 1.0.0, it's highly recommended to uninstall it before downloading `azion` CLI.

There are two ways to download and use the `azion` CLI:

- Cloning this repository and [building](#building) the project manually.
- Package managers, since `azion` is also available as `homebrew`, `rpm`, `deb` and `apk` packages.

To use `rpm`, `deb` and `apk` packages, please visit our [releases](https://github.com/aziontech/azion-cli/releases) page, and download the desired package.

To download azion CLI through Homebrew, run:

```sh
brew install aziontech/tap/azioncli
``````

## Building Locally

```sh
# Build project, by default it will connect to Production APIs
$ make build

# Cross-Build for multiple platforms and architectures
$ make cross-build
```

---


## How to Use

### Authentication

In order to perform network operations it is *mandatory* to provide [an authentication token](https://www.azion.com/en/documentation/products/accounts/personal-tokens/).

You can provide the token in two ways:

- Using `azion-cli token [tokenvalue]` command, which saves the token in a configuration file for further use:

```
$ azion -t <authentication token>
```

- Using environment variable, which the token is cleared when the terminal is closed:

```sh
$ export AZIONCLI_TOKEN=<authentication token>
```

### Help

You can just run `azion -h` and see it's options

```sh
$ azion -h

DESCRIPTION
  The Azion Command Line Interface is a unified tool to manage your Azion projects and resources

SYNOPSIS
  azion <command> <subcommand> [flags]

EXAMPLES
  $ azion
  $ azion -t azionxxxxxx
  $ azion --debug
  $ azion -h
  

AVAILABLE COMMANDS
  build          Builds an edge application locally
  deploy         Deploys an application on the Azion platform
  dev            Starts a local development server
  help           Help about any command
  init           Initializes an edge application from a template
  link           Links a local application to an Azion edge application
  personal_token Manages the personal tokens configured on the Azion platform

LOCAL OPTIONS
  -c, --config string      Sets the Azion configuration folder for the current command only, without changing persistent settings.
  -d, --debug              Displays log at a debug level
  -h, --help               Displays more information about the Azion CLI
  -l, --log-level string   Displays log at a debug level (default "info")
  -s, --silent             Silences log completely; mostly used for automation purposes
  -t, --token string       Saves a given personal token locally to authorize CLI commands
  -v, --version            version for azion
  -y, --yes                Answers all yes/no interactions automatically with yes
  

LEARN MORE
  
  Use 'azion <command> <subcommand> --help' for more information about a command
```

For each command or subcommand use the `-h|--help` flag to learn more about it:

```sh
$ azion dev --help
Azion CLI 1.0.0

DESCRIPTION
  Start a development server locally

SYNOPSIS
  azion dev [flags]

EXAMPLES
         
  $ azion dev
  $ azion dev --help
  

LOCAL OPTIONS
  -h, --help   Displays more information about the dev command
  

GLOBAL OPTIONS
  -c, --config string      Sets the Azion configuration folder for the current command only, without changing persistent settings.
  -d, --debug              Displays log at a debug level
  -l, --log-level string   Displays log at a debug level (default "info")
  -s, --silent             Silences log completely; mostly used for automation purposes
  -t, --token string       Saves a given personal token locally to authorize CLI commands
  -y, --yes                Answers all yes/no interactions automatically with yes
  

LEARN MORE
  
  Use 'azion <command> <subcommand> --help' for more information about a command
```

### Autocompletion

It's possible to enable the autocompletion to be used with the `azion` CLI. To learn more about its settings and installation based on your OS, check the [autocompletion page](https://github.com/aziontech/azion-cli/wiki/Azion-CLI-autocompletion).

## License

This project is licensed under the terms of the [MIT](LICENSE) license.
