package xerr

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Wrap(ctx context.Context, err error) error {
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}

func Errorf(ctx context.Context, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	return Wrap(ctx, err)
}
