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
- [Setup Autocomplete](https://github.com/aziontech/azion-cli/wiki/Azion-CLI-autocomplete)
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
brew install azion
``````

##### Test The Azion CLI with docker 
get link release https://github.com/aziontech/azion/releases
```sh
docker run -it --rm ubuntu sh -c '
    cd
    apt-get update && apt-get install -y wget tar && \
    wget https://github.com/aziontech/azion/releases/download/1.28.1/azion_1.28.1_linux_amd64.tar.gz && \
    tar -xzf azion_1.28.1_linux_amd64.tar.gz && \
    chmod +x azion && \
    ./azion version; \
    exec sh'
```

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

- Using `azion -t [tokenvalue]` command, which saves the token in a configuration file for further use:

```
$ azion -t <authentication token>
```

- Using environment variable, which the token is cleared when the terminal is closed:

```sh
$ export AZIONCLI_TOKEN=<authentication token>
```

### Commands

Check all reference documentation for the available [commands](https://github.com/aziontech/azion-cli/wiki/azion).

### Autocomplete

It's possible to enable the autocompletion to be used with the `azion` CLI. To learn more about its settings and installation based on your OS, check the [autocompletion page](https://github.com/aziontech/azion-cli/wiki/Azion-CLI-autocomplete).

### Go Lang Update Policy

Our update policy for the Go language aims to ensure that all developers benefit from new features, bug fixes and improvements. To this end, we have adopted the practice of maintaining our project at [Release](https://go.dev/doc/devel/release) - 1 (Go 1.21 -1). We allow exceptions when vulnerabilities are found. 

By following this approach, we prioritize the stability and consistency of the project, while keeping up to date with the latest language updates. As new versions of Go are released, we will evaluate and adopt the new "[Release](https://go.dev/doc/devel/release)" version, after having sufficient time to test and validate compatibility with our existing code.

This way, our developers can take advantage of the new features, bug fixes and improvements made available by the Go community, while minimizing the possible impacts that sudden updates could have on our project.

We appreciate the understanding and collaboration of all developers who follow our update policy. Together, we guarantee a stable and up-to-date development environment for our Go project. 

## License

This project is licensed under the terms of the [MIT](LICENSE) license.
