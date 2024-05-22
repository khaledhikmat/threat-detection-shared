package config

import (
	"os"
)

const (
	runtimeModeKey = "RUN_TIME_MODE"
)

var providers map[string]IService

type configService struct {
}

func New() IService {
	providers = map[string]IService{
		"dapr": newEnvVars(),
		"aws":  newEnvVars(),
	}

	p := configService{}

	return &p
}

func (s *configService) GetRuntime() string {
	return os.Getenv(runtimeModeKey)
}

func (s *configService) GetSupportedAIModel() string {
	r, ok := providers[s.GetRuntime()]
	if !ok {
		return ""
	}

	return r.GetSupportedAIModel()
}

func (s *configService) GetSupportedAlertType() string {
	r, ok := providers[s.GetRuntime()]
	if !ok {
		return ""
	}

	return r.GetSupportedAlertType()
}

func (s *configService) GetSupportedMediaIndexType() string {
	r, ok := providers[s.GetRuntime()]
	if !ok {
		return ""
	}

	return r.GetSupportedMediaIndexType()
}

func (s *configService) GetIndexerType() string {
	r, ok := providers[s.GetRuntime()]
	if !ok {
		return ""
	}

	return r.GetIndexerType()
}

func (s *configService) GetCapturer() Capturer {
	r, ok := providers[s.GetRuntime()]
	if !ok {
		return Capturer{}
	}

	return r.GetCapturer()
}

func (s *configService) Finalize() {
	r, ok := providers[s.GetRuntime()]
	if !ok {
		return
	}

	r.Finalize()
}
