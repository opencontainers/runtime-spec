package specs_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/docker/go/canonical/json"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func TestConfigRoundTrip(t *testing.T) {
	for i, configString := range []string{
		// canonical version of the config.md example
		`{"ociVersion":"0.5.0-dev","platform":{"os":"linux","arch":"amd64"},"process":{"terminal":true,"consoleSize":{"height":0,"width":0},"user":{"uid":1,"gid":1,"additionalGids":[5,6]},"args":["sh"],"env":["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","TERM=xterm"],"cwd":"/","capabilities":{"bounding":["CAP_AUDIT_WRITE","CAP_KILL","CAP_NET_BIND_SERVICE"],"effective":["CAP_AUDIT_WRITE","CAP_KILL"],"inheritable":["CAP_AUDIT_WRITE","CAP_KILL","CAP_NET_BIND_SERVICE"],"permitted":["CAP_AUDIT_WRITE","CAP_KILL","CAP_NET_BIND_SERVICE"],"ambient":["CAP_NET_BIND_SERVICE"]},"rlimits":[{"type":"RLIMIT_CORE","hard":1024,"soft":1024},{"type":"RLIMIT_NOFILE","hard":1024,"soft":1024}],"noNewPrivileges":true,"apparmorProfile":"acme_secure_profile","selinuxLabel":"system_u:system_r:svirt_lxc_net_t:s0:c124,c675"},"root":{"path":"rootfs","readonly":true},"hostname":"slartibartfast","mounts":[{"destination":"/proc","type":"proc","source":"proc"},{"destination":"/dev","type":"tmpfs","source":"tmpfs","options":["nosuid","strictatime","mode=755","size=65536k"]},{"destination":"/dev/pts","type":"devpts","source":"devpts","options":["nosuid","noexec","newinstance","ptmxmode=0666","mode=0620","gid=5"]},{"destination":"/dev/shm","type":"tmpfs","source":"shm","options":["nosuid","noexec","nodev","mode=1777","size=65536k"]},{"destination":"/dev/mqueue","type":"mqueue","source":"mqueue","options":["nosuid","noexec","nodev"]},{"destination":"/sys","type":"sysfs","source":"sysfs","options":["nosuid","noexec","nodev"]},{"destination":"/sys/fs/cgroup","type":"cgroup","source":"cgroup","options":["nosuid","noexec","nodev","relatime","ro"]}],"hooks":{"prestart":[{"path":"/usr/bin/fix-mounts","args":["fix-mounts","arg1","arg2"],"env":["key1=value1"]},{"path":"/usr/bin/setup-network"}],"poststart":[{"path":"/usr/bin/notify-start","timeout":5}],"poststop":[{"path":"/usr/sbin/cleanup.sh","args":["cleanup.sh","-f"]}]},"annotations":{"com.example.key1":"value1","com.example.key2":"value2"},"linux":{"uidMappings":[{"hostID":1000,"containerID":0,"size":32000}],"gidMappings":[{"hostID":1000,"containerID":0,"size":32000}],"sysctl":{"net.core.somaxconn":"256","net.ipv4.ip_forward":"1"},"resources":{"devices":[{"allow":false,"access":"rwm"},{"allow":true,"type":"c","major":10,"minor":229,"access":"rw"},{"allow":true,"type":"b","major":8,"minor":0,"access":"r"}],"disableOOMKiller":false,"oomScoreAdj":100,"memory":{"limit":536870912,"reservation":536870912,"swap":536870912,"kernel":0,"kernelTCP":0,"swappiness":0},"cpu":{"shares":1024,"quota":1000000,"period":500000,"realtimeRuntime":950000,"realtimePeriod":1000000,"cpus":"2-3","mems":"0-7"},"pids":{"limit":32771},"blockIO":{"blkioWeight":10,"blkioLeafWeight":10,"blkioWeightDevice":[{"major":8,"minor":0,"weight":500,"leafWeight":300},{"major":8,"minor":16,"weight":500}],"blkioThrottleReadBpsDevice":[{"major":8,"minor":0,"rate":600}],"blkioThrottleWriteIOPSDevice":[{"major":8,"minor":16,"rate":300}]},"hugepageLimits":[{"pageSize":"2MB","limit":9223372036854772000}],"network":{"classID":1048577,"priorities":[{"name":"eth0","priority":500},{"name":"eth1","priority":1000}]}},"cgroupsPath":"/myRuntime/myContainer","namespaces":[{"type":"pid"},{"type":"network"},{"type":"ipc"},{"type":"uts"},{"type":"mount"},{"type":"user"},{"type":"cgroup"}],"devices":[{"path":"/dev/fuse","type":"c","major":10,"minor":229,"fileMode":438,"uid":0,"gid":0},{"path":"/dev/sda","type":"b","major":8,"minor":0,"fileMode":432,"uid":0,"gid":0}],"seccomp":{"defaultAction":"SCMP_ACT_ALLOW","architectures":["SCMP_ARCH_X86","SCMP_ARCH_X32"],"syscalls":[{"names":["getcwd","chmod"],"action":"SCMP_ACT_ERRNO","args":null}]},"rootfsPropagation":"slave","maskedPaths":["/proc/kcore","/proc/latency_stats","/proc/timer_stats","/proc/sched_debug"],"readonlyPaths":["/proc/asound","/proc/bus","/proc/fs","/proc/irq","/proc/sys","/proc/sysrq-trigger"],"mountLabel":"system_u:object_r:svirt_sandbox_file_t:s0:c715,c811"}}`,

		// minimal Linux example (removing optional fields from the config.md example)
		`{"ociVersion":"1.0.0","platform":{"os":"linux","arch":"amd64"},"process":{"user":{"uid":1,"gid":1},"args":["sh"],"cwd":"/"},"root":{"path":"rootfs"}}`,
	} {
		t.Run(fmt.Sprintf("config %d", i), func(t *testing.T) {
			configBytes := []byte(configString)
			var configStruct specs.Spec
			err := json.NewDecoder(bytes.NewReader(configBytes)).Decode(&configStruct)
			if err != nil {
				t.Fatalf("failed to decode: %v", err)
			}
			outBytes, err := json.Marshal(configStruct)
			if err != nil {
				t.Fatalf("failed to encode: %v", err)
			}
			if bytes.Compare(configBytes, outBytes) != 0 {
				t.Fatalf("failed to round-trip:\n%s", string(outBytes))
			}
		})
	}
}
