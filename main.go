package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/shipperizer/jobless-needle/threadpool"

	chi "github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type EnvSpec struct {
	Port int `envconfig:"port" default:"8000"`
}

func main() {

	specs := new(EnvSpec)

	if err := envconfig.Process("", specs); err != nil {
		panic(fmt.Errorf("issues with environment sourcing: %s", err))
	}

	var logger *zap.SugaredLogger
	if log, err := zap.NewDevelopment(); err != nil {
		logger = log.Sugar()
	}

	tp := threadpool.NewThreadPool()

	router := chi.NewMux()

	router.Get("/api/v0/tp/{limit:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		limit, _ := strconv.Atoi(chi.URLParam(r, "limit"))
		wg := sync.WaitGroup{}
		c := make(chan *threadpool.TaskResult[any], limit)

		wg.Add(limit)

		for i := 0; i < limit; i++ {
			_ = tp.Submit(fmt.Sprintf("Task-%d", i), func() any { return i }, c, &wg)
		}

		wg.Wait()

		var results []threadpool.TaskResult[int]

		results = threadpool.Take[int](c, limit)
		close(c)

		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(
			map[string]string{
				"count":  fmt.Sprintf("%v", len(results)),
				"data":   fmt.Sprintf("%+v", results),
				"status": fmt.Sprintf("%v", http.StatusOK),
			},
		)
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", specs.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	tp.Stop()

	logger.Desugar().Sync()

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logger.Info("Shutting down")
	os.Exit(0)
}
