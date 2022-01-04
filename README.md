# azion-cli

The developer friendly with to interact with Azion!

## How to Use

You can just run `azioncli` and see it's options

```sh
$ azioncli
This is a placeholder description used while the actual description is still not ready.

Usage:
  azioncli [flags]
  azioncli [command]

Available Commands:
  configure     Configure parameters and credentials
  edge_services Manages edge services of an Azion account
  help          Help about any command
  version       Returns bin version

Flags:
  -h, --help           help for azioncli
  -t, --token string   Use provided token
  -v, --verbose        Makes azioncli verbose during the operation
      --version        version for azioncli

Use "azioncli [command] --help" for more information about a command.
```

For each subcommand you can use the `help` subcommand or with the `--help` flag
```sh
$ azioncli help edge_services
You may create, update, delete, list and describe services of an Azion account.

Usage:
  azioncli edge_services [flags]
  azioncli edge_services [command]

Available Commands:
  create      Creates a new edge service
  delete      Deletes a service based on a given service_id
  describe    Describes a service based on a given service_id
  list        Lists services of an Azion account
  resources   Manages resources in a given edge-service
  update      Updates parameters of an edge service

Flags:
  -h, --help   help for edge_services

Global Flags:
  -t, --token string   Use provided token
  -v, --verbose        Makes azioncli verbose during the operation

Use "azioncli edge_services [command] --help" for more information about a command.

```

## Building

```sh
# Build project, by default it will connect to the Stage APIs
$ make build

# Building Production version
$ make build ENVFILE=./env/prod

# Cross-Build for multiple platforms and architectures
$ make cross-build
```
