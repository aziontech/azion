# Azion CLI
[![MIT License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![CLI Reference](https://img.shields.io/badge/cli-reference-green.svg)](https://github.com/aziontech/azion-cli/wiki/azioncli)


The Azion CLI (command-line interface) is an open source tool that enables you to manage any Azion service via command line. Through it, you can manage all Azion products, create automations using CI/CD scripts or pipelines, provision multiple services that make up your application with a few commands, and also manage your Azion configurations as code.

The developer friendly way to interact with Azion!

## Quick links
* [Building](#building)
* [How to Use](#How-to-Use)
* [Commands Reference](https://github.com/aziontech/azion-cli/wiki/azioncli)
* [Contributing](CONTRIBUTING.md)
* [Code of Conduct](CODE_OF_CONDUCT.md)
* [License](#License)


## Building

```sh
# Build project, by default it will connect to the Stage APIs
$ make build

# Building Production version
$ make build ENVFILE=./env/prod

# Cross-Build for multiple platforms and architectures
$ make cross-build
```

## How to Use

In order to perform network operations it is *mandatory* to provide an authentication token

You can provide token in two ways.
* Using azion-cli configure command (this command saves the token in a configuration file for further use):
$ azioncli configure -t <authentication token>

OR

* Using environment variable (in this way the token will be cleared when the terminal is closed):
$ export AZIONCLI_TOKEN=<authentication token>


You can just run `azioncli` and see it's options

```text
$ azioncli
Interact easily with Azion services

USAGE
  azioncli [flags]

API COMMANDS
  edge_functions: Manages your Azion account's Edge Functions
  edge_services: Manages your Azion account's Edge Services

ADDITIONAL COMMANDS
  configure:     Configure parameters and credentials
  help:          Help about any command
  version:       Returns the binary version

FLAGS
  -h, --help      help for azioncli
  -v, --verbose   Makes azioncli verbose during the operation
      --version   version for azioncli

LEARN MORE
  Use 'azioncli <command> <subcommand> --help' for more information about a command
```

For each subcommand you the `-h|--help` flag to learn more about it:

```text
$ ./bin/azioncli edge_functions --help
You can create, update, delete, list and describe your Azion account's Edge Functions

USAGE
  azioncli edge_functions [flags]

COMMANDS
  create:     Create a new Edge Function
  delete:     Deletes an Edge Function
  describe:   Describes an Edge Function
  list:       Lists your account's Edge Functions
  update:     Updates an Edge Function

FLAGS
  -h, --help   help for edge_functions

INHERITED FLAGS
  -v, --verbose   Makes azioncli verbose during the operation

LEARN MORE
  Use 'azioncli <command> <subcommand> --help' for more information about a command
```

## License

This project is licensed under the terms of the [MIT](LICENSE) license.

