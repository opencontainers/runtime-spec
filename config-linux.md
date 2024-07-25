# <a name="linuxContainerConfiguration" />Linux Container Configuration

This document describes the schema for the [Linux-specific section](config.md#platform-specific-configuration) of the [container configuration](config.md).
The Linux container specification uses various kernel features like namespaces, cgroups, capabilities, LSM, and filesystem jails to fulfill the spec.

## <a name="configLinuxDefaultFilesystems" />Default Filesystems

The Linux ABI includes both syscalls and several special file paths.
Applications expecting a Linux environment will very likely expect these file paths to be set up correctly.

The following filesystems SHOULD be made available in each container's filesystem:

| Path     | Type   |
| -------- | ------ |
| /proc    | [proc][] |
| /sys     | [sysfs][]  |
| /dev/pts | [devpts][] |
| /dev/shm | [tmpfs][]  |

## <a name="configLinuxNamespaces" />Namespaces

A namespace wraps a global system resource in an abstraction that makes it appear to the processes within the namespace that they have their own isolated instance of the global resource.
Changes to the global resource are visible to other processes that are members of the namespace, but are invisible to other processes.
For more information, see the [namespaces(7)][namespaces.7_2] man page.

Namespaces are specified as an array of entries inside the `namespaces` root field.
The following parameters can be specified to set up namespaces:

* **`type`** *(string, REQUIRED)* - namespace type. The following namespace types SHOULD be supported:
    * **`pid`** processes inside the container will only be able to see other processes inside the same container or inside the same pid namespace.
    * **`network`** the container will have its own network stack.
    * **`mount`** the container will have an isolated mount table.
    * **`ipc`** processes inside the container will only be able to communicate to other processes inside the same container via system level IPC.
    * **`uts`** the container will be able to have its own hostname and domain name.
    * **`user`** the container will be able to remap user and group IDs from the host to local users and groups within the container.
    * **`cgroup`** the container will have an isolated view of the cgroup hierarchy.
    * **`time`** the container will be able to have its own clocks.
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
        "type": "network",
        "path": "/var/run/netns/neta"
    },
    {
        "type": "mount"
    },
    {
        "type": "ipc"
    },
    {
        "type": "uts"
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
]
```

## <a name="configLinuxUserNamespaceMappings" />User namespace mappings

**`uidMappings`** (array of objects, OPTIONAL) describes the user namespace uid mappings from the host to the container.
**`gidMappings`** (array of objects, OPTIONAL) describes the user namespace gid mappings from the host to the container.

Each entry has the following structure:

* **`containerID`** *(uint32, REQUIRED)* - is the starting uid/gid in the container.
* **`hostID`** *(uint32, REQUIRED)* - is the starting uid/gid on the host to be mapped to *containerID*.
* **`size`** *(uint32, REQUIRED)* - is the number of ids to be mapped.

The runtime SHOULD NOT modify the ownership of referenced filesystems to realize the mapping.
Note that the number of mapping entries MAY be limited by the [kernel][user-namespaces].

### Example

```json
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
]
```

## <a name="configLinuxTimeOffset" />Offset for Time Namespace

**`timeOffsets`** (object, OPTIONAL) sets the offset for Time Namespace. For more information
see the [time_namespaces][time_namespaces.7].

The name of the clock is the entry key.
Entry values are objects with the following properties:

* **`secs`** *(int64, OPTIONAL)* - is the offset of clock (in seconds) in the container.
* **`nanosecs`** *(uint32, OPTIONAL)* - is the offset of clock (in nanoseconds) in the container.

## <a name="configLinuxDevices" />Devices

**`devices`** (array of objects, OPTIONAL) lists devices that MUST be available in the container.
The runtime MAY supply them however it likes (with [`mknod`][mknod.2], by bind mounting from the runtime mount namespace, using symlinks, etc.).

Each entry has the following structure:

* **`type`** *(string, REQUIRED)* - type of device: `c`, `b`, `u` or `p`.
    More info in [mknod(1)][mknod.1].
* **`path`** *(string, REQUIRED)* - full path to device inside container.
    If a [file][] already exists at `path` that does not match the requested device, the runtime MUST generate an error.
    The path MAY be anywhere in the container filesystem, notably outside of `/dev`.
* **`major, minor`** *(int64, REQUIRED unless `type` is `p`)* - [major, minor numbers][devices] for the device.
* **`fileMode`** *(uint32, OPTIONAL)* - file mode for the device.
    You can also control access to devices [with cgroups](#configLinuxDeviceAllowedlist).
* **`uid`** *(uint32, OPTIONAL)* - id of device owner in the [container namespace](glossary.md#container-namespace).
* **`gid`** *(uint32, OPTIONAL)* - id of device group in the [container namespace](glossary.md#container-namespace).

The same `type`, `major` and `minor` SHOULD NOT be used for multiple devices.

Containers MAY NOT access any device node that is not either explicitly
referenced in the **`devices`** array or listed as being part of the
[default devices](#configLinuxDefaultDevices).
Rationale: runtimes based on virtual machines need to be able to adjust the node
devices, and accessing device nodes that were not adjusted could have undefined
behaviour.


### Example

```json
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
]
```

### <a name="configLinuxDefaultDevices" />Default Devices

In addition to any devices configured with this setting, the runtime MUST also supply:

* [`/dev/null`][null.4]
* [`/dev/zero`][zero.4]
* [`/dev/full`][full.4]
* [`/dev/random`][random.4]
* [`/dev/urandom`][random.4]
* [`/dev/tty`][tty.4]
* `/dev/console` is set up if [`terminal`](config.md#process) is enabled in the config by bind mounting the pseudoterminal pty to `/dev/console`.
* [`/dev/ptmx`][pts.4].
  A [bind-mount or symlink of the container's `/dev/pts/ptmx`][devpts].

## <a name="configLinuxControlGroups" />Control groups

Also known as cgroups, they are used to restrict resource usage for a container and handle device access.
cgroups provide controls (through controllers) to restrict cpu, memory, IO, pids, network and RDMA resources for the container.
For more information, see the [kernel cgroups documentation][cgroup-v1].

A runtime MAY, during a particular [container operation](runtime.md#operation),
such as [create](runtime.md#create), [start](runtime.md#start), or
[exec](runtime.md#exec), check if the container cgroup is fit for purpose,
and MUST [generate an error](runtime.md#errors) if such a check fails.
For example, a frozen cgroup or (for [create](runtime.md#create) operation)
a non-empty cgroup. The reason for this is that accepting such configurations
could cause container operation outcomes that users may not anticipate or
understand, such as operation on one container inadvertently affecting other
containers.

### <a name="configLinuxCgroupsPath" />Cgroups Path

**`cgroupsPath`** (string, OPTIONAL) path to the cgroups.
It can be used to either control the cgroups hierarchy for containers or to run a new process in an existing container.

The value of `cgroupsPath` MUST be either an absolute path or a relative path.

* In the case of an absolute path (starting with `/`), the runtime MUST take the path to be relative to the cgroups mount point.
* In the case of a relative path (not starting with `/`), the runtime MAY interpret the path relative to a runtime-determined location in the cgroups hierarchy.

If the value is specified, the runtime MUST consistently attach to the same place in the cgroups hierarchy given the same value of `cgroupsPath`.
If the value is not specified, the runtime MAY define the default cgroups path.
Runtimes MAY consider certain `cgroupsPath` values to be invalid, and MUST generate an error if this is the case.

Implementations of the Spec can choose to name cgroups in any manner.
The Spec does not include naming schema for cgroups.
The Spec does not support per-controller paths for the reasons discussed in the [cgroupv2 documentation][cgroup-v2].
The cgroups will be created if they don't exist.

You can configure a container's cgroups via the `resources` field of the Linux configuration.
Do not specify `resources` unless limits have to be updated.
For example, to run a new process in an existing container without updating limits, `resources` need not be specified.

Runtimes MAY attach the container process to additional cgroup controllers beyond those necessary to fulfill the `resources` settings.

### Cgroup ownership

Runtimes MAY, according to the following rules, change (or cause to
be changed) the owner of the container's cgroup to the host uid that
maps to the value of `process.user.uid` in the [container
namespace](glossary.md#container-namespace); that is, the user that
will execute the container process.

Runtimes SHOULD NOT change the ownership of container cgroups when
cgroups v1 is in use.  Cgroup delegation is not secure in cgroups
v1.

A runtime SHOULD NOT change the ownership of a container cgroup
unless it will also create a new cgroup namespace for the container.
Typically this occurs when the `linux.namespaces` array contains an
object with `type` equal to `"cgroup"` and `path` unset.

Runtimes SHOULD change the cgroup ownership if and only if the
cgroup filesystem is to be mounted read/write; that is, when the
configuration's `mounts` array contains an object where:

- The `source` field is equal to `"cgroup"`
- The `destination` field is equal to `"/sys/fs/cgroup"`
- The `options` field does not contain the value `"ro"`

If the configuration does not specify such a mount, the runtime
SHOULD NOT change the cgroup ownership.

A runtime that changes the cgroup ownership SHOULD only change the
ownership of the container's cgroup directory and files within that
directory that are listed in `/sys/kernel/cgroup/delegate`.  See
`cgroups(7)` for details about this file.  Note that not all files
listed in `/sys/kernel/cgroup/delegate` necessarily exist in every
cgroup.  Runtimes MUST NOT fail in this scenario, and SHOULD change
the ownership of the listed files that do exist in the cgroup.

If the `/sys/kernel/cgroup/delegate` file does not exist, the
runtime MUST fall back to using the following list of files:

```
cgroup.procs
cgroup.subtree_control
cgroup.threads
```

The runtime SHOULD NOT change the ownership of any other files.
Changing other files may allow the container to elevate its own
resource limits or perform other unwanted behaviour.

### Example

```json
"cgroupsPath": "/myRuntime/myContainer",
"resources": {
    "memory": {
    "limit": 100000,
    "reservation": 200000
    },
    "devices": [
        {
            "allow": false,
            "access": "rwm"
        }
    ]
}
```

### <a name="configLinuxDeviceAllowedlist" />Allowed Device list

**`devices`** (array of objects, OPTIONAL) configures the [allowed device list][cgroup-v1-devices].
The runtime MUST apply entries in the listed order.

Each entry has the following structure:

* **`allow`** *(boolean, REQUIRED)* - whether the entry is allowed or denied.
* **`type`** *(string, OPTIONAL)* - type of device: `a` (all), `c` (char), or `b` (block).
    Unset values mean "all", mapping to `a`.
* **`major, minor`** *(int64, OPTIONAL)* - [major, minor numbers][devices] for the device.
    Unset values mean "all", mapping to [`*` in the filesystem API][cgroup-v1-devices].
* **`access`** *(string, OPTIONAL)* - cgroup permissions for device.
    A composition of `r` (read), `w` (write), and `m` (mknod).

#### Example

```json
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
]
```

### <a name="configLinuxMemory" />Memory

**`memory`** (object, OPTIONAL) represents the cgroup subsystem `memory` and it's used to set limits on the container's memory usage.
For more information, see the kernel cgroups documentation about [memory][cgroup-v1-memory].

Values for memory specify the limit in bytes, or `-1` for unlimited memory.

* **`limit`** *(int64, OPTIONAL)* - sets limit of memory usage
* **`reservation`** *(int64, OPTIONAL)* - sets soft limit of memory usage
* **`swap`** *(int64, OPTIONAL)* - sets limit of memory+Swap usage
* **`kernel`** *(int64, OPTIONAL, NOT RECOMMENDED)* - sets hard limit for kernel memory
* **`kernelTCP`** *(int64, OPTIONAL, NOT RECOMMENDED)* - sets hard limit for kernel TCP buffer memory

The following properties do not specify memory limits, but are covered by the `memory` controller:

* **`swappiness`** *(uint64, OPTIONAL)* - sets swappiness parameter of vmscan (See sysctl's vm.swappiness)
    The values are from 0 to 100. Higher means more swappy.
* **`disableOOMKiller`** *(bool, OPTIONAL)* - enables or disables the OOM killer.
    If enabled (`false`), tasks that attempt to consume more memory than they are allowed are immediately killed by the OOM killer.
    The OOM killer is enabled by default in every cgroup using the `memory` subsystem.
    To disable it, specify a value of `true`.
* **`useHierarchy`** *(bool, OPTIONAL)* - enables or disables hierarchical memory accounting.
    If enabled (`true`), child cgroups will share the memory limits of this cgroup.
* **`checkBeforeUpdate`** *(bool, OPTIONAL)* - enables container memory usage check before setting a new limit.
    If enabled (`true`), runtime MAY check if a new memory limit is lower than the current usage, and MUST
    reject the new limit. Practically, when cgroup v1 is used, the kernel rejects the limit lower than the
    current usage, and when cgroup v2 is used, an OOM killer is invoked. This setting can be used on
    cgroup v2 to mimic the cgroup v1 behavior.

#### Example

```json
"memory": {
    "limit": 536870912,
    "reservation": 536870912,
    "swap": 536870912,
    "kernel": -1,
    "kernelTCP": -1,
    "swappiness": 0,
    "disableOOMKiller": false
}
```

### <a name="configLinuxCPU" />CPU

**`cpu`** (object, OPTIONAL) represents the cgroup subsystems `cpu` and `cpusets`.
For more information, see the kernel cgroups documentation about [cpusets][cgroup-v1-cpusets].

The following parameters can be specified to set up the controller:

* **`shares`** *(uint64, OPTIONAL)* - specifies a relative share of CPU time available to the tasks in a cgroup
* **`quota`** *(int64, OPTIONAL)* - specifies the total amount of time in microseconds for which all tasks in a cgroup can run during one period (as defined by **`period`** below)
    If specified with any (valid) positive value, it MUST be no smaller than `burst` (runtimes MAY generate an error).
* **`burst`** *(uint64, OPTIONAL)* - specifies the maximum amount of accumulated time in microseconds for which all tasks in a cgroup can run additionally for burst during one period (as defined by **`period`** below)
    If specified, this value MUST be no larger than any positive `quota` (runtimes MAY generate an error).
* **`period`** *(uint64, OPTIONAL)* - specifies a period of time in microseconds for how regularly a cgroup's access to CPU resources should be reallocated (CFS scheduler only)
* **`realtimeRuntime`** *(int64, OPTIONAL)* - specifies a period of time in microseconds for the longest continuous period in which the tasks in a cgroup have access to CPU resources
* **`realtimePeriod`** *(uint64, OPTIONAL)* - same as **`period`** but applies to realtime scheduler only
* **`cpus`** *(string, OPTIONAL)* - list of CPUs the container will run on. This is a comma-separated list, with dashes to represent ranges. For example, `0-3,7` represents CPUs 0,1,2,3, and 7.
* **`mems`** *(string, OPTIONAL)* - list of memory nodes the container will run on. This is a comma-separated list, with dashes to represent ranges. For example, `0-3,7` represents memory nodes 0,1,2,3, and 7.
* **`idle`** *(int64, OPTIONAL)* - cgroups are configured with minimum weight, 0: default behavior, 1: SCHED_IDLE.

#### Example

```json
"cpu": {
    "shares": 1024,
    "quota": 1000000,
    "burst": 1000000,
    "period": 500000,
    "realtimeRuntime": 950000,
    "realtimePeriod": 1000000,
    "cpus": "2-3",
    "mems": "0-7",
    "idle": 0
}
```

### <a name="configLinuxBlockIO" />Block IO

**`blockIO`** (object, OPTIONAL) represents the cgroup subsystem `blkio` which implements the block IO controller.
For more information, see the kernel cgroups documentation about [blkio][cgroup-v1-blkio] of cgroup v1 or [io][cgroup-v2-io] of cgroup v2, .

Note that I/O throttling settings in cgroup v1 apply only to Direct I/O due to kernel implementation constraints, while this limitation does not exist in cgroup v2.

The following parameters can be specified to set up the controller:

* **`weight`** *(uint16, OPTIONAL)* - specifies per-cgroup weight. This is default weight of the group on all devices until and unless overridden by per-device rules.
* **`leafWeight`** *(uint16, OPTIONAL)* - equivalents of `weight` for the purpose of deciding how much weight tasks in the given cgroup has while competing with the cgroup's child cgroups.
* **`weightDevice`** *(array of objects, OPTIONAL)* - an array of per-device bandwidth weights.
    Each entry has the following structure:
    * **`major, minor`** *(int64, REQUIRED)* - major, minor numbers for device.
        For more information, see the [mknod(1)][mknod.1] man page.
    * **`weight`** *(uint16, OPTIONAL)* - bandwidth weight for the device.
    * **`leafWeight`** *(uint16, OPTIONAL)* - bandwidth weight for the device while competing with the cgroup's child cgroups, CFQ scheduler only

    You MUST specify at least one of `weight` or `leafWeight` in a given entry, and MAY specify both.

* **`throttleReadBpsDevice`**, **`throttleWriteBpsDevice`** *(array of objects, OPTIONAL)* - an array of per-device bandwidth rate limits.
    Each entry has the following structure:
    * **`major, minor`** *(int64, REQUIRED)* - major, minor numbers for device.
        For more information, see the [mknod(1)][mknod.1] man page.
    * **`rate`** *(uint64, REQUIRED)* - bandwidth rate limit in bytes per second for the device

* **`throttleReadIOPSDevice`**, **`throttleWriteIOPSDevice`** *(array of objects, OPTIONAL)* - an array of per-device IO rate limits.
    Each entry has the following structure:
    * **`major, minor`** *(int64, REQUIRED)* - major, minor numbers for device.
        For more information, see the [mknod(1)][mknod.1] man page.
    * **`rate`** *(uint64, REQUIRED)* - IO rate limit for the device

#### Example

```json
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
}
```

### <a name="configLinuxVTPMs" />vTPMs

**`vtpms`** (array of objects, OPTIONAL) lists a number of emulated TPMs that will be made available to the container.

Each entry has the following structure:

* **`statePath`** *(string, REQUIRED)* - Unique path where vTPM writes its state into.
* **`statePathIsManaged`** *(string, OPTIONAL)* - Whether runc is allowed to delete the TPM's state path upon destroying the TPM, defaults to false.
* **`vtpmVersion`** *(string, OPTIONAL)* - The version of TPM to emulate, either 1.2 or 2, defaults to 1.2.
* **`createCerts`** *(boolean, OPTIONAL)* - If true then create certificates for the vTPM, defaults to false.
* **`runAs`** *(string, OPTIONAL)* - Under which user to run the vTPM, e.g.  'tss'.
Contributor
@mrunalp mrunalp on Aug 7, 2020
Does it make sense to run this as the container user or it is typically set to a separate tss user?

@KevinLi1020	Reply...
* **`pcrBanks`** *(string, OPTIONAL)* - Comma-separated list of PCR banks to activate, default depends on `swtpm`.
* **`encryptionPassword`** *(string, OPTIONAL)* - Write state encrypted with a key derived from the password, defaults to not encrypted.

#### Example

```json
    "vtpms": [
        {
            "statePath": "/var/lib/runc/myvtpm1",
            "statePathIsManaged": false,
            "vtpmVersion": "2",
            "createCerts": false,
            "runAs": "tss",
            "pcrBanks": "sha1,sha512",
            "encryptionPassword": "mysecret"
        }
    ]
```

### <a name="configLinuxHugePageLimits" />Huge page limits

**`hugepageLimits`** (array of objects, OPTIONAL) represents the `hugetlb` controller which allows to limit the HugeTLB reservations (if supported) or usage (page fault).
By default if supported by the kernel, `hugepageLimits` defines the hugepage sizes and limits for HugeTLB controller
reservation accounting, which allows to limit the HugeTLB reservations per control group and enforces the controller
limit at reservation time and at the fault of HugeTLB memory for which no reservation exists.
Otherwise if not supported by the kernel, this should fallback to the page fault accounting, which allows users to limit
the HugeTLB usage (page fault) per control group and enforces the limit during page fault.

Note that reservation limits are superior to page fault limits, since reservation limits are enforced at reservation
time (on mmap or shget), and never causes the application to get SIGBUS signal if the memory was reserved before hand.
This allows for easier fallback to alternatives such as non-HugeTLB memory for example. In the case of page fault
accounting, it's very hard to avoid processes getting SIGBUS since the sysadmin needs precisely know the HugeTLB usage
of all the tasks in the system and make sure there is enough pages to satisfy all requests. Avoiding tasks getting
SIGBUS on overcommited systems is practically impossible with page fault accounting.

For more information, see the kernel cgroups documentation about [HugeTLB][cgroup-v1-hugetlb].

Each entry has the following structure:

* **`pageSize`** *(string, REQUIRED)* - hugepage size.
    The value has the format `<size><unit-prefix>B` (64KB, 2MB, 1GB), and must match the `<hugepagesize>` of the
    corresponding control file found in `/sys/fs/cgroup/hugetlb/hugetlb.<hugepagesize>.rsvd.limit_in_bytes` (if
    hugetlb_cgroup reservation is supported) or `/sys/fs/cgroup/hugetlb/hugetlb.<hugepagesize>.limit_in_bytes` (if not
    supported).
    Values of `<unit-prefix>` are intended to be parsed using base 1024 ("1KB" = 1024, "1MB" = 1048576, etc).
* **`limit`** *(uint64, REQUIRED)* - limit in bytes of *hugepagesize* HugeTLB reservations (if supported) or usage.

#### Example

```json
"hugepageLimits": [
    {
        "pageSize": "2MB",
        "limit": 209715200
    },
    {
        "pageSize": "64KB",
        "limit": 1000000
    }
]
```

### <a name="configLinuxNetwork" />Network

**`network`** (object, OPTIONAL) represents the cgroup subsystems `net_cls` and `net_prio`.
For more information, see the kernel cgroups documentations about [net\_cls cgroup][cgroup-v1-net-cls] and [net\_prio cgroup][cgroup-v1-net-prio].

The following parameters can be specified to set up the controller:

* **`classID`** *(uint32, OPTIONAL)* - is the network class identifier the cgroup's network packets will be tagged with
* **`priorities`** *(array of objects, OPTIONAL)* - specifies a list of objects of the priorities assigned to traffic originating from processes in the group and egressing the system on various interfaces.
    The following parameters can be specified per-priority:
    * **`name`** *(string, REQUIRED)* - interface name in [runtime network namespace](glossary.md#runtime-namespace)
    * **`priority`** *(uint32, REQUIRED)* - priority applied to the interface

#### Example

```json
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
}
```

### <a name="configLinuxPIDS" />PIDs

**`pids`** (object, OPTIONAL) represents the cgroup subsystem `pids`.
For more information, see the kernel cgroups documentation about [pids][cgroup-v1-pids].

The following parameters can be specified to set up the controller:

* **`limit`** *(int64, REQUIRED)* - specifies the maximum number of tasks in the cgroup

#### Example

```json
"pids": {
    "limit": 32771
}
```

### <a name="configLinuxRDMA" />RDMA

**`rdma`** (object, OPTIONAL) represents the cgroup subsystem `rdma`.
For more information, see the kernel cgroups documentation about [rdma][cgroup-v1-rdma].

The name of the device to limit is the entry key.
Entry values are objects with the following properties:

* **`hcaHandles`** *(uint32, OPTIONAL)* - specifies the maximum number of hca_handles in the cgroup
* **`hcaObjects`** *(uint32, OPTIONAL)* - specifies the maximum number of hca_objects in the cgroup

You MUST specify at least one of the `hcaHandles` or `hcaObjects` in a given entry, and MAY specify both.

#### Example

```json
"rdma": {
    "mlx5_1": {
        "hcaHandles": 3,
        "hcaObjects": 10000
    },
    "mlx4_0": {
        "hcaObjects": 1000
    },
    "rxe3": {
        "hcaObjects": 10000
    }
}
```

## <a name="configLinuxUnified" />Unified

**`unified`** (object, OPTIONAL) allows cgroup v2 parameters to be to be set and modified for the container.

Each key in the map refers to a file in the cgroup unified hierarchy.

The OCI runtime MUST ensure that the needed cgroup controllers are enabled for the cgroup.

Configuration unknown to the runtime MUST still be written to the relevant file.

The runtime MUST generate an error when the configuration refers to a cgroup controller that is not present or that cannot be enabled.

### Example

```json
"unified": {
    "io.max": "259:0 rbps=2097152 wiops=120\n253:0 rbps=2097152 wiops=120",
    "hugetlb.1GB.max": "1073741824"
}
```

If a controller is enabled on the cgroup v2 hierarchy but the configuration is provided for the cgroup v1 equivalent controller, the runtime MAY attempt a conversion.

If the conversion is not possible the runtime MUST generate an error.

## <a name="configLinuxIntelRdt" />IntelRdt

**`intelRdt`** (object, OPTIONAL) represents the [Intel Resource Director Technology][intel-rdt-cat-kernel-interface].
If `intelRdt` is set, the runtime MUST write the container process ID to the `tasks` file in a proper sub-directory in a mounted `resctrl` pseudo-filesystem. That sub-directory name is specified by `closID` parameter.
If no mounted `resctrl` pseudo-filesystem is available in the [runtime mount namespace](glossary.md#runtime-namespace), the runtime MUST [generate an error](runtime.md#errors).

If `intelRdt` is not set, the runtime MUST NOT manipulate any `resctrl` pseudo-filesystems.

The following parameters can be specified for the container:

* **`closID`** *(string, OPTIONAL)* - specifies the identity for RDT Class of Service (CLOS).

* **`l3CacheSchema`** *(string, OPTIONAL)* - specifies the schema for L3 cache id and capacity bitmask (CBM).
    The value SHOULD start with `L3:` and SHOULD NOT contain newlines.
* **`memBwSchema`** *(string, OPTIONAL)* - specifies the schema of memory bandwidth per L3 cache id.
    The value MUST start with `MB:` and MUST NOT contain newlines.

The following rules on parameters MUST be applied:

* If both `l3CacheSchema` and `memBwSchema` are set, runtimes MUST write the combined value to the `schemata` file in that sub-directory discussed in `closID`.

* If `l3CacheSchema` contains a line beginning with `MB:`, the value written to `schemata` file MUST be the non-`MB:` line(s) from `l3CacheSchema` and the line from `memBWSchema`.

* If either `l3CacheSchema` or `memBwSchema` is set, runtimes MUST write the value to the `schemata` file in the that sub-directory discussed in `closID`.

* If neither `l3CacheSchema` nor `memBwSchema` is set, runtimes MUST NOT write to `schemata` files in any `resctrl` pseudo-filesystems.

* If `closID` is not set, runtimes MUST use the container ID from [`start`](runtime.md#start) and create the `<container-id>` directory.

* If `closID` is set, `l3CacheSchema` and/or `memBwSchema` is set
  * if `closID` directory in a mounted `resctrl` pseudo-filesystem doesn't exist, the runtimes MUST create it.
  * if `closID` directory in a mounted `resctrl` pseudo-filesystem exists, runtimes MUST compare `l3CacheSchema` and/or `memBwSchema` value with `schemata` file, and [generate an error](runtime.md#errors) if doesn't match.

* If `closID` is set, and neither of `l3CacheSchema` and `memBwSchema` are set, runtime MUST check if corresponding pre-configured directory `closID` is present in mounted `resctrl`. If such pre-configured directory `closID` exists, runtime MUST assign container to this `closID` and [generate an error](runtime.md#errors) if directory does not exist.

* **`enableCMT`** *(boolean, OPTIONAL)* - specifies if Intel RDT CMT should be enabled:
    * CMT (Cache Monitoring Technology) supports monitoring of the last-level cache (LLC) occupancy
      for the container.

* **`enableMBM`** *(boolean, OPTIONAL)* - specifies if Intel RDT MBM should be enabled:
    * MBM (Memory Bandwidth Monitoring) supports monitoring of total and local memory bandwidth
      for the container.

### Example

Consider a two-socket machine with two L3 caches where the default CBM is 0x7ff and the max CBM length is 11 bits,
and minimum memory bandwidth of 10% with a memory bandwidth granularity of 10%.

Tasks inside the container only have access to the "upper" 7/11 of L3 cache on socket 0 and the "lower" 5/11 L3 cache on socket 1,
and may use a maximum memory bandwidth of 20% on socket 0 and 70% on socket 1.

```json
"linux": {
    "intelRdt": {
        "closID": "guaranteed_group",
        "l3CacheSchema": "L3:0=7f0;1=1f",
        "memBwSchema": "MB:0=20;1=70"
    }
}
```

## <a name="configLinuxSysctl" />Sysctl

**`sysctl`** (object, OPTIONAL) allows kernel parameters to be modified at runtime for the container.
For more information, see the [sysctl(8)][sysctl.8] man page.

### Example

```json
"sysctl": {
    "net.ipv4.ip_forward": "1",
    "net.core.somaxconn": "256"
}
```

## <a name="configLinuxSeccomp" />Seccomp

Seccomp provides application sandboxing mechanism in the Linux kernel.
Seccomp configuration allows one to configure actions to take for matched syscalls and furthermore also allows matching on values passed as arguments to syscalls.
For more information about Seccomp, see [Seccomp][seccomp] kernel documentation.
The actions, architectures, and operators are strings that match the definitions in seccomp.h from [libseccomp][] and are translated to corresponding values.

**`seccomp`** (object, OPTIONAL)

The following parameters can be specified to set up seccomp:

* **`defaultAction`** *(string, REQUIRED)* - the default action for seccomp. Allowed values are the same as `syscalls[].action`.
* **`defaultErrnoRet`** *(uint, OPTIONAL)* - the errno return code to use.
    Some actions like `SCMP_ACT_ERRNO` and `SCMP_ACT_TRACE` allow to specify the errno code to return.
    When the action doesn't support an errno, the runtime MUST print and error and fail.
    If not specified then its default value is `EPERM`.
* **`architectures`** *(array of strings, OPTIONAL)* - the architecture used for system calls.
    A valid list of constants as of libseccomp v2.5.0 is shown below.

    * `SCMP_ARCH_X86`
    * `SCMP_ARCH_X86_64`
    * `SCMP_ARCH_X32`
    * `SCMP_ARCH_ARM`
    * `SCMP_ARCH_AARCH64`
    * `SCMP_ARCH_MIPS`
    * `SCMP_ARCH_MIPS64`
    * `SCMP_ARCH_MIPS64N32`
    * `SCMP_ARCH_MIPSEL`
    * `SCMP_ARCH_MIPSEL64`
    * `SCMP_ARCH_MIPSEL64N32`
    * `SCMP_ARCH_PPC`
    * `SCMP_ARCH_PPC64`
    * `SCMP_ARCH_PPC64LE`
    * `SCMP_ARCH_S390`
    * `SCMP_ARCH_S390X`
    * `SCMP_ARCH_PARISC`
    * `SCMP_ARCH_PARISC64`
    * `SCMP_ARCH_RISCV64`

* **`flags`** *(array of strings, OPTIONAL)* - list of flags to use with seccomp(2).

    A valid list of constants is shown below.

    * `SECCOMP_FILTER_FLAG_TSYNC`
    * `SECCOMP_FILTER_FLAG_LOG`
    * `SECCOMP_FILTER_FLAG_SPEC_ALLOW`
    * `SECCOMP_FILTER_FLAG_WAIT_KILLABLE_RECV`

* **`listenerPath`** *(string, OPTIONAL)* - specifies the path of UNIX domain socket over which the runtime will send the [container process state](#containerprocessstate) data structure when the `SCMP_ACT_NOTIFY` action is used.
    This socket MUST use `AF_UNIX` domain and `SOCK_STREAM` type.
    The runtime MUST send exactly one [container process state](#containerprocessstate) per connection.
    The connection MUST NOT be reused and it MUST be closed after sending a seccomp state.
    If sending to this socket fails, the runtime MUST [generate an error](runtime.md#errors).
    If the `SCMP_ACT_NOTIFY` action is not used this value is ignored.

    The runtime sends the following file descriptors using `SCM_RIGHTS` and set their names in the `fds` array of the [container process state](#containerprocessstate):

    * **`seccompFd`** (string, REQUIRED) is the seccomp file descriptor returned by the seccomp syscall.

* **`listenerMetadata`** *(string, OPTIONAL)* - specifies an opaque data to pass to the seccomp agent.
    This string will be sent as the `metadata` field in the [container process state](#containerprocessstate).
    This field MUST NOT be set if `listenerPath` is not set.

* **`syscalls`** *(array of objects, OPTIONAL)* - match a syscall in seccomp.
    While this property is OPTIONAL, some values of `defaultAction` are not useful without `syscalls` entries.
    For example, if `defaultAction` is `SCMP_ACT_KILL` and `syscalls` is empty or unset, the kernel will kill the container process on its first syscall.
    Each entry has the following structure:

    * **`names`** *(array of strings, REQUIRED)* - the names of the syscalls.
        `names` MUST contain at least one entry.
    * **`action`** *(string, REQUIRED)* - the action for seccomp rules.
        A valid list of constants as of libseccomp v2.5.0 is shown below.

        * `SCMP_ACT_KILL`
        * `SCMP_ACT_KILL_PROCESS`
        * `SCMP_ACT_KILL_THREAD`
        * `SCMP_ACT_TRAP`
        * `SCMP_ACT_ERRNO`
        * `SCMP_ACT_TRACE`
        * `SCMP_ACT_ALLOW`
        * `SCMP_ACT_LOG`
        * `SCMP_ACT_NOTIFY`

    * **`errnoRet`** *(uint, OPTIONAL)* - the errno return code to use.
        Some actions like `SCMP_ACT_ERRNO` and `SCMP_ACT_TRACE` allow to specify the errno code to return.
        When the action doesn't support an errno, the runtime MUST print and error and fail.
        If not specified its default value is `EPERM`.

    * **`args`** *(array of objects, OPTIONAL)* - the specific syscall in seccomp.
        Each entry has the following structure:

        * **`index`** *(uint, REQUIRED)* - the index for syscall arguments in seccomp.
        * **`value`** *(uint64, REQUIRED)* - the value for syscall arguments in seccomp.
        * **`valueTwo`** *(uint64, OPTIONAL)* - the value for syscall arguments in seccomp.
        * **`op`** *(string, REQUIRED)* - the operator for syscall arguments in seccomp.
            A valid list of constants as of libseccomp v2.3.2 is shown below.

            * `SCMP_CMP_NE`
            * `SCMP_CMP_LT`
            * `SCMP_CMP_LE`
            * `SCMP_CMP_EQ`
            * `SCMP_CMP_GE`
            * `SCMP_CMP_GT`
            * `SCMP_CMP_MASKED_EQ`

### Example

```json
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
}
```

### <a name="containerprocessstate" />The Container Process State

The container process state is a data structure passed via a UNIX socket.
The container runtime MUST send the container process state over the UNIX socket as regular payload serialized in JSON and file descriptors MUST be sent using `SCM_RIGHTS`.
The container runtime MAY use several `sendmsg(2)` calls to send the aforementioned data.
If more than one `sendmsg(2)` is used, the file descriptors MUST be sent only in the first call.

The container process state includes the following properties:

* **`ociVersion`** (string, REQUIRED) is version of the Open Container Initiative Runtime Specification with which the container process state complies.
* **`fds`** (array, OPTIONAL) is a string array containing the names of the file descriptors passed.
    The index of the name in this array corresponds to index of the file descriptors in the `SCM_RIGHTS` array.
* **`pid`** (int, REQUIRED) is the container process ID, as seen by the runtime.
* **`metadata`** (string, OPTIONAL) opaque metadata.
* **`state`** ([state](runtime.md#state), REQUIRED) is the state of the container.

Example sending a single `seccompFd` file descriptor in the `SCM_RIGHTS` array:

```json
{
    "ociVersion": "1.0.2",
    "fds": [
        "seccompFd"
    ],
    "pid": 4422,
    "metadata": "MKNOD=/dev/null,/dev/net/tun;BPF_MAP_TYPES=hash,array",
    "state": {
        "ociVersion": "1.0.2",
        "id": "oci-container1",
        "status": "creating",
        "pid": 4422,
        "bundle": "/containers/redis",
        "annotations": {
            "myKey": "myValue"
        }
    }
}
```

## <a name="configLinuxRootfsMountPropagation" />Rootfs Mount Propagation

**`rootfsPropagation`** (string, OPTIONAL) sets the rootfs's mount propagation.
Its value is either `shared`, `slave`, `private` or `unbindable`.
It's worth noting that a peer group is defined as a group of VFS mounts that propagate events to each other.
A nested container is defined as a container launched inside an existing container.

* **`shared`**: the rootfs mount belongs to a new peer group.
    This means that further mounts (e.g. nested containers) will also belong to that peer group and will propagate events to the rootfs.
    Note this does not mean that it's shared with the host.
* **`slave`**: the rootfs mount receives propagation events from the host (e.g. if something is mounted on the host it will also appear in the container) but not the other way around.
* **`private`**: the rootfs mount doesn't receive mount propagation events from the host and further mounts in nested containers will be isolated from the host and from the rootfs (even if the nested container `rootfsPropagation` option is shared).
* **`unbindable`**: the rootfs mount is a private mount that cannot be bind-mounted.

The [Shared Subtrees][sharedsubtree] article in the kernel documentation has more information about mount propagation.

### Example

```json
"rootfsPropagation": "slave",
```

## <a name="configLinuxMaskedPaths" />Masked Paths

**`maskedPaths`** (array of strings, OPTIONAL) will mask over the provided paths inside the container so that they cannot be read.
The values MUST be absolute paths in the [container namespace](glossary.md#container_namespace).

### Example

```json
"maskedPaths": [
    "/proc/kcore"
]
```

## <a name="configLinuxReadonlyPaths" />Readonly Paths

**`readonlyPaths`** (array of strings, OPTIONAL) will set the provided paths as readonly inside the container.
The values MUST be absolute paths in the [container namespace](glossary.md#container-namespace).

### Example

```json
"readonlyPaths": [
    "/proc/sys"
]
```

## <a name="configLinuxMountLabel" />Mount Label

**`mountLabel`** (string, OPTIONAL) will set the Selinux context for the mounts in the container.

### Example

```json
"mountLabel": "system_u:object_r:svirt_sandbox_file_t:s0:c715,c811"
```

## <a name="configLinuxPersonality" />Personality

**`personality`** (object, OPTIONAL) sets the Linux execution personality. For more information
see the [personality][personality.2] syscall documentation. As most of the options are
obsolete and rarely used, and some reduce security, the currently supported set is a small
subset of the available options.

* **`domain`** *(string, REQUIRED)* - the execution domain.
    The valid list of constants is shown below. `LINUX32` will set the `uname` system call to show
    a 32 bit CPU type, such as `i686`.

    * `LINUX`
    * `LINUX32`

* **`flags`** *(array of strings, OPTIONAL)* - the additional flags to apply.
    Currently no flag values are supported.


[cgroup-v1]: https://www.kernel.org/doc/Documentation/cgroup-v1/cgroups.txt
[cgroup-v1-blkio]: https://www.kernel.org/doc/Documentation/cgroup-v1/blkio-controller.txt
[cgroup-v1-cpusets]: https://www.kernel.org/doc/Documentation/cgroup-v1/cpusets.txt
[cgroup-v1-devices]: https://www.kernel.org/doc/Documentation/cgroup-v1/devices.txt
[cgroup-v1-hugetlb]: https://www.kernel.org/doc/Documentation/cgroup-v1/hugetlb.txt
[cgroup-v1-memory]: https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
[cgroup-v1-net-cls]: https://www.kernel.org/doc/Documentation/cgroup-v1/net_cls.txt
[cgroup-v1-net-prio]: https://www.kernel.org/doc/Documentation/cgroup-v1/net_prio.txt
[cgroup-v1-pids]: https://www.kernel.org/doc/Documentation/cgroup-v1/pids.txt
[cgroup-v1-rdma]: https://www.kernel.org/doc/Documentation/cgroup-v1/rdma.txt
[cgroup-v2]: https://www.kernel.org/doc/Documentation/cgroup-v2.txt
[cgroup-v2-io]: https://docs.kernel.org/admin-guide/cgroup-v2.html#io
[devices]: https://www.kernel.org/doc/Documentation/admin-guide/devices.txt
[devpts]: https://www.kernel.org/doc/Documentation/filesystems/devpts.txt
[file]: http://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap03.html#tag_03_164
[libseccomp]: https://github.com/seccomp/libseccomp
[proc]: https://www.kernel.org/doc/Documentation/filesystems/proc.txt
[seccomp]: https://www.kernel.org/doc/Documentation/prctl/seccomp_filter.txt
[sharedsubtree]: https://www.kernel.org/doc/Documentation/filesystems/sharedsubtree.txt
[sysfs]: https://www.kernel.org/doc/Documentation/filesystems/sysfs.txt
[tmpfs]: https://www.kernel.org/doc/Documentation/filesystems/tmpfs.txt

[full.4]: http://man7.org/linux/man-pages/man4/full.4.html
[mknod.1]: http://man7.org/linux/man-pages/man1/mknod.1.html
[mknod.2]: http://man7.org/linux/man-pages/man2/mknod.2.html
[namespaces.7_2]: http://man7.org/linux/man-pages/man7/namespaces.7.html
[null.4]: http://man7.org/linux/man-pages/man4/null.4.html
[personality.2]: http://man7.org/linux/man-pages/man2/personality.2.html
[pts.4]: http://man7.org/linux/man-pages/man4/pts.4.html
[random.4]: http://man7.org/linux/man-pages/man4/random.4.html
[sysctl.8]: http://man7.org/linux/man-pages/man8/sysctl.8.html
[tty.4]: http://man7.org/linux/man-pages/man4/tty.4.html
[zero.4]: http://man7.org/linux/man-pages/man4/zero.4.html
[user-namespaces]: http://man7.org/linux/man-pages/man7/user_namespaces.7.html
[intel-rdt-cat-kernel-interface]: https://www.kernel.org/doc/Documentation/x86/intel_rdt_ui.txt
[time_namespaces.7]: https://man7.org/linux/man-pages/man7/time_namespaces.7.html
