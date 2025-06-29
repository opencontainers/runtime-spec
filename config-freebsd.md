# <a name="FreeBSDContainerConfiguration" />FreeBSD Container Configuration

This document describes the schema for the [FreeBSD-specific section](config.md#platform-specific-configuration) of the [container configuration](config.md).

## <a name="configFreeBSDDevices" />Devices

Devices in FreeBSD are accessed via the `devfs` filesystem. Typically, each container will have a `devfs` filesystem mounted into its `/dev` directory. Often, a minimal set of devices is exposed to the container using ruleset 4 from `/etc/defaults/devfs.rules` - the ruleset is specified as a mount option.

Optionally, additional devices can be exposed to the container using an array of entries inside the `devices` root field:

* **`path`** _(string, REQUIRED)_ - the device path relative to `/dev`
* **`mode`** _(uint32, OPTIONAL)_ - file mode for the device.

### Example
```json
"devices": [
	{
        "path": "pf",
        "mode": 448
    }
]
```

## <a name="configFreeBSDJail" />Jail

On FreeBSD, containers are implemented using the platform's jail subsystem. Configuration for the container jail is via fields inside the `jail` root field:

* **`parent`** _(string, OPTIONAL)_ - parent jail.
    The value is the name of a jail which should be this container's parent (defaults to none). This can be used to share namespaces such as `vnet` with another container.
* **`host`** _(string, OPTIONAL)_ - allow overriding hostname, domainname, hostuuid and hostid.
    The value can be "new" which allows these values to be overridden in the container or "inherit" to use the host values (or parent container values). If set to "new", the values for hostname and domainname are taken from the base config, if present.
* **`ip4`** _(string, OPTIONAL)_ - control the availability of IPv4 addresses.
    This is typically left unset if the container has a vnet, set to "inherit" to allow access to host (or parent container) addresses or set to "disable" to stop use of IPv4 entirely.
* **`ip6`** _(string, OPTIONAL)_ - control the availability of IPv6 addresses.
    This is typically left unset if the container has a vnet, set to "inherit" to allow access to host (or parent container) addresses or set to "disable" to stop use of IPv6 entirely.
* **`vnet`** _(string, OPTIONAL)_ - control the vnet used for this container.
    The value can be "new" which causes a new vnet to be created for the container or "inherit" which shares the vnet for the parent container (or host if there is no parent).
* **`sysvmsg`** _(string, OPTIONAL)_ - allow access to SYSV IPC message primitives.
    If set to "inherit", all IPC objects in the host (or parent container) are visible to this container, whether they were created by the container itself, the base system, or other containers.  If set to "new", the container will have its own key namespace, and can only see the objects that it has created; the system (or parent container) has access to the container's objects, but not to its keys.  If set to "disable", the container cannot perform any sysvmsg-related system calls. Defaults to "new".
* **`sysvsem`** _(string, OPTIONAL)_ - allow access to SYSV IPC semaphore primitives, in the same manner as sysvmsg. Defaults to "new".
* **`sysvshm`** _(string, OPTIONAL)_ - allow access to SYSV IPC shared memory primitives, in the same manner as sysvmsg. Defaults to "new".
* **`enforceStatfs`** _(integer, OPTIONAL)_ - control visibility of mounts in the container.
    A value of 0 allows visibility of all host mounts, 1 allows visibility of mounts nested under the container's root and 2 only allows the container root to be visible. If unset, the default value is 2.
* **`allow`** _(object, OPTIONAL)_ - Some restrictions of the container environment may be set on a per-container basis.  With the exception of **`setHostname`** and **`reservedPorts`**, these boolean parameters are off by default.
  - **`setHostname`** _(bool, OPTIONAL)_ - Allow the container's hostname to be changed. Defaults to `false`.
  - **`rawSockets`** _(bool, OPTIONAL)_ - Allow the container to use raw sockets to support network utilities such as ping and traceroute. Defaults to `false`.
  - **`chflags`** _(bool, OPTIONAL)_ - Allow the system file flags to be changed. Defaults to `false`.
  - **`mount`** _(array of strings, OPTIONAL)_ - Allow the listed filesystem types to be mounted and unmounted in the container.
  - **`quotas`** _(bool, OPTIONAL)_ - Allow the filesystem quotas to be changed in the container. Defaults to `false`.
  - **`socketAf`** _(bool, OPTIONAL)_ - Allow socket types other than IPv4, IPv6 and unix. Defaults to `false`.
  - **`reservedPorts`** _(bool, OPTIONAL)_ - Allow the jail to bind to ports lower than 1024. Defaults to `false`.
  - **`suser`** _(bool, OPTIONAL)_ - The value of the jail's security.bsd.suser_enabled sysctl. The super-user will be disabled automatically if its parent system has it disabled.  The super-user is enabled by default.

### Mapping from jail(8) config file

This table defines the mappings from a typical `jail(8)` config file to the container configuration:

| Jail parameter   | JSON equivalent      |
| --------------   | -------------------- |
| `jid`            | -                    |
| `name`           | see below            |
| `path`           | `root.path`          |
| `ip4.addr`       | -                    |
| `ip4.saddrsel`   | -                    |
| `ip4`            | `freebsd.jail.ip4`   |
| `ip6.addr`       | -                    |
| `ip6.saddrsel`   | -                    |
| `ip6`            | `freebsd.jail.ip6`   |
| `vnet`           | `freebsd.jail.vnet`  |
| `host.hostname`  | `hostname`           |
| `host`           | `freebsd.jail.host`  |
| `sysvmsg`        | `freebsd.jail.sysvmsg` |
| `sysvsem`        | `freebsd.jail.sysvsem` |
| `sysvshm`        | `freebsd.jail.sysvshm` |
| `securelevel`    | -                    |
| `devfs_ruleset`  | see below            |
| `children.max`   | see below            |
| `enforce_statfs` | `freebsd.jail.enforceStatfs` |
| `persist`        | -                    |
| `parent`         | `freebsd.jail.parent`  |
| `osrelease`      | -                    |
| `osreldate`      | -                    |
| `allow.set_hostname` | `freebsd.jail.allow.setHostname` |
| `allow.sysvipc`  | `freebsd.jail.allow.sysvipc` |
| `allow.raw_sockets`  | `freebsd.jail.allow.rawSockets` |
| `allow.chflags`  | `freebsd.jail.allow.chflags` |
| `allow.mount`    | `freebsd.jail.allow.mount` |
| `allow.quotas`    | `freebsd.jail.allow.quotas` |
| `allow.read_msgbuf` | -                       |
| `allow.socket_af` | `freebsd.jail.allow.socketAf` |
| `allow.mlock`    | - |
| `allow.nfsd`     | - |
| `allow.reserved_ports` | `freebsd.jail.allow.reservedPorts` |
| `allow.unprivileged_proc_debug` | - |
| `allow.suser`    | `freebsd.jail.allow.suser` |
| `allow.mount.*`  | see below            |

The jail name SHOULD be set to the create command's `container-id` argument.

Network addresses are typically managed by the host (e.g. using CNI or netavark) so we do not include a mapping for `ip4.addr` or `ip6.addr`.

The `devfs_ruleset` parameter is only required for jails which create new `devfs` mounts - typically OCI runtimes will mount `devfs` on the host. The value is a rule set number - these rule sets are defined on the host, typically via `/etc/defaults/devfs.rules` or using the `devfs` command line utility.

The `children.max` parameter SHOULD be managed by the OCI runtime e.g. when a new container shares namespaces with an existing container.

The `allow.mount.*` parameter set is extensible - allowed mount types are listed as an array. As with `devfs`, typically the OCI runtime will manage mounts for the container by performing mount operations on the host.

Jail parameters not supported by this runtime extension are marked with "-". These parameters will have their default values - see the `jail(8)` man page for details.

### Example
```json
"jail": {
    "host": "new",
    "vnet": "new",
    "enforceStatfs": 1,
	"allow": {
		"rawSockets": true,
		"chflags": true,
		"mount": [
			"tmpfs"
		]
	}
}
```
