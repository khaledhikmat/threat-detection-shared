package publisher

import (
	"context"
	"fmt"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/khaledhikmat/threat-detection-shared/equates"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

type publishFunction[T any] func(ctx context.Context, cfgsvc config.IService, client dapr.Client, pubsub, topic string, entity T) error

var providerFunctions = map[string]publishFunction[any]{
	"dapr": publishViaDapr[any],
}

func New(c dapr.Client, cfgsvc config.IService) IService {
	return &publisher{
		DaprClient: c,
		CfgSvc:     cfgsvc,
	}
}

type publisher struct {
	DaprClient dapr.Client
	CfgSvc     config.IService
}

func (s *publisher) PublishRecordingClip(ctx context.Context, pubsub, topic string, clip equates.RecordingClip) error {
	return publish[equates.RecordingClip](ctx, s.CfgSvc, s.DaprClient, pubsub, topic, clip)
}

func (s *publisher) Finalize() {
}

func publish[T any](ctx context.Context, cfgsvc config.IService, client dapr.Client, pubsub, topic string, entity T) error {
	fn, ok := providerFunctions[cfgsvc.GetPublisherProvider()]
	if !ok {
		return fmt.Errorf("provider %s not supported", cfgsvc.GetPublisherProvider())
	}

	return fn(ctx, cfgsvc, client, pubsub, topic, entity)
}
