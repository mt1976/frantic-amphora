# testentity

`testentity` is a comprehensive example DAO package demonstrating the frantic-amphora framework's capabilities. It serves as a modernized template showcasing typed database operations, cache integration, background workers, import/export functionality, and comprehensive entity type examples.

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

The `TestEntity` struct includes examples of all entity types available in the framework, organized by category:

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

The `TestEntity` struct contains the following fields:

| Field Name | Field Reference | Type | Tags | Purpose |
|------------|----------------|------|------|---------|
| **ID** (required) | `Fields.ID` | `int` | `storm:"id,increment=100"` | Primary key with auto-increment |
| **Key** (required) | `Fields.Key` | `string` | `storm:"index,unique"` | Encoded unique identifier |
| **Raw** (required) | `Fields.Raw` | `string` | `storm:"index,unique"` | Raw unique identifier |
| **Audit** (required) | `Fields.Audit` | `audit.Audit` | `csv:"-"` | Audit trail information |
| ExampleString | `Fields.ExampleString` | `string` |  | Example fields grouped by type String types |
| ExampleBool | `Fields.ExampleBool` | `entities.Bool` |  | Boolean types |
| ExampleStormBool | `Fields.ExampleStormBool` | `entities.StormBool` |  |  |
| ExampleInt | `Fields.ExampleInt` | `entities.Int` |  | Integer types |
| ExampleInt32 | `Fields.ExampleInt32` | `entities.Int32` |  |  |
| ExampleInt64 | `Fields.ExampleInt64` | `entities.Int64` |  |  |
| ExampleUint | `Fields.ExampleUint` | `entities.UInt` |  |  |
| ExampleUint32 | `Fields.ExampleUint32` | `entities.UInt32` |  |  |
| ExampleUint64 | `Fields.ExampleUint64` | `entities.UInt64` |  |  |
| ExampleFloat | `Fields.ExampleFloat` | `entities.Float` |  | Float types |
| ExampleFloat32 | `Fields.ExampleFloat32` | `entities.Float32` |  |  |
| ExampleFloat64 | `Fields.ExampleFloat64` | `entities.Float64` |  |  |
| ExampleDecimal | `Fields.ExampleDecimal` | `entities.Decimal` |  | Specialized numeric types |
| ExamplePercentage | `Fields.ExamplePercentage` | `entities.Percentage` |  |  |
| ExampleRate | `Fields.ExampleRate` | `entities.Rate` |  |  |
| ExampleMoney | `Fields.ExampleMoney` | `entities.Money` |  | Money types |
| ExampleCurrency | `Fields.ExampleCurrency` | `entities.Currency` |  |  |
| ExampleDate | `Fields.ExampleDate` | `time.Time` |  | Date/Time types |
| ExampleField | `Fields.ExampleField` | `entities.Field` |  | Entity framework types |
| ExampleTable | `Fields.ExampleTable` | `entities.Table` |  |  |
| LastLogin | `Fields.LastLogin` | `time.Time` |  | Additional fields for demonstration |
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

- `type TestEntity struct { ... }`
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
- `func GetBy(field entities.Field, value any) (TestEntity, error)`
- `func GetAll() ([]TestEntity, error)`
- `func GetAllWhere(field entities.Field, value any) ([]TestEntity, error)`

### Mutations

- `func Delete(ctx context.Context, id int, note string) error`
- `func DeleteBy(ctx context.Context, field entities.Field, value any, note string) error`
- `func Drop() error`
- `func ClearDown(ctx context.Context) error`

### Record methods

- `func (record *TestEntity) Validate() error`
- `func (record *TestEntity) Update(ctx context.Context, note string) error`
- `func (record *TestEntity) UpdateWithAction(ctx context.Context, auditAction audit.Action, note string) error`
- `func (record *TestEntity) Create(ctx context.Context, note string) error`
- `func (record *TestEntity) Clone(ctx context.Context) (TestEntity, error)`

### Lookups

- `func GetDefaultLookup() (lookup.Lookup, error)`
- `func GetLookup(field, value entities.Field) (lookup.Lookup, error)`

### Cache integration

- `func CacheHydrator(ctx context.Context) func() ([]any, error)`
- `func CacheSynchroniser(ctx context.Context) func(any) error`

### Construction

- `func New() TestEntity`
- `func Create(ctx context.Context, basis TestEntity) (TestEntity, error)`

### Import / Export

- `func (record *TestEntity) ExportRecordAsJSON(name string)`
- `func ExportAllAsJSON(message string)`
- `func (record *TestEntity) ExportRecordAsCSV(name string) error`
- `func ExportAllAsCSV(msg string) error`
- `func ImportAllFromCSV() error`

### Worker

- `func Worker(j jobs.Job, db *database.DB)`

### Debug

- `func (record *TestEntity) Spew()`

## Regenerate

- From this package directory, run: `go generate ./...`

## Next edits

- Adjust the domain fields in the model file.
- Update validation/defaulting hooks.
- Replace any placeholder logic (e.g. clone, import processor) with real implementations.
