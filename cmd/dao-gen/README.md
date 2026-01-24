# DAO Code Generator

A code generation tool for creating Data Access Object (DAO) packages with the frantic-amphora framework.

## Overview

The `dao-gen` tool generates complete DAO packages with:

- Type-safe database operations
- Audit trail integration
- Cache management
- Background workers
- Import/Export capabilities
- Field validation

## Usage

### Basic Command

```bash
./dao-gen -pkg <package-name> -type <TypeName> -out <output-directory>
```

### Command Line Options

- `-pkg` (required) - Package name for the generated code
- `-type` (required) - Entity type name (e.g., "User", "Product")
- `-table` (optional) - Database table name (defaults to type name)
- `-namespace` (optional) - Cache namespace (defaults to "cheeseOnToast")
- `-out` (optional) - Output directory (defaults to current directory)
- `-force` - Overwrite existing files without rotating to .OLD
- `-with-worker` - Generate worker file (default: true)
- `-with-impex` - Generate import/export file (default: true)
- `-with-debug` - Generate debug file (default: true)

### Example

```bash
./dao-gen -pkg user -type User -table users -out ../../dao/user
```

## Using .definition Files

The generator supports `.definition` files to customize the domain fields of your entity.

### Creating a .definition File

1. Create a file named `<TypeName>.definition` in your output directory
2. Define your fields after the "Domain specific fields, starts" marker
3. Use standard Go struct syntax

#### Example: User.definition

```go
// User.definition

// Domain specific fields, starts.
//
// User account fields
UserName     string             `validate:"required,min=3"`
Email        string             `storm:"index,unique" validate:"required,email"`
PasswordHash string             `validate:"required"`
FirstName    string             `validate:"required"`
LastName     string             `validate:"required"`
IsActive     entities.Bool
IsAdmin      entities.Bool
CreatedAt    time.Time
UpdatedAt    time.Time
LastLoginAt  time.Time
LoginCount   entities.Int
```

### Field Definition Rules

1. **Mandatory Fields** (automatically included):
   - `ID int` - Primary key
   - `Key string` - Indexed unique identifier
   - `Raw string` - Raw unique identifier
   - `Audit audit.Audit` - Audit trail

2. **Your Domain Fields**:
   - Define only your business-specific fields in the .definition file
   - Use any Go type or framework entity type (entities.Bool, entities.Int, etc.)
   - Add struct tags for Storm indexing and validation
   - Comments and blank lines are preserved

3. **Automatic Generation**:
   - Each field in your struct automatically gets:
     - A `fieldNames` struct entry
     - A `Fields` variable initialization
   - This enables type-safe queries: `GetBy(Fields.Email, "user@example.com")`

### Available Entity Types

The framework provides special entity types that can be marshalled to/from strings:

**Boolean Types:**

- `entities.Bool` - Framework boolean
- `entities.StormBool` - Storm-compatible boolean

**Integer Types:**

- `entities.Int`, `entities.Int32`, `entities.Int64`
- `entities.UInt`, `entities.UInt32`, `entities.UInt64`

**Float Types:**

- `entities.Float`, `entities.Float32`, `entities.Float64`

**Specialized Numeric Types:**

- `entities.Decimal` - High-precision decimal
- `entities.Percentage` - Percentage values
- `entities.Rate` - Rate values

**Money Types:**

- `entities.Money` - Monetary amount
- `entities.Currency` - Currency with code and amount

**Other Types:**

- `time.Time` - Date/time values
- `string` - Standard strings
- `entities.Field` - Field name type
- `entities.Table` - Table name type

## Generated Files

Running `dao-gen` creates the following files:

- `<type>Model.go` - Entity struct definition (uses .definition file)
- `<type>.go` - Main DAO operations (Count, Get, Create, Update, Delete)
- `<type>DB.go` - Database lifecycle (Initialise, Close, connections)
- `<type>Cache.go` - Cache integration (hydration, synchronization)
- `<type>Internals.go` - Internal helper functions
- `<type>Helpers.go` - Custom business logic hooks
- `<type>Worker.go` - Background job processing (optional)
- `<type>Impex.go` - Import/Export functionality (optional)
- `<type>Debug.go` - Debug utilities (optional)
- `README.md` - Package documentation

## File Rotation

Certain files that commonly contain custom code (`*Model.go`, `*Helpers.go`) will be rotated to `.OLD` instead of being overwritten during regeneration. This preserves your customizations.

To force overwriting without rotation, use the `-force` flag.

## Workflow

### Initial Generation

```bash
# 1. Create your .definition file
cat > User.definition << 'EOF'
// Domain specific fields, starts.
UserName string `validate:"required"`
Email    string `storm:"index,unique" validate:"email"`
IsActive entities.Bool
EOF

# 2. Generate the DAO package
./dao-gen -pkg user -type User -out ../../dao/user

# 3. Customize business logic in *Helpers.go as needed
```

### Regeneration After Changes

```bash
# 1. Update your .definition file with new fields
# 2. Regenerate (existing Model.go will be rotated to .OLD)
./dao-gen -pkg user -type User -out ../../dao/user

# 3. Merge any custom code from .OLD file if needed
```

## Template Customization

Templates are embedded in the binary from the `templates/` directory:

- `model.tmpl` - Entity model structure
- `dao.tmpl` - Main DAO operations
- `db.tmpl` - Database management
- `cache.tmpl` - Cache integration
- `internals.tmpl` - Internal functions
- `helpers.tmpl` - Business logic hooks
- `worker.tmpl` - Background worker
- `impex.tmpl` - Import/Export
- `debug.tmpl` - Debug utilities
- `readme.tmpl` - Package documentation

To customize templates, modify the `.tmpl` files and rebuild the generator.

## Example Complete Workflow

```bash
# Generate a Product DAO
cd cmd/dao-gen

# Build the generator
go build -o dao-gen

# Create the definition file
cat > Product.definition << 'EOF'
// Domain specific fields, starts.
SKU         string           `storm:"index,unique" validate:"required"`
Name        string           `validate:"required,min=3"`
Description string           `validate:"max=500"`
Price       entities.Decimal `validate:"required,gt=0"`
Currency    string           `validate:"required,len=3"`
StockLevel  entities.Int
IsAvailable entities.Bool
CategoryID  string           `storm:"index"`
CreatedAt   time.Time
UpdatedAt   time.Time
EOF

# Generate the DAO package
mkdir -p ../../dao/product
./dao-gen -pkg product -type Product -table products -out ../../dao/product

# The generated package is ready to use!
```

## Version

Current version: 0.5.0 (Updated: 2026-01-24)

## See Also

- [TemplateStore.definition](./TemplateStore.definition) - Example definition file with all entity types
- [../dao/test/templateStoreV2/](../dao/test/templateStoreV2/) - Reference implementation
