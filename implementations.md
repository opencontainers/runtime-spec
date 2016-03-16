# Implementations

The following sections link to associated projects, some of which are maintained by the OCI and some of which are maintained by external organizations.
If you know of any associated projects that are not listed here, please file a pull request adding a link to that project.

## Runtime (Container)

* [opencontainers/runc](https://github.com/opencontainers/runc) - Reference implementation of OCI runtime

## Runtime (Virtual Machine)

* [hyperhq/runv](https://github.com/hyperhq/runv) - Hypervisor-based runtime for OCI

## Bundle authoring

* [kunalkushwaha/octool](https://github.com/kunalkushwaha/octool) - A config linter and validator.
* [mrunalp/ocitools](https://github.com/mrunalp/ocitools) - A config generator.

## Testing

* [huawei-openlab/oct](https://github.com/huawei-openlab/oct) - Open Container Testing framework for OCI configuration and runtime

## Hooks

For addressing container issues that are outside the scope of this specification, the following hooks may be useful.

### Linux

#### Networking

* [CNI][], the Container Network Interface.
  The [`exec-plugins.sh`][cni-exec-plugins.sh] script can be used in both [prestart][] (with `add`) and [poststop][] (with `del`) hooks insert and remove a container from networks [configured][cni-conf] in a host-side directory like `/etc/cni/net.d`.

[prestart]: runtime-config.md#prestart
[poststop]: runtime-config.md#poststop

[CNI]: https://github.com/appc/cni
[cni-exec-plugins.sh]: https://github.com/appc/cni/blob/v0.1.0/scripts/exec-plugins.sh
[cni-conf]: https://github.com/appc/cni/blob/v0.1.0/SPEC.md
