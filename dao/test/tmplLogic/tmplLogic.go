package tmpllogic

import (
	"context"
	"fmt"
	"time"

	"github.com/goforj/godump"
	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/test/templateStoreV3"
	"github.com/mt1976/frantic-core/dateHelpers"
	"github.com/mt1976/frantic-core/idHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func Creator(ctx context.Context, basics *templateStoreV3.TemplateStoreV3) (string, bool, *templateStoreV3.TemplateStoreV3, error) {
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
func BuildUserCode(u *templateStoreV3.TemplateStoreV3) string {
	return fmt.Sprintf("%v%v%v", u.UID, "_", u.UserName)
}

func DuplicateCheck(record *templateStoreV3.TemplateStoreV3) (bool, error) {
	logHandler.TraceLogger.Printf("Performing duplicate check for %v record %v", templateStoreV3.TableName, record.Key)

	if record.Key == "" {
		logHandler.InfoLogger.Printf("Duplicate check failed for %v record: Key is empty", templateStoreV3.TableName)

		return false, nil
	}
	responseRecord, err := templateStoreV3.GetBy(templateStoreV3.Fields.Key, record.Key)
	if err != nil {
		logHandler.TraceLogger.Printf("Duplicate check failed for %v record: %v - NO RECORD EXISTS, DONT SKIP CREATE", templateStoreV3.TableName, err)
		return false, nil
	}
	if responseRecord.Audit.DeletedBy != "" {
		logHandler.TraceLogger.Printf("Duplicate check passed for %v record %v: Existing record is deleted, DONT SKIP CREATE", templateStoreV3.TableName, record.Key)
		return false, nil
	}
	logHandler.TraceLogger.Printf("Duplicate check found existing %v record %v: SKIP CREATE", templateStoreV3.TableName, record.Key)
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

func PostCreate(ctx context.Context, record *templateStoreV3.TemplateStoreV3) error {
	// Custom post-create logic can be added here
	// logHandler.WarningLogger.Printf("PostCreate logic executed for TemplateStore Key: %v %v %+v Locks: %+v", record.Key, godump.DumpStr(ctx), record.PostTest, godump.DumpStr(record.Lock))

	record.PostTest = append(record.PostTest, "CREATE@"+time.Now().Format(dateHelpers.Format.DMY))
	// logHandler.WarningLogger.Printf("PostCreate after create: %+v Lock:%v", record.PostTest, godump.DumpStr(record.Lock))
	// Need to store record also to capture the update to PostTest
	// fmt.Printf("DOING POST CREATE")
	actionError := record.PostUpdateUpdate(ctx, "Post Create Processing")
	if actionError != nil {
		logHandler.ErrorLogger.Printf("Error updating record during post-create processing for TemplateStore Key: %v: %v", record.Key, actionError.Error())
		return actionError
	}
	// logHandler.WarningLogger.Printf("PostCreate after update: %+v %v", record.PostTest, record.ID)
	// logHandler.WarningLogger.Printf("PostCreate after update: %+v %v", record.PostTest, record.ID)
	// logHandler.WarningLogger.Printf("PostCreate after update: %+v %v", record.PostTest, record.ID)
	// logHandler.WarningLogger.Printf("PostCreate after update: %+v %v", record.PostTest, record.ID)

	return nil
}

func PostUpdate(ctx context.Context, record *templateStoreV3.TemplateStoreV3) error {
	// Custom post-update logic can be added here
	logHandler.TraceLogger.Printf("PostUpdate logic executed for TemplateStore Key: %v %v %+v", record.Key, godump.DumpStr(ctx), record.PostTest)
	record.PostTest = append(record.PostTest, "UPDATE@"+time.Now().Format(dateHelpers.Format.DMY))
	logHandler.TraceLogger.Printf("PostUpdate after update: %+v", record.PostTest)
	// Need to store record also to capture the update to PostTest
	// fmt.Printf("DOING POST CREATE")

	actionError := record.PostUpdateUpdate(ctx, "Post Update Processing")
	if actionError != nil {
		logHandler.ErrorLogger.Printf("Error updating record during post-update processing for TemplateStore Key: %v: %v", record.Key, actionError.Error())
		return actionError
	}
	logHandler.TraceLogger.Printf("PostUpdate after update: %+v %v", record.PostTest, record.ID)
	return nil
}

func PostDelete(ctx context.Context, record *templateStoreV3.TemplateStoreV3) error {
	// Custom post-delete logic can be added here
	logHandler.TraceLogger.Printf("PostDelete logic executed for TemplateStore Key: %v", record.Key)
	return nil
}
