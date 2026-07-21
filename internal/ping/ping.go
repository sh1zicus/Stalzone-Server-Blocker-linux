package ping

import (
	"net"
	"time"
)

func TCP(address string, timeout time.Duration) (time.Duration, bool) {

	start := time.Now()

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return 0, false
	}

	_ = conn.Close()

	return time.Since(start), true
}
