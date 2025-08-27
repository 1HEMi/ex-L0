package main

import (
	"Ex-L0/internal/app"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()

		time.Sleep(200 * time.Millisecond)
	}()

	app.Run(ctx)
}
