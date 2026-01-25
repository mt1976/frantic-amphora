// Data Access Object template
// Version: 0.5.0
// Updated on: 2026-01-24

package templateStoreV2

import (
	"context"

	"github.com/mt1976/frantic-amphora/dao"
	"github.com/mt1976/frantic-core/logHandler"
)

type creatorFunc func(TemplateStore) (string, TemplateStore, error)
type upgraderFunc func(TemplateStore) (TemplateStore, error)
type defaulterFunc func(*TemplateStore) error
type validatorFunc func(*TemplateStore) error
type preDeleteFunc func(*TemplateStore) error
type postGetFunc func(*TemplateStore) error
type clonerFunc func(ctx context.Context, source TemplateStore) (TemplateStore, error)
type duplicateCheckFunc func(*TemplateStore) (bool, error)
type workerFunc func(string, string)
type postCreateFunc func(ctx context.Context, record *TemplateStore) error
type postUpdateFunc func(ctx context.Context, record *TemplateStore) error
type postDeleteFunc func(ctx context.Context, record *TemplateStore) error
type postCloneFunc func(ctx context.Context, record *TemplateStore) error
type postDropFunc func(ctx context.Context) error

var creator creatorFunc
var upgrader upgraderFunc
var defaulter defaulterFunc
var validator validatorFunc
var preDelete preDeleteFunc
var postGet postGetFunc
var cloner clonerFunc
var duplicateCheck duplicateCheckFunc
var worker workerFunc
var postCreate postCreateFunc
var postUpdate postUpdateFunc
var postDelete postDeleteFunc
var postClone postCloneFunc
var postDrop postDropFunc
var postClearDown postDropFunc

// RegisterCreator registers a creator function for TemplateStore.
func RegisterCreator(fn creatorFunc) {
	logHandler.DatabaseLogger.Printf("Registering Creator for %v (%v)", tableName, dao.GetFunctionName(fn))
	creator = fn
}

// RegisterPostCreate registers a post-create function for TemplateStore.
func RegisterPostCreate(fn postCreateFunc) {
	logHandler.DatabaseLogger.Printf("Registering PostCreate for %v (%v)", tableName, dao.GetFunctionName(fn))
	postCreate = fn
}

// RegisterPostUpdate registers a post-update function for TemplateStore.
func RegisterPostUpdate(fn postUpdateFunc) {
	logHandler.DatabaseLogger.Printf("Registering PostUpdate for %v (%v)", tableName, dao.GetFunctionName(fn))
	postUpdate = fn
}

// RegisterPostDelete registers a post-delete function for TemplateStore.
func RegisterPostDelete(fn postDeleteFunc) {
	logHandler.DatabaseLogger.Printf("Registering PostDelete for %v (%v)", tableName, dao.GetFunctionName(fn))
	postDelete = fn
}

// RegisterPostClone registers a post-clone function for TemplateStore.
func RegisterPostClone(fn postCloneFunc) {
	logHandler.DatabaseLogger.Printf("Registering PostClone for %v (%v)", tableName, dao.GetFunctionName(fn))
	postClone = fn
}

// RegisterPostDrop registers a post-drop function for TemplateStore.
func RegisterPostDrop(fn postDropFunc) {
	logHandler.DatabaseLogger.Printf("Registering PostDrop for %v (%v)", tableName, dao.GetFunctionName(fn))
	postDrop = fn
}

// RegisterPostClearDown registers a post-clear-down function for TemplateStore.
func RegisterPostClearDown(fn postDropFunc) {
	logHandler.DatabaseLogger.Printf("Registering PostClearDown for %v (%v)", tableName, dao.GetFunctionName(fn))
	postClearDown = fn
}

// RegisterUpgrader registers an upgrader function for TemplateStore.
func RegisterUpgrader(fn upgraderFunc) {
	logHandler.DatabaseLogger.Printf("Registering Upgrader for %v (%v)", tableName, dao.GetFunctionName(fn))
	upgrader = fn
}

// RegisterDefaulter registers a defaulter function for TemplateStore.
func RegisterDefaulter(fn defaulterFunc) {
	logHandler.DatabaseLogger.Printf("Registering Defaulter for %v (%v)", tableName, dao.GetFunctionName(fn))
	defaulter = fn
}

// RegisterValidator registers a validator function for TemplateStore.
func RegisterValidator(fn validatorFunc) {
	logHandler.DatabaseLogger.Printf("Registering Validator for %v (%v)", tableName, dao.GetFunctionName(fn))
	validator = fn
}

// RegisterPreDelete registers a pre-delete function for TemplateStore.
func RegisterPreDelete(fn preDeleteFunc) {
	logHandler.DatabaseLogger.Printf("Registering PreDelete for %v (%v)", tableName, dao.GetFunctionName(fn))
	preDelete = fn
}

// RegisterPostGet registers a post-get function for TemplateStore.
func RegisterPostGet(fn postGetFunc) {
	logHandler.DatabaseLogger.Printf("Registering PostGet for %v (%v)", tableName, dao.GetFunctionName(fn))
	postGet = fn
}

// RegisterCloner registers a cloner function for TemplateStore.
func RegisterCloner(fn clonerFunc) {
	logHandler.DatabaseLogger.Printf("Registering Cloner for %v (%v)", tableName, dao.GetFunctionName(fn))
	cloner = fn
}

// RegisterDuplicateCheck registers a duplicate check function for TemplateStore.
func RegisterDuplicateCheck(fn duplicateCheckFunc) {
	logHandler.DatabaseLogger.Printf("Registering DuplicateCheck for %v (%v)", tableName, dao.GetFunctionName(fn))
	duplicateCheck = fn
}

// RegisterWorker registers a worker function for TemplateStore.
func RegisterWorker(fn workerFunc) {
	logHandler.DatabaseLogger.Printf("Registering Worker for %v (%v)", tableName, dao.GetFunctionName(fn))
	worker = fn
}

// upgradeProcessing performs any one-time upgrade or migration logic on the record.
func (record *TemplateStore) upgradeProcessing() error {
	if upgrader != nil {
		updatedRecord, err := upgrader(*record)
		if err != nil {
			return err
		}
		*record = updatedRecord
	}
	return nil
}

// defaultProcessing applies any default values prior to validation and persistence.
func (record *TemplateStore) defaultProcessing() error {
	if defaulter != nil {
		return defaulter(record)
	}
	return nil
}

// validationProcessing validates the record and returns an error if it is invalid.
func (record *TemplateStore) validationProcessing() error {
	if validator != nil {
		return validator(record)
	}
	return nil
}

// postGetProcessing runs any post-load processing after a record is retrieved.
func (h *TemplateStore) postGetProcessing() error {
	if postGet != nil {
		return postGet(h)
	}
	return nil
}

// preDeleteProcessing runs any checks or actions required before delete.
func (record *TemplateStore) preDeleteProcessing() error {
	if preDelete != nil {
		return preDelete(record)
	}
	return nil
}

// templateClone contains the package's clone logic.
func templateClone(ctx context.Context, source TemplateStore) (TemplateStore, error) {
	if cloner != nil {
		return cloner(ctx, source)
	}
	return New(), nil
}

// // assertTemplateStore asserts that an `any` returned by lower layers is a *TemplateStore.
// func assertTemplateStore(result any, field entities.Field, value any) (*TemplateStore, error) {
// 	x, ok := result.(*TemplateStore)
// 	if !ok {
// 		return nil, ce.ErrDAOAssertWrapper(tableName, field.String(), value,
// 			ce.ErrInvalidTypeWrapper(field.String(), fmt.Sprintf("%T", result), "*TemplateStore"))
// 	}
// 	return x, nil
// }

// PostCreate runs any post-create processing after a record is created.
func (record *TemplateStore) postCreateProcessing() error {
	if postCreate != nil {
		return postCreate(context.Background(), record)
	}
	return nil
}

func (record *TemplateStore) postUpdateProcessing() error {
	if postUpdate != nil {
		return postUpdate(context.Background(), record)
	}
	return nil
}

// postDeleteProcessing runs any post-delete processing after a record is deleted.
func (record *TemplateStore) postDeleteProcessing() error {
	if postDelete != nil {
		return postDelete(context.Background(), record)
	}
	return nil
}

// // postCloneProcessing runs any post-clone processing after a record is cloned.
// func (record *TemplateStore) postCloneProcessing() error {
// 	if postClone != nil {
// 		return postClone(context.Background(), record)
// 	}
// 	return nil
// }

// // postDropProcessing runs any post-drop processing after the table is dropped.
// func postDropProcessing() error {
// 	if postDrop != nil {
// 		return postDrop(context.Background())
// 	}
// 	return nil
// }
