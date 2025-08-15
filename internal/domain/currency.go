package domain

import (
	"fmt"
	"sync"
	"time"
)

// Currency represents a currency with its properties
type Currency struct {
	Code         string    `json:"code"`          // ISO 4217 code (USD, EUR, TRY)
	Name         string    `json:"name"`          // Full name (US Dollar, Euro, Turkish Lira)
	Symbol       string    `json:"symbol"`        // Symbol ($, €, ₺)
	Decimals     int       `json:"decimals"`      // Decimal places (2 for most currencies)
	ExchangeRate float64   `json:"exchange_rate"` // Rate relative to base currency
	LastUpdated  time.Time `json:"last_updated"`
}

// CurrencyConverter handles currency conversions
type CurrencyConverter struct {
	rates map[string]float64
	mu    sync.RWMutex
}

// NewCurrencyConverter creates a new currency converter
func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]float64),
	}
}

// SetExchangeRate sets the exchange rate for a currency
func (cc *CurrencyConverter) SetExchangeRate(from, to string, rate float64) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	key := fmt.Sprintf("%s_%s", from, to)
	cc.rates[key] = rate
}

// GetExchangeRate gets the exchange rate between two currencies
func (cc *CurrencyConverter) GetExchangeRate(from, to string) (float64, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	if from == to {
		return 1.0, nil
	}

	key := fmt.Sprintf("%s_%s", from, to)
	if rate, exists := cc.rates[key]; exists {
		return rate, nil
	}

	// Try reverse rate
	reverseKey := fmt.Sprintf("%s_%s", to, from)
	if rate, exists := cc.rates[reverseKey]; exists {
		return 1.0 / rate, nil
	}

	return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
}

// Convert converts an amount from one currency to another
func (cc *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	rate, err := cc.GetExchangeRate(from, to)
	if err != nil {
		return 0, err
	}

	return amount * rate, nil
}

// MultiCurrencyAmount represents an amount in multiple currencies
type MultiCurrencyAmount struct {
	Amount    float64            `json:"amount"`
	Currency  string             `json:"currency"`
	Converted map[string]float64 `json:"converted,omitempty"`
}

// ConvertTo converts the amount to multiple currencies
func (mca *MultiCurrencyAmount) ConvertTo(converter *CurrencyConverter, targetCurrencies []string) error {
	if mca.Converted == nil {
		mca.Converted = make(map[string]float64)
	}

	for _, target := range targetCurrencies {
		if target == mca.Currency {
			mca.Converted[target] = mca.Amount
			continue
		}

		converted, err := converter.Convert(mca.Amount, mca.Currency, target)
		if err != nil {
			return fmt.Errorf("failed to convert to %s: %w", target, err)
		}

		mca.Converted[target] = converted
	}

	return nil
}

// Supported currencies
var (
	USD = &Currency{Code: "USD", Name: "US Dollar", Symbol: "$", Decimals: 2}
	EUR = &Currency{Code: "EUR", Name: "Euro", Symbol: "€", Decimals: 2}
	TRY = &Currency{Code: "TRY", Name: "Turkish Lira", Symbol: "₺", Decimals: 2}
	GBP = &Currency{Code: "GBP", Name: "British Pound", Symbol: "£", Decimals: 2}
	JPY = &Currency{Code: "JPY", Name: "Japanese Yen", Symbol: "¥", Decimals: 0}
)

// GetSupportedCurrencies returns all supported currencies
func GetSupportedCurrencies() []*Currency {
	return []*Currency{USD, EUR, TRY, GBP, JPY}
}

// ValidateCurrencyCode validates if a currency code is supported
func ValidateCurrencyCode(code string) bool {
	for _, currency := range GetSupportedCurrencies() {
		if currency.Code == code {
			return true
		}
	}
	return false
}

