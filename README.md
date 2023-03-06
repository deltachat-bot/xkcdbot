#  deltabot-cli for Go

[![CI](https://github.com/deltachat-bot/xkcdbot/actions/workflows/ci.yml/badge.svg)](https://github.com/deltachat-bot/xkcdbot/actions/workflows/ci.yml)
![Go version](https://img.shields.io/github/go-mod/go-version/deltachat-bot/xkcdbot)

Small bot that allows to get [XKCD](https://xkcd.com) comics in Delta Chat.

## Install

```sh
go install github.com/deltachat-bot/xkcdbot@latest
```

### Installing deltachat-rpc-server

This program depends on a standalone Delta Chat RPC server `deltachat-rpc-server` program that must be
available in your `PATH`. For installation instructions check:
https://github.com/deltachat/deltachat-core-rust/tree/master/deltachat-rpc-server

## Usage

Configure the bot:

```sh
xkcdbot init bot@example.com PASSWORD
```

Start the bot:

```sh
xkcdbot serve
```

Run `xkcdbot --help` to see all available options.
