package publisher

import (
	"context"
	"fmt"

	dapr "github.com/dapr/go-sdk/client"

	"github.com/khaledhikmat/threat-detection-shared/equates"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

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
	if !s.CfgSvc.IsDapr() && !s.CfgSvc.IsDiagrid() {
		return nil
	}

	if (s.CfgSvc.IsDapr() || s.CfgSvc.IsDiagrid()) && s.DaprClient == nil {
		return fmt.Errorf("Dapr client is nil")
	}

	return s.DaprClient.PublishEvent(ctx, pubsub, topic, clip)
}

func (s *publisher) Finalize() {
}
