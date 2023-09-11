# Azion CLI
[![MIT License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![CLI Reference](https://img.shields.io/badge/cli-reference-green.svg)](https://github.com/aziontech/azion-cli/wiki/azion)
[![Go Report Card](https://goreportcard.com/badge/github.com/aziontech/azion-cli)](https://goreportcard.com/report/github.com/aziontech/azion-cli)


The Azion CLI (command-line interface) is an open source tool that enables you to manage any Azion service via command line. Through it, you can manage all Azion products, create automations using CI/CD scripts or pipelines, provision multiple services that make up your application with a few commands, and also manage your Azion configurations as code.

The developer friendly way to interact with Azion!

## Quick links

- [Downloading](#downloading)
- [Building](#building)
- [Setup Autocomplete](#setup-autocomplete)
- [How to Use](#How-to-Use)
- [Commands Reference](https://github.com/aziontech/azion-cli/wiki/azion)
- [Contributing](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [License](#License)


## Downloading

>**Attention**: if you've downloaded `azioncli` in an older version than 1.0.0, it's highly recommended to uninstall it before downloading `azion` cli.

There are two ways to download and use the `azion` CLI.
The first, is the regular way of cloning this repository and [building](#building) the project manually.
However, `azion` is also available as `homebrew`, `rpm`, `deb` and `apk` packages.

To use `rpm`, `deb` and `apk` packages, please visit our [releases](https://github.com/aziontech/azion-cli/releases) page, and download the desired package. 

To download azion CLI through Homebrew, run:

- `brew install aziontech/tap/azioncli`


## Building Locally

```sh
# Build project, by default it will connect to Production APIs
$ make build

# Cross-Build for multiple platforms and architectures
$ make cross-build
```

---

## Setup Autocomplete

> Please verify if you have autocomplete enabled globally, otherwise the autocomplete for cli won't work

## Dependencies zsh

You need to install zsh-autosuggestions

MacOs:

```shell
brew install zsh-autosuggestions
```

Linux:

```shell
git clone https://github.com/zsh-users/zsh-autosuggestions \${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
```

---

## Installing autocomplete for azion cli only

### MacOs / Linux

Run
```shell
echo "autoload -U compinit; compinit" >> ~/.zshrc 
```
Then run
```shell
echo "source <(azioncli completion zsh); compdef _azioncli azioncli" >> ~/.zshrc
```

If you uninstall azion cli, please edit your `~/.zshrc` file and remove the line `source <(azioncli completion zsh); compdef _azioncli azioncli`

---

## Installing autocomplete globally

### MacOs

Run:

```shell
echo "autoload -U compinit; compinit" >> ~/.zshrc 
```

Then open your `~/.zshrc` file, and add the following content to it

```shell
if type brew &>/dev/null
then
  FPATH="$(brew --prefix)/share/zsh/site-functions:${FPATH}"

  autoload -Uz compinit
  compinit
fi
```

### Linux

Open your `~/.zshrc` file and add the following content to it

```shell
plugins=(zsh-autosuggestions)
````

If you have other plugins, just add zsh-zutosuggestions to the end.

> Whether you chose to activate autocomplete globally or for azion cli only, the steps of each section should only be run once. After that, autocomplete will work every time you open a new terminal.

### Dependencies bash

You need to install bash-completion

MacOs:

```shell
brew install bash-completion
```

Centos/RHEL 7:

```shell
yum install bash-completion bash-completion-extras
```

Debian/Ubuntu:

```shell
apt-get install bash-completion
```

Alpine:

```shell
apk add bash-completion
```

---

## Installing autocomplete for azion cli only
### MacOs / Linux

Run:

```shell
echo "source <(azioncli completion bash)" >> ~/.bashrc 
```

If you uninstall azion cli, please edit your `~/.bashrc` file and remove the line `source <(azioncli completion bash)`

---

## Installing autocomplete globally
### MacOS

Open your `~/.bashrc` file, and add the following content to it
```shell
BREW_PREFIX=$(brew --prefix)
[[ -r "${BREW_PREFIX}/etc/profile.d/bash_completion.sh" ]] && . ${BREW_PREFIX}/etc/profile.d/bash_completion.sh
```

#### Linux

Run:

```shell
echo "source /etc/profile.d/bash_completion.sh" >> ~/.bashrc
```

> Whether you chose to activate autocomplete globally or for azion cli only, the steps of each section should only be run once. After that, autocomplete will work every time you open a new terminal.

---
### Dependencies fish

Run the command below once:

```shell
echo "azioncli completion fish | source" >> ~/.config/fish/config.fish
```

If you uninstall azion cli, please edit your `~/.config/fish/config.fish` file and remove the line `azioncli completion fish | source`

## How to Use

In order to perform network operations it is *mandatory* to provide an authentication token.

You can provide token in two ways:

- Using azion-cli token command (this command saves the token in a configuration file for further use):
$ azion -t <authentication token>

- Using environment variable (in this way the token will be cleared when the terminal is closed):
$ export AZIONCLI_TOKEN=<authentication token>


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

## License

This project is licensed under the terms of the [MIT](LICENSE) license.


