package storage

import (
	"context"
	"fmt"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/khaledhikmat/threat-detection-shared/equates"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

var providerFileFunctions = map[string]func(path2File string) (string, error){
	"aws":   storeFileViaAWS,
	"azure": storeFileViaAzure,
}

var providerKeyValueFunctions = map[string]func(ctx context.Context, cfgsvc config.IService, client dapr.Client, store, key, value string) error{
	"dapr": storeKeyValueViaDapr,
}

func New(c dapr.Client, cfgsvc config.IService) IService {
	return &storage{
		DaprClient: c,
		CfgSvc:     cfgsvc,
	}
}

type storage struct {
	DaprClient dapr.Client
	CfgSvc     config.IService
}

func (s *storage) StoreKeyValue(ctx context.Context, provider, store, key, value string) error {
	fn, ok := providerKeyValueFunctions[provider]
	if !ok {
		return fmt.Errorf("provider %s not supported", provider)
	}

	return fn(ctx, s.CfgSvc, s.DaprClient, store, key, value)
}

func (s *storage) StoreRecordingClip(_ context.Context, provider string, clip equates.RecordingClip) (string, error) {
	fn, ok := providerFileFunctions[provider]
	if !ok {
		return "", fmt.Errorf("provider %s not supported", provider)
	}

	return fn(clip.LocalReference)
}

func (s *storage) Finalize() {
}
