package main

import (
	"context"
	"kafka-workers/internal/manager"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	input := make(chan int, 100)
	output := make(chan int, 100)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	go consumer(input, ctx)
	go consumer(input, ctx)

	go producer(output)

	workerManager := manager.NewManager(
		input,
		output,
	)

	go workerManager.ManagerWorker(ctx)

	<-ctx.Done()
	close(input)
	close(output)
	stop()
}

func consumer(
	input chan int,
	ctx context.Context,
) {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			value := rand.Intn(100)
			log.Printf("consumer: %d\n", value)
			input <- value
		}
	}
}

func producer(
	output chan int,
) {
	for item := range output {
		log.Println(item)
	}
}
