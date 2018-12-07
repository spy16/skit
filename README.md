# skit

[![GoDoc](https://godoc.org/github.com/spy16/skit?status.svg)](https://godoc.org/github.com/spy16/skit) [![Build Status](https://travis-ci.org/spy16/skit.svg?branch=master)](https://travis-ci.org/spy16/skit) [![Go Report Card](https://goreportcard.com/badge/github.com/spy16/skit)](https://goreportcard.com/report/github.com/spy16/skit)

Skit is a simple tool/library written in Go (or `Golang`) for building Slack bots.
Skit pre-compiled binary is good enough to build simple slack bots. For more complex
usecases skit can be used as a library as well.

## Installation

Simply download the pre-built binary for your platform from the
[Releases](https://github.com/spy16/skit/releases) section.


## Usage

### Pre-compiled Binary

Release archive will contain a `skit.toml` file with some sample handlers
setup. To run this default setup:

1. Create a bot on slack by following [Documentation](https://api.slack.com/bot-users#creating-bot-user)
2. Set slack token generate for the bot in `skit.toml`
3. Run `skit -c skit.toml`
4. Go to slack and find the bot which you created and chat!

### As a library

Following sample shows how to build a simple bot that echo's everything
you say to it!

```go
config := skit.Config{
    Token: "your-token-here",
}
sk, err := skit.New(config, logrus.New())
if err != nil {
    panic(err)
}

sk.Register("echo_all", handlers.Echo(".*"))

sk.Listen(context.Background())
```
