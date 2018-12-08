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
2. Set slack token generated for the bot in `skit.toml`
3. Run `skit -c skit.toml`
4. Go to slack and find the bot which you created and chat!

#### `skit.toml` configuration file

Following `skit.toml` file can be used to setup a simple slack bot that
passes any message sent to it to a shell script `bot.sh`. The output of
this script will be sent back as the response.

```toml
token = "your-token-here"
log_level = "info"

[[handlers]]
name = "matcha_all"
type = "command"
match = [
  ".*"
]
cmd = "./bot.sh"
args = [
  "{{ .event.Text }}"
]
```

See `examples/` for samples of different handler configurations.

### As a library

Following sample shows how to build a simple bot that echo's everything
you say to it!

```go
sk:= skit.New("your-token", logrus.New())
sk.Register("echo_all", skit.SimpleHandler("{{.event.Text}}",  ".*"))
sk.Listen(context.Background())
```


## Handlers

A handler is an implementation of `skit.Handler` interface. A handler is responsible
for consuming event, processing and responding to it when applicable. Currently 3 types
of handlers are available.

### 1. `simple`

- `simple` handler can be used to build simple regex based conversational bot.
- Following sample shows how to configure `simple` handler:

    ```toml
    [[handlers]]
    name = "simple_bot"
    type = "simple"
    match = [
      "my name is (?P<name>.*)"
    ]
    message = "Hello there {{ .name }}!"
    ```

### 2. `command`

- `command` handler can be used to delegate the message handling responsibility to external command
- This allows building bots which are more complex than the ones built with `simple`
- Following sample shows how to configure `command` handler:

    ```toml
    [[handlers]]
    name = "external_command"
    type = "command"
    match = [
      ".*"
    ]
    cmd = "./hello.sh"
    args = [
      "{{ .event.Text }}"
    ]
    timeout = "5s"
    ```

### 3. `lua`

- `lua` allows writing handlers using Lua script which will be executed using embedded [Gopher-Lua](https://gitter.im/yuin/gopher-lua)
- Lua code will have access to skit instance and its APIs and can be used build complex bots.
- You can provide paths to be added to `LUA_PATH` as `paths` in the following config.
- In the config sample below, `source` can be any valid lua code. So you can put your handler code
  in a file (e.g., `handler.lua`) under one of the `paths` and use `source="require('handler')"` in
  the handler config.
- Following sample shows how to configure `lua` handler:

    ```toml
    [[handlers]]
    name = "lua_handler"
    type = "lua"
    handler = "handle_event"
    paths = []
    source = """
      function handle_event(ctx, sk, event)
        sk:SendText(ctx, "Hello from Lua!", event.Channel)
        return true
      end
    """
    ```