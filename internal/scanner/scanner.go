package scanner

import (
	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/option"
)

type ScannedResult struct {
	Result string `json:"result"`
}

type Scanner interface {
	Scan(ctx context.Context, fileContent []byte, mimeType string) (*ScannedResult, error)
}

type GoogleScanner struct {
	client *documentai.DocumentProcessorClient
}

func (gs *GoogleScanner) Scan(ctx context.Context, fileContent []byte, mimeType string) (*ScannedResult, error) {
	req := &documentaipb.ProcessRequest{
		SkipHumanReview: true,
		Name:            fmt.Sprintf("projects/%s/locations/%s/processors/%s", "235872245316", "eu", "5cfedb6edd696fb"),
		Source: &documentaipb.ProcessRequest_RawDocument{
			RawDocument: &documentaipb.RawDocument{
				Content:  fileContent,
				MimeType: mimeType,
			},
		},
	}
	resp, err := gs.client.ProcessDocument(ctx, req)
	if err != nil {
		fmt.Println(fmt.Errorf("processDocument: %w", err))
	}
	// Handle the results.
	document := resp.GetDocument()
	entities, err := json.Marshal(document.GetEntities())
	if err != nil {
		panic(err)
	}
	return &ScannedResult{
		Result: string(entities),
	}, nil
}

func NewGoogleScanner(ctx context.Context) Scanner {
	endpoint := fmt.Sprintf("%s-documentai.googleapis.com:443", "eu")
	client, err := documentai.NewDocumentProcessorClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		fmt.Println(fmt.Errorf("error creating Document AI client: %w", err))
	}
	return &GoogleScanner{
		client: client,
	}
}
