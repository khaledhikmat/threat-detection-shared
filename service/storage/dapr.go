package storage

import (
	"context"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/khaledhikmat/threat-detection-shared/models"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

func NewDaprStorage(c dapr.Client, cfgsvc config.IService) IService {
	return &daprStorage{
		DaprClient: c,
		CfgSvc:     cfgsvc,
	}
}

type daprStorage struct {
	DaprClient dapr.Client
	CfgSvc     config.IService
}

func (s *daprStorage) StoreKeyValue(ctx context.Context, store, key, value string) error {
	return s.DaprClient.SaveState(ctx,
		store,
		key,
		[]byte(value),
		nil)
}

func (s *daprStorage) StoreRecordingClip(_ context.Context, _ models.RecordingClip) (string, error) {
	return "", nil
}

func (s *daprStorage) RetrieveRecordingClip(_ context.Context, _ models.RecordingClip) ([]byte, error) {
	return nil, nil
}

func (s *daprStorage) DownloadRecordingClip(_ context.Context, _ models.RecordingClip) ([]byte, error) {
	return nil, nil
}

func (s *daprStorage) Finalize() {
}
