package config

import (
	"os"
)

const (
	runtimeEnvKey  = "RUN_TIME_ENV"
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

func (s *configService) GetRuntimeEnv() string {
	return os.Getenv(runtimeEnvKey)
}

func (s *configService) GetRuntimeMode() string {
	return os.Getenv(runtimeModeKey)
}

func (s *configService) GetOtelProvider() string {
	r, ok := providers[s.GetRuntimeMode()]
	if !ok {
		return ""
	}

	return r.GetOtelProvider()
}

func (s *configService) GetSupportedAIModel() string {
	r, ok := providers[s.GetRuntimeMode()]
	if !ok {
		return ""
	}

	return r.GetSupportedAIModel()
}

func (s *configService) GetSupportedAlertType() string {
	r, ok := providers[s.GetRuntimeMode()]
	if !ok {
		return ""
	}

	return r.GetSupportedAlertType()
}

func (s *configService) GetSupportedMediaIndexType() string {
	r, ok := providers[s.GetRuntimeMode()]
	if !ok {
		return ""
	}

	return r.GetSupportedMediaIndexType()
}

func (s *configService) GetIndexerType() string {
	r, ok := providers[s.GetRuntimeMode()]
	if !ok {
		return ""
	}

	return r.GetIndexerType()
}

func (s *configService) GetCapturer() Capturer {
	r, ok := providers[s.GetRuntimeMode()]
	if !ok {
		return Capturer{}
	}

	return r.GetCapturer()
}

func (s *configService) Finalize() {
	r, ok := providers[s.GetRuntimeMode()]
	if !ok {
		return
	}

	r.Finalize()
}
