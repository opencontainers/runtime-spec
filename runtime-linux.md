# <a name="linuxRuntime" />Linux Runtime

## <a name="runtimeLinuxDevSymbolicLinks" /> Dev symbolic links

After the container has `/proc` mounted, the following standard symlinks MUST be setup within `/dev/` for the IO.

|    Source       | Destination |
| --------------- | ----------- |
| /proc/self/fd   | /dev/fd     |
| /proc/self/fd/0 | /dev/stdin  |
| /proc/self/fd/1 | /dev/stdout |
| /proc/self/fd/2 | /dev/stderr |
