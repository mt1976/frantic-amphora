// Package dao provides data access primitives and lifecycle helpers, including
// initialization, maintenance, lookups, and auditing support.
package dao

import (
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/asdine/storm/v3"
	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-core/commonConfig"
	ce "github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

var name = "DAO"
var DBVersion = 1
var DB *storm.DB
var DBName string = "default"

// Initialise sets up the DAO package with the provided configuration.
func Initialise(cfg *commonConfig.Settings) error {
	clock := timing.Start(name, "Initialise", "")
	logHandler.InfoLogger.Printf("[%v] Initialising...", strings.ToUpper(name))

	DBVersion = cfg.GetDatabase_Version()
	DBName = cfg.GetDatabase_Name()

	logHandler.InfoLogger.Printf("[%v] Initialised", strings.ToUpper(name))
	clock.Stop(1)
	return nil
}

// GetDBNameFromPath extracts the database name from a given file path.
func GetDBNameFromPath(t string) string {
	dbName := t
	// split dbName on "/"
	dbNameArr := strings.Split(dbName, string(os.PathSeparator))
	noparts := len(dbNameArr)
	dbName = dbNameArr[noparts-1]
	logHandler.TraceLogger.Printf("dbName: %v\n", dbName)
	return dbName
}

// CheckDAOReadyState verifies that the DAO is initialized before performing any operations.
func CheckDAOReadyState(table string, action audit.Action, isDaoReady bool) {
	if !isDaoReady {
		err := ce.ErrDAONotInitialisedWrapper(table, action.Description())
		logHandler.ErrorLogger.Panic(err)
	}
}

// GetFunctionName returns the name of a function passed as an interface.
func GetFunctionName(i any) string {
	// Function to get the clean name of a function stored in an interface
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
