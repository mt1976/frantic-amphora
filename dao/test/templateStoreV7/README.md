# templateStoreV7

`templateStoreV7` is a DAO package for the TemplateStoreV7 table in the frantic-amphora framework.

## Overview

This package provides:

- **Type-safe database operations** using strongly-typed field queries
- **Audit trail integration** for all CRUD operations
- **Cache management** with automatic hydration and synchronization
- **Background worker** support for async operations
- **Import/Export** capabilities (JSON and CSV formats)
- **Validation** using struct tags

## Entity Definition

The `TemplateStoreV7` struct represents records in the TemplateStoreV7 table.

## Field Definitions

The `TemplateStoreV7` struct contains the following fields:

| Field Name | Field Reference | Type | Tags | Purpose |
|------------|----------------|------|------|---------|
| **ID** (required) | `Fields.ID` | `int` | `storm:"id,increment=100"` | Primary key with auto-increment |
| **Key** (required) | `Fields.Key` | `string` | `storm:"index,unique"` | Encoded unique identifier |
| **Raw** (required) | `Fields.Raw` | `string` | `storm:"index,unique"` | Raw unique identifier |
| **Audit** (required) | `Fields.Audit` | `audit.Audit` | `csv:"-"` | Audit trail information |
| **Lock** (required) | `Fields.Lock` | `sync.Mutex` | `csv:"-"` | Record locking for concurrent updates |
| EntityInt | `Fields.EntityInt` | `entities.Int` |  |  |
| EntityInt64 | `Fields.EntityInt64` | `entities.Int64` |  |  |
| EntityUInt64 | `Fields.EntityUInt64` | `entities.UInt64` |  |  |
| EntityUInt | `Fields.EntityUInt` | `entities.UInt` |  |  |
| EntityInt32 | `Fields.EntityInt32` | `entities.Int32` |  |  |
| EntityUInt32 | `Fields.EntityUInt32` | `entities.UInt32` |  |  |
| EntityFloat | `Fields.EntityFloat` | `entities.Float` |  |  |
| EntityFloat32 | `Fields.EntityFloat32` | `entities.Float32` |  |  |
| EntityFloat64 | `Fields.EntityFloat64` | `entities.Float64` |  |  |
| EntityDecimal | `Fields.EntityDecimal` | `entities.Decimal` |  |  |
| EntityMoney | `Fields.EntityMoney` | `entities.Money` |  |  |
| EntityPercentage | `Fields.EntityPercentage` | `entities.Percentage` |  |  |
| EntityRate | `Fields.EntityRate` | `entities.Rate` |  |  |
| EntityCurrency | `Fields.EntityCurrency` | `entities.Currency` |  |  |
| EntityBool | `Fields.EntityBool` | `entities.Bool` |  |  |
| EntityStormBool | `Fields.EntityStormBool` | `entities.StormBool` |  |  |
| EntityDate | `Fields.EntityDate` | `entities.Date` |  |  |
| EntityTime | `Fields.EntityTime` | `entities.Time` |  |  |
| EntityDateTime | `Fields.EntityDateTime` | `entities.DateTime` |  |  |
| EntityField | `Fields.EntityField` | `entities.Field` |  |  |
| EntityURI | `Fields.EntityURI` | `entities.URI` |  |  |
| EntityAction | `Fields.EntityAction` | `entities.ACTION` |  |  |
| EntityKey | `Fields.EntityKey` | `entities.KEY` |  |  |
| EntityTable | `Fields.EntityTable` | `entities.Table` |  |  |
| Name | `Fields.Name` | `string` | `validate:"required,max=25,min=5"` |  |
| Profile | `Fields.Profile` | `string` | `storm:"index"` |  |
| ProfileKey | `Fields.ProfileKey` | `string` | `storm:"index"` |  |
| ProfileEnrichment | `Fields.ProfileEnrichment` | `string // Used to Display Date // profile name of the trip` |  |  |
| StatusKey | `Fields.StatusKey` | `string` | `storm:"index"` |  |
| StatusEnrichment | `Fields.StatusEnrichment` | `string // Only used during the display of records` |  |  |
| Destination | `Fields.Destination` | `string` | `storm:"index" validate:"required,max=35,min=4"` |  |
| DestinationKey | `Fields.DestinationKey` | `string` | `storm:"index"` |  |
| DestinationEnrichment | `Fields.DestinationEnrichment` | `string // Only used during the display of records` |  |  |
| Notes | `Fields.Notes` | `string` | `validate:"max=75"` |  |
| StartDate | `Fields.StartDate` | `time.Time // Start Date of the Trip` |  |  |
| EndDate | `Fields.EndDate` | `time.Time // End Date of the Trip` |  |  |
| DisplayStartDate | `Fields.DisplayStartDate` | `string // External Start Date of the Trip` |  |  |
| DisplayEndDate | `Fields.DisplayEndDate` | `string // External End Date of the Trip` |  |  |
| DurationInDays | `Fields.DurationInDays` | `entities.Int // Duration of the Trip in Days` |  |  |
| MinimumNumberOfItems | `Fields.MinimumNumberOfItems` | `entities.Int // Minimum Number of Items` |  |  |
| Countdown | `Fields.Countdown` | `entities.Int // Days remaining until the trip starts` |  |  |
| TotalItems | `Fields.TotalItems` | `entities.Int // Total number of items for the trip` |  |  |
| TotalItemsPacked | `Fields.TotalItemsPacked` | `entities.Int // Total number of items packed for the trip` |  |  |
| Year | `Fields.Year` | `string` | `storm:"index"` |  |
| DisplayID | `Fields.DisplayID` | `string` |  |  |
| SequenceNumber | `Fields.SequenceNumber` | `entities.Int // Sequence Number of the Trip` |  |  |
| PackingProgress | `Fields.PackingProgress` | `entities.Int // Packing Progress percentage` |  |  |
| TripProgress | `Fields.TripProgress` | `entities.Int // Trip Progress percentage` |  |  |
| Completed | `Fields.Completed` | `entities.Bool // Is the trip completed?` |  |  |
| StartNotified | `Fields.StartNotified` | `entities.Bool // Has the start notification been sent?` |  |  |
| EndWarned | `Fields.EndWarned` | `entities.Bool // Has the end warning been sent?` |  |  |
| EndNotified | `Fields.EndNotified` | `entities.Bool // Has the end notification been sent?` |  |  |
| LastHost | `Fields.LastHost` | `string` | `storm:"index"` |  |
| PostTest | `Fields.PostTest` | `[]string // For testing post processing hooks` |  |  |
| NewTimeStart | `Fields.NewTimeStart` | `entities.Time` |  |  |
| NewDateTimeStart | `Fields.NewDateTimeStart` | `entities.DateTime` |  |  |
| NewDateStart | `Fields.NewDateStart` | `entities.Date` |  |  |

| EntityInt | `Fields.EntityInt` | `entities.Int` |  |  |
| EntityInt64 | `Fields.EntityInt64` | `entities.Int64` |  |  |
| EntityUInt64 | `Fields.EntityUInt64` | `entities.UInt64` |  |  |
| EntityUInt | `Fields.EntityUInt` | `entities.UInt` |  |  |
| EntityInt32 | `Fields.EntityInt32` | `entities.Int32` |  |  |
| EntityUInt32 | `Fields.EntityUInt32` | `entities.UInt32` |  |  |
| EntityFloat | `Fields.EntityFloat` | `entities.Float` |  |  |
| EntityFloat32 | `Fields.EntityFloat32` | `entities.Float32` |  |  |
| EntityFloat64 | `Fields.EntityFloat64` | `entities.Float64` |  |  |
| EntityDecimal | `Fields.EntityDecimal` | `entities.Decimal` |  |  |
| EntityMoney | `Fields.EntityMoney` | `entities.Money` |  |  |
| EntityPercentage | `Fields.EntityPercentage` | `entities.Percentage` |  |  |
| EntityRate | `Fields.EntityRate` | `entities.Rate` |  |  |
| EntityCurrency | `Fields.EntityCurrency` | `entities.Currency` |  |  |
| EntityBool | `Fields.EntityBool` | `entities.Bool` |  |  |
| EntityStormBool | `Fields.EntityStormBool` | `entities.StormBool` |  |  |
| EntityDate | `Fields.EntityDate` | `entities.Date` |  |  |
| EntityTime | `Fields.EntityTime` | `entities.Time` |  |  |
| EntityDateTime | `Fields.EntityDateTime` | `entities.DateTime` |  |  |
| EntityField | `Fields.EntityField` | `entities.Field` |  |  |
| EntityURI | `Fields.EntityURI` | `entities.URI` |  |  |
| EntityAction | `Fields.EntityAction` | `entities.ACTION` |  |  |
| EntityKey | `Fields.EntityKey` | `entities.KEY` |  |  |
| EntityTable | `Fields.EntityTable` | `entities.Table` |  |  |
| Name | `Fields.Name` | `string` | `validate:"required,max=25,min=5"` |  |
| Profile | `Fields.Profile` | `string` | `storm:"index"` |  |
| ProfileKey | `Fields.ProfileKey` | `string` | `storm:"index"` |  |
| ProfileEnrichment | `Fields.ProfileEnrichment` | `string // Used to Display Date // profile name of the trip` |  |  |
| StatusKey | `Fields.StatusKey` | `string` | `storm:"index"` |  |
| StatusEnrichment | `Fields.StatusEnrichment` | `string // Only used during the display of records` |  |  |
| Destination | `Fields.Destination` | `string` | `storm:"index" validate:"required,max=35,min=4"` |  |
| DestinationKey | `Fields.DestinationKey` | `string` | `storm:"index"` |  |
| DestinationEnrichment | `Fields.DestinationEnrichment` | `string // Only used during the display of records` |  |  |
| Notes | `Fields.Notes` | `string` | `validate:"max=75"` |  |
| StartDate | `Fields.StartDate` | `time.Time // Start Date of the Trip` |  |  |
| EndDate | `Fields.EndDate` | `time.Time // End Date of the Trip` |  |  |
| DisplayStartDate | `Fields.DisplayStartDate` | `string // External Start Date of the Trip` |  |  |
| DisplayEndDate | `Fields.DisplayEndDate` | `string // External End Date of the Trip` |  |  |
| DurationInDays | `Fields.DurationInDays` | `entities.Int // Duration of the Trip in Days` |  |  |
| MinimumNumberOfItems | `Fields.MinimumNumberOfItems` | `entities.Int // Minimum Number of Items` |  |  |
| Countdown | `Fields.Countdown` | `entities.Int // Days remaining until the trip starts` |  |  |
| TotalItems | `Fields.TotalItems` | `entities.Int // Total number of items for the trip` |  |  |
| TotalItemsPacked | `Fields.TotalItemsPacked` | `entities.Int // Total number of items packed for the trip` |  |  |
| Year | `Fields.Year` | `string` | `storm:"index"` |  |
| DisplayID | `Fields.DisplayID` | `string` |  |  |
| SequenceNumber | `Fields.SequenceNumber` | `entities.Int // Sequence Number of the Trip` |  |  |
| PackingProgress | `Fields.PackingProgress` | `entities.Int // Packing Progress percentage` |  |  |
| TripProgress | `Fields.TripProgress` | `entities.Int // Trip Progress percentage` |  |  |
| Completed | `Fields.Completed` | `entities.Bool // Is the trip completed?` |  |  |
| StartNotified | `Fields.StartNotified` | `entities.Bool // Has the start notification been sent?` |  |  |
| EndWarned | `Fields.EndWarned` | `entities.Bool // Has the end warning been sent?` |  |  |
| EndNotified | `Fields.EndNotified` | `entities.Bool // Has the end notification been sent?` |  |  |
| LastHost | `Fields.LastHost` | `string` | `storm:"index"` |  |
| PostTest | `Fields.PostTest` | `[]string // For testing post processing hooks` |  |  |
| NewTimeStart | `Fields.NewTimeStart` | `entities.Time` |  |  |
| NewDateTimeStart | `Fields.NewDateTimeStart` | `entities.DateTime` |  |  |
| NewDateStart | `Fields.NewDateStart` | `entities.Date` |  |  |


**Note:** Fields marked as **(required)** are mandatory framework fields and must not be modified or removed.

### Using Field References

Field references enable type-safe queries throughout the DAO:

```go
// Get a record by a specific field
record, err := templateStoreV7.GetBy(templateStoreV7.Fields.Key, "abc123")

// Query with WHERE conditions
records, err := templateStoreV7.GetAllWhere(templateStoreV7.Fields.SomeField, value)

// Count records matching criteria
count, err := templateStoreV7.CountWhere(templateStoreV7.Fields.Active, true)

// Fetch a single page (1-based) - pageInfo.String() gives "Page 1 of 10"
records, pageInfo, err := templateStoreV7.GetAllPaged(1, 25)
records, pageInfo, err = templateStoreV7.GetAllPagedWhere(templateStoreV7.Fields.SomeField, value, 2, 25)
```

## Public API

### Exported types/vars

- `type TemplateStoreV7 struct { ... }`
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
- `func GetBy(field entities.Field, value any) (*TemplateStoreV7, error)`
- `func GetAll() ([]TemplateStoreV7, error)`
- `func GetAllWhere(field entities.Field, value any) ([]TemplateStoreV7, error)`
- `func GetAllPaged(page, pageSize int) ([]*TemplateStoreV7, database.PageInfo, error)`
- `func GetAllPagedWhere(field entities.Field, value any, page, pageSize int) ([]*TemplateStoreV7, database.PageInfo, error)`
- `func PageCount(pageSize int) (int, error)`
- `func PageCountWhere(field entities.Field, value any, pageSize int) (int, error)`

### Mutations

- `func Delete(ctx context.Context, id int, note string) error`
- `func DeleteBy(ctx context.Context, field entities.Field, value any, note string) error`
- `func Drop() error`
- `func ClearDown(ctx context.Context) error`

### Record methods

- `func (record *TemplateStoreV7) Validate() error`
- `func (record *TemplateStoreV7) Update(ctx context.Context, note string) (*TemplateStoreV7, error)`
- `func (record *TemplateStoreV7) UpdateWithAction(ctx context.Context, auditAction audit.Action, note string) (*TemplateStoreV7, error)`
- `func (record *TemplateStoreV7) Clone(ctx context.Context) (*TemplateStoreV7, error)`

### Lookups

- `func GetDefaultLookup() (lookup.Lookup, error)`
- `func GetLookup(field, value entities.Field) (lookup.Lookup, error)`

### Cache integration

- `func CacheHydrator(ctx context.Context) func() ([]any, error)`
- `func CacheSynchroniser(ctx context.Context) func(any) error`

### Construction

- `func New() *TemplateStoreV7`
- `func Create(ctx context.Context, basis *TemplateStoreV7) (*TemplateStoreV7, error)`

### Import / Export

- `func (record *TemplateStoreV7) ExportRecordToJSON(name string)`
- `func ExportAllToJSON(message string)`
- `func (record *TemplateStoreV7) ExportRecordToCSV(name string) error`
- `func ExportAllToCSV(msg string) error`
- `func ImportAllFromCSV() error`

### Worker

- `func Worker(j jobs.Job, db *database.DB)`

### Debug

- `func (record *TemplateStoreV7) Spew()`

## Regenerate

- From this package directory, run: `go generate ./...`

## Next edits

- Adjust the domain fields in the model file.
- Update validation/defaulting hooks.
- Replace any placeholder logic (e.g. clone, import processor) with real implementations.

---

## Generation Information

**Generated Date:** 30/10/2026 & 08:43  
**Generated By:** matttownsend (orion)  
**Generated From Template Version:** 0.5.24 - 2026-01-31
