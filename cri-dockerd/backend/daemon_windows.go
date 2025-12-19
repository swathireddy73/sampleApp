//go:build windows
// +build windows

package backend

import (
	"errors"
	"net"
)

// listenFD returns an error on Windows because systemd socket activation is not supported.
func listenFD(addr string) (net.Listener, error) {
	// TODO addressed: Listening on a file descriptor is not supported on Windows.
	// This function intentionally returns an error to indicate that.
	return nil, errors.New("listening on a file descriptor is not supported on Windows")
}

// handleNotify is intentionally empty on Windows because systemd notifications do not exist.
func handleNotify() {
	// TODO addressed: On Windows, systemd notifications are not supported.
	// This function is intentionally left empty for cross-platform compatibility.
}
