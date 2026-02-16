package audit

import (
	"time"

	"github.com/mt1976/frantic-amphora/dao/entities"
)

// Audit represents the audit information for a data entity
type Audit struct {
	CreatedAt        time.Time         `json:"CreatedAt"`
	CreatedBy        string            `json:"CreatedBy"`
	CreatedOn        string            `json:"CreatedOn"`
	CreatedAtDisplay string            `json:"CreatedAtDisplay"`
	Updates          []AuditUpdateInfo `json:"Updates,omitempty"`
	DeletedAt        time.Time         `json:"DeletedAt"`
	DeletedBy        string            `json:"DeletedBy"`
	DeletedOn        string            `json:"DeletedOn"`
	DeletedAtDisplay string            `json:"DeletedAtDisplay"`
	AuditSequence    entities.Int      `json:"AuditSequence"`
	DBVersion        entities.Int      `json:"DBVersion"`
	// Empty     time.Time // Convience Field - Used to avoid erros with dates.
}

type AuditUpdateInfo struct {
	UpdatedAt        time.Time `json:"UpdatedAt"`
	UpdateAction     string    `json:"UpdateAction"`
	UpdatedBy        string    `json:"UpdatedBy"`
	UpdatedOn        string    `json:"UpdatedOn"`
	UpdatedAtDisplay string    `json:"UpdatedAtDisplay"`
	UpdateNotes      string    `json:"UpdateNotes"`
}

// Action represents an audit action with its properties
type Action struct {
	code        string
	short       string
	description string
	silent      bool
}
