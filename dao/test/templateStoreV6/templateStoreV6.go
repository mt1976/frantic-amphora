// Data Access Object for the TripStore table
// Template Version: 0.5.24 - 2026-01-31
// Generated
// Date: 10/02/2026 & 13:16
// Who : matttownsend (orion)

package templateStoreV6

import (
	"context"
	"fmt"
	"reflect"

	"github.com/mt1976/frantic-amphora/dao"
	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/database"
	"github.com/mt1976/frantic-amphora/dao/entities"
	"github.com/mt1976/frantic-amphora/dao/lookup"
	ce "github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

// Count returns the total number of records in the table.
func Count() (int, error) {
	logHandler.DatabaseLogger.Printf("COUNT %v", tableName)
	return activeDBConnection.Count(&TemplateStoreV6{})
}

// CountWhere returns the number of records matching a field/value filter.
func CountWhere(field entities.Field, value any) (int, error) {
	logHandler.DatabaseLogger.Printf("COUNT %v WHERE (%v=%v)", tableName, field.String(), value)
	clock := timing.Start(tableName, "Count", fmt.Sprintf("%v=%v", field.String(), value))
	count, err := activeDBConnection.CountWhere(field, value, &TemplateStoreV6{})
	if err != nil {
		clock.Stop(0)
		return 0, err
	}
	clock.Stop(count)
	return count, nil
}

// GetBy returns a single record matching the given field/value.
func GetBy(field entities.Field, value any) (*TemplateStoreV6, error) {
	logHandler.DatabaseLogger.Printf("SELECT %v WHERE (%v=%v)", tableName, field.String(), value)

	clock := timing.Start(tableName, "Get", fmt.Sprintf("%v=%v", field, value))

	dao.CheckDAOReadyState(tableName, audit.GET, databaseConnectionActive)

	logHandler.TraceLogger.Printf("Data type for %v where %v=%v %v", tableName, field.String(), value, reflect.TypeOf(value).Name())

	if field == Fields.ID && reflect.TypeOf(value).Name() != "int" {
		logHandler.DatabaseLogger.Printf("Invalid data type  Check for %v where %v=%v %v", tableName, field.String(), value, reflect.TypeOf(value).Name())
		msg := "invalid data type. Expected type of %v is int"
		clock.Stop(0)
		logHandler.ErrorLogger.Printf("Invalid data type for %v where %v=%v", tableName, field.String(), value)
		return New(), ce.ErrGetWrapper(tableName, field.String(), value, fmt.Errorf(msg, value))
	}
	logHandler.TraceLogger.Printf("Getting %v record where %v=%v", tableName, field.String(), value)
	record, err := database.GetTyped[TemplateStoreV6](activeDBConnection, field, value)
	if err != nil {
		clock.Stop(0)
		logHandler.DatabaseLogger.Printf("Record not found for %v where %v=%v: %v", tableName, field.String(), value, err)
		return New(), ce.ErrRecordNotFoundWrapper(tableName, field.String(), fmt.Sprintf("%v", value))
	}
	postGetRecord := *record
	logHandler.TraceLogger.Printf("Post Get Processing for %v record %v", tableName, postGetRecord.Key)
	if err := postGetRecord.postGet(context.Background()); err != nil {
		logHandler.DatabaseLogger.Printf("Post Get processing error for %v record %v: %v", tableName, postGetRecord.Key, err)
		clock.Stop(0)
		return New(), ce.ErrGetWrapper(tableName, field.String(), value, err)
	}

	logHandler.TraceLogger.Printf("Post Get Processing completed for %v record %v", tableName, postGetRecord.Key)
	record = &postGetRecord

	clock.Stop(1)
	return &postGetRecord, nil
}

// GetAll returns all TripStore records.
func GetAll() ([]*TemplateStoreV6, error) {
	logHandler.DatabaseLogger.Printf("SELECT %v ALL", tableName)
	dao.CheckDAOReadyState(tableName, audit.GET, databaseConnectionActive)

	clock := timing.Start(tableName, "GetAll", "ALL")
	records, err := database.GetAllTyped[TemplateStoreV6](activeDBConnection)

	logHandler.TraceLogger.Printf("Retrieved %v records from %v", len(records), tableName)

	if err != nil {
		logHandler.ErrorLogger.Print(ce.ErrNotFoundWrapper(tableName, err).Error())
		clock.Stop(0)
		return nil, ce.ErrNotFoundWrapper(tableName, err)
	}

	logHandler.TraceLogger.Printf("Post Get All Processing for %v records", len(records))

	result, err := postGetList(context.Background(), records)
	if err != nil {
		clock.Stop(0)
		return nil, err
	}
	logHandler.TraceLogger.Printf("Post Get All Processing completed for %v records", len(result))
	clock.Stop(len(result))
	return result, nil
}

// GetAllUncached returns all TripStore records without cache.
func GetAllUncached() ([]*TemplateStoreV6, error) {
	logHandler.DatabaseLogger.Printf("SELECT %v ALL", tableName)
	dao.CheckDAOReadyState(tableName, audit.GET, databaseConnectionActive)

	clock := timing.Start(tableName, "GetAllUncached", "ALL")
	records, err := database.GetAllTyped[TemplateStoreV6](activeDBConnection)
	if err != nil {
		clock.Stop(0)
		return nil, ce.ErrNotFoundWrapper(tableName, err)
	}
	result, err := postGetList(context.Background(), records)
	if err != nil {
		clock.Stop(0)
		return nil, err
	}
	clock.Stop(len(result))
	return result, nil
}

// GetAllWhere returns all records matching a field/value filter.
func GetAllWhere(field entities.Field, value any) ([]*TemplateStoreV6, error) {
	logHandler.DatabaseLogger.Printf("SELECT %v WHERE (%v=%v)", tableName, field.String(), value)
	dao.CheckDAOReadyState(tableName, audit.GET, databaseConnectionActive)

	clock := timing.Start(tableName, "GetAllWhere", fmt.Sprintf("%v=%v", field, value))
	records, err := database.GetAllWhereTyped[TemplateStoreV6](activeDBConnection, field, value)
	if err != nil {
		clock.Stop(0)
		return nil, err
	}
	result, err := postGetList(context.Background(), records)
	if err != nil {
		clock.Stop(0)
		return nil, err
	}
	clock.Stop(len(result))
	return result, nil
}

// New returns an empty TripStore record.
func New() *TemplateStoreV6 {
	logHandler.DatabaseLogger.Printf("NEW %v", tableName)
	return &TemplateStoreV6{}
}

// Create constructs and inserts a new TripStore record.
func Create(ctx context.Context, basis *TemplateStoreV6) (*TemplateStoreV6, error) {
	logHandler.DatabaseLogger.Printf("CREATE %v ...", tableName)
	dao.CheckDAOReadyState(tableName, audit.CREATE, databaseConnectionActive)
	logHandler.TraceLogger.Printf("**** Create %v Record: %v %+v", tableName, basis.Key, basis)
	basis, err := (*basis).insertOrUpdate(ctx, fmt.Sprintf("New %v Record", tableName), audit.CREATE, CREATE, false)
	logHandler.TraceLogger.Printf("**** Created %v Record: %v %+v %+v %+v %v", tableName, basis.Key, *basis, basis, &basis, err)
	if err != nil {
		logHandler.ErrorLogger.Panic(ce.ErrDAOCreateWrapper(tableName, basis.ID, err))
		return basis, err
	}
	// godump.Dump(basis)
	logHandler.TraceLogger.Printf("**** Created %v Record: %v %+v", tableName, basis.Key, &basis)
	logHandler.EventLogger.Printf("Created %v record %v", tableName, basis.Key)
	return basis, nil
}

// Create constructs and inserts a new TripStore record.
func importRecord(ctx context.Context, basis *TemplateStoreV6) (*TemplateStoreV6, error) {
	logHandler.DatabaseLogger.Printf("IMPORT %v ...", tableName)
	dao.CheckDAOReadyState(tableName, audit.IMPORT, databaseConnectionActive)
	logHandler.TraceLogger.Printf("**** Import %v Record: %v %+v", tableName, basis.Key, basis)
	basis, err := (*basis).insertOrUpdate(ctx, fmt.Sprintf("Import %v Record", tableName), audit.IMPORT, IMPORT, false)
	logHandler.TraceLogger.Printf("**** Import %v Record: %v %+v %+v %+v %v", tableName, basis.Key, *basis, basis, &basis, err)
	if err != nil {
		logHandler.ErrorLogger.Printf("Error importing %v record %v: %v", tableName, basis.Key, err.Error())
		logHandler.ErrorLogger.Panic(ce.ErrDAOCreateWrapper(tableName, basis.ID, err))
		return basis, err
	}
	logHandler.TraceLogger.Printf("**** Import %v Record: %v %+v", tableName, basis.Key, &basis)
	// logHandler.EventLogger.Printf("Imported %v record %v", tableName, basis.Key)
	return basis, nil
}

// Delete deletes a record by ID.
func Delete(ctx context.Context, id int, note string) error {
	logHandler.DatabaseLogger.Printf("DELETE %v WHERE %v=%v (%v)", tableName, Fields.ID, id, note)

	err := DeleteBy(ctx, Fields.ID, id, note)
	if err != nil {
		return ce.ErrDAODeleteWrapper(tableName, Fields.ID.String(), id, err)
	}
	// logHandler.EventLogger.Printf("Deleted %v record %v", tableName, id)
	return err
}

// DeleteBy deletes a record by field/value.
func DeleteBy(ctx context.Context, field entities.Field, value any, note string) error {
	logHandler.DatabaseLogger.Printf("DELETE %v WHERE %v=%v (%v)", tableName, field, value, note)
	dao.CheckDAOReadyState(tableName, audit.DELETE, databaseConnectionActive)

	clock := timing.Start(tableName, "Delete", fmt.Sprintf("%v=%v", field.String(), value))

	recordList, err := GetAllWhere(field, value)
	if err != nil {
		clock.Stop(0)
		return ce.ErrDAODeleteWrapper(tableName, field.String(), value, err)
	}

	if len(recordList) == 0 {
		clock.Stop(0)
		return ce.ErrRecordNotFoundWrapper(tableName, field.String(), fmt.Sprintf("%v", value))
	}

	for _, record := range recordList {
		logHandler.TraceLogger.Printf("Deleting %v record %v", tableName, record.Key)
		logHandler.TraceLogger.Printf("Pre-Delete Audit Processing for %v record %v", tableName, record.Key)
		if err := record.Audit.Action(ctx, audit.DELETE.WithMessage(note)); err != nil {
			clock.Stop(0)
			return ce.ErrDAOUpdateAuditWrapper(tableName, value, err)
		}

		logHandler.TraceLogger.Printf("Delete Pre-Delete Processing for %v record %v", tableName, record.Key)
		if err := record.preDeleteProcessing(ctx); err != nil {
			clock.Stop(0)
			return ce.ErrDAODeleteWrapper(tableName, field.String(), value, err)
		}

		logHandler.TraceLogger.Printf("Deleting %v record %v", tableName, record.Key)
		if err := activeDBConnection.Delete(record); err != nil {
			clock.Stop(0)
			return ce.ErrDAODeleteWrapper(tableName, field.String(), value, err)
		}
		logHandler.TraceLogger.Printf("Deleting Post-Delete Processing for %v record %v", tableName, record.Key)
		if err := record.postDeleteProcessing(ctx); err != nil {
			clock.Stop(0)
			return ce.ErrDAODeleteWrapper(tableName, field.String(), value, err)
		}
	}

	clock.Stop(1)
	return nil
}

// Validate runs record validation and returns an error if invalid.
func (record *TemplateStoreV6) Validate() error {
	logHandler.DatabaseLogger.Printf("Validating %v record %v", tableName, record.Key)
	return record.validationProcessing()
}

// Update persists changes to an existing record.
func (record *TemplateStoreV6) Update(ctx context.Context, note string) error {
	logHandler.DatabaseLogger.Printf("Updating %v record %v (%v)", tableName, record.Key, note)
	record, err := record.insertOrUpdate(ctx, note, audit.UPDATE, UPDATE, false)
	logHandler.DatabaseLogger.Printf("Update record %v %v", record.Key, note)

	return err
}

func (record *TemplateStoreV6) PostUpdateUpdate(ctx context.Context, note string) error {
	logHandler.DatabaseLogger.Printf("Updating %v record %v during post-update processing (%v)", tableName, record.Key, note)
	record, err := record.insertOrUpdate(ctx, note, audit.UPDATE, UPDATE, true)
	logHandler.DatabaseLogger.Printf("Post-update Update record %v %v", record.Key, note)
	return err
}

// UpdateWithAction persists changes using the provided audit action.
func (record *TemplateStoreV6) UpdateWithAction(ctx context.Context, auditAction audit.Action, note string) error {
	logHandler.DatabaseLogger.Printf("%ving %v record %v with action %v (%v)", UPDATE, tableName, record.Key, auditAction.Code(), note)
	record, err := record.insertOrUpdate(ctx, note, auditAction, UPDATE, false)
	logHandler.DatabaseLogger.Printf("%ved %v record %v with action %v (%v)", UPDATE, tableName, record.Key, auditAction.Code(), note)
	return err
}

// Clone returns a copy of the record using templateClone.
func (record *TemplateStoreV6) Clone(ctx context.Context) (*TemplateStoreV6, error) {
	logHandler.DatabaseLogger.Printf("Cloning %v record %v", tableName, record.Key)
	return templateClone(ctx, record)
}

// GetDefaultLookup returns the default lookup for this table.
func GetDefaultLookup() (lookup.Lookup, error) {
	logHandler.TraceLogger.Printf("Getting default lookup for %v", tableName)
	return GetLookup(Fields.Key, Fields.Raw)
}

// GetLookup builds a lookup of key/value pairs from all records.
func GetLookup(field, value entities.Field) (lookup.Lookup, error) {
	logHandler.TraceLogger.Printf("Building lookup for %v field %v value %v", tableName, field.String(), value.String())
	dao.CheckDAOReadyState(tableName, audit.PROCESS, databaseConnectionActive)

	clock := timing.Start(tableName, "Lookup", "BUILD")

	recordList, err := GetAll()
	if err != nil {
		lkpErr := ce.ErrDAOLookupWrapper(tableName, field.String(), value, err)
		logHandler.ErrorLogger.Print(lkpErr.Error())
		clock.Stop(0)
		return lookup.Lookup{}, lkpErr
	}

	var rtnLookup lookup.Lookup
	rtnLookup.Data = make([]lookup.LookupData, 0)

	for _, record := range recordList {
		key := reflect.ValueOf(record).Elem().FieldByName(field.String()).Interface().(string)
		val := reflect.ValueOf(record).Elem().FieldByName(value.String()).Interface().(string)
		rtnLookup.Data = append(rtnLookup.Data, lookup.LookupData{Key: key, Value: val})
	}

	clock.Stop(len(rtnLookup.Data))
	return rtnLookup, nil
}

// Drop drops the underlying database bucket/table for this entity.
func Drop() error {
	logHandler.DatabaseLogger.Printf("Drop %v", tableName)
	err := activeDBConnection.Drop(TemplateStoreV6{})
	if err != nil {
		return err
	}
	if postDrop != nil {
		if err := postDrop(context.Background()); err != nil {
			return err
		}
	}
	return nil
}

// ClearDown deletes all records from this table.
func ClearDown(ctx context.Context) error {
	logHandler.DatabaseLogger.Printf("ClearDown %v", tableName)

	dao.CheckDAOReadyState(tableName, audit.PROCESS, databaseConnectionActive)

	clock := timing.Start(tableName, "Clear", "INITIALISE")

	recordList, err := GetAll()
	if err != nil {
		logHandler.ErrorLogger.Print(ce.ErrDAOInitialisationWrapper(tableName, err).Error())
		clock.Stop(0)
		return ce.ErrDAOInitialisationWrapper(tableName, err)
	}

	count := 0
	logHandler.TraceLogger.Printf("Clearing %v records", len(recordList))

	for i, record := range recordList {
		logHandler.TraceLogger.Printf("(%v/%v) DELETE %v WHERE %v=%v", i+1, len(recordList), tableName, Fields.ID, record.Key)

		delErr := Delete(ctx, record.ID, fmt.Sprintf("Clearing %v %v", tableName, record.Key))
		if delErr != nil {
			logHandler.ErrorLogger.Print(ce.ErrDAOInitialisationWrapper(tableName, delErr).Error())
			continue
		}
		count++
	}

	if postClearDown != nil {
		if err := postClearDown(ctx); err != nil {
			logHandler.ErrorLogger.Print(ce.ErrDAOInitialisationWrapper(tableName, err).Error())
			clock.Stop(0)
			return ce.ErrDAOInitialisationWrapper(tableName, err)
		}
	}
	clock.Stop(count)
	//	logHandler.DatabaseLogger.Printf("Cleared down %v", tableName)
	return nil
}
