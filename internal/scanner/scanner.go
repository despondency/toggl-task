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

type Model struct {
	Currency        *string
	DeliveryDate    *string
	DueDate         *string
	InvoiceDate     *string
	ReceiptDate     *string
	NetAmount       *string
	TotalAmount     *string
	TotalTaxAmount  *string
	VatAmount       *string
	VatCategoryCode *string
	VatTaxAmount    *string
	VatTaxRate      *string
	Items           []*Item
}

type Item struct {
	Amount        *int
	Txt           *string
	Description   *string
	ProductCode   *string
	PurchaseOrder *string
	Quantity      *string
	Unit          *string
	UnitPrice     *string
}

func (m *Model) String() string {
	return fmt.Sprintln(m.TotalAmount)
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
		// this is hardcoded, in a real life project i'd inject from env variables or read from Vault.
		Name: fmt.Sprintf("projects/%s/locations/%s/processors/%s", "235872245316", "eu", "5cfedb6edd696fb"),
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
	model := createModel(document)

	_ = model

	entities, err := json.Marshal(document.GetEntities())
	if err != nil {
		panic(err)
	}
	return &ScannedResult{
		Result: string(entities),
	}, nil
}

func createModel(document *documentaipb.Document) *Model {
	model := &Model{
		Items: make([]*Item, 0),
	}
	var i *Item
	for _, entity := range document.GetEntities() {
		switch entity.GetType() {
		case "currency":
			model.Currency = strPtr(entity.GetMentionText())
		case "delivery_date":
			model.DeliveryDate = strPtr(entity.GetMentionText())
		case "due_date":
			model.DueDate = strPtr(entity.GetMentionText())
		case "invoice_date":
			model.InvoiceDate = strPtr(entity.GetMentionText())
		case "receipt_date":
			model.ReceiptDate = strPtr(entity.GetMentionText())
		case "net_amount":
			model.NetAmount = strPtr(entity.GetMentionText())
		case "total_amount":
			model.TotalAmount = strPtr(entity.GetMentionText())
		case "total_tax_amount":
			model.TotalTaxAmount = strPtr(entity.GetMentionText())
		case "vat_amount":
			model.VatAmount = strPtr(entity.GetMentionText())
		case "vat_category_code":
			vatCategoryCode := entity.GetMentionText()
			model.VatCategoryCode = &vatCategoryCode
		case "vat_tax_amount":
			model.VatTaxAmount = strPtr(entity.GetMentionText())
		case "vat_tax_rate":
			model.VatTaxRate = strPtr(entity.GetMentionText())
		case "line_item":
			i = &Item{}
			i.Txt = strPtr(entity.GetMentionText())
			model.Items = append(model.Items, i)
		case "line_item/description":
			i.Description = strPtr(entity.GetMentionText())
		case "line_item/product_code":
			i.ProductCode = strPtr(entity.GetMentionText())
		case "line_item/purchase_order":
			i.PurchaseOrder = strPtr(entity.GetMentionText())
		case "line_item/quantity":
			i.Quantity = strPtr(entity.GetMentionText())
		case "line_item/unit":
			i.Unit = strPtr(entity.GetMentionText())
		case "line_item/unit_price":
			i.UnitPrice = strPtr(entity.GetMentionText())
		}
	}
	return model
}

func strPtr(s string) *string {
	return &s
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
