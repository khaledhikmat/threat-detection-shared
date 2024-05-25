package storage

import (
	"context"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/khaledhikmat/threat-detection-shared/models"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

func NewStorage() IService {
	return &storage{}
}

type storage struct {
	DaprClient dapr.Client
	CfgSvc     config.IService
}

func (s *storage) StoreKeyValue(_ context.Context, _, _, _ string) error {
	return nil
}

func (s *storage) StoreRecordingClip(_ context.Context, _ models.RecordingClip) (string, error) {
	return "", nil
}

func (s *storage) RetrieveRecordingClip(_ context.Context, _ models.RecordingClip) ([]byte, error) {
	return nil, nil
}

func (s *storage) DownloadRecordingClip(_ context.Context, _ models.RecordingClip) ([]byte, error) {
	return nil, nil
}

func (s *storage) Finalize() {
}
