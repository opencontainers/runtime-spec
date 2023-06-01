# <a name="VirtualMachineSpecificContainerConfiguration" /> Virtual-machine-specific Container Configuration

This section describes the schema for the [virtual-machine-specific section](config.md#platform-specific-configuration) of the [container configuration](config.md).
The virtual-machine container specification provides additional configuration for the hypervisor, kernel, and image.

## <a name="HypervisorObject" /> Hypervisor Object

**`hypervisor`** (object, OPTIONAL) specifies details of the hypervisor that manages the container virtual machine.
* **`path`** (string, REQUIRED) path to the hypervisor binary that manages the container virtual machine.
    This value MUST be an absolute path in the [runtime mount namespace](glossary.md#runtime-namespace).
* **`parameters`** (array of strings, OPTIONAL) specifies an array of parameters to pass to the hypervisor.

### Example

```json
    "hypervisor": {
        "path": "/path/to/vmm",
        "parameters": ["opts1=foo", "opts2=bar"]
    }
```

## <a name="KernelObject" /> Kernel Object

**`kernel`** (object, REQUIRED) specifies details of the kernel to boot the container virtual machine with.
* **`path`** (string, REQUIRED) path to the kernel used to boot the container virtual machine.
    This value MUST be an absolute path in the [runtime mount namespace](glossary.md#runtime-namespace).
* **`parameters`** (array of strings, OPTIONAL) specifies an array of parameters to pass to the kernel.
* **`initrd`** (string, OPTIONAL) path to an initial ramdisk to be used by the container virtual machine.
    This value MUST be an absolute path in the [runtime mount namespace](glossary.md#runtime-namespace).

### Example

```json
    "kernel": {
        "path": "/path/to/vmlinuz",
        "parameters": ["foo=bar", "hello world"],
        "initrd": "/path/to/initrd.img"
    }
```

## <a name="ImageObject" /> Image Object

**`image`** (object, OPTIONAL) specifies details of the image that contains the root filesystem for the container virtual machine.
* **`path`** (string, REQUIRED) path to the container virtual machine root image.
    This value MUST be an absolute path in the [runtime mount namespace](glossary.md#runtime-namespace).
* **`format`** (string, REQUIRED) format of the container virtual machine root image. Commonly supported formats are:
    * **`raw`** [raw disk image format][raw-image-format]. Unset values for `format` will default to that format.
    * **`qcow2`** [QEMU image format][qcow2-image-format].
    * **`vdi`** [VirtualBox 1.1 compatible image format][vdi-image-format].
    * **`vmdk`** [VMware compatible image format][vmdk-image-format].
    * **`vhd`** [Virtual Hard Disk image format][vhd-image-format].

This image contains the root filesystem that the virtual machine **`kernel`** will boot into, not to be confused with the container root filesystem itself. The latter, as specified by **`path`** from the [Root Configuration](config.md#Root-Configuration) section, will be mounted inside the virtual machine at a location chosen by the virtual-machine-based runtime.

### Example

```json
    "image": {
        "path": "/path/to/vm/rootfs.img",
	"format": "raw"
    }
```

## <a name="HwConfigObject" /> HWConfig Object

**`hwConfig`** (object OPTIONAL) Specifies the hardware configuration that should be passed to the VM.
* **`deviceTree`** (string OPTIONAL) Path to the container device-tree file that should be passed to the VM.
* **`vcpus`** (int OPTIONAL) Number of virtual cpus for the VM.
* **`memory`** (int OPTIONAL) Maximum memory in bytes allocated to the VM.
* **`dtdevs`** (array OPTIONAL) Host device tree nodes to passthrough to the VM, see [Xen Config][xl-config-format] for the details.
* **`iomems`** (array OPTIONAL) Allow auto-translated domains to access specific hardware I/O memory pages, see [Xen Config][xl-config-format].
    * **`firstGFN`** (int OPTIONAL) Guest Frame Number to map the iomem range.
        If GFN is not specified, the mapping will be done to the same Frame Number as was provided in firstMFN, see [Xen Config][xl-config-format] for the details.
    * **`firstMFN`** (int REQUIRED) Physical page number of iomem regions, see [Xen Config][xl-config-format] for the details.
    * **`nrMFNs`** (int REQUIRED) Number of pages to be mapped, see [Xen Config][xl-config-format] for the details.
* **`irqs`** (array OPTIONAL) Allows VM to access specific physical IRQs, see [Xen Config][xl-config-format] for the details.

This hwConfig object contains the description of the hardware that can be safely passed through to the VM. Where **`deviceTree`** is the path to the device-tree blob, which conains description of the isolated hardware and paravirtualized hardware that should be used by VM. **`dtdevs`**, **`iomems`** and **`irqs`** parameters describing the minimun set of the parameters, needed for VM to access the hardware.

### Example

```json
    "hwConfig": {
        "deviceTree": "/path/to/vm/devicetree.dtb",
        "vcpus": 1,
        "memory": 4194304,
        "dtdevs": [
            "path/to/dev1_node",
            "path/to/dev2_node"
        ],
        "iomems": [
            {
                "firstMFN": 12288,
                "nrMFNs": 1
            },
            {
                "firstGFN": 12544,
                "firstMFN": 33024,
                "nrMFNs": 2
            }
        ],
        "irqs": [
            11,
            22
        ]
    }
```

[raw-image-format]: https://en.wikipedia.org/wiki/IMG_(file_format)
[qcow2-image-format]: https://git.qemu.org/?p=qemu.git;a=blob_plain;f=docs/interop/qcow2.txt;hb=HEAD
[vdi-image-format]: https://forensicswiki.org/wiki/Virtual_Disk_Image_(VDI)
[vmdk-image-format]: http://www.vmware.com/app/vmdk/?src=vmdk
[vhd-image-format]: https://github.com/libyal/libvhdi/blob/master/documentation/Virtual%20Hard%20Disk%20(VHD)%20image%20format.asciidoc
[xl-config-format]: https://xenbits.xen.org/docs/4.10-testing/man/xl.cfg.5.html