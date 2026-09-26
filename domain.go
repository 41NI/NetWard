package netward

import (
	"errors"
	"net/netip"
)

// ErrNotSupported indicates the operating system is not supported.
var ErrNotSupported = errors.New("local port discovery is not supported on this operating system")

// UnresolvedPID indicates the socket's owning process ID could not be resolved.
const UnresolvedPID = 0

// Listener is an active local listening port.
type Listener struct {
	LocalIP netip.Addr

	LocalPort uint16

	Protocol string

	ProcessID int
}
