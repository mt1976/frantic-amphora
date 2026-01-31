// Data Access Object for the TemplateStoreV3 table
// Template Version: 0.5.12 - 2026-01-27
// Generated
// Date: 27/01/2026 & 15:01
// Who : matttownsend (orion)

package templateStoreV3

import (
	"context"

	"github.com/mt1976/frantic-amphora/dao"
	"github.com/mt1976/frantic-core/logHandler"
)

type (
	creatorFunc        func(ctx context.Context, sourceData *TemplateStoreV3) (string, bool, *TemplateStoreV3, error)
	upgraderFunc       func(*TemplateStoreV3) (*TemplateStoreV3, error)
	defaulterFunc      func(*TemplateStoreV3) error
	validatorFunc      func(*TemplateStoreV3) error
	preDeleteFunc      func(ctx context.Context, record *TemplateStoreV3) error
	postGetFunc        func(ctx context.Context, record *TemplateStoreV3) error
	clonerFunc         func(ctx context.Context, sourceRecord *TemplateStoreV3) (*TemplateStoreV3, error)
	duplicateCheckFunc func(*TemplateStoreV3) (bool, error)
	workerFunc         func(string, string)
	postCreateFunc     func(ctx context.Context, record *TemplateStoreV3) (error, bool, *TemplateStoreV3, string)
	postUpdateFunc     func(ctx context.Context, record *TemplateStoreV3) (error, bool, *TemplateStoreV3, string)
	postDeleteFunc     func(ctx context.Context, record *TemplateStoreV3) error
	postCloneFunc      func(ctx context.Context, record *TemplateStoreV3) error
	postDropFunc       func(ctx context.Context) error
)

var (
	creator        creatorFunc
	upgrader       upgraderFunc
	defaulter      defaulterFunc
	validator      validatorFunc
	preDelete      preDeleteFunc
	postGet        postGetFunc
	cloner         clonerFunc
	duplicateCheck duplicateCheckFunc
	worker         workerFunc
	postCreate     postCreateFunc
	postUpdate     postUpdateFunc
	postDelete     postDeleteFunc
	postClone      postCloneFunc
	postDrop       postDropFunc
	postClearDown  postDropFunc
)

// RegisterCreator registers a creator function for TemplateStoreV3.
func RegisterCreator(fn creatorFunc) {
	logHandler.EventLogger.Printf("[REGISTER] Creator for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] Creator for %v (%v)", tableName, dao.GetFunctionName(fn))
	creator = fn
}

// RegisterPostCreate registers a post-create function for TemplateStoreV3.
func RegisterPostCreate(fn postCreateFunc) {
	logHandler.EventLogger.Printf("[REGISTER] PostCreate for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] PostCreate for %v (%v)", tableName, dao.GetFunctionName(fn))
	postCreate = fn
}

// RegisterPostUpdate registers a post-update function for TemplateStoreV3.
func RegisterPostUpdate(fn postUpdateFunc) {
	logHandler.EventLogger.Printf("[REGISTER] PostUpdate for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] PostUpdate for %v (%v)", tableName, dao.GetFunctionName(fn))
	postUpdate = fn
}

// RegisterPostDelete registers a post-delete function for TemplateStoreV3.
func RegisterPostDelete(fn postDeleteFunc) {
	logHandler.EventLogger.Printf("[REGISTER] PostDelete for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] PostDelete for %v (%v)", tableName, dao.GetFunctionName(fn))
	postDelete = fn
}

// RegisterPostClone registers a post-clone function for TemplateStoreV3.
func RegisterPostClone(fn postCloneFunc) {
	logHandler.EventLogger.Printf("[REGISTER] PostClone for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] PostClone for %v (%v)", tableName, dao.GetFunctionName(fn))
	postClone = fn
}

// RegisterPostDrop registers a post-drop function for TemplateStoreV3.
func RegisterPostDrop(fn postDropFunc) {
	logHandler.EventLogger.Printf("[REGISTER] PostDrop for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] PostDrop for %v (%v)", tableName, dao.GetFunctionName(fn))
	postDrop = fn
}

// RegisterPostClearDown registers a post-clear-down function for TemplateStoreV3.
func RegisterPostClearDown(fn postDropFunc) {
	logHandler.EventLogger.Printf("[REGISTER] PostClearDown for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] PostClearDown for %v (%v)", tableName, dao.GetFunctionName(fn))
	postClearDown = fn
}

// RegisterUpgrader registers an upgrader function for TemplateStoreV3.
func RegisterUpgrader(fn upgraderFunc) {
	logHandler.EventLogger.Printf("[REGISTER] Upgrader for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] Upgrader for %v (%v)", tableName, dao.GetFunctionName(fn))
	upgrader = fn
}

// RegisterDefaulter registers a defaulter function for TemplateStoreV3.
func RegisterDefaulter(fn defaulterFunc) {
	logHandler.EventLogger.Printf("[REGISTER] Defaulter for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] Defaulter for %v (%v)", tableName, dao.GetFunctionName(fn))
	defaulter = fn
}

// RegisterValidator registers a validator function for TemplateStoreV3.
func RegisterValidator(fn validatorFunc) {
	logHandler.EventLogger.Printf("[REGISTER] Validator for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] Validator for %v (%v)", tableName, dao.GetFunctionName(fn))
	validator = fn
}

// RegisterPreDelete registers a pre-delete function for TemplateStoreV3.
func RegisterPreDelete(fn preDeleteFunc) {
	logHandler.EventLogger.Printf("[REGISTER] PreDelete for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] PreDelete for %v (%v)", tableName, dao.GetFunctionName(fn))
	preDelete = fn
}

// RegisterPostGet registers a post-get function for TemplateStoreV3.
func RegisterPostGet(fn postGetFunc) {
	logHandler.EventLogger.Printf("[REGISTER] PostGet for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] PostGet for %v (%v)", tableName, dao.GetFunctionName(fn))
	postGet = fn
}

// RegisterCloner registers a cloner function for TemplateStoreV3.
func RegisterCloner(fn clonerFunc) {
	logHandler.EventLogger.Printf("[REGISTER] Cloner for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] Cloner for %v (%v)", tableName, dao.GetFunctionName(fn))
	cloner = fn
}

// RegisterDuplicateCheck registers a duplicate check function for TemplateStoreV3.
func RegisterDuplicateCheck(fn duplicateCheckFunc) {
	logHandler.EventLogger.Printf("[REGISTER] DuplicateCheck for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] DuplicateCheck for %v (%v)", tableName, dao.GetFunctionName(fn))
	duplicateCheck = fn
}

// RegisterWorker registers a worker function for TemplateStoreV3.
func RegisterWorker(fn workerFunc) {
	logHandler.EventLogger.Printf("[REGISTER] Worker for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.DatabaseLogger.Printf("[REGISTER] Worker for %v (%v)", tableName, dao.GetFunctionName(fn))
	worker = fn
}

// upgradeProcessing performs any one-time upgrade or migration logic on the record.
func (record *TemplateStoreV3) upgradeProcessing() error {
	if upgrader != nil {
		logHandler.DatabaseLogger.Printf("[UPGRADE] record %v of %v", record.Key, TableName.String())
		updatedRecord, err := upgrader(record)
		if err != nil {
			return err
		}
		*record = *updatedRecord
		logHandler.DatabaseLogger.Printf("[UPGRADE] Upgrade complete for record %v of %v", record.Key, TableName.String())
	}
	return nil
}

// defaultProcessing applies any default values prior to validation and persistence.
func (record *TemplateStoreV3) defaultProcessing() error {
	if defaulter != nil {
		logHandler.DatabaseLogger.Printf("[DEFAULT] Applying defaults to record %v of %v", record.Key, TableName.String())
		err := defaulter(record)
		logHandler.DatabaseLogger.Printf("[DEFAULT] Defaults applied to record %v of %v", record.Key, TableName.String())
		return err
	}
	return nil
}

// validationProcessing validates the record and returns an error if it is invalid.
func (record *TemplateStoreV3) validationProcessing() error {
	if validator != nil {
		logHandler.DatabaseLogger.Printf("[VALIDATE] Validating record %v of %v", record.Key, TableName.String())
		err := validator(record)
		logHandler.DatabaseLogger.Printf("[VALIDATE] Record %v of %v is valid", record.Key, TableName.String())
		return err
	}
	return nil
}

// postGetProcessing runs any post-load processing after a record is retrieved.
func (h *TemplateStoreV3) postGetProcessing(ctx context.Context) error {
	if postGet != nil {
		logHandler.DatabaseLogger.Printf("[POSTGET] Processing for %v Record: %v", TableName.String(), h.Key)
		err := postGet(ctx, h)
		logHandler.DatabaseLogger.Printf("[POSTGET] Processing complete for %v Record: %v", TableName.String(), h.Key)
		return err
	}
	return nil
}

// preDeleteProcessing runs any checks or actions required before delete.
func (record *TemplateStoreV3) preDeleteProcessing(ctx context.Context) error {
	if preDelete != nil {
		logHandler.DatabaseLogger.Printf("[PREDELETE] Processing for %v Record: %v", TableName.String(), record.Key)
		err := preDelete(ctx, record)
		logHandler.DatabaseLogger.Printf("[PREDELETE] Processing complete for %v Record: %v", TableName.String(), record.Key)
		return err
	}
	return nil
}

// templateClone contains the package's clone logic.
func templateClone(ctx context.Context, source *TemplateStoreV3) (*TemplateStoreV3, error) {
	if cloner != nil {
		logHandler.DatabaseLogger.Printf("[CLONE] Cloning record %v of %v", source.Key, TableName.String())
		nr, err := cloner(ctx, source)
		logHandler.DatabaseLogger.Printf("[CLONE] Cloning complete for record %v of %v", source.Key, TableName.String())
		return nr, err
	}
	return New(), nil
}

// PostCreate runs any post-create processing after a record is created.
func (record *TemplateStoreV3) postCreateProcessing(ctx context.Context) (error, bool, *TemplateStoreV3, string) {
	if postCreate != nil {

		logHandler.DatabaseLogger.Printf("[POSTCREATE] Processing for %v Record: %v", TableName.String(), record.Key)

		err, recordUpdated, updatedRecord, feedbackMessage := postCreate(ctx, record)
		if err != nil {
			return err, false, New(), ""
		}
		if !recordUpdated {
			return nil, false, New(), ""
		}
		if feedbackMessage == "" {
			feedbackMessage = "Post Create Processing"
		}
		// Return the TemplateStoreV3 record with the new profile and notes
		logHandler.DatabaseLogger.Printf("[POSTCREATE] Processing complete for %v Record: %v", TableName.String(), record.Key)
		return nil, true, updatedRecord, feedbackMessage
	}
	return nil, false, New(), ""
}

func (record *TemplateStoreV3) postUpdateProcessing(ctx context.Context) (error, bool, *TemplateStoreV3, string) {
	if postUpdate != nil {
		logHandler.DatabaseLogger.Printf("[POSTUPDATE] Processing for %v Record: %v", TableName.String(), record.Key)

		err, recordUpdated, updatedRecord, feedbackMessage := postUpdate(ctx, record)
		if err != nil {
			logHandler.ErrorLogger.Printf("Error during postUpdateProcessing for %v Record: %v Error: %v", TableName.String(), record.Key, err.Error())
			return err, false, New(), ""
		}
		if !recordUpdated {
			logHandler.DatabaseLogger.Printf("[POSTUPDATE] No update required for %v Record: %v", TableName.String(), record.Key)
			return nil, false, New(), ""
		}
		if feedbackMessage == "" {
			feedbackMessage = "Post Update Processing"
		}
		// Return the TemplateStoreV3 record with the new profile and notes
		logHandler.DatabaseLogger.Printf("[POSTUPDATE] Processing complete for %v Record: %v", TableName.String(), record.Key)
		return nil, true, updatedRecord, feedbackMessage

	}
	logHandler.DatabaseLogger.Printf("[POSTUPDATE] No post-update function registered for %v Record: %v", TableName.String(), record.Key)
	return nil, false, New(), ""
}

// postDeleteProcessing runs any post-delete processing after a record is deleted.
func (record *TemplateStoreV3) postDeleteProcessing(ctx context.Context) error {
	if postDelete != nil {
		logHandler.DatabaseLogger.Printf("[POSTDELETE] Processing for %v Record: %v", TableName.String(), record.Key)
		err := postDelete(ctx, record)
		logHandler.DatabaseLogger.Printf("[POSTDELETE] Processing complete for %v Record: %v", TableName.String(), record.Key)
		return err
	}
	return nil
}
