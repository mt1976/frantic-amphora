# Running Code-Gen (`dao-gen`)

This repo includes a small code generator at `cmd/dao-gen` that produces a new DAO package using the TemplateStoreV2 template set.

See also: [README.md](README.md)

## Quick start

From the repo root:

```bash
mkdir -p dao/test/fred
go run ./cmd/dao-gen -out dao/test/fred -pkg fred -type Fred -table Fred -namespace main -force
```

What those flags mean (high level):

- `-out`: where the generated package files are written
- `-pkg`: Go package name (usually lowercase)
- `-type`: exported Go struct/type name (usually PascalCase)
- `-table`: table name used by the DAO (defaults to `-type` if omitted)
- `-namespace`: value passed into `database.WithNameSpace(...)`
- `-force`: allow overwriting existing generated files (with some safety rules below)

## Usage

Run:

```bash
go run ./cmd/dao-gen [flags]
```

You can also see built-in help:

```bash
go run ./cmd/dao-gen -h
```

## Options

All flags supported by `dao-gen`:

- `-out` (string, default `.`)
  - Output directory for generated files.
- `-pkg` (string, required)
  - Go package name to generate (example: `fred`).
- `-type` (string, required)
  - Exported type/struct name (example: `Fred`).
- `-table` (string, optional)
  - Logical table name used to create `TableName`.
  - Defaults to `-type`.
- `-namespace` (string, optional)
  - Value passed into `database.WithNameSpace(...)` when connecting.
  - Defaults to `cheeseOnToast` if omitted.
- `-force` (bool, default `false`)
  - Allows overwriting existing generated files.
  - Without `-force`, the generator will refuse to overwrite any existing files.
- `-with-worker` (bool, default `true`)
  - Generate the worker file (`*Worker.go`).
- `-with-impex` (bool, default `true`)
  - Generate the import/export file (`*Impex.go`).
- `-with-debug` (bool, default `true`)
  - Generate the debug file (`*Debug.go`).

## Custom field definitions

You can customize the generated model by creating a `.definition` file in the output directory with the same name as your type. For example, for `Fred`, create `Fred.definition`:

```go
// Domain specific fields, starts.

// User's display name
Name string `storm:"index"`

// User's email address
Email string `storm:"index,unique"`

// Whether the user is active
IsActive entities.Bool

// User's age
Age entities.Int

// Account balance
Balance entities.Money
```

The generator will:

1. Inject these fields into the generated `*Model.go` struct
2. Generate corresponding `Fields` entries for type-safe queries
3. Create a field definitions table in the generated `README.md`

All lines before "Domain specific fields, starts." are ignored. Comments immediately before a field are used as documentation in the README.

## Using with `go generate`

The recommended workflow is to commit a small `generate.go` file inside your target package:

```go
package fred

//go:generate go run ../../../cmd/dao-gen -out . -pkg fred -type Fred -table Fred -namespace main -force
```

Then you can regenerate with:

```bash
go generate ./...
```

## Using `dao-gen` from another project (Option 2)

If you want to use the generator from a different repo, you can run it directly from this module without copying any files.

Important: the generated code imports `database` and `cache` from frantic-amphora:

- `github.com/mt1976/frantic-amphora/dao/database`
- `github.com/mt1976/frantic-amphora/dao/cache`

So the project you generate into must also depend on frantic-amphora (it will compile against those packages).

### One-off usage

From your other project:

```bash
go get github.com/mt1976/frantic-amphora@latest

go run github.com/mt1976/frantic-amphora/cmd/dao-gen@latest -out ./app/dao/destinationStore -pkg destinationStore -type DestinationStore -table DestinationStore -namespace main -force

go run github.com/mt1976/frantic-amphora/cmd/dao-gen@latest -out ./app/dao/itemStore -pkg itemStore -type ItemStore -table ItemStore -namespace main -force

go run github.com/mt1976/frantic-amphora/cmd/dao-gen@latest -out ./app/dao/lockStore -pkg lockStore -type LockStore -table LockStore -namespace main -force

go run github.com/mt1976/frantic-amphora/cmd/dao-gen@latest -out ./app/dao/profileStore -pkg profileStore -type ProfileStore -table ProfileStore -namespace main -force

go run github.com/mt1976/frantic-amphora/cmd/dao-gen@latest -out ./app/dao/sequenceStore -pkg sequenceStore -type SequenceStore -table SequenceStore -namespace main -force

go run github.com/mt1976/frantic-amphora/cmd/dao-gen@latest -out ./app/dao/settingsStore -pkg settingsStore -type SettingsStore -table SettingsStore -namespace main -force

go run github.com/mt1976/frantic-amphora/cmd/dao-gen@latest -out ./app/dao/statusStore -pkg statusStore -type StatusStore -table StatusStore -namespace main -force

go run github.com/mt1976/frantic-amphora/cmd/dao-gen@latest -out ./app/dao/templateStore -pkg templateStore -type TemplateStore -table TemplateStore -namespace main -force

go run github.com/mt1976/frantic-amphora/cmd/dao-gen@latest -out ./app/dao/tripStore -pkg tripStore -type TripStore -table TripStore -namespace main -force

go run github.com/mt1976/frantic-amphora/cmd/dao-gen@latest -out ./app/dao/userStore -pkg userStore -type UserStore -table UserStore -namespace main -force

```

In the example above replace fred/Fred with the name of the DAO, for example userStore/UserStore and the namespace used by the db instance, normally this will be "main"

For repeatable builds, pin a specific version instead of `@latest`.

### With `go generate`

In your other project’s package, you can use a directive like:

```go
package fred

//go:generate go run github.com/mt1976/frantic-amphora/cmd/dao-gen@latest -out . -pkg fred -type Fred -table Fred -namespace main -force
```

If you prefer to pin the generator version, replace `@latest` with a tag/commit (and keep your project’s `go.mod` aligned).

## Output files

`dao-gen` writes a set of `.go` files similar to TemplateStoreV2, including:

- `*Model.go` (type + Fields + TableName, can be customized via `.definition` file)
- `*DB.go` (database lifecycle management)
- `*Cache.go` (cache hydrator and synchronizer)
- `*.go` (main DAO CRUD/query functions)
- `*Internals.go` (internal validation and processing)
- `*Helpers.go` (function registration system for custom hooks)
- Optional: `*Worker.go`, `*Impex.go`, `*Debug.go`
- `README.md` (package documentation with field definitions table)
