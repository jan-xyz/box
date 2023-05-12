module github.com/jan-xyz/box/example/strsvc

go 1.20

require (
	github.com/aws/aws-lambda-go v1.41.0
	github.com/jan-xyz/box v0.2.0
	go.opentelemetry.io/otel v1.15.1
	go.opentelemetry.io/otel/metric v0.38.1
	go.opentelemetry.io/otel/trace v1.15.1
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
)

replace github.com/jan-xyz/box => ../../
