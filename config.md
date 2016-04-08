# Container Configuration file

The container's top-level directory MUST contain a configuration file called `config.json`.
The canonical schema is defined in this document, but there is a JSON Schema in [`schema/schema.json`](schema/schema.json) and Go bindings in [`specs-go/config.go`](specs-go/config.go).

The configuration file contains metadata necessary to implement standard operations against the container.
This includes the process to run, environment variables to inject, sandboxing features to use, etc.

Below is a detailed description of each field defined in the configuration format.

## Specification version

* **`ociRuntimeVersion`** (string, required) must be in [SemVer v2.0.0](http://semver.org/spec/v2.0.0.html) format and specifies the version of the OpenContainer specification with which the bundle complies.
The OpenContainer spec follows semantic versioning and retains forward and backward compatibility within major versions.
For example, if an implementation is compliant with version 1.0.1 of the spec, it is compatible with the complete 1.x series.
NOTE that there is no guarantee for forward or backward compatibility for version 0.x.

*Example*

```json
    "ociRuntimeVersion": "0.1.0"
```

## Root Configuration

Each container has exactly one *root filesystem*, specified in the *root* object:

* **`path`** (string, required) Specifies the path to the root filesystem for the container, relative to the path where the manifest is. A directory MUST exist at the relative path declared by the field.
* **`readonly`** (bool, optional) If true then the root filesystem MUST be read-only inside the container. Defaults to false.

*Example*

```json
"root": {
    "path": "rootfs",
    "readonly": true
}
```

## Mounts

You can add array of mount points inside container as `mounts`.
The runtime MUST mount entries in the listed order.
The parameters are similar to the ones in [the Linux mount system call](http://man7.org/linux/man-pages/man2/mount.2.html).

* **`destination`** (string, required) Destination of mount point: path inside container.
* **`type`** (string, required) Linux, *filesystemtype* argument supported by the kernel are listed in */proc/filesystems* (e.g., "minix", "ext2", "ext3", "jfs", "xfs", "reiserfs", "msdos", "proc", "nfs", "iso9660"). Windows: ntfs
* **`source`** (string, required) a device name, but can also be a directory name or a dummy. Windows, the volume name that is the target of the mount point. \\?\Volume\{GUID}\ (on Windows source is called target)
* **`options`** (list of strings, optional) in the fstab format [https://wiki.archlinux.org/index.php/Fstab](https://wiki.archlinux.org/index.php/Fstab).

### Linux Example

```json
"mounts": [
    {
        "destination": "/tmp",
        "type": "tmpfs",
        "source": "tmpfs",
        "options": ["nosuid","strictatime","mode=755","size=65536k"]
    },
    {
        "destination": "/data",
        "type": "bind",
        "source": "/volumes/testing",
        "options": ["rbind","rw"]
    }
]
```

### Windows Example

```json
"mounts": [
    "myfancymountpoint": {
        "destination": "C:\\Users\\crosbymichael\\My Fancy Mount Point\\",
        "type": "ntfs",
        "source": "\\\\?\\Volume\\{2eca078d-5cbc-43d3-aff8-7e8511f60d0e}\\",
        "options": []
    }
]
```

See links for details about [mountvol](http://ss64.com/nt/mountvol.html) and [SetVolumeMountPoint](https://msdn.microsoft.com/en-us/library/windows/desktop/aa365561(v=vs.85).aspx) in Windows.


## Process configuration

* **`terminal`** (bool, optional) specifies whether you want a terminal attached to that process. Defaults to false.
* **`cwd`** (string, required) is the working directory that will be set for the executable. This value MUST be an absolute path.
* **`env`** (array of strings, optional) contains a list of variables that will be set in the process's environment prior to execution. Elements in the array are specified as Strings in the form "KEY=value". The left hand side must consist solely of letters, digits, and underscores `_` as outlined in [IEEE Std 1003.1-2001](http://pubs.opengroup.org/onlinepubs/009695399/basedefs/xbd_chap08.html).
* **`args`** (string, required) executable to launch and any flags as an array. The executable is the first element and must be available at the given path inside of the rootfs. If the executable path is not an absolute path then the search $PATH is interpreted to find the executable.

For Linux-based systems the process structure supports the following process specific fields:

* **`capabilities`** (array of strings, optional) capabilities is an array that specifies Linux capabilities that can be provided to the process inside the container.
Valid values are the strings for capabilities defined in [the man page](http://man7.org/linux/man-pages/man7/capabilities.7.html)
* **`rlimits`** (array of rlimits, optional) rlimits is an array of rlimits that allows setting resource limits for a process inside the container.
The kernel enforces the `soft` limit for a resource while the `hard` limit acts as a ceiling for that value that could be set by an unprivileged process.
Valid values for the 'type' field are the resources defined in [the man page](http://man7.org/linux/man-pages/man2/setrlimit.2.html).
* **`apparmorProfile`** (string, optional) apparmor profile specifies the name of the apparmor profile that will be used for the container.
For more information about Apparmor, see [Apparmor documentation](https://wiki.ubuntu.com/AppArmor)
* **`selinuxLabel`** (string, optional) SELinux process label specifies the label with which the processes in a container are run.
For more information about SELinux, see  [Selinux documentation](http://selinuxproject.org/page/Main_Page)
* **`noNewPrivileges`** (bool, optional) setting `noNewPrivileges` to true prevents the processes in the container from gaining additional privileges.
[The kernel doc](https://www.kernel.org/doc/Documentation/prctl/no_new_privs.txt) has more information on how this is achieved using a prctl system call.

### User

The user for the process is a platform-specific structure that allows specific control over which user the process runs as.

#### Linux User

For Linux-based systems the user structure has the following fields:

* **`uid`** (int, required) specifies the user id.
* **`gid`** (int, required) specifies the group id.
* **`additionalGids`** (array of ints, optional) specifies additional group ids to be added to the process.

_Note: symbolic name for uid and gid, such as uname and gname respectively, are left to upper levels to derive (i.e. `/etc/passwd` parsing, NSS, etc)_

*Example (Linux)*

```json
"process": {
    "terminal": true,
    "user": {
        "uid": 1,
        "gid": 1,
        "additionalGids": [5, 6]
    },
    "env": [
        "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
        "TERM=xterm"
    ],
    "cwd": "/root",
    "args": [
        "sh"
    ],
    "apparmorProfile": "acme_secure_profile",
    "selinuxLabel": "system_u:system_r:svirt_lxc_net_t:s0:c124,c675",
    "noNewPrivileges": true,
    "capabilities": [
        "CAP_AUDIT_WRITE",
        "CAP_KILL",
        "CAP_NET_BIND_SERVICE"
    ],
    "rlimits": [
        {
            "type": "RLIMIT_NOFILE",
            "hard": 1024,
            "soft": 1024
        }
    ]
}
```


## Hostname

* **`hostname`** (string, optional) as it is accessible to processes running inside.  On Linux, you can only set this if your bundle creates a new [UTS namespace][uts-namespace].

*Example*

```json
"hostname": "mrsdalloway"
```

## Platform-specific configuration

* **`os`** (string, required) specifies the operating system family this image must run on. Values for os must be in the list specified by the Go Language document for [`$GOOS`](https://golang.org/doc/install/source#environment).
* **`arch`** (string, required) specifies the instruction set for which the binaries in the image have been compiled. Values for arch must be in the list specified by the Go Language document for [`$GOARCH`](https://golang.org/doc/install/source#environment).

```json
"platform": {
    "os": "linux",
    "arch": "amd64"
}
```

Interpretation of the platform section of the JSON file is used to find which platform-specific sections may be available in the document.
For example, if `os` is set to `linux`, then a JSON object conforming to the [Linux-specific schema](config-linux.md) SHOULD be found at the key `linux` in the `config.json`.

## Hooks

Lifecycle hooks allow custom events for different points in a container's runtime.
Presently there are `Prestart`, `Poststart` and `Poststop`.

* [`Prestart`](#prestart) is a list of hooks to be run before the container process is executed
* [`Poststart`](#poststart) is a list of hooks to be run immediately after the container process is started
* [`Poststop`](#poststop) is a list of hooks to be run after the container process exits

Hooks allow one to run code before/after various lifecycle events of the container.
Hooks MUST be called in the listed order.
The state of the container is passed to the hooks over stdin, so the hooks could get the information they need to do their work.

Hook paths are absolute and are executed from the host's filesystem.

### Prestart

The pre-start hooks are called after the container process is spawned, but before the user supplied command is executed.
They are called after the container namespaces are created on Linux, so they provide an opportunity to customize the container.
In Linux, for e.g., the network namespace could be configured in this hook.

If a hook returns a non-zero exit code, then an error including the exit code and the stderr is returned to the caller and the container is torn down.

### Poststart

The post-start hooks are called after the user process is started.
For example this hook can notify user that real process is spawned.

If a hook returns a non-zero exit code, then an error is logged and the remaining hooks are executed.

### Poststop

The post-stop hooks are called after the container process is stopped.
Cleanup or debugging could be performed in such a hook.
If a hook returns a non-zero exit code, then an error is logged and the remaining hooks are executed.

*Example*

```json
    "hooks" : {
        "prestart": [
            {
                "path": "/usr/bin/fix-mounts",
                "args": ["fix-mounts", "arg1", "arg2"],
                "env":  [ "key1=value1"]
            },
            {
                "path": "/usr/bin/setup-network"
            }
        ],
        "poststart": [
            {
                "path": "/usr/bin/notify-start"
                "timeout": 5
            }
        ],
        "poststop": [
            {
                "path": "/usr/sbin/cleanup.sh",
                "args": ["cleanup.sh", "-f"]
            }
        ]
    }
```

`path` is required for a hook.
`args` and `env` are optional.
`timeout` is the number of seconds before aborting the hook.
The semantics are the same as `Path`, `Args` and `Env` in [golang Cmd](https://golang.org/pkg/os/exec/#Cmd).

## Annotations

Annotations are optional arbitrary non-identifying metadata that can be attached to containers.
This information may be large, may be structured or unstructured.
Annotations are key-value maps.

```json
"annotations": {
	"key1" : "value1",
	"key2" : "value2"
}
```

## Configuration Schema Example

Here is a full example `config.json` for reference.

```json
{
    "ociRuntimeVersion": "0.3.0",
    "platform": {
        "os": "linux",
        "arch": "amd64"
    },
    "process": {
        "terminal": true,
        "user": {
            "uid": 1,
            "gid": 1,
            "additionalGids": [
                5,
                6
            ]
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
            "CAP_KILL",
            "CAP_NET_BIND_SERVICE"
        ],
        "rlimits": [
            {
                "type": "RLIMIT_NOFILE",
                "hard": 1024,
                "soft": 1024
            }
        ],
        "apparmorProfile": "",
        "selinuxLabel": ""
    },
    "root": {
        "path": "rootfs",
        "readonly": true
    },
    "hostname": "slartibartfast",
    "mounts": [
        {
            "destination": "/proc",
            "type": "proc",
            "source": "proc"
        },
        {
            "destination": "/dev",
            "type": "tmpfs",
            "source": "tmpfs",
            "options": [
                "nosuid",
                "strictatime",
                "mode=755",
                "size=65536k"
            ]
        },
        {
            "destination": "/dev/pts",
            "type": "devpts",
            "source": "devpts",
            "options": [
                "nosuid",
                "noexec",
                "newinstance",
                "ptmxmode=0666",
                "mode=0620",
                "gid=5"
            ]
        },
        {
            "destination": "/dev/shm",
            "type": "tmpfs",
            "source": "shm",
            "options": [
                "nosuid",
                "noexec",
                "nodev",
                "mode=1777",
                "size=65536k"
            ]
        },
        {
            "destination": "/dev/mqueue",
            "type": "mqueue",
            "source": "mqueue",
            "options": [
                "nosuid",
                "noexec",
                "nodev"
            ]
        },
        {
            "destination": "/sys",
            "type": "sysfs",
            "source": "sysfs",
            "options": [
                "nosuid",
                "noexec",
                "nodev"
            ]
        },
        {
            "destination": "/sys/fs/cgroup",
            "type": "cgroup",
            "source": "cgroup",
            "options": [
                "nosuid",
                "noexec",
                "nodev",
                "relatime",
                "ro"
            ]
        }
    ],
    "hooks": {
        "prestart": [
            {
                "path": "/usr/bin/uptime",
                "args": [
                    "/usr/bin/uptime"
                ],
                "env": []
            }
        ]
    },
    "linux": {
        "resources": {
            "devices": [
                {
                    "allow": false,
                    "access": "rwm"
                }
            ]
        },
        "namespaces": [
            {
                "type": "pid"
            },
            {
                "type": "network"
            },
            {
                "type": "ipc"
            },
            {
                "type": "uts"
            },
            {
                "type": "mount"
            }
        ]
    }
}
```


[uts-namespace]: http://man7.org/linux/man-pages/man7/namespaces.7.html
