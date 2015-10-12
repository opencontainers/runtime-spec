# Runtime and Lifecycle

## State

Runtime MUST store container metadata on disk so that external tools can consume and act on this information.
It is recommended that this data be stored in a temporary filesystem so that it can be removed on a system reboot.
On Linux/Unix based systems the metadata MUST be stored under `/run/opencontainer/containers`.
For non-Linux/Unix based systems the location of the root metadata directory is currently undefined.
Within that directory there MUST be one directory for each container created, where the name of the directory MUST be the ID of the container.
For example: for a Linux container with an ID of `173975398351`, there will be a corresponding directory: `/run/opencontainer/containers/173975398351`.
Within each container's directory, there MUST be a JSON encoded file called `state.json` that contains the runtime state of the container.
For example: `/run/opencontainer/containers/173975398351/state.json`.

The `state.json` file MUST contain all of the following properties:

* **`version`**: (string) is the OCF specification version used when creating the container.
* **`id`**: (string) is the container's ID.
This MUST be unique across all containers on this host.
There is no requirement that it be unique across hosts.
The ID is provided in the state because hooks will be executed with the state as the payload.
This allows the hooks to perform cleanup and teardown logic after the runtime destroys its own state.
* **`pid`**: (int) is the ID of the main process within the container, as seen by the host.
* **`bundlePath`**: (string) is the absolute path to the container's bundle directory.
This is provided so that consumers can find the container's configuration and root filesystem on the host.

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
7. The runtime executes any [post-start hooks](runtime-config.md#post-start)
8. A user tells the runtime to send a termination signal to the container process
9. The runtime [sends a termination signal to the container process](#stop-process)
10. The container process exits
11. The runtime [terminates any other processes in the container](#stop-process)
12. The runtime executes any [post-stop hooks](runtime-config.md#post-stop)
13. The runtime [removes the container](#cleanup)

With steps 7 and 8, the user is explicitly stopping the container process (via the runtime), but it's also possible that the container process could exit for other reasons.
In that case we skip directly from 6 to [10](#stop-process), skipping any post-start hooks that hadn't been launched and terminating any in-progress post-start hook.

Failure in a pre-start hook or other setup task can cause a jump straight to [12](runtime-config.md#post-stop).

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
