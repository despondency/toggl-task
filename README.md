# Introduction

This service allows the upload of a receipt in png,jpeg,pdf.
Uses Google Document AI to scan it, stores the raw receipt on local disk and uploads scanned results to MongoDB

## Dependencies

Install the following

```
golang 1.21
docker
makefile  
```
### Mock generation dependencies
```
go install github.com/vektra/mockery/v2@v2.32.4
```

## Run Service

```shell
gcloud auth application-default login
make run
```

## Run unit tests

```
make test-unit
```

## Run integration tests

```
gcloud auth application-default login
make run
make test-integration
```

## Generate mocks 
```
mockery 
```

## Example curl requests to get started

Upload a receipt
```
curl --location 'http://localhost:8080/v1/receipt' \
--form 'file=@"<RECEIPTS FOLDER>/document2.pdf"'
```

Get a receipt result
```
curl --location 'http://localhost:8080/v1/receipt?id=<ID returned by the upload>'
```


## Improvements beyond version 1 (make it production ready)

Separate the persistence and scan layers.
Create persistence service, which will handle the upload solely (implement it to support any distributed storage (GCS, S3 etc)).
After the persistence service persists the receipt in its raw format (png, jpeg, pdf), it will send a message to a broker (NATS, Kafka) 
with payload of the form:
```
{
  "receipt": <path to the receipt in the distributed storage layer(S3, GCS, etc)> 
}
```
Then a scanner service will pick up the message and talk with Google Document AI. Scan the document, get whatever we need
from it and will persist the results to a NoSQL storage (I've chosen MongoDB, but we can go with BigQuery, Cassandra, DynamoDB)

Then we can have a receipt-api that sole purpose will be to just read and serve from our NoSQL storage

By splitting the whole task in three layers (receipt-persister, receipt-scanner, receipt-api), 
we allow ourselves to mature each one on their own without impacting all of the layers.

All of those can be deployed in Kubernetes as deployments. 
Of course all this requires to have distributed tracing, metrics and alerting. 

Also, general stuff around code sanitization, adding more configurability to application.go by introducing Config per environment

## Few words about the testing

Small sanity table unit test for the receipt_service that verify order of operations in this case since we don't have much business logic that is to be verified right now
I'd write httptest unit tests for the handlers even though there are integration tests. For the persisters and mongo as well. It's going to be mostly mocking and verification of ops.

I have not focused to go for 100% unit test coverage since I think the idea of this exercise is to show ideas and ways to do stuff rather than writing 20+ test cases that cover every possible err scenario.
I agree that in a real-life project I'd aim for 80%+. Teammates like to sleep without interruption, good for team morale.

Integration tests on the endpoints use the provided Google Document AI processor, some people may argue that external dependencies in integration test should be mocked
and I agree to some extent, because using the real Google Document AI processor in the test turns into e2e test, its a really thin line tho.

I'd experiment more with testcontainers for integration tests + getting the app up and running in docker rather than starting an in memory server
I've split the integration tests from the unit tests by using a ./tests folder for integration tests. Its stylistic opinion and not something to enforce, since in a 
real world scenario i'd have httptest validation unit tests on the handlers with the same name, and don't really wanna put both integration test and unit test in same file, it gets cluttered when there are a lot of tests