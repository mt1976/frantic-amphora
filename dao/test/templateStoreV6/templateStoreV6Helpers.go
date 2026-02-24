// Data Access Object for the TripStore table
// Template Version: 0.5.24 - 2026-01-31
// Generated
// Date: 10/02/2026 & 13:16
// Who : matttownsend (orion)

package templateStoreV6

import (
	"context"
	"fmt"
	"os"

	"github.com/goforj/godump"
	"github.com/mt1976/frantic-amphora/dao"
	"github.com/mt1976/frantic-core/contextHandler"
	"github.com/mt1976/frantic-core/logHandler"
)

type (
	creatorFunc        func(ctx context.Context, source *TemplateStoreV6) (id string, skipPostCreate bool, record *TemplateStoreV6, err error)
	upgraderFunc       func(*TemplateStoreV6) error
	defaulterFunc      func(*TemplateStoreV6) error
	validatorFunc      func(*TemplateStoreV6) error
	preDeleteFunc      func(ctx context.Context, record *TemplateStoreV6) error
	postGetFunc        func(ctx context.Context, record *TemplateStoreV6) error
	clonerFunc         func(ctx context.Context, source *TemplateStoreV6) (clonedRecord *TemplateStoreV6, err error)
	duplicateCheckFunc func(*TemplateStoreV6) (found bool, err error)
	workerFunc         func(string, string)
	postCreateFunc     func(ctx context.Context, record *TemplateStoreV6) error
	postUpdateFunc     func(ctx context.Context, record *TemplateStoreV6) error
	postDeleteFunc     func(ctx context.Context, record *TemplateStoreV6) error
	postCloneFunc      func(ctx context.Context, record *TemplateStoreV6) error
	postDropFunc       func(ctx context.Context) error
	importerFunc       func(ctx context.Context, record *TemplateStoreV6) (err error)
	exporterFunc       func(ctx context.Context, record *TemplateStoreV6) (err error)
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
	importer       importerFunc
	exporter       exporterFunc
)

// RegisterCreator registers a creator function for TripStore.
func RegisterCreator(fn creatorFunc) {
	logHandler.Event.Printf("[REGISTER] Creator for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] Creator for %v (%v)", tableName, dao.GetFunctionName(fn))
	creator = fn
}

// RegisterPostCreate registers a post-create function for TripStore.
func RegisterPostCreate(fn postCreateFunc) {
	logHandler.Event.Printf("[REGISTER] PostCreate for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] PostCreate for %v (%v)", tableName, dao.GetFunctionName(fn))
	postCreate = fn
}

// RegisterPostUpdate registers a post-update function for TripStore.
func RegisterPostUpdate(fn postUpdateFunc) {
	logHandler.Event.Printf("[REGISTER] PostUpdate for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] PostUpdate for %v (%v)", tableName, dao.GetFunctionName(fn))
	postUpdate = fn
}

// RegisterPostDelete registers a post-delete function for TripStore.
func RegisterPostDelete(fn postDeleteFunc) {
	logHandler.Event.Printf("[REGISTER] PostDelete for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] PostDelete for %v (%v)", tableName, dao.GetFunctionName(fn))
	postDelete = fn
}

// RegisterPostClone registers a post-clone function for TripStore.
func RegisterPostClone(fn postCloneFunc) {
	logHandler.Event.Printf("[REGISTER] PostClone for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] PostClone for %v (%v)", tableName, dao.GetFunctionName(fn))
	postClone = fn
}

// RegisterPostDrop registers a post-drop function for TripStore.
func RegisterPostDrop(fn postDropFunc) {
	logHandler.Event.Printf("[REGISTER] PostDrop for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] PostDrop for %v (%v)", tableName, dao.GetFunctionName(fn))
	postDrop = fn
}

// RegisterPostClearDown registers a post-clear-down function for TripStore.
func RegisterPostClearDown(fn postDropFunc) {
	logHandler.Event.Printf("[REGISTER] PostClearDown for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] PostClearDown for %v (%v)", tableName, dao.GetFunctionName(fn))
	postClearDown = fn
}

// RegisterUpgrader registers an upgrader function for TripStore.
func RegisterUpgrader(fn upgraderFunc) {
	logHandler.Event.Printf("[REGISTER] Upgrader for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] Upgrader for %v (%v)", tableName, dao.GetFunctionName(fn))
	upgrader = fn
}

// RegisterDefaulter registers a defaulter function for TripStore.
func RegisterDefaulter(fn defaulterFunc) {
	logHandler.Event.Printf("[REGISTER] Defaulter for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] Defaulter for %v (%v)", tableName, dao.GetFunctionName(fn))
	defaulter = fn
}

// RegisterValidator registers a validator function for TripStore.
func RegisterValidator(fn validatorFunc) {
	logHandler.Event.Printf("[REGISTER] Validator for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] Validator for %v (%v)", tableName, dao.GetFunctionName(fn))
	validator = fn
}

// RegisterPreDelete registers a pre-delete function for TripStore.
func RegisterPreDelete(fn preDeleteFunc) {
	logHandler.Event.Printf("[REGISTER] PreDelete for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] PreDelete for %v (%v)", tableName, dao.GetFunctionName(fn))
	preDelete = fn
}

// RegisterPostGet registers a post-get function for TripStore.
func RegisterPostGet(fn postGetFunc) {
	logHandler.Event.Printf("[REGISTER] PostGet for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] PostGet for %v (%v)", tableName, dao.GetFunctionName(fn))
	postGet = fn
}

// RegisterCloner registers a cloner function for TripStore.
func RegisterCloner(fn clonerFunc) {
	logHandler.Event.Printf("[REGISTER] Cloner for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] Cloner for %v (%v)", tableName, dao.GetFunctionName(fn))
	cloner = fn
}

// RegisterDuplicateCheck registers a duplicate check function for TripStore.
func RegisterDuplicateCheck(fn duplicateCheckFunc) {
	logHandler.Event.Printf("[REGISTER] DuplicateCheck for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] DuplicateCheck for %v (%v)", tableName, dao.GetFunctionName(fn))
	duplicateCheck = fn
}

// RegisterWorker registers a worker function for TripStore.
func RegisterWorker(fn workerFunc) {
	logHandler.Event.Printf("[REGISTER] Worker for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] Worker for %v (%v)", tableName, dao.GetFunctionName(fn))
	worker = fn
}

// RegisterImporter registers an importer function for TripStore.
func RegisterImporter(fn importerFunc) {
	logHandler.Event.Printf("[REGISTER] Importer for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] Importer for %v (%v)", tableName, dao.GetFunctionName(fn))
	importer = fn
}

// RegisterExporter registers an exporter function for TripStore.
func RegisterExporter(fn exporterFunc) {
	logHandler.Event.Printf("[REGISTER] Exporter for %v (%v)", tableName, dao.GetFunctionName(fn))
	logHandler.Database.Printf("[REGISTER] Exporter for %v (%v)", tableName, dao.GetFunctionName(fn))
	exporter = fn
}

// upgradeProcessing performs any one-time upgrade or migration logic on the record.
func (record *TemplateStoreV6) upgradeProcessing() error {
	if upgrader != nil {
		logHandler.Database.Printf("[UPGRADE] record %v of %v", record.Key, TableName.String())
		err := upgrader(record)
		if err != nil {
			return err
		}
		//*record = *updatedRecord
		logHandler.Database.Printf("[UPGRADE] Upgrade complete for record %v of %v", record.Key, TableName.String())
	}
	return nil
}

// defaultProcessing applies any default values prior to validation and persistence.
func (record *TemplateStoreV6) defaultProcessing() error {
	if defaulter != nil {
		logHandler.Database.Printf("[DEFAULT] Applying defaults to record %v of %v", record.Key, TableName.String())
		err := defaulter(record)
		logHandler.Database.Printf("[DEFAULT] Defaults applied to record %v of %v", record.Key, TableName.String())
		return err
	}
	return nil
}

// validationProcessing validates the record and returns an error if it is invalid.
func (record *TemplateStoreV6) validationProcessing() error {
	if validator != nil {
		logHandler.Database.Printf("[VALIDATE] Validating record %v of %v", record.Key, TableName.String())
		err := validator(record)
		logHandler.Database.Printf("[VALIDATE] Record %v of %v is valid", record.Key, TableName.String())
		return err
	}
	return nil
}

// postGetProcessing runs any post-load processing after a record is retrieved.
func (h *TemplateStoreV6) postGetProcessing(ctx context.Context) error {
	if postGet != nil {
		logHandler.Database.Printf("[POSTGET] Processing for %v Record: %v", TableName.String(), h.Key)
		err := postGet(ctx, h)
		logHandler.Database.Printf("[POSTGET] Processing complete for %v Record: %v", TableName.String(), h.Key)
		return err
	}
	return nil
}

// preDeleteProcessing runs any checks or actions required before delete.
func (record *TemplateStoreV6) preDeleteProcessing(ctx context.Context) error {
	if preDelete != nil {
		logHandler.Database.Printf("[PREDELETE] Processing for %v Record: %v", TableName.String(), record.Key)
		err := preDelete(ctx, record)
		logHandler.Database.Printf("[PREDELETE] Processing complete for %v Record: %v", TableName.String(), record.Key)
		return err
	}
	return nil
}

// templateClone contains the package's clone logic.
func templateClone(ctx context.Context, source *TemplateStoreV6) (clone *TemplateStoreV6, err error) {
	if cloner != nil {
		logHandler.Database.Printf("[CLONE] Cloning %v Record: %v", TableName.String(), source.Key)
		return cloner(ctx, source)
	}
	return New(), nil
}

// PostCreate runs any post-create processing after a record is created.
func (record *TemplateStoreV6) postCreateProcessing(ctx context.Context) error {
	if postCreate != nil {
		logHandler.Trace.Printf("[POSTCREATE] Processing for %v Record: %v", TableName.String(), record.Key)
		// logHandler.LockLogger.Printf("[POSTCREATE] LOCKING RECORD: %v", record.Raw)

		// Get the workerPool from the context and run the post-create processing in a worker to avoid holding locks during potentially long-running processing.
		// WARNING :: Doesnt seem like context is at this point is the correct context (the one we get doenst have the worker pool in it) - need to investigate further and ensure we have the correct context here with the worker pool included. This is important to ensure we can run post-create processing in a worker and avoid holding locks during long-running processing which can lead to contention and performance issues.
		// godump.Dump(ctx, "Context at start of post-create processing for record "+record.Key)

		workerPool, err := contextHandler.GetWorkerPool(ctx)
		if err != nil {
			logHandler.Error.Printf("Error retrieving worker pool from context for post-create processing for %v record %v: %v", TableName.String(), record.Key, err.Error())
			// Proceed with synchronous processing if worker pool is not available in context, but log the error. This ensures that post-create processing still occurs even if there is an issue with the worker pool in the context, while also providing visibility into the issue through logging.
			logHandler.Trace.Printf("[POSTCREATE] Worker pool not found in context, running post-create processing for %v Record: %v synchronously", TableName.String(), record.Key)
			os.Exit(0)
		}
		if workerPool != nil {
			logHandler.Trace.Printf("[POSTCREATE] Running post-create processing for %v Record: %v in worker pool Lock:%v", TableName.String(), record.Key, godump.DumpStr(record.Lock))
			task := workerPool.SubmitErr(func() error {
				// Lock record for the duration of the post-create processing to prevent concurrent updates to the same record. This is necessary because the record may be updated during post-processing and we want to prevent concurrent updates to the same record during that time.
				logHandler.Lock.Printf("[%v,%v] Locked - POSTCREATE", TableName.String(), record.Raw)
				logHandler.Trace.Printf("CONTEXT=[%v]", godump.DumpStr(ctx))
				record.Lock.Lock()
				err := postCreate(ctx, record)
				if err != nil {
					logHandler.Error.Printf("Error during postCreateProcessing for %v Record: %v Error: %v", TableName.String(), record.Key, err.Error())
					record.Lock.Unlock()
					logHandler.Unlock.Printf("[%v,%v] Unlocked - POSTCREATE with error %v", TableName.String(), record.Raw, err.Error())
					return err
				}
				record.Lock.Unlock()
				logHandler.Unlock.Printf("[%v,%v] Unlocked - POSTCREATE", TableName.String(), record.Raw)
				logHandler.Trace.Printf("[POSTCREATE] Post-create processing complete for %v Record: %v", TableName.String(), record.Key)
				return nil
			})
			err := task.Wait() // Wait for the task to complete before returning to ensure post-create processing is complete before any subsequent operations that may depend on it. This is important to maintain data integrity and consistency, especially if the post-create processing involves updates to the record or related records.
			if err != nil {
				logHandler.Error.Printf("Error in post-create worker task for %v Record: %v ", TableName.String(), record.Key)
				return fmt.Errorf("error in post-create worker task for %v Record: %v %+v", TableName.String(), record.Key, godump.DumpStr(err))
			}
			return nil
		} else {
			logHandler.Trace.Printf("[POSTCREATE] No worker pool found in context, running post-create processing for %v Record: %v synchronously", TableName.String(), record.Key)
		}
	}
	return nil
}

func (record *TemplateStoreV6) postUpdateProcessing(ctx context.Context) error {
	if postUpdate != nil {
		logHandler.Trace.Printf("[POSTUPDATE] Processing for %v Record: %v (%v)", TableName.String(), record.Key, godump.DumpStr(ctx))
		workerPool, err := contextHandler.GetWorkerPool(ctx)
		if err != nil {
			logHandler.Error.Printf("Error retrieving worker pool from context for post-create processing for %v record %v: %v", TableName.String(), record.Key, err.Error())
			// Proceed with synchronous processing if worker pool is not available in context, but log the error. This ensures that post-create processing still occurs even if there is an issue with the worker pool in the context, while also providing visibility into the issue through logging.
			logHandler.Trace.Printf("[POSTCREATE] Worker pool not found in context, running post-create processing for %v Record: %v synchronously", TableName.String(), record.Key)
			os.Exit(0)
		}
		if workerPool != nil {

			logHandler.Trace.Printf("[POSTUPDATE] Running post-update processing for %v Record: %v in worker pool", TableName.String(), record.Key)
			task := workerPool.SubmitErr(func() error {
				logHandler.Trace.Printf("[POSTUPDATE] Starting post-update processing for %v Record: %v Lock: %v", TableName.String(), record.Key, godump.DumpStr(record.Lock))
				// Lock record for the duration of the post-update processing to prevent concurrent updates to the same record. This is necessary because the record may be updated during post-processing and we want to prevent
				record.Lock.Lock()
				logHandler.Lock.Printf("[%v,%v] Locked - POSTUPDATE", TableName.String(), record.Raw)
				logHandler.Trace.Printf("CONTEXT=[%v]", godump.DumpStr(ctx))
				err := postUpdate(ctx, record)
				if err != nil {
					record.Lock.Unlock()
					logHandler.Unlock.Printf("[%v,%v] Unlocked - POSTUPDATE with error %v", TableName.String(), record.Raw, err.Error())
					logHandler.Error.Printf("error during postUpdateProcessing task for %v Record: %v Error: %v", TableName.String(), record.Key, err.Error())
					// panic(err)
					return err
				}
				record.Lock.Unlock()
				logHandler.Unlock.Printf("[%v,%v] Unlocked - POSTUPDATE", TableName.String(), record.Raw)
				logHandler.Database.Printf("[POSTUPDATE] Post-update processing complete for %v Record: %v", TableName.String(), record.Key)
				return nil
			})
			err := task.Wait() // Wait for the task to complete before returning to ensure post-update processing is complete before any subsequent operations that may depend on it. This is important to maintain data integrity and consistency, especially if the post-update processing involves updates to the record or related records.
			if err != nil {
				logHandler.Error.Printf("Error in post-update worker task for %v Record: %v ", TableName.String(), record.Key)
				return fmt.Errorf("error in post-update worker task for %v Record: %v ", TableName.String(), record.Key)
			}
			return nil
		} else {
			logHandler.Database.Printf("[POSTUPDATE] No worker pool found in context, running post-update processing for %v Record: %v synchronously", TableName.String(), record.Key)
		}

		logHandler.Database.Printf("[POSTUPDATE] Processing complete for %v Record: %v", TableName.String(), record.Key)
		return nil

	}
	logHandler.Database.Printf("[POSTUPDATE] No post-update function registered for %v Record: %v", TableName.String(), record.Key)
	return nil
}

// postDeleteProcessing runs any post-delete processing after a record is deleted.
func (record *TemplateStoreV6) postDeleteProcessing(ctx context.Context) error {
	if postDelete != nil {
		logHandler.Database.Printf("[POSTDELETE] Processing for %v Record: %v", TableName.String(), record.Key)
		err := postDelete(ctx, record)
		logHandler.Database.Printf("[POSTDELETE] Processing complete for %v Record: %v", TableName.String(), record.Key)
		return err
	}
	return nil
}
