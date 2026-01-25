// Data Access Object template
// Version: 0.5.0
// Updated on: 2026-01-24

package testentity

import (
	"context"
)

type creatorFunc func(TestEntity) (string, TestEntity, error)
type upgraderFunc func(TestEntity) (TestEntity, error)
type defaulterFunc func(*TestEntity) error
type validatorFunc func(*TestEntity) error
type preDeleteFunc func(*TestEntity) error
type postGetFunc func(*TestEntity) error
type clonerFunc func(ctx context.Context, source TestEntity) (TestEntity, error)
type duplicateCheckFunc func(*TestEntity) (bool, error)
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

// RegisterCreator registers a creator function for TestEntity.
func RegisterCreator(fn creatorFunc) {
	creator = fn
}

// RegisterUpgrader registers an upgrader function for TestEntity.
func RegisterUpgrader(fn upgraderFunc) {
	upgrader = fn
}

// RegisterDefaulter registers a defaulter function for TestEntity.
func RegisterDefaulter(fn defaulterFunc) {
	defaulter = fn
}

// RegisterValidator registers a validator function for TestEntity.
func RegisterValidator(fn validatorFunc) {
	validator = fn
}

// RegisterPreDelete registers a pre-delete function for TestEntity.
func RegisterPreDelete(fn preDeleteFunc) {
	preDelete = fn
}

// RegisterPostGet registers a post-get function for TestEntity.
func RegisterPostGet(fn postGetFunc) {
	postGet = fn
}

// RegisterCloner registers a cloner function for TestEntity.
func RegisterCloner(fn clonerFunc) {
	cloner = fn
}

// RegisterDuplicateCheck registers a duplicate check function for TestEntity.
func RegisterDuplicateCheck(fn duplicateCheckFunc) {
	duplicateCheck = fn
}

// RegisterWorker registers a worker function for TestEntity.
func RegisterWorker(fn workerFunc) {
	worker = fn
}

// upgradeProcessing performs any one-time upgrade or migration logic on the record.
func (record *TestEntity) upgradeProcessing() error {
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
func (record *TestEntity) defaultProcessing() error {
	if defaulter != nil {
		return defaulter(record)
	}
	return nil
}

// validationProcessing validates the record and returns an error if it is invalid.
func (record *TestEntity) validationProcessing() error {
	if validator != nil {
		return validator(record)
	}
	return nil
}

// postGetProcessing runs any post-load processing after a record is retrieved.
func (h *TestEntity) postGetProcessing() error {
	if postGet != nil {
		return postGet(h)
	}
	return nil
}

// preDeleteProcessing runs any checks or actions required before delete.
func (record *TestEntity) preDeleteProcessing() error {
	if preDelete != nil {
		return preDelete(record)
	}
	return nil
}

// templateClone contains the package's clone logic.
func templateClone(ctx context.Context, source TestEntity) (TestEntity, error) {
	if cloner != nil {
		return cloner(ctx, source)
	}
	return New(), nil
}

// // assertTestEntity asserts that an `any` returned by lower layers is a *TestEntity.
// func assertTestEntity(result any, field entities.Field, value any) (*TestEntity, error) {
// 	x, ok := result.(*TestEntity)
// 	if !ok {
// 		return nil, ce.ErrDAOAssertWrapper(tableName, field.String(), value,
// 			ce.ErrInvalidTypeWrapper(field.String(), fmt.Sprintf("%T", result), "*TestEntity"))
// 	}
// 	return x, nil
// }
