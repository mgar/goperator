
## Installation

You can install the lastest version of the source code and running:

`make install`

`goperator` binary will be placed in your `/usr/local/bin` folder

## Tags

If you want to list your instances or access them using SSH, instances must have the following tags:

- *component*
- *environment*

All the other tags such as _service_ and _working_version_ are optional.

PEM files to allow ssh access to those instances will have the format `<environment>-<component>.pem` and should be stored within the `ssh-keys` folder.

## Usage

Set some environment variables if you want to override the current ones:

```Bash
export AWS_REGION=us-west-2
```

- `start` Start one or many instances which are currently stopped
- `stop` Stop one or many instances which are currently started
- `terminate` Terminate one or many instances
- `list` List instances based on its tags [environment] and [component]
- `command` Execute a command on one or many instances