package tmpllogic

import (
	"context"
	"fmt"
	"time"

	"github.com/goforj/godump"
	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/test/templateStoreV6"
	"github.com/mt1976/frantic-core/dateHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func Creator(ctx context.Context, basics *templateStoreV6.TemplateStoreV6) (id string, skipPostCreate bool, record *templateStoreV6.TemplateStoreV6, err error) {
	// Custom creation logic can be added here
	logHandler.Service.Printf("Creator logic executed for TemplateStore ")
	record = templateStoreV6.New()
	// record.Key = idHelpers.Encode(sessionID)
	// record.Raw = sessionID
	record.Name = basics.Name
	record.Destination = basics.Destination
	record.Profile = basics.Profile
	record.ProfileKey = basics.ProfileKey
	record.ProfileEnrichment = basics.ProfileEnrichment
	record.LastHost = ""
	record.PostTest = append(record.PostTest, "CREATE@"+time.Now().Format(dateHelpers.Format.Detail))

	id = record.Name
	skipPostCreate = false
	// logHandler.ServiceLogger.Printf("Creator generated ID %v for TemplateStore record with Name: %v", id, record.Name)
	err = nil
	return
}

func DuplicateCheck(record *templateStoreV6.TemplateStoreV6) (found bool, err error) {
	logHandler.Warning.Printf("Performing duplicate check for %v record %v", templateStoreV6.TableName, record.Key)

	if record.Name == "" {
		logHandler.Warning.Printf("Duplicate check failed for %v record: Key is empty", templateStoreV6.TableName)
		found = false
		err = fmt.Errorf("key is empty")
		return
	}
	responseRecord, err := templateStoreV6.GetBy(templateStoreV6.Fields.Name, record.Name)
	if err != nil {
		logHandler.Warning.Printf("Duplicate check failed for %v record: %v - NO RECORD EXISTS, DONT SKIP CREATE", templateStoreV6.TableName, err)
		found = false
		err = nil
		return
	}
	if responseRecord.Audit.DeletedBy != "" {
		logHandler.Warning.Printf("Duplicate check passed for %v record %v: Existing record is deleted, DONT SKIP CREATE", templateStoreV6.TableName, record.Name)
		found = false
		err = nil
		return
	}
	logHandler.Warning.Printf("Duplicate check found existing %v record %v: SKIP CREATE", templateStoreV6.TableName, record.Name)
	found = true
	err = nil
	return
}

// jobProcessor processes jobs related to the TemplateStore tableName entity.
func JobProcessor(name, desc string) {
	clock := timing.Start(name, "Process", desc)
	count := 0
	templateEntries, err := templateStoreV6.GetAll()
	if err != nil {
		logHandler.Error.Printf("[%v] Error: '%v'", name, err.Error())
		return
	}

	notemplateEntries := len(templateEntries)
	if notemplateEntries == 0 {
		logHandler.Service.Printf("[%v] No %vs to process", name, templateStoreV6.TableName)
		clock.Stop(0)
		return
	}

	for templateEntryIndex, templateRecord := range templateEntries {
		logHandler.Service.Printf("[%v] %v(%v/%v) %v", name, templateStoreV6.TableName, templateEntryIndex+1, notemplateEntries, templateRecord.Raw)
		_ = templateRecord.UpdateWithAction(context.Background(), audit.SERVICE, "Job Processing "+desc)
		count++
	}
	clock.Stop(count)
}

func PostCreate(ctx context.Context, record *templateStoreV6.TemplateStoreV6) error {
	// Custom post-create logic can be added here
	// logHandler.WarningLogger.Printf("PostCreate logic executed for TemplateStore Key: %v %v %+v Locks: %+v", record.Key, godump.DumpStr(ctx), record.PostTest, godump.DumpStr(record.Lock))
	logHandler.Warning.Printf("PostCreate logic executed for TemplateStore Key: %v %v %+v Locks: %+v", record.Key, godump.DumpStr(ctx), record.PostTest, godump.DumpStr(record.Lock))
	record.PostTest = append(record.PostTest, "POSTCREATE@"+time.Now().Format(dateHelpers.Format.Detail))
	// logHandler.WarningLogger.Printf("PostCreate after create: %+v Lock:%v", record.PostTest, godump.DumpStr(record.Lock))
	// Need to store record also to capture the update to PostTest
	// fmt.Printf("DOING POST CREATE")
	actionError := record.PostUpdateUpdate(ctx, "Post Create Processing")
	if actionError != nil {
		logHandler.Error.Printf("Error updating record during post-create processing for TemplateStore Key: %v: %v", record.Key, actionError.Error())
		return actionError
	}
	// logHandler.WarningLogger.Printf("PostCreate after update: %+v %v", record.PostTest, record.ID)
	// logHandler.WarningLogger.Printf("PostCreate after update: %+v %v", record.PostTest, record.ID)
	// logHandler.WarningLogger.Printf("PostCreate after update: %+v %v", record.PostTest, record.ID)
	// logHandler.WarningLogger.Printf("PostCreate after update: %+v %v", record.PostTest, record.ID)

	return nil
}

func PostUpdate(ctx context.Context, record *templateStoreV6.TemplateStoreV6) error {
	// Custom post-update logic can be added here
	logHandler.Trace.Printf("PostUpdate logic executed for TemplateStore Key: %v %v %+v", record.Key, godump.DumpStr(ctx), record.PostTest)
	record.PostTest = append(record.PostTest, "UPDATE@"+time.Now().Format(dateHelpers.Format.DMY))
	logHandler.Trace.Printf("PostUpdate after update: %+v", record.PostTest)
	// Need to store record also to capture the update to PostTest
	// fmt.Printf("DOING POST CREATE")

	actionError := record.PostUpdateUpdate(ctx, "Post Update Processing")
	if actionError != nil {
		logHandler.Error.Printf("Error updating record during post-update processing for TemplateStore Key: %v: %v", record.Key, actionError.Error())
		return actionError
	}
	logHandler.Trace.Printf("PostUpdate after update: %+v %v", record.PostTest, record.ID)
	return nil
}

func PostDelete(ctx context.Context, record *templateStoreV6.TemplateStoreV6) error {
	// Custom post-delete logic can be added here
	logHandler.Trace.Printf("PostDelete logic executed for TemplateStore Name: %v", record.Name)
	return nil
}

func Cloner(ctx context.Context, source *templateStoreV6.TemplateStoreV6) (cloneRecord *templateStoreV6.TemplateStoreV6, err error) {
	// Custom cloning logic can be added here
	logHandler.Warning.Printf("Cloner logic executed for TemplateStore Name: %v", source.Name)
	cloneRecord = templateStoreV6.New()
	cloneRecord.Name = source.Name + "_clone"
	cloneRecord.Destination = source.Destination
	cloneRecord.Profile = source.Profile
	cloneRecord.ProfileKey = source.ProfileKey
	cloneRecord.ProfileEnrichment = source.ProfileEnrichment
	cloneRecord.LastHost = ""
	cloneRecord.PostTest = append(source.PostTest, "CLONED@"+time.Now().Format(dateHelpers.Format.Detail))
	err = nil
	logHandler.Warning.Printf("Cloner created clone for TemplateStore Name: %v with new Name: %v (%v)", source.Name, cloneRecord.Name, godump.DumpStr(cloneRecord))
	return
}
