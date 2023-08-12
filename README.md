# Box

[![codecov](https://codecov.io/gh/jan-xyz/box/branch/main/graph/badge.svg?token=261hqevLHY)](https://codecov.io/gh/jan-xyz/box)
[![Go Reference](https://pkg.go.dev/badge/github.com/jan-xyz/box.svg)](https://pkg.go.dev/github.com/jan-xyz/box)

Boxing everything in neat little boxes, stack them, combine them, hide them.

The goal of this project is to provide simple abstractions that avoid leaking
models across boundaries and provide tools to create applications that are convenient
to reason about and maintainable with a strict compatibility with existing libraries.

The main idea stems from my experience with [go-kit](https://github.com/go-kit/kit)
and it inspired the majority of this project. However, over the years I noticed
some things where go-kit was deviating from standard tools, and with the release
generics in Go, I thought about providing a type-safe version of what go-kit taught
me.

## Philosophy

The basic idea is that most applications have to deal with different model layers
when they interface with the outside world:

1. transport protocol - e.g. HTTP requests/responses, async message buses,
  gRPC, lambda trigger etc.
1. data-transfer-object (DTO) - Protocolbuffers, JSON, AVRO, database models
1. internal model - the domain model used inside of the application

For each of the models there is a distinct layer that provides conversion logic,
domain specific handling and encapsulation.

Terms used in this applications to refer to the different layers are:

| Layer | Model | Concerns |
|-------|--------|---------|
| Transport | transport protocol | Communication patterns, status codes, meta-data extraction, error conversion |
| Endpoint | DTO | DTO decoding, validation & sanitisation, meta-data extraction |
| Service | internal model | businesss logic |

A model dependency must only point inwards, such that the Service has no dependency
on the DTO and the Endpoint has no dependency on the transport protocol. In the
same way a Layer must not jump a layer, such that the Transport knows nothing about
the internal model. In fact, the Transport can in most cases be a generic implementation
that doesn't need to know anything about the underlying models and all dependencies
can be injected via conversion functions. This allows sharing Transports, which
provides for easier implementation of infrastructure best-practices.

### Why should you separate layers?

#### Different Evolution Cadences

The layers have different evolution cadences and creating a lose coupling
allows for easier migration and maintenance of the individual parts.

The transport is usually very slow moving, especially established protocols like
HTTP. For the DTO it's not always up to you to decide on changes and coupling
it strongly with the internal model can turn into month-long refactoring work.

#### Migration to different models

By decoupling DTO from your internal model, you can connect your service
to different endpoints and transports. That is useful when the source of your information
changes, e.g. by consuming from a different message bus. From experience, this
happens a lot when micro-services migrate from one team to another or new major versions
are introduced.

#### Re using infrastructure

<!--this is still missing and needs more information -->

This repository also provides a collection of Transports that can be used off-the-shelf.

## Usage in services

See the [strsvc](./example/strsvc/) for a full example.

A basic example for encapsulating your service in an endpoint and transport

```go
import (
...
  "github.com/aws/aws-lambda-go/lambda"
  awslambdago "github.com/jan-xyz/box/transports/github.com/aws/aws-lambda-go"
...
)

func main() {
  db := database.New()
  s := service.New(db)
  ep := endpoint.New(s)
  h := awslambdago.NewAPIGatewayTransport(
    endpoint.Decode,
    endpoint.Encode,
    endpoint.EncodeError,
    ep,
  )

  lambda.StartHandlerFunc(h)
}
```

The above example is a lambda that is triggered by an API Gateway request. The layers
could look like this:

| Layer | Model |
|-------|-------|
| Transport | API Gateway Request & Response |
| Endpoint | JSON |
| Service | internal model |

Another example could be a lambda triggered by a DynamoDB Stream

 Layer | Model |
|-------|--------|
| Transport | DynamoDB Stream Record |
| Endpoint | database model |
| Service | internal model |

or an HTTP Server with a Protobuf body

| Layer | Model |
|-------|--------|
| Transport | HTTP Request & Response |
| Endpoint | protobuf model |
| Service | internal model |

## Usage for clients

For clients the same principle applies as above but from the rule
that dependencies only point inwards, we can conclude that instead of wrapping
the inner layers with the next outer layer, we provide clients that can deal
with the internal model. That leads to thin wrappers around provided libraries
like the aws-sdk or HTTP clients. The layers in this case look a bit different
and I inverted the order to better represent the flow of data through them.

| Layer | Model |
|-------|-------|
| Service | internal model | businesss logic |
| Endpoint | DTO | DTO encoding, meta-data injection |
| Transport | transport protocol | retries, backoffs, meta-data injection, error handling |

An example for an aws-sdk implementation could look like this

```go
import (
...
  "github.com/aws/aws-sdk-go-v2/service/dynamodb"
  "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
...
)
  type dynamoDBClient interface {
    GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
  }

  type db struct {
    client dynamoDBClient
  }

  func New() *db {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

    client := dynamodb.NewFromConfig(cfg)
    return &db{
      client: client, 
    }
  }

  func (d *db) GetCustomer(ctx context.Context, custmerID string) (service.Customer, error) {
    input := dynamodb.GetItemInput{
      TableName: "My-Table",
      Key: map[string]types.AttributeValue{
       "id": customerID,
      }
    }

    out, err := d.client.GetItem(ctx, input)
    if err != nil {
      return service.Customer{}, err
    }

    return dbModelToInternalModel(out)
  }
```

In most SDK examples you will not have to deal with the transport protocol and it
is hidden inside of the SDK. In other cases you will need to provide that yourself
and it makes sense to separate it in the same way. For example, the HTTP layer
of a service with retries, jitter, exponential back-offs should be separated
from the DTO and can be shared across many clients and even teams and companies.

<!-- I would like to provide an example for that as well. -->

## Further Reading and Inspiration

1. <https://jeffreypalermo.com/2008/07/the-onion-architecture-part-1/>
1. <https://alistair.cockburn.us/hexagonal-architecture/>
1. <https://gokit.io/faq/#architecture-and-design>
1. <https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html>
1. <https://en.wikipedia.org/wiki/Dependency_inversion_principle>
