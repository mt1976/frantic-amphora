// Data Access Object for the TemplateStoreV3 table
// Template Version: 0.5.10 - 2026-01-26
// Generated 
// Date: 27/01/2026 & 10:17
// Who : matttownsend (orion)

package templateStoreV3

import (
	"github.com/goforj/godump"
	"github.com/mt1976/frantic-core/logHandler"
)

// Spew outputs the contents of the record to the Trace log.
func (record *TemplateStoreV3) Spew() {
	logHandler.TraceLogger.Printf("[%v] Record=[%+v]", tableName, godump.DumpStr(record))
}
