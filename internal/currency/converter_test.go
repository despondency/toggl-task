package currency_test

import (
	"fmt"
	"github.com/despondency/toggl-task/internal/currency"
	currencymock "github.com/despondency/toggl-task/mocks/currency"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestUnitConverter_Convert(t *testing.T) {
	testCases := []struct {
		name           string
		expectedErr    error
		amount         *big.Float
		expected       *big.Float
		createInstance func(t *testing.T) currency.Converter
	}{
		{
			name:     "simple case, no error",
			amount:   big.NewFloat(100),
			expected: big.NewFloat(89),
			createInstance: func(t *testing.T) currency.Converter {
				currencyCaller := currencymock.NewConverterCaller(t)
				currencyCaller.EXPECT().ConvertCurrency(map[string]string{
					"base_currency": "USD",
					"currencies":    "EUR",
				}).Return(&currency.ConvertResponse{
					Meta: currency.Meta{
						LastUpdatedAt: time.Now(),
					},
					Data: map[string]*currency.Currency{
						"EUR": {
							Code:  "EUR",
							Value: 0.89,
						},
					},
				}, nil)
				return currency.NewCurrencyConverter(currencyCaller)
			},
		},
		{
			name:        "simple case, error",
			expectedErr: fmt.Errorf("error calling currencyapi"),
			createInstance: func(t *testing.T) currency.Converter {
				currencyCaller := currencymock.NewConverterCaller(t)
				currencyCaller.EXPECT().ConvertCurrency(map[string]string{
					"base_currency": "USD",
					"currencies":    "EUR",
				}).Return(&currency.ConvertResponse{
					Meta: currency.Meta{
						LastUpdatedAt: time.Now(),
					},
					Data: map[string]*currency.Currency{
						"EUR": {
							Code:  "EUR",
							Value: 0.89,
						},
					},
				}, fmt.Errorf("error calling currencyapi"))
				return currency.NewCurrencyConverter(currencyCaller)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			instance := tc.createInstance(t)
			resp, err := instance.Convert(tc.amount, "USD")
			if tc.expectedErr != nil {
				assert.EqualErrorf(t, err, tc.expectedErr.Error(), "")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 0, tc.expected.Cmp(resp), fmt.Sprintf("expected %s, to equal %s", tc.expected.String(), resp.String()))
			}
		})
	}
}
