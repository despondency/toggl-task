package scanner

import (
	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/api/option"
	"strconv"
)

type ScannedResult struct {
	Result string `json:"result"`
}

type Model struct {
	currency        *string
	deliveryDate    *string
	dueDate         *string
	invoiceDate     *string
	receiptDate     *string
	netAmount       *float64
	totalAmount     *float64
	totalTaxAmount  *float64
	vatAmount       *float64
	vatCategoryCode *string
	vatTaxAmount    *float64
	vatTaxRate      *float64
	items           []*Item
}

type Item struct {
	amount        *int
	txt           *string
	description   *string
	productCode   *string
	purchaseOrder *string
	quantity      *int64
	unit          *string
	unitPrice     *float64
}

func (m *Model) String() string {
	return fmt.Sprintln(m.totalAmount)
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

	log.Infof("%s ", model)

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
		items: make([]*Item, 0),
	}
	var i *Item
	for _, entity := range document.GetEntities() {
		switch entity.GetType() {
		case "currency":
			model.currency = strPtr(entity.GetMentionText())
		case "delivery_date":
			model.deliveryDate = strPtr(entity.GetMentionText())
		case "due_date":
			model.dueDate = strPtr(entity.GetMentionText())
		case "invoice_date":
			model.invoiceDate = strPtr(entity.GetMentionText())
		case "receipt_date":
			model.receiptDate = strPtr(entity.GetMentionText())
		case "net_amount":
			model.netAmount = parseFloat64(entity)
		case "total_amount":
			model.totalAmount = parseFloat64(entity)
		case "total_tax_amount":
			model.totalTaxAmount = parseFloat64(entity)
		case "vat_amount":
			model.vatAmount = parseFloat64(entity)
		case "vat_category_code":
			vatCategoryCode := entity.GetMentionText()
			model.vatCategoryCode = &vatCategoryCode
		case "vat_tax_amount":
			model.vatTaxAmount = parseFloat64(entity)
		case "vat_tax_rate":
			model.vatTaxRate = parseFloat64(entity)
		case "line_item":
			i = &Item{}
			i.txt = strPtr(entity.GetMentionText())
			model.items = append(model.items, i)
		case "line_item/description":
			i.description = strPtr(entity.GetMentionText())
		case "line_item/product_code":
			i.productCode = strPtr(entity.GetMentionText())
		case "line_item/purchase_order":
			i.purchaseOrder = strPtr(entity.GetMentionText())
		case "line_item/quantity":
			i.quantity = parseInt(entity)
		case "line_item/unit":
			i.unit = strPtr(entity.GetMentionText())
		case "line_item/unit_price":
			i.unitPrice = parseFloat64(entity)
		}
	}
	return model
}

func strPtr(s string) *string {
	return &s
}

func parseInt(entity *documentaipb.Document_Entity) *int64 {
	var parsed int64
	var err error
	if entity.GetMentionText() != "" {
		parsed, err = strconv.ParseInt(entity.GetMentionText(), 10, 64)
		if err != nil {
			return nil
		}
	}
	return &parsed
}

func parseFloat64(entity *documentaipb.Document_Entity) *float64 {
	var parsed float64
	var err error
	if entity.GetMentionText() != "" {
		parsed, err = strconv.ParseFloat(entity.GetMentionText(), 64)
		if err != nil {
			return nil
		}
	}
	return &parsed
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
