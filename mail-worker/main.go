package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/gemsnote/gemsnote/internal/mailqueue"
	"github.com/hibiken/asynq"
)

func main() {
	concurrency := flag.Int("concurrency", 2, "Maximum concurrent mail tasks")
	flag.Parse()
	if *concurrency < 1 || flag.NArg() != 0 {
		log.Fatal("usage: gemsnote-mail-worker [-concurrency N]")
	}
	client, err := mailqueue.NewRedis(os.Getenv("GEMSNOTE_MAIL_REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	server := asynq.NewServerFromRedisClient(client, asynq.Config{
		Concurrency: *concurrency, Queues: map[string]int{mailqueue.Queue: 1}, ShutdownTimeout: 35 * time.Second,
		RetryDelayFunc: func(n int, _ error, _ *asynq.Task) time.Duration { return time.Duration(n+1) * 5 * time.Second },
	})
	mux := asynq.NewServeMux()
	mux.Handle(mailqueue.TaskType, &mailqueue.Handler{Redis: client})
	if err := server.Run(mux); err != nil {
		log.Fatal("mail worker stopped")
	}
}
