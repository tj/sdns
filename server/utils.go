package server

import "strings"

// Default port to the standard.
func defaultPort(addr string) string {
	if !strings.Contains(addr, ":") {
		return addr + ":53"
	}

	return addr
}
