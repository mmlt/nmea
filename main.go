package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/mmlt/nmea/cmd/tool"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// trap SIGINT and cancel context
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	err := tool.NewRootCommand().ExecuteContext(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
