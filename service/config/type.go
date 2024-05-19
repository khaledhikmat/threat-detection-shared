package config

type IService interface {
	IsDapr() bool
	IsDiagrid() bool
	GetSupportedAIModel() string
	GetSupportedAlertType() string
	GetSupportedMediaIndexType() string
	GetCapturer() Capturer
	GetPublisherProvider() string
	GetKeyValStorageProvider() string
	GetFileStorageProvider() string
	GetIndexProvider() string
	GetCloudStorage(provider string) CloudStorage
	Finalize()
}
