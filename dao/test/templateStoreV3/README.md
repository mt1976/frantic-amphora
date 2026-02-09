# templateStoreV3

`templateStoreV3` is a DAO package for the TemplateStoreV3 table in the frantic-amphora framework.

## Overview

This package provides:

- **Type-safe database operations** using strongly-typed field queries
- **Audit trail integration** for all CRUD operations
- **Cache management** with automatic hydration and synchronization
- **Background worker** support for async operations
- **Import/Export** capabilities (JSON and CSV formats)
- **Validation** using struct tags

## Entity Definition

The `TemplateStoreV3` struct represents records in the TemplateStoreV3 table.

## Field Definitions

The `TemplateStoreV3` struct contains the following fields:

| Field Name | Field Reference | Type | Tags | Purpose |
|------------|----------------|------|------|---------|
| **ID** (required) | `Fields.ID` | `int` | `storm:"id,increment=100"` | Primary key with auto-increment |
| **Key** (required) | `Fields.Key` | `string` | `storm:"index,unique"` | Encoded unique identifier |
| **Raw** (required) | `Fields.Raw` | `string` | `storm:"index,unique"` | Raw unique identifier |
| **Audit** (required) | `Fields.Audit` | `audit.Audit` | `csv:"-"` | Audit trail information |
| ExampleString | `Fields.ExampleString` | `string // Example string field` |  | Example fields grouped by type String types |
| ExampleBool | `Fields.ExampleBool` | `entities.Bool // Example boolean field` |  | Boolean types |
| ExampleStormBool | `Fields.ExampleStormBool` | `entities.StormBool // Example storm boolean field` |  |  |
| ExampleInt | `Fields.ExampleInt` | `entities.Int // Example integer field` |  | Integer types |
| ExampleInt32 | `Fields.ExampleInt32` | `entities.Int32 // Example int32 field` |  |  |
| ExampleInt64 | `Fields.ExampleInt64` | `entities.Int64 // Example int64 field` |  |  |
| ExampleUint | `Fields.ExampleUint` | `entities.UInt // Example unsigned integer field` |  |  |
| ExampleUint32 | `Fields.ExampleUint32` | `entities.UInt32 // Example unsigned int32 field` |  |  |
| ExampleUint64 | `Fields.ExampleUint64` | `entities.UInt64 // Example unsigned int64 field` |  |  |
| ExampleFloat | `Fields.ExampleFloat` | `entities.Float // Example float field` |  | Float types |
| ExampleFloat32 | `Fields.ExampleFloat32` | `entities.Float32 // Example float32 field` |  |  |
| ExampleFloat64 | `Fields.ExampleFloat64` | `entities.Float64 // Example float64 field` |  |  |
| ExampleDecimal | `Fields.ExampleDecimal` | `entities.Decimal // Example decimal field` |  | Specialized numeric types |
| ExamplePercentage | `Fields.ExamplePercentage` | `entities.Percentage // Example percentage field` |  |  |
| ExampleRate | `Fields.ExampleRate` | `entities.Rate // Example rate field` |  |  |
| ExampleMoney | `Fields.ExampleMoney` | `entities.Money // Example money field` |  | Money types |
| ExampleCurrency | `Fields.ExampleCurrency` | `entities.Currency // Example currency field` |  |  |
| ExampleDate | `Fields.ExampleDate` | `time.Time // Example date field` |  | Date/Time types |
| ExampleField | `Fields.ExampleField` | `entities.Field // Example field type1` |  | Entity framework types |
| ExampleTable | `Fields.ExampleTable` | `entities.Table // Example table type1` |  |  |
| UID | `Fields.UID` | `string` | `validate:"required"` | User Management fields |
| GID | `Fields.GID` | `string` | `storm:"index" validate:"required"` |  |
| RealName | `Fields.RealName` | `string` | `validate:"required,min=5"` |  |
| UserName | `Fields.UserName` | `string` | `validate:"required,min=5"` |  |
| UserCode | `Fields.UserCode` | `string` | `storm:"index" validate:"required,min=5"` |  |
| Email | `Fields.Email` | `string` |  |  |
| Notes | `Fields.Notes` | `string` | `validate:"max=75"` |  |
| Active | `Fields.Active` | `entities.Bool` |  |  |
| LastLogin | `Fields.LastLogin` | `time.Time // Last login time` |  |  |
| LastHost | `Fields.LastHost` | `string` | `storm:"index"` |  |
| PostTest | `Fields.PostTest` | `[]string // For testing post processing hooks` |  |  |
| Lock | `Fields.Lock` | `sync.Mutex // For testing concurrent updates, not persisted to DB` |  |  |


**Note:** Fields marked as **(required)** are mandatory framework fields and must not be modified or removed.

### Using Field References

Field references enable type-safe queries throughout the DAO:

```go
// Get a record by a specific field
record, err := templateStoreV3.GetBy(templateStoreV3.Fields.Key, "abc123")

// Query with WHERE conditions
records, err := templateStoreV3.GetAllWhere(templateStoreV3.Fields.SomeField, value)

// Count records matching criteria
count, err := templateStoreV3.CountWhere(templateStoreV3.Fields.Active, true)
```

## Public API

### Exported types/vars

- `type TemplateStoreV3 struct { ... }`
- `var TableName entities.Table`
- `var Fields fieldNames`

### Database lifecycle

- `func Initialise(ctx context.Context, cached bool)`
- `func IsInitialised() bool`
- `func Close()`
- `func GetDatabaseConnections() func() ([]*database.DB, error)`

### Queries

- `func Count() (int, error)`
- `func CountWhere(field entities.Field, value any) (int, error)`
- `func GetBy(field entities.Field, value any) (*TemplateStoreV3, error)`
- `func GetAll() ([]TemplateStoreV3, error)`
- `func GetAllWhere(field entities.Field, value any) ([]TemplateStoreV3, error)`

### Mutations

- `func Delete(ctx context.Context, id int, note string) error`
- `func DeleteBy(ctx context.Context, field entities.Field, value any, note string) error`
- `func Drop() error`
- `func ClearDown(ctx context.Context) error`

### Record methods

- `func (record *TemplateStoreV3) Validate() error`
- `func (record *TemplateStoreV3) Update(ctx context.Context, note string) (*TemplateStoreV3, error)`
- `func (record *TemplateStoreV3) UpdateWithAction(ctx context.Context, auditAction audit.Action, note string) (*TemplateStoreV3, error)`
- `func (record *TemplateStoreV3) Clone(ctx context.Context) (*TemplateStoreV3, error)`

### Lookups

- `func GetDefaultLookup() (lookup.Lookup, error)`
- `func GetLookup(field, value entities.Field) (lookup.Lookup, error)`

### Cache integration

- `func CacheHydrator(ctx context.Context) func() ([]any, error)`
- `func CacheSynchroniser(ctx context.Context) func(any) error`

### Construction

- `func New() *TemplateStoreV3`
- `func Create(ctx context.Context, basis *TemplateStoreV3) (*TemplateStoreV3, error)`

### Import / Export

- `func (record *TemplateStoreV3) ExportRecordToJSON(name string)`
- `func ExportAllToJSON(message string)`
- `func (record *TemplateStoreV3) ExportRecordToCSV(name string) error`
- `func ExportAllToCSV(msg string) error`
- `func ImportAllFromCSV() error`

### Worker

- `func Worker(j jobs.Job, db *database.DB)`

### Debug

- `func (record *TemplateStoreV3) Spew()`

## Regenerate

- From this package directory, run: `go generate ./...`

## Next edits

- Adjust the domain fields in the model file.
- Update validation/defaulting hooks.
- Replace any placeholder logic (e.g. clone, import processor) with real implementations.

---

## Generation Information

**Generated Date:** 09/02/2026 & 10:22  
**Generated By:** matttownsend (orion)  
**Generated From Template Version:** 0.5.24 - 2026-01-31
