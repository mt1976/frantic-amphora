// Data Access Object template
// Version: 0.5.0
// Updated on: 2026-01-24

package templateStoreV2

import (
	"context"
	"fmt"

	"github.com/mt1976/frantic-amphora/dao/entities"
	ce "github.com/mt1976/frantic-core/commonErrors"
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

var creator creatorFunc
var upgrader upgraderFunc
var defaulter defaulterFunc
var validator validatorFunc
var preDelete preDeleteFunc
var postGet postGetFunc
var cloner clonerFunc
var duplicateCheck duplicateCheckFunc
var worker workerFunc

// RegisterCreator registers a creator function for TemplateStore.
func RegisterCreator(fn creatorFunc) {
	creator = fn
}

// RegisterUpgrader registers an upgrader function for TemplateStore.
func RegisterUpgrader(fn upgraderFunc) {
	upgrader = fn
}

// RegisterDefaulter registers a defaulter function for TemplateStore.
func RegisterDefaulter(fn defaulterFunc) {
	defaulter = fn
}

// RegisterValidator registers a validator function for TemplateStore.
func RegisterValidator(fn validatorFunc) {
	validator = fn
}

// RegisterPreDelete registers a pre-delete function for TemplateStore.
func RegisterPreDelete(fn preDeleteFunc) {
	preDelete = fn
}

// RegisterPostGet registers a post-get function for TemplateStore.
func RegisterPostGet(fn postGetFunc) {
	postGet = fn
}

// RegisterCloner registers a cloner function for TemplateStore.
func RegisterCloner(fn clonerFunc) {
	cloner = fn
}

// RegisterDuplicateCheck registers a duplicate check function for TemplateStore.
func RegisterDuplicateCheck(fn duplicateCheckFunc) {
	duplicateCheck = fn
}

// RegisterWorker registers a worker function for TemplateStore.
func RegisterWorker(fn workerFunc) {
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

// assertTemplateStore asserts that an `any` returned by lower layers is a *TemplateStore.
func assertTemplateStore(result any, field entities.Field, value any) (*TemplateStore, error) {
	x, ok := result.(*TemplateStore)
	if !ok {
		return nil, ce.ErrDAOAssertWrapper(tableName, field.String(), value,
			ce.ErrInvalidTypeWrapper(field.String(), fmt.Sprintf("%T", result), "*TemplateStore"))
	}
	return x, nil
}
