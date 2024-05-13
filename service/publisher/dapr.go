package publisher

import (
	"context"
	"fmt"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

func publishViaDapr[T any](ctx context.Context, cfgsvc config.IService, client dapr.Client, pubsub, topic string, entity T) error {
	if !cfgsvc.IsDapr() && !cfgsvc.IsDiagrid() {
		return nil
	}

	if (cfgsvc.IsDapr() || cfgsvc.IsDiagrid()) && client == nil {
		return fmt.Errorf("Dapr client is nil")
	}

	return client.PublishEvent(ctx, pubsub, topic, entity)

}
