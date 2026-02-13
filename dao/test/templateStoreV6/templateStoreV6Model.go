// Data Access Object for the TripStore table
// Template Version: 0.5.24 - 2026-01-31
// Generated
// Date: 10/02/2026 & 13:16
// Who : matttownsend (orion)

package templateStoreV6

import (
	"sync"
	"time"

	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/entities"
)

// TableName is the canonical DAO table identifier for this package.
var (
	TableName              = entities.Table("TemplateStoreV6")
	tableName              = TableName.String()
	URI       entities.URI = entities.URI("trip")
	PK        entities.KEY = entities.KEY("tripkey")
	RAW       entities.KEY = entities.KEY("tripalt")
)

// The TemplateStoreV6 struct defines the data model for the TemplateStoreV6 table.
// Adjust domain fields and tags as required in the TemplateStoreV6.definitions file.
type TemplateStoreV6 struct {
	// The primary key field(s), managed by the framework, DO NOT MODIFY
	ID  int    `storm:"id,increment=100"`
	Key string `storm:"index,unique"`
	Raw string `storm:"index,unique"`
	// Audit information, managed by the framework, DO NOT MODIFY
	Audit audit.Audit `csv:"-"`
	// Locking information, managed by the framework, DO NOT MODIFY
	Lock sync.Mutex `csv:"-"` // Add this field to enable record locking for concurrent updates

	// Domain specific fields
	//
	Name                  string        `validate:"required,max=25,min=5"` // this field will not be indexed
	Profile               string        `storm:"index"`                    // profile of the trip
	ProfileKey            string        `storm:"index"`                    // profile of the trip
	ProfileEnrichment     string        // Used to Display Date            // profile name of the trip
	StatusKey             string        `storm:"index"` // Status key of the trip
	StatusEnrichment      string        // Only used during the display of records
	Destination           string        `storm:"index" validate:"required,max=35,min=4"` // Destination of the trip
	DestinationKey        string        `storm:"index"`                                  // Internal Key for Destination
	DestinationEnrichment string        // Only used during the display of records
	Notes                 string        `validate:"max=75"` // Notes about the trip
	StartDate             time.Time     // Start Date of the Trip
	EndDate               time.Time     // End Date of the Trip
	DisplayStartDate      string        // External Start Date of the Trip
	DisplayEndDate        string        // External End Date of the Trip
	DurationInDays        entities.Int  // Duration of the Trip in Days
	MinimumNumberOfItems  entities.Int  // Minimum Number of Items
	Countdown             entities.Int  // Days remaining until the trip starts
	TotalItems            entities.Int  // Total number of items for the trip
	TotalItemsPacked      entities.Int  // Total number of items packed for the trip
	Year                  string        `storm:"index"` // Year of the Trip
	DisplayID             string        ``              // Display ID of the Trip
	SequenceNumber        entities.Int  // Sequence Number of the Trip
	PackingProgress       entities.Int  // Packing Progress percentage
	TripProgress          entities.Int  // Trip Progress percentage
	Completed             entities.Bool // Is the trip completed?
	StartNotified         entities.Bool // Has the start notification been sent?
	EndWarned             entities.Bool // Has the end warning been sent?
	EndNotified           entities.Bool // Has the end notification been sent?

	// Add no more fields below this line
}

type fieldNames struct {
	// The primary key field(s), managed by the framework, DO NOT MODIFY
	ID  entities.Field
	Key entities.Field
	Raw entities.Field
	// The audit information, managed by the framework, DO NOT MODIFY
	Audit entities.Field
	// Domain specific fields
	Name                  entities.Field
	Profile               entities.Field
	ProfileKey            entities.Field
	ProfileEnrichment     entities.Field
	StatusKey             entities.Field
	StatusEnrichment      entities.Field
	Destination           entities.Field
	DestinationKey        entities.Field
	DestinationEnrichment entities.Field
	Notes                 entities.Field
	StartDate             entities.Field
	EndDate               entities.Field
	DisplayStartDate      entities.Field
	DisplayEndDate        entities.Field
	DurationInDays        entities.Field
	MinimumNumberOfItems  entities.Field
	Countdown             entities.Field
	TotalItems            entities.Field
	TotalItemsPacked      entities.Field
	Year                  entities.Field
	DisplayID             entities.Field
	SequenceNumber        entities.Field
	PackingProgress       entities.Field
	TripProgress          entities.Field
	Completed             entities.Field
	StartNotified         entities.Field
	EndWarned             entities.Field
	EndNotified           entities.Field
	Travel                entities.Field
	Hotel                 entities.Field
	Desk                  entities.Field

	// Add no more fields below this line
}

// Fields provides strongly-typed field names for use with GetBy/GetAllWhere/etc.
//
// Example: GetBy(Fields.Key, "abc")
//
// Note: the values are the struct field names as stored in Storm.
var Fields = fieldNames{
	// The primary key field(s), managed by the framework, DO NOT MODIFY
	ID:  "ID",
	Key: "Key",
	Raw: "Raw",
	// The audit information, managed by the framework, DO NOT MODIFY
	Audit: "Audit",
	// tableName-specific fields, please modify as required
	Name:                  "Name",
	Profile:               "Profile",
	ProfileKey:            "ProfileKey",
	ProfileEnrichment:     "ProfileEnrichment",
	StatusKey:             "StatusKey",
	StatusEnrichment:      "StatusEnrichment",
	Destination:           "Destination",
	DestinationKey:        "DestinationKey",
	DestinationEnrichment: "DestinationEnrichment",
	Notes:                 "Notes",
	StartDate:             "StartDate",
	EndDate:               "EndDate",
	DisplayStartDate:      "DisplayStartDate",
	DisplayEndDate:        "DisplayEndDate",
	DurationInDays:        "DurationInDays",
	MinimumNumberOfItems:  "MinimumNumberOfItems",
	Countdown:             "Countdown",
	TotalItems:            "TotalItems",
	TotalItemsPacked:      "TotalItemsPacked",
	Year:                  "Year",
	DisplayID:             "DisplayID",
	SequenceNumber:        "SequenceNumber",
	PackingProgress:       "PackingProgress",
	TripProgress:          "TripProgress",
	Completed:             "Completed",
	StartNotified:         "StartNotified",
	EndWarned:             "EndWarned",
	EndNotified:           "EndNotified",
	Travel:                "Travel",
	Hotel:                 "Hotel",
	Desk:                  "Desk",
	// Add no more fields below this line
}
