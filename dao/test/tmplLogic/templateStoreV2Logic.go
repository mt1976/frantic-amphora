package tmpllogic

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/goforj/godump"
	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/test/templateStoreV3"
	"github.com/mt1976/frantic-core/logHandler"
)

// Login logs a user in (creating a record if needed) and updates last login metadata.
func Login(ctx context.Context, sq string) (templateStoreV3.TemplateStoreV3, error) {
	temp := buildUserStub(sq)
	var usr templateStoreV3.TemplateStoreV3

	logHandler.TraceLogger.Printf("%v", godump.DumpStr(temp))

	usrList, err := templateStoreV3.GetAllWhere(templateStoreV3.Fields.UserCode, temp.UserCode)
	if err != nil || len(usrList) == 0 {
		if err != nil {
			logHandler.WarningLogger.Printf("Warning=[%v] User=[%v]", err.Error(), temp.UserName)
			return templateStoreV3.New(), err
		}
		usr, err = Add(ctx, sq)
		if err != nil {
			logHandler.ErrorLogger.Printf("Warning=[%v] User=[%v]", err.Error(), temp.UserName)
			return templateStoreV3.New(), err
		}
		return usr, nil
	}

	// Existing user(s): update login details, return first updated.
	u := usrList[0]
	u.LastHost, _ = os.Hostname()
	u.LastLogin = time.Now()
	if err = u.UpdateWithAction(ctx, audit.LOGIN, fmt.Sprintf("User %v logged in", u.UserName)); err != nil {
		logHandler.WarningLogger.Printf("Warning=[%v] User=[%v]", err.Error(), u.UserName)
		return templateStoreV3.New(), err
	}
	usr = u
	return usr, nil
}

// Add creates and persists a new user record based on the current OS user.
func Add(ctx context.Context, sq string) (templateStoreV3.TemplateStoreV3, error) {
	testu := buildUserStub(sq)

	newUser := templateStoreV3.New()
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

	newUser.UserName = testu.UserName
	newUser.UID = testu.UID
	newUser.RealName = testu.RealName
	newUser.Email = testu.Email
	newUser.GID = testu.GID
	u, err := templateStoreV3.Create(ctx, newUser)
	if err != nil {
		logHandler.ErrorLogger.Printf("Error: '%v'", err.Error())
		return templateStoreV3.New(), err
	}

	return u, nil
}

// buildUserStub builds a record stub for the current OS user.
func buildUserStub(sq string) templateStoreV3.TemplateStoreV3 {
	currentUser, _ := user.Current()
	hostname, _ := os.Hostname()

	stub := templateStoreV3.New()
	stub.ID = 0
	stub.UID = fmt.Sprintf("%v%04v", currentUser.Uid, sq)
	stub.UserName = currentUser.Username
	stub.RealName = currentUser.Name
	stub.GID = currentUser.Gid
	stub.Email = strings.ToLower(fmt.Sprintf("%v@%v.com", currentUser.Username, hostname))
	stub.UserCode = BuildUserCode(stub)
	return stub
}
