// Data Access Object for the TemplateStoreV3 table
// Template Version: 0.5.10 - 2026-01-26
// Generated 
// Date: 27/01/2026 & 10:17
// Who : matttownsend (orion)

package templateStoreV3

import (
	"context"

	"github.com/mt1976/frantic-amphora/dao/database"
	"github.com/mt1976/frantic-core/commonConfig"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

var activeDBConnection *database.DB
var databaseConnectionActive bool
var cfg *commonConfig.Settings

// Initialise opens the database connection for TemplateStoreV3 and optionally enables caching.
func Initialise(ctx context.Context, cached bool) {
	//logHandler.DatabaseLogger.Printf("Opening connection to %v", tableName)
	logHandler.TraceLogger.Printf("Initialising %v DAO Caching: %t", tableName, cached)

	clock := timing.Start(tableName, "Initialise", "Initialise")
	cfg = commonConfig.Get()
	_ = cfg

	activeDBConnection = database.Connect(TemplateStoreV3{}, database.WithVerbose(false), database.WithCaching(cached), database.WithCacheKey(Fields.Key), database.WithNameSpace("main"))
	databaseConnectionActive = true

	clock.Stop(1)
	//logHandler.DatabaseLogger.Printf("Opened connection to %v", tableName)
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
