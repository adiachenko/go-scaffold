package conf

var _ = Load()

// App
var (
	AppENV   = GetEnv("APP_ENV", "local")
	AppDebug = GetEnvAsBool("APP_DEBUG", false)
)

// RabbitMQ
var (
	RabbitMqUser     = GetEnv("RABBITMQ_USERNAME", "guest")
	RabbitMqPassword = GetEnv("RABBITMQ_PASSWORD", "guest")
	RabbitMqHost     = GetEnv("RABBITMQ_HOST", "localhost")
	RabbitMqPort     = GetEnvAsInt("RABBITMQ_PORT", 5672)
)

// Jaeger
var (
	TracingDriver      = GetEnv("TRACING_DRIVER", "noop")
	TracingServiceName = GetEnv("TRACING_SERVICE_NAME", "SQL Ingest Talent Insights")
	ZipkinHost         = GetEnv("ZIPKIN_HOST", "localhost")
	ZipkinPort         = GetEnv("ZIPKIN_PORT", "9411")
)
