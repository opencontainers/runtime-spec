# Linux-specific configuration

The Linux container specification uses various kernel features like namespaces,
cgroups, capabilities, LSM, and file system jails to fulfill the spec.
Additional information is needed for Linux over the [default spec configuration](config.md)
in order to configure these various kernel features.

## Linux namespaces

A namespace wraps a global system resource in an abstraction that makes it 
appear to the processes within the namespace that they have their own isolated 
instance of the global resource.  Changes to the global resource are visible to 
other processes that are members of the namespace, but are invisible to other 
processes. For more information, see [the man page](http://man7.org/linux/man-pages/man7/namespaces.7.html)

Namespaces are specified in the spec as an array of entries. Each entry has a 
type field with possible values described below and an optional path element. 
If a path is specified, that particular file is used to join that type of namespace.

```json
    "namespaces": [
        {
            "type": "pid",
            "path": "/proc/1234/ns/pid"
        },
        {
            "type": "net",
            "path": "/var/run/netns/neta"
        },
        {
            "type": "mnt",
        },
        {
            "type": "ipc",
        },
        {
            "type": "uts",
        },
        {
            "type": "user",
        },
    ]
```

#### Namespace types

* **pid** processes inside the container will only be able to see other processes inside the same container.
* **network** the container will have it's own network stack.
* **mnt** the container will have an isolated mount table.
* **ipc** processes inside the container will only be able to communicate to other processes inside the same
container via system level IPC.
* **uts** the container will be able to have it's own hostname and domain name.
* **user** the container will be able to remap user and group IDs from the host to local users and groups
within the container.

### Access to devices

Devices is an array specifying the list of devices to be created in the container.
Next parameters can be specified:

* type - type of device: 'c', 'b', 'u' or 'p'. More info in `man mknod`
* path - full path to device inside container
* major, minor - major, minor numbers for device. More info in `man mknod`.
                 There is special value: `-1`, which means `*` for `device`
                 cgroup setup.
* permissions - cgroup permissions for device. A composition of 'r'
                (read), 'w' (write), and 'm' (mknod).
* fileMode - file mode for device file
* uid - uid of device owner
* gid - gid of device owner

```json
   "devices": [
        {
            "type": "c",
            "path": "/dev/random",
            "major": 1,
            "minor": 8,
            "permissions": "rwm",
            "fileMode": 0666,
            "uid": 0,
            "gid": 0
        },
        {
            "type": "c",
            "path": "/dev/urandom",
            "major": 1,
            "minor": 9,
            "permissions": "rwm",
            "fileMode": 0666,
            "uid": 0,
            "gid": 0
        },
        {
            "type": "c",
            "path": "/dev/null",
            "major": 1,
            "minor": 3,
            "permissions": "rwm",
            "fileMode": 0666,
            "uid": 0,
            "gid": 0
        },
        {
            "type": "c",
            "path": "/dev/zero",
            "major": 1,
            "minor": 5,
            "permissions": "rwm",
            "fileMode": 0666,
            "uid": 0,
            "gid": 0
        },
        {
            "type": "c",
            "path": "/dev/tty",
            "major": 5,
            "minor": 0,
            "permissions": "rwm",
            "fileMode": 0666,
            "uid": 0,
            "gid": 0
        },
        {
            "type": "c",
            "path": "/dev/full",
            "major": 1,
            "minor": 7,
            "permissions": "rwm",
            "fileMode": 0666,
            "uid": 0,
            "gid": 0
        }
    ]
```

## Linux control groups

Also known as cgroups, they are used to restrict resource usage for a container and handle
device access.  cgroups provide controls to restrict cpu, memory, IO, and network for
the container. For more information, see the [kernel cgroups documentation](https://www.kernel.org/doc/Documentation/cgroups/cgroups.txt)

## Linux capabilities

Capabilities is an array that specifies Linux capabilities that can be provided to the process
inside the container. Valid values are the string after `CAP_` for capabilities defined
in [the man page](http://man7.org/linux/man-pages/man7/capabilities.7.html)

```json
   "capabilities": [
        "AUDIT_WRITE",
        "KILL",
        "NET_BIND_SERVICE"
    ]
```

## Linux sysctl

sysctl allows kernel parameters to be modified at runtime for the container.
For more information, see [the man page](http://man7.org/linux/man-pages/man8/sysctl.8.html)

```json
   "sysctl": {
        "net.ipv4.ip_forward": "1",
        "net.core.somaxconn": "256"
   }
```

## Linux rlimits

```json
   "rlimits": [
        {
            "type": "RLIMIT_NPROC",
            "soft": 1024,
            "hard": 102400
        }
   ]
```

rlimits allow setting resource limits. The type is from the values defined in [the man page](http://man7.org/linux/man-pages/man2/setrlimit.2.html). The kernel enforces the soft limit for a resource while the hard limit acts as a ceiling for that value that could be set by an unprivileged process.

## Linux user namespace mappings

```json
    "uidMappings": [
        {
            "hostID": 1000,
            "containerID": 0,
            "size": 10
        }
    ],
    "gidMappings": [
        {
            "hostID": 1000,
            "containerID": 0,
            "size": 10
        }
    ]
```

uid/gid mappings describe the user namespace mappings from the host to the container. *hostID* is the starting uid/gid on the host to be mapped to *containerID* which is the starting uid/gid in the container and *size* refers to the number of ids to be mapped. The Linux kernel has a limit of 5 such mappings that can be specified.

## Rootfs Mount Propagation
rootfsPropagation sets the rootfs's mount propagation. Its value is either slave, private, or shared. [The kernel doc](https://www.kernel.org/doc/Documentation/filesystems/sharedsubtree.txt) has more information about mount propagation.

```json
    "rootfsPropagation": "slave",
```

## Selinux process label

Selinux process label specifies the label with which the processes in a container are run.
For more information about SELinux, see  [Selinux documentation](http://selinuxproject.org/page/Main_Page)
```json
   "selinuxProcessLabel": "system_u:system_r:svirt_lxc_net_t:s0:c124,c675"
```

## Apparmor profile

Apparmor profile specifies the name of the apparmor profile that will be used for the container.
For more information about Apparmor, see [Apparmor documentation](https://wiki.ubuntu.com/AppArmor)

```json
   "apparmorProfile": "acme_secure_profile"
```

## Seccomp

Seccomp provides application sandboxing mechanism in the Linux kernel.
Seccomp configuration allows one to configure actions to take for matched syscalls and furthermore also allows
matching on values passed as arguments to syscalls.
For more information about Seccomp, see [Seccomp kernel documentation](https://www.kernel.org/doc/Documentation/prctl/seccomp_filter.txt)
The actions and operators are strings that match the definitions in seccomp.h from [libseccomp](https://github.com/seccomp/libseccomp) and are translated to corresponding values.

```json
   "seccomp": {
       "defaultAction": "SCMP_ACT_ALLOW",
       "syscalls": [
           {
               "name": "getcwd",
               "action": "SCMP_ACT_ERRNO"
           }
       ]
   }
```
