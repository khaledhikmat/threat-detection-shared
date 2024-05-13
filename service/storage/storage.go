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

func (s *storage) StoreKeyValue(ctx context.Context, store, key, value string) error {
	fn, ok := providerKeyValueFunctions[s.CfgSvc.GetKeyValStorageProvider()]
	if !ok {
		return fmt.Errorf("provider %s not supported", s.CfgSvc.GetKeyValStorageProvider())
	}

	return fn(ctx, s.CfgSvc, s.DaprClient, store, key, value)
}

func (s *storage) StoreRecordingClip(_ context.Context, clip equates.RecordingClip) (string, error) {
	fn, ok := providerFileFunctions[s.CfgSvc.GetFileStorageProvider()]
	if !ok {
		return "", fmt.Errorf("provider %s not supported", s.CfgSvc.GetFileStorageProvider())
	}

	return fn(clip.LocalReference)
}

func (s *storage) Finalize() {
}
