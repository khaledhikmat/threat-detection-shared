package storage

import (
	"context"

	"github.com/khaledhikmat/threat-detection-shared/equates"
)

type IService interface {
	StoreKeyValue(ctx context.Context, provider, store, key, value string) error
	StoreRecordingClip(ctx context.Context, provider string, clip equates.RecordingClip) (string, error)
	Finalize()
}
