package config

type IService interface {
	IsDapr() bool
	IsDiagrid() bool
	GetCapturer() Capturer
	GetPublisherProvider() string
	GetKeyValStorageProvider() string
	GetFileStorageProvider() string
	GetCloudStorage(provider string) CloudStorage
	Finalize()
}
