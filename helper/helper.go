package helper

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdown() {
	log.Println("Gracefully shutting down...")
	// Channel to listen for termination signals
	signals := make(chan os.Signal, 1)
	// Channel to signal that cleanup is complete
	done := make(chan bool, 1)

	// Notify signals channel on SIGINT or SIGTERM
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Block until a signal is received
		sig := <-signals
		fmt.Println("Signal received:", sig)

		// Perform cleanup or run other functions

		// Notify the main function that cleanup is complete
		done <- true
	}()

	// Block until the cleanup function has signaled completion
	<-done
	fmt.Println("Exiting program.")
}
