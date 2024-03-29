# CLI Utility to add configurations as per git remote.

This is a CLI Utility which stores configurations according to git remotes. And applies your configurations according to current remotes. This commands automatically detects  github, bitbucket and gitlab and applies properties in the repository. 


## How to use

Download CLI according to your platform. And run

## All options

`gitconfig-provider --help`

This will list all available commands.

## Main commands to note

These are list of commands along with their use cases. To know the inputs and other details about commands please run `gitconfig-provider [command] --help`.

- `listProviders`: Lists all providers along with properties we have added using this CLI,

- `addProvider`: Adds new provider with empty properties

- `addConfig`: Adds properties as key value pair to given git provider.

- `deleteConfig`: Deletes property from given git provider.

- `reset`: Reset all git providers and properties to the original state.

- `apply`: Detects remote from current git repository, for remote detection it prefers `origin` as remote name, if no `origin` found, it prefers first remote from `git remote -v` command, and applies configuration according to detected remote.


## Test locally or on your broswer.

[![Open in Visual Studio Code](https://img.shields.io/static/v1?logo=visualstudiocode&label=&message=Open%20in%20Visual%20Studio%20Code&labelColor=2c2c32&color=007acc&logoColor=007acc)](https://open.vscode.dev/PrashamTrivedi/gitconfig-provider)

[Open In github.dev](https://github.dev/PrashamTrivedi/gitconfig-provider) 
