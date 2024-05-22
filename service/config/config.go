package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	RuntimeModeKey  = "THREAT_DETECTION_MODE"
	RuntimeModeDapr = "dapr"
	RuntimeModeAWS  = "aws"

	AIModelKey          = "AI_MODEL"
	AlertTypeKey        = "ALERT_TYPE"
	MediaIndexerTypeKey = "MEDIA_INDEXER_TYPE"
	IndexerTypeKey      = "INDEXER_TYPE"
)

type configService struct {
	Capturer Capturer  `json:"capturer"`
	FsData   *embed.FS `json:"-"`
}

func New(fs *embed.FS) IService {
	p := configService{
		FsData: fs,
	}

	err := json.Unmarshal(read(p.FsData, fmt.Sprintf("data/%s.json", "dev")), &p)
	if err != nil {
		panic(err)
	}

	return &p
}

func (s *configService) GetRuntime() string {
	return strings.ToLower(os.Getenv(RuntimeModeKey))
}

func (s *configService) GetSupportedAIModel() string {
	return os.Getenv(AIModelKey)
}

func (s *configService) GetSupportedAlertType() string {
	return os.Getenv(AlertTypeKey)
}

func (s *configService) GetSupportedMediaIndexType() string {
	return os.Getenv(MediaIndexerTypeKey)
}

func (s *configService) GetIndexerType() string {
	return os.Getenv(IndexerTypeKey)
}

func (s *configService) GetCapturer() Capturer {
	return s.Capturer
}

func (s *configService) Finalize() {
}

func read(fs *embed.FS, file string) []byte {
	fd, err := fs.ReadFile(file)
	if err != nil {
		fd, _ = fs.ReadFile("data/dev.json")
	}

	return fd
}
