package procfs

import (
	"bufio"
	"encoding/hex"
	"io"
	"net/netip"
	"strconv"
	"strings"
)

type SocketEntry struct {
	IP       netip.Addr
	Port     uint16
	Protocol string
	Inode    uint64
}

const (
	tcpStateListen      = "0A"
	udpStateUnconnected = "07"
)

func ParseListeners(r io.Reader, protocol string) ([]SocketEntry, error) {
	var entries []SocketEntry
	scanner := bufio.NewScanner(r)

	isTCP := strings.HasPrefix(protocol, "tcp")
	isUDP := strings.HasPrefix(protocol, "udp")

	if scanner.Scan() {
		_ = scanner.Text()
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		state := fields[3]
		if isTCP && state != tcpStateListen {
			continue
		}
		if isUDP && state != udpStateUnconnected {
			continue
		}

		parts := strings.Split(fields[1], ":")
		if len(parts) != 2 {
			continue
		}

		ipHex, portHex := parts[0], parts[1]

		port, err := strconv.ParseUint(portHex, 16, 16)
		if err != nil {
			continue
		}

		ipBytes, err := hex.DecodeString(ipHex)
		if err != nil {
			continue
		}

		var ip netip.Addr
		if len(ipBytes) == 4 {
			ip = netip.AddrFrom4([4]byte{ipBytes[3], ipBytes[2], ipBytes[1], ipBytes[0]})
		} else if len(ipBytes) == 16 {
			var v6 [16]byte
			for i := 0; i < 4; i++ {
				v6[i*4+0] = ipBytes[i*4+3]
				v6[i*4+1] = ipBytes[i*4+2]
				v6[i*4+2] = ipBytes[i*4+1]
				v6[i*4+3] = ipBytes[i*4+0]
			}
			ip = netip.AddrFrom16(v6)
		} else {
			continue
		}

		inodeStr := fields[9]
		inode, err := strconv.ParseUint(inodeStr, 10, 64)
		if err != nil {
			continue
		}

		entries = append(entries, SocketEntry{
			IP:       ip,
			Port:     uint16(port),
			Protocol: protocol,
			Inode:    inode,
		})
	}

	return entries, scanner.Err()
}
