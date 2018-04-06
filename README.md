# Open Container Initiative Runtime Specification

The [Open Container Initiative][oci] develops specifications for standards on Operating System process and application containers.

The specification can be found [here](spec.md).

## Table of Contents

Additional documentation about how this group operates:

- [Code of Conduct][code-of-conduct]
- [Style and Conventions](style.md)
- [Implementations](implementations.md)
- [Releases](RELEASES.md)
- [project](project.md)
- [charter][charter]

## Use Cases

To provide context for users the following section gives example use cases for each part of the spec.

### Application Bundle Builders

Application bundle builders can create a [bundle](bundle.md) directory that includes all of the files required for launching an application as a container.
The bundle contains an OCI [configuration file](config.md) where the builder can specify host-independent details such as [which executable to launch](config.md#process) and host-specific settings such as [mount](config.md#mounts) locations, [hook](config.md#posix-platform-hooks) paths, Linux [namespaces](config-linux.md#namespaces) and [cgroups](config-linux.md#control-groups).
Because the configuration includes host-specific settings, application bundle directories copied between two hosts may require configuration adjustments.

### Hook Developers

[Hook](config.md#posix-platform-hooks) developers can extend the functionality of an OCI-compliant runtime by hooking into a container's lifecycle with an external application.
Example use cases include sophisticated network configuration, volume garbage collection, etc.

### Runtime Developers

Runtime developers can build runtime implementations that run OCI-compliant bundles and container configuration, containing low-level OS and host-specific details, on a particular platform.

## Contributing

See [our contribution documentation](CONTRIBUTING.md).

[charter]: https://www.opencontainers.org/about/governance
[code-of-conduct]: https://github.com/opencontainers/tob/blob/d2f9d68c1332870e40693fe077d311e0742bc73d/code-of-conduct.md
[oci]: https://www.opencontainers.org
