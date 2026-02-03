// Package entities provides type-safe wrappers for common database types.
// It includes custom types for integers, floats, booleans, currencies, and database metadata
// that can be marshalled to and from strings for compatibility with various storage backends,
// particularly Storm database which has limitations with native Go types.
package entities

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/beorn7/floats"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/shopspring/decimal"
)

// Field represents a database field name used for queries and schema definitions.
type (
	field string
	Field field
)

// URI represents a Uniform Resource Identifier for entity references and links.
type (
	uri string
	URI uri
)

// KEY represents a primary key value for database entities.
type (
	pk  string
	KEY pk
)

// Table represents a database table name.
type (
	table string
	Table table
)

// Int is an integer type that can be marshalled to and from a string.
// This allows for safe storage in databases that have issues with native integer types.
type Int struct {
	Value string // String representation of the integer value
}

// Integer type aliases for different bit widths and signedness.
type (
	Int64  Int // 64-bit signed integer
	UInt64 Int // 64-bit unsigned integer
	UInt   Int // Platform-dependent unsigned integer
	Int32  Int // 32-bit signed integer
	UInt32 Int // 32-bit unsigned integer
)

// Float is a floating-point type that can be marshalled to and from a string.
// This allows for safe storage in databases that have issues with native float types.
type Float struct {
	Value string // String representation of the floating-point value
}

// Floating-point type aliases and related financial types.
type (
	Float32  Float    // 32-bit floating-point number
	Float64  Float    // 64-bit floating-point number
	Decimal  Float    // High-precision decimal number
	Currency struct { // Currency value with ISO currency code
		Value Float  // Monetary amount
		CCY   string // 3-letter ISO currency code (e.g., "USD", "GBP")
	}
)

// Financial and rate type aliases.
type (
	Money      Float // Monetary amount
	Percentage Float // Percentage value (e.g., 25.5 for 25.5%)
	Rate       Float // Interest rate or exchange rate
)

// Bool is a boolean type that can be marshalled to and from a string.
// This has been created because Storm database does not support boolean types properly.
type Bool struct {
	Value string // String representation: "true" or "false"
}

// StormBool is an alias for Bool, explicitly indicating compatibility with Storm database.
// Storm database does not support boolean types properly, so this string-backed type is used.
type StormBool Bool

// String returns the Field as a string.
func (f Field) String() string {
	return string(f)
}

// String constants for boolean representation.
const (
	constTrue  = "true"  // String representation of boolean true
	constFalse = "false" // String representation of boolean false
)

// Set sets the StormBool value from a native Go bool.
func (sb *StormBool) Set(b bool) {
	if b {
		sb.Value = constTrue
	} else {
		sb.Value = constFalse
	}
}

// Bool returns the StormBool as a native Go bool.
func (sb *StormBool) Bool() bool {
	return sb.Value == constTrue
}

// String returns the StormBool as a string ("true" or "false").
func (sb *StormBool) String() string {
	return sb.Value
}

// IsTrue returns true if the StormBool represents a true value.
func (sb *StormBool) IsTrue() bool {
	return sb.Bool()
}

// IsFalse returns true if the StormBool represents a false value.
func (sb *StormBool) IsFalse() bool {
	return !sb.Bool()
}

// Set sets the Int value from a native Go int and returns the updated Int.
func (i *Int) Set(in int) Int {
	i.Value = strconv.Itoa(in)
	return *i
}

// Int converts the Int to a native Go int.
// Returns 0 if the value is empty. Panics if the value cannot be parsed.
func (i *Int) Int() int {
	if i.Value == "" {
		return 0
	}
	val, err := strconv.Atoi(i.Value)
	if err != nil {
		logHandler.ErrorLogger.Panic(commonErrors.ErrInvalidTypeWrapper("Int", i.Value, "int"))
	}
	// logHandler.InfoLogger.Printf("val: '%v' int: '%d'", i.Value, val)
	return val
}

// Int64 converts the Int to a 64-bit signed integer.
func (i *Int) Int64() int64 {
	return int64(i.Int())
}

// UInt64 converts the Int to a 64-bit unsigned integer.
func (i *Int) UInt64() uint64 {
	return uint64(i.Int())
}

// UInt converts the Int to a platform-dependent unsigned integer.
func (i *Int) UInt() uint {
	return uint(i.Int())
}

// Int32 converts the Int to a 32-bit signed integer.
func (i *Int) Int32() int32 {
	return int32(i.Int())
}

// UInt32 converts the Int to a 32-bit unsigned integer.
func (i *Int) UInt32() uint32 {
	return uint32(i.Int())
}

// Get returns the Int value as a native Go int.
func (i *Int) Get() int {
	return i.Int()
}

// String returns the Int as its string representation.
func (i *Int) String() string {
	return i.Value
}

// Equals returns true if this Int equals another Int.
func (i *Int) Equals(other Int) bool {
	return i.Value == other.Value
}

// LessThan returns true if this Int is less than another Int.
func (i *Int) LessThan(other Int) bool {
	return i.Int() < other.Int()
}

// LessThanOrEqual returns true if this Int is less than or equal to another Int.
func (i *Int) LessThanOrEqual(other Int) bool {
	return i.Int() <= other.Int()
}

// GreaterThan returns true if this Int is greater than another Int.
func (i *Int) GreaterThan(other Int) bool {
	return i.Int() > other.Int()
}

// GreaterThanOrEqual returns true if this Int is greater than or equal to another Int.
func (i *Int) GreaterThanOrEqual(other Int) bool {
	return i.Int() >= other.Int()
}

// Add adds another Int to this Int and returns the result.
// This method modifies the receiver and returns it.
func (i *Int) Add(other Int) Int {
	sum := i.Int() + other.Int()
	return i.Set(sum)
}

// Subtract subtracts another Int from this Int and returns the result.
// This method modifies the receiver and returns it.
func (i *Int) Subtract(other Int) Int {
	diff := i.Int() - other.Int()
	return i.Set(diff)
}

// MultiplyBy multiplies this Int by another Int and returns the result.
// This method modifies the receiver and returns it.
func (i *Int) MultiplyBy(other Int) Int {
	prod := i.Int() * other.Int()
	return i.Set(prod)
}

// DivideBy divides this Int by another Int and returns the result.
// Panics if attempting to divide by zero.
// This method modifies the receiver and returns it.
func (i *Int) DivideBy(other Int) Int {
	if other.Int() == 0 {
		logHandler.ErrorLogger.Panic(fmt.Errorf("division by zero in Int.Divide"))
		return *i
	}
	quot := i.Int() / other.Int()
	return i.Set(quot)
}

// IncrementBy increases this Int by another Int value and returns the result.
// This method modifies the receiver and returns it.
func (i *Int) IncrementBy(other Int) Int {
	sum := i.Int() + other.Int()
	return i.Set(sum)
}

// Increment increases this Int by 1 and returns the result.
// This method modifies the receiver and returns it.
func (i *Int) Increment() Int {
	sum := i.Int() + 1
	return i.Set(sum)
}

// DecrementBy decreases this Int by another Int value and returns the result.
// This method modifies the receiver and returns it.
func (i *Int) DecrementBy(other Int) Int {
	diff := i.Int() - other.Int()
	return i.Set(diff)
}

// Decrement decreases this Int by 1 and returns the result.
// This method modifies the receiver and returns it.
func (i *Int) Decrement() Int {
	diff := i.Int() - 1
	return i.Set(diff)
}

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
		logHandler.ErrorLogger.Panic(commonErrors.ErrInvalidTypeWrapper("Float", f.Value, "float64"))
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

// Set sets the Bool value from a native Go bool.
func (b *Bool) Set(in bool) {
	if in == true {
		b.Value = constTrue
	} else {
		b.Value = constFalse
	}
}

// SetTrue sets the Bool value to true.
func (b *Bool) SetTrue() {
	b.Value = constTrue
}

// SetFalse sets the Bool value to false.
func (b *Bool) SetFalse() {
	b.Value = constFalse
}

// SetFromString sets the Bool value by parsing a string.
// Recognizes "true", "1", "yes", "y", "t" (case-insensitive) as true, everything else as false.
func (b *Bool) SetFromString(in string) {
	in = strings.ToLower(strings.TrimSpace(in))
	if in == "true" || in == "1" || in == "yes" || in == "y" || in == "t" {
		b.Value = constTrue
	} else {
		b.Value = constFalse
	}
}

// Toggle inverts the Bool value (true becomes false, false becomes true).
func (b *Bool) Toggle() {
	if b.Bool() {
		b.Value = constFalse
	} else {
		b.Value = constTrue
	}
}

// HtmlChecked returns "checked" if the Bool is true, otherwise it returns an empty string
func (b *Bool) HtmlChecked() string {
	if b.Bool() {
		return "checked" // Checked
	}
	return "" // Not Checked
}

// HtmlSelected returns "selected" if the Bool is true, otherwise returns an empty string.
// Useful for HTML select elements.
func (b *Bool) HtmlSelected() string {
	if b.Bool() {
		return "selected" // Selected
	}
	return "" // Not Selected
}

// HtmlDisabled returns "disabled" if the Bool is true, otherwise returns an empty string.
// Useful for HTML form elements.
func (b *Bool) HtmlDisabled() string {
	if b.Bool() {
		return "disabled" // Disabled
	}
	return "" // Not Disabled
}

// HtmlReadOnly returns "readonly" if the Bool is true, otherwise returns an empty string.
// Useful for HTML form elements.
func (b *Bool) HtmlReadOnly() string {
	if b.Bool() {
		return "readonly" // ReadOnly
	}
	return "" // Not ReadOnly
}

// Bool converts the Bool to a native Go bool.
// Returns false if the value is empty, otherwise true only if value equals "true".
func (b *Bool) Bool() bool {
	if b.Value == "" {
		// Needs reversing based on current implementation
		return false
	}
	return b.Value == constTrue
}

// Get returns the Bool value as a native Go bool.
func (b *Bool) Get() bool {
	return b.Bool()
}

// String returns the Bool as a string ("true" or "false").
func (b *Bool) String() string {
	return b.Value
}

// IsTrue returns true if the Bool represents a true value.
func (b *Bool) IsTrue() bool {
	return b.Bool()
}

// IsFalse returns true if the Bool represents a false value.
func (b *Bool) IsFalse() bool {
	return !b.Bool()
}

// String returns the Table name as a string.
func (t *Table) String() string {
	return fmt.Sprintf("%v", *t)
}

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
		logHandler.ErrorLogger.Panic(commonErrors.ErrInvalidTypeWrapper("Currency Code", code, "3-letter ISO currency code"))
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

// String returns the URI as a string.
func (u *URI) String() string {
	return string(*u)
}

// Set sets the URI value from a string.
func (u *URI) Set(in string) {
	*u = URI(in)
}

// Get returns the URI as a string.
func (u *URI) Get() string {
	return u.String()
}

// String returns the primary key as a string.
func (p *KEY) String() string {
	return string(*p)
}

// Set sets the primary key value from a string.
func (p *KEY) Set(in string) {
	*p = KEY(in)
}

// Get returns the primary key as a string.
func (p *KEY) Get() string {
	return p.String()
}

// Key returns the primary key wrapped in curly braces.
// Example: if the key is "123", this returns "{123}".
func (p *KEY) Key() string {
	return "{" + p.String() + "}"
}
