package importExportHelper

import (
	"encoding/csv"
	"fmt"
	"io"
	"os/user"
	"reflect"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/goforj/godump"
	"github.com/mt1976/frantic-amphora/dao/entities"
	"github.com/mt1976/frantic-core/application"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/paths"
	"github.com/mt1976/frantic-core/timing"
)

func ExportCSV[T any](exportName string, exportList []*T, idField entities.Field, mode mode) error {
	clock := timing.Start(exportName, EXPORT, "")

	// Do nothing if there are no records to export
	if len(exportList) == 0 {
		logHandler.Export.Printf("No records to export for %v", exportName)
		clock.Stop(0)
		return nil
	}

	logHandler.Export.Printf("Exporting %v record(s) as CSV in %v mode", len(exportList), mode)

	switch mode {
	case SINGLE:
		logHandler.Trace.Printf("Exporting to Single Record to Dumps path")
		exportName = buildName(exportName, exportList, idField, mode)
	case BULK:
		logHandler.Trace.Printf("Exporting to Bulk Records to Dumps path")
		exportName = buildName(exportName, exportList, idField, mode)
	case DEFAULTS:
		logHandler.Trace.Printf("Exporting to Defaults path")
		exportName = buildName("", exportList, idField, mode)
	default:
		logHandler.Error.Panicf("Unknown export mode %v, defaulting to BULK", mode)
		return fmt.Errorf("Unknown export mode %v, defaulting to BULK", mode)
	}

	path := paths.Dumps()
	if mode == DEFAULTS {
		path = paths.Defaults()
	}

	logHandler.Export.Printf("Exporting %v to [%v/%v.csv]", exportName, path.String(), exportName)

	exportFile := openTargetFile(exportName, EXPORT, logHandler.Export, CSV, path.String())
	defer exportFile.Close()

	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = FIELDSEPARATOR // Use tab-delimited format
		writer.UseCRLF = true
		return gocsv.NewSafeCSVWriter(writer)
	})

	_, err := gocsv.MarshalString(exportList) // Get all texts as CSV string
	if err != nil {
		logHandler.Export.Panicf("error exporting %v: %v", exportName, err.Error())
	}

	err = gocsv.MarshalFile(exportList, exportFile) // Get all texts as CSV string
	if err != nil {
		logHandler.Export.Panicf("error exporting %v: %v", exportName, err.Error())
	}

	// Example: # Generated (10) TestEntity Records at 15:04:05 2006-01-02 by 501/username on MacOS(mode=BULK)
	for i, item := range exportList {
		idFieldContent := getFieldValue(item, idField)
		logHandler.Export.Printf("Export (%v/%v) %v %v", i+1, len(exportList), getTypeName(item), idFieldContent)
	}
	noItems := len(exportList)
	// plurality := "s"
	// if noItems == 1 {
	// 	plurality = ""
	// }
	u, _ := user.Current()
	var by string
	if u != nil {
		by = u.Uid + SEP + u.Username
	} else {
		by = "sys_" + application.SystemIdentity()
	}
	on := application.SystemIdentity()
	os := application.OS()
	table := reflect.TypeOf(exportList)
	if table.Kind() == reflect.Slice {
		table = table.Elem()
	}
	msg := fmt.Sprintf("Exported (%v/%v) %v(s) to [%v] at %v %v by %v on %v(%v) in %v mode.", noItems, noItems, exportName, exportFile.Name(), time.Now().Format("15:04:05"), time.Now().Format("2006-01-02"), by, on, os, mode)
	exportFile.WriteString("# " + msg)

	exportFile.Close()

	logHandler.Export.Println(msg)
	// logHandler.EventLogger.Printf("Exported (%v/%v) %v(s) to [%v]", len(exportList), len(exportList), exportName, exportFile.Name())
	clock.Stop(len(exportList))
	return nil
}

func ExportDefaults[T any](exportList []*T, idField entities.Field) error {
	return ExportCSV("", exportList, idField, DEFAULTS)
}

func ExportJSON[T any](exportName string, exportList []*T, idField entities.Field) error {
	clock := timing.Start(exportName, EXPORT, "")

	//if exportName == "" {
	//	exportName = buildName(exportName, exportList, idField)
	//}
	logHandler.Trace.Printf("Exporting %v record(s) as JSON '%v'", len(exportList), exportName)

	for _, record := range exportList {
		// ID := reflect.ValueOf(record).FieldByName(idField.String())
		outputName := buildNameForRecord(exportName, record, idField)
		logHandler.Trace.Printf("Exporting %v.json", outputName)

		exportJSON(outputName, paths.Dumps(), record)
	}
	clock.Stop(1)
	return nil
}

func buildName[T any](baseName string, exportList []*T, idField entities.Field, mode mode) string {
	// Set default base name if not provided
	if baseName == "" {
		baseName = "Export"
	}

	// Handle empty list
	if len(exportList) == 0 {
		return idHelpers.GetUUID() + SEP + baseName
	}

	// Extract type name from first record
	firstRecord := exportList[0]
	typeName := getTypeName(firstRecord)
	uuid := idHelpers.GetUUID()

	// Handle DEFAULTS mode - return only type name
	if mode == DEFAULTS {
		return typeName
	}

	// Build base domain name with UUID
	domainName := uuid + SEP + typeName

	// Handle BULK mode with custom base name
	if mode == BULK && baseName != "" && idField.String() == "" {
		return domainName + SEP + baseName
	}

	if mode == BULK && baseName != "" {
		return domainName + SEP + baseName + SEP + "Bulk"
	}

	// For multiple records, return domain name only
	if len(exportList) > 1 {
		return domainName
	}

	// For single record, append field value if available
	fieldValue := getFieldValue(firstRecord, idField)
	if fieldValue != "" {
		return domainName + SEP + fieldValue + SEP + baseName
	}

	return domainName
}

// getTypeName extracts the type name from a potentially pointer-wrapped value
func getTypeName[T any](record T) string {
	typ := reflect.TypeOf(record)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ.Name()
}

// getFieldValue safely extracts a field value from a record
func getFieldValue[T any](record T, idField entities.Field) string {
	val := reflect.ValueOf(record)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if !val.IsValid() || val.Kind() != reflect.Struct {
		return ""
	}

	field := val.FieldByName(idField.String())
	if !field.IsValid() {
		return ""
	}

	return fmt.Sprintf("%v", field.Interface())
}

func buildNameForRecord[T any](baseName string, record T, idField entities.Field) string {
	logHandler.Trace.Printf("buildNameForRecord IN: baseName=[%v]", baseName)
	if baseName == "" {
		baseName = EXPORT
	}
	var domainName string
	typ := reflect.TypeOf(record)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	domainName = typ.Name()
	if domainName == "" {
		domainName = "Record"
	}
	domainName = idHelpers.GetUUID() + SEP + domainName

	val := reflect.ValueOf(record)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if !val.IsValid() || val.Kind() != reflect.Struct {
		return domainName
	}
	field := val.FieldByName(idField.String())
	if !field.IsValid() {
		return domainName
	}
	xx := field.Interface()
	domainName = domainName + SEP + fmt.Sprintf("%v", xx) + SEP + baseName
	logHandler.Trace.Printf("buildNameForRecord OUT: domainName=[%v]", domainName)
	return domainName
}

func exportJSON[T any](exportName string, where paths.FileSystemPath, record T) {
	logHandler.Trace.Printf("Exporting %v.json", exportName)
	logHandler.Export.Printf("Exporting %v %v.json", entities.GetStructType(record), exportName)
	exportFile := openTargetFile(exportName, EXPORT, logHandler.Export, JSON, where.String())
	defer exportFile.Close()

	exportFile.WriteString(godump.DumpJSONStr(record))

	exportFile.Close()
}
