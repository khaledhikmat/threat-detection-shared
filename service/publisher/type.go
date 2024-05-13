package publisher

import (
	"context"

	"github.com/khaledhikmat/threat-detection-shared/equates"
)

type IService interface {
	PublishRecordingClip(ctx context.Context, pubsub, topic string, clip equates.RecordingClip) error
}
