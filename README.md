# eporto [![Build Status](https://travis-ci.org/fabian-z/eporto.svg?branch=master)](https://travis-ci.org/fabian-z/eporto) [![goreportcard](https://goreportcard.com/badge/github.com/fabian-z/eporto?update=1)](https://goreportcard.com/report/github.com/fabian-z/eporto)

`eporto` is a simple web application allowing you to buy and print digital stamps (Internetmarke) from Deutsche Post.

Inspired by [frank](https://github.com/gsauthof/frank), but updates product list and prices automatically.

## Installation

Currently, no pre-compiled release binaries are available. To install `eporto`, install the [Go toolchain](https://golang.org/) and run

```
go get github.com/fabian-z/eporto
```

The resulting binary will be located in `$GOPATH/bin` or in `$HOME/go/bin`.

## Usage

Check out the `docs` folder in this repository for details on how to request API credentials. The application will not work if you don't have your own API access.
Configuration is expected in a file `eporto.conf`, see the example file.

This software is currently optimized for use on Linux and uses the `lpr` command with a configured printer name to print stamps after buying.
PRs for cross-platform compatibility welcome.

## Contribution

If you find any bugs or would like to see (or contribute to) a feature, please don't hesitate to open an issue or PR.

## License

This project is licensed under the MIT License (see `LICENSE.md`)
