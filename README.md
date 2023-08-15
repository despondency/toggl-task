# Introduction

This service is allowing the upload of a receipt in png,jpeg,pdf.
Uses Google Document AI to scan it, stores the raw receipt on local disk and uploads scanned results to MongoDB

***I am not including the receipts you've sent me here in the project root, 
since they may or may not contain sensitive information.
Extract them somewhere on your machine and use the curl commands below to do a request***

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
make run
```

## Run unit tests

```
make unit-tests
```

## Run integration tests

```
make integration-tests
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


