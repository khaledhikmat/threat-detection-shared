Threat detection shared library lives in a separate repository.

```bash
go mod init github.com/khaledhikmat/threat-detection-shared
go get -u github.com/joho/godotenv
go get -u github.com/google/uuid

go get -u go.opentelemetry.io/contrib/propagators/aws/xray
go get -u go.opentelemetry.io/otel
go get -u go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
go get -u go.opentelemetry.io/otel/sdk/resource
go get -u go.opentelemetry.io/otel/sdk/trace
go get -u go.opentelemetry.io/otel/sdk/metric
go get -u go.opentelemetry.io/otel/sdk/metric/aggregation
go get -u go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc
```

More about telemetry....

