package scanner

import (
	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/api/option"
	"io/ioutil"
)

type ScannedResult struct {
	Result string `json:"result"`
}

type Scanner interface {
	Scan(fileContent []byte) (string, error)
}

type GoogleScanner struct {
}

func NewGoogleScanner() Scanner {
	projectID := "235872245316"
	location := "eu"
	// Create a Processor before running sample
	processorID := "5cfedb6edd696fb"
	filePath := "?"
	mimeType := "image/png"
	flag.Parse()

	ctx := context.Background()

	endpoint := fmt.Sprintf("%s-documentai.googleapis.com:443", location)
	client, err := documentai.NewDocumentProcessorClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		fmt.Println(fmt.Errorf("error creating Document AI client: %w", err))
	}
	defer client.Close()

	// Open local file.
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(fmt.Errorf("ioutil.ReadFile: %w", err))
	}

	req := &documentaipb.ProcessRequest{
		SkipHumanReview: true,
		Name:            fmt.Sprintf("projects/%s/locations/%s/processors/%s", projectID, location, processorID),
		Source: &documentaipb.ProcessRequest_RawDocument{
			RawDocument: &documentaipb.RawDocument{
				Content:  data,
				MimeType: mimeType,
			},
		},
	}
	resp, err := client.ProcessDocument(ctx, req)
	if err != nil {
		fmt.Println(fmt.Errorf("processDocument: %w", err))
	}

	// Handle the results.
	document := resp.GetDocument()
	entities, err := json.Marshal(document.GetEntities())
	if err != nil {
		panic(err)
	}
	fmt.Println(string(entities))
	_ = ScannedResult{
		Result: string(entities),
	}
	return nil
}
