# <a name="linuxFeatures" />Linux Features Structure

This document describes the [Linux-specific section](features.md#platform-specific-features) of the [Features structure](features.md).

## <a name="linuxFeaturesNamespaces" />Namespaces

* **`namespaces`** (array of strings, OPTIONAL) The recognized names of the namespaces, including namespaces that might not be supported by the host operating system.
  The runtime MUST recognize the elements in this array as the [`type` of `linux.namespaces` objects in `config.json`](config-linux.md#namespaces).

### Example

```json
"namespaces": [
  "cgroup",
  "ipc",
  "mount",
  "network",
  "pid",
  "user",
  "uts"
]
```

## <a name="linuxFeaturesCapabilities" />Capabilities

* **`capabilities`** (array of strings, OPTIONAL) The recognized names of the capabilities, including capabilities that might not be supported by the host operating system.
  The runtime MUST recognize the elements in this array in the [`process.capabilities` object of `config.json`](config.md#linux-process).

### Example

```json
"capabilities": [
  "CAP_CHOWN",
  "CAP_DAC_OVERRIDE",
  "CAP_DAC_READ_SEARCH",
  "CAP_FOWNER",
  "CAP_FSETID",
  "CAP_KILL",
  "CAP_SETGID",
  "CAP_SETUID",
  "CAP_SETPCAP",
  "CAP_LINUX_IMMUTABLE",
  "CAP_NET_BIND_SERVICE",
  "CAP_NET_BROADCAST",
  "CAP_NET_ADMIN",
  "CAP_NET_RAW",
  "CAP_IPC_LOCK",
  "CAP_IPC_OWNER",
  "CAP_SYS_MODULE",
  "CAP_SYS_RAWIO",
  "CAP_SYS_CHROOT",
  "CAP_SYS_PTRACE",
  "CAP_SYS_PACCT",
  "CAP_SYS_ADMIN",
  "CAP_SYS_BOOT",
  "CAP_SYS_NICE",
  "CAP_SYS_RESOURCE",
  "CAP_SYS_TIME",
  "CAP_SYS_TTY_CONFIG",
  "CAP_MKNOD",
  "CAP_LEASE",
  "CAP_AUDIT_WRITE",
  "CAP_AUDIT_CONTROL",
  "CAP_SETFCAP",
  "CAP_MAC_OVERRIDE",
  "CAP_MAC_ADMIN",
  "CAP_SYSLOG",
  "CAP_WAKE_ALARM",
  "CAP_BLOCK_SUSPEND",
  "CAP_AUDIT_READ",
  "CAP_PERFMON",
  "CAP_BPF",
  "CAP_CHECKPOINT_RESTORE"
]
```

## <a name="linuxFeaturesCgroup" />Cgroup

**`cgroup`** (object, OPTIONAL) represents the runtime's implementation status of cgroup managers.
Irrelevant to the cgroup version of the host operating system.

* **`v1`** (bool, OPTIONAL) represents whether the runtime supports cgroup v1.
* **`v2`** (bool, OPTIONAL) represents whether the runtime supports cgroup v2.
* **`systemd`** (bool, OPTIONAL) represents whether the runtime supports system-wide systemd cgroup manager.
* **`systemdUser`** (bool, OPTIONAL) represents whether the runtime supports user-scoped systemd cgroup manager.
* **`rdma`** (bool, OPTIONAL) represents whether the runtime supports RDMA cgroup controller.

### Example

```json
"cgroup": {
  "v1": true,
  "v2": true,
  "systemd": true,
  "systemdUser": true,
  "rdma": false
}
```

## <a name="linuxFeaturesSeccomp" />Seccomp

**`seccomp`** (object, OPTIONAL) represents the runtime's implementation status of seccomp.
Irrelevant to the kernel version of the host operating system.

* **`enabled`** (bool, OPTIONAL) represents whether the runtime supports seccomp.
* **`actions`** (array of strings, OPTIONAL) The recognized names of the seccomp actions.
  The runtime MUST recognize the elements in this array in the [`syscalls[].action` property of the `linux.seccomp` object in `config.json`](config-linux.md#seccomp).
* **`operators`** (array of strings, OPTIONAL) The recognized names of the seccomp operators.
  The runtime MUST recognize the elements in this array in the [`syscalls[].args[].op` property of the `linux.seccomp` object in `config.json`](config-linux.md#seccomp).
* **`archs`** (array of strings, OPTIONAL) The recognized names of the seccomp architectures.
  The runtime MUST recognize the elements in this array in the [`architectures` property of the `linux.seccomp` object in `config.json`](config-linux.md#seccomp).
* **`knownFlags`** (array of strings, OPTIONAL) The recognized names of the seccomp flags.
  The runtime MUST recognize the elements in this array in the [`flags` property of the `linux.seccomp` object in `config.json`](config-linux.md#seccomp).
* **`supportedFlags`** (array of strings, OPTIONAL) The recognized and supported names of the seccomp flags.
  This list may be a subset of `knownFlags` due to some flags not supported by the current kernel and/or libseccomp.
  The runtime MUST recognize and support the elements in this array in the [`flags` property of the `linux.seccomp` object in `config.json`](config-linux.md#seccomp).

### Example

```json
"seccomp": {
  "enabled": true,
  "actions": [
    "SCMP_ACT_ALLOW",
    "SCMP_ACT_ERRNO",
    "SCMP_ACT_KILL",
    "SCMP_ACT_LOG",
    "SCMP_ACT_NOTIFY",
    "SCMP_ACT_TRACE",
    "SCMP_ACT_TRAP"
  ],
  "operators": [
    "SCMP_CMP_EQ",
    "SCMP_CMP_GE",
    "SCMP_CMP_GT",
    "SCMP_CMP_LE",
    "SCMP_CMP_LT",
    "SCMP_CMP_MASKED_EQ",
    "SCMP_CMP_NE"
  ],
  "archs": [
    "SCMP_ARCH_AARCH64",
    "SCMP_ARCH_ARM",
    "SCMP_ARCH_MIPS",
    "SCMP_ARCH_MIPS64",
    "SCMP_ARCH_MIPS64N32",
    "SCMP_ARCH_MIPSEL",
    "SCMP_ARCH_MIPSEL64",
    "SCMP_ARCH_MIPSEL64N32",
    "SCMP_ARCH_PPC",
    "SCMP_ARCH_PPC64",
    "SCMP_ARCH_PPC64LE",
    "SCMP_ARCH_S390",
    "SCMP_ARCH_S390X",
    "SCMP_ARCH_X32",
    "SCMP_ARCH_X86",
    "SCMP_ARCH_X86_64"
  ],
  "knownFlags": [
    "SECCOMP_FILTER_FLAG_LOG"
  ],
  "supportedFlags": [
    "SECCOMP_FILTER_FLAG_LOG"
  ]
}
```

## <a name="linuxFeaturesApparmor" />AppArmor

**`apparmor`** (object, OPTIONAL) represents the runtime's implementation status of AppArmor.
Irrelevant to the availability of AppArmor on the host operating system.

* **`enabled`** (bool, OPTIONAL) represents whether the runtime supports AppArmor.

### Example

```json
"apparmor": {
  "enabled": true
}
```

## <a name="linuxFeaturesApparmor" />SELinux

**`selinux`** (object, OPTIONAL) represents the runtime's implementation status of SELinux.
Irrelevant to the availability of SELinux on the host operating system.

* **`enabled`** (bool, OPTIONAL) represents whether the runtime supports SELinux.

### Example

```json
"selinux": {
  "enabled": true
}
```

## <a name="linuxFeaturesIntelRdt" />Intel RDT

**`intelRdt`** (object, OPTIONAL) represents the runtime's implementation status of Intel RDT.
Irrelevant to the availability of Intel RDT on the host operating system.

* **`enabled`** (bool, OPTIONAL) represents whether the runtime supports Intel RDT.

### Example

```json
"intelRdt": {
  "enabled": true
}
```
