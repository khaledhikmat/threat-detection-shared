package pubsub

import (
	"context"

	"github.com/khaledhikmat/threat-detection-shared/models"
)

func NewPubsub() IService {
	return &pubsub{}
}

type pubsub struct {
}

func (s *pubsub) CreateQueue(_ context.Context, queue, _ string) (string, string, error) {
	return queue, queue, nil
}

func (s *pubsub) CreateTopic(_ context.Context, topic string) (string, error) {
	return topic, nil
}

func (s *pubsub) PublishRecordingClip(_ context.Context, _, _ string, _ models.RecordingClip) error {
	return nil
}

func (s *pubsub) Subscribe(_ context.Context, _, _, _ string) (chan string, error) {
	return nil, nil
}

func (s *pubsub) Finalize() {
}
