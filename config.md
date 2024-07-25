# <a name="configuration" />Configuration

This configuration file contains metadata necessary to implement [standard operations](runtime.md#operations) against the container.
This includes the process to run, environment variables to inject, sandboxing features to use, etc.

The canonical schema is defined in this document, but there is a JSON Schema in [`schema/config-schema.json`](schema/config-schema.json) and Go bindings in [`specs-go/config.go`](specs-go/config.go).
[Platform](spec.md#platforms)-specific configuration schema are defined in the [platform-specific documents](#platform-specific-configuration) linked below.
For properties that are only defined for some [platforms](spec.md#platforms), the Go property has a `platform` tag listing those protocols (e.g. `platform:"linux,solaris"`).

Below is a detailed description of each field defined in the configuration format and valid values are specified.
Platform-specific fields are identified as such.
For all platform-specific configuration values, the scope defined below in the [Platform-specific configuration](#platform-specific-configuration) section applies.


## <a name="configSpecificationVersion" />Specification version

* **`ociVersion`** (string, REQUIRED) MUST be in [SemVer v2.0.0][semver-v2.0.0] format and specifies the version of the Open Container Initiative Runtime Specification with which the bundle complies.
    The Open Container Initiative Runtime Specification follows semantic versioning and retains forward and backward compatibility within major versions.
    For example, if a configuration is compliant with version 1.1 of this specification, it is compatible with all runtimes that support any 1.1 or later release of this specification, but is not compatible with a runtime that supports 1.0 and not 1.1.

### Example

```json
"ociVersion": "0.1.0"
```

## <a name="configRoot" />Root

**`root`** (object, OPTIONAL) specifies the container's root filesystem.
On Windows, for Windows Server Containers, this field is REQUIRED.
For [Hyper-V Containers](config-windows.md#hyperv), this field MUST NOT be set.

On all other platforms, this field is REQUIRED.

* **`path`** (string, REQUIRED) Specifies the path to the root filesystem for the container.
    * On Windows, `path` MUST be a [volume GUID path][naming-a-volume].
    * On POSIX platforms, `path` is either an absolute path or a relative path to the bundle.
        For example, with a bundle at `/to/bundle` and a root filesystem at `/to/bundle/rootfs`, the `path` value can be either `/to/bundle/rootfs` or `rootfs`.
        The value SHOULD be the conventional `rootfs`.

    A directory MUST exist at the path declared by the field.

* **`readonly`** (bool, OPTIONAL) If true then the root filesystem MUST be read-only inside the container, defaults to false.
    * On Windows, this field MUST be omitted or false.

### Example (POSIX platforms)

```json
"root": {
    "path": "rootfs",
    "readonly": true
}
```

### Example (Windows)

```json
"root": {
    "path": "\\\\?\\Volume{ec84d99e-3f02-11e7-ac6c-00155d7682cf}\\"
}
```

## <a name="configMounts" />Mounts

**`mounts`** (array of objects, OPTIONAL) specifies additional mounts beyond [`root`](#root).
The runtime MUST mount entries in the listed order.
For Linux, the parameters are as documented in [mount(2)][mount.2] system call man page.
For Solaris, the mount entry corresponds to the 'fs' resource in the [zonecfg(1M)][zonecfg.1m] man page.

* **`destination`** (string, REQUIRED) Destination of mount point: path inside container.
    * Linux: This value SHOULD be an absolute path.
      For compatibility with old tools and configurations, it MAY be a relative path, in which case it MUST be interpreted as relative to "/".
      Relative paths are **deprecated**.
    * Windows: This value MUST be an absolute path.
      One mount destination MUST NOT be nested within another mount (e.g., c:\\foo and c:\\foo\\bar).
    * Solaris: This value MUST be an absolute path.
      Corresponds to "dir" of the fs resource in [zonecfg(1M)][zonecfg.1m].
    * For all other platforms: This value MUST be an absolute path.
* **`source`** (string, OPTIONAL) A device name, but can also be a file or directory name for bind mounts or a dummy.
    Path values for bind mounts are either absolute or relative to the bundle.
    A mount is a bind mount if it has either `bind` or `rbind` in the options.
    * Windows: a local directory on the filesystem of the container host. UNC paths and mapped drives are not supported.
    * Solaris: corresponds to "special" of the fs resource in [zonecfg(1M)][zonecfg.1m].
* **`options`** (array of strings, OPTIONAL) Mount options of the filesystem to be used.
    * Linux: See [Linux mount options](#configLinuxMountOptions) below.
    * Solaris: corresponds to "options" of the fs resource in [zonecfg(1M)][zonecfg.1m].
    * Windows: runtimes MUST support `ro`, mounting the filesystem read-only when `ro` is given.

### <a name="configLinuxMountOptions" />Linux mount options

Runtimes MUST/SHOULD/MAY implement the following option strings for Linux:

 Option name      | Requirement | Description
------------------|-------------|-----------------------------------------------------
 `async`          | MUST        | [^1]
 `atime`          | MUST        | [^1]
 `bind`           | MUST        | Bind mount [^2]
 `defaults`       | MUST        | [^1]
 `dev`            | MUST        | [^1]
 `diratime`       | MUST        | [^1]
 `dirsync`        | MUST        | [^1]
 `exec`           | MUST        | [^1]
 `iversion`       | MUST        | [^1]
 `lazytime`       | MUST        | [^1]
 `loud`           | MUST        | [^1]
 `mand`           | MAY         | [^1] (Deprecated in kernel 5.15, util-linux 2.38)
 `noatime`        | MUST        | [^1]
 `nodev`          | MUST        | [^1]
 `nodiratime`     | MUST        | [^1]
 `noexec`         | MUST        | [^1]
 `noiversion`     | MUST        | [^1]
 `nolazytime`     | MUST        | [^1]
 `nomand`         | MAY         | [^1]
 `norelatime`     | MUST        | [^1]
 `nostrictatime`  | MUST        | [^1]
 `nosuid`         | MUST        | [^1]
 `nosymfollow`    | SHOULD      | [^1] (Introduced in kernel 5.10, util-linux 2.38)
 `private`        | MUST        | Bind mount propagation [^2]
 `ratime`         | SHOULD      | Recursive `atime` [^3]
 `rbind`          | MUST        | Recursive bind mount [^2]
 `rdev`           | SHOULD      | Recursive `dev` [^3]
 `rdiratime`      | SHOULD      | Recursive `diratime` [^3]
 `relatime`       | MUST        | [^1]
 `remount`        | MUST        | [^1]
 `rexec`          | SHOULD      | Recursive `dev` [^3]
 `rnoatime`       | SHOULD      | Recursive `noatime` [^3]
 `rnodiratime`    | SHOULD      | Recursive `nodiratime` [^3]
 `rnoexec`        | SHOULD      | Recursive `noexec` [^3]
 `rnorelatime`    | SHOULD      | Recursive `norelatime` [^3]
 `rnostrictatime` | SHOULD      | Recursive `nostrictatime` [^3]
 `rnosuid`        | SHOULD      | Recursive `nosuid` [^3]
 `rnosymfollow`   | SHOULD      | Recursive `nosymfollow` [^3]
 `ro`             | MUST        | [^1]
 `rprivate`       | MUST        | Bind mount propagation [^2]
 `rrelatime  `    | SHOULD      | Recursive `relatime` [^3]
 `rro`            | SHOULD      | Recursive `ro` [^3]
 `rrw`            | SHOULD      | Recursive `rw` [^3]
 `rshared`        | MUST        | Bind mount propagation [^2]
 `rslave`         | MUST        | Bind mount propagation [^2]
 `rstrictatime`   | SHOULD      | Recursive `strictatime` [^3]
 `rsuid`          | SHOULD      | Recursive `suid` [^3]
 `rsymfollow`     | SHOULD      | Recursive `symfollow` [^3]
 `runbindable`    | MUST        | Bind mount propagation [^2]
 `rw`             | MUST        | [^1]
 `shared`         | MUST        | [^1]
 `silent`         | MUST        | [^1]
 `slave`          | MUST        | Bind mount propagation [^2]
 `strictatime`    | MUST        | [^1]
 `suid`           | MUST        | [^1]
 `symfollow`      | SHOULD      | Opposite of `nosymfollow`
 `sync`           | MUST        | [^1]
 `tmpcopyup`      | MAY         | copy up the contents to a tmpfs
 `unbindable`     | MUST        | Bind mount propagation [^2]
 `idmap`          | SHOULD      | Indicates that the mount MUST have an idmapping applied. This option SHOULD NOT be passed to the underlying [`mount(2)`][mount.2] call. If `uidMappings` or `gidMappings` are specified for the mount, the runtime MUST use those values for the mount's mapping. If they are not specified, the runtime MAY use the container's user namespace mapping, otherwise an [error MUST be returned](runtime.md#errors).  If there are no `uidMappings` and `gidMappings` specified and the container isn't using user namespaces, an [error MUST be returned](runtime.md#errors). This SHOULD be implemented using [`mount_setattr(MOUNT_ATTR_IDMAP)`][mount_setattr.2], available since Linux 5.12.
 `ridmap`         | SHOULD      | Indicates that the mount MUST have an idmapping applied, and the mapping is applied recursively [^3]. This option SHOULD NOT be passed to the underlying [`mount(2)`][mount.2] call. If `uidMappings` or `gidMappings` are specified for the mount, the runtime MUST use those values for the mount's mapping. If they are not specified, the runtime MAY use the container's user namespace mapping, otherwise an [error MUST be returned](runtime.md#errors).  If there are no `uidMappings` and `gidMappings` specified and the container isn't using user namespaces, an [error MUST be returned](runtime.md#errors). This SHOULD be implemented using [`mount_setattr(MOUNT_ATTR_IDMAP)`][mount_setattr.2], available since Linux 5.12.

[^1]: Corresponds to [`mount(8)` (filesystem-independent)][mount.8-filesystem-independent].
[^2]: Corresponds to [bind mounts and shared subtrees][mount-bind].
[^3]: These `AT_RECURSIVE` options need kernel 5.12 or later. See [`mount_setattr(2)`][mount_setattr.2]

The "MUST" options correspond to [`mount(8)`][mount.8].

Runtimes MAY also implement custom option strings that are not listed in the table above.
If a custom option string is already recognized by [`mount(8)`][mount.8], the runtime SHOULD follow the behavior of [`mount(8)`][mount.8].

Runtimes SHOULD treat unknown options as [filesystem-specific ones][mount.8-filesystem-specific])
and pass those as a comma-separated string to the fifth (`const void *data`) argument of [`mount(2)`][mount.2].

### Example (Windows)

```json
"mounts": [
    {
        "destination": "C:\\folder-inside-container",
        "source": "C:\\folder-on-host",
        "options": ["ro"]
    }
]
```

### <a name="configPOSIXMounts" />POSIX-platform Mounts

For POSIX platforms the `mounts` structure has the following fields:

* **`type`** (string, OPTIONAL) The type of the filesystem to be mounted.
    * Linux: filesystem types supported by the kernel as listed in */proc/filesystems* (e.g., "minix", "ext2", "ext3", "jfs", "xfs", "reiserfs", "msdos", "proc", "nfs", "iso9660"). For bind mounts (when `options` include either `bind` or `rbind`), the type is a dummy, often "none" (not listed in */proc/filesystems*).
    * Solaris: corresponds to "type" of the fs resource in [zonecfg(1M)][zonecfg.1m].
* **`uidMappings`** (array of type LinuxIDMapping, OPTIONAL) The mapping to convert UIDs from the source file system to the destination mount point.
  This SHOULD be implemented using [`mount_setattr(MOUNT_ATTR_IDMAP)`][mount_setattr.2], available since Linux 5.12.
  If specified, the `options` field of the `mounts` structure SHOULD contain either `idmap` or `ridmap` to specify whether the mapping should be applied recursively for `rbind` mounts, as well as to ensure that older runtimes will not silently ignore this field.
  The format is the same as [user namespace mappings](config-linux.md#user-namespace-mappings).
  If specified, it MUST be specified along with `gidMappings`.
* **`gidMappings`** (array of type LinuxIDMapping, OPTIONAL) The mapping to convert GIDs from the source file system to the destination mount point.
  This SHOULD be implemented using [`mount_setattr(MOUNT_ATTR_IDMAP)`][mount_setattr.2], available since Linux 5.12.
  If specified, the `options` field of the `mounts` structure SHOULD contain either `idmap` or `ridmap` to specify whether the mapping should be applied recursively for `rbind` mounts, as well as to ensure that older runtimes will not silently ignore this field.
  For more details see `uidMappings`.
  If specified, it MUST be specified along with `uidMappings`.


### Example (Linux)

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
        "type": "none",
        "source": "/volumes/testing",
        "options": ["rbind","rw"]
    }
]
```

### Example (Solaris)

```json
"mounts": [
    {
        "destination": "/opt/local",
        "type": "lofs",
        "source": "/usr/local",
        "options": ["ro","nodevices"]
    },
    {
        "destination": "/opt/sfw",
        "type": "lofs",
        "source": "/opt/sfw"
    }
]
```

## <a name="configProcess" />Process

**`process`** (object, OPTIONAL) specifies the container process.
This property is REQUIRED when [`start`](runtime.md#start) is called.

* **`terminal`** (bool, OPTIONAL) specifies whether a terminal is attached to the process, defaults to false.
    As an example, if set to true on Linux a pseudoterminal pair is allocated for the process and the pseudoterminal pty is duplicated on the process's [standard streams][stdin.3].
* **`consoleSize`** (object, OPTIONAL) specifies the console size in characters of the terminal.
    Runtimes MUST ignore `consoleSize` if `terminal` is `false` or unset.
    * **`height`** (uint, REQUIRED)
    * **`width`** (uint, REQUIRED)
* **`cwd`** (string, REQUIRED) is the working directory that will be set for the executable.
    This value MUST be an absolute path.
* **`env`** (array of strings, OPTIONAL) with the same semantics as [IEEE Std 1003.1-2008's `environ`][ieee-1003.1-2008-xbd-c8.1].
* **`args`** (array of strings, OPTIONAL) with similar semantics to [IEEE Std 1003.1-2008 `execvp`'s *argv*][ieee-1003.1-2008-functions-exec].
    This specification extends the IEEE standard in that at least one entry is REQUIRED (non-Windows), and that entry is used with the same semantics as `execvp`'s *file*. This field is OPTIONAL on Windows, and `commandLine` is REQUIRED if this field is omitted.
* **`commandLine`** (string, OPTIONAL) specifies the full command line to be executed on Windows.
    This is the preferred means of supplying the command line on Windows. If omitted, the runtime will fall back to escaping and concatenating fields from `args` before making the system call into Windows.


### <a name="configPOSIXProcess" />POSIX process

For systems that support POSIX rlimits (for example Linux and Solaris), the `process` object supports the following process-specific properties:

* **`rlimits`** (array of objects, OPTIONAL) allows setting resource limits for the process.
    Each entry has the following structure:

    * **`type`** (string, REQUIRED) the platform resource being limited.
        * Linux: valid values are defined in the [`getrlimit(2)`][getrlimit.2] man page, such as `RLIMIT_MSGQUEUE`.
        * Solaris: valid values are defined in the [`getrlimit(3)`][getrlimit.3] man page, such as `RLIMIT_CORE`.

        The runtime MUST [generate an error](runtime.md#errors) for any values which cannot be mapped to a relevant kernel interface.
        For each entry in `rlimits`, a [`getrlimit(3)`][getrlimit.3] on `type` MUST succeed.
        For the following properties, `rlim` refers to the status returned by the `getrlimit(3)` call.

    * **`soft`** (uint64, REQUIRED) the value of the limit enforced for the corresponding resource.
        `rlim.rlim_cur` MUST match the configured value.
    * **`hard`** (uint64, REQUIRED) the ceiling for the soft limit that could be set by an unprivileged process.
        `rlim.rlim_max` MUST match the configured value.
        Only a privileged process (e.g. one with the `CAP_SYS_RESOURCE` capability) can raise a hard limit.

    If `rlimits` contains duplicated entries with same `type`, the runtime MUST [generate an error](runtime.md#errors).

### <a name="configLinuxProcess" />Linux Process

For Linux-based systems, the `process` object supports the following process-specific properties.

* **`apparmorProfile`** (string, OPTIONAL) specifies the name of the AppArmor profile for the process.
    For more information about AppArmor, see [AppArmor documentation][apparmor].
* **`capabilities`** (object, OPTIONAL) is an object containing arrays that specifies the sets of capabilities for the process.
    Valid values are defined in the [capabilities(7)][capabilities.7] man page, such as `CAP_CHOWN`.
    Any value which cannot be mapped to a relevant kernel interface, or cannot
    be granted otherwise MUST be [logged as a warning](runtime.md#warnings) by
    the runtime. Runtimes SHOULD NOT fail if the container configuration requests
    capabilities that cannot be granted, for example, if the runtime operates in
    a restricted environment with a limited set of capabilities.
    `capabilities` contains the following properties:

    * **`effective`** (array of strings, OPTIONAL) the `effective` field is an array of effective capabilities that are kept for the process.
    * **`bounding`** (array of strings, OPTIONAL) the `bounding` field is an array of bounding capabilities that are kept for the process.
    * **`inheritable`** (array of strings, OPTIONAL) the `inheritable` field is an array of inheritable capabilities that are kept for the process.
    * **`permitted`** (array of strings, OPTIONAL) the `permitted` field is an array of permitted capabilities that are kept for the process.
    * **`ambient`** (array of strings, OPTIONAL) the `ambient` field is an array of ambient capabilities that are kept for the process.
* **`noNewPrivileges`** (bool, OPTIONAL) setting `noNewPrivileges` to true prevents the process from gaining additional privileges.
    As an example, the [`no_new_privs`][no-new-privs] article in the kernel documentation has information on how this is achieved using a `prctl` system call on Linux.
* **`oomScoreAdj`** *(int, OPTIONAL)* adjusts the oom-killer score in `[pid]/oom_score_adj` for the process's `[pid]` in a [proc pseudo-filesystem][proc_2].
    If `oomScoreAdj` is set, the runtime MUST set `oom_score_adj` to the given value.
    If `oomScoreAdj` is not set, the runtime MUST NOT change the value of `oom_score_adj`.

    This is a per-process setting, where as [`disableOOMKiller`](config-linux.md#memory) is scoped for a memory cgroup.
    For more information on how these two settings work together, see [the memory cgroup documentation section 10. OOM Contol][cgroup-v1-memory_2].
* **`scheduler`** (object, OPTIONAL) is an object describing the scheduler properties for the process.  The `scheduler` contains the following properties:

    * **`policy`** (string, REQUIRED) represents the scheduling policy.  A valid list of values is:

        * `SCHED_OTHER`
        * `SCHED_FIFO`
        * `SCHED_RR`
        * `SCHED_BATCH`
        * `SCHED_ISO`
        * `SCHED_IDLE`
        * `SCHED_DEADLINE`

    * **`nice`** (int32, OPTIONAL) is the nice value for the process, affecting its priority. A lower nice value corresponds to a higher priority. If not set, the runtime must use the value 0.
    * **`priority`** (int32, OPTIONAL) represents the static priority of the process, used by real-time policies like SCHED_FIFO and SCHED_RR. If not set, the runtime must use the value 0.
    * **`flags`** (array of strings, OPTIONAL) is an array of strings representing scheduling flags.  A valid list of values is:

        * `SCHED_FLAG_RESET_ON_FORK`
        * `SCHED_FLAG_RECLAIM`
        * `SCHED_FLAG_DL_OVERRUN`
        * `SCHED_FLAG_KEEP_POLICY`
        * `SCHED_FLAG_KEEP_PARAMS`
        * `SCHED_FLAG_UTIL_CLAMP_MIN`
        * `SCHED_FLAG_UTIL_CLAMP_MAX`

    * **`runtime`** (uint64, OPTIONAL) represents the amount of time in nanoseconds during which the process is allowed to run in a given period, used by the deadline scheduler. If not set, the runtime must use the value 0.
    * **`deadline`** (uint64, OPTIONAL) represents the absolute deadline for the process to complete its execution, used by the deadline scheduler. If not set, the runtime must use the value 0.
    * **`period`** (uint64, OPTIONAL) represents the length of the period in nanoseconds used for determining the process runtime, used by the deadline scheduler. If not set, the runtime must use the value 0.
* **`selinuxLabel`** (string, OPTIONAL) specifies the SELinux label for the process.
    For more information about SELinux, see  [SELinux documentation][selinux].
* **`ioPriority`** (object, OPTIONAL) configures the I/O priority settings for the container's processes within the process group.
    The I/O priority settings will be automatically applied to the entire process group, affecting all processes within the container.
    The following properties are available:

    * **`class`** (string, REQUIRED) specifies the I/O scheduling class. Possible values are `IOPRIO_CLASS_RT`, `IOPRIO_CLASS_BE`, and `IOPRIO_CLASS_IDLE`.
    * **`priority`** (int, REQUIRED) specifies the priority level within the class. The value should be an integer ranging from 0 (highest) to 7 (lowest).
* **`execCPUAffinity`** (object, OPTIONAL) specifies CPU affinity used to execute the process.
    This setting is not applicable to the container's init process.
    The following properties are available:
    * **`initial`** (string, OPTIONAL) is a list of CPUs a runtime parent
      process to be run on initially, before the transition to container's
      cgroup. This is a a comma-separated list, with dashes to represent
      ranges. For example, `0-3,7` represents CPUs 0,1,2,3, and 7.
    * **`final`** (string, OPTIONAL) is a list of CPUs the process will be run
      on after the transition to container's cgroup. The format is the same as
      for `initial`. If omitted or empty, the container's default CPU affinity,
      as defined by [cpu.cpus property](./config.md#configLinuxCPUs)), is used.

### <a name="configUser" />User

The user for the process is a platform-specific structure that allows specific control over which user the process runs as.

#### <a name="configPOSIXUser" />POSIX-platform User

For POSIX platforms the `user` structure has the following fields:

* **`uid`** (int, REQUIRED) specifies the user ID in the [container namespace](glossary.md#container-namespace).
* **`gid`** (int, REQUIRED) specifies the group ID in the [container namespace](glossary.md#container-namespace).
* **`umask`** (int, OPTIONAL) specifies the [umask][umask_2] of the user. If unspecified, the umask should not be changed from the calling process' umask.
* **`additionalGids`** (array of ints, OPTIONAL) specifies additional group IDs in the [container namespace](glossary.md#container-namespace) to be added to the process.

_Note: symbolic name for uid and gid, such as uname and gname respectively, are left to upper levels to derive (i.e. `/etc/passwd` parsing, NSS, etc)_

### Example (Linux)

```json
"process": {
    "terminal": true,
    "consoleSize": {
        "height": 25,
        "width": 80
    },
    "user": {
        "uid": 1,
        "gid": 1,
        "umask": 63,
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
    "ioPriority": {
        "class": "IOPRIO_CLASS_IDLE",
        "priority": 4
    },
    "noNewPrivileges": true,
    "capabilities": {
        "bounding": [
            "CAP_AUDIT_WRITE",
            "CAP_KILL",
            "CAP_NET_BIND_SERVICE"
        ],
       "permitted": [
            "CAP_AUDIT_WRITE",
            "CAP_KILL",
            "CAP_NET_BIND_SERVICE"
        ],
       "inheritable": [
            "CAP_AUDIT_WRITE",
            "CAP_KILL",
            "CAP_NET_BIND_SERVICE"
        ],
        "effective": [
            "CAP_AUDIT_WRITE",
            "CAP_KILL"
        ],
        "ambient": [
            "CAP_NET_BIND_SERVICE"
        ]
    },
    "rlimits": [
        {
            "type": "RLIMIT_NOFILE",
            "hard": 1024,
            "soft": 1024
        }
    ],
    "execCPUAffinity": {
        "initial": "7",
        "final": "0-3,7"
    }
}
```
### Example (Solaris)

```json
"process": {
    "terminal": true,
    "consoleSize": {
        "height": 25,
        "width": 80
    },
    "user": {
        "uid": 1,
        "gid": 1,
        "umask": 7,
        "additionalGids": [2, 8]
    },
    "env": [
        "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
        "TERM=xterm"
    ],
    "cwd": "/root",
    "args": [
        "/usr/bin/bash"
    ]
}
```

#### <a name="configWindowsUser" />Windows User

For Windows based systems the user structure has the following fields:

* **`username`** (string, OPTIONAL) specifies the user name for the process.

### Example (Windows)

```json
"process": {
    "terminal": true,
    "user": {
        "username": "containeradministrator"
    },
    "env": [
        "VARIABLE=1"
    ],
    "cwd": "c:\\foo",
    "args": [
        "someapp.exe",
    ]
}
```


## <a name="configHostname" />Hostname

* **`hostname`** (string, OPTIONAL) specifies the container's hostname as seen by processes running inside the container.
    On Linux, for example, this will change the hostname in the [container](glossary.md#container-namespace) [UTS namespace][uts-namespace.7].
    Depending on your [namespace configuration](config-linux.md#namespaces), the container UTS namespace may be the [runtime](glossary.md#runtime-namespace) [UTS namespace][uts-namespace.7].

### Example

```json
"hostname": "mrsdalloway"
```

## <a name="configDomainname" />Domainname

* **`domainname`** (string, OPTIONAL) specifies the container's domainname as seen by processes running inside the container.
    On Linux, for example, this will change the domainname in the [container](glossary.md#container-namespace) [UTS namespace][uts-namespace.7].
    Depending on your [namespace configuration](config-linux.md#namespaces), the container UTS namespace may be the [runtime](glossary.md#runtime-namespace) [UTS namespace][uts-namespace.7].

### Example

```json
"domainname": "foobarbaz.test"
```

## <a name="configPlatformSpecificConfiguration" />Platform-specific configuration

* **`linux`** (object, OPTIONAL) [Linux-specific configuration](config-linux.md).
    This MAY be set if the target platform of this spec is `linux`.
* **`windows`** (object, OPTIONAL) [Windows-specific configuration](config-windows.md).
    This MUST be set if the target platform of this spec is `windows`.
* **`solaris`** (object, OPTIONAL) [Solaris-specific configuration](config-solaris.md).
    This MAY be set if the target platform of this spec is `solaris`.
* **`vm`** (object, OPTIONAL) [Virtual-machine-specific configuration](config-vm.md).
    This MAY be set if the target platform and architecture of this spec support hardware virtualization.
* **`zos`** (object, OPTIONAL) [z/OS-specific configuration](config-zos.md).
    This MAY be set if the target platform of this spec is `zos`.

### Example (Linux)

```json
{
    "linux": {
        "namespaces": [
            {
                "type": "pid"
            }
        ]
    }
}
```

## <a name="configHooks" />POSIX-platform Hooks

For POSIX platforms, the configuration structure supports `hooks` for configuring custom actions related to the [lifecycle](runtime.md#lifecycle) of the container.

* **`hooks`** (object, OPTIONAL) MAY contain any of the following properties:
    * **`prestart`** (array of objects, OPTIONAL, **DEPRECATED**) is an array of [`prestart` hooks](#prestart).
        * Entries in the array contain the following properties:
            * **`path`** (string, REQUIRED) with similar semantics to [IEEE Std 1003.1-2008 `execv`'s *path*][ieee-1003.1-2008-functions-exec].
                This specification extends the IEEE standard in that **`path`** MUST be absolute.
            * **`args`** (array of strings, OPTIONAL) with the same semantics as [IEEE Std 1003.1-2008 `execv`'s *argv*][ieee-1003.1-2008-functions-exec].
            * **`env`** (array of strings, OPTIONAL) with the same semantics as [IEEE Std 1003.1-2008's `environ`][ieee-1003.1-2008-xbd-c8.1].
            * **`timeout`** (int, OPTIONAL) is the number of seconds before aborting the hook.
                If set, `timeout` MUST be greater than zero.
        * The value of `path` MUST resolve in the [runtime namespace](glossary.md#runtime-namespace).
        * The `prestart` hooks MUST be executed in the [runtime namespace](glossary.md#runtime-namespace).
    * **`createRuntime`** (array of objects, OPTIONAL) is an array of [`createRuntime` hooks](#createRuntime-hooks).
        * Entries in the array contain the following properties (the entries are identical to the entries in the deprecated `prestart` hooks):
            * **`path`** (string, REQUIRED) with similar semantics to [IEEE Std 1003.1-2008 `execv`'s *path*][ieee-1003.1-2008-functions-exec].
                This specification extends the IEEE standard in that **`path`** MUST be absolute.
            * **`args`** (array of strings, OPTIONAL) with the same semantics as [IEEE Std 1003.1-2008 `execv`'s *argv*][ieee-1003.1-2008-functions-exec].
            * **`env`** (array of strings, OPTIONAL) with the same semantics as [IEEE Std 1003.1-2008's `environ`][ieee-1003.1-2008-xbd-c8.1].
            * **`timeout`** (int, OPTIONAL) is the number of seconds before aborting the hook.
                If set, `timeout` MUST be greater than zero.
        * The value of `path` MUST resolve in the [runtime namespace](glossary.md#runtime-namespace).
        * The `createRuntime` hooks MUST be executed in the [runtime namespace](glossary.md#runtime-namespace).
    * **`createContainer`** (array of objects, OPTIONAL) is an array of [`createContainer` hooks](#createContainer-hooks).
        * Entries in the array have the same schema as `createRuntime` entries.
        * The value of `path` MUST resolve in the [runtime namespace](glossary.md#runtime-namespace).
        * The `createContainer` hooks MUST be executed in the [container namespace](glossary.md#container-namespace).
    * **`startContainer`** (array of objects, OPTIONAL) is an array of [`startContainer` hooks](#startContainer-hooks).
        * Entries in the array have the same schema as `createRuntime` entries.
        * The value of `path` MUST resolve in the [container namespace](glossary.md#container-namespace).
        * The `startContainer` hooks MUST be executed in the [container namespace](glossary.md#container-namespace).
    * **`poststart`** (array of objects, OPTIONAL) is an array of [`poststart` hooks](#poststart).
        * Entries in the array have the same schema as `createRuntime` entries.
        * The value of `path` MUST resolve in the [runtime namespace](glossary.md#runtime-namespace).
        * The `poststart` hooks MUST be executed in the [runtime namespace](glossary.md#runtime-namespace).
    * **`poststop`** (array of objects, OPTIONAL) is an array of [`poststop` hooks](#poststop).
        * Entries in the array have the same schema as `createRuntime` entries.
        * The value of `path` MUST resolve in the [runtime namespace](glossary.md#runtime-namespace).
        * The `poststop` hooks MUST be executed in the [runtime namespace](glossary.md#runtime-namespace).

Hooks allow users to specify programs to run before or after various lifecycle events.
Hooks MUST be called in the listed order.
The [state](runtime.md#state) of the container MUST be passed to hooks over stdin so that they may do work appropriate to the current state of the container.

### <a name="configHooksPrestart" />Prestart

The `prestart` hooks MUST be called as part of the [`create`](runtime.md#create) operation after the runtime environment has been created (according to the configuration in config.json) but before the `pivot_root` or any equivalent operation has been executed.
On Linux, for example, they are called after the container namespaces are created, so they provide an opportunity to customize the container (e.g. the network namespace could be specified in this hook).
The `prestart` hooks MUST be called before the `createRuntime` hooks.

Note: `prestart` hooks were deprecated in favor of `createRuntime`, `createContainer` and `startContainer` hooks, which allow more granular hook control during the create and start phase.

The `prestart` hooks' path MUST resolve in the [runtime namespace](glossary.md#runtime-namespace).
The `prestart` hooks MUST be executed in the [runtime namespace](glossary.md#runtime-namespace).

### <a name="configHooksCreateRuntime" />CreateRuntime Hooks

The `createRuntime` hooks MUST be called as part of the [`create`](runtime.md#create) operation after the runtime environment has been created (according to the configuration in config.json) but before the `pivot_root` or any equivalent operation has been executed.

The `createRuntime` hooks' path MUST resolve in the [runtime namespace](glossary.md#runtime-namespace).
The `createRuntime` hooks MUST be executed in the [runtime namespace](glossary.md#runtime-namespace).

On Linux, for example, they are called after the container namespaces are created, so they provide an opportunity to customize the container (e.g. the network namespace could be specified in this hook).

The definition of `createRuntime` hooks is currently underspecified and hooks authors, should only expect from the runtime that the mount namespace have been created and the mount operations performed. Other operations such as cgroups and SELinux/AppArmor labels might not have been performed by the runtime.

### <a name="configHooksCreateContainer" />CreateContainer Hooks

The `createContainer` hooks MUST be called as part of the [`create`](runtime.md#create) operation after the runtime environment has been created (according to the configuration in config.json) but before the `pivot_root` or any equivalent operation has been executed.
The `createContainer` hooks MUST be called after the `createRuntime` hooks.

The `createContainer` hooks' path MUST resolve in the [runtime namespace](glossary.md#runtime-namespace).
The `createContainer` hooks MUST be executed in the [container namespace](glossary.md#container-namespace).

For example, on Linux this would happen before the `pivot_root` operation is executed but after the mount namespace was created and setup.

The definition of `createContainer` hooks is currently underspecified and hooks authors, should only expect from the runtime that the mount namespace and different mounts will be setup. Other operations such as cgroups and SELinux/AppArmor labels might not have been performed by the runtime.

### <a name="configHooksStartContainer" />StartContainer Hooks

The `startContainer` hooks MUST be called [before the user-specified process is executed](runtime.md#lifecycle) as part of the [`start`](runtime.md#start) operation.
This hook can be used to execute some operations in the container, for example running the `ldconfig` binary on linux before the container process is spawned.

The `startContainer` hooks' path MUST resolve in the [container namespace](glossary.md#container-namespace).
The `startContainer` hooks MUST be executed in the [container namespace](glossary.md#container-namespace).

### <a name="configHooksPoststart" />Poststart

The `poststart` hooks MUST be called [after the user-specified process is executed](runtime.md#lifecycle) but before the [`start`](runtime.md#start) operation returns.
For example, this hook can notify the user that the container process is spawned.

The `poststart` hooks' path MUST resolve in the [runtime namespace](glossary.md#runtime-namespace).
The `poststart` hooks MUST be executed in the [runtime namespace](glossary.md#runtime-namespace).

### <a name="configHooksPoststop" />Poststop

The `poststop` hooks MUST be called [after the container is deleted](runtime.md#lifecycle) but before the [`delete`](runtime.md#delete) operation returns.
Cleanup or debugging functions are examples of such a hook.

The `poststop` hooks' path MUST resolve in the [runtime namespace](glossary.md#runtime-namespace).
The `poststop` hooks MUST be executed in the [runtime namespace](glossary.md#runtime-namespace).

### Summary

See the below table for a summary of hooks and when they are called:

|           Name          | Namespace |                                                            When                                                                    |
| ----------------------- | --------- | -----------------------------------------------------------------------------------------------------------------------------------|
| `prestart` (Deprecated) | runtime   | After the start  operation is called but before the user-specified program command is executed.                                    |
| `createRuntime`         | runtime   | During the create operation, after the runtime environment has been created and before the pivot root or any equivalent operation. |
| `createContainer`       | container | During the create operation, after the runtime environment has been created and before the pivot root or any equivalent operation. |
| `startContainer`        | container | After the start operation is called but before the user-specified program command is executed.                                     |
| `poststart`             | runtime   | After the user-specified process is executed but before the start operation returns.                                               |
| `poststop`              | runtime   | After the container is deleted but before the delete operation returns.                                                            |

### Example

```json
"hooks": {
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
    "createRuntime": [
        {
            "path": "/usr/bin/fix-mounts",
            "args": ["fix-mounts", "arg1", "arg2"],
            "env":  [ "key1=value1"]
        },
        {
            "path": "/usr/bin/setup-network"
        }
    ],
    "createContainer": [
        {
            "path": "/usr/bin/mount-hook",
            "args": ["-mount", "arg1", "arg2"],
            "env":  [ "key1=value1"]
        }
    ],
    "startContainer": [
        {
            "path": "/usr/bin/refresh-ldcache"
        }
    ],
    "poststart": [
        {
            "path": "/usr/bin/notify-start",
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

## <a name="configAnnotations" />Annotations

**`annotations`** (object, OPTIONAL) contains arbitrary metadata for the container.
This information MAY be structured or unstructured.
Annotations MUST be a key-value map.
If there are no annotations then this property MAY either be absent or an empty map.

Keys MUST be strings.
Keys MUST NOT be an empty string.
Keys SHOULD be named using a reverse domain notation - e.g. `com.example.myKey`.

The `org.opencontainers` namespace for keys is reserved for use by this specification, annotations using keys in this namespace MUST be as described in this section.
The following keys in the `org.opencontainers` namespaces MAY be used:
|                   Key                   | Definition                                                         |
| --------------------------------------- | -----------------------------------------------------------------------------------------------------------------------------------|
| `org.opencontainers.image.os`           | Indicates the operating system the container image was built to run on. The annotation value MUST have a valid value for the `os` property as defined in [the OCI image specification][oci-image-config-properties]. This annotation SHOULD only be used in accordance with the [OCI image specification's runtime conversion specification][oci-image-conversion]. |
| `org.opencontainers.image.os.version`   | Indicates the operating system version targeted by the container image. The annotation value MUST have a valid value for the `os.version` property as defined in [the OCI image specification][oci-image-config-properties]. This annotation SHOULD only be used in accordance with the [OCI image specification's runtime conversion specification][oci-image-conversion]. |
| `org.opencontainers.image.os.features`  | Indicates mandatory operating system features required by the container image. The annotation value MUST have a valid value for the `os.features` property as defined in [the OCI image specification][oci-image-config-properties]. This annotation SHOULD only be used in accordance with the [OCI image specification's runtime conversion specification][oci-image-conversion]. |
| `org.opencontainers.image.architecture` | Indicates the architecture that binaries in the container image are built to run on. The annotation value MUST have a valid value for the `architecture` property as defined in [the OCI image specification][oci-image-config-properties]. This annotation SHOULD only be used in accordance with the [OCI image specification's runtime conversion specification][oci-image-conversion]. |
| `org.opencontainers.image.variant`      | Indicates the variant of the architecture that binaries in the container image are built to run on. The annotation value MUST have a valid value for the `variant` property as defined in [the OCI image specification][oci-image-config-properties]. This annotation SHOULD only be used in accordance with the [OCI image specification's runtime conversion specification][oci-image-conversion]. |
| `org.opencontainers.image.author`       | Indicates the author of the container image. The annotation value MUST have a valid value for the `author` property as defined in [the OCI image specification][oci-image-config-properties]. This annotation SHOULD only be used in accordance with the [OCI image specification's runtime conversion specification][oci-image-conversion]. |
| `org.opencontainers.image.created`      | Indicates the date and time when the container image was created. The annotation value MUST have a valid value for the `created` property as defined in [the OCIimage specification][oci-image-config-properties]. This annotation SHOULD only be used in accordance with the [OCI image specification's runtime conversion specification][oci-image-conversion]. |
| `org.opencontainers.image.stopSignal`   | Indicates signal that SHOULD be sent by the container runtimes to [kill the container](runtime.md#kill). The annotation value MUST have a valid value for the `config.StopSignal` property as defined in [the OCI image specification][oci-image-config-properties]. This annotation SHOULD only be used in accordance with the [OCI image specification's runtime conversion specification][oci-image-conversion]. |

All other keys in the `org.opencontainers` namespace not specified in this above table are reserved and MUST NOT be used by subsequent specifications.
Runtimes MUST handle unknown annotation keys like any other [unknown property](#extensibility).

Values MUST be strings.
Values MAY be an empty string.

```json
"annotations": {
    "com.example.gpu-cores": "2"
}
```

## <a name="configExtensibility" />Extensibility

Runtimes MAY [log](runtime.md#warnings) unknown properties but MUST otherwise ignore them.
That includes not [generating errors](runtime.md#errors) if they encounter an unknown property.

## Valid values

Runtimes MUST generate an error when invalid or unsupported values are encountered.
Unless support for a valid value is explicitly required, runtimes MAY choose which subset of the valid values it will support.

## Configuration Schema Example

Here is a full example `config.json` for reference.

```json
{
    "ociVersion": "1.0.1",
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
        "capabilities": {
            "bounding": [
                "CAP_AUDIT_WRITE",
                "CAP_KILL",
                "CAP_NET_BIND_SERVICE"
            ],
            "permitted": [
                "CAP_AUDIT_WRITE",
                "CAP_KILL",
                "CAP_NET_BIND_SERVICE"
            ],
            "inheritable": [
                "CAP_AUDIT_WRITE",
                "CAP_KILL",
                "CAP_NET_BIND_SERVICE"
            ],
            "effective": [
                "CAP_AUDIT_WRITE",
                "CAP_KILL"
            ],
            "ambient": [
                "CAP_NET_BIND_SERVICE"
            ]
        },
        "rlimits": [
            {
                "type": "RLIMIT_CORE",
                "hard": 1024,
                "soft": 1024
            },
            {
                "type": "RLIMIT_NOFILE",
                "hard": 1024,
                "soft": 1024
            }
        ],
        "apparmorProfile": "acme_secure_profile",
        "oomScoreAdj": 100,
        "selinuxLabel": "system_u:system_r:svirt_lxc_net_t:s0:c124,c675",
        "ioPriority": {
            "class": "IOPRIO_CLASS_IDLE",
            "priority": 4
        },
        "noNewPrivileges": true
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
                "path": "/usr/bin/fix-mounts",
                "args": [
                    "fix-mounts",
                    "arg1",
                    "arg2"
                ],
                "env": [
                    "key1=value1"
                ]
            },
            {
                "path": "/usr/bin/setup-network"
            }
        ],
        "poststart": [
            {
                "path": "/usr/bin/notify-start",
                "timeout": 5
            }
        ],
        "poststop": [
            {
                "path": "/usr/sbin/cleanup.sh",
                "args": [
                    "cleanup.sh",
                    "-f"
                ]
            }
        ]
    },
    "linux": {
        "devices": [
            {
                "path": "/dev/fuse",
                "type": "c",
                "major": 10,
                "minor": 229,
                "fileMode": 438,
                "uid": 0,
                "gid": 0
            },
            {
                "path": "/dev/sda",
                "type": "b",
                "major": 8,
                "minor": 0,
                "fileMode": 432,
                "uid": 0,
                "gid": 0
            }
        ],
        "uidMappings": [
            {
                "containerID": 0,
                "hostID": 1000,
                "size": 32000
            }
        ],
        "gidMappings": [
            {
                "containerID": 0,
                "hostID": 1000,
                "size": 32000
            }
        ],
        "sysctl": {
            "net.ipv4.ip_forward": "1",
            "net.core.somaxconn": "256"
        },
        "cgroupsPath": "/myRuntime/myContainer",
        "resources": {
            "network": {
                "classID": 1048577,
                "priorities": [
                    {
                        "name": "eth0",
                        "priority": 500
                    },
                    {
                        "name": "eth1",
                        "priority": 1000
                    }
                ]
            },
            "pids": {
                "limit": 32771
            },
            "hugepageLimits": [
                {
                    "pageSize": "2MB",
                    "limit": 9223372036854772000
                },
                {
                    "pageSize": "64KB",
                    "limit": 1000000
                }
            ],
            "memory": {
                "limit": 536870912,
                "reservation": 536870912,
                "swap": 536870912,
                "kernel": -1,
                "kernelTCP": -1,
                "swappiness": 0,
                "disableOOMKiller": false
            },
            "cpu": {
                "shares": 1024,
                "quota": 1000000,
                "period": 500000,
                "realtimeRuntime": 950000,
                "realtimePeriod": 1000000,
                "cpus": "2-3",
                "idle": 1,
                "mems": "0-7"
            },
            "devices": [
                {
                    "allow": false,
                    "access": "rwm"
                },
                {
                    "allow": true,
                    "type": "c",
                    "major": 10,
                    "minor": 229,
                    "access": "rw"
                },
                {
                    "allow": true,
                    "type": "b",
                    "major": 8,
                    "minor": 0,
                    "access": "r"
                }
            ],
            "blockIO": {
                "weight": 10,
                "leafWeight": 10,
                "weightDevice": [
                    {
                        "major": 8,
                        "minor": 0,
                        "weight": 500,
                        "leafWeight": 300
                    },
                    {
                        "major": 8,
                        "minor": 16,
                        "weight": 500
                    }
                ],
                "throttleReadBpsDevice": [
                    {
                        "major": 8,
                        "minor": 0,
                        "rate": 600
                    }
                ],
                "throttleWriteIOPSDevice": [
                    {
                        "major": 8,
                        "minor": 16,
                        "rate": 300
                    }
                ]
            },
            "vtpms": [
                {
                    "statePath": "/var/lib/runc/myvtpm1",
                    "vtpmVersion": "2",
                    "createCerts": false,
                    "runAs": "tss",
                    "pcrBanks": "sha1,sha512"
                }
            ]
        },
        "rootfsPropagation": "slave",
        "seccomp": {
            "defaultAction": "SCMP_ACT_ALLOW",
            "architectures": [
                "SCMP_ARCH_X86",
                "SCMP_ARCH_X32"
            ],
            "syscalls": [
                {
                    "names": [
                        "getcwd",
                        "chmod"
                    ],
                    "action": "SCMP_ACT_ERRNO"
                }
            ]
        },
        "timeOffsets": {
            "monotonic": {
                "secs": 172800,
                "nanosecs": 0
            },
            "boottime": {
                "secs": 604800,
                "nanosecs": 0
            }
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
            },
            {
                "type": "user"
            },
            {
                "type": "cgroup"
            },
            {
                "type": "time"
            }
        ],
        "maskedPaths": [
            "/proc/kcore",
            "/proc/latency_stats",
            "/proc/timer_stats",
            "/proc/sched_debug"
        ],
        "readonlyPaths": [
            "/proc/asound",
            "/proc/bus",
            "/proc/fs",
            "/proc/irq",
            "/proc/sys",
            "/proc/sysrq-trigger"
        ],
        "mountLabel": "system_u:object_r:svirt_sandbox_file_t:s0:c715,c811"
    },
    "annotations": {
        "com.example.key1": "value1",
        "com.example.key2": "value2"
    }
}
```


[apparmor]: https://wiki.ubuntu.com/AppArmor
[cgroup-v1-memory_2]: https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
[selinux]:http://selinuxproject.org/page/Main_Page
[no-new-privs]: https://www.kernel.org/doc/Documentation/prctl/no_new_privs.txt
[proc_2]: https://www.kernel.org/doc/Documentation/filesystems/proc.txt
[umask.2]: http://pubs.opengroup.org/onlinepubs/009695399/functions/umask.html
[semver-v2.0.0]: http://semver.org/spec/v2.0.0.html
[ieee-1003.1-2008-xbd-c8.1]: http://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap08.html#tag_08_01
[ieee-1003.1-2008-functions-exec]: http://pubs.opengroup.org/onlinepubs/9699919799/functions/exec.html
[naming-a-volume]: https://aka.ms/nb3hqb
[oci-image-config-properties]: https://github.com/opencontainers/image-spec/blob/v1.1.0-rc2/config.md#properties
[oci-image-conversion]: https://github.com/opencontainers/image-spec/blob/v1.1.0-rc2/conversion.md

[capabilities.7]: http://man7.org/linux/man-pages/man7/capabilities.7.html
[mount.2]: http://man7.org/linux/man-pages/man2/mount.2.html
[mount.8]: http://man7.org/linux/man-pages/man8/mount.8.html
[mount.8-filesystem-independent]: http://man7.org/linux/man-pages/man8/mount.8.html#FILESYSTEM-INDEPENDENT_MOUNT_OPTIONS
[mount.8-filesystem-specific]: http://man7.org/linux/man-pages/man8/mount.8.html#FILESYSTEM-SPECIFIC_MOUNT_OPTIONS
[mount_setattr.2]: http://man7.org/linux/man-pages/man2/mount_setattr.2.html
[mount-bind]: https://docs.kernel.org/filesystems/sharedsubtree.html
[getrlimit.2]: http://man7.org/linux/man-pages/man2/getrlimit.2.html
[getrlimit.3]: http://pubs.opengroup.org/onlinepubs/9699919799/functions/getrlimit.html
[stdin.3]: http://man7.org/linux/man-pages/man3/stdin.3.html
[uts-namespace.7]: http://man7.org/linux/man-pages/man7/namespaces.7.html
[zonecfg.1m]: http://docs.oracle.com/cd/E86824_01/html/E54764/zonecfg-1m.html
