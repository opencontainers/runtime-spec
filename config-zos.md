_This document is a work in progress._

# <a name="ZOSContainerConfiguration" />z/OS Container Configuration

This document describes the schema for the [z/OS-specific section](config.md#platform-specific-configuration) of the [container configuration](config.md).

## <a name="configZOSDevices" />Devices

**`devices`** (array of objects, OPTIONAL) lists devices that MUST be available in the container.
The runtime MAY supply them however it likes.

Each entry has the following structure:

* **`type`** *(string, REQUIRED)* - type of device: `c`, `b`, `u` or `p`.
* **`path`** *(string, REQUIRED)* - full path to device inside container.
    If a file already exists at `path` that does not match the requested device, the runtime MUST generate an error.
* **`major, minor`** *(int64, REQUIRED unless `type` is `p`)* - major, minor numbers for the device.
* **`fileMode`** *(uint32, OPTIONAL)* - file mode for the device.

The same `type`, `major` and `minor` SHOULD NOT be used for multiple devices.
