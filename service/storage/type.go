package storage

import (
	"context"

	"github.com/khaledhikmat/threat-detection-shared/equates"
)

type IService interface {
	StoreKeyValue(ctx context.Context, store, key, value string) error
	StoreRecordingClip(ctx context.Context, clip equates.RecordingClip) (string, error)
	Finalize()
}
