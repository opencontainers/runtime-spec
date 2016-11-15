# JSON schema

## Overview

This directory contains the [JSON Schema](http://json-schema.org/) for validating JSON covered by this specification.

The layout of the files is as follows:

* [config-schema.json](config-schema.json) - the primary entrypoint for the [configuration](../config.asc) schema
* [config-linux.json](config-linux.json) - the [Linux-specific configuration sub-structure](../config-linux.asc)
* [config-solaris.json](config-solaris.json) - the [Solaris-specific configuration sub-structure](../config-solaris.asc)
* [config-windows.json](config-windows.json) - the [Windows-specific configuration sub-structure](../config-windows.asc)
* [state-schema.json](state-schema.json) - the primary entrypoint for the [state JSON](../runtime.asc#state) schema
* [defs.json](defs.json) - definitions for general types
* [defs-linux.json](defs-linux.json) - definitions for Linux-specific types
* [validate.go](validate.go) - validation utility source code


## Utility

There is also included a simple utility for facilitating validation.
To build it:

```bash
export GOPATH=`mktemp -d`
go get -d ./...
go build ./validate.go
rm -rf $GOPATH
```

Or you can just use make command to create the utility:

```bash
make validate
```

Then use it like:

```bash
./validate config-schema.json <yourpath>/config.json
```
