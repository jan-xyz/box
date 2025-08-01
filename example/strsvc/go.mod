module github.com/jan-xyz/box/example/strsvc

go 1.23.0

require (
	github.com/aws/aws-lambda-go v1.49.0
	github.com/jan-xyz/box v0.2.0
	go.opentelemetry.io/otel v1.37.0
	go.opentelemetry.io/otel/trace v1.37.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
)

replace github.com/jan-xyz/box => ../../
