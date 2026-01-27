// Data Access Object for the TemplateStoreV3 table
// Template Version: 0.5.10 - 2026-01-26
// Generated 
// Date: 27/01/2026 & 10:17
// Who : matttownsend (orion)

package templateStoreV3

import (
	"context"
	"fmt"
	"strings"

	"github.com/goforj/godump"
	ce "github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-amphora/dao"
	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

// insertOrUpdate performs shared validation/audit and then creates or updates the record.
func (record *TemplateStoreV3) insertOrUpdate(ctx context.Context, note, activity string, auditAction audit.Action, operation string) error {
	isCreateOperation := false
	if strings.EqualFold(operation, "Create") {
		isCreateOperation = true
		if !strings.EqualFold(auditAction.Code(), "Create") {
			return ce.ErrDAOUpdateWrapper(tableName, ce.ErrValidationFailed)
		}
	}

	dao.CheckDAOReadyState(tableName, auditAction, databaseConnectionActive)

	clock := timing.Start(tableName, activity, fmt.Sprintf("%v", record.ID))
	if isCreateOperation {
		if err := record.checkForDuplicate(); err != nil {
			clock.Stop(0)
			return ce.ErrDAOCreateWrapper(tableName, record.ID, err)
		}
	}

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
		logHandler.TraceLogger.Printf("Creating %v record %v %v", tableName, record.Key, record.ID)
		actionError = activeDBConnection.Create(record)
	} else {
		logHandler.TraceLogger.Printf("Updating %v record %v %v", tableName, record.Key, record.ID)
		actionError = activeDBConnection.Update(record)
	}
	if actionError != nil {
		godump.Dump(record)
		updErr := ce.ErrDAOUpdateWrapper(tableName, actionError)
		logHandler.ErrorLogger.Panic(updErr.Error(), actionError)
		clock.Stop(0)
		return updErr
	}

	if !isCreateOperation {
		if err := record.postUpdateProcessing(ctx); err != nil {
			updProcErr := ce.ErrDAOUpdateWrapper(tableName, err)
			logHandler.ErrorLogger.Print(updProcErr.Error())
			clock.Stop(0)
			return updProcErr
		}
	} else {
		if err := record.postCreateProcessing(ctx); err != nil {
			createProcErr := ce.ErrDAOCreateWrapper(tableName, record.ID, err)
			logHandler.ErrorLogger.Print(createProcErr.Error())
			clock.Stop(0)
			return createProcErr
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

	return nil
}
