package publisher

import (
	"context"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

func publishViaDapr[T any](ctx context.Context, cfgsvc config.IService, client dapr.Client, pubsub, topic string, entity T) error {
	return client.PublishEvent(ctx, pubsub, topic, entity)

}
