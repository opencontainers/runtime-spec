# JSON schema

## Overview

This directory contains the [JSON Schema](https://json-schema.org) for validating JSON covered by this specification.

The layout of the files is as follows:

* [config-schema.json](config-schema.json) - the primary entrypoint for the [configuration](../config.md) schema
* [config-linux.json](config-linux.json) - the [Linux-specific configuration sub-structure](../config-linux.md)
* [config-solaris.json](config-solaris.json) - the [Solaris-specific configuration sub-structure](../config-solaris.md)
* [config-windows.json](config-windows.json) - the [Windows-specific configuration sub-structure](../config-windows.md)
* [config-freebsd.json](config-freebsd.json) - the [FreeBSD-specific configuration sub-structure](../config-freebsd.md)
* [state-schema.json](state-schema.json) - the primary entrypoint for the [state JSON](../runtime.md#state) schema
* [defs.json](defs.json) - definitions for general types
* [defs-linux.json](defs-linux.json) - definitions for Linux-specific types
* [defs-windows.json](defs-windows.json) - definitions for Windows-specific types
* [validate.go](validate.go) - validation utility source code


## Utility

There is also included a simple utility for facilitating validation.
To build it:

```bash
go get github.com/xeipuuv/gojsonschema
go build ./validate.go
```

Or you can just use make command to create the utility:

```bash
make validate
```

Then use it like:

```bash
./validate config-schema.json <yourpath>/config.json
```

Or like:

```bash
./validate https://raw.githubusercontent.com/opencontainers/runtime-spec/<runtime-spec-version>/schema/config-schema.json <yourpath>/config.json
```
