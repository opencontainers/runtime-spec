package specs

import (
	"strconv"
	"testing"
)

func TestVersionSynced(t *testing.T) {
	wantVer := strconv.Itoa(VersionMajor) + "." + strconv.Itoa(VersionMinor) + "." + strconv.Itoa(VersionPatch) + VersionDev
	if Version != wantVer {
		t.Fatalf("Version mismatch: got %q, want %q", Version, wantVer)
	}
}
