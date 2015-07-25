# Bundle Container Format

This section defines a format for encoding a container as a *bundle* - a directory organized in a certain way, and containing all the necessary data and metadata for any compliant runtime to perform all standard operations against it. See also [OS X application bundles](http://en.wikipedia.org/wiki/Bundle_%28OS_X%29) for a similar use of the term *bundle*.

The format does not define distribution. In other words, it only specifies how a container must be stored on a local filesystem, for consumption by a runtime. It does not specify how to transfer a container between computers, how to discover containers, or assign names or versions to them. Any distribution method capable of preserving the original layout of a container, as specified here, is considered compliant.

A standard container bundle is made of the following 3 parts:

- A top-level directory holding everything else
- One or more content directories
- A configuration file

# Directory layout

A Standard Container bundle is a directory containing all the content needed to load and run a container. This includes its configuration file(s) and content directories. The main property of this directory layout is that it can be moved as a unit to another machine and run the same container.

*Example*

```
/
|-- config.json
`-- rootfs
```

## Configuration

The config file's syntax and semantics are described in [this specification](config.md).  By default, containers will use `config.json` in the bundle root:

```
/
|-- config.json
`-- rootfs
```

However, sometimes you need a more flexible configuration than you can get from a single static file.  Runtime's should use the following logic to select which config file to use:

1. If a `config.json` file exists (or you were passed a file path), use that.
2. If a `config` directory exists (or you were passed a directory path), walk it looking for the best platform match and use that.

### Alternative configurations

Runtimes like [runC][] allow you to specify a different config file explicitly, so you may find it convenient to place additional config files somewhere in your bundle.  For example, with a bundle like:

```
/
|-- config.json
|-- config-shell.json
`-- rootfs
```

Then you could use `runc` to launch your application, and `runc config-shell.json` to launch a shell in a similar container environment for poking around.

### Multiple platforms

In some situations, it's convenient to have a single bundle for multiple platforms.  For example, a Python bundle template could be written with the Python interpreter in `rootfs` with `process` configs designed to launch a `app/main.py`.  This template could be shared by many developers, who create bundles by dropping their application into the `app` directory without having to worry about platform idiosyncrasies.

```
/
|-- config
|   |-- linux.json
|   `-- windows.json
|-- rootfs
|   |-- linux
|   `-- windows
`-- app
```

When the runtime loads a config from a directory, it walks the directory recursively to find all `*.json` files, and checks those files for compatible [`version`s][version].  Of the compatible files, it chooses the config with:

1. The best [`platform`][platform] match for the runtime system, breaking ties with
2. The newest [`version`][version].

Each config file in the directory should stand alone, so there may be some information that is duplicated among several config files.  Bundle authors who want a [DRYer][DRY] system are free to use an independent tool to generate the config files.

## Content

One or more *content directories* and *auxiliary files* may be adjacent to the configuration file or configuration directory. This must include at least the root filesystem (referenced in the configuration file by the [`root` field][root]) and may include other related content (signatures, other configs, etc.). The interpretation of these resources may be specified in the configuration (e.g. the [`root` field][root]) or they may be runtime extensions. The names of the non-config directories are arbitrary, but users should consider using conventional names.

[runC]: https://github.com/opencontainers/runc
[version]: ./config.md#manifest-version
[platform]: ./config.md#platform-specific-configuration
[DRY]: https://en.wikipedia.org/wiki/Don%27t_repeat_yourself
[root]: ./config.md#root-configuration
