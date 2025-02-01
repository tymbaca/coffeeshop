package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tymbaca/coffeeshop/barista/barista"
	"github.com/tymbaca/coffeeshop/barista/logger"
	"github.com/tymbaca/coffeeshop/barista/model"
	"github.com/tymbaca/coffeeshop/barista/tracer"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := tracer.Init("localhost:4318"); err != nil {
		logger.Fatal(err.Error())
	}

	b := barista.NewBarista(ctx, 4)

	// server := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: http.DefaultServeMux,
	// }
	//
	// http.HandleFunc("POST /order", b.HandleOrder)
	//
	// go func() {
	// 	<-ctx.Done()
	// 	server.Shutdown(context.Background())
	// }()
	//
	// logger.Info("starting the barista!")
	// if err := server.ListenAndServe(); err != nil {
	// 	if !errors.Is(err, http.ErrServerClosed) {
	// 		logger.Fatal(err.Error())
	// 	}
	// }

	natsConn, err := nats.Connect("nast://localhost:4222")
	if err != nil {
		panic(err)
	}

	sub, err := natsConn.Subscribe("order", handleNatsOrder(b))
	if err != nil {
		panic(err)
	}
	_ = sub

	<-ctx.Done()
	logger.Info("exiting...")
	time.Sleep(1 * time.Second)
}

func handleNatsOrder(b *barista.Barista) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		ctx := context.Background()

		var order model.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			err := msg.Ack()
			if err != nil {
				panic(err)
			}
		}

		traceMap := map[string]string{}
		for k := range msg.Header {
			traceMap[k] = msg.Header.Get(k)
		}
		ctx = tracer.FromMap(ctx, traceMap)

		b.HandleOrder(ctx, order)
	}
}
