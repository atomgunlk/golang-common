package graceful

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// DefaultShutdownTimeout represents default shutdown timeout value
	DefaultShutdownTimeout = 10 * time.Second
)

// Shutdown represents callback function to do gracefully shutdown
type Shutdown func() error

// ListenSignal represents han
func ListenSignal(s Shutdown) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	fmt.Println("[ListenSignal]: listen for the signal of os.Interrupt and syscall.SIGTERM")

	<-quit

	return s()
}
