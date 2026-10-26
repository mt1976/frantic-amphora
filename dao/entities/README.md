# Entities

The `entities` package provides small types for representing values in DAO models and
database operations. Several numeric and boolean types store their value as a string,
which helps with storage backends that do not handle Go's native numeric or boolean
types consistently.

## Numeric types

| Type | Use |
| --- | --- |
| `Int` | Integer value |
| `Int64`, `Int32` | Integer values named for their signed bit width |
| `UInt64`, `UInt32`, `UInt` | Unsigned integer values |
| `Float` | Floating-point value |
| `Float32`, `Float64` | Floating-point values named for their bit width |
| `Decimal` | Float-backed value with decimal-oriented conversion helpers |
| `Money`, `Percentage`, `Rate` | Float-backed values whose names describe their domain |
| `Currency` | Monetary value and a three-letter currency code, such as `USD` |

`Int` and `Float` keep their underlying value in a `Value` string. Their `Set` methods
accept native Go numbers, and methods such as `Int()`, `Float()`, `Get()`, and `String()`
provide access to the value. Empty numeric values convert to zero. Invalid stored numeric
strings cause the conversion methods to panic.

`Currency` contains a `Float` amount and a `CCY` currency code. An empty code defaults to
`GBP`; non-empty codes must be exactly three characters and are stored in uppercase.
Use `Set`, `SetValue`, and `SetCode` to populate it, or the `New`, `NewCurrency`, and
`NewAmount` methods to create values.

Although `Decimal` has decimal-oriented helpers, it is based on `Float` and stores a
floating-point value. Use an appropriate decimal arithmetic type separately if exact
decimal arithmetic is required.

## Date and time types

`Date`, `Time` and `DateTime` store their `Value` as a native Go `time.Time`, not a string.
This lets the database store them as time values directly. Values are normalized to UTC
when set through the type methods.

| Type | Use |
| --- | --- |
| `Date` | Calendar date; its value is midnight UTC |
| `Time` | Time of day; its value also keeps a reference date for correct timezone conversion |
| `DateTime` | Point in time |

Each type can be set and read in UTC, in the machine's local timezone, or in any IANA timezone
you name:

| Purpose | Set | Get |
| --- | --- | --- |
| UTC | `Set(t)`, `SetNow()` | `Get()` |
| Local timezone | `SetLocal(t)` | `GetLocal()` |
| Any timezone | `SetIn("America/New_York", t)` | `GetIn("America/New_York")` |

`SetIn` and `GetIn` return an error for an unknown timezone name. Empty values return the
zero `time.Time`. All three also have `String`, `IsZero`, `Equals`, `Before` and `After`.

```go
var when entities.DateTime
when.SetIn("America/New_York", time.Date(2026, 1, 15, 9, 30, 0, 0, time.UTC)) // wall clock: 09:30 in New York
fmt.Println(when.String())                 // 2026-01-15T14:30:00Z (UTC)
tokyo, _ := when.GetIn("Asia/Tokyo")       // 23:30 on 15 Jan in Tokyo
```

- `DateTime` and `Time`: `SetLocal` and `SetIn` interpret the wall-clock fields of the `time.Time`
  you pass (including its date, ignoring its own location) as being in that zone, then store the
  resulting instant in UTC.
- `Time` is intended for time-of-day values, but retains the input date as a reference so timezone
  offsets and daylight-saving rules can be applied correctly. Its `String` method shows only the
  UTC time of day; `Get` returns the full stored UTC `time.Time`. `Equals`, `Before` and `After`
  compare only its UTC time of day, not the reference date.
- `Date` is never shifted between zones: `SetIn` stores the calendar date the instant falls on
  in that zone, and `GetIn` returns midnight of the stored date in that zone, so a date set
  in a zone reads back unchanged in that zone.

## Boolean types

| Type | Use |
| --- | --- |
| `Bool` | String-backed boolean with helpers for native Go `bool` values |
| `StormBool` | Boolean type named for compatibility with Storm |

Use `Set`, `SetTrue`, `SetFalse`, or `Toggle` to change a `Bool`; use `Bool()` or `Get()`
to read it as a Go boolean. `SetFromString` treats `true`, `1`, `yes`, `y`, and `t`
(case-insensitive, with surrounding whitespace ignored) as true; other values are false.
`Bool` also provides `HtmlChecked`, `HtmlSelected`, `HtmlDisabled`, and `HtmlReadOnly`
for use in HTML attributes.

`StormBool` has methods to set, read, and stringify its value.

## Identifiers and metadata

| Type | Use |
| --- | --- |
| `Field` | Database field name, with a `String()` method |
| `Table` | Database table name |
| `KEY` | Primary key value, including a `Key()` helper that wraps it in braces |
| `URI` | URI string value |
| `ACTION` | Action string value, such as a create, update, or delete action |

These types are string-backed and provide `Set`, `Get`, and/or `String` methods for
convenient conversion, depending on the type. For example:

```go
var field entities.Field = "Email"
var table entities.Table = "users"

var key entities.KEY
key.Set("123")
fmt.Println(key.Key()) // {123}
```

## Example in a model

```go
type User struct {
    ID        entities.Int
    IsActive  entities.Bool
    Balance   entities.Currency
    Email     string
}

var id entities.Int
id.Set(42)

var active entities.Bool
active.Set(true)

var balance entities.Currency
balance.Set("USD", 19.95)
```

Use the wrapped types when their storage behavior or semantic names are useful to your
model; ordinary Go types such as `string` and `time.Time` remain suitable for other fields.

## Source layout

| File | Contents |
| --- | --- |
| `entitiesModel.go` | Type definitions and constants |
| `entitiesInt.go` | `Int` methods |
| `entitiesFloat.go` | `Float` methods |
| `entitiesBool.go` | `Bool` and `StormBool` methods |
| `entitiesCurrency.go` | `Currency` methods |
| `entitiesDateTime.go` | `Date`, `Time` and `DateTime` methods |
| `entitiesIdentifiers.go` | `Field`, `Table`, `KEY`, `URI` and `ACTION` methods |
| `entitiesValidators.go` | Reflection-based field and type validators |
