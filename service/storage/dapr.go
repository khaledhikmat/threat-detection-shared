package storage

import (
	"context"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

func storeKeyValueViaDapr(ctx context.Context, cfgsvc config.IService, client dapr.Client, store, key, value string) error {
	return client.SaveState(ctx,
		store,
		key,
		[]byte(value),
		nil)
}
