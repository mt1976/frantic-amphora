// This file contains the integer entity methods.

package entities

import (
	"fmt"
	"strconv"

	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/logHandler"
)

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
		logHandler.Error.Panic(commonErrors.ErrInvalidTypeWrapper("Int", i.Value, "int"))
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
		logHandler.Error.Panic(fmt.Errorf("division by zero in Int.Divide"))
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
