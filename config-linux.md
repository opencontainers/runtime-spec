# <a name="linuxContainerConfiguration" />Linux Container Configuration

This document describes the schema for the [Linux-specific section](config.md#platform-specific-configuration) of the [container configuration](config.md).
The Linux container specification uses various kernel features like namespaces, cgroups, capabilities, LSM, and filesystem jails to fulfill the spec.

## <a name="configLinuxDefaultFilesystems" />Default Filesystems

The Linux ABI includes both syscalls and several special file paths.
Applications expecting a Linux environment will very likely expect these file paths to be setup correctly.

The following filesystems SHOULD be made available in each container's filesystem:

| Path     | Type   |
| -------- | ------ |
| /proc    | [procfs][procfs]   |
| /sys     | [sysfs][sysfs]     |
| /dev/pts | [devpts][devpts]   |
| /dev/shm | [tmpfs][tmpfs]     |

## <a name="configLinuxNamespaces" />Namespaces

A namespace wraps a global system resource in an abstraction that makes it appear to the processes within the namespace that they have their own isolated instance of the global resource.
Changes to the global resource are visible to other processes that are members of the namespace, but are invisible to other processes.
For more information, see the [namespaces(7)][namespaces.7_2] man page.

Namespaces are specified as an array of entries inside the `namespaces` root field.
The following parameters can be specified to setup namespaces:

* **`type`** *(string, REQUIRED, linux)* - namespace type. The following namespace types are supported:
    * **`pid`** processes inside the container will only be able to see other processes inside the same container.
    * **`network`** the container will have its own network stack.
    * **`mount`** the container will have an isolated mount table.
    * **`ipc`** processes inside the container will only be able to communicate to other processes inside the same container via system level IPC.
    * **`uts`** the container will be able to have its own hostname and domain name.
    * **`user`** the container will be able to remap user and group IDs from the host to local users and groups within the container.
    * **`cgroup`** the container will have an isolated view of the cgroup hierarchy.

* **`path`** *(string, OPTIONAL, linux)* - an absolute path to namespace file in the [runtime mount namespace](glossary.md#runtime-namespace)

If a path is specified, that particular file is used to join that type of namespace.
If a namespace type is not specified in the `namespaces` array, the container MUST inherit the [runtime namespace](glossary.md#runtime-namespace) of that type.
If a `namespaces` field contains duplicated namespaces with same `type`, the runtime MUST error out.

###### Example

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
        }
    ]
```

## <a name="configLinuxUserNamespaceMappings" />User namespace mappings

**`uidMappings`** (array of objects, OPTIONAL, linux) describes the user namespace uid mappings from the host to the container.
**`gidMappings`** (array of objects, OPTIONAL, linux) describes the user namespace gid mappings from the host to the container.

Each entry has the following structure:

* **`hostID`** *(uint32, REQUIRED, linux)* - is the starting uid/gid on the host to be mapped to *containerID*.
* **`containerID`** *(uint32, REQUIRED, linux)* - is the starting uid/gid in the container.
* **`size`** *(uint32, REQUIRED, linux)* - is the number of ids to be mapped.

The runtime SHOULD NOT modify the ownership of referenced filesystems to realize the mapping.
Note that the number of mapping entries MAY be limited by the [kernel][user-namespaces].

###### Example

```json
    "uidMappings": [
        {
            "hostID": 1000,
            "containerID": 0,
            "size": 32000
        }
    ],
    "gidMappings": [
        {
            "hostID": 1000,
            "containerID": 0,
            "size": 32000
        }
    ]
```

## <a name="configLinuxDevices" />Devices

**`devices`** (array of objects, OPTIONAL, linux) lists devices that MUST be available in the container.
The runtime may supply them however it likes (with [mknod][mknod.2], by bind mounting from the runtime mount namespace, etc.).

Each entry has the following structure:

* **`type`** *(string, REQUIRED, linux)* - type of device: `c`, `b`, `u` or `p`.
  More info in [mknod(1)][mknod.1].
* **`path`** *(string, REQUIRED, linux)* - full path to device inside container.
  If a [file][file.1] already exists at `path` that does not match the requested device, the runtime MUST generate an error.
* **`major, minor`** *(int64, REQUIRED unless `type` is `p`, linux)* - [major, minor numbers][devices] for the device.
* **`fileMode`** *(uint32, OPTIONAL, linux)* - file mode for the device.
  You can also control access to devices [with cgroups](#device-whitelist).
* **`uid`** *(uint32, OPTIONAL, linux)* - id of device owner.
* **`gid`** *(uint32, OPTIONAL, linux)* - id of device group.

The same `type`, `major` and `minor` SHOULD NOT be used for multiple devices.

###### Example

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

###### <a name="configLinuxDefaultDevices" />Default Devices

In addition to any devices configured with this setting, the runtime MUST also supply:

* [`/dev/null`][null.4]
* [`/dev/zero`][zero.4]
* [`/dev/full`][full.4]
* [`/dev/random`][random.4]
* [`/dev/urandom`][random.4]
* [`/dev/tty`][tty.4]
* [`/dev/console`][console.4] is setup if terminal is enabled in the config by bind mounting the pseudoterminal slave to /dev/console.
* [`/dev/ptmx`][pts.4].
  A [bind-mount or symlink of the container's `/dev/pts/ptmx`][devpts].

## <a name="configLinuxControlGroups" />Control groups

Also known as cgroups, they are used to restrict resource usage for a container and handle device access.
cgroups provide controls (through controllers) to restrict cpu, memory, IO, pids and network for the container.
For more information, see the [kernel cgroups documentation][cgroup-v1].

The path to the cgroups can be specified in the Spec via `cgroupsPath`.
`cgroupsPath` can be used to either control the cgroup hierarchy for containers or to run a new process in an existing container.
If `cgroupsPath` is:
* ... an absolute path (starting with `/`), the runtime MUST take the path to be relative to the cgroup mount point.
* ... a relative path (not starting with `/`), the runtime MAY interpret the path relative to a runtime-determined location in the cgroup hierarchy.
* ... not specified, the runtime MAY define the default cgroup path.
Runtimes MAY consider certain `cgroupsPath` values to be invalid, and MUST generate an error if this is the case.
If a `cgroupsPath` value is specified, the runtime MUST consistently attach to the same place in the cgroup hierarchy given the same value of `cgroupsPath`.

Implementations of the Spec can choose to name cgroups in any manner.
The Spec does not include naming schema for cgroups.
The Spec does not support per-controller paths for the reasons discussed in the [cgroupv2 documentation][cgroup-v2].
The cgroups will be created if they don't exist.

You can configure a container's cgroups via the `resources` field of the Linux configuration.
Do not specify `resources` unless limits have to be updated.
For example, to run a new process in an existing container without updating limits, `resources` need not be specified.

A runtime MUST at least use the minimum set of cgroup controllers required to fulfill the `resources` settings.
However, a runtime MAY attach the container process to additional cgroup controllers supported by the system.

###### Example

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

#### <a name="configLinuxDeviceWhitelist" />Device whitelist

**`devices`** (array of objects, OPTIONAL, linux) configures the [device whitelist][cgroup-v1-devices].
The runtime MUST apply entries in the listed order.

Each entry has the following structure:

* **`allow`** *(boolean, REQUIRED, linux)* - whether the entry is allowed or denied.
* **`type`** *(string, OPTIONAL, linux)* - type of device: `a` (all), `c` (char), or `b` (block).
  `null` or unset values mean "all", mapping to `a`.
* **`major, minor`** *(int64, OPTIONAL, linux)* - [major, minor numbers][devices] for the device.
  `null` or unset values mean "all", mapping to [`*` in the filesystem API][cgroup-v1-devices].
* **`access`** *(string, OPTIONAL, linux)* - cgroup permissions for device.
  A composition of `r` (read), `w` (write), and `m` (mknod).

###### Example

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

#### <a name="configLinuxDisableOutOfMemoryKiller" />Disable out-of-memory killer

`disableOOMKiller` contains a boolean (`true` or `false`) that enables or disables the Out of Memory killer for a cgroup.
If enabled (`false`), tasks that attempt to consume more memory than they are allowed are immediately killed by the OOM killer.
The OOM killer is enabled by default in every cgroup using the `memory` subsystem.
To disable it, specify a value of `true`.
For more information, see [the memory cgroup man page][cgroup-v1-memory].

* **`disableOOMKiller`** *(bool, OPTIONAL, linux)* - enables or disables the OOM killer

###### Example

```json
    "disableOOMKiller": false
```

#### <a name="configLinuxSetOomScoreAdj" />Set oom_score_adj

`oomScoreAdj` sets heuristic regarding how the process is evaluated by the kernel during memory pressure.
For more information, see [the proc filesystem documentation section 3.1][procfs].
This is a kernel/system level setting, where as `disableOOMKiller` is scoped for a memory cgroup.
For more information on how these two settings work together, see [the memory cgroup documentation section 10. OOM Contol][cgroup-v1-memory].

* **`oomScoreAdj`** *(int, OPTIONAL, linux)* - adjust the oom-killer score

###### Example

```json
    "oomScoreAdj": 100
```

#### <a name="configLinuxMemory" />Memory

**`memory`** (object, OPTIONAL, linux) represents the cgroup subsystem `memory` and it's used to set limits on the container's memory usage.
For more information, see [the memory cgroup man page][cgroup-v1-memory].

The following parameters can be specified to setup the controller:

* **`limit`** *(uint64, OPTIONAL, linux)* - sets limit of memory usage in bytes

* **`reservation`** *(uint64, OPTIONAL, linux)* - sets soft limit of memory usage in bytes

* **`swap`** *(uint64, OPTIONAL, linux)* - sets limit of memory+Swap usage

* **`kernel`** *(uint64, OPTIONAL, linux)* - sets hard limit for kernel memory

* **`kernelTCP`** *(uint64, OPTIONAL, linux)* - sets hard limit in bytes for kernel TCP buffer memory

* **`swappiness`** *(uint64, OPTIONAL, linux)* - sets swappiness parameter of vmscan (See sysctl's vm.swappiness)

###### Example

```json
    "memory": {
        "limit": 536870912,
        "reservation": 536870912,
        "swap": 536870912,
        "kernel": 0,
        "kernelTCP": 0,
        "swappiness": 0
    }
```

#### <a name="configLinuxCPU" />CPU

**`cpu`** (object, OPTIONAL, linux) represents the cgroup subsystems `cpu` and `cpusets`.
For more information, see [the cpusets cgroup man page][cgroup-v1-cpusets].

The following parameters can be specified to setup the controller:

* **`shares`** *(uint64, OPTIONAL, linux)* - specifies a relative share of CPU time available to the tasks in a cgroup

* **`quota`** *(int64, OPTIONAL, linux)* - specifies the total amount of time in microseconds for which all tasks in a cgroup can run during one period (as defined by **`period`** below)

* **`period`** *(uint64, OPTIONAL, linux)* - specifies a period of time in microseconds for how regularly a cgroup's access to CPU resources should be reallocated (CFS scheduler only)

* **`realtimeRuntime`** *(int64, OPTIONAL, linux)* - specifies a period of time in microseconds for the longest continuous period in which the tasks in a cgroup have access to CPU resources

* **`realtimePeriod`** *(uint64, OPTIONAL, linux)* - same as **`period`** but applies to realtime scheduler only

* **`cpus`** *(string, OPTIONAL, linux)* - list of CPUs the container will run in

* **`mems`** *(string, OPTIONAL, linux)* - list of Memory Nodes the container will run in

###### Example

```json
    "cpu": {
        "shares": 1024,
        "quota": 1000000,
        "period": 500000,
        "realtimeRuntime": 950000,
        "realtimePeriod": 1000000,
        "cpus": "2-3",
        "mems": "0-7"
    }
```

#### <a name="configLinuxBlockIO" />Block IO

**`blockIO`** (object, OPTIONAL, linux) represents the cgroup subsystem `blkio` which implements the block IO controller.
For more information, see [the kernel cgroups documentation about blkio][cgroup-v1-blkio].

The following parameters can be specified to setup the controller:

* **`blkioWeight`** *(uint16, OPTIONAL, linux)* - specifies per-cgroup weight. This is default weight of the group on all devices until and unless overridden by per-device rules. The range is from 10 to 1000.

* **`blkioLeafWeight`** *(uint16, OPTIONAL, linux)* - equivalents of `blkioWeight` for the purpose of deciding how much weight tasks in the given cgroup has while competing with the cgroup's child cgroups. The range is from 10 to 1000.

* **`blkioWeightDevice`** *(array of objects, OPTIONAL, linux)* - specifies the list of devices which will be bandwidth rate limited. The following parameters can be specified per-device:
    * **`major, minor`** *(int64, REQUIRED, linux)* - major, minor numbers for device. More info in `man mknod`.
    * **`weight`** *(uint16, OPTIONAL, linux)* - bandwidth rate for the device, range is from 10 to 1000
    * **`leafWeight`** *(uint16, OPTIONAL, linux)* - bandwidth rate for the device while competing with the cgroup's child cgroups, range is from 10 to 1000, CFQ scheduler only

    You MUST specify at least one of `weight` or `leafWeight` in a given entry, and MAY specify both.

* **`blkioThrottleReadBpsDevice`**, **`blkioThrottleWriteBpsDevice`**, **`blkioThrottleReadIOPSDevice`**, **`blkioThrottleWriteIOPSDevice`** *(array of objects, OPTIONAL, linux)* - specify the list of devices which will be IO rate limited.
  The following parameters can be specified per-device:
    * **`major, minor`** *(int64, REQUIRED, linux)* - major, minor numbers for device. More info in `man mknod`.
    * **`rate`** *(uint64, REQUIRED, linux)* - IO rate limit for the device

###### Example

```json
    "blockIO": {
        "blkioWeight": 10,
        "blkioLeafWeight": 10,
        "blkioWeightDevice": [
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
        "blkioThrottleReadBpsDevice": [
            {
                "major": 8,
                "minor": 0,
                "rate": 600
            }
        ],
        "blkioThrottleWriteIOPSDevice": [
            {
                "major": 8,
                "minor": 16,
                "rate": 300
            }
        ]
    }
```

#### <a name="configLinuxHugePageLimits" />Huge page limits

**`hugepageLimits`** (array of objects, OPTIONAL, linux) represents the `hugetlb` controller which allows to limit the
HugeTLB usage per control group and enforces the controller limit during page fault.
For more information, see the [kernel cgroups documentation about HugeTLB][cgroup-v1-hugetlb].

Each entry has the following structure:

* **`pageSize`** *(string, REQUIRED, linux)* - hugepage size

* **`limit`** *(uint64, REQUIRED, linux)* - limit in bytes of *hugepagesize* HugeTLB usage

###### Example

```json
   "hugepageLimits": [
        {
            "pageSize": "2MB",
            "limit": 209715200
        }
   ]
```

#### <a name="configLinuxNetwork" />Network

**`network`** (object, OPTIONAL, linux) represents the cgroup subsystems `net_cls` and `net_prio`.
For more information, see [the net\_cls cgroup man page][cgroup-v1-net-cls] and [the net\_prio cgroup man page][cgroup-v1-net-prio].

The following parameters can be specified to setup the controller:

* **`classID`** *(uint32, OPTIONAL, linux)* - is the network class identifier the cgroup's network packets will be tagged with
* **`priorities`** *(array of objects, OPTIONAL, linux)* - specifies a list of objects of the priorities assigned to traffic originating from processes in the group and egressing the system on various interfaces.
  The following parameters can be specified per-priority:
    * **`name`** *(string, REQUIRED, linux)* - interface name
    * **`priority`** *(uint32, REQUIRED, linux)* - priority applied to the interface

###### Example

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

#### <a name="configLinuxPIDS" />PIDs

**`pids`** (object, OPTIONAL, linux) represents the cgroup subsystem `pids`.
For more information, see [the pids cgroup man page][cgroup-v1-pids].

The following parameters can be specified to setup the controller:

* **`limit`** *(int64, REQUIRED, linux)* - specifies the maximum number of tasks in the cgroup

###### Example

```json
   "pids": {
        "limit": 32771
   }
```

## <a name="configLinuxIntelRdt" />IntelRdt

Intel platforms with new Xeon CPU support Intel Resource Director Technology
(RDT). Cache Allocation Technology (CAT) is a sub-feature of RDT, which
currently supports L3 cache resource allocation.

This feature provides a way for the software to restrict cache allocation to a
defined 'subset' of L3 cache which may be overlapping with other 'subsets'.
The different subsets are identified by class of service (CLOS) and each CLOS
has a capacity bitmask (CBM).

In Linux kernel, it is exposed via "resource control" filesystem, which is a
"cgroup-like" interface.

Comparing with cgroups, it has similar process management lifecycle and
interfaces in a container. But unlike cgroups' hierarchy, it has single level
filesystem layout.

Intel RDT "resource control" filesystem hierarchy:
```
mount -t resctrl resctrl /sys/fs/resctrl
tree /sys/fs/resctrl
/sys/fs/resctrl/
|-- info
|   |-- L3
|       |-- cbm_mask
|       |-- min_cbm_bits
|       |-- num_closids
|-- cpus
|-- schemata
|-- tasks
|-- <container_id>
    |-- cpus
    |-- schemata
    |-- tasks

```

For containers, we can make use of `tasks` and `schemata` configuration for
L3 cache resource constraints if hardware and kernel support Intel RDT/CAT.

The file `tasks` has a list of tasks that belongs to this group (e.g.,
<container_id>" group). Tasks can be added to a group by writing the task ID
to the "tasks" file  (which will automatically remove them from the previous
group to which they belonged). New tasks created by fork(2) and clone(2) are
added to the same group as their parent. If a pid is not in any sub group, it
is in root group.

The file `schemata` has allocation masks/values for L3 cache on each socket,
which contains L3 cache id and capacity bitmask (CBM).
```
	Format: "L3:<cache_id0>=<cbm0>;<cache_id1>=<cbm1>;..."
```
For example, on a two-socket machine, L3's schema line could be `L3:0=ff;1=c0`
Which means L3 cache id 0's CBM is 0xff, and L3 cache id 1's CBM is 0xc0.

The valid L3 cache CBM is a *contiguous bits set* and number of bits that can
be set is less than the max bit. The max bits in the CBM is varied among
supported Intel Xeon platforms. In Intel RDT "resource control" filesystem
layout, the CBM in a group should be a subset of the CBM in root. Kernel will
check if it is valid when writing. e.g., 0xfffff in root indicates the max bits
of CBM is 20 bits, which mapping to entire L3 cache capacity. Some valid CBM
values to set in a group: 0xf, 0xf0, 0x3ff, 0x1f00 and etc.

**`intelRdt`** (object, OPTIONAL, linux) represents the L3 cache resource constraints in Intel Xeon platforms.

For more information, see [Intel RDT/CAT kernel interface][intel-rdt-cat-kernel-interface].

The following parameters can be specified for the container:

* **`l3CacheSchema`** *(string, OPTIONAL, linux)* - specifies the schema for L3 cache id and capacity bitmask (CBM)

###### Example
```json
There are two L3 caches in the two-socket machine, the default CBM is 0xfffff
and the max CBM length is 20 bits. This configuration assigns 4/5 of L3 cache
id 0 and the whole L3 cache id 1 for the container:

"linux": {
	"intelRdt": {
		"l3CacheSchema": "L3:0=ffff0;1=fffff"
	}
}
```

## <a name="configLinuxSysctl" />Sysctl

**`sysctl`** (object, OPTIONAL, linux) allows kernel parameters to be modified at runtime for the container.
For more information, see the [sysctl(8)][sysctl.8] man page.

###### Example

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

**`seccomp`** (object, OPTIONAL, linux)

The following parameters can be specified to setup seccomp:

* **`defaultAction`** *(string, REQUIRED, linux)* - the default action for seccomp. Allowed values are the same as `syscalls[].action`.

* **`architectures`** *(array of strings, OPTIONAL, linux)* - the architecture used for system calls.
    A valid list of constants as of libseccomp v2.3.2 is shown below.

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

* **`syscalls`** *(array of objects, REQUIRED, linux)* - match a syscall in seccomp.

    Each entry has the following structure:

    * **`names`** *(array of strings, REQUIRED, linux)* - the names of the syscalls.

    * **`action`** *(string, REQUIRED, linux)* - the action for seccomp rules.
        A valid list of constants as of libseccomp v2.3.2 is shown below.

        * `SCMP_ACT_KILL`
        * `SCMP_ACT_TRAP`
        * `SCMP_ACT_ERRNO`
        * `SCMP_ACT_TRACE`
        * `SCMP_ACT_ALLOW`

    * **`args`** *(array of objects, OPTIONAL, linux)* - the specific syscall in seccomp.

        Each entry has the following structure:

        * **`index`** *(uint, REQUIRED, linux)* - the index for syscall arguments in seccomp.

        * **`value`** *(uint64, REQUIRED, linux)* - the value for syscall arguments in seccomp.

        * **`valueTwo`** *(uint64, REQUIRED, linux)* - the value for syscall arguments in seccomp.

        * **`op`** *(string, REQUIRED, linux)* - the operator for syscall arguments in seccomp.
            A valid list of constants as of libseccomp v2.3.2 is shown below.

            * `SCMP_CMP_NE`
            * `SCMP_CMP_LT`
            * `SCMP_CMP_LE`
            * `SCMP_CMP_EQ`
            * `SCMP_CMP_GE`
            * `SCMP_CMP_GT`
            * `SCMP_CMP_MASKED_EQ`

###### Example

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

## <a name="configLinuxRootfsMountPropagation" />Rootfs Mount Propagation

**`rootfsPropagation`** (string, OPTIONAL, linux) sets the rootfs's mount propagation.
Its value is either slave, private, or shared.
The [Shared Subtrees][sharedsubtree] article in the kernel documentation has more information about mount propagation.

###### Example

```json
    "rootfsPropagation": "slave",
```

## <a name="configLinuxMaskedPaths" />Masked Paths

**`maskedPaths`** (array of strings, OPTIONAL, linux) will mask over the provided paths inside the container so that they cannot be read.
The values MUST be absolute paths in the [container namespace][container-namespace2].

###### Example

```json
    "maskedPaths": [
        "/proc/kcore"
    ]
```

## <a name="configLinuxReadonlyPaths" />Readonly Paths

**`readonlyPaths`** (array of strings, OPTIONAL, linux) will set the provided paths as readonly inside the container.
The values MUST be absolute paths in the [container namespace][container-namespace2].

###### Example

```json
    "readonlyPaths": [
        "/proc/sys"
    ]
```

## <a name="configLinuxMountLabel" />Mount Label

**`mountLabel`** (string, OPTIONAL, linux) will set the Selinux context for the mounts in the container.

###### Example

```json
    "mountLabel": "system_u:object_r:svirt_sandbox_file_t:s0:c715,c811"
```


[container-namespace2]: glossary.md#container_namespace

[cgroup-v1]: https://www.kernel.org/doc/Documentation/cgroup-v1/cgroups.txt
[cgroup-v1-blkio]: https://www.kernel.org/doc/Documentation/cgroup-v1/blkio-controller.txt
[cgroup-v1-cpusets]: https://www.kernel.org/doc/Documentation/cgroup-v1/cpusets.txt
[cgroup-v1-devices]: https://www.kernel.org/doc/Documentation/cgroup-v1/devices.txt
[cgroup-v1-hugetlb]: https://www.kernel.org/doc/Documentation/cgroup-v1/hugetlb.txt
[cgroup-v1-memory]: https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
[cgroup-v1-net-cls]: https://www.kernel.org/doc/Documentation/cgroup-v1/net_cls.txt
[cgroup-v1-net-prio]: https://www.kernel.org/doc/Documentation/cgroup-v1/net_prio.txt
[cgroup-v1-pids]: https://www.kernel.org/doc/Documentation/cgroup-v1/pids.txt
[cgroup-v2]: https://www.kernel.org/doc/Documentation/cgroup-v2.txt
[devices]: https://www.kernel.org/doc/Documentation/devices.txt
[devpts]: https://www.kernel.org/doc/Documentation/filesystems/devpts.txt
[file]: http://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap03.html#tag_03_164
[libseccomp]: https://github.com/seccomp/libseccomp
[procfs]: https://www.kernel.org/doc/Documentation/filesystems/proc.txt
[seccomp]: https://www.kernel.org/doc/Documentation/prctl/seccomp_filter.txt
[sharedsubtree]: https://www.kernel.org/doc/Documentation/filesystems/sharedsubtree.txt
[sysfs]: https://www.kernel.org/doc/Documentation/filesystems/sysfs.txt
[tmpfs]: https://www.kernel.org/doc/Documentation/filesystems/tmpfs.txt

[console.4]: http://man7.org/linux/man-pages/man4/console.4.html
[full.4]: http://man7.org/linux/man-pages/man4/full.4.html
[mknod.1]: http://man7.org/linux/man-pages/man1/mknod.1.html
[mknod.2]: http://man7.org/linux/man-pages/man2/mknod.2.html
[namespaces.7_2]: http://man7.org/linux/man-pages/man7/namespaces.7.html
[null.4]: http://man7.org/linux/man-pages/man4/null.4.html
[pts.4]: http://man7.org/linux/man-pages/man4/pts.4.html
[random.4]: http://man7.org/linux/man-pages/man4/random.4.html
[sysctl.8]: http://man7.org/linux/man-pages/man8/sysctl.8.html
[tty.4]: http://man7.org/linux/man-pages/man4/tty.4.html
[zero.4]: http://man7.org/linux/man-pages/man4/zero.4.html
[user-namespaces]: http://man7.org/linux/man-pages/man7/user_namespaces.7.html
[intel-rdt-cat-kernel-interface]: https://www.kernel.org/doc/Documentation/x86/intel_rdt_ui.txt
