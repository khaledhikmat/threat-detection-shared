package publisher

import (
	"context"

	"github.com/khaledhikmat/threat-detection-shared/models"
)

type IService interface {
	PublishRecordingClip(ctx context.Context, pubsub, topic string, clip models.RecordingClip) error
}
