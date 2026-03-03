package audit

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/mt1976/frantic-core/application"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/contextHandler"
	"github.com/mt1976/frantic-core/logHandler"
)

var (
	hostNameOnce   sync.Once
	cachedHostName string
)

func getHostName() string {
	hostNameOnce.Do(func() {
		cachedHostName = application.HostName()
	})
	return cachedHostName
}

// getDBVersion retrieves the current database version
func getDBVersion() string {
	// Implement the logic to get the DB version without importing the dao package
	return cfg.GetDatabase_Version()
}

func (a *Action) popMessage() string {
	message := a.description
	a.description = ""
	return message
}

func (a *Audit) Spew() error {
	// Spew the Audit Data
	noAudit := len(a.Updates)
	// logger.InfoLogger.Printf(" No Updates=[%v]", noAudit)
	if noAudit > 0 {
		for i := range a.Updates {
			upd := a.Updates[i]
			logHandler.Trace.Printf(AUDITMSG, upperName, upd.UpdateAction, upd.UpdatedAtDisplay, upd.UpdatedBy, upd.UpdatedOn, upd.UpdateNotes)
		}
	}
	return nil
}

func (a *Audit) VersionString() string {
	if a.Version != "" {
		return a.Version
	} else {
		version := strconv.Itoa(a.DBVersion.Int())
		a.Version = version
		return a.Version
	}
}

func (a *Audit) VersionAbsolute() int {
	ver := a.VersionString()
	// Split the version on "."
	verParts := strings.Split(ver, ".")
	if len(verParts) == 0 {
		return 0
	}
	// Get the first part of the version and convert to int
	majorVersion, err := strconv.Atoi(verParts[0])
	if err != nil {
		majorVersion = 0
	}
	minorVersion, err := strconv.Atoi(verParts[1])
	if err != nil {
		minorVersion = 0
	}
	patchVersion, err := strconv.Atoi(verParts[2])
	if err != nil {
		patchVersion = 0
	}
	// Calculate the absolute version as (major * 10000) + (minor * 100) + patch
	absoluteVersion := (majorVersion * 10000) + (minorVersion * 100) + patchVersion
	return absoluteVersion
}

func getAuditUserCode(ctx context.Context) (string, error) {
	defaultUser := cfg.GetServiceUser_UserCode()

	if ctx == context.Background() {
		usr := "svc_" + getHostName()
		return usr, nil
	}

	// Implement the logic to get the user without importing the dao package
	if ctx == context.TODO() || ctx == nil {
		usr := "sys_" + getHostName()
		return usr, nil
	}

	// Get the current user from the context
	sessionUser := contextHandler.GetSession_UserCode(ctx)
	// ctx.Value(cfg.GetSecuritySessionKey_UserCode())
	if sessionUser != "" {
		return sessionUser, nil
	}
	return defaultUser, commonErrors.ErrContextCannotGetUserCode
}
