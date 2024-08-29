package main

import (
	"context"
	"fullcycle-client-server-api/client"
	"fullcycle-client-server-api/server"
	"log"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	srv := server.InitServer(ctx)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown server: %v", err)
		}

		log.Println("server shutdown")
	}()
	log.Println("Server initialized")

	client.GetDollarPrice(ctx)
	log.Println("Client done")

	log.Println("Waiting for server to close...")
	cancel()
	wg.Wait()

	log.Println("Server closed, exiting...")
}
