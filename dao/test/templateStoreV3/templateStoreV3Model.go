// Data Access Object for the TemplateStoreV3 table
// Template Version: 0.5.10 - 2026-01-26
// Generated
// Date: 27/01/2026 & 15:01
// Who : matttownsend (orion)

package templateStoreV3

import (
	"time"

	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/entities"
)

// TableName is the canonical DAO table identifier for this package.
var (
	TableName = entities.Table("TemplateStoreV3")
	tableName = TableName.String()
)

// The TemplateStoreV3 struct defines the data model for the TemplateStoreV3 table.
// Adjust domain fields and tags as required in the TemplateStoreV3.definitions file.
type TemplateStoreV3 struct {
	// The primary key field(s), managed by the framework, DO NOT MODIFY
	ID  int    `storm:"id,increment=100"`
	Key string `storm:"index,unique"`
	Raw string `storm:"index,unique"`
	// Audit information, managed by the framework, DO NOT MODIFY
	Audit audit.Audit `csv:"-"`

	// Domain specific fields
	//
	// Example fields grouped by type
	// String types
	ExampleString string // Example string field
	// Boolean types
	ExampleBool      entities.Bool      // Example boolean field
	ExampleStormBool entities.StormBool // Example storm boolean field
	// Integer types
	ExampleInt    entities.Int    // Example integer field
	ExampleInt32  entities.Int32  // Example int32 field
	ExampleInt64  entities.Int64  // Example int64 field
	ExampleUint   entities.UInt   // Example unsigned integer field
	ExampleUint32 entities.UInt32 // Example unsigned int32 field
	ExampleUint64 entities.UInt64 // Example unsigned int64 field
	// Float types
	ExampleFloat   entities.Float   // Example float field
	ExampleFloat32 entities.Float32 // Example float32 field
	ExampleFloat64 entities.Float64 // Example float64 field
	// Specialized numeric types
	ExampleDecimal    entities.Decimal    // Example decimal field
	ExamplePercentage entities.Percentage // Example percentage field
	ExampleRate       entities.Rate       // Example rate field
	// Money types
	ExampleMoney    entities.Money    // Example money field
	ExampleCurrency entities.Currency // Example currency field
	// Date/Time types
	ExampleDate time.Time // Example date field
	// Entity framework types
	ExampleField entities.Field // Example field type1
	ExampleTable entities.Table // Example table type1
	// User Management fields
	UID      string `validate:"required"`
	GID      string `storm:"index" validate:"required"`
	RealName string `validate:"required,min=5"`
	UserName string `validate:"required,min=5"`
	UserCode string `storm:"index" validate:"required,min=5"`
	Email    string
	Notes    string `validate:"max=75"`
	Active   entities.Bool

	LastLogin time.Time // Last login time
	LastHost  string    `storm:"index"` // Last host with index
	// Add no more fields below this line
	PostTest []string
}

type fieldNames struct {
	// The primary key field(s), managed by the framework, DO NOT MODIFY
	ID  entities.Field
	Key entities.Field
	Raw entities.Field
	// The audit information, managed by the framework, DO NOT MODIFY
	Audit entities.Field
	// Domain specific fields
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
	UID               entities.Field
	GID               entities.Field
	RealName          entities.Field
	UserName          entities.Field
	UserCode          entities.Field
	Email             entities.Field
	Notes             entities.Field
	Active            entities.Field
	LastLogin         entities.Field
	LastHost          entities.Field

	// Add no more fields below this line
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
	UID:               "UID",
	GID:               "GID",
	RealName:          "RealName",
	UserName:          "UserName",
	UserCode:          "UserCode",
	Email:             "Email",
	Notes:             "Notes",
	Active:            "Active",
	LastLogin:         "LastLogin",
	LastHost:          "LastHost",
	// Add no more fields below this line
}
