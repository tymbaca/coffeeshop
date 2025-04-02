package barista

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/tymbaca/coffeeshop/barista/logger"
	"github.com/tymbaca/coffeeshop/barista/model"
	"github.com/tymbaca/coffeeshop/barista/tracer"
	"github.com/tymbaca/coffeeshop/barista/xerr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Barista struct {
	orderCh chan model.Order

	mu          sync.Mutex
	milk        int // ml
	coffeeBeans int // g
}

func NewBarista(ctx context.Context, workerCount int) *Barista {
	b := &Barista{
		orderCh:     make(chan model.Order, 1),
		milk:        1_000,
		coffeeBeans: 200,
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

func (b *Barista) HTTPHandleOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	b.HandleOrder(r.Context(), order)
}

func (b *Barista) HandleOrder(ctx context.Context, order model.Order) {
	ctx, span := tracer.Start(ctx, "HandleOrder")
	defer span.End()

	order.Ctx = ctx

	b.orderCh <- order
}

func (b *Barista) runWorker(ctx context.Context) {
	for order := range b.orderCh {
		ctx = order.Ctx

		err := b.cook(ctx, order)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		logger.Infof("resources left: milk: %dml, coffee: %dg", b.milk, b.coffeeBeans)
	}
}

func (b *Barista) cook(ctx context.Context, order model.Order) (err error) {
	ctx, span := tracer.Start(ctx, "cook coffee", trace.WithAttributes(attribute.String("type", string(order.Type))))
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	if order.Type == "" {
		return xerr.Errorf(ctx, "say something, please")
	}

	orderID := span.SpanContext().TraceID().String()

	logger.Infof("got order %s: %v", orderID, order)
	start := time.Now()
	defer func() {
		logger.Infof("order %s completed, duration: %s", orderID, time.Since(start))
	}()

	switch order.Type {
	case model.Cappuccino:
		if err := b.brewCappuccino(ctx); err != nil {
			return xerr.Errorf(ctx, "can't brew cappuccino: %w", err)
		}
	case model.Latte:
		if err := b.brewLatte(ctx); err != nil {
			return xerr.Errorf(ctx, "can't brew latte: %w", err)
		}
	case model.Espresso:
		if err := b.brewEspresso(ctx); err != nil {
			return xerr.Errorf(ctx, "can't brew espresso: %w", err)
		}
	default:
		return xerr.Errorf(ctx, "we don't brew %s", order.Type)
	}

	return nil
}

func (b *Barista) brewCappuccino(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "brewCappuccino")
	defer span.End()

	if err := b.brewEspresso(ctx); err != nil {
		return xerr.Errorf(ctx, "can't brew espresso: %w", err)
	}

	if err := b.mixWithMilk(ctx, 200); err != nil {
		return xerr.Errorf(ctx, "can't pout milk: %w", err)
	}

	return nil
}

func (b *Barista) brewLatte(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "brewLatte")
	defer span.End()

	if err := b.brewEspresso(ctx); err != nil {
		return xerr.Errorf(ctx, "can't brew espresso: %w", err)
	}

	if err := b.mixWithMilk(ctx, 350); err != nil {
		return xerr.Errorf(ctx, "can't pout milk: %w", err)
	}

	return nil
}

func (b *Barista) brewEspresso(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "brewEspresso")
	defer span.End()

	const amount = 8

	if err := b.getPowder(ctx, amount); err != nil {
		return xerr.Errorf(ctx, "can't get beans: %w", err)
	}

	b.boilAndBrew(ctx)

	return nil
}

func (b *Barista) boilAndBrew(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "boilAndBrew")
	defer span.End()

	// do some work
	sleep(5000, 8000)
}

func (b *Barista) mixWithMilk(ctx context.Context, amount int) error {
	ctx, span := tracer.Start(ctx, "mixWithMilk")
	defer span.End()

	if err := b.getMilk(ctx, amount); err != nil {
		return xerr.Errorf(ctx, "can't get milk: %w", err)
	}

	b.pourMilk(ctx, amount)

	sleep(1000, 2000)
	return nil
}

func (b *Barista) pourMilk(ctx context.Context, _ int) {
	_, span := tracer.Start(ctx, "pourMilk")
	defer span.End()

	sleep(1000, 2000)
}

func (b *Barista) getPowder(ctx context.Context, amount int) error {
	ctx, span := tracer.Start(ctx, "getPowder")
	defer span.End()

	b.mu.Lock()
	defer b.mu.Unlock()

	sleep(1000, 2000)

	if b.coffeeBeans < amount {
		return xerr.Errorf(ctx, "not enough coffee beans")
	}

	b.coffeeBeans -= amount

	return nil
}

func (b *Barista) getMilk(ctx context.Context, amount int) error {
	ctx, span := tracer.Start(ctx, "getMilk")
	defer span.End()

	b.mu.Lock()
	defer b.mu.Unlock()

	sleep(1000, 2000)

	if b.milk < amount {
		return xerr.Errorf(ctx, "not enough milk")
	}

	b.milk -= amount
	return nil
}

const speed = 4

func sleep(min, max int) {
	dur := time.Duration(gofakeit.IntRange(min, max)) * time.Millisecond
	time.Sleep(dur / speed)
}
