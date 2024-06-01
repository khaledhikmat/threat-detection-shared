package config

type IService interface {
	// Runtime environment
	GetRuntimeEnv() string
	GetRuntimeMode() string

	// Otel Provider type i.e. aws, noop, etc
	GetOtelProvider() string

	// Model Invoker type i.e. weapon, fire, etc
	GetSupportedAIModel() string
	// Alert Notifier type i.e. ccure, snow, pers, slack, email, etc
	GetSupportedAlertType() string
	// Alert Media Indexer type i.e. database, search service, etc
	GetSupportedMediaIndexType() string
	// Indexer type i.e. sqllite, AWS rds, etc
	GetIndexerType() string

	// Capturer configuration
	GetCapturer() Capturer

	Finalize()
}
