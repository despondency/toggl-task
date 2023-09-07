package scanner

import (
	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"context"
	"fmt"
	"github.com/despondency/toggl-task/internal/model"
	"google.golang.org/api/option"
)

type Scanner interface {
	Scan(ctx context.Context, fileContent []byte, mimeType string) (*model.Receipt, error)
}

type GoogleScanner struct {
	client *documentai.DocumentProcessorClient
}

func (gs *GoogleScanner) Scan(ctx context.Context, fileContent []byte, mimeType string) (*model.Receipt, error) {
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
		return nil, err
	}
	// Handle the results.
	document := resp.GetDocument()
	m := createModel(document)
	return m, nil
}

func createModel(document *documentaipb.Document) *model.Receipt {
	m := &model.Receipt{
		Items: make([]*model.Item, 0),
	}
	var i *model.Item
	for _, entity := range document.GetEntities() {
		switch entity.GetType() {
		case "currency":
			m.Currency = strPtr(entity.GetMentionText())
		case "delivery_date":
			m.DeliveryDate = strPtr(entity.GetMentionText())
		case "due_date":
			m.DueDate = strPtr(entity.GetMentionText())
		case "invoice_date":
			m.InvoiceDate = strPtr(entity.GetMentionText())
		case "receipt_date":
			m.ReceiptDate = strPtr(entity.GetMentionText())
		case "net_amount":
			m.NetAmount = strPtr(entity.GetMentionText())
		case "total_amount":
			m.TotalAmount = strPtr(entity.GetMentionText())
		case "total_tax_amount":
			m.TotalTaxAmount = strPtr(entity.GetMentionText())
		case "vat_amount":
			m.VatAmount = strPtr(entity.GetMentionText())
		case "vat_category_code":
			vatCategoryCode := entity.GetMentionText()
			m.VatCategoryCode = &vatCategoryCode
		case "vat_tax_amount":
			m.VatTaxAmount = strPtr(entity.GetMentionText())
		case "vat_tax_rate":
			m.VatTaxRate = strPtr(entity.GetMentionText())
		case "line_item":
			i = &model.Item{}
			i.Txt = strPtr(entity.GetMentionText())
			m.Items = append(m.Items, i)
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
	return m
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
