// Data Access Object template
// Version: 0.5.0
// Updated on: 2026-01-24

package templateStoreV2

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

// ExportRecordToJSON exports the TemplateStore record as a JSON file.
func (record *TemplateStore) ExportRecordToJSON(name string) {
	ID := reflect.ValueOf(*record).FieldByName(Fields.ID.String())
	clock := timing.Start(tableName, "Export", fmt.Sprintf("%v", ID))

	err := importExportHelper.ExportJSON(name, []TemplateStore{*record}, Fields.ID)
	if err != nil {
		logHandler.ExportLogger.Panicf("error exporting %v record %v: %v", tableName, ID, err.Error())
	}

	clock.Stop(1)
}

// ExportAllToJSON exports all TemplateStore records as JSON files.
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

// ExportRecordToCSV exports the TemplateStore record as a CSV file.
func (record *TemplateStore) ExportRecordToCSV(name string) error {
	ID := reflect.ValueOf(*record).FieldByName(Fields.ID.String())
	clock := timing.Start(tableName, "Export", fmt.Sprintf("%v", ID))

	err := importExportHelper.ExportCSV(name, []TemplateStore{*record}, Fields.ID)
	if err != nil {
		logHandler.ExportLogger.Printf("Error exporting %v record %v: %v", tableName, ID, err.Error())
		clock.Stop(0)
		return err
	}

	clock.Stop(1)
	return nil
}

// ExportAllToCSV exports all TemplateStore records as a CSV file.
func ExportAllToCSV(msg string) error {
	exportListData, err := GetAll()
	if err != nil {
		logHandler.ExportLogger.Panicf("error Getting all %v's: %v", tableName, err.Error())
	}
	return importExportHelper.ExportCSV(msg, exportListData, Fields.ID)
}

// ImportAllFromCSV imports records for this table from a CSV file.
func ImportAllFromCSV() error {
	return importExportHelper.ImportCSV(tableName, &TemplateStore{}, templateImportProcessor)
}

// templateImportProcessor is called for each CSV row during import.
func templateImportProcessor(inOriginal **TemplateStore) (string, error) {
	importedData := **inOriginal
	stringField1 := strconv.Itoa(importedData.ID)

	_, err := Create(context.TODO(), importedData)
	if err != nil {
		logHandler.ImportLogger.Panicf("Error importing %v: %v", tableName, err.Error())
		return stringField1, err
	}

	return stringField1, nil
}
