package tmpllogic

import (
	"context"
	"fmt"
	"time"

	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/test/templateStoreV2"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func Creator(basics templateStoreV2.TemplateStore) (string, templateStoreV2.TemplateStore, error) {
	// Custom creation logic can be added here

	id := idHelpers.GetUUID()

	record := templateStoreV2.New()
	// record.Key = idHelpers.Encode(sessionID)
	// record.Raw = sessionID
	record.UserName = basics.UserName
	record.UID = basics.UID
	record.RealName = basics.RealName
	record.Email = basics.Email
	record.GID = basics.GID
	record.Active.IsTrue()
	record.LastLogin = time.Time{}
	record.LastHost = ""
	record.UserCode = BuildUserCode(record)

	return id, record, nil
}

// BuildUserCode creates a stable user code string used for lookups.
func BuildUserCode(u templateStoreV2.TemplateStore) string {
	return fmt.Sprintf("%v%v%v", u.UID, "_", u.UserName)
}

func DuplicateCheck(record *templateStoreV2.TemplateStore) (bool, error) {

	responseRecord, err := templateStoreV2.GetBy(templateStoreV2.Fields.Key, record.Key)
	if err != nil {
		return false, err
	}
	if responseRecord.Audit.DeletedBy != "" {
		return false, nil
	}
	return true, nil
}

// jobProcessor processes jobs related to the TemplateStore tableName entity.
func JobProcessor(name, desc string) {
	clock := timing.Start(name, "Process", desc)
	count := 0

	templateEntries, err := templateStoreV2.GetAll()
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] Error: '%v'", name, err.Error())
		return
	}

	notemplateEntries := len(templateEntries)
	if notemplateEntries == 0 {
		logHandler.ServiceLogger.Printf("[%v] No %vs to process", name, templateStoreV2.TableName)
		clock.Stop(0)
		return
	}

	for templateEntryIndex, templateRecord := range templateEntries {
		logHandler.ServiceLogger.Printf("[%v] %v(%v/%v) %v", name, templateStoreV2.TableName, templateEntryIndex+1, notemplateEntries, templateRecord.Raw)
		_ = templateRecord.UpdateWithAction(context.Background(), audit.SERVICE, "Job Processing "+desc)
		count++
	}
	clock.Stop(count)
}
