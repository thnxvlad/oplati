package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	hserver "github.com/thnxvlad/oplati/internal/server"
	"github.com/thnxvlad/oplati/internal/server/hmiddlewares"
)

func main() {
	publicServer := hserver.NewPublicServer(":8080", hmiddlewares.LoggingMiddleware)
	privateServer := hserver.NewPublicServer(":8081", hmiddlewares.LoggingMiddleware)

	go func() {
		fmt.Println("public server started...")
		err := publicServer.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		fmt.Println("private server started...")
		err := privateServer.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	<-ctx.Done()
	fmt.Println("system closed")
}
