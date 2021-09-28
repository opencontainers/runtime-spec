# <a name="implementations" />Implementations

The following sections link to associated projects, some of which are maintained by the OCI and some of which are maintained by external organizations.
If you know of any associated projects that are not listed here, please file a pull request adding a link to that project.

## <a name="implementationsRuntimeContainer" />Runtime (Container)

* [alibaba/inclavare-containers][rune] - Enclave OCI runtime for confidential computing
* [containers/crun][crun] - Runtime implementation in C
* [containers/youki][youki] - Runtime implementation in Rust
* [opencontainers/runc][runc] - Reference implementation of OCI runtime
* [projectatomic/bwrap-oci][bwrap-oci] - Convert the OCI spec file to a command line for [bubblewrap][bubblewrap]

## <a name="implementationsRuntimeVirtualMachine" />Runtime (Virtual Machine)

* [clearcontainers/runtime][cc-runtime] - Hypervisor-based OCI runtime utilising [virtcontainers][virtcontainers] by Intel®.
* [google/gvisor][gvisor] - gVisor is a user-space kernel, contains runsc to run sandboxed containers.
* [hyperhq/runv][runv] - Hypervisor-based runtime for OCI
* [kata-containers/runtime][kata-runtime] - Hypervisor-based OCI runtime combining technology from [clearcontainers/runtime][cc-runtime] and [hyperhq/runv][runv].

## <a name="implementationsTestingTools" />Testing & Tools

* [huawei-openlab/oct][oct] - Open Container Testing framework for OCI configuration and runtime
* [kunalkushwaha/octool][octool] - A config linter and validator.
* [opencontainers/runtime-tools][runtime-tools] - A config generator and runtime/bundle testing framework.

[bubblewrap]: https://github.com/projectatomic/bubblewrap
[bwrap-oci]: https://github.com/projectatomic/bwrap-oci
[cc-runtime]: https://github.com/clearcontainers/runtime
[crun]: https://github.com/containers/crun
[gvisor]: https://github.com/google/gvisor
[kata-runtime]: https://github.com/kata-containers/runtime
[oct]: https://github.com/huawei-openlab/oct
[octool]: https://github.com/kunalkushwaha/octool
[runc]: https://github.com/opencontainers/runc
[rune]: https://github.com/alibaba/inclavare-containers
[runtime-tools]: https://github.com/opencontainers/runtime-tools
[runv]: https://github.com/hyperhq/runv
[virtcontainers]: https://github.com/containers/virtcontainers
[youki]: https://github.com/containers/youki