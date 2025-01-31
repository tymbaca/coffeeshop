package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/tymbaca/coffeeshop/barista/logger"
	"github.com/tymbaca/coffeeshop/barista/model"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b := NewBarista(ctx, 4)

	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	http.HandleFunc("/order", b.order)

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

type barista struct {
	orderCh chan model.Order

	mu          sync.Mutex
	milk        int // ml
	coffeeBeans int // g
}

func NewBarista(ctx context.Context, workerCount int) *barista {
	b := &barista{
		orderCh:     make(chan model.Order, 1),
		milk:        10_000,
		coffeeBeans: 5_000,
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

func (b *barista) cook(ctx context.Context, order model.Order) error {
	if order.Type == "" {
		return fmt.Errorf("say something, please")
	}

	switch order.Type {
	case model.Cappuccino:
		if err := b.brewCappuccino(ctx); err != nil {
			return fmt.Errorf("can't brew cappuccino: %w", err)
		}
	case model.Latte:
		if err := b.brewLatte(ctx); err != nil {
			return fmt.Errorf("can't brew latte: %w", err)
		}
	case model.Espresso:
		if err := b.brewEspresso(ctx); err != nil {
			return fmt.Errorf("can't brew espresso: %w", err)
		}
	default:
		return fmt.Errorf("we don't brew %s", order.Type)
	}

	return nil
}

func (b *barista) brewCappuccino(ctx context.Context) error {
	if err := b.brewEspresso(ctx); err != nil {
		return fmt.Errorf("can't brew espresso: %w", err)
	}

	if err := b.pourMilk(ctx, 200); err != nil {
		return fmt.Errorf("can't pout milk: %w", err)
	}

	return nil
}

func (b *barista) brewLatte(ctx context.Context) error {
	if err := b.brewEspresso(ctx); err != nil {
		return fmt.Errorf("can't brew espresso: %w", err)
	}

	if err := b.pourMilk(ctx, 350); err != nil {
		return fmt.Errorf("can't pout milk: %w", err)
	}

	return nil
}

func (b *barista) brewEspresso(ctx context.Context) error {
	const amount = 8

	if err := b.getBeans(ctx, amount); err != nil {
		return fmt.Errorf("can't get beans: %w", err)
	}

	sleep(5000, 8000)

	return nil
}

func (b *barista) pourMilk(ctx context.Context, amount int) error {
	if err := b.getMilk(ctx, amount); err != nil {
		return fmt.Errorf("can't get milk: %w", err)
	}

	sleep(1000, 2000)
	return nil
}

func (b *barista) getBeans(_ context.Context, amount int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.coffeeBeans < amount {
		return fmt.Errorf("not enough coffee beans")
	}

	b.coffeeBeans -= amount

	return nil
}

func (b *barista) getMilk(_ context.Context, amount int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.milk < amount {
		return fmt.Errorf("not enough milk")
	}

	b.milk -= amount
	return nil
}

func sleep(min, max int) {
	time.Sleep(time.Duration(gofakeit.IntRange(min, max)) * time.Microsecond)
}
