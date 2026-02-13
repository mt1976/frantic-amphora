// Data Access Object for the TripStore table
// Template Version: 0.5.24 - 2026-01-31
// Generated
// Date: 10/02/2026 & 13:16
// Who : matttownsend (orion)

package templateStoreV6

import (
	"context"

	"github.com/mt1976/frantic-amphora/dao/database"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

var (
	activeDBConnection       *database.DB
	databaseConnectionActive bool
	cfg                      *commonConfig.Settings
)

// Initialise opens the database connection for TripStore and optionally enables caching.
func Initialise(ctx context.Context, cached bool) {
	// logHandler.DatabaseLogger.Printf("Opening connection to %v", tableName)
	logHandler.Trace.Printf("Initialising %v DAO Caching: %t", tableName, cached)

	clock := timing.Start(tableName, "Initialise", "Initialise")
	cfg = commonConfig.Get()
	_ = cfg

	activeDBConnection = database.Connect(TemplateStoreV6{}, database.WithVerbose(false), database.WithCaching(cached), database.WithCacheKey(Fields.Key), database.WithNameSpace("main"))
	databaseConnectionActive = true

	clock.Stop(1)
	// logHandler.DatabaseLogger.Printf("Opened connection to %v", tableName)
}

// IsInitialised reports whether the DAO has an active database connection.
func IsInitialised() bool {
	return databaseConnectionActive
}

// GetDatabaseConnections returns a function that supplies the database connections used by this DAO.
func GetDatabaseConnections() func() ([]*database.DB, error) {
	return func() ([]*database.DB, error) {
		return []*database.DB{activeDBConnection}, nil
	}
}
