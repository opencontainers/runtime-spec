# Runtime and Lifecycle

## State

The runtime state for a container is persisted on disk so that external tools can consume and act on this information.
The runtime state is stored in a JSON encoded file.
It is recommended that this file is stored in a temporary filesystem so that it can be removed on a system reboot.
On Linux based systems the state information should be stored in `/run/opencontainer/containers`.
The directory structure for a container is `/run/opencontainer/containers/<containerID>/state.json`.
By providing a default location that container state is stored external applications can find all containers running on a system.

* **`version`** (string) Version of the OCI specification used when creating the container.
* **`id`** (string) ID is the container's ID.
* **`pid`** (int) Pid is the ID of the main process within the container.
* **`bundlePath`** (string) BundlePath is the path to the container's bundle directory.

The ID is provided in the state because hooks will be executed with the state as the payload.
This allows the hook to perform clean and teardown logic after the runtime destroys its own state.

The root directory to the bundle is provided in the state so that consumers can find the container's configuration and rootfs where it is located on the host's filesystem.

*Example*

```json
{
    "id": "oc-container",
    "pid": 4422,
    "root": "/containers/redis"
}
```

## Typical lifecycle

A typical lifecyle progresses like this:

1. There is no container
2. A user tells the runtime to start a container and launch a process inside it
3. The runtime [creates the container](#create)
4. The runtime executes any [pre-start hooks](runtime-config.md#pre-start)
5. The runtime [executes the container process](#start-process)
6. The container process is running
7. A user tells the runtime to send a termination signal to the container process
8. The runtime [sends a termination signal to the container process](#stop-process)
9. The container process exits
10. The runtime [terminates any other processes in the container](#stop-process)
11. The runtime executes any [post-stop hooks](runtime-config.md#post-stop)
12. The runtime [removes the container](#cleanup)

With steps 7 and 8, the user is explicitly stopping the container process (via the runtime), but it's also possible that the container process could exit for other reasons.
In that case we skip directly from 6 to [10](#stop-process).

Failure in a pre-start hook or other setup task can cause a jump straight to [11](runtime-config.md#post-stop).

### Create

Create the container: file system, namespaces, cgroups, capabilities, etc.
The invoked process forks, with one branch that stays in the host namespace and another that enters the container.
The host process carries out all container setup actions, and continues running for the life of the container so it can perform teardown after the container process exits.
The container process changes users and drops privileges in preparation for the container process start.
At this point, the host process writes the [`state.json`](#state) file with the host-side version of the container-process's PID (the container process may be in a PID namespace).

### Start (process)

After the pre-start hooks complete, the host process signals the container process to execute the runtime.
The runtime execs the process defined in `config.json`'s [**`process`** attribute](config.md#process-configuration).
On Linux hosts, some information for this execution may come from outside the `config.json` and `runtime.json` specifications.
See the [Linux-specific notes for details](runtime-linux.md#file-descriptors).

### Stop (process)

Send a termination signal to the container process (can optionally send other signals to the container process, e.g. a kill signal).
When the process exits, the host process collects it's exit status to return as its own exit status.
If there are any remaining processes in the container's cgroup (and [we only support unified-hierarchies](runtime-config-linux.md#control-groups)), the host process kills and reaps them.

### Cleanup

The host process removes the [`state.json`](#state) file and the container: unmounting file systems, removing namespaces, etc.
This is the inverse of create.
The host process then exits with the container processes's exit status.

## Joining existing containers

Joining an existing container looks just like the usual workflow, except that the container process [joins the target container](runtime-config-linux.md#control-groups) at the beginning of step 3.
It can then, depending on its configuration, continue to create an additional child cgroup underneath the one it joined.

When exiting, the reaping logic in the [stop phase](#stop-process) is the same.
If the container process created a child cgroup, all other processes in that child cgroup are reaped, but no other processes in the joined cgroup (which the container process did not create) are reaped.
