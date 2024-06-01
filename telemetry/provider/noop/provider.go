package noopotelprovider

import (
	"context"
)

// Processor is a no-op OpenTelemetry provider.
func Processor(_ context.Context, _ string) (shutdown func(context.Context) error, err error) {

	err = nil
	shutdown = func(_ context.Context) error {
		return nil
	}

	return
}
