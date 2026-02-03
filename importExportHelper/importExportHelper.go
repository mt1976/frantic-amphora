package importExportHelper

import (
	"fmt"
	"log"
	"os"

	"github.com/mt1976/frantic-core/paths"
)

// TemplateImportData is a struct to hold the data from the CSV file
// it is used to import the data into the database
// The struct tags are used to map the fields to the CSV columns
// this struct should be customised to suit the specific requirements of the entryination table/DAO.

type (
	mode          string
	fileExtension string
)

const (
	DEFAULTS       mode          = "DEFAULTS"
	SINGLE         mode          = "SINGLE"
	BULK           mode          = "BULK"
	FIELDSEPARATOR rune          = '|'
	SEP            string        = "-"
	IMPORT         string        = "Import"
	EXPORT         string        = "Export"
	CSV            fileExtension = "csv"
	JSON           fileExtension = "json"
)

func openTargetFile(in, _ string, useLog *log.Logger, extention fileExtension, path string) *os.File {
	// defaultPath := paths.Defaults()
	templateDataFileName := in + "." + string(extention)
	fileName := fmt.Sprintf("%s%s/%s", paths.Application().String(), path, templateDataFileName)

	dataFileHandle, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		useLog.Fatalf("error opening file: %v", err)
		panic(err)
	}
	// useLog.Printf("%ving %vs from File=[%v]", action, in, dataFileHandle.Name())
	return dataFileHandle
}
