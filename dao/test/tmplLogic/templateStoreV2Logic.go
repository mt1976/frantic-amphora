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
	"github.com/mt1976/frantic-amphora/dao/test/templateStoreV2"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
)

// Login logs a user in (creating a record if needed) and updates last login metadata.
func Login(ctx context.Context, sq string) (templateStoreV2.TemplateStore, error) {
	temp := buildUserStub(sq)
	var usr templateStoreV2.TemplateStore

	logHandler.TraceLogger.Printf("%v", godump.DumpStr(temp))

	usrList, err := templateStoreV2.GetAllWhere(templateStoreV2.Fields.UserCode, temp.UserCode)
	if err != nil || len(usrList) == 0 {
		if err != nil {
			logHandler.WarningLogger.Printf("Warning=[%v] User=[%v]", err.Error(), temp.UserName)
			return templateStoreV2.New(), err
		}
		usr, err = Add(ctx, sq)
		if err != nil {
			logHandler.ErrorLogger.Printf("Warning=[%v] User=[%v]", err.Error(), temp.UserName)
			return templateStoreV2.New(), err
		}
		return usr, nil
	}

	// Existing user(s): update login details, return first updated.
	u := usrList[0]
	u.LastHost, _ = os.Hostname()
	u.LastLogin = time.Now()
	if err = u.UpdateWithAction(ctx, audit.LOGIN, fmt.Sprintf("User %v logged in", u.UserName)); err != nil {
		logHandler.WarningLogger.Printf("Warning=[%v] User=[%v]", err.Error(), u.UserName)
		return templateStoreV2.New(), err
	}
	usr = u
	return usr, nil
}

// Add creates and persists a new user record based on the current OS user.
func Add(ctx context.Context, sq string) (templateStoreV2.TemplateStore, error) {
	testu := buildUserStub(sq)

	newUser := templateStoreV2.New()
	// use the creator to build the new record
	id, u, err := Creator(ctx, newUser)
	if err != nil {
		logHandler.ErrorLogger.Printf("Error: '%v'", err.Error())
		return templateStoreV2.New(), err
	}

	u.Key = idHelpers.Encode(id)
	u.Raw = id
	u.UserName = testu.UserName
	u.UID = testu.UID
	u.RealName = testu.RealName
	u.Email = testu.Email
	u.GID = testu.GID

	u, err = templateStoreV2.Create(ctx, u)
	if err != nil {
		logHandler.ErrorLogger.Printf("Error: '%v'", err.Error())
		return templateStoreV2.New(), err
	}

	return u, nil
}

// buildUserStub builds a record stub for the current OS user.
func buildUserStub(sq string) templateStoreV2.TemplateStore {
	currentUser, _ := user.Current()
	hostname, _ := os.Hostname()

	stub := templateStoreV2.New()
	stub.ID = 0
	stub.UID = fmt.Sprintf("%v%04v", currentUser.Uid, sq)
	stub.UserName = currentUser.Username
	stub.RealName = currentUser.Name
	stub.GID = currentUser.Gid
	stub.Email = strings.ToLower(fmt.Sprintf("%v@%v.com", currentUser.Username, hostname))
	stub.UserCode = BuildUserCode(stub)
	return stub
}
