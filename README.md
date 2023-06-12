# README

This tool is a Command Line Interface to work with Seeed devices, the *Wio Link* and the *Wio Node*. This is a re-write of the original [Wio Link CLI](
https://github.com/Seeed-Studio/wio-cli) in Python. The purpose of rewriting the CLI in Python is to make it easier to install and use on different platforms.
There is also some additional functionality added to the CLI.

This tool will use the same format of configuration file as the origional CLI. The location is `~/.wio/config.json` on Linux and Mac OS X and `%USERPROFILE%\.wio\config.json` on Windows.

## Installation

To get started simply run the following command:

```bash
go install github.com/gabeduke/wio-cli-go@main

wio-cli-go --help
This is a rewrite of the Wio CLI in Go. It is a CLI for interacting with the Wio Link IoT Platform.

The original Wio CLI was written in Python and can be found here: https://github.com/Seeed-Studio/wio-cli

This CLI is a work in progress and is not yet feature complete. It is being developed by Gabriel Duke (gabeduke@gmail.com)
The goal of this rewrite is to make the CLI more performant and easier to maintain.

Usage:
  wio [command]

Aliases:
  wio, wio-cli-go, wio-cli

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  nodes       Manage your Wio Nodes
  user        Manage user settings

Flags:
      --config string            config file (default is $HOME/.wio.json)
  -h, --help                     help for wio
  -l, --log-level logLevelEnum   log level: "info", "debug", "warn", "error" (default is warn)
  -t, --toggle                   Help message for toggle

Use "wio [command] --help" for more information about a command.
```

## Usage

### User

The CLI is organized into subcommands. The first subcommand is `user`. This command is used to setup the configuration 
file. The configuration file is used to store your Wio API key and the default Wio Node to use. The configuration file 
is located at `~/.wio/config.json` on Linux and Mac OS X and `%USERPROFILE%\.wio\config.json` on Windows. You can define
a custom path by passing the `--config` flag to the CLI.

The tool will first search in command line flags and configuration file for an email address, if not found then it will 
prompt for user input. It will then prompt for the location of the Wio server being used and search for the server IP address.
Once the information has been gathered it will write everything to the configuration file.

You may also call `login` directly or `create` to create a new user account.

### Nodes

The `nodes` subcommand is used to manage your Wio Nodes. You can add, remove, and list your nodes. You can also set the
default node to use for other commands.

To Register a nodes you will need to use the `register` subcommand. This will prompt you for a name for the node and then
it will search for the node on the local network. Once the node is found it will register the node with the Wio server, then
configure the Access Point (AP) mode on the node. The node will then reboot and connect to the Wio server. Once the node
is connected to the Wio server it will be available for use.