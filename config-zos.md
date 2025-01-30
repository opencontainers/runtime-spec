# <a name="ZOSContainerConfiguration" />z/OS Container Configuration

This document describes the schema for the [z/OS-specific section](config.md#platform-specific-configuration) of the [container configuration](config.md).
The z/OS container specification uses z/OS UNIX kernel features like namespaces and filesystem jails to fulfill the spec.

Applications expecting a z/OS environment will very likely expect these file paths to be set up correctly.

The following filesystems SHOULD be made available in each container's filesystem:

| Path     | Type   |
| -------- | ------ |
| /proc    | [proc][] |

## <a name="configZOSNamespaces" />Namespaces

A namespace wraps a global system resource in an abstraction that makes it appear to the processes within the namespace that they have their own isolated instance of the global resource.
Changes to the global resource are visible to other processes that are members of the namespace, but are invisible to other processes.
For more information, see https://www.ibm.com/docs/zos/latest?topic=planning-namespaces-zos-unix.

Namespaces are specified as an array of entries inside the `namespaces` root field.
The following parameters can be specified to set up namespaces:

* **`type`** *(string, REQUIRED)* - namespace type. The following namespace types SHOULD be supported:
    * **`pid`** processes inside the container will only be able to see other processes inside the same container or inside the same pid namespace.
    * **`mount`** the container will have an isolated mount table.
    * **`ipc`** processes inside the container will only be able to communicate to other processes inside the same container via system level IPC.
    * **`uts`** the container will be able to have its own hostname and domain name.
* **`path`** *(string, OPTIONAL)* - namespace file.
    This value MUST be an absolute path in the [runtime mount namespace](glossary.md#runtime-namespace).
    The runtime MUST place the container process in the namespace associated with that `path`.
    The runtime MUST [generate an error](runtime.md#errors) if `path` is not associated with a namespace of type `type`.

    If `path` is not specified, the runtime MUST create a new [container namespace](glossary.md#container-namespace) of type `type`.

If a namespace type is not specified in the `namespaces` array, the container MUST inherit the [runtime namespace](glossary.md#runtime-namespace) of that type.
If a `namespaces` field contains duplicated namespaces with same `type`, the runtime MUST [generate an error](runtime.md#errors).

### Example

```json
"namespaces": [
    {
        "type": "pid",
        "path": "/proc/1234/ns/pid"
    },
    {
        "type": "mount"
    },
    {
        "type": "ipc"
    },
    {
        "type": "uts"
    }
]
```
