package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	RuntimeModeKey     = "THREAT_DETECTION_MODE"
	RuntimeModeDapr    = "dapr"
	RuntimeModeDiagrid = "diagrid"

	AIModelKey          = "AI_MODEL"
	AlertTypeKey        = "ALERT_TYPE"
	MediaIndexerTypeKey = "MEDIA_INDEXER_TYPE"
	SqlliteFilePath     = "SQLLITE_FILE_PATH"
)

type configService struct {
	Dapr                  bool                    `json:"isDapr"`
	Diagrid               bool                    `json:"isDiagrid"`
	Capturer              Capturer                `json:"capturer"`
	PublisherProvider     string                  `json:"publisherProvider"`
	KeyValStorageProvider string                  `json:"keyValStorageProvider"`
	FileStorageProvider   string                  `json:"fileStorageProvider"`
	IndexProvider         string                  `json:"indexProvider"`
	CloudStorageProviders map[string]CloudStorage `json:"cloudStorageProviders"`
	FsData                *embed.FS               `json:"-"`
}

func New(fs *embed.FS) IService {
	p := configService{
		FsData: fs,
	}

	err := json.Unmarshal(read(p.FsData, fmt.Sprintf("data/%s.json", "dev")), &p)
	if err != nil {
		panic(err)
	}

	// Resolve storage provider access keys
	for key, value := range p.CloudStorageProviders {
		if value.AccessKeyEnvVar != "" {
			value.AccessKey = os.Getenv(value.AccessKeyEnvVar)
			p.CloudStorageProviders[key] = value
		}
	}

	return &p
}

func (s *configService) IsDapr() bool {
	mode := strings.ToLower(os.Getenv(RuntimeModeKey))
	if mode == RuntimeModeDapr || mode == RuntimeModeDiagrid {
		return true
	}

	return false
}

func (s *configService) IsDiagrid() bool {
	return strings.ToLower(os.Getenv(RuntimeModeKey)) == RuntimeModeDiagrid
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

func (s *configService) GetCapturer() Capturer {
	return s.Capturer
}

func (s *configService) GetPublisherProvider() string {
	return s.PublisherProvider
}

func (s *configService) GetKeyValStorageProvider() string {
	return s.KeyValStorageProvider
}

func (s *configService) GetFileStorageProvider() string {
	return s.FileStorageProvider
}

func (s *configService) GetIndexProvider() string {
	return s.IndexProvider
}

func (s *configService) GetCloudStorage(provider string) CloudStorage {
	value, ok := s.CloudStorageProviders[provider]
	if !ok {
		return CloudStorage{}
	}

	return value
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
