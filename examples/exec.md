# Executing additional processes inside an existing container

The [start operation][start] is described for creating new containers, but you can also use it to execute additional processes inside an existing container.
This sort of execution is similar to [`nsenter`][nsenter.1]—which joins existing namespaces and executes a new process inside them—but also allows you to easily join existing [control groups][cgroups], etc.

Because [start][] creates a new container, you'll need to give it a container ID, even if your new process is joining the existing container's sandbox with no additional isolation.

The [start][] configuration will look like:

## Example (Linux)

```json
{
    "ociVersion": "0.5.0",
    "platform": {
        "os": "linux",
        "arch": "amd64"
    },
    "process": {
        "terminal": true,
        "user": {
            "uid": 0,
            "gid": 0
        },
        "args": [
            "sh"
        ],
        "env": [
            "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
            "TERM=xterm"
        ],
        "cwd": "/",
        "capabilities": [
            "CAP_AUDIT_WRITE",
            "CAP_KILL"
        ]
    },
    "root": {
        "path": "rootfs"
    },
    "linux": {
         "cgroupsPath": "/myRuntime/myContainer",
         "namespaces": [
             {
                 "type": "user",
                 "path": "/proc/1234/ns/user"
             },
             {
                 "type": "pid",
                 "path": "/proc/1234/ns/pid"
             },
             {
                 "type": "network",
                 "path": "/proc/1234/ns/net"
             },
             {
                 "type": "ipc",
                 "path": "/proc/1234/ns/ipc"
             },
             {
                 "type": "uts",
                 "path": "/proc/1234/ns/uts"
             },
             {
                 "type": "mount",
                 "path": "/proc/1234/ns/mnt"
             }
         ]
    }
}
```

### Process

You can use whichever [**`process`**][process] settings you like.

### Root

The [**`root.path`**][root] value is arbitrary.
You must set a value because the field is currently [required][root], but [unless you are creating a new mount namespace the value is ignored][ignored-root].

### Control group and namespace paths

The [**`linux.cgroupsPath`**][cgroups] and [**`linux.namespaces[].path`**][namespaces] values can be extracted from [proc][proc.5]:

1. Get the [state JSON][state] for the target container using a [state query][state-query].
   The process ID for the target container is the **`pid`** field (e.g. `1234`).
2. Using a [proc][proc.5] mount for the appropriate [PID namespace][pid_namespaces.7] (possibly just `/proc`), find the `{pid}` entry for the target container (e.g. `/proc/1234`).
3. The namespace paths are in the [`ns`][proc.5] subdirectory (e.g. `/proc/1234/ns/user`, `/proc/1234/ns/pid`, …).
4. The cgroup paths are in the [`cgroup`][proc.5] file (e.g. `/proc/1234/cgroup`).
   Pick whichever path makes the most sense to you; containers created via OCI runtimes will [have the same path for all controllers][cgroups].

[cgroups]: ../config-linux.md#control-groups
[namespaces]: ../config-linux.md#namespaces
[process]: ../config.md#process-configuration
[root]: ../config.md#root-configuration
[start]: ../runtime.md#start
[state]: ../runtime.md#state
[state-query]: ../runtime.md#query-state

[ignored-root]: https://github.com/opencontainers/runtime-spec/pull/388#issuecomment-212188661

[nsenter.1]: http://man7.org/linux/man-pages/man1/nsenter.1.html
[proc.5]: http://man7.org/linux/man-pages/man5/proc.5.html
[pid_namespaces.7]: http://man7.org/linux/man-pages/man7/pid_namespaces.7.html
