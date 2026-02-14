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
	"github.com/mt1976/frantic-core/timing"

	"github.com/mt1976/frantic-core/logHandler"
)

// Retrieve retrieves a single record from the database based on the specified fields.Field and value.
//
// Parameters:
//   - fields.Field: The fields.Field to be used for filtering the record.
//   - value: The value of the specified fields.Field to filter the record.
//   - to: A pointer to the struct where the retrieved record will be stored.
//
// Returns:
//   - any: The retrieved record.
//   - error: An error object if any issues occur during the retrieval process; otherwise, nil.
//
// DEPRECATED: Use Get instead.
func (db *DB) Retrieve(field entities.Field, value, to any) (any, error) {
	logHandler.Warning.Printf("Retrieve is DEPRECATED, use Get instead")
	panic("Retrieve is DEPRECATED, use Get instead")
	// return db.get(field, value, to)
}

// Get retrieves a single record from the database based on the specified fields.Field and value.
//
// Parameters:
//   - fields.Field: The fields.Field to be used for filtering the record.
//   - value: The value of the specified fields.Field to filter the record.
//   - to: A pointer to the struct where the retrieved record will be stored.
//
// Returns:
//   - any: The retrieved record.
//   - error: An error object if any issues occur during the retrieval process; otherwise, nil.
func (db *DB) Get(field entities.Field, value, to *any) (*any, error) {
	// logHandler.DatabaseLogger.Printf("[GET]<%v> (%+v=%+v)[%+v] [...%v.db]", entities.GetStructType(to), fields.Field, value, entities.GetStructType(to), db.Name)
	return db.get(field, value, to)
}

// get is the internal implementation for retrieving a single record from the database.
//
// Parameters:
//   - fields.Field: The fields.Field to be used for filtering the record.
//   - value: The value of the specified fields.Field to filter the record.
//   - to: A pointer to the struct where the retrieved record will be stored.
//
// Returns:
//   - any: The retrieved record.
//   - error: An error object if any issues occur during the retrieval process; otherwise, nil.
func (db *DB) get(field entities.Field, value, to *any) (*any, error) {
	logHandler.Database.Printf("[GET] %v WHERE %+v=%+v) [...%v.db]", entities.GetStructType(to), field.String(), value, db.Name)

	if cache.IsEnabled(to) {
		cachedValue, err := cache.GetWhere(to, field, value)
		if err == nil {
			reflect.ValueOf(to).Elem().Set(reflect.ValueOf(cachedValue).Elem())
			logHandler.Database.Printf("[GET] %v WHERE %+v=%+v) [...%v.db] - From Cache", entities.GetStructType(to), field.String(), value, db.Name)
			return *cachedValue, nil
		}
		logHandler.Database.Printf("[GET] %v WHERE %+v=%+v) [...%v.db] - Not Found in Cache", entities.GetStructType(to), field.String(), value, db.Name)
	}

	logHandler.Database.Printf("[GET] %v WHERE %+v=%+v) [...%v.db] - From Database", entities.GetStructType(to), field.String(), value, db.Name)

	// [GET] from database
	err := db.connection.One(field.String(), value, to)

	logHandler.Database.Printf("[GET] %v WHERE %+v=%+v) [...%v.db] - Completed", entities.GetStructType(to), field.String(), value, db.Name)
	if err != nil {
		// On error, do not attempt to use or populate the cache
		logHandler.Error.Printf("[GET] %v WHERE %+v=%+v) [...%v.db] - Error from DB: %v", entities.GetStructType(to), field.String(), value, db.Name, err)
		return nil, err
	}

	if cache.IsEnabled(to) {
		logHandler.Database.Printf("[GET] %v WHERE %+v=%+v) [...%v.db] - Populating Cache", entities.GetStructType(to), field.String(), value, db.Name)
		cache.AddEntry(to)
	} else {
		logHandler.Database.Printf("[GET] %v WHERE %+v=%+v) [...%v.db] - Caching Disabled or Not Initialised", entities.GetStructType(to), field.String(), value, db.Name)
	}

	logHandler.Database.Printf("[GET] %v WHERE %+v=%+v) [...%v.db] - Returning", entities.GetStructType(to), field.String(), value, db.Name)

	return to, err
}

// GetAll retrieves all records of the specified type from the database.
//
// Parameters:
//   - to: A pointer to a slice where the retrieved records will be stored.
//
// Returns:
//   - []any: A slice of all retrieved records.
//   - error: An error object if any issues occur during the retrieval process; otherwise, nil.
func (db *DB) GetAll(to any, options ...func(*index.Options)) ([]*any, error) {
	logHandler.Info.Printf("[GET] %v ALL [%+v] [...%v.db]", entities.GetStructType(to), options, db.Name)

	if cache.IsEnabled(to) {
		var resultList []*any
		// Get all records from cache
		allRecords, err := cache.GetAll(to)
		if err != nil {
			logHandler.Error.Printf("[GET] %v ALL [%+v] [...%v.db] - Error retrieving from cache: %v", entities.GetStructType(to), options, db.Name, err)
			return nil, err
		}
		// Convert cached records to []any, but only return if there are records
		for _, record := range allRecords {
			resultList = append(resultList, record)
		}
		logHandler.Info.Printf("[GET] %v ALL [%+v] [...%v.db] - Returning %d cached entries", entities.GetStructType(to), options, db.Name, len(resultList))
		if len(resultList) > 0 {
			return resultList, nil
		}
	}

	logHandler.Info.Printf("[GET] %v ALL [%+v] [...%v.db] - From Database", entities.GetStructType(to), options, db.Name)
	// [GET] from database
	err := db.connection.All(to, options...)
	if err != nil {
		// On error, do not attempt to use or populate the cache
		logHandler.Error.Printf("[GET] %v ALL [%+v] [...%v.db] - Error from DB: %v", entities.GetStructType(to), options, db.Name, err)
		return nil, err
	}

	logHandler.Info.Printf("[GET] %v ALL [%+v] [...%v.db] - Completed", entities.GetStructType(to), options, db.Name)

	// Use reflection to iterate through the slice without assuming its concrete type
	sliceValue := reflect.ValueOf(to).Elem()
	if sliceValue.Kind() != reflect.Slice {
		logHandler.Error.Printf("[GET] %v ALL - Expected slice, got %v from DB", entities.GetStructType(to), sliceValue.Kind())
		return nil, commonErrors.ErrInvalidTypeWrapper("GetAll", string(entities.GetStructType(to)), sliceValue.Kind().String())
	}

	logHandler.Info.Printf("[GET] %v ALL [%+v] [...%v.db] - Retrieved %d entries from DB", entities.GetStructType(to), options, db.Name, sliceValue.Len())

	// Convert the typed slice (e.g. []TemplateStore) into []any
	result := make([]*any, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i).Interface()
		result[i] = &elem
	}

	// Populate cache if enabled
	if cache.IsEnabled(to) {
		logHandler.Info.Printf("[GET] %v ALL [%+v] [...%v.db] - Populating Cache", entities.GetStructType(to), options, db.Name)
		err = cache.AddEntries(result)
		if err != nil {
			logHandler.Error.Printf("[GET] %v ALL [%+v] [...%v.db] - Error populating Cache: %v", entities.GetStructType(to), options, db.Name, err)
			return nil, err
		}
	} else {
		logHandler.Info.Printf("[GET] %v ALL [%+v] [...%v.db] - Caching Disabled or Not Initialised", entities.GetStructType(to), options, db.Name)
	}

	logHandler.Info.Printf("[GET] %v ALL [%+v] [...%v.db] on %v - Returning %d entries from cache", entities.GetStructType(to), entities.GetStructType(to), db.Name, "GetAll", sliceValue.Len())
	return result, nil
}

// GetAllWhere retrieves all TemplateStore records that match the specified fields.Field and value.
//
// Parameters:
//   - fields.Field: The fields.Field to be used for filtering records.
//   - value: The value of the specified fields.Field to filter records.
//
// Returns:
//   - []TemplateStore: A slice of TemplateStore records that match the specified criteria.
//   - error: An error object if any issues occur during the retrieval process; otherwise, nil.
func (db *DB) GetAllWhere(field entities.Field, value, to *any) ([]*any, error) {
	tableName := entities.GetStructType(to)
	logHandler.Database.Printf("[GET] %v WHERE %v=%v", tableName, field.String(), value)
	zero := []*any{}
	clock := timing.Start(tableName.String(), "GetAll", fmt.Sprintf("%v=%v", field.String(), value))

	// logHandler.DatabaseLogger.Printf("SELECT %v WHERE %v=%v", Domain, fields.Field, value)
	if err := entities.IsValidFieldInStruct(field, to); err != nil {
		logHandler.Error.Printf("Field validation error for fields.Field '%v': %v", field.String(), err)
		clock.Stop(0)
		return nil, err
	}

	if err := entities.IsValidTypeForField(field, value, to); err != nil {
		logHandler.Error.Printf("Type validation error for fields.Field '%v': %v", field.String(), err)
		clock.Stop(0)
		return nil, err
	}

	// If caching is enabled, attempt to retrieve records from cache
	if cache.IsEnabled(to) {
		cachedValues, err := cache.GetAllWhere(to, field, value)
		if err == nil && len(cachedValues) > 0 {
			logHandler.Database.Printf("[GET] %v WHERE %v=%v - From Cache", tableName, field.String(), value)
			// Convert []**any to []*any
			result := make([]*any, len(cachedValues))
			for i, v := range cachedValues {
				result[i] = *v
			}
			clock.Stop(len(result))
			return result, nil
		}
		logHandler.Database.Printf("[GET] %v WHERE %v=%v - Not Found in Cache", tableName, field.String(), value)
	}

	// Otherwise, use Storm's indexed query to retrieve matching records directly.
	query := db.connection.Select(q.Eq(field.String(), value))
	err := query.Find(to)
	if err != nil {
		if err == storm.ErrNotFound {
			clock.Stop(0)
			return zero, nil
		}
		logHandler.Error.Printf("Error querying %v where %v=%v: %v", tableName, field.String(), value, err)
		clock.Stop(0)
		return nil, err
	}

	sliceValue := reflect.ValueOf(to).Elem()
	if sliceValue.Kind() != reflect.Slice {
		logHandler.Error.Printf("[GET]<%v>{WHERE} - Expected slice pointer, got %v", tableName, sliceValue.Kind())
		clock.Stop(0)
		return nil, commonErrors.ErrInvalidTypeWrapper("GetAllWhere", tableName.String(), sliceValue.Kind().String())
	}

	resultList := make([]*any, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i).Interface()
		resultList[i] = &elem
	}

	// Populate cache if enabled
	if cache.IsEnabled(to) {
		logHandler.Database.Printf("[GET] %v WHERE %v=%v - Populating Cache", tableName, field.String(), value)
		err = cache.AddEntries(resultList)
		if err != nil {
			logHandler.Error.Printf("[GET] %v WHERE %v=%v - Error populating Cache: %v", tableName, field.String(), value, err)
			return nil, err
		}
	} else {
		logHandler.Database.Printf("[GET] %v WHERE %v=%v - Caching Disabled or Not Initialised", tableName, field.String(), value)
	}

	// hydrateerr := db.hydrateCacheBulk(resultList)
	// if hydrateerr != nil {
	// 	logHandler.ErrorLogger.Printf("[CCH]<%v>{AddBulk} Error hydrating cache: %v", tableName, hydrateerr)
	// }

	clock.Stop(len(resultList))
	return resultList, nil
}

// Delete removes the specified record from the database.
//
// Parameters:
//   - data: A pointer to the struct representing the record to be deleted.
//
// Returns:
//   - error: An error object if any issues occur during the deletion process; otherwise, nil.
func (db *DB) Delete(data any) error {
	logHandler.Database.Printf("[DELETE] %v [...%v.db] (%.10s)", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
	err := db.connection.DeleteStruct(data)
	if err != nil {
		logHandler.Error.Printf("[DELETE] %v [...%v.db] (%.10s) - Error: %v", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data), err)
		return err
	}

	if cache.IsEnabled(data) {
		logHandler.Database.Printf("[DELETE] %v [...%v.db] (%.10s) - Removing from Cache", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
		err := cache.RemoveEntry(data)
		if err != nil {
			logHandler.Error.Printf("[DELETE] %v [...%v.db] (%.10s) - Error removing from Cache: %v", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data), err)
			return err
		}
	} else {
		logHandler.Database.Printf("[DELETE] %v [...%v.db] (%.10s) - Caching Disabled or Not Initialised", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
	}
	// removeFromCache(db, data, "Delete", entities.GetStructType(data))
	return nil
}

// Drop removes the entire bucket or collection associated with the specified struct from the database.
//
// Parameters:
//   - data: A pointer to the struct representing the type whose bucket or collection is to be dropped.
//
// Returns:
//   - error: An error object if any issues occur during the drop process; otherwise, nil.
func (db *DB) Drop(data any) error {
	logHandler.Database.Printf("[DROP] %v [...%v.db] (%.10s)", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
	err := db.connection.Drop(data)
	if err != nil {
		logHandler.Error.Printf("[DROP] %v [...%v.db] (%.10s) - Error: %v", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data), err)
		return err
	}
	logHandler.Database.Printf("[DROP] %v [...%v.db] (%.10s) - Removing from Cache", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
	err = cache.ClearCacheForType(data)
	if err != nil {
		logHandler.Error.Printf("[DROP] %v [...%v.db] (%.10s) - Error clearing Cache: %v", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data), err)
		return err
	}
	// removeFromCache(db, data, "Drop", entities.GetStructType(data))
	return nil
}

// Update modifies an existing record in the database.
//
// Parameters:
//   - data: A pointer to the struct representing the record to be updated.
//
// Returns:
//   - error: An error object if any issues occur during the update process; otherwise, nil.
func (db *DB) Update(data any) error {
	logHandler.Database.Printf("[UPDATE] %v [...%v.db] (%.10s) - Start", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
	err := validate(data, db)
	if err != nil {
		logHandler.Error.Printf("[UPDATE] %v [...%v.db] (%.10s) - Error", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
		return commonErrors.ErrWrapper(err)
	}
	logHandler.Database.Printf("[UPDATE] %v [...%v.db] (%.10s) - End", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))

	if cache.IsEnabled(data) {
		logHandler.Info.Printf("[UPDATE] %v [...%v.db] (%.10s) - Updating Cache", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
		//		godump.Dump(data, "Updating Cache with data:")
		err := cache.AddEntry(data)
		if err != nil {
			logHandler.Error.Printf("[UPDATE] %v [...%v.db] (%.10s) - Error updating Cache: %v", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data), err)
			return err
		}
	}
	// Concurrently update the database
	// 	go bgUpdate(data, db)
	// } else {
	logHandler.Database.Printf("[UPDATE] %v [...%v.db] (%.10s) - Caching Disabled or Not Initialised", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
	err = db.connection.Update(data)
	if err != nil {
		logHandler.Error.Printf("[UPDATE] %v [...%v.db] (%.10s) - Error updating DB: %v", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data), err)
		return err
	}
	//	}

	return err
}

func bgUpdate(data any, db *DB) {
	logHandler.Database.Printf("[UPDATE] %v [...%v.db] (%.10s) - Caching Disabled or Not Initialised", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
	err := db.connection.Update(data)
	if err != nil {
		logHandler.Error.Panicf("[UPDATE] %v [...%v.db] (%.10s) - Error updating DB: %v", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data), err)
		return
	}
}

// Create adds a new record to the database.
//
// Parameters:
//   - data: A pointer to the struct representing the record to be created.
//
// Returns:
//   - error: An error object if any issues occur during the creation process; otherwise, nil.
func (db *DB) Create(data any) error {
	logHandler.Database.Printf("[CREATE] %v [...%v.db] (%.10s) - Start", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
	err := validate(data, db)
	if err != nil {
		logHandler.Error.Printf("[CREATE] %v [...%v.db] (%.10s) - Error", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
		return commonErrors.ErrCreateWrapper(err)
	}
	logHandler.Database.Printf("[CREATE] %v [...%v.db] (%.10s) - End", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))

	if cache.IsEnabled(data) {
		logHandler.Info.Printf("[CREATE] %v [...%v.db] (%.10s) - Adding to Cache", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
		err := cache.AddEntry(data)
		if err != nil {
			logHandler.Error.Printf("[CREATE] %v [...%v.db] (%.10s) - Error adding to Cache: %v", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data), err)
			return err
		}
		// Save to database asynchronously when cache is enabled
		go func() {
			saveErr := db.connection.Save(data)
			if saveErr != nil {
				logHandler.Error.Printf("[CREATE] %v [...%v.db] (%.10s) - Error: %v", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data), saveErr)
			}
		}()
	} else {
		logHandler.Database.Printf("[CREATE] %v [...%v.db] (%.10s) - Caching Disabled or Not Initialised", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data))
		err = db.connection.Save(data)
		if err != nil {
			logHandler.Error.Printf("[CREATE] %v [...%v.db] (%.10s) - Error: %v", entities.GetStructType(data), db.Name, fmt.Sprintf("%+v", data), err)
		}
	}
	return err
}

// Count returns the total number of records of the specified type in the database.
//
// Parameters:
//   - data: A pointer to the struct representing the type whose records are to be counted.
//
// Returns:
//   - int: The total number of records.
//   - error: An error object if any issues occur during the counting process; otherwise, nil.
func (db *DB) Count(data any) (int, error) {
	logHandler.Database.Printf("[COUNT] %v [...%v.db]", entities.GetStructType(data), db.Name)
	// if db.isCaching(data) {
	// 	logHandler.Cache.Printf("[CNT]<%v>{SKIP} Count [%+v] [...%v.db] - Caching Enabled", entities.GetStructType(data), entities.GetStructType(data), db.Name)
	// 	return len(inMemoryCache[entities.GetStructType(data)]), nil
	// }
	// for key, value := range connectionPool {
	// 	logHandler.DatabaseLogger.Printf("[CON]<%v>{CONNECTION POOL} Connection Pool [%v] [%v] [codec=%v]", entities.GetStructType(data), key, value.databaseName, value.connection.Node.Codec().Name())
	// }
	return db.connection.Count(data)
}

// CountWhere returns the number of records that match the specified fields.Field and value.
//
// Parameters:
//   - fields.FieldName: The fields.Field to be used for filtering records.
//   - value: The value of the specified fields.Field to filter records.
//   - to: A pointer to the struct representing the type whose records are to be counted.
//
// Returns:
//   - int: The number of records that match the specified criteria.
//   - error: An error object if any issues occur during the counting process; otherwise, nil.
func (db *DB) CountWhere(field entities.Field, value any, to any) (int, error) {
	logHandler.Database.Printf("[COUNT] %v WHERE %+v=%+v [...%v.db]", entities.GetStructType(to), field.String(), value, db.Name)
	if err := entities.IsValidFieldInStruct(field, to); err != nil {
		logHandler.Error.Printf("[COUNT] %v WHERE %+v=%+v [...%v.db] - Error (%e)", entities.GetStructType(to), field.String(), value, db.Name, err)
		return 0, err
	}
	// if db.isCaching(to) {
	// 	logHandler.Cache.Printf("[CNT]<%v>{SKIP} CountWhere (%+v=%+v)[%+v] [...%v.db] - Caching Enabled", entities.GetStructType(to), field.String(), value, entities.GetStructType(to), db.Name)
	// 	// Range through inMemoryCache and count matching entries
	// 	count := 0
	// 	for _, v := range inMemoryCache[entities.GetStructType(to)] {
	// 		val := reflect.ValueOf(v).Elem().FieldByName(string(fieldName))
	// 		if val.IsValid() && val.Interface() == value {
	// 			count++
	// 		}
	// 	}
	// 	return count, nil
	// }
	query := db.connection.Select(q.Eq(field.String(), value))
	count, err := query.Count(to)
	logHandler.Database.Printf("[COUNT] %v WHERE %+v=%+v [...%v.db] - Result: %d", entities.GetStructType(to), field.String(), value, db.Name, count)
	return count, err
}
