package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Receipt struct {
	Id              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Currency        *string            `json:"currency,omitempty" bson:"currency"`
	DeliveryDate    *string            `json:"delivery_date,omitempty" bson:"delivery_date"`
	DueDate         *string            `json:"due_date" bson:"due_date"`
	InvoiceDate     *string            `json:"invoice_data" bson:"invoice_date"`
	ReceiptDate     *string            `json:"receipt_date" bson:"receipt_date"`
	NetAmount       *string            `json:"net_amount" bson:"net_amount"`
	TotalAmount     *string            `json:"total_amount" bson:"total_amount"`
	TotalTaxAmount  *string            `json:"total_tax_amount" bson:"total_tax_amount"`
	VatAmount       *string            `json:"vat_amount" bson:"vat_amount"`
	VatCategoryCode *string            `json:"vat_category_code" bson:"vat_category_code"`
	VatTaxAmount    *string            `json:"vat_tax_amount" bson:"vat_tax_amount"`
	VatTaxRate      *string            `json:"vat_tax_rate" bson:"vat_tax_rate"`
	Tags            []string           `json:"tags" bson:"tags"`
	Items           []*Item            `json:"items" bson:"items"`
}

type Item struct {
	Amount        *int    `json:"amount" bson:"amount"`
	Txt           *string `json:"txt" bson:"txt"`
	Description   *string `json:"description" bson:"description"`
	ProductCode   *string `json:"product_code" bson:"product_code"`
	PurchaseOrder *string `json:"purchase_order" bson:"purchase_order"`
	Quantity      *string `json:"quantity" bson:"quantity"`
	Unit          *string `json:"unit" bson:"unit"`
	UnitPrice     *string `json:"unit_price" bson:"unit_price"`
}
