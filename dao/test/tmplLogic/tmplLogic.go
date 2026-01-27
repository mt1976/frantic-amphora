package tmpllogic

import (
	"context"
	"fmt"
	"time"

	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/test/templateStoreV3"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func Creator(ctx context.Context, basics templateStoreV3.TemplateStoreV3) (string, bool, templateStoreV3.TemplateStoreV3, error) {
	// Custom creation logic can be added here

	id := idHelpers.GetUUID()

	record := templateStoreV3.New()
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

	return id, false, record, nil
}

// BuildUserCode creates a stable user code string used for lookups.
func BuildUserCode(u templateStoreV3.TemplateStoreV3) string {
	return fmt.Sprintf("%v%v%v", u.UID, "_", u.UserName)
}

func DuplicateCheck(record *templateStoreV3.TemplateStoreV3) (bool, error) {
	logHandler.EventLogger.Printf("Performing duplicate check for %v record %v", templateStoreV3.TableName, record.Key)
	if record.Key == "" {
		logHandler.InfoLogger.Printf("Duplicate check failed for %v record: Key is empty", templateStoreV3.TableName)
		return false, nil
	}
	responseRecord, err := templateStoreV3.GetBy(templateStoreV3.Fields.Key, record.Key)
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

	templateEntries, err := templateStoreV3.GetAll()
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] Error: '%v'", name, err.Error())
		return
	}

	notemplateEntries := len(templateEntries)
	if notemplateEntries == 0 {
		logHandler.ServiceLogger.Printf("[%v] No %vs to process", name, templateStoreV3.TableName)
		clock.Stop(0)
		return
	}

	for templateEntryIndex, templateRecord := range templateEntries {
		logHandler.ServiceLogger.Printf("[%v] %v(%v/%v) %v", name, templateStoreV3.TableName, templateEntryIndex+1, notemplateEntries, templateRecord.Raw)
		_ = templateRecord.UpdateWithAction(context.Background(), audit.SERVICE, "Job Processing "+desc)
		count++
	}
	clock.Stop(count)
}

func PostCreate(ctx context.Context, record *templateStoreV3.TemplateStoreV3) (error, bool, string) {
	// Custom post-create logic can be added here
	logHandler.ServiceLogger.Printf("PostCreate logic executed for TemplateStore Key: %v", record.Key)
	update := true
	message := "post create processing completed"
	return nil, update, message
}

func PostUpdate(ctx context.Context, record *templateStoreV3.TemplateStoreV3) (error, bool, string) {
	// Custom post-update logic can be added here
	logHandler.ServiceLogger.Printf("PostUpdate logic executed for TemplateStore Key: %v", record.Key)
	update := true
	message := "post update processing completed"
	return nil, update, message
}
