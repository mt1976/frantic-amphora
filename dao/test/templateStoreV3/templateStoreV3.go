// Data Access Object for the TemplateStoreV3 table
// Template Version: 0.5.10 - 2026-01-26
// Generated
// Date: 27/01/2026 & 10:17
// Who : matttownsend (orion)

package templateStoreV3

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
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

// Count returns the total number of records in the table.
func Count() (int, error) {
	//logHandler.DatabaseLogger.Printf("COUNT %v", tableName)
	return activeDBConnection.Count(&TemplateStoreV3{})
}

// CountWhere returns the number of records matching a field/value filter.
func CountWhere(field entities.Field, value any) (int, error) {
	//logHandler.DatabaseLogger.Printf("COUNT %v WHERE (%v=%v)", tableName, field.String(), value)
	clock := timing.Start(tableName, "Count", fmt.Sprintf("%v=%v", field.String(), value))
	count, err := activeDBConnection.CountWhere(field, value, &TemplateStoreV3{})
	if err != nil {
		clock.Stop(0)
		return 0, err
	}
	clock.Stop(count)
	return count, nil
}

// GetBy returns a single record matching the given field/value.
func GetBy(field entities.Field, value any) (TemplateStoreV3, error) {
	//	logHandler.DatabaseLogger.Printf("SELECT %v WHERE (%v=%v)", tableName, field.String(), value)
	clock := timing.Start(tableName, "Get", fmt.Sprintf("%v=%v", field, value))

	dao.CheckDAOReadyState(tableName, audit.GET, databaseConnectionActive)

	if field == Fields.ID && reflect.TypeOf(value).Name() != "int" {
		msg := "invalid data type. Expected type of %v is int"
		clock.Stop(0)
		return TemplateStoreV3{}, ce.ErrGetWrapper(tableName, field.String(), value, fmt.Errorf(msg, value))
	}

	record, err := database.GetTyped[TemplateStoreV3](activeDBConnection, field, value)
	if err != nil {
		clock.Stop(0)
		return TemplateStoreV3{}, ce.ErrRecordNotFoundWrapper(tableName, field.String(), fmt.Sprintf("%v", value))
	}
	if err := record.postGet(context.Background()); err != nil {
		clock.Stop(0)
		return TemplateStoreV3{}, ce.ErrGetWrapper(tableName, field.String(), value, err)
	}

	clock.Stop(1)
	return record, nil
}

// GetAll returns all TemplateStoreV3 records.
func GetAll() ([]TemplateStoreV3, error) {
	//	logHandler.DatabaseLogger.Printf("SELECT %v ALL", tableName)
	dao.CheckDAOReadyState(tableName, audit.GET, databaseConnectionActive)

	clock := timing.Start(tableName, "GetAll", "ALL")
	records, err := database.GetAllTyped[TemplateStoreV3](activeDBConnection)
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

// GetAll returns all TemplateStoreV3 records.
func GetAllUncached() ([]TemplateStoreV3, error) {
	//	logHandler.DatabaseLogger.Printf("SELECT %v ALL", tableName)
	dao.CheckDAOReadyState(tableName, audit.GET, databaseConnectionActive)

	clock := timing.Start(tableName, "GetAllUncached", "ALL")
	records, err := database.GetAllTyped[TemplateStoreV3](activeDBConnection)
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
func GetAllWhere(field entities.Field, value any) ([]TemplateStoreV3, error) {
	//	logHandler.DatabaseLogger.Printf("SELECT %v WHERE (%v=%v)", tableName, field.String(), value)
	dao.CheckDAOReadyState(tableName, audit.GET, databaseConnectionActive)

	clock := timing.Start(tableName, "GetAllWhere", fmt.Sprintf("%v=%v", field, value))
	records, err := database.GetAllWhereTyped[TemplateStoreV3](activeDBConnection, field, value)
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

// New returns an empty TemplateStoreV3 record.
func New() TemplateStoreV3 {
	return TemplateStoreV3{}
}

// Create constructs and inserts a new TemplateStoreV3 record.
func Create(ctx context.Context, basis TemplateStoreV3) (TemplateStoreV3, error) {
	dao.CheckDAOReadyState(tableName, audit.CREATE, databaseConnectionActive)

	clock := timing.Start(tableName, "Create", "Inserting new record")

	id, skip, record, err := creator(ctx, basis)
	if err != nil {
		logHandler.ErrorLogger.Panic(ce.ErrDAOCreateWrapper(tableName, fmt.Sprintf("%v", basis), err))
	}

	if skip {
		logHandler.InfoLogger.Printf("Creation of %v (%v) record has been skipped", tableName, id)

		createdRecord, err := doPostProcessing(ctx, record)
		if err != nil {
			logHandler.ErrorLogger.Panic(ce.ErrDAOCreateWrapper(tableName, record.ID, err))
		}

		clock.Stop(1)
		return createdRecord, nil
	}

	record.Key = idHelpers.Encode(id)
	record.Raw = id

	auditErr := record.Audit.Action(ctx, audit.CREATE.WithMessage(fmt.Sprintf("New %v created %v", tableName, basis)))
	if auditErr != nil {
		logHandler.ErrorLogger.Panic(ce.ErrDAOUpdateAuditWrapper(tableName, record.ID, auditErr))
	}

	writeErr := activeDBConnection.Create(&record)
	if writeErr != nil {
		logHandler.ErrorLogger.Panic(ce.ErrDAOCreateWrapper(tableName, record.ID, writeErr))
	}

	createdRecord, err := doPostProcessing(ctx, record)
	if err != nil {
		logHandler.ErrorLogger.Panic(ce.ErrDAOCreateWrapper(tableName, record.ID, err))
	}
	clock.Stop(1)
	return createdRecord, nil
}

func doPostProcessing(ctx context.Context, record TemplateStoreV3) (TemplateStoreV3, error) {
	if err := record.postCreateProcessing(ctx); err != nil {
		logHandler.ErrorLogger.Panic(ce.ErrDAOCreateWrapper(tableName, record.ID, err))
	}

	createdRecord, getErr := GetBy(Fields.Key, record.Key)
	if getErr != nil {
		logHandler.ErrorLogger.Panic(ce.ErrDAOCreateWrapper(tableName, record.ID, getErr))
		return record, getErr
	}
	return createdRecord, nil
}

// Delete deletes a record by ID.
func Delete(ctx context.Context, id int, note string) error {
	return DeleteBy(ctx, Fields.ID, id, note)
}

// DeleteBy deletes a record by field/value.
func DeleteBy(ctx context.Context, field entities.Field, value any, note string) error {
	//	logHandler.DatabaseLogger.Printf("DELETE %v WHERE %v=%v", tableName, field, value)
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
		if err := record.Audit.Action(ctx, audit.DELETE.WithMessage(note)); err != nil {
			clock.Stop(0)
			return ce.ErrDAOUpdateAuditWrapper(tableName, value, err)
		}

		if err := record.preDeleteProcessing(ctx); err != nil {
			clock.Stop(0)
			return ce.ErrDAODeleteWrapper(tableName, field.String(), value, err)
		}

		if err := activeDBConnection.Delete(&record); err != nil {
			clock.Stop(0)
			return ce.ErrDAODeleteWrapper(tableName, field.String(), value, err)
		}

		if err := record.postDeleteProcessing(ctx); err != nil {
			clock.Stop(0)
			return ce.ErrDAODeleteWrapper(tableName, field.String(), value, err)
		}
	}

	clock.Stop(1)
	return nil
}

// Validate runs record validation and returns an error if invalid.
func (record *TemplateStoreV3) Validate() error {
	return record.validationProcessing()
}

// Update persists changes to an existing record.
func (record *TemplateStoreV3) Update(ctx context.Context, note string) error {
	return record.insertOrUpdate(ctx, note, "Update", audit.UPDATE, "Update")
}

// UpdateWithAction persists changes using the provided audit action.
func (record *TemplateStoreV3) UpdateWithAction(ctx context.Context, auditAction audit.Action, note string) error {
	return record.insertOrUpdate(ctx, note, "Update", auditAction, "Update")
}

// Create inserts a new record.
func (record *TemplateStoreV3) Create(ctx context.Context, note string) error {
	return record.insertOrUpdate(ctx, note, "Create", audit.CREATE, "Create")
}

// Clone returns a copy of the record using templateClone.
func (record *TemplateStoreV3) Clone(ctx context.Context) (TemplateStoreV3, error) {
	logHandler.DatabaseLogger.Printf("Clone %v id: %v", tableName, record.Key)
	return templateClone(ctx, *record)
}

// GetDefaultLookup returns the default lookup for this table.
func GetDefaultLookup() (lookup.Lookup, error) {
	return GetLookup(Fields.Key, Fields.Raw)
}

// GetLookup builds a lookup of key/value pairs from all records.
func GetLookup(field, value entities.Field) (lookup.Lookup, error) {
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

	for _, a := range recordList {
		key := reflect.ValueOf(a).FieldByName(field.String()).Interface().(string)
		val := reflect.ValueOf(a).FieldByName(value.String()).Interface().(string)
		rtnLookup.Data = append(rtnLookup.Data, lookup.LookupData{Key: key, Value: val})
	}

	clock.Stop(len(rtnLookup.Data))
	return rtnLookup, nil
}

// Drop drops the underlying database bucket/table for this entity.
func Drop() error {
	logHandler.TraceLogger.Printf("Drop %v", tableName)
	err := activeDBConnection.Drop(TemplateStoreV3{})
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
	logHandler.TraceLogger.Printf("ClearDown %v", tableName)

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
		logHandler.TraceLogger.Printf("(%v/%v) DELETE %v WHERE %v=%v", i+1, len(recordList), tableName, Fields.ID, record.ID)

		delErr := Delete(ctx, record.ID, fmt.Sprintf("Clearing %v %v @ initialisation ", tableName, record.ID))
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
