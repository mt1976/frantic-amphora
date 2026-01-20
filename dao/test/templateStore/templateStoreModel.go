package templateStore

// Data Access Object template
// Version: 0.3.0
// Updated on: 2025-12-31

//TODO: RENAME "template" TO THE NAME OF THE DOMAIN ENTITY
//TODO: Update the template_Store struct to match the domain entity
//TODO: Update the Fields. constants to match the domain entity

import (
	"time"

	audit "github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/entities"
)

var Domain = "Template"
var TableName = Domain + "Store"

// TemplateStore represents a User entity.
type TemplateStore struct {
	// First three fields are mandatory for all DAO entities
	ID  int    `storm:"id,increment=100"` // primary key with auto increment
	Key string `storm:"unique"`           // key
	Raw string `storm:"unique"`           // raw ID before encoding
	// Add your domain entity fields below
	UID           string `validate:"required"`
	GID           string `storm:"index" validate:"required"`
	RealName      string `validate:"required,min=5"` // this field will not be indexed
	UserName      string `validate:"required,min=5"`
	UserCode      string `storm:"index" validate:"required,min=5"`
	Email         string
	Notes         string `validate:"max=75"`
	Active        entities.Bool
	ExampleInt    entities.Int
	ExampleFloat  entities.Float
	ExampleBool   entities.Bool
	ExampleDate   time.Time
	ExampleString string
	LastLogin     time.Time
	LastHost      string `storm:"index"`
	// Last field is mandatory for all DAO entities
	Audit audit.Audit `csv:"-"` // audit data
}

// Fields provides a structured way to reference model field names.
type fieldNames struct {
	// First four fields are mandatory for all DAO entities
	ID    entities.Field
	Key   entities.Field
	Raw   entities.Field
	Audit entities.Field
	// Add your domain entity fields below
	UID           entities.Field
	GID           entities.Field
	RealName      entities.Field
	UserName      entities.Field
	UserCode      entities.Field
	Email         entities.Field
	Notes         entities.Field
	Active        entities.Field
	ExampleInt    entities.Field
	ExampleFloat  entities.Field
	ExampleBool   entities.Field
	ExampleDate   entities.Field
	ExampleString entities.Field
	LastLogin     entities.Field
	LastHost      entities.Field
}

// Fields provides a structured way to reference model field names.
var Fields = fieldNames{
	// First four fields are mandatory for all DAO entities
	ID:    "ID",
	Key:   "Key",
	Raw:   "Raw",
	Audit: "Audit",
	// Add your domain entity fields below
	UID:           "UID",
	GID:           "GID",
	RealName:      "RealName",
	UserName:      "UserName",
	UserCode:      "UserCode",
	Email:         "Email",
	Notes:         "Notes",
	Active:        "Active",
	ExampleInt:    "ExampleInt",
	ExampleFloat:  "ExampleFloat",
	ExampleBool:   "ExampleBool",
	ExampleDate:   "ExampleDate",
	ExampleString: "ExampleString",
	LastLogin:     "LastLogin",
	LastHost:      "LastHost",
}
