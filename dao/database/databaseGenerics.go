package database

import (
	"fmt"
	"reflect"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/index"
	"github.com/asdine/storm/v3/q"
	"github.com/mt1976/frantic-amphora/dao/cache"
	"github.com/mt1976/frantic-amphora/dao/entities"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/logHandler"
)

// GetTyped retrieves a single record from the database and returns it as a concrete type.
//
// NOTE: Go does not currently allow methods with type parameters, so these helpers are
// package-level functions.
//
// T is expected to be a struct type (not a pointer). Storm expects a pointer to a struct
// destination; if you call this with a pointer type, it will return an invalid-type error.
func GetTyped[T any](db *DB, field entities.Field, value any) (*T, error) {
	logHandler.Database.Printf("GetTyped called for %v where %v=%T(%v)", entities.GetStructType(new(T)), field.String(), value, value)
	var empty *T
	record := new(T)

	// Check if a record exists in the cache
	if cache.IsEnabled(&record) {
		cachedRecord, err := cache.GetWhere(record, field, value)
		if err == nil {
			logHandler.Database.Printf("[GET] %v WHERE %+v=%+v [...%v.db] - From Cache", entities.GetStructType(&record), field.String(), value, db.Name)
			return *cachedRecord, nil
		}
		logHandler.Database.Printf("[GET] %v WHERE %+v=%+v [...%v.db] - Not Found in Cache", entities.GetStructType(&record), field.String(), value, db.Name)
	}

	logHandler.Trace.Printf("[GET] %v WHERE %+v=%+v [...%v.db] %v", entities.GetStructType(&record), field.String(), value, db.Name, reflect.TypeOf(record))
	if err := db.connection.One(field.String(), value, record); err != nil {
		logHandler.Trace.Printf("[GET]Error in GetTyped for %v where %v=%v: %v", entities.GetStructType(&record), field.String(), value, err)
		return empty, err
	}
	return record, nil
}

// GetAllTyped retrieves all records for type T and returns a typed slice.
//
// NOTE: T is expected to be a struct type (not a pointer).
func GetAllTyped[T any](db *DB, options ...func(*index.Options)) ([]*T, error) {
	logHandler.Database.Printf("GetAllTyped called for %v", entities.GetStructType(new(T)))

	var record *T

	//if reflect.TypeOf(record) != nil && reflect.TypeOf(record).Kind() == reflect.Ptr {
	//	logHandler.ErrorLogger.Printf("Invalid type passed to GetAllTyped: %T", record)
	//	return nil, commonErrors.ErrInvalidTypeWrapper("GetAllTyped", fmt.Sprintf("%T", record), "non-pointer struct")
	//}

	// Check if cache is enabled and retrieve from cache if available
	if cache.IsEnabled(record) {
		logHandler.Event.Printf("Cache is enabled for %v", entities.GetStructType(record))
		cachedResult, err := cache.GetAll(record)
		if err == nil {
			logHandler.Database.Printf("[GET] %v ALL [...%v.db] - From Cache", entities.GetStructType(record), db.Name)
			// Dereference double pointers to match return type []*T
			result := make([]*T, len(cachedResult))
			for i, item := range cachedResult {
				result[i] = *item
			}
			return result, nil
		}
		logHandler.Database.Printf("[GET] %v ALL [...%v.db] - Not Found in Cache", entities.GetStructType(record), db.Name)
	}

	logHandler.Trace.Printf("[GET] %v ALL [...%v.db]", entities.GetStructType(record), db.Name)

	result := []*T{}
	if err := db.connection.All(&result, options...); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAllWhereTyped retrieves all matching records for type T filtered by field/value.
//
// NOTE: T is expected to be a struct type (not a pointer).
func GetAllWhereTyped[T any](db *DB, field entities.Field, value any) ([]*T, error) {
	logHandler.Database.Printf("GetAllWhereTyped called for %v where %v=%T(%v)", entities.GetStructType(new(T)), field.String(), value, value)
	var record T
	var zero []*T
	if reflect.TypeOf(record) != nil && reflect.TypeOf(record).Kind() == reflect.Ptr {
		return nil, commonErrors.ErrInvalidTypeWrapper("GetAllWhereTyped", fmt.Sprintf("%T", record), "non-pointer struct")
	}

	if err := entities.IsValidFieldInStruct(field, record); err != nil {
		return nil, err
	}
	if err := entities.IsValidTypeForField(field, value, record); err != nil {
		return nil, err
	}
	logHandler.Database.Printf("Valid field/type check passed for %v.%v=%T(%v)", entities.GetStructType(record), field.String(), value, value)

	// Check if cache is enabled and retrieve from cache if available
	if cache.IsEnabled(record) {
		cachedResult, err := cache.GetAllWhere(record, field, value)
		if err == nil {
			logHandler.Database.Printf("[GET] %v WHERE (%+v=%+v) ALL [...%v.db] - From Cache", entities.GetStructType(record), field.String(), value, db.Name)
			return cachedResult, nil
		}
		logHandler.Database.Printf("[GET] %v WHERE (%+v=%+v) ALL [...%v.db] - Not Found in Cache", entities.GetStructType(record), field.String(), value, db.Name)
	}

	logHandler.Database.Printf("[GET] %v WHERE (%+v=%+v) ALL [...%v.db]", entities.GetStructType(record), field.String(), value, db.Name)
	result := []*T{}
	query := db.connection.Select(q.Eq(field.String(), value))
	if err := query.Find(&result); err != nil {
		if err == storm.ErrNotFound {
			logHandler.Database.Printf("No records found for %v where %v=%v", entities.GetStructType(record), field.String(), value)
			return zero, nil
		}
		logHandler.Error.Printf("Error in GetAllWhereTyped for %v where %v=%v: %v", entities.GetStructType(record), field.String(), value, err)
		return nil, err
	}
	logHandler.Database.Printf("Found %d records for %v where %v=%v", len(result), entities.GetStructType(record), field.String(), value)
	return result, nil
}

// GetAllPagedTyped retrieves a single page of records for type T, along with paging
// information (including the total number of pages) for display purposes.
//
// Pages are 1-based. Paged reads always go to the database (bypassing the cache) so
// that record ordering is stable between pages.
//
// NOTE: T is expected to be a struct type (not a pointer).
func GetAllPagedTyped[T any](db *DB, page, pageSize int) ([]*T, PageInfo, error) {
	logHandler.Database.Printf("GetAllPagedTyped called for %v page=%d pageSize=%d", entities.GetStructType(new(T)), page, pageSize)
	record := new(T)

	page, offset, err := validatePaging("GetAllPagedTyped", page, pageSize)
	if err != nil {
		return nil, PageInfo{}, err
	}

	total, err := db.connection.Count(record)
	if err != nil {
		logHandler.Error.Printf("Error counting %v for paging: %v", entities.GetStructType(record), err)
		return nil, PageInfo{}, err
	}
	pageInfo := NewPageInfo(page, pageSize, total)

	logHandler.Database.Printf("[GET] %v ALL PAGE %d/%d (size=%d, skip=%d) [...%v.db] - From Database", entities.GetStructType(record), pageInfo.Page, pageInfo.TotalPages, pageSize, offset, db.Name)
	result := []*T{}
	query := db.connection.Select(q.True()).Skip(offset).Limit(pageSize)
	if err := query.Find(&result); err != nil {
		if err == storm.ErrNotFound {
			logHandler.Database.Printf("No records found for %v on page %d/%d", entities.GetStructType(record), pageInfo.Page, pageInfo.TotalPages)
			return []*T{}, pageInfo, nil
		}
		logHandler.Error.Printf("Error in GetAllPagedTyped for %v page %d: %v", entities.GetStructType(record), page, err)
		return nil, PageInfo{}, err
	}
	logHandler.Database.Printf("Found %d records for %v on page %d/%d", len(result), entities.GetStructType(record), pageInfo.Page, pageInfo.TotalPages)
	return result, pageInfo, nil
}

// GetAllPagedWhereTyped retrieves a single page of records for type T filtered by
// field/value, along with paging information (including the total number of pages).
//
// Pages are 1-based. Paged reads always go to the database (bypassing the cache) so
// that record ordering is stable between pages.
//
// NOTE: T is expected to be a struct type (not a pointer).
func GetAllPagedWhereTyped[T any](db *DB, field entities.Field, value any, page, pageSize int) ([]*T, PageInfo, error) {
	logHandler.Database.Printf("GetAllPagedWhereTyped called for %v where %v=%T(%v) page=%d pageSize=%d", entities.GetStructType(new(T)), field.String(), value, value, page, pageSize)
	var record T
	if reflect.TypeOf(record) != nil && reflect.TypeOf(record).Kind() == reflect.Ptr {
		return nil, PageInfo{}, commonErrors.ErrInvalidTypeWrapper("GetAllPagedWhereTyped", fmt.Sprintf("%T", record), "non-pointer struct")
	}

	if err := entities.IsValidFieldInStruct(field, record); err != nil {
		return nil, PageInfo{}, err
	}
	if err := entities.IsValidTypeForField(field, value, record); err != nil {
		return nil, PageInfo{}, err
	}
	logHandler.Database.Printf("Valid field/type check passed for %v.%v=%T(%v)", entities.GetStructType(record), field.String(), value, value)

	page, offset, err := validatePaging("GetAllPagedWhereTyped", page, pageSize)
	if err != nil {
		return nil, PageInfo{}, err
	}

	total, err := db.connection.Select(q.Eq(field.String(), value)).Count(&record)
	if err != nil {
		logHandler.Error.Printf("Error counting %v where %v=%v for paging: %v", entities.GetStructType(record), field.String(), value, err)
		return nil, PageInfo{}, err
	}
	pageInfo := NewPageInfo(page, pageSize, total)

	logHandler.Database.Printf("[GET] %v WHERE (%+v=%+v) PAGE %d/%d (size=%d, skip=%d) [...%v.db] - From Database", entities.GetStructType(record), field.String(), value, pageInfo.Page, pageInfo.TotalPages, pageSize, offset, db.Name)
	result := []*T{}
	query := db.connection.Select(q.Eq(field.String(), value)).Skip(offset).Limit(pageSize)
	if err := query.Find(&result); err != nil {
		if err == storm.ErrNotFound {
			logHandler.Database.Printf("No records found for %v where %v=%v on page %d/%d", entities.GetStructType(record), field.String(), value, pageInfo.Page, pageInfo.TotalPages)
			return []*T{}, pageInfo, nil
		}
		logHandler.Error.Printf("Error in GetAllPagedWhereTyped for %v where %v=%v page %d: %v", entities.GetStructType(record), field.String(), value, page, err)
		return nil, PageInfo{}, err
	}
	logHandler.Database.Printf("Found %d records for %v where %v=%v on page %d/%d", len(result), entities.GetStructType(record), field.String(), value, pageInfo.Page, pageInfo.TotalPages)
	return result, pageInfo, nil
}
