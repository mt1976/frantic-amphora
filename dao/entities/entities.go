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

// Field represents a database field used for queries
type field string
type Field field

// Table represents a database table
type table string
type Table table

// Int is an integer type that can be marshalled to and from a string
type Int struct {
	Value string
}
type Int64 Int
type UInt64 Int
type UInt Int
type Int32 Int
type UInt32 Int

// Float is a float type that can be marshalled to and from a string
type Float struct {
	Value string
}

type Float32 Float
type Float64 Float
type Decimal Float
type Currency struct {
	Value Float
	CCY   string
}
type Money Float
type Percentage Float
type Rate Float

// Bool is a boolean type that can be marshalled to and from a string, this has been created as Storm does not support boolean types properly
type Bool struct {
	Value string
}

// StormBool is a boolean type that can be marshalled to and from a string, this has been created as Storm does not support boolean types properly
type StormBool Bool

func (f Field) String() string {
	return string(f)
}

const constTrue = "true"
const constFalse = "false"

func (sb *StormBool) Set(b bool) {
	if b {
		sb.Value = constTrue
	} else {
		sb.Value = constFalse
	}
}

func (sb *StormBool) Bool() bool {
	return sb.Value == constTrue
}

func (sb *StormBool) String() string {
	return sb.Value
}

func (sb *StormBool) IsTrue() bool {
	return sb.Bool()
}

func (sb *StormBool) IsFalse() bool {
	return !sb.Bool()
}

func (i *Int) Set(in int) Int {
	i.Value = strconv.Itoa(in)
	return *i
}

func (i *Int) Int() int {
	if i.Value == "" {
		return 0
	}
	val, err := strconv.Atoi(i.Value)
	if err != nil {
		logHandler.ErrorLogger.Panic(commonErrors.ErrInvalidTypeWrapper("Int", i.Value, "int"))
	}
	//logHandler.InfoLogger.Printf("val: '%v' int: '%d'", i.Value, val)
	return val
}

func (i *Int) Int64() int64 {
	return int64(i.Int())
}

func (i *Int) UInt64() uint64 {
	return uint64(i.Int())
}

func (i *Int) UInt() uint {
	return uint(i.Int())
}

func (i *Int) Int32() int32 {
	return int32(i.Int())
}

func (i *Int) UInt32() uint32 {
	return uint32(i.Int())
}

func (i *Int) Get() int {
	return i.Int()
}

func (i *Int) String() string {
	return i.Value
}

func (i *Int) Equals(other Int) bool {
	return i.Value == other.Value
}

func (i *Int) LessThan(other Int) bool {
	return i.Int() < other.Int()
}

func (i *Int) LessThanOrEqual(other Int) bool {
	return i.Int() <= other.Int()
}

func (i *Int) GreaterThan(other Int) bool {
	return i.Int() > other.Int()
}

func (i *Int) GreaterThanOrEqual(other Int) bool {
	return i.Int() >= other.Int()
}

func (i *Int) Add(other Int) Int {
	sum := i.Int() + other.Int()
	return i.Set(sum)
}

func (i *Int) Subtract(other Int) Int {
	diff := i.Int() - other.Int()
	return i.Set(diff)
}

func (i *Int) MultiplyBy(other Int) Int {
	prod := i.Int() * other.Int()
	return i.Set(prod)
}

func (i *Int) DivideBy(other Int) Int {
	if other.Int() == 0 {
		logHandler.ErrorLogger.Panic(fmt.Errorf("division by zero in Int.Divide"))
		return *i
	}
	quot := i.Int() / other.Int()
	return i.Set(quot)
}

func (i *Int) IncrementBy(other Int) Int {
	sum := i.Int() + other.Int()
	return i.Set(sum)
}

func (i *Int) Increment() Int {
	sum := i.Int() + 1
	return i.Set(sum)
}

func (i *Int) DecrementBy(other Int) Int {
	diff := i.Int() - other.Int()
	return i.Set(diff)
}

func (i *Int) Decrement() Int {
	diff := i.Int() - 1
	return i.Set(diff)
}

func (f *Float) Set(in float64) Float {
	f.Value = strconv.FormatFloat(in, 'f', -1, 64)
	return *f
}

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

func (f *Float) Get() float64 {
	return f.Float()
}

func (f *Float) String() string {
	return f.Value
}

func (f *Float) Float32() float32 {
	return float32(f.Float())
}

func (f *Float) Float64() float64 {
	return f.Float()
}

func (f *Float) Decimal() decimal.Decimal {
	return decimal.NewFromFloat(f.Float())
}

func (f *Float) Currency() decimal.Decimal {
	return decimal.NewFromFloat(f.Float())
}

func (f *Float) Money() decimal.Decimal {
	return decimal.NewFromFloat(f.Float())
}

func (f *Float) Percentage() decimal.Decimal {
	return decimal.NewFromFloat(f.Float())
}

func (f *Float) Equals(other Float) bool {
	return floats.AlmostEqual(f.Float64(), other.Float64(), floats.MinNormal)
}

func (f *Float) LessThan(other Float) bool {
	return f.Float64() < other.Float64()
}

func (f *Float) LessThanOrEqual(other Float) bool {
	return f.Float64() <= other.Float64()
}

func (f *Float) GreaterThan(other Float) bool {
	return f.Float64() > other.Float64()
}
func (f *Float) GreaterThanOrEqual(other Float) bool {
	return f.Float64() >= other.Float64()
}

func (b *Bool) Set(in bool) {
	if in == true {
		b.Value = constTrue
	} else {
		b.Value = constFalse
	}
}

func (b *Bool) SetTrue() {
	b.Value = constTrue
}

func (b *Bool) SetFalse() {
	b.Value = constFalse
}

func (b *Bool) SetFromString(in string) {
	in = strings.ToLower(strings.TrimSpace(in))
	if in == "true" || in == "1" || in == "yes" || in == "y" || in == "t" {
		b.Value = constTrue
	} else {
		b.Value = constFalse
	}
}

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

func (b *Bool) HtmlSelected() string {
	if b.Bool() {
		return "selected" // Selected
	}
	return "" // Not Selected
}

func (b *Bool) HtmlDisabled() string {
	if b.Bool() {
		return "disabled" // Disabled
	}
	return "" // Not Disabled
}

func (b *Bool) HtmlReadOnly() string {
	if b.Bool() {
		return "readonly" // ReadOnly
	}
	return "" // Not ReadOnly
}

func (b *Bool) Bool() bool {
	if b.Value == "" {
		// Needs reversing based on current implementation
		return false
	}
	return b.Value == constTrue
}

func (b *Bool) Get() bool {
	return b.Bool()
}

func (b *Bool) String() string {
	return b.Value
}

func (b *Bool) IsTrue() bool {
	return b.Bool()
}

func (b *Bool) IsFalse() bool {
	return !b.Bool()
}

func (t *Table) String() string {
	return fmt.Sprintf("%v", *t)
}

func (c *Currency) Set(code string, value float64) Currency {
	c.SetValue(value)
	c.SetCode(code)
	return *c
}

func (c *Currency) SetValue(value float64) {
	c.Value.Set(value)
}

func (c *Currency) GetCode() string {
	return c.CCY
}

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

func (c *Currency) New() Currency {
	x := Currency{}
	x.SetValue(0.0)
	x.SetCode("")
	return x
}

func (c *Currency) NewCurrency(code string) Currency {
	x := Currency{}
	x.SetValue(0.0)
	x.SetCode(code)
	return x
}

func (c *Currency) NewAmount(value float64) Currency {
	x := Currency{}
	x.SetValue(value)
	x.SetCode("")
	return x
}

func (c *Currency) Code() string {
	return c.GetCode()
}
func (c *Currency) GetValue() float64 {
	return c.Value.Float()
}

func (c *Currency) Amount() float64 {
	return c.GetValue()
}

func (c *Currency) String() string {
	return fmt.Sprintf("%s %.2f", c.CCY, c.Value.Float())
}

func (c *Currency) Get() (string, float64) {
	return c.CCY, c.Value.Float()
}
