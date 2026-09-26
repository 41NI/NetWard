package procfs

import (
	"net/netip"
	"strings"
	"testing"
)

func TestParseListeners(t *testing.T) {
	// IPv4 TCP listening (0A); established (01) and malformed entries are ignored.
	tcp4Data := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode                                                     
   0: 00000000:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0
   1: 0100007F:1388 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 67890 1 0000000000000000 100 0 0 10 0
   2: 0100007F:1388 0100007F:9999 01 00000000:00000000 00:00000000 00000000  1000        0 67891 1 0000000000000000 100 0 0 10 0
   3: malformed line here
`

	entries, err := ParseListeners(strings.NewReader(tcp4Data), "tcp4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].IP != netip.MustParseAddr("0.0.0.0") || entries[0].Port != 80 || entries[0].Inode != 12345 {
		t.Errorf("unexpected entry 0: %+v", entries[0])
	}

	if entries[1].IP != netip.MustParseAddr("127.0.0.1") || entries[1].Port != 5000 || entries[1].Inode != 67890 {
		t.Errorf("unexpected entry 1: %+v", entries[1])
	}

	// IPv6 TCP listening (0A)
	tcp6Data := `  sl  local_address                         rem_address                           st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000000000000000000000000000:0050 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 11111 1 0000000000000000 100 0 0 10 0
`
	entries, err = ParseListeners(strings.NewReader(tcp6Data), "tcp6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].IP != netip.MustParseAddr("::") || entries[0].Port != 80 || entries[0].Inode != 11111 {
		t.Errorf("unexpected entry: %+v", entries[0])
	}

	// IPv4 UDP unconnected (07)
	udp4Data := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode                                                     
   0: 0100007F:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 22222 1 0000000000000000 100 0 0 10 0
`
	entries, err = ParseListeners(strings.NewReader(udp4Data), "udp4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].IP != netip.MustParseAddr("127.0.0.1") || entries[0].Port != 53 || entries[0].Inode != 22222 {
		t.Errorf("unexpected entry: %+v", entries[0])
	}

	// IPv6 UDP unconnected (07)
	udp6Data := `  sl  local_address                         rem_address                           st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000000000000000000001000000:0035 00000000000000000000000000000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 33333 1 0000000000000000 100 0 0 10 0
`
	entries, err = ParseListeners(strings.NewReader(udp6Data), "udp6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].IP != netip.MustParseAddr("::1") || entries[0].Port != 53 || entries[0].Inode != 33333 {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestParseListeners_UDPStateZeroA_NotTreatedAsListening(t *testing.T) {
	// Synthetic UDP entry with invalid listening state 0A.
	udpDataZeroA := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0035 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 22222 1 0000000000000000 100 0 0 10 0
`
	entries, err := ParseListeners(strings.NewReader(udpDataZeroA), "udp4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries because state 0A is not valid for UDP listening, got %d", len(entries))
	}
}
