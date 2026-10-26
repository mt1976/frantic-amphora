// Package entities provides type-safe wrappers for common database types.
// It includes custom types for integers, floats, booleans, currencies, and database metadata
// that can be marshalled to and from strings for compatibility with various storage backends,
// particularly Storm database which has limitations with native Go types.
package entities

import "time"

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

// ACTION represents a web action type (e.g., CREATE, UPDATE, DELETE).
type (
	action string
	ACTION action
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

// String constants for boolean representation.
const (
	constTrue  = "true"  // String representation of boolean true
	constFalse = "false" // String representation of boolean false
)

// Date is a calendar date stored as midnight UTC in a time.Time value.
type Date struct {
	Value time.Time // Calendar date represented at midnight UTC
}

// Time is a time of day stored in a UTC time.Time value.
type Time struct {
	Value time.Time // UTC time, including its reference date
}

// DateTime is a point in time stored as a UTC time.Time value.
type DateTime struct {
	Value time.Time // UTC timestamp
}
