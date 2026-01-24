# TemplateStoreV2

`templateStoreV2` is a comprehensive example DAO package demonstrating the frantic-amphora framework's capabilities. It serves as a modernized template showcasing typed database operations, cache integration, background workers, import/export functionality, and comprehensive entity type examples.

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

The `TemplateStore` struct includes examples of all entity types available in the framework, organized by category:

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

## Public API

### Exported types/vars

- `type TemplateStore struct { ... }`
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
- `func GetBy(field entities.Field, value any) (TemplateStore, error)`
- `func GetAll() ([]TemplateStore, error)`
- `func GetAllWhere(field entities.Field, value any) ([]TemplateStore, error)`

### Mutations

- `func Delete(ctx context.Context, id int, note string) error`
- `func DeleteBy(ctx context.Context, field entities.Field, value any, note string) error`
- `func Drop() error`
- `func ClearDown(ctx context.Context) error`

### Record methods

- `func (record *TemplateStore) Validate() error`
- `func (record *TemplateStore) Update(ctx context.Context, note string) error`
- `func (record *TemplateStore) UpdateWithAction(ctx context.Context, auditAction audit.Action, note string) error`
- `func (record *TemplateStore) Create(ctx context.Context, note string) error`
- `func (record *TemplateStore) Clone(ctx context.Context) (TemplateStore, error)`

### Lookups

- `func GetDefaultLookup() (lookup.Lookup, error)`
- `func GetLookup(field, value entities.Field) (lookup.Lookup, error)`

### Cache integration

- `func CacheHydrator(ctx context.Context) func() ([]any, error)`
- `func CacheSynchroniser(ctx context.Context) func(any) error`

### Construction

- `func New() TemplateStore`
- `func Create(ctx context.Context, userName, uid, realName, email, gid string) (TemplateStore, error)`

### Import / Export

- `func (record *TemplateStore) ExportRecordAsJSON(name string)`
- `func ExportAllAsJSON(message string)`
- `func (record *TemplateStore) ExportRecordAsCSV(name string) error`
- `func ExportAllAsCSV(msg string) error`
- `func ImportAllFromCSV() error`

### Worker

- `func Worker(j jobs.Job, db *database.DB)`

### Example business logic

- `func Login(ctx context.Context, sq string) (TemplateStore, error)`
- `func Add(ctx context.Context, sq string) (TemplateStore, error)`
- `func (u *TemplateStore) SetName(name string) error`

### Debug

- `func (record *TemplateStore) Spew()`

