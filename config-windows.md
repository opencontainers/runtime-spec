# <a name="windowsSpecificContainerConfiguration" />Windows-specific Container Configuration

This document describes the schema for the [Windows-specific section](config.md#platform-specific-configuration) of the [container configuration](config.md).
The Windows container specification uses APIs provided by the Windows Host Compute Service (HCS) to fulfill the spec.

## <a name="configWindowsResources" />Resources

**`resources`** (object, OPTIONAL, windows) configurs the container's resource limits.

The following subsections specify parameters for `resources`.

### <a name="configWindowsMemory" />Memory

**`memory`** *(object, OPTIONAL, windows)* configures the container's memory usage.

The following parameters can be specified:

* **`limit`** *(uint64, OPTIONAL, windows)* - sets limit of memory usage in bytes.

* **`reservation`** *(uint64, OPTIONAL, windows)* - sets the guaranteed minimum amount of memory for a container in bytes.

#### Example

```json
    "windows": {
        "resources": {
            "memory": {
                "limit": 2097152,
                "reservation": 524288
            }
        }
    }
```

### <a name="configWindowsCpu" />CPU

**`cpu`** *(object, OPTIONAL, windows)* configures for the container's CPU usage.

The following parameters can be specified:

* **`count`** *(uint64, OPTIONAL, windows)* - specifies the number of CPUs available to the container.

* **`shares`** *(uint16, OPTIONAL, windows)* - specifies the relative weight to other containers with CPU shares. The range is from 1 to 10000.

* **`percent`** *(uint, OPTIONAL, windows)* - specifies the percentage of available CPUs usable by the container.

#### Example

```json
    "windows": {
        "resources": {
            "cpu": {
                "percent": 50
            }
        }
    }
```

### <a name="configWindowsStorage" />Storage

**`storage`** *(object, OPTIONAL, windows)* configures the container's storage usage.

The following parameters can be specified:

* **`iops`** *(uint64, OPTIONAL, windows)* - specifies the maximum IO operations per second for the system drive of the container.

* **`bps`** *(uint64, OPTIONAL, windows)* - specifies the maximum bytes per second for the system drive of the container.

* **`sandboxSize`** *(uint64, OPTIONAL, windows)* - specifies the minimum size of the system drive in bytes.

#### Example

```json
    "windows": {
        "resources": {
            "storage": {
                "iops": 50
            }
        }
    }
```

### <a name="configWindowsNetwork" />Network

**`network`** *(object, OPTIONAL, windows)* configures the container's network usage.

The following parameters can be specified:

* **`egressBandwidth`** *(uint64, OPTIONAL, windows)* - specified the maximum egress bandwidth in bytes per second for the container.

#### Example

```json
    "windows": {
        "resources": {
            "network": {
                "egressBandwidth": 1048577
            }
        }
   }
```
