# Virtual-machine-specific Configuration

Virtual-machine-based runtimes require additional configuration to that specified in the [default spec configuration](config.md).

This optional configuration is specified in a "VM" object:

* **`imagePath`** (string, required) path to file that represents the root filesystem for the virtual machine.
* **`kernel`** (object, required) specifies details of the kernel to boot the virtual machine with.

Note that `imagePath` refers to a path on the host (outside of the virtual machine).
This field is distinct from the **`path`** field in the [Root Configuration](config.md#Root-Configuration) section since in the context of a virtual-machine-based runtime:

* **`imagePath`** will represent the root filesystem for the virtual machine.
* The container root filesystem specified by **`path`** from the [Root Configuration](config.md#Root-Configuration) section will be mounted inside the virtual machine at a location chosen by the virtual-machine-based runtime.

The virtual-machine-based runtime will use these two path fields to arrange for the **`path`** from the [Root Configuration](config.md#Root-Configuration) section to be presented to the process to run as the root filesystem.

## Kernel

Used by virtual-machine-based runtimes only.

* **`path`** (string, required) specifies the path to the kernel used to boot the virtual machine.
* **`parameters`** (string, optional) specifies a space-separated list of parameters to pass to the kernel.
* **`initrd`** (string, optional) specifies the path to an initial ramdisk to be used by the virtual machine.

## Example of a fully-populated `VM` object

```json
"vm": {
    "imagePath": "path/to/rootfs.img",
    "kernel": {
        "path": "path/to/vmlinuz",
        "parameters": "foo=bar hello world",
        "initrd": "path/to/initrd.img"
    },
}
```
