## Hooks

Hooks allow one to run code before/after various lifecycle events of the container.
The state of the container is passed to the hooks over stdin, so the hooks could get the information they need to do their work.

Hook paths are absolute and are executed from the host's filesystem.

### Pre-start

The pre-start hooks are called after the container process is spawned, but before the user supplied command is executed.
They are called after the container namespaces are created on Linux, so they provide an opportunity to customize the container.
In Linux, for e.g., the network namespace could be configured in this hook.

If a hook returns a non-zero exit code, then an error including the exit code and the stderr is returned to the caller and the container is torn down.

### Post-stop

The post-stop hooks are called after the container process is stopped. Cleanup or debugging could be performed in such a hook.
If a hook returns a non-zero exit code, then an error is logged and the remaining hooks are executed.

*Example*

```json
    "hooks" : {
        "prestart": [
            {
                "path": "/usr/bin/fix-mounts",
                "args": ["arg1", "arg2"],
                "env":  [ "key1=value1"]
            },
            {
                "path": "/usr/bin/setup-network"
            }
        ],
        "poststop": [
            {
                "path": "/usr/sbin/cleanup.sh",
                "args": ["-f"]
            }
        ]
    }
```

`path` is required for a hook. `args` and `env` are optional.

## Mount Configuration

Additional filesystems can be declared as "mounts", specified in the *mounts* array. The parameters are similar to the ones in Linux mount system call. [http://linux.die.net/man/2/mount](http://linux.die.net/man/2/mount)

* **type** (string, required) Linux, *filesystemtype* argument supported by the kernel are listed in */proc/filesystems* (e.g., "minix", "ext2", "ext3", "jfs", "xfs", "reiserfs", "msdos", "proc", "nfs", "iso9660"). Windows: ntfs
* **source** (string, required) a device name, but can also be a directory name or a dummy. Windows, the volume name that is the target of the mount point. \\?\Volume\{GUID}\ (on Windows source is called target)
* **destination** (string, required) where the source filesystem is mounted relative to the container rootfs.
* **options** (list of strings, optional) in the fstab format [https://wiki.archlinux.org/index.php/Fstab](https://wiki.archlinux.org/index.php/Fstab).

*Example (Linux)*

```json
"mounts": [
    {
        "type": "proc",
        "source": "proc",
        "destination": "/proc",
        "options": []
    },
    {
        "type": "tmpfs",
        "source": "tmpfs",
        "destination": "/dev",
        "options": ["nosuid","strictatime","mode=755","size=65536k"]
    },
    {
        "type": "devpts",
        "source": "devpts",
        "destination": "/dev/pts",
        "options": ["nosuid","noexec","newinstance","ptmxmode=0666","mode=0620","gid=5"]
    },
    {
        "type": "bind",
        "source": "/volumes/testing",
        "destination": "/data",
        "options": ["rbind","rw"]
    }
]
```

*Example (Windows)*

```json
"mounts": [
    {
        "type": "ntfs",
        "source": "\\\\?\\Volume\\{2eca078d-5cbc-43d3-aff8-7e8511f60d0e}\\",
        "destination": "C:\\Users\\crosbymichael\\My Fancy Mount Point\\",
        "options": []
    }
]
```

See links for details about [mountvol](http://ss64.com/nt/mountvol.html) and [SetVolumeMountPoint](https://msdn.microsoft.com/en-us/library/windows/desktop/aa365561(v=vs.85).aspx) in Windows.
