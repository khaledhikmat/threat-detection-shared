package storage

import (
	"context"
	"fmt"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/khaledhikmat/threat-detection-shared/models"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

var providerClipStoreFunctions = map[string]func(ctx context.Context, clip models.RecordingClip) (string, error){
	"dapr": storeClipViaAWS,
	"aws":  storeClipViaAWS,
}

var providerClipRetrieveFunctions = map[string]func(ctx context.Context, clip models.RecordingClip) ([]byte, error){
	"dapr": retrieveClipFromAWS,
	"aws":  retrieveClipFromAWS,
}

var providerClipDownloadFunctions = map[string]func(ctx context.Context, clip models.RecordingClip) ([]byte, error){
	"dapr": downloadClipFromAWS,
	"aws":  downloadClipFromAWS,
}

var providerKeyValueStoreFunctions = map[string]func(ctx context.Context, cfgsvc config.IService, client dapr.Client, store, key, value string) error{
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
	fn, ok := providerKeyValueStoreFunctions[s.CfgSvc.GetRuntime()]
	if !ok {
		return fmt.Errorf("provider %s not supported", s.CfgSvc.GetRuntime())
	}

	return fn(ctx, s.CfgSvc, s.DaprClient, store, key, value)
}

func (s *storage) StoreRecordingClip(ctx context.Context, clip models.RecordingClip) (string, error) {
	fn, ok := providerClipStoreFunctions[s.CfgSvc.GetRuntime()]
	if !ok {
		return "", fmt.Errorf("provider %s not supported", s.CfgSvc.GetRuntime())
	}

	return fn(ctx, clip)
}

func (s *storage) RetrieveRecordingClip(ctx context.Context, clip models.RecordingClip) ([]byte, error) {
	fn, ok := providerClipRetrieveFunctions[s.CfgSvc.GetRuntime()]
	if !ok {
		return []byte{}, fmt.Errorf("provider %s not supported", s.CfgSvc.GetRuntime())
	}

	return fn(ctx, clip)
}

func (s *storage) DownloadRecordingClip(ctx context.Context, clip models.RecordingClip) ([]byte, error) {
	fn, ok := providerClipDownloadFunctions[s.CfgSvc.GetRuntime()]
	if !ok {
		return []byte{}, fmt.Errorf("provider %s not supported", s.CfgSvc.GetRuntime())
	}

	return fn(ctx, clip)
}

func (s *storage) Finalize() {
}
