// This file contains the currency entity methods.

package entities

import (
	"fmt"
	"strings"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/logHandler"
)

// Set sets both the currency code and value, and returns the updated Currency.
func (c *Currency) Set(code string, value float64) Currency {
	c.SetValue(value)
	c.SetCode(code)
	return *c
}

// SetValue sets the monetary value of the Currency.
func (c *Currency) SetValue(value float64) {
	c.Value.Set(value)
}

// GetCode returns the 3-letter ISO currency code.
func (c *Currency) GetCode() string {
	return c.CCY
}

// SetCode sets the 3-letter ISO currency code.
// Defaults to "GBP" if an empty code is provided.
// Panics if the code is not exactly 3 characters long.
func (c *Currency) SetCode(code string) {
	if code == "" {
		// That is, no currency specified, so we default to GBP
		// That will teach Trump!
		// (Just kidding, of course. :-) )
		code = "GBP"
	}
	if len(code) != 3 {
		logHandler.Error.Panic(commonErrors.ErrInvalidTypeWrapper("Currency Code", code, "3-letter ISO currency code"))
	}
	code = strings.ToUpper(code)
	c.CCY = code
}

// New creates a new Currency with zero value and default currency code (GBP).
func (c *Currency) New() Currency {
	x := Currency{}
	x.SetValue(0.0)
	x.SetCode("")
	return x
}

// NewCurrency creates a new Currency with zero value and the specified currency code.
func (c *Currency) NewCurrency(code string) Currency {
	x := Currency{}
	x.SetValue(0.0)
	x.SetCode(code)
	return x
}

// NewAmount creates a new Currency with the specified value and default currency code (GBP).
func (c *Currency) NewAmount(value float64) Currency {
	x := Currency{}
	x.SetValue(value)
	x.SetCode("")
	return x
}

// Code returns the 3-letter ISO currency code.
// This is an alias for GetCode().
func (c *Currency) Code() string {
	return c.GetCode()
}

// GetValue returns the monetary value as a float64.
func (c *Currency) GetValue() float64 {
	return c.Value.Float()
}

// Amount returns the currency amount as a float64.
// This is an alias for GetValue().
func (c *Currency) Amount() float64 {
	return c.GetValue()
}

// String returns the currency as a string in the format "CCY amount".
// Example: "USD 123.45"
func (c *Currency) String() string {
	return fmt.Sprintf("%s %.2f", c.CCY, c.Value.Float())
}

// Get returns both the currency code and value as a tuple.
func (c *Currency) Get() (string, float64) {
	return c.CCY, c.Value.Float()
}
