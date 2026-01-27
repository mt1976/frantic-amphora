// Data Access Object for the TemplateStoreV3 table
// Template Version: 0.5.10 - 2026-01-26
// Generated 
// Date: 27/01/2026 & 12:54
// Who : matttownsend (orion)

package templateStoreV3

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/mt1976/frantic-amphora/dao"
	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/importExportHelper"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

// ExportRecordToJSON exports the record as a JSON file.
func (record *TemplateStoreV3) ExportRecordToJSON(name string) {
	ID := reflect.ValueOf(*record).FieldByName(Fields.ID.String())
	clock := timing.Start(tableName, "Export", fmt.Sprintf("%v", ID))

	err := importExportHelper.ExportJSON(name, []TemplateStoreV3{*record}, Fields.ID)
	if err != nil {
		logHandler.ExportLogger.Panicf("error exporting %v record %v: %v", tableName, ID, err.Error())
	}

	clock.Stop(1)
}

// ExportAllToJSON exports all records as JSON files.
func ExportAllToJSON(message string) {
	dao.CheckDAOReadyState(tableName, audit.EXPORT, databaseConnectionActive)

	clock := timing.Start(tableName, "Export", "ALL")
	recordList, _ := GetAll()
	if len(recordList) == 0 {
		logHandler.WarningLogger.Printf("[%v] %v data not found", tableName, tableName)
		clock.Stop(0)
		return
	}

	err := importExportHelper.ExportJSON(message, recordList, Fields.ID)
	if err != nil {
		logHandler.ExportLogger.Panicf("error exporting all %v's: %v", tableName, err.Error())
	}
	clock.Stop(len(recordList))
}

// ExportRecordToCSV exports the record as a CSV file.
func (record *TemplateStoreV3) ExportRecordToCSV(name string) error {
	ID := reflect.ValueOf(*record).FieldByName(Fields.ID.String())
	clock := timing.Start(tableName, "Export", fmt.Sprintf("%v", ID))

	err := importExportHelper.ExportCSV(name, []TemplateStoreV3{*record}, Fields.ID)
	if err != nil {
		logHandler.ExportLogger.Printf("Error exporting %v record %v: %v", tableName, ID, err.Error())
		clock.Stop(0)
		return err
	}

	clock.Stop(1)
	return nil
}

// ExportAllToCSV exports all records as a CSV file.
func ExportAllToCSV(msg string) error {
	exportListData, err := GetAll()
	if err != nil {
		logHandler.ExportLogger.Panicf("error Getting all %v's: %v", tableName, err.Error())
	}
	return importExportHelper.ExportCSV(msg, exportListData, Fields.ID)
}

// ImportAllFromCSV imports records for this table from a CSV file.
func ImportAllFromCSV() error {
	return importExportHelper.ImportCSV(tableName, &TemplateStoreV3{}, templateImportProcessor)
}

// templateImportProcessor is called for each CSV row during import.
func templateImportProcessor(inOriginal **TemplateStoreV3) (string, error) {
	importedData := **inOriginal
	stringField1 := strconv.Itoa(importedData.ID)

	_, err := Create(context.TODO(), importedData)
	if err != nil {
		logHandler.ImportLogger.Panicf("Error importing %v: %v", tableName, err.Error())
		return stringField1, err
	}

	return stringField1, nil
}
