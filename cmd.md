# Runtime Command Line Interface

This document specifies the command line syntax for runtime operations defined in the [Runtime and Lifecycle](runtime.md) specification.
It is STRONGLY RECOMMENDED that implementations adhere to the command line definition defined below to ensure interoperability.
This document will not duplicate the information/semantics specified in the [Runtime and Lifecycle](runtime.md) document, rather it just focuses on syntax of the command line.

## Convention

For the purpose of this document, the following syntactical conventions are used:

* `[ ... ]` denotes an optional field
* `< ... >` denotes a substitution field. 
The word(s), or phrase, within the `<>` describes the information to be passed-in.

## General Format and Behavior

The general format of all commands MUST be:
```
runtime [global-options] action [action-specific-options] [arguments]
```

Unknown options (global and action specific) MUST generate an error and exit with a non-zero exit code, without changing the state of the environment.

Upon successful running of an action, the exit code MUST be zero.

If there is an error during the running of an action, then:
* the exit code MUST be non-zero
* any error text MUST be displayed on stderr
* the state of the environment SHOULD be the same as if the action was never attempted (modulo any possible trivial ancillary changes such as logging)

### Global Options

Global options are ones that apply to all actions.
This specification doesn't define any global options.
Implementation MAY define their own.

### Action Specific Options

All actions MUST support the `--help` action-specific-option, which:
* MUST display some help text to stdout
* MUST NOT perform the `action` itself
* MUST have an exit code of zero if the help text is successfully displayed

Implementations MAY define their own action-specific options.

## Actions

This section defines the actions defined by this specification.

### State

Format: `runtime [global-options] state [options] <container-id>`

Options: None

This action MUST display the state of the specific container to stdout.
Unless otherwise specified by the user, the format of the state MUST be in JSON.

See [Query State](runtime.md#query-state).

### Create

Format: `runtime [global-options] create [options] <container-id>`

Options:
* `-b <dir>`,  `--bundle <dir>` The path to the root of the bundle directory. If not specified the default value MUST be the current working directory.
* `--console <path>` The PTY slave path for the newly created container.
* `--pid-file <path>` The file path into which the process ID is written. If not specified then a file MUST NOT be created.

See [Create](runtime.md#create).

### Start

Format: `runtime [global-options] start [options] <container-id>`

Options: None

See [Start](runtime.md#start).

### Kill

Format: `runtime [global-options] kill [options] <container-id> <signal>`

Options: None

The `signal` MUST either be the signal's numerical value (e.g. `15`) or the signal's name (e.g. `SIGTERM`).

See [Kill](runtime.md#kill).

### Delete

Format: `runtime [global-options] delete [options] <container-id>`

Options: None

See [Delete](runtime.md#delete).
