// Data Access Object for the TripStore table
// Template Version: 0.6.00 - 2026-02-13
// Generated
// Date: 10/02/2026 & 13:16
// Who : matttownsend (orion)

package templateStoreV6

import (
	"context"

	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-core/logHandler"
)

// CacheHydrator returns the cache hydrator function for this table.
func CacheHydrator(ctx context.Context) func() ([]any, error) {
	_ = ctx
	return func() ([]any, error) {
		records, err := GetAll()
		if err != nil {
			return nil, err
		}
		result := make([]any, len(records))
		for i := range records {
			result[i] = records[i]
		}
		return result, nil
	}
}

// CacheSynchroniser returns the cache synchroniser function for this table.
func CacheSynchroniser(ctx context.Context) func(any) error {
	logHandler.InfoLogger.Printf("Defining Sync function for %v", tableName)
	return func(data any) error {
		switch rec := data.(type) {
		case TemplateStoreV6:
			logHandler.CacheLogger.Printf("Sync cache for %v Key: %v", tableName, rec.Key)
			return rec.UpdateWithAction(ctx, audit.SYNC, "Cache Sync Update")
		case *TemplateStoreV6:
			if rec == nil {
				return nil
			}
			logHandler.CacheLogger.Printf("Sync cache for %v Key: %v", tableName, rec.Key)
			return rec.UpdateWithAction(ctx, audit.SYNC, "Cache Sync Update")
		default:
			logHandler.WarningLogger.Printf("Sync cache for %v received unexpected type %T", tableName, data)
			return nil
		}
	}
}
