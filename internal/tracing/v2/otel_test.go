// Copyright 2025 Redpanda Data, Inc.

package tracing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/redpanda-data/benthos/v4/public/service"
)

func TestInitSpansFromParentTextMap(t *testing.T) {
	t.Run("it will update the context for each message in the batch", func(t *testing.T) {
		textMap := map[string]any{
			"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		}

		msgOne := service.NewMessage([]byte("hello"))
		msgTwo := service.NewMessage([]byte("world"))

		batch := service.MessageBatch{msgOne, msgTwo}

		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))
		tp := noop.NewTracerProvider()

		err := InitSpansFromParentTextMap(tp, "test", textMap, batch)
		assert.NoError(t, err)

		spanOne := trace.SpanFromContext(batch[0].Context())
		assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", spanOne.SpanContext().TraceID().String())
		assert.Equal(t, "00f067aa0ba902b7", spanOne.SpanContext().SpanID().String())

		spanTwo := trace.SpanFromContext(batch[1].Context())
		assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", spanTwo.SpanContext().TraceID().String())
		assert.Equal(t, "00f067aa0ba902b7", spanTwo.SpanContext().SpanID().String())
	})
}
