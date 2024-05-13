package storage

import (
	"context"
	"fmt"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

func storeKeyValueViaDapr(ctx context.Context, cfgsvc config.IService, client dapr.Client, store, key, value string) error {
	if !cfgsvc.IsDapr() && !cfgsvc.IsDiagrid() {
		return nil
	}

	if (cfgsvc.IsDapr() || cfgsvc.IsDiagrid()) && client == nil {
		return fmt.Errorf("Dapr client is nil")
	}

	return client.SaveState(ctx,
		store,
		key,
		[]byte(value),
		nil)
}
