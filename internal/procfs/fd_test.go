package procfs

import (
	"net"
	"os"
	"testing"
)

func TestResolvePIDs_EmptyTarget(t *testing.T) {
	res := ResolvePIDs(make(map[uint64]struct{}))
	if len(res) != 0 {
		t.Fatalf("expected empty map")
	}
}

func TestResolvePIDs_UnresolvedTarget(t *testing.T) {
	targets := map[uint64]struct{}{
		999999999999999999: {},
	}
	res := ResolvePIDs(targets)
	if len(res) != 0 {
		t.Fatalf("expected unresolved target to remain absent from result")
	}
	if len(targets) != 1 {
		t.Fatalf("expected unresolved target to remain in target set")
	}
}

func TestResolvePIDs_MatchingTarget(t *testing.T) {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	f, err := os.Open("/proc/net/tcp")
	if err != nil {
		t.Skip("cannot open /proc/net/tcp")
	}
	defer f.Close()

	entries, err := ParseListeners(f, "tcp4")
	if err != nil {
		t.Fatal(err)
	}

	targets := make(map[uint64]struct{})
	for _, e := range entries {
		targets[e.Inode] = struct{}{}
	}

	res := ResolvePIDs(targets)

	myPID := os.Getpid()
	found := false
	for _, pid := range res {
		if pid == myPID {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected at least one inode to resolve to my PID")
	}
}
