# goperator

goperator is a CLI tool to manage your AWS infrastructure.

## Installation

You can install the lastest version of the source code and running:

`make install`

`goperator` binary will be placed in your `/usr/local/bin` folder

## Tags

To run goperator is very important that all your EC2 instances have the following tags:

- *component*
- *environment*

all the other tags such as _service_ and _working_version_ are optional.

PEM files to allow ssh access to those instances will have the format `<environment>-<component>.pem` and should be stored within the `ssh-keys` folder.

## Usage

Set some environment variables if you want to override the current ones:

```Bash
export AWS_REGION=us-west-2
```