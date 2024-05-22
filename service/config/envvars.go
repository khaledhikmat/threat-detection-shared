package config

import (
	"os"
	"strconv"
)

const (
	aiModelKey          = "AI_MODEL"
	alertTypeKey        = "ALERT_TYPE"
	mediaIndexerTypeKey = "MEDIA_INDEXER_TYPE"
	indexerTypeKey      = "INDEXER_TYPE"

	agentMode             = "AGENT_MODE"
	agentRecordingsFolder = "AGENT_RECORDINGS_FOLDER"
	agentSamplesFolder    = "AGENT_SAMPLES_FOLDER"
	capturerMaxCameras    = "CAPTURER_MAX_CAMERAS"
)

func newEnvVars() IService {
	p := envVarsService{}

	return &p
}

type envVarsService struct {
}

func (s *envVarsService) GetRuntime() string {
	return os.Getenv(runtimeModeKey)
}

func (s *envVarsService) GetSupportedAIModel() string {
	return os.Getenv(aiModelKey)
}

func (s *envVarsService) GetSupportedAlertType() string {
	return os.Getenv(alertTypeKey)
}

func (s *envVarsService) GetSupportedMediaIndexType() string {
	return os.Getenv(mediaIndexerTypeKey)
}

func (s *envVarsService) GetIndexerType() string {
	return os.Getenv(indexerTypeKey)
}

func (s *envVarsService) GetCapturer() Capturer {
	m, err := strconv.Atoi(os.Getenv(capturerMaxCameras))
	if err != nil {
		m = 3
	}

	return Capturer{
		AgentMode:        os.Getenv(agentMode),
		MaxCameras:       m,
		RecordingsFolder: os.Getenv(agentRecordingsFolder),
		SamplesFolder:    os.Getenv(agentSamplesFolder),
	}
}

func (s *envVarsService) Finalize() {
}
