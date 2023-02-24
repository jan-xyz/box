module github.com/jan-xyz/box/example/strsvc

go 1.19

require (
	github.com/aws/aws-lambda-go v1.37.0
	github.com/jan-xyz/box v0.2.0
	go.opentelemetry.io/otel v1.13.0
	go.opentelemetry.io/otel/trace v1.13.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
)

replace github.com/jan-xyz/box => ../../
