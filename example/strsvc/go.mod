module github.com/jan-xyz/box/example/strsvc

go 1.20

require (
	github.com/aws/aws-lambda-go v1.43.0
	github.com/jan-xyz/box v0.2.0
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
)

replace github.com/jan-xyz/box => ../../
