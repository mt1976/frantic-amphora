// Data Access Object for the TripStore table
// Template Version: 0.5.24 - 2026-01-31
// Generated
// Date: 10/02/2026 & 13:16
// Who : matttownsend (orion)

package templateStoreV6

import (
	"github.com/goforj/godump"
	"github.com/mt1976/frantic-core/logHandler"
)

// Spew outputs the contents of the record to the Trace log.
func (record *TemplateStoreV6) Spew() {
	logHandler.TraceLogger.Printf("[%v] Record=[%+v]", tableName, godump.DumpStr(record))
}
