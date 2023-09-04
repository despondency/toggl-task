package currency

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"
)

const baseCurrencyAPIURL = "https://api.currencyapi.com/v3/"

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
	client *http.Client
	apiKey string
}

func WithAPIKey(apiKey string) APICurrencyConverterCallerOptions {
	return func(a *APICurrencyConverterCaller) {
		a.apiKey = apiKey
	}
}

func WithHTTPClient(client *http.Client) APICurrencyConverterCallerOptions {
	return func(a *APICurrencyConverterCaller) {
		a.client = client
	}
}

type APICurrencyConverterCallerOptions func(*APICurrencyConverterCaller)

func NewAPICurrencyConverterCaller(ops ...APICurrencyConverterCallerOptions) (ConverterCaller, error) {
	caller := &APICurrencyConverterCaller{}
	for _, op := range ops {
		op(caller)
	}
	if caller.apiKey == "" {
		return nil, fmt.Errorf("no api key present")
	}
	if caller.client == nil {
		caller.client = http.DefaultClient
	}
	return caller, nil
}

func (m *APICurrencyConverterCaller) ConvertCurrency(req map[string]string) (*ConvertResponse, error) {
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest(http.MethodGet, baseCurrencyAPIURL+"latest", bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("apikey", m.apiKey)
	response, err := m.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	currencyConvertResp := &ConvertResponse{}
	bodyBytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, currencyConvertResp)
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
