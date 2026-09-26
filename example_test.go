package netward_test

import (
	"fmt"

	"github.com/41NI/NetWard"
)

func ExampleDiscover() {
	listeners, err := netward.Discover()
	if err != nil {
		if err == netward.ErrNotSupported {
			fmt.Println("NetWard is not supported on this OS")
			return
		}
		fmt.Printf("Discovery failed: %v\n", err)
		return
	}

	for _, l := range listeners {
		pidStr := "Unresolved"
		if l.ProcessID != netward.UnresolvedPID {
			pidStr = fmt.Sprintf("%d", l.ProcessID)
		}

		fmt.Printf("Discovered %s listener on %s:%d (PID: %s)\n",
			l.Protocol,
			l.LocalIP,
			l.LocalPort,
			pidStr,
		)
	}
}
