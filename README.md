# frantic-amphora

**frantic-amphora** is a comprehensive Data Access Object (DAO) framework for Go applications, providing type-safe database operations, code generation, and enterprise features built on top of Storm (BoltDB). It's designed to reduce boilerplate while maintaining flexibility through extensible hooks and validators.

## What is frantic-amphora?

frantic-amphora provides:

- **Type-safe database operations** with strongly-typed field queries and generic helpers
- **Automatic code generation** (`dao-gen`) to scaffold complete DAO packages
- **Built-in auditing** tracking all CRUD operations with user context
- **Cache integration** with automatic synchronization and hydration
- **Background job scheduling** using cron expressions
- **Import/Export capabilities** for JSON and CSV formats
- **Validation framework** with customizable hooks
- **Extensible entity types** including typed numbers, money, decimals, and more

## Core Components

### DAO Framework (`dao/`)

The foundation of the framework with several key packages:

- **[database](dao/database/README.md)** - Storm-backed database layer with typed generic helpers
- **[cache](dao/cache/)** - Cache management and synchronization
- **[entities](dao/entities/)** - Typed entity definitions (Bool, Int, Money, Decimal, etc.)
- **[audit](dao/audit/)** - Audit trail integration for tracking changes
- **[lookup](dao/lookup/)** - Lookup table support
- **[maintenance](dao/maintenance/)** - Database backup and pruning utilities

### Code Generation (`cmd/dao-gen`)

The **[dao-gen](cmd/dao-gen/README.md)** tool generates complete DAO packages with:

- Type-safe CRUD operations
- Field definitions for query building
- Cache hydration and synchronization
- Background worker support
- Import/Export functionality
- Validation hooks and custom logic registration

See **[CODE-GEN.md](CODE-GEN.md)** for complete usage instructions and command-line options.

### Supporting Packages

- **[importExportHelper](importExportHelper/README.md)** - CSV and JSON import/export utilities
- **[jobs](jobs/README.md)** - Background job scheduling with cron support
- **[data/config](data/config/README.md)** - Sample configuration files

## Quick Start

### Using the Code Generator

Generate a new DAO package:

```bash
go run ./cmd/dao-gen -pkg user -type User -table users -out ./dao/user -force
```

### Custom Field Definitions

Create a `User.definition` file to define domain-specific fields:

```go
// Domain specific fields, starts.

// User's full name
Name string `storm:"index"`

// User's email address
Email string `storm:"index,unique"`

// User's age
Age entities.Int

// Account balance
Balance entities.Money
```

See the **[dao-gen README](cmd/dao-gen/README.md)** for details on field types and customization.

## Example Implementation

The **[templateStoreV2](dao/test/templateStoreV2/README.md)** package provides a complete reference implementation demonstrating:

- All available entity types (Bool, Int variants, Float variants, Decimal, Money, Currency, etc.)
- Function registration pattern for custom logic
- Cache integration
- Worker implementation
- Import/Export usage
- Validation patterns

## Architecture Pattern

Generated DAO packages follow a consistent structure:

1. **Model** (`*Model.go`) - Entity definition with typed fields
2. **DAO** (`*.go`) - CRUD operations (Count, Get, Create, Update, Delete)
3. **Database** (`*DB.go`) - Database lifecycle management
4. **Cache** (`*Cache.go`) - Cache hydration and synchronization
5. **Helpers** (`*Helpers.go`) - Function registration for custom hooks
6. **Internals** (`*Internals.go`) - Internal validation and processing
7. **Worker** (`*Worker.go`) - Background job processing (optional)
8. **Import/Export** (`*Impex.go`) - Data import/export (optional)
9. **Debug** (`*Debug.go`) - Debug utilities (optional)

Custom business logic is implemented in separate `*Logic.go` files and registered via the Helpers functions.

## Key Features

### Type-Safe Queries

```go
// Query using typed fields
user, err := userStore.GetBy(ctx, userStore.Fields.Email, "user@example.com")
```

### Audit Trail

All operations automatically track:

- Who performed the action
- When it occurred
- What changed

### Extensibility via Registration

```go
// Register custom validation
userStore.RegisterValidator(func(ctx context.Context, user *User) error {
    if user.Age.Int() < 18 {
        return errors.New("must be 18 or older")
    }
    return nil
})
```

### Generic Helpers

```go
import "github.com/mt1976/frantic-amphora/dao/database"

// Type-safe retrieval
user, err := database.GetTyped[User](db, userStore.Fields.ID, 123)
users, err := database.GetAllTyped[User](db)
```

## Documentation

- **[Code Generation Guide](CODE-GEN.md)** - Complete guide to using `dao-gen`
- **[dao-gen Tool Reference](cmd/dao-gen/README.md)** - Tool-specific documentation
- **[Database Package](dao/database/README.md)** - Database layer and generic helpers
- **[TemplateStoreV2 Example](dao/test/templateStoreV2/README.md)** - Reference implementation
- **[Import/Export Helper](importExportHelper/README.md)** - Data import/export utilities
- **[Jobs Package](jobs/README.md)** - Background job scheduling

## Dependencies

Built on:

- [Storm v3](https://github.com/asdine/storm) - Embedded key/value database (BoltDB wrapper)
- [frantic-core](https://github.com/mt1976/frantic-core) - Core utilities (logging, timing, paths, etc.)
- [go-playground/validator](https://github.com/go-playground/validator) - Struct validation
- [robfig/cron](https://github.com/robfig/cron) - Cron job scheduling (via jobs package)

## License

See [LICENSE](LICENSE) for details.
