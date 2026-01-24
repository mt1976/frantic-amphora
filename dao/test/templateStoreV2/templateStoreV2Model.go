// Data Access Object template
// Version: 0.5.0
// Updated on: 2026-01-24

package templateStoreV2

import (
	"time"

	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/entities"
)

// TableName is the canonical DAO table identifier for this package.
var TableName = entities.Table("TemplateStore")

var tableName = TableName.String()

// TemplateStore represents a sample entity for demonstrating reduced DAO boilerplate.
// Replace this struct and Fields as needed for your real entity.
type TemplateStore struct {
	// The primary key field(s), managed by the framework, DO NOT MODIFY
	ID  int    `storm:"id,increment=100"`
	Key string `storm:"index,unique"`
	Raw string `storm:"index,unique"`
	// Audit information, managed by the framework, DO NOT MODIFY
	Audit audit.Audit `csv:"-"`

	// Domain specific fields
	UID      string `validate:"required"`
	GID      string `storm:"index" validate:"required"`
	RealName string `validate:"required,min=5"`
	UserName string `validate:"required,min=5"`
	UserCode string `storm:"index" validate:"required,min=5"`
	Email    string
	Notes    string `validate:"max=75"`
	Active   entities.Bool

	// Example fields grouped by type
	// String types
	ExampleString string
	// Boolean types
	ExampleBool      entities.Bool
	ExampleStormBool entities.StormBool
	// Integer types
	ExampleInt    entities.Int
	ExampleInt32  entities.Int32
	ExampleInt64  entities.Int64
	ExampleUint   entities.UInt
	ExampleUint32 entities.UInt32
	ExampleUint64 entities.UInt64
	// Float types
	ExampleFloat   entities.Float
	ExampleFloat32 entities.Float32
	ExampleFloat64 entities.Float64
	// Specialized numeric types
	ExampleDecimal    entities.Decimal
	ExamplePercentage entities.Percentage
	ExampleRate       entities.Rate
	// Money types
	ExampleMoney    entities.Money
	ExampleCurrency entities.Currency
	// Date/Time types
	ExampleDate time.Time
	// Entity framework types
	ExampleField entities.Field
	ExampleTable entities.Table

	LastLogin time.Time
	LastHost  string `storm:"index"`
}

type fieldNames struct {
	// The primary key field(s), managed by the framework, DO NOT MODIFY
	ID  entities.Field
	Key entities.Field
	Raw entities.Field
	// The audit information, managed by the framework, DO NOT MODIFY
	Audit entities.Field
	// Domain specific fields, please modify as required
	UID      entities.Field
	GID      entities.Field
	RealName entities.Field
	UserName entities.Field
	UserCode entities.Field
	Email    entities.Field
	Notes    entities.Field
	Active   entities.Field
	// Example fields
	ExampleString     entities.Field
	ExampleBool       entities.Field
	ExampleStormBool  entities.Field
	ExampleInt        entities.Field
	ExampleInt32      entities.Field
	ExampleInt64      entities.Field
	ExampleUint       entities.Field
	ExampleUint32     entities.Field
	ExampleUint64     entities.Field
	ExampleFloat      entities.Field
	ExampleFloat32    entities.Field
	ExampleFloat64    entities.Field
	ExampleDecimal    entities.Field
	ExamplePercentage entities.Field
	ExampleRate       entities.Field
	ExampleMoney      entities.Field
	ExampleCurrency   entities.Field
	ExampleDate       entities.Field
	ExampleField      entities.Field
	ExampleTable      entities.Field
	LastLogin         entities.Field
	LastHost          entities.Field
}

// Fields provides strongly-typed field names for use with GetBy/GetAllWhere/etc.
//
// Example: GetBy(Fields.Key, "abc")
//
// Note: the values are the struct field names as stored in Storm.
var Fields = fieldNames{
	// The primary key field(s), managed by the framework, DO NOT MODIFY
	ID:  "ID",
	Key: "Key",
	Raw: "Raw",
	// The audit information, managed by the framework, DO NOT MODIFY
	Audit: "Audit",
	// tableName-specific fields, please modify as required
	UID:      "UID",
	GID:      "GID",
	RealName: "RealName",
	UserName: "UserName",
	UserCode: "UserCode",
	Email:    "Email",
	Notes:    "Notes",
	Active:   "Active",
	// Example fields
	ExampleString:     "ExampleString",
	ExampleBool:       "ExampleBool",
	ExampleStormBool:  "ExampleStormBool",
	ExampleInt:        "ExampleInt",
	ExampleInt32:      "ExampleInt32",
	ExampleInt64:      "ExampleInt64",
	ExampleUint:       "ExampleUint",
	ExampleUint32:     "ExampleUint32",
	ExampleUint64:     "ExampleUint64",
	ExampleFloat:      "ExampleFloat",
	ExampleFloat32:    "ExampleFloat32",
	ExampleFloat64:    "ExampleFloat64",
	ExampleDecimal:    "ExampleDecimal",
	ExamplePercentage: "ExamplePercentage",
	ExampleRate:       "ExampleRate",
	ExampleMoney:      "ExampleMoney",
	ExampleCurrency:   "ExampleCurrency",
	ExampleDate:       "ExampleDate",
	ExampleField:      "ExampleField",
	ExampleTable:      "ExampleTable",
	LastLogin:         "LastLogin",
	LastHost:          "LastHost",
}
