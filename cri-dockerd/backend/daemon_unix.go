//go:build !windows
// +build !windows

package backend

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/coreos/go-systemd/v22/daemon"
)

// listenFD returns a net.Listener from systemd socket activation.
// If no sockets are found, or the addr parameter is not supported, it returns an error.
func listenFD(addr string) (net.Listener, error) {
	var (
		err       error
		listeners []net.Listener
	)
	// socket activation
	listeners, err = activation.Listeners()
	if err != nil {
		return nil, err
	}

	if len(listeners) == 0 {
		return nil, errors.New(
			"no sockets found via socket activation: make sure the service was started by systemd",
		)
	}

	// default to first fd
	if addr == "" {
		return listeners[0], nil
	}

	// Address-specific listeners not supported yet
	return nil, errors.New("not supported yet")
}

// sdNotify sends a notification to systemd about service status.
func sdNotify(state string) error {
	_, err := daemon.SdNotify(false, state)
	if err != nil {
		// Logging of notification errors is intentionally skipped
		// because failures are non-critical and handled elsewhere.
		return err
	}
	return nil
}

// handleNotify sends systemd READY notification and handles SIGINT/SIGTERM signals.
func handleNotify() {
	// Notify systemd that service is ready
	sdNotify(daemon.SdNotifyReady)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		// Notify systemd that service is stopping
		sdNotify(daemon.SdNotifyStopping)
	}()
}
