//go:build linux

package netward

import (
	"net"
	"net/netip"
	"os"
	"strconv"
	"testing"
)

func TestDiscoverIntegration(t *testing.T) {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	expectedPort := uint16(addr.Port)
	expectedIP := netip.MustParseAddr("127.0.0.1")
	expectedPID := os.Getpid()

	listeners, err := Discover()
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}

	found := false
	for _, l := range listeners {
		if l.LocalIP == expectedIP && l.LocalPort == expectedPort {
			found = true
			if l.ProcessID != expectedPID {
				t.Errorf("expected PID %d for our test socket, got %d", expectedPID, l.ProcessID)
			}
			if l.Protocol != "tcp4" {
				t.Errorf("expected protocol tcp4, got %s", l.Protocol)
			}
			break
		}
	}

	if !found {
		t.Errorf("did not find our test socket (127.0.0.1:%d) in Discover output", expectedPort)
	}
}

func BenchmarkDiscover_EmptyTarget(b *testing.B) {
	b.ReportAllocs()

	var totalProcs, totalFDs int
	procDirs, _ := os.ReadDir("/proc")
	for _, p := range procDirs {
		if !p.IsDir() {
			continue
		}
		if _, err := strconv.Atoi(p.Name()); err == nil {
			totalProcs++
			fds, _ := os.ReadDir("/proc/" + p.Name() + "/fd")
			totalFDs += len(fds)
		}
	}

	listeners, _ := Discover()
	b.Logf("Environment: Procs=%d, TotalFDs=%d, Listeners=%d", totalProcs, totalFDs, len(listeners))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Discover()
	}
}

func BenchmarkDiscover_WithListener(b *testing.B) {
	b.ReportAllocs()

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		b.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	var totalProcs, totalFDs int
	procDirs, _ := os.ReadDir("/proc")
	for _, p := range procDirs {
		if !p.IsDir() {
			continue
		}
		if _, err := strconv.Atoi(p.Name()); err == nil {
			totalProcs++
			fds, _ := os.ReadDir("/proc/" + p.Name() + "/fd")
			totalFDs += len(fds)
		}
	}

	listeners, _ := Discover()
	b.Logf("Environment: Procs=%d, TotalFDs=%d, Listeners=%d", totalProcs, totalFDs, len(listeners))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Discover()
	}
}
