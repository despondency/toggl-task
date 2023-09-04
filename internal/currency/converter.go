package currency

import (
	"encoding/json"
	"github.com/everapihq/currencyapi-go"
	"math/big"
	"time"
)

type ConvertResponse struct {
	Meta Meta                 `json:"meta"`
	Data map[string]*Currency `json:"data"`
}
type Meta struct {
	LastUpdatedAt time.Time `json:"last_updated_at"`
}
type Currency struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

type Converter interface {
	Convert(amount *big.Float, currency string) (*big.Float, error)
}

type ConverterCaller interface {
	ConvertCurrency(req map[string]string) (*ConvertResponse, error)
}

type APICurrencyConverter struct {
	convertCaller ConverterCaller
}

func NewCurrencyConverter(caller ConverterCaller) Converter {
	return &APICurrencyConverter{convertCaller: caller}
}

type APICurrencyConverterCaller struct {
}

func NewAPICurrencyConverterCaller(apiKey string) ConverterCaller {
	if apiKey != "" {
		currencyapi.Init(apiKey)
	}
	return &APICurrencyConverterCaller{}
}

func (m *APICurrencyConverterCaller) ConvertCurrency(req map[string]string) (*ConvertResponse, error) {
	byteResp := currencyapi.Latest(req)
	currencyConvertResp := &ConvertResponse{}
	err := json.Unmarshal(byteResp, currencyConvertResp)
	if err != nil {
		return nil, err
	}
	return currencyConvertResp, nil
}

// Convert we currently exchange all currencies to EUR
func (m *APICurrencyConverter) Convert(amount *big.Float, currency string) (*big.Float, error) {
	convertTo := "EUR"
	req := map[string]string{
		"base_currency": currency,
		"currencies":    convertTo,
	}
	resp, err := m.convertCaller.ConvertCurrency(req)
	if err != nil {
		return nil, err
	}
	exchangeRate := big.NewFloat(resp.Data[convertTo].Value)
	return new(big.Float).Mul(amount, exchangeRate), nil
}
