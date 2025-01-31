package main

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/tymbaca/coffeeshop/barista/logger"
	"github.com/tymbaca/coffeeshop/barista/model"
)

func main() {
	b := barista{}
	http.HandleFunc("/order", b.order)
}

type barista struct {
	orderCh chan model.Order

	mu          sync.Mutex
	milk        int // ml
	coffeeBeans int // g
}

func NewBarista(ctx context.Context, workerCount int) *barista {
	b := &barista{
		orderCh: make(chan model.Order, 1),
		milk:    10_000,
		coffee:  5_000,
	}

	for range workerCount {
		go b.runWorker(ctx)
	}

	go func() {
		<-ctx.Done()
		close(b.orderCh)
	}()

	return b
}

func (b *barista) order(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
}

func (b *barista) runWorker(ctx context.Context) {
	for order := range b.orderCh {
		err := b.cook(ctx, order)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}
}

func (b *barista) cook(_ context.Context, order model.Order) error {
	switch order.Type {
	case model.Cappuccino:

	case model.Latte:
	}

	time.Sleep(5*time.Second + time.Duration(rand.Intn(5000))*time.Second)

	return errors.New("not implemented")
}

func (b *barista) pourMilk(_ context.Context, amount int) error {
	b.mu.Lock()
	defer b.mu.Unlock()
}
