package config

type IService interface {
	IsDapr() bool
	IsDiagrid() bool
	GetCapturer() Capturer
	GetPublisherProvider() string
	GetStorageProvider() string
	GetCloudStorage(provider string) CloudStorage
	Finalize()
}
