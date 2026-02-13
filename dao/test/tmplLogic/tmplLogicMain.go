package tmpllogic

import (
	"context"
	"fmt"
	"os"
	"os/user"

	"github.com/goforj/godump"
	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/test/templateStoreV6"
	"github.com/mt1976/frantic-core/logHandler"
)

// Login logs a user in (creating a record if needed) and updates last login metadata.
func Login(ctx context.Context, sq string) (*templateStoreV6.TemplateStoreV6, error) {
	usr := templateStoreV6.New()

	logHandler.Trace.Printf("%v", godump.DumpStr(usr))

	usr, err := Add(ctx, sq)
	if err != nil {
		logHandler.Error.Printf("Warning=[%v] User=[%v]", err.Error(), usr.Name)
		return templateStoreV6.New(), err
	}
	return usr, nil

	// Existing user(s): update login details, return first updated.

	usr.LastHost, _ = os.Hostname()
	// u.LastLogin = time.Now()
	if err = usr.UpdateWithAction(ctx, audit.LOGIN, fmt.Sprintf("User %v logged in", usr.Name)); err != nil {
		logHandler.Warning.Printf("Warning=[%v] User=[%v]", err.Error(), usr.Name)
		return templateStoreV6.New(), err
	}

	return usr, nil
}

// Add creates and persists a new user record based on the current OS user.
func Add(ctx context.Context, sq string) (*templateStoreV6.TemplateStoreV6, error) {
	logHandler.Info.Printf("Adding new user to TemplateStoreV6: SQ=%v", sq)

	newUser := templateStoreV6.New()
	// use the creator to build the new record
	// _, skip, u, err := Creator(ctx, newUser)
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Error: '%v'", err.Error())
	// 	return templateStoreV2.New(), err
	// }
	// if skip {
	// 	logHandler.WarningLogger.Printf("Creation of %v record skipped by creator function", templateStoreV2.TableName)
	// 	return templateStoreV2.New(), nil
	// }

	newUser.Name = os.Getenv("USERNAME")
	if newUser.Name == "" {
		// Fallback to current OS user if USERNAME env var is not set
		if currentUser, err := user.Current(); err == nil {
			newUser.Name = currentUser.Username
		} else {
			logHandler.Error.Printf("Error retrieving current user: '%v'", err.Error())
			return templateStoreV6.New(), err
		}
	}
	newUser.Name = newUser.Name + "_" + sq
	newUser.Destination = newUser.Name + "'s Destination"
	newUser.Profile = "Standard"
	newUser.ProfileKey = "standard"
	newUser.ProfileEnrichment = "Standard profile with no special requirements"
	u, err := templateStoreV6.Create(ctx, newUser)
	if err != nil {
		logHandler.Error.Printf("Error: '%v'", err.Error())
		return templateStoreV6.New(), err
	}
	logHandler.Info.Printf("New user added to TemplateStoreV6: Name=%v, UID=%v", u.Name, u.Raw)
	return u, nil
}
