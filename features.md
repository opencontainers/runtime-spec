# <a name="features" />Features Structure

A [runtime](glossary.md#runtime) MAY provide a JSON structure about its implemented features to [runtime callers](glossary.md#runtime-caller).
This JSON structure is called ["Features structure"](glossary.md#features-structure).

The Features structure is irrelevant to the actual availability of the features in the host operating system.
Hence, the content of the Features structure SHOULD be determined on the compilation time of the runtime, not on the execution time.

All properties in the Features structure except `ociVersionMin` and `ociVersionMax` MAY either be absent or have the `null` value.
The `null` value MUST NOT be confused with an empty value such as `0`, `false`, `""`, `[]`, and `{}`.

## <a name="featuresSpecificationVersion" />Specification version

* **`ociVersionMin`** (string, REQUIRED) The minimum recognized version of the Open Container Initiative Runtime Specification.
  The runtime MUST accept this value as the [`ociVersion` property of `config.json`](config.md#specification-version).

* **`ociVersionMax`** (string, REQUIRED) The maximum recognized version of the Open Container Initiative Runtime Specification.
  The runtime MUST accept this value as the [`ociVersion` property of `config.json`](config.md#specification-version).
  The value MUST NOT be less than the value of the `ociVersionMin` property.
  The Features structure MUST NOT contain properties that are not defined in this version of the Open Container Initiative Runtime Specification.

### Example
```json
{
  "ociVersionMin": "1.0.0",
  "ociVersionMax": "1.1.0"
}
```

## <a name="featuresHooks" />Hooks
* **`hooks`** (array of strings, OPTIONAL) The recognized names of the [hooks](config.md#hooks).
  The runtime MUST support the elements in this array as the [`hooks` property of `config.json`](config.md#hooks).

### Example
```json
"hooks": [
  "prestart",
  "createRuntime",
  "createContainer",
  "startContainer",
  "poststart",
  "poststop"
]
```

## <a name="featuresMountOptions" />Mount Options

* **`mountOptions`** (array of strings, OPTIONAL) The recognized names of the mount options, including options that might not be supported by the host operating system.
  The runtime MUST recognize the elements in this array as the [`options` of `mounts` objects in `config.json`](config.md#mounts).
  * Linux: this array SHOULD NOT contain filesystem-specific mount options that are passed to the [mount(2)][mount.2] syscall as `const void *data`.

### Example

```json
"mountOptions": [
  "acl",
  "async",
  "atime",
  "bind",
  "defaults",
  "dev",
  "diratime",
  "dirsync",
  "exec",
  "iversion",
  "lazytime",
  "loud",
  "mand",
  "noacl",
  "noatime",
  "nodev",
  "nodiratime",
  "noexec",
  "noiversion",
  "nolazytime",
  "nomand",
  "norelatime",
  "nostrictatime",
  "nosuid",
  "nosymfollow",
  "private",
  "ratime",
  "rbind",
  "rdev",
  "rdiratime",
  "relatime",
  "remount",
  "rexec",
  "rnoatime",
  "rnodev",
  "rnodiratime",
  "rnoexec",
  "rnorelatime",
  "rnostrictatime",
  "rnosuid",
  "rnosymfollow",
  "ro",
  "rprivate",
  "rrelatime",
  "rro",
  "rrw",
  "rshared",
  "rslave",
  "rstrictatime",
  "rsuid",
  "rsymfollow",
  "runbindable",
  "rw",
  "shared",
  "silent",
  "slave",
  "strictatime",
  "suid",
  "symfollow",
  "sync",
  "tmpcopyup",
  "unbindable"
]
```


## <a name="featuresPlatformSpecificFeatures" />Platform-specific features

* **`linux`** (object, OPTIONAL) [Linux-specific features](features-linux.md).
  This MAY be set if the runtime supports `linux` platform.

## <a name="featuresAnnotations" />Annotations

**`annotations`** (object, OPTIONAL) contains arbitrary metadata of the runtime.
This information MAY be structured or unstructured.
Annotations MUST be a key-value map that follows the same convention as the Key and Values of the [`annotations` property of `config.json`](config.md#annotations).
However, annotations do not need to contain the possible values of the [`annotations` property of `config.json`](config.md#annotations).
The current version of the spec do not provide a way to enumerate the possible values of the [`annotations` property of `config.json`](config.md#annotations).

### Example
```json
"annotations": {
  "org.opencontainers.runc.checkpoint.enabled": "true",
  "org.opencontainers.runc.version": "1.1.0"
}
```

# Example

Here is a full example for reference.

```json
{
  "ociVersionMin": "1.0.0",
  "ociVersionMax": "1.1.0-rc.2",
  "hooks": [
    "prestart",
    "createRuntime",
    "createContainer",
    "startContainer",
    "poststart",
    "poststop"
  ],
  "mountOptions": [
    "async",
    "atime",
    "bind",
    "defaults",
    "dev",
    "diratime",
    "dirsync",
    "exec",
    "iversion",
    "lazytime",
    "loud",
    "mand",
    "noatime",
    "nodev",
    "nodiratime",
    "noexec",
    "noiversion",
    "nolazytime",
    "nomand",
    "norelatime",
    "nostrictatime",
    "nosuid",
    "nosymfollow",
    "private",
    "ratime",
    "rbind",
    "rdev",
    "rdiratime",
    "relatime",
    "remount",
    "rexec",
    "rnoatime",
    "rnodev",
    "rnodiratime",
    "rnoexec",
    "rnorelatime",
    "rnostrictatime",
    "rnosuid",
    "rnosymfollow",
    "ro",
    "rprivate",
    "rrelatime",
    "rro",
    "rrw",
    "rshared",
    "rslave",
    "rstrictatime",
    "rsuid",
    "rsymfollow",
    "runbindable",
    "rw",
    "shared",
    "silent",
    "slave",
    "strictatime",
    "suid",
    "symfollow",
    "sync",
    "tmpcopyup",
    "unbindable"
  ],
  "linux": {
    "namespaces": [
      "cgroup",
      "ipc",
      "mount",
      "network",
      "pid",
      "user",
      "uts"
    ],
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
    ],
    "cgroup": {
      "v1": true,
      "v2": true,
      "systemd": true,
      "systemdUser": true,
      "rdma": true
    },
    "seccomp": {
      "enabled": true,
      "actions": [
        "SCMP_ACT_ALLOW",
        "SCMP_ACT_ERRNO",
        "SCMP_ACT_KILL",
        "SCMP_ACT_KILL_PROCESS",
        "SCMP_ACT_KILL_THREAD",
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
        "SCMP_ARCH_RISCV64",
        "SCMP_ARCH_S390",
        "SCMP_ARCH_S390X",
        "SCMP_ARCH_X32",
        "SCMP_ARCH_X86",
        "SCMP_ARCH_X86_64"
      ],
      "knownFlags": [
        "SECCOMP_FILTER_FLAG_TSYNC",
        "SECCOMP_FILTER_FLAG_SPEC_ALLOW",
        "SECCOMP_FILTER_FLAG_LOG"
      ],
      "supportedFlags": [
        "SECCOMP_FILTER_FLAG_TSYNC",
        "SECCOMP_FILTER_FLAG_SPEC_ALLOW",
        "SECCOMP_FILTER_FLAG_LOG"
      ]
    },
    "apparmor": {
      "enabled": true
    },
    "selinux": {
      "enabled": true
    },
    "intelRdt": {
      "enabled": true
    }
  },
  "annotations": {
    "io.github.seccomp.libseccomp.version": "2.5.4",
    "org.opencontainers.runc.checkpoint.enabled": "true",
    "org.opencontainers.runc.commit": "v1.1.0-534-g26851168",
    "org.opencontainers.runc.version": "1.1.0+dev"
  }
}
```

[mount.2]: https://man7.org/linux/man-pages/man2/mount.2.html
