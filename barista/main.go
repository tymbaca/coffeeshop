package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/tymbaca/coffeeshop/barista/barista"
	"github.com/tymbaca/coffeeshop/barista/logger"
	"github.com/tymbaca/coffeeshop/barista/tracer"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := tracer.Init("localhost:4318"); err != nil {
		logger.Fatal(err.Error())
	}

	b := barista.NewBarista(ctx, 4)

	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	http.HandleFunc("POST /order", b.HandleOrder)

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	logger.Info("starting the barista!")
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
	}

	logger.Info("exiting...")
}
