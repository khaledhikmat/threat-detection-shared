package config

type Capturer struct {
	AgentMode          string `json:"agentMode"`
	MaxCameras         int    `json:"maxCameras"`
	RecordingsFolder   string `json:"recordingsFolder"`
	SamplesFolder      string `json:"samplesFolder"`
	StorageDestination string `json:"storageDestination"`
}

type CloudStorage struct {
	AccessKeyEnvVar string `json:"accessKeyEnvVar"`
	AccessKey       string `json:"accessKey"`
}
