//go:build linux

package netward

import (
	"os"

	"github.com/41NI/NetWard/internal/procfs"
)

// Discover returns active local TCP and UDP listening ports.
func Discover() ([]Listener, error) {
	files := []struct {
		path     string
		protocol string
	}{
		{"/proc/net/tcp", "tcp4"},
		{"/proc/net/tcp6", "tcp6"},
		{"/proc/net/udp", "udp4"},
		{"/proc/net/udp6", "udp6"},
	}

	var allEntries []procfs.SocketEntry

	for _, f := range files {
		file, err := os.Open(f.path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		entries, err := procfs.ParseListeners(file, f.protocol)
		file.Close()
		if err != nil {
			return nil, err
		}

		allEntries = append(allEntries, entries...)
	}

	targets := make(map[uint64]struct{}, len(allEntries))
	for _, entry := range allEntries {
		targets[entry.Inode] = struct{}{}
	}
	inodeToPID := procfs.ResolvePIDs(targets)

	var listeners []Listener
	for _, entry := range allEntries {
		pid := UnresolvedPID
		if resolved, ok := inodeToPID[entry.Inode]; ok {
			pid = resolved
		}

		listeners = append(listeners, Listener{
			LocalIP:   entry.IP,
			LocalPort: entry.Port,
			Protocol:  entry.Protocol,
			ProcessID: pid,
		})
	}

	return listeners, nil
}
