package otelprovider

import (
	"context"

	aws "github.com/khaledhikmat/threat-detection-shared/telemetry/provider/aws"
	noop "github.com/khaledhikmat/threat-detection-shared/telemetry/provider/noop"
)

type OPType string

const (
	NoOp            OPType = "noop"
	AwsOtelProvider OPType = "aws"
)

type otelProviderConfig struct {
	Ctx            context.Context
	ServiceName    string
	ServiceVersion string
	Type           OPType
}

type OPOption func(f *otelProviderConfig)

func WithServiceName(serviceName string) OPOption {
	return func(f *otelProviderConfig) {
		f.ServiceName = serviceName
	}
}

func WithProviderType(opype OPType) OPOption {
	return func(f *otelProviderConfig) {
		f.Type = opype
	}
}

// Telemetry Provider Processor Signature
type otelProviderProcessor func(ctx context.Context, service string) (func(context.Context) error, error)

// Available Processors
var expProcs = map[OPType]otelProviderProcessor{
	NoOp:            noop.Processor,
	AwsOtelProvider: aws.Processor,
}

func New(ctx context.Context, service string, opts ...OPOption) (func(context.Context) error, error) {
	config := &otelProviderConfig{
		Ctx:         ctx,
		ServiceName: service,
		Type:        NoOp,
	}

	for _, opt := range opts {
		opt(config)
	}

	providerProc, ok := expProcs[config.Type]
	if !ok {
		providerProc = expProcs[NoOp]
	}

	// Run the provider to return shutdown function
	return providerProc(config.Ctx,
		config.ServiceName,
	)
}
