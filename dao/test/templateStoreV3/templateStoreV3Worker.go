// Data Access Object for the TemplateStoreV3 table
// Template Version: 0.5.10 - 2026-01-26
// Generated 
// Date: 27/01/2026 & 10:17
// Who : matttownsend (orion)

package templateStoreV3

import (
	"github.com/mt1976/frantic-amphora/dao/database"
	"github.com/mt1976/frantic-amphora/jobs"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

// Worker is a job that is scheduled to run at a predefined interval.
func Worker(j jobs.Job, db *database.DB) {
	clock := timing.Start(jobs.CodedName(j), "Initialise", j.Description())
	oldDB := activeDBConnection
	dbSwitched := false

	if db != nil {
		if activeDBConnection.Name != db.Name {
			logHandler.EventLogger.Printf("Switching to %v.db", db.Name)
			activeDBConnection = db
			dbSwitched = true
		}
	}

	if worker != nil {
		worker(jobs.CodedName(j), j.Description())
	}

	if dbSwitched {
		logHandler.EventLogger.Printf("Switching back to %v.db from %v.db", oldDB.Name, activeDBConnection.Name)
		activeDBConnection = oldDB
	}
	clock.Stop(1)
}
