// Data Access Object for the TripStore table
// Template Version: 0.5.24 - 2026-01-31
// Generated
// Date: 10/02/2026 & 13:16
// Who : matttownsend (orion)

package templateStoreV6

import (
	"context"
	"fmt"
	"strings"

	"github.com/goforj/godump"
	"github.com/mt1976/frantic-amphora/dao"
	"github.com/mt1976/frantic-amphora/dao/audit"
	ce "github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

type op string

const (
	UPDATE op = "Update"
	CREATE op = "Create"
	IMPORT op = "Import"
)

// insertOrUpdate performs shared validation/audit and then creates or updates the record.
func (record *TemplateStoreV6) insertOrUpdate(ctx context.Context, note string, auditAction audit.Action, operation op, isPostProcessingRun bool) (*TemplateStoreV6, error) {
	logHandler.TraceLogger.Printf("INSERTORUPDATE called for %v record: %v operation: %v action: %v isPostProcessing: %t", tableName, record.Raw, operation, auditAction.Code(), isPostProcessingRun)
	logHandler.DatabaseLogger.Printf("INSERTORUPDATE called for %v record %v operation %v action %v postProcessing %t", tableName, record.Key, operation, auditAction.Code(), isPostProcessingRun)
	logHandler.EventLogger.Printf("INSERTORUPDATE called for %v record %v operation %v action %v postProcessing %t", tableName, record.Key, operation, auditAction.Code(), isPostProcessingRun)

	isCreateOperation := false
	if operation == CREATE || operation == IMPORT {
		isCreateOperation = true
		logHandler.TraceLogger.Printf("INSERTORUPDATE determined operation is CREATE type for %v record %v", tableName, record.Key)
		if !strings.EqualFold(auditAction.Code(), "Create") && !strings.EqualFold(auditAction.Code(), "Import") {
			logHandler.WarningLogger.Printf("Audit action '%v' does not match operation type 'Create' for %v record %v. This may lead to incorrect audit records.", auditAction.Code(), tableName, record.Key)
			return New(), ce.ErrDAOUpdateWrapper(tableName, ce.ErrValidationFailed)
		}
	}

	logHandler.TraceLogger.Printf("INSERTORUPDATE for %v record %v isCreate %t operation %v", tableName, record.Key, isCreateOperation, operation)

	locked := false

	if !isCreateOperation && !(operation == UPDATE && isPostProcessingRun) {
		logHandler.TraceLogger.Printf("Attempting to lock %v record %v for %v operation", tableName, record.Key, operation)
		logHandler.LockLogger.Printf("[%v,%v] Locking for %v", tableName, record.Raw, operation)
		record.Lock.Lock()
		locked = true
	} else {
		// For create operations, we will lock the record later in the process after the record has been created and we have a valid key. This is to prevent locking a record that may not be created if there are validation errors or duplicate key issues.
		logHandler.LockLogger.Printf("[%v,%v] Deferring %v lock for new record until after creation", tableName, record.Raw, operation)
	}

	logHandler.TraceLogger.Printf("Starting %v processing for %v record %v isCreate %t", operation, tableName, record.Key, isCreateOperation)

	dao.CheckDAOReadyState(tableName, auditAction, databaseConnectionActive)

	clock := timing.Start(tableName, string(operation), fmt.Sprintf("%v", record.Key))

	logHandler.TraceLogger.Printf("Processing %v record %v", tableName, record.Key)
	// Invoke custom creator logic if defined
	if isCreateOperation {
		if creator != nil {
			logHandler.WarningLogger.Printf("Invoking custom creator for %v record %v", tableName, record.Key)
			id, skip, createdRecord, err := creator(ctx, record)
			if err != nil {
				logHandler.ErrorLogger.Panic(ce.ErrDAOCreateWrapper(tableName, fmt.Sprintf("%v", record.Key), err))
			}
			if skip {
				logHandler.EventLogger.Printf("Custom creator recomends skipping update for %v record %v", tableName, record.Key)
				// No more processing required
				clock.Stop(0)
				return record, nil
			}

			record = createdRecord
			logHandler.EventLogger.Printf("Custom creator completed for %v record %v", tableName, record.Key)

			// Update record with keys
			if id != "" {
				record.Raw = id
				record.Key = idHelpers.Encode(id)
			} else {
				record.Raw = ""
				record.Key = ""
			}
		}
		if record.Key == "" {
			logHandler.WarningLogger.Printf("No Key provided/found, Generating new UUID for %v record", tableName)
			record.Raw = idHelpers.GetUUID()
			record.Key = idHelpers.Encode(record.Raw)
		}

		// godump.Dump(record, "Record after creator processing", record.Key)
		// Check for duplicates on create
		logHandler.TraceLogger.Printf("Checking for duplicate %v record %v", tableName, record.Key)
		err := record.checkForDuplicate()
		if err != nil {
			logHandler.ErrorLogger.Printf("Duplicate check failed for %v record %v: %v", tableName, record.Key, err)
			clock.Stop(0)
			return New(), ce.ErrDAOCreateWrapper(tableName, record.Key, err)
		}
	}
	logHandler.TraceLogger.Printf("Running default/validation processing for %v record %v", tableName, record.Key)
	if calculationError := record.defaultProcessing(); calculationError != nil {
		rtnErr := ce.ErrDAOCaclulationWrapper(tableName, calculationError)
		logHandler.ErrorLogger.Print(rtnErr.Error())
		clock.Stop(0)
		return New(), rtnErr
	}

	logHandler.TraceLogger.Printf("Running validation processing for %v record %v", tableName, record.Key)
	if validationError := record.validationProcessing(); validationError != nil {
		valErr := ce.ErrDAOValidationWrapper(tableName, validationError)
		logHandler.ErrorLogger.Print(valErr.Error())
		clock.Stop(0)
		return New(), valErr
	}

	logHandler.TraceLogger.Printf("Performing audit action for %v record %v", tableName, record.Key)
	auditErr := record.Audit.Action(ctx, auditAction.WithMessage(note))
	if auditErr != nil {
		audErr := ce.ErrDAOUpdateAuditWrapper(tableName, record.Key, auditErr)
		logHandler.ErrorLogger.Print(audErr.Error())
		clock.Stop(0)
		return New(), audErr
	}

	var actionError error
	if isCreateOperation {
		logHandler.TraceLogger.Printf("Creating %v record %v %v", tableName, record.Key, record.Raw)
		actionError = activeDBConnection.Create(record)
		logHandler.TraceLogger.Printf("Created %v record %v %v", tableName, record.Key, record.Raw)

	} else {
		logHandler.TraceLogger.Printf("Updating %v record %v %v", tableName, record.Key, record.Raw)
		actionError = activeDBConnection.Update(record)
		logHandler.TraceLogger.Printf("Updated %v record %v %v", tableName, record.Key, record.Raw)
	}

	logHandler.TraceLogger.Printf("POST %v operation completed for %v record %v %v", operation, tableName, record.Key, record.Raw)
	logHandler.TraceLogger.Printf("POST %v operation completed for %v record %v %v", operation, tableName, record.Key, record.Raw)
	if actionError != nil {
		updErr := ce.ErrDAOUpdateWrapper(tableName, actionError)
		logHandler.ErrorLogger.Panic(updErr.Error(), actionError)
		clock.Stop(0)
		return New(), updErr
	}

	if locked {
		record.Lock.Unlock()
		logHandler.LockLogger.Printf("[%v,%v] Unlocked record after post-processing: %v", tableName, record.Raw, operation)
	}

	// // Unlock record before post-processing to avoid deadlocks if post-processing updates the record
	// if !isCreateOperation {
	// 	logHandler.TraceLogger.Printf("UNLOCKING RECORD before post-processing: %v", record.Raw)
	// 	record.Lock.Unlock()
	// 	logHandler.LockLogger.Printf("UNLOCKED RECORD before post-processing: %v", record.Raw)
	// } else {
	// 	logHandler.TraceLogger.Printf("Deferring unlock for new record until after creation: %v", record.Raw)
	// }

	if isPostProcessingRun {
		// Skip post-processing to avoid infinite loop when record is updated during post-processing
		logHandler.TraceLogger.Printf("Skipping post-processing for %v record %v to avoid infinite loop", tableName, record.Key)
		clock.Stop(1)
		return record, nil
	}

	var err error
	// logHandler.LockLogger.Printf("Post Processing: LOCKING RECORD: %v", record.Raw)
	// record.Lock.Lock()
	// logHandler.TraceLogger.Printf("Post Processing: LOCKED RECORD: %v", record.Raw)

	if !isCreateOperation {
		logHandler.TraceLogger.Printf("Starting post-update processing for %v record %v", tableName, record.Key)
		err = record.postUpdateProcessing(ctx)
		logHandler.TraceLogger.Printf("Post-Update processing completed for %v record %v err %e", tableName, record.Key, err)
	} else {
		logHandler.TraceLogger.Printf("Starting post-create processing for %v record %v", tableName, record.Key)
		err = record.postCreateProcessing(ctx)
		logHandler.TraceLogger.Printf("Post-Create processing completed for %v record %v err %e", tableName, record.Key, err)
	}

	// logHandler.TraceLogger.Printf("Post Processing: UNLOCKING RECORD: %v", record.Raw)
	// record.Lock.Unlock()
	// logHandler.TraceLogger.Printf("Post Processing: UNLOCKED RECORD: %v", record.Raw)

	logHandler.TraceLogger.Printf("POST %v processing completed for %v record %v err %e %+v", operation, tableName, record.Key, err, record)
	if err != nil {
		createProcErr := ce.ErrDAOCreateWrapper(tableName, record.Key, err)
		logHandler.ErrorLogger.Print(createProcErr.Error())
		clock.Stop(0)
		return New(), createProcErr
	}

	// Reset record to updated version (copy values from newRec back to record)
	//logHandler.TraceLogger.Printf("Updating record with values from post %v processing for %v record %v", operation, tableName, record.Key)
	//*record = newRec

	// godump.Dump(record)
	// godump.Dump(newRec)
	// if update {
	// 	if message == "" {
	// 		message = "Post " + string(operation) + " Processing"
	// 	}
	// 	logHandler.TraceLogger.Printf("Post %v processing requires update for %v record %v %v", operation, tableName, record.Key, record.Raw)
	// 	actionError = activeDBConnection.Update(record)
	// 	logHandler.TraceLogger.Printf("Post %v processing requires update for %v record %v %v", operation, tableName, record.Key, record.Raw)
	// 	// err = record.UpdateWithAction(ctx, audit.UPDATE, message)
	// 	if actionError != nil {
	// 		updErr := ce.ErrDAOCreateWrapper(tableName, record.Key, actionError)
	// 		logHandler.ErrorLogger.Panic(updErr.Error())
	// 		clock.Stop(0)
	// 		return New(), updErr
	// 	}
	// }

	logHandler.TraceLogger.Printf("%v operation completed for %v record %v %v %+v", operation, tableName, record.Key, record.Raw, record)
	logHandler.TraceLogger.Printf("%v", godump.DumpJSONStr(record))
	clock.Stop(1)
	return record, nil
}

func (record *TemplateStoreV6) unlock(isCreate bool) error {
	record.Lock.Unlock()
	if isCreate {
		logHandler.LockLogger.Printf("UNLOCKED NEW RECORD: %v", record.Raw)
	} else {
		logHandler.LockLogger.Printf("UNLOCKED RECORD: %v", record.Raw)
	}
	return nil
}

// postGetList runs post-get processing for each record in the list.
func postGetList(ctx context.Context, recordList []*TemplateStoreV6) ([]*TemplateStoreV6, error) {
	clock := timing.Start(tableName, "Process", "POSTGET")
	returnList := []*TemplateStoreV6{}
	for _, record := range recordList {
		if err := record.postGet(ctx); err != nil {
			clock.Stop(0)
			return nil, err
		}
		returnList = append(returnList, record)
	}
	clock.Stop(len(returnList))
	return returnList, nil
}

// postGet runs upgrade/default/validation processing after a record is loaded.
func (record *TemplateStoreV6) postGet(ctx context.Context) error {
	logHandler.TraceLogger.Printf("PostGet processing for %v record %v", tableName, record.Key)
	if upgradeError := record.upgradeProcessing(); upgradeError != nil {
		return upgradeError
	}
	err := record.postGetProcessing(ctx)
	if err != nil {
		logHandler.ErrorLogger.Printf("PostGet processing error for %v record %v: %v", tableName, record.Key, err)
		return err
	}
	return nil
}

// checkForDuplicate checks whether the record key already exists.
func (record *TemplateStoreV6) checkForDuplicate() error {
	dao.CheckDAOReadyState(tableName, audit.PROCESS, databaseConnectionActive)
	logHandler.TraceLogger.Printf("Checking for duplicate %v record %v", tableName, record.Key)
	if duplicateCheck != nil {
		found, err := duplicateCheck(record)
		if err != nil {
			return err
		}
		if found {
			logHandler.TraceLogger.Printf("A duplicate match for %v, %v has been found", tableName, record.Key)
			return ce.ErrDuplicate
		}
		return nil
	}
	logHandler.TraceLogger.Printf("No duplicate check function defined for %v", tableName)
	return nil
}
