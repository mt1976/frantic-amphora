package testentity

import (
	"time"

	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/entities"
)

// TableName is the canonical DAO table identifier for this package.
var TableName = entities.Table("TestEntityTable")

var tableName = TableName.String()

// TestEntity represents an example entity.
//
// This is generated from the TemplateStoreV2 template set. Adjust domain fields and tags as required.
type TestEntity struct {
	// The primary key field(s), managed by the framework, DO NOT MODIFY
	ID  int    `storm:"id,increment=100"`
	Key string `storm:"index,unique"`
	Raw string `storm:"index,unique"`
	//
	// Domain specific fields, starts.
    //
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
    // Additional fields for demonstration
	LastLogin time.Time
	LastHost  string `storm:"index"`
	// Domain specific fields, ends.	
	//
	// Audit information, managed by the framework, DO NOT MODIFY
	Audit audit.Audit `csv:"-"`	
}

type fieldNames struct {
	// The primary key field(s), managed by the framework, DO NOT MODIFY
	ID  entities.Field
	Key entities.Field
	Raw entities.Field
	//
	// Domain specific fields, starts.
	ExampleString entities.Field
	ExampleBool entities.Field
	ExampleStormBool entities.Field
	ExampleInt entities.Field
	ExampleInt32 entities.Field
	ExampleInt64 entities.Field
	ExampleUint entities.Field
	ExampleUint32 entities.Field
	ExampleUint64 entities.Field
	ExampleFloat entities.Field
	ExampleFloat32 entities.Field
	ExampleFloat64 entities.Field
	ExampleDecimal entities.Field
	ExamplePercentage entities.Field
	ExampleRate entities.Field
	ExampleMoney entities.Field
	ExampleCurrency entities.Field
	ExampleDate entities.Field
	ExampleField entities.Field
	ExampleTable entities.Field
	LastLogin entities.Field
	LastHost entities.Field
	// Domain specific fields, ends.
	//
	// The audit information, managed by the framework, DO NOT MODIFY
	Audit entities.Field
	
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
	//
	// Domain specific fields, starts.
	ExampleString: "ExampleString",
	ExampleBool: "ExampleBool",
	ExampleStormBool: "ExampleStormBool",
	ExampleInt: "ExampleInt",
	ExampleInt32: "ExampleInt32",
	ExampleInt64: "ExampleInt64",
	ExampleUint: "ExampleUint",
	ExampleUint32: "ExampleUint32",
	ExampleUint64: "ExampleUint64",
	ExampleFloat: "ExampleFloat",
	ExampleFloat32: "ExampleFloat32",
	ExampleFloat64: "ExampleFloat64",
	ExampleDecimal: "ExampleDecimal",
	ExamplePercentage: "ExamplePercentage",
	ExampleRate: "ExampleRate",
	ExampleMoney: "ExampleMoney",
	ExampleCurrency: "ExampleCurrency",
	ExampleDate: "ExampleDate",
	ExampleField: "ExampleField",
	ExampleTable: "ExampleTable",
	LastLogin: "LastLogin",
	LastHost: "LastHost",
	// Domain specific fields, ends.
	//
	// The audit information, managed by the framework, DO NOT MODIFY
	Audit: "Audit",
}
