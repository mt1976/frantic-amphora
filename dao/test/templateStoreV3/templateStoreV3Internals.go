// Data Access Object for the TemplateStoreV3 table
// Template Version: 0.5.11 - 2026-01-27
// Generated
// Date: 27/01/2026 & 15:01
// Who : matttownsend (orion)

package templateStoreV3

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
)

// insertOrUpdate performs shared validation/audit and then creates or updates the record.
func (record *TemplateStoreV3) insertOrUpdate(ctx context.Context, note string, auditAction audit.Action, operation op) error {
	isCreateOperation := false
	if operation == CREATE {
		isCreateOperation = true
		if !strings.EqualFold(auditAction.Code(), "Create") {
			return ce.ErrDAOUpdateWrapper(tableName, ce.ErrValidationFailed)
		}
	}

	logHandler.DatabaseLogger.Printf("Starting %v processing for %v record %v isCreate %t", operation, tableName, record.Key, isCreateOperation)

	dao.CheckDAOReadyState(tableName, auditAction, databaseConnectionActive)

	clock := timing.Start(tableName, string(operation), fmt.Sprintf("%v", record.ID))
	if isCreateOperation {
		// Check for duplicates on create
		logHandler.DatabaseLogger.Printf("Checking for duplicate %v record %v", tableName, record.Key)
		if err := record.checkForDuplicate(); err != nil {
			clock.Stop(0)
			return ce.ErrDAOCreateWrapper(tableName, record.ID, err)
		}
	}

	logHandler.DatabaseLogger.Printf("Processing %v record %v", tableName, record.Key)
	// Invoke custom creator logic if defined
	if isCreateOperation {
		if creator != nil {
			logHandler.DatabaseLogger.Printf("Invoking custom creator for %v record %v", tableName, record.Key)
			id, skip, createdRecord, err := creator(ctx, *record)
			if err != nil {
				logHandler.ErrorLogger.Panic(ce.ErrDAOCreateWrapper(tableName, fmt.Sprintf("%v", record.Key), err))
			}
			if !skip {
				record = &createdRecord
				logHandler.DatabaseLogger.Printf("Custom creator completed for %v record %v", tableName, record.Key)
			} else {
				logHandler.DatabaseLogger.Printf("Custom creator skipped for %v record %v", tableName, record.Key)
			}
			record.Raw = id
			record.Key = idHelpers.Encode(id)
		}
	}
	logHandler.DatabaseLogger.Printf("Running default/validation processing for %v record %v", tableName, record.Key)
	if calculationError := record.defaultProcessing(); calculationError != nil {
		rtnErr := ce.ErrDAOCaclulationWrapper(tableName, calculationError)
		logHandler.ErrorLogger.Print(rtnErr.Error())
		clock.Stop(0)
		return rtnErr
	}

	if validationError := record.validationProcessing(); validationError != nil {
		valErr := ce.ErrDAOValidationWrapper(tableName, validationError)
		logHandler.ErrorLogger.Print(valErr.Error())
		clock.Stop(0)
		return valErr
	}

	auditErr := record.Audit.Action(ctx, auditAction.WithMessage(note))
	if auditErr != nil {
		audErr := ce.ErrDAOUpdateAuditWrapper(tableName, record.ID, auditErr)
		logHandler.ErrorLogger.Print(audErr.Error())
		clock.Stop(0)
		return audErr
	}

	var actionError error
	if isCreateOperation {
		logHandler.DatabaseLogger.Printf("Creating %v record %v %v", tableName, record.Key, record.ID)
		actionError = activeDBConnection.Create(record)

	} else {
		logHandler.DatabaseLogger.Printf("Updating %v record %v %v", tableName, record.Key, record.ID)
		actionError = activeDBConnection.Update(record)
	}
	logHandler.DatabaseLogger.Printf("%v operation completed for %v record %v", operation, tableName, record.Key)
	if actionError != nil {
		godump.Dump(record)
		updErr := ce.ErrDAOUpdateWrapper(tableName, actionError)
		logHandler.ErrorLogger.Panic(updErr.Error(), actionError)
		clock.Stop(0)
		return updErr
	}
	var err error
	var update bool = false
	var message string = ""
	if !isCreateOperation {
		logHandler.DatabaseLogger.Printf("Starting post-update processing for %v record %v", tableName, record.Key)
		err, update, message = record.postUpdateProcessing(ctx)
		logHandler.DatabaseLogger.Printf("Post-Update processing completed for %v record %v err %e", tableName, record.Key, err)
	} else {
		logHandler.DatabaseLogger.Printf("Starting post-create processing for %v record %v", tableName, record.Key)
		err, update, message = record.postCreateProcessing(ctx)
		logHandler.DatabaseLogger.Printf("Post-Create processing completed for %v record %v err %e", tableName, record.Key, err)
	}
	if err != nil {
		createProcErr := ce.ErrDAOCreateWrapper(tableName, record.ID, err)
		logHandler.ErrorLogger.Print(createProcErr.Error())
		clock.Stop(0)
		return createProcErr
	}
	if update {
		if message == "" {
			message = "Post " + string(operation) + " Processing"
		}
		logHandler.DatabaseLogger.Printf("Post %v processing requires update for %v record %v %v", operation, tableName, record.Key, record.ID)
		actionError = activeDBConnection.Update(record)
		//err = record.UpdateWithAction(ctx, audit.UPDATE, message)
		if actionError != nil {
			updErr := ce.ErrDAOCreateWrapper(tableName, record.ID, actionError)
			logHandler.ErrorLogger.Panic(updErr.Error())
			clock.Stop(0)
			return updErr
		}
	}

	clock.Stop(1)
	return nil
}

// postGetList runs post-get processing for each record in the list.
func postGetList(ctx context.Context, recordList []TemplateStoreV3) ([]TemplateStoreV3, error) {
	clock := timing.Start(tableName, "Process", "POSTGET")
	returnList := []TemplateStoreV3{}
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
func (record *TemplateStoreV3) postGet(ctx context.Context) error {
	if upgradeError := record.upgradeProcessing(); upgradeError != nil {
		return upgradeError
	}
	//if defaultingError := record.defaultProcessing(); defaultingError != nil {
	//	return defaultingError
	//}
	//if validationError := record.validationProcessing(); validationError != nil {
	//	return validationError
	//}
	return record.postGetProcessing(ctx)
}

// checkForDuplicate checks whether the record key already exists.
func (record *TemplateStoreV3) checkForDuplicate() error {
	dao.CheckDAOReadyState(tableName, audit.PROCESS, databaseConnectionActive)
	logHandler.TraceLogger.Printf("Checking for duplicate %v record %v", tableName, record.Key)
	if duplicateCheck != nil {
		found, err := duplicateCheck(record)
		if err != nil {
			return err
		}
		if found {
			logHandler.WarningLogger.Printf("Duplicate %v, %v already in use", tableName, record.Key)
			return ce.ErrDuplicate
		}
		return nil
	}
	logHandler.TraceLogger.Printf("No duplicate check function defined for %v", tableName)

	return nil
}
