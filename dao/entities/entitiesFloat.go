// This file contains the floating-point entity methods.

package entities

import (
	"strconv"

	"github.com/beorn7/floats"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/shopspring/decimal"
)

// Set sets the Float value from a native Go float64 and returns the updated Float.
func (f *Float) Set(in float64) Float {
	f.Value = strconv.FormatFloat(in, 'f', -1, 64)
	return *f
}

// Float converts the Float to a native Go float64.
// Returns 0.0 if the value is empty. Panics if the value cannot be parsed.
func (f *Float) Float() float64 {
	if f.Value == "" {
		return 0.0
	}
	val, err := strconv.ParseFloat(f.Value, 64)
	if err != nil {
		logHandler.Error.Panic(commonErrors.ErrInvalidTypeWrapper("Float", f.Value, "float64"))
	}
	return val
}

// Get returns the Float value as a native Go float64.
func (f *Float) Get() float64 {
	return f.Float()
}

// String returns the Float as its string representation.
func (f *Float) String() string {
	return f.Value
}

// Float32 converts the Float to a 32-bit floating-point number.
func (f *Float) Float32() float32 {
	return float32(f.Float())
}

// Float64 converts the Float to a 64-bit floating-point number.
func (f *Float) Float64() float64 {
	return f.Float()
}

// Decimal converts the Float to a high-precision decimal.Decimal type.
func (f *Float) Decimal() decimal.Decimal {
	return decimal.NewFromFloat(f.Float())
}

// Currency converts the Float to a decimal.Decimal for currency calculations.
func (f *Float) Currency() decimal.Decimal {
	return decimal.NewFromFloat(f.Float())
}

// Money converts the Float to a decimal.Decimal for monetary calculations.
func (f *Float) Money() decimal.Decimal {
	return decimal.NewFromFloat(f.Float())
}

// Percentage converts the Float to a decimal.Decimal for percentage calculations.
func (f *Float) Percentage() decimal.Decimal {
	return decimal.NewFromFloat(f.Float())
}

// Equals returns true if this Float is almost equal to another Float.
// Uses floating-point tolerance to account for precision issues.
func (f *Float) Equals(other Float) bool {
	return floats.AlmostEqual(f.Float64(), other.Float64(), floats.MinNormal)
}

// LessThan returns true if this Float is less than another Float.
func (f *Float) LessThan(other Float) bool {
	return f.Float64() < other.Float64()
}

// LessThanOrEqual returns true if this Float is less than or equal to another Float.
func (f *Float) LessThanOrEqual(other Float) bool {
	return f.Float64() <= other.Float64()
}

// GreaterThan returns true if this Float is greater than another Float.
func (f *Float) GreaterThan(other Float) bool {
	return f.Float64() > other.Float64()
}

// GreaterThanOrEqual returns true if this Float is greater than or equal to another Float.
func (f *Float) GreaterThanOrEqual(other Float) bool {
	return f.Float64() >= other.Float64()
}
