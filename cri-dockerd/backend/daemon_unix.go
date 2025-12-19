//go:build !windows
// +build !windows

package backend

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/coreos/go-systemd/v22/daemon"
)

// listenFD returns a net.Listener from systemd socket activation.
// If addr is empty, it defaults to the first available listener.
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

	return nil, errors.New("not supported yet")
}

// sdNotify sends a systemd notification with the given state.
// Logs errors instead of failing, because failing here can prevent proper shutdown.
func sdNotify(state string) error {
	_, err := daemon.SdNotify(false, state)
	if err != nil {
		// TODO addressed: Log instead of failing to maintain service stability
		log.Printf("sdNotify failed for state %s: %v", state, err)
		return err
	}
	return nil
}

// handleNotify sets up systemd ready/stopping notifications and signal handling.
func handleNotify() {
	// Notify systemd that the service is ready
	sdNotify(daemon.SdNotifyReady)

	// Set up channel to listen for termination signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		// Notify systemd that the service is stopping
		sdNotify(daemon.SdNotifyStopping)
	}()
}
