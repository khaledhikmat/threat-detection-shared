package pubsub

import (
	"context"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/khaledhikmat/threat-detection-shared/models"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

func NewDaprPubsub(c dapr.Client, cfgsvc config.IService) IService {
	return &daprPubsub{
		DaprClient: c,
		CfgSvc:     cfgsvc,
	}
}

type daprPubsub struct {
	DaprClient dapr.Client
	CfgSvc     config.IService
}

func (s *daprPubsub) CreateTopic(_ context.Context, topic string) (string, error) {
	return topic, nil
}

func (s *daprPubsub) CreateQueue(_ context.Context, queue, _ string) (string, string, error) {
	return queue, queue, nil
}

func (s *daprPubsub) PublishRecordingClip(ctx context.Context, pubsub, topic string, clip models.RecordingClip) error {
	return s.DaprClient.PublishEvent(ctx, pubsub, topic, clip)
}

func (s *daprPubsub) Subscribe(_ context.Context, _, _, _ string) (chan string, error) {
	return nil, nil
}

func (s *daprPubsub) Finalize() {
}
