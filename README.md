# gitlab-automation : Automation of Gitlab management

gitlab-automation is a command line tool to automate operations on Gitlab. Most of automation process required Gitlab token.

## Usage

The usage of gitlab-automation is based on the subcommand.

For example: ```gitlab-automation help``` is returning the list of available commands.

For example: ```gitlab-automation permissions``` is the command to audit the permissions at a group level.

Each command could require additional parameters.

For example: ```gitlab-automation permissions``` requires the token and the group identifier as a minimum.

By requesting the command without arguments, the command is self-explanatory:
```
Usage of permissions:
  -gid string
    	Gitlab group id to check
  -token string
    	Gitlab token
```

## Build

There is a Makefile to build. By default, the resulting binary is being placed in the ```build``` directory.

Processes available by the Makefile are show by default if no argument is passed to the make command.