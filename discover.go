//go:build !linux

package netward

// Discover returns active local TCP and UDP listening ports.
func Discover() ([]Listener, error) {
	return nil, ErrNotSupported
}
