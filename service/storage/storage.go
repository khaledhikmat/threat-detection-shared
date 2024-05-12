package storage

import (
	"context"
	"fmt"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/khaledhikmat/threat-detection-shared/equates"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

var providerFunctions = map[string]func(path2File string) (string, error){
	"aws":   store2aws,
	"azure": store2azure,
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
	if !s.CfgSvc.IsDapr() && !s.CfgSvc.IsDiagrid() {
		return nil
	}

	if (s.CfgSvc.IsDapr() || s.CfgSvc.IsDiagrid()) && s.DaprClient == nil {
		return fmt.Errorf("Dapr client is nil")
	}

	return s.DaprClient.SaveState(ctx,
		store,
		key,
		[]byte(value),
		nil)
}

func (s *storage) StoreRecordingClip(_ context.Context, provider string, clip equates.RecordingClip) (string, error) {
	fn, ok := providerFunctions[provider]
	if !ok {
		return "", fmt.Errorf("provider not found")
	}

	return fn(clip.LocalReference)
}

func (s *storage) Finalize() {
}

// Cloud Provider Functions
func store2aws(path2File string) (string, error) {
	return fmt.Sprintf("https://%s/aws_url", path2File), nil
}

func store2azure(path2File string) (string, error) {
	return fmt.Sprintf("https://%s/azure_url", path2File), nil
}
