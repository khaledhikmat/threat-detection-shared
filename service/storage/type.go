package storage

import (
	"context"

	"github.com/khaledhikmat/threat-detection-shared/models"
)

type IService interface {
	StoreKeyValue(ctx context.Context, store, key, value string) error
	StoreRecordingClip(ctx context.Context, clip models.RecordingClip) (string, error)
	RetrieveRecordingClip(ctx context.Context, clip models.RecordingClip) ([]byte, error)
	DownloadRecordingClip(ctx context.Context, clip models.RecordingClip) ([]byte, error)
	Finalize()
}
