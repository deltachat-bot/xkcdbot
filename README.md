#  XKCD Bot

![Latest release](https://img.shields.io/github/v/tag/deltachat-bot/xkcdbot?label=release)
![Go version](https://img.shields.io/github/go-mod/go-version/deltachat-bot/xkcdbot)
[![CI](https://github.com/deltachat-bot/xkcdbot/actions/workflows/ci.yml/badge.svg)](https://github.com/deltachat-bot/xkcdbot/actions/workflows/ci.yml)

Small bot that allows to get [XKCD](https://xkcd.com) comics in Delta Chat.

## Install

Binary releases can be found at: https://github.com/deltachat-bot/xkcdbot/releases

To install from source:

```sh
go install github.com/deltachat-bot/xkcdbot@latest
```

### Installing deltachat-rpc-server

This program depends on a standalone Delta Chat RPC server `deltachat-rpc-server` program.
For installation instructions check:
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
