package cli

import (
	"fmt"
	"io"
	"os/signal"
	"syscall"

	"github.com/doubledutch/quantum"
)

// Run runs a client
func Run(conn quantum.ClientConn, request quantum.Request, w io.Writer) error {
	// Used to ensure that all lines are captured
	doneCh := make(chan interface{}, 1)
	defer close(doneCh)

	go func() {
		// Print lines received on outCh
		for line := range conn.Logs() {
			fmt.Fprint(w, line)
		}
		doneCh <- struct{}{}
	}()

	signal.Notify(conn.Signals(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL)

	err := conn.Run(request)
	// Ensure that printing lines has finished
	<-doneCh

	return err
}
