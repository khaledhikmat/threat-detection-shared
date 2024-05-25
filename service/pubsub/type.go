package pubsub

import (
	"context"

	"github.com/khaledhikmat/threat-detection-shared/models"
)

type IService interface {
	CreateTopic(ctx context.Context, topic string) (string, error)
	CreateQueue(ctx context.Context, queue, topicARN string) (string, string, error)
	PublishRecordingClip(ctx context.Context, pubsub, topic string, clip models.RecordingClip) error
	Subscribe(ctx context.Context, topic, queueURL, queue string) (chan string, error)
}
