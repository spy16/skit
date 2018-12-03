> WIP

# skit

[![GoDoc](https://godoc.org/github.com/spy16/skit?status.svg)](https://godoc.org/github.com/spy16/skit) [![Build Status](https://travis-ci.org/spy16/skit.svg?branch=master)](https://travis-ci.org/spy16/skit) [![Go Report Card](https://goreportcard.com/badge/github.com/spy16/skit)](https://goreportcard.com/report/github.com/spy16/skit)

skit is a simple Slackbot Kit in Go (or `Golang`). Skit pre-compiled
binary is good enough to build simple slack bots. For more complex usecases
skit can be used as a library as well.

## Installation

Simply download the pre-built binary for your platform from the
[Releases](https://github.com/spy16/skit/releases) section.


## Usage

### Pre-compiled Binary

1. Create a custom configuration file by refering to `./examples`
2. Add your slack api token into `skit.yaml` file (`token` field)
3. Run skit as `skit -c <your-config-file>` or `TOKEN=<slack-token> skit -c <config-file>`

> Environment variable `TOKEN` if present will override the value of `token`
> from the configuration file.

### As a library

```go
config := skit.Config{
    Token: "your-token-here",
}
sk, err := skit.New(config, logrus.New())
if err != nil {
    panic(err)
}

sk.Listen(context.Background())
```
