# templateStoreV3

`templateStoreV3` is a comprehensive example DAO package demonstrating the frantic-amphora framework's capabilities. It serves as a modernized template showcasing typed database operations, cache integration, background workers, import/export functionality, and comprehensive entity type examples.

## Overview

This package provides a complete reference implementation for building DAO layers with:

- **Type-safe database operations** using strongly-typed field queries
- **Audit trail integration** for all CRUD operations
- **Cache management** with automatic hydration and synchronization
- **Background worker** support for async operations
- **Import/Export** capabilities (JSON and CSV formats)
- **Validation** using struct tags
- **Comprehensive entity type examples** demonstrating all available framework types

## Entity Type Examples

The `TemplateStoreV3` struct includes examples of all entity types available in the framework, organized by category:

### String Types

- `ExampleString` - Standard Go string

### Boolean Types

- `ExampleBool` - Framework boolean type (entities.Bool)
- `ExampleStormBool` - Storm-compatible boolean type (entities.StormBool)

### Integer Types

- `ExampleInt` - Framework integer (entities.Int)
- `ExampleInt32` - 32-bit integer (entities.Int32)
- `ExampleInt64` - 64-bit integer (entities.Int64)
- `ExampleUint` - Unsigned integer (entities.UInt)
- `ExampleUint32` - Unsigned 32-bit integer (entities.UInt32)
- `ExampleUint64` - Unsigned 64-bit integer (entities.UInt64)

### Float Types

- `ExampleFloat` - Framework float (entities.Float)
- `ExampleFloat32` - 32-bit float (entities.Float32)
- `ExampleFloat64` - 64-bit float (entities.Float64)

### Specialized Numeric Types

- `ExampleDecimal` - High-precision decimal (entities.Decimal)
- `ExamplePercentage` - Percentage values (entities.Percentage)
- `ExampleRate` - Rate values (entities.Rate)

### Money Types

- `ExampleMoney` - Monetary amount (entities.Money)
- `ExampleCurrency` - Currency with code and amount (entities.Currency)

### Date/Time Types

- `ExampleDate` - Standard Go time.Time

### Entity Framework Types

- `ExampleField` - Field name type (entities.Field)
- `ExampleTable` - Table name type (entities.Table)

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
| UID | `Fields.UID` | `string` | validate:"required" | User Management fields |
| GID | `Fields.GID` | `string` | storm:"index" validate:"required" |  |
| RealName | `Fields.RealName` | `string` | validate:"required,min=5" |  |
| UserName | `Fields.UserName` | `string` | validate:"required,min=5" |  |
| UserCode | `Fields.UserCode` | `string` | storm:"index" validate:"required,min=5" |  |
| Email | `Fields.Email` | `string` |  |  |
| Notes | `Fields.Notes` | `string` | validate:"max=75" |  |
| Active | `Fields.Active` | `entities.Bool` |  |  |
| LastLogin | `Fields.LastLogin` | `time.Time // Last login time` |  |  |
| LastHost | `Fields.LastHost` | `string` | storm:"index" |  |


**Note:** Fields marked as **(required)** are mandatory framework fields and must not be modified or removed.

### Using Field References

Field references enable type-safe queries throughout the DAO:

```go
// Get a record by a specific field
user, err := GetBy(Fields.UserName, "john.doe")

// Query with WHERE conditions
activeUsers, err := GetAllWhere(Fields.Active, true)

// Count records matching criteria
count, err := CountWhere(Fields.GID, "admin-group")
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
- `func GetBy(field entities.Field, value any) (TemplateStoreV3, error)`
- `func GetAll() ([]TemplateStoreV3, error)`
- `func GetAllWhere(field entities.Field, value any) ([]TemplateStoreV3, error)`

### Mutations

- `func Delete(ctx context.Context, id int, note string) error`
- `func DeleteBy(ctx context.Context, field entities.Field, value any, note string) error`
- `func Drop() error`
- `func ClearDown(ctx context.Context) error`

### Record methods

- `func (record *TemplateStoreV3) Validate() error`
- `func (record *TemplateStoreV3) Update(ctx context.Context, note string) error`
- `func (record *TemplateStoreV3) UpdateWithAction(ctx context.Context, auditAction audit.Action, note string) error`
- `func (record *TemplateStoreV3) Create(ctx context.Context, note string) error`
- `func (record *TemplateStoreV3) Clone(ctx context.Context) (TemplateStoreV3, error)`

### Lookups

- `func GetDefaultLookup() (lookup.Lookup, error)`
- `func GetLookup(field, value entities.Field) (lookup.Lookup, error)`

### Cache integration

- `func CacheHydrator(ctx context.Context) func() ([]any, error)`
- `func CacheSynchroniser(ctx context.Context) func(any) error`

### Construction

- `func New() TemplateStoreV3`
- `func Create(ctx context.Context, basis TemplateStoreV3) (TemplateStoreV3, error)`

### Import / Export

- `func (record *TemplateStoreV3) ExportRecordAsJSON(name string)`
- `func ExportAllAsJSON(message string)`
- `func (record *TemplateStoreV3) ExportRecordAsCSV(name string) error`
- `func ExportAllAsCSV(msg string) error`
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

**Generated Date:** 27/01/2026 & 12:22  
**Generated By:** matttownsend (orion)
**Generated From Template Version:** 0.5.10 - 2026-01-26