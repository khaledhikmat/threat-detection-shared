package config

type IService interface {
	IsDapr() bool
	IsDiagrid() bool
	GetCapturer() Capturer
	GetStorageProvider() string
	GetCloudStorage(provider string) CloudStorage
	Finalize()
}
