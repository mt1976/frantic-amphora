package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/alitto/pond/v2"
	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/database"
	"github.com/mt1976/frantic-amphora/dao/entities"
	"github.com/mt1976/frantic-amphora/dao/test/templateStoreV7"
	logic "github.com/mt1976/frantic-amphora/dao/test/tmplLogic"
	"github.com/mt1976/frantic-core/contextHandler"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/mathHelpers"
)

// Were going to be upgrading and testing the DAO module

type Sausage struct {
	Field1 string
	Field2 entities.Int
	When   time.Time
}
type Supper struct {
	FieldA string
	FieldB int
	FieldC float64
}

var (
	SampleKey  entities.Field
	Sample2Key entities.Field
)

func main() {
	ctx := context.Background()
	shutdown, stopSignals := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stopSignals()

	logHandler.Info.Println("Creating Worker Pool...")
	workerPool := pond.NewPool(100, pond.WithContext(ctx), pond.WithQueueSize(100))
	// Now add the pool to the context so it can be used by the jobs

	ctx = contextHandler.AddWorkerPoolToContext(ctx, workerPool)
	defer workerPool.StopAndWait()
	logHandler.Info.Printf("Worker Pool Created with %d workers", 10)

	ctx = contextHandler.SetSession_UserKey(ctx, "UserKey123")
	ctx = contextHandler.SetSession_UserCode(ctx, "UserCode123")

	logHandler.Banner("INFO", "START", "Starting DAO Test Application - Phase 1")

	logHandler.Info.Println("Initialize User Store")
	templateStoreV7.Initialise(ctx, false)
	templateStoreV7.RegisterCreator(logic.Creator)
	templateStoreV7.RegisterDuplicateCheck(logic.DuplicateCheck)
	templateStoreV7.RegisterWorker(logic.JobProcessor)
	templateStoreV7.RegisterPostCreate(logic.PostCreate)
	templateStoreV7.RegisterPostUpdate(logic.PostUpdate)
	templateStoreV7.RegisterPostDelete(logic.PostDelete)
	templateStoreV7.RegisterCloner(logic.Cloner)

	// tripStore.Initialise(ctx, false)

	logHandler.Info.Println("Clear Down User Store")

	templateStoreV7.ClearDown(ctx)

	totalElapsed := time.Duration(0)
	// start := time.Now()
	// for i := 0; i < 5; i++ {
	in_start := time.Now()
	phase := strconv.Itoa(0)
	msg2 := test(ctx, phase, 5)
	logHandler.Info.Printf("Phase 2 Test Message: %v", msg2)

	if failures := pagingTest(ctx, "P", 20, 5); failures > 0 {
		logHandler.Error.Printf("Paging Test FAILED with %d failure(s)", failures)
	} else {
		logHandler.Info.Printf("Paging Test PASSED")
	}
	if err := entityTypesTest(ctx); err != nil {
		logHandler.Error.Printf("Entity Types Test FAILED: %v", err)
	} else {
		logHandler.Info.Printf("Entity Types Test PASSED")
	}

	stop := time.Now()
	in_elapsed := stop.Sub(in_start)
	totalElapsed += in_elapsed
	//		logHandler.ErrorLogger.Printf("P_%v Test Duration: %v Start: %v Stop: %v", i+1, in_elapsed, start.Format(time.RFC3339), stop.Format(time.RFC3339))
	// Cache.PurgeExpiredEntries()
	// xx := tripStore.New()
	// xx.Name = "Test Trip"
	// xx.Destination = "Somewhere"
	// xx.StartDate = time.Now().Add(48 * time.Hour)
	// xx.EndDate = time.Now().Add(120 * time.Hour)

	// newTrip, err := tripStore.Create(ctx, xx)
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Error creating trip: %v", err)
	// } else {
	// 	logHandler.InfoLogger.Printf("Created Trip: %v (%v to %v)", newTrip.Name, newTrip.StartDate.Format("2006-01-02"), newTrip.EndDate.Format("2006-01-02"))
	// }
	//}
	// stop := time.Now()
	// recs, _ := templateStoreV2.GetAll()
	// for _, r := range recs {
	// 	logHandler.InfoLogger.Printf("P1 User: %v - %v", r.ID, r.RealName)
	// }

	// templateStoreV2.Close()

	// //Cache.Spew()

	// //Cache.SynchroniseAll()
	// //Cache.Disable(templateStoreV2.TemplateStore{})
	logHandler.Info.Printf("Total Test Duration: %v", totalElapsed)

	logHandler.Banner("INFO", "STOP", "Stopping DAO Test Application - Phase 1")

	// Need to keep program running to allow time to inspect logs and outputs before the deferred workerPool.StopAndWait() is called and the worker pool is stopped.
	// godump.DumpJSON("Test Complete - Entering Idle State", workerPool)
	// Log all the info from the worker pool to verify it is still active and processing any remaining tasks before we exit.
	debugPool(workerPool)

	logHandler.Info.Printf("Use Ctrl+c to stop application.")
	<-shutdown.Done()
	logHandler.Info.Printf("Shutdown requested; waiting for worker pool to finish.")
}

func debugPool(workerPool pond.Pool) {
	logHandler.Info.Printf("Worker Pool - Completed Tasks: %v", workerPool.CompletedTasks())
	logHandler.Info.Printf("Worker Pool - Running Workers: %v", workerPool.RunningWorkers())
	logHandler.Info.Printf("Worker Pool - Dropped Tasks: %v", workerPool.DroppedTasks())
	logHandler.Info.Printf("Worker Pool - Waiting Tasks: %v", workerPool.WaitingTasks())
	logHandler.Info.Printf("Worker Pool - Submitted Tasks: %v", workerPool.SubmittedTasks())
	logHandler.Info.Printf("Worker Pool - Succefully Tasks: %v", workerPool.SuccessfulTasks())
	logHandler.Info.Printf("Worker Pool - Failed Tasks: %v", workerPool.FailedTasks())
	logHandler.Info.Printf("Worker Pool - Max Concurrency: %v", workerPool.MaxConcurrency())
	logHandler.Info.Printf("Worker Pool - Queue Size: %v", workerPool.QueueSize())
}

func test(ctx context.Context, phase string, baselineThings int) string {
	logHandler.Warning.Printf("Phase %v Creating %v Baseline Things", phase, baselineThings)
	// //Cache.Activate(templateStoreV2.TemplateStore{})
	// //Cache.RegisterExpiry(templateStoreV2.TemplateStore{}, time.Duration(baselineThings)*time.Second)
	// //Cache.RegisterKey(templateStoreV2.TemplateStore{}, templateStoreV2.Fields.Key)
	// //Cache.RegisterSynchroniser(templateStoreV2.TemplateStore{}, templateStoreV2.CacheSynchroniser(ctx))
	// //Cache.RegisterHydrator(templateStoreV2.TemplateStore{}, templateStoreV2.CacheHydrator(ctx))

	logHandler.Info.Printf("Phase %v Adding %d Baseline Things to Store", phase, baselineThings)
	var newRecords []*templateStoreV7.TemplateStoreV7
	for i := 0; i < baselineThings; i++ {
		logHandler.Info.Printf("Phase %v Creating Baseline User %v", phase, i+1)
		newRecord, info := logic.Login(ctx, fmt.Sprintf("%04v", i+1))
		if info != nil {
			logHandler.Error.Printf("Phase %v Error creating Baseline User %v: %v", phase, i+1, info)
		}
		logHandler.Info.Printf("Phase %v Created Baseline User %v %v", phase, i+1, newRecord.Name)
		// fmt.Print(".")
		//	//Cache.AddEntry(usr)
		newRecords = append(newRecords, newRecord)

		clone, err := newRecord.Clone(ctx)
		if err != nil {
			logHandler.Error.Printf("Phase %v Error cloning Baseline User %v: %v", phase, i+1, info)
		}
		logHandler.Info.Printf("Phase %v Cloned Baseline User %v %v", phase, i+1, clone.Name)

		dbVersion := clone.Audit.VersionString()
		appVersion := clone.Audit.VersionAbsolute()
		logHandler.Event.Printf("Phase %v Cloned Baseline User '%v' Version: '%v' Absolute Version: %v", phase, clone.Name, dbVersion, appVersion)
		logHandler.Event.Printf("Phase %v Cloned Baseline User '%v' Version: '%v' Absolute Version: %v", phase, clone.Name, dbVersion, appVersion)
		logHandler.Event.Printf("Phase %v Cloned Baseline User '%v' Version: '%v' Absolute Version: %v", phase, clone.Name, dbVersion, appVersion)
	}

	logHandler.Info.Printf("Phase %v Baseline %v Things Added to Store %d", phase, baselineThings, len(newRecords))
	//	logHandler.InfoLogger.Printf("Phase %v Hydrating Cache for Things", phase)

	////Cache.HydrateForType(templateStoreV2.TemplateStore{})

	setupTemplates, err := templateStoreV7.GetAll()
	if err != nil {
		logHandler.Error.Printf("Phase %v Error getting all Things: %v", phase, err)
	}
	logHandler.Info.Printf("Phase %v Setup Templates Loaded: %v", phase, len(setupTemplates))
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Phase %v Error getting all Things: %v", phase, err)
	// }
	// //Cache.SpewForType(templateStoreV2.TemplateStore{})

	// logHandler.InfoLogger.Printf("Phase %v Setup Templates Loaded: %v", phase, len(setupTemplates))
	// if len(setupTemplates) != baselineThings {
	// 	logHandler.ErrorLogger.Printf("Phase %v Setup Template count mismatch: expected %v, got %v", phase, baselineThings, len(setupTemplates))
	// 	logHandler.ErrorLogger.Printf("Phase %v Setup Template count mismatch: expected %v, got %v", phase, baselineThings, len(setupTemplates))
	// 	os.Exit(0)
	// 	return fmt.Sprintf("Phase %v Setup Template count mismatch: expected %v, got %v", phase, baselineThings, len(setupTemplates))
	// }
	uKey := ""
	for x, u := range newRecords {
		logHandler.Info.Printf("Phase %v Evaluating : (%v/%v) %v %v", phase, x, baselineThings, u.Name, u.Key)
		if mathHelpers.CoinToss() {
			logHandler.Info.Printf("Phase %v User: (%v/%v) %v %v", phase, x, baselineThings, u.Name, u.Key)
			uKey = u.Key
			break
		}
		logHandler.Info.Printf("PostTest: %+v", u.PostTest)
	}
	if uKey == "" && len(newRecords) > 0 {
		uKey = newRecords[0].Key
	}
	logHandler.Info.Printf("Phase %v Selected User Key: %v %v", phase, uKey, len(newRecords))

	// Read back in the created user with key uKey, and update its LastHost field to "orion+datetime"
	userRec, getErr := templateStoreV7.GetBy(templateStoreV7.Fields.Key, uKey)
	if getErr != nil {
		logHandler.Error.Printf("Phase %v Error getting user by key %v: %v", phase, uKey, getErr)
	} else {
		logHandler.Info.Printf("Phase %v Retrieved User by Key %v: %v", phase, uKey, userRec.Raw)
		userRec.LastHost = fmt.Sprintf("orion-%v", time.Now().Format("150405"))
		updateErr := userRec.UpdateWithAction(ctx, audit.UPDATE, "Test Update of LastHost field "+userRec.LastHost)
		if updateErr != nil {
			logHandler.Error.Printf("Phase %v Error updating user %v: %v", phase, uKey, updateErr)
		} else {
			logHandler.Info.Printf("Phase %v Updated User %v LastHost to %v", phase, uKey, userRec.LastHost)
		}
	}

	// os.Exit(0)

	err = userRec.ExportRecordToCSV("SingleMode")
	if err != nil {
		logHandler.Error.Printf("Phase %v Error exporting user %v as CSV: %v", phase, uKey, err)
	}

	logHandler.Info.Printf("Phase %v Exporting All Things as CSV", phase)
	err = templateStoreV7.ExportAllToCSV("END" + phase)
	if err != nil {
		logHandler.Error.Printf("Phase %v Error exporting all Things as CSV: %v", phase, err)
	}
	logHandler.Info.Printf("Phase %v Exporting All Things as JSON", phase)

	logHandler.Info.Printf("Phase %v Exporting All Things to Defaults", phase)
	err = templateStoreV7.ExportDefaults()
	if err != nil {
		logHandler.Error.Printf("Phase %v Error exporting all Things to Defaults: %v", phase, err)
	}

	logHandler.Info.Printf("Cleardown User Store for Phase %v Imports from Defaults", phase)
	err = templateStoreV7.ClearDown(ctx)
	if err != nil {
		logHandler.Error.Printf("Phase %v Error cleardown user store: %v", phase, err)
	}

	// Check its empty before we import
	templateEntries, err := templateStoreV7.GetAll()
	if err != nil {
		logHandler.Error.Printf("Phase %v Error getting all Things before import: %v", phase, err)
	}
	notemplateEntries := len(templateEntries)
	logHandler.Info.Printf("Phase %v Things before Import: %v", phase, notemplateEntries)
	if notemplateEntries != 0 {
		logHandler.Error.Printf("Phase %v User store not empty before import: expected 0, got %v", phase, notemplateEntries)
		os.Exit(0)
	}

	err = templateStoreV7.ClearDown(ctx)
	if err != nil {
		logHandler.Error.Printf("Phase %v Error cleardown user store: %v", phase, err)
	}

	logHandler.Info.Printf("Phase %v Importing All Things from Defaults", phase)
	err = templateStoreV7.ImportDefaults(ctx)
	if err != nil {
		logHandler.Error.Printf("Phase %v Error importing all Things from Defaults: %v", phase, err)
	}
	logHandler.Info.Printf("Phase %v Completed Imports from Defaults", phase)

	// Get all Things again to verify count
	allThings, err := templateStoreV7.GetAll()
	if err != nil {
		logHandler.Error.Printf("Phase %v Error getting all Things after import: %v", phase, err)
	}
	logHandler.Info.Printf("Phase %v Total Things after Import: %v", phase, len(allThings))
	// // Benchmark: retrieve a cached record repeatedly and report timings.
	// // This is intentionally simple and uses FindByKey which should hit the cache after hydration.
	// fetchIterations := 1
	// fetchesPerIteration := 1
	// if fetchIterations < 1 {
	// 	fetchIterations = 1
	// }
	// if fetchesPerIteration < 1 {
	// 	fetchesPerIteration = 1
	// }

	// var overallTotal time.Duration
	// var overallMin time.Duration
	// var overallMax time.Duration
	// var overallSamples int
	// for iter := 1; iter <= fetchIterations; iter++ {
	// 	var iterTotal time.Duration
	// 	var iterMin time.Duration
	// 	var iterMax time.Duration
	// 	var iterSamples int

	// 	for j := 0; j < fetchesPerIteration; j++ {
	// 		fetchStart := time.Now()
	// 		_, fetchErr := templateStoreV2.GetBy(templateStoreV2.Fields.Key, uKey)
	// 		fetchElapsed := time.Since(fetchStart)
	// 		if fetchErr != nil {
	// 			logHandler.ErrorLogger.Printf("Phase %v Cache fetch error (iter=%d, n=%d): %v", phase, iter, j+1, fetchErr)
	// 			continue
	// 		}

	// 		if iterSamples == 0 {
	// 			iterMin, iterMax = fetchElapsed, fetchElapsed
	// 		} else {
	// 			if fetchElapsed < iterMin {
	// 				iterMin = fetchElapsed
	// 			}
	// 			if fetchElapsed > iterMax {
	// 				iterMax = fetchElapsed
	// 			}
	// 		}
	// 		iterTotal += fetchElapsed
	// 		iterSamples++

	// 		if overallSamples == 0 {
	// 			overallMin, overallMax = fetchElapsed, fetchElapsed
	// 		} else {
	// 			if fetchElapsed < overallMin {
	// 				overallMin = fetchElapsed
	// 			}
	// 			if fetchElapsed > overallMax {
	// 				overallMax = fetchElapsed
	// 			}
	// 		}
	// 		overallTotal += fetchElapsed
	// 		overallSamples++
	// 	}

	// 	if iterSamples == 0 {
	// 		logHandler.WarningLogger.Printf("Phase %v Cache fetch iteration %d: no successful samples", phase, iter)
	// 		continue
	// 	}
	// 	iterAvg := iterTotal / time.Duration(iterSamples)
	// 	logHandler.InfoLogger.Printf(
	// 		"Phase %v Cache fetch iteration %d: avg=%v min=%v max=%v samples=%d",
	// 		phase, iter, iterAvg, iterMin, iterMax, iterSamples,
	// 	)
	// }

	// if overallSamples > 0 {
	// 	overallAvg := overallTotal / time.Duration(overallSamples)
	// 	logHandler.InfoLogger.Printf(
	// 		"Phase %v Cache fetch overall: avg=%v min=%v max=%v samples=%d",
	// 		phase, overallAvg, overallMin, overallMax, overallSamples,
	// 	)
	// } else {
	// 	logHandler.WarningLogger.Printf("Phase %v Cache fetch overall: no successful samples", phase)
	// }

	// Things, err := templateStoreV2.GetAll()
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Phase %v Error getting all Things: %v", phase, err)
	// }
	// logHandler.InfoLogger.Printf("Phase %v Total Things: %v", phase, len(Things))
	// // for _, u := range Things {
	// // 	logHandler.InfoLogger.Printf("Phase %v User: %v", phase, u.RealName)
	// // }
	// logHandler.InfoLogger.Printf("Phase %v Counting all Things", phase)
	// count, err := templateStoreV2.Count()
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Phase %v Error counting Things: %v", phase, err)
	// } else {
	// 	logHandler.InfoLogger.Printf("Phase %v User count: %d", phase, count)
	// }
	// if count != baselineThings {
	// 	logHandler.ErrorLogger.Printf("Phase %v User count mismatch: expected %d, got %d", phase, baselineThings, count)
	// }

	// logHandler.InfoLogger.Printf("Phase %v Counting active Things with LastHost='orion'", phase)
	// countw, err := templateStoreV2.CountWhere(templateStoreV2.Fields.LastHost, "orion")
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Phase %v Error counting active Things: %v", phase, err)
	// } else {
	// 	logHandler.InfoLogger.Printf("Phase %v Active user count: %d", phase, countw)
	// }
	// if countw > baselineThings {
	// 	logHandler.ErrorLogger.Printf("Phase %v Active user count exceeds baseline: %d > %d", phase, countw, baselineThings)
	// }
	// if countw < 0 {
	// 	logHandler.ErrorLogger.Printf("Phase %v Active user count is negative: %d", phase, countw)
	// }
	// logHandler.InfoLogger.Printf("Phase %v Get %v", phase, uKey)
	// rec, err := templateStoreV2.GetBy(templateStoreV2.Fields.Key, uKey)
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Phase %v Error getting user by key: %v", phase, err)
	// } else {
	// 	logHandler.InfoLogger.Printf("Phase %v User by key: %v", phase, rec.RealName)
	// }
	// logHandler.InfoLogger.Printf("Phase %v Repeated user load %v", phase, len(Things))
	// for _, u := range Things {
	// 	_, err := templateStoreV2.GetAllWhere(templateStoreV2.Fields.UserName, u.UserName)
	// 	if err != nil {
	// 		logHandler.ErrorLogger.Printf("Phase %v Error getting Things by LastHost: %v", phase, err)
	// 	} else {
	// 		//			logHandler.InfoLogger.Printf("Phase %v Things by LastHost: %v/%v (%v) %v", phase, x+1, len(Things), len(dum), u.UserName)
	// 	}
	// }

	// // logHandler.InfoLogger.Printf("Phase %v Flushing Cache", phase)
	// // err = //Cache.SynchroniseForType(templateStoreV2.TemplateStore{})
	// // if err != nil {
	// // 	logHandler.ErrorLogger.Printf("Phase %v Error flushing cache: %v", phase, err)
	// // } else {
	// // 	logHandler.InfoLogger.Printf("Phase %v Cache flushed successfully", phase)
	// // }
	logHandler.Info.Printf("Phase %v Completed", phase)

	// created, updated, noTables, noCacheEntries := //Cache.Stats()
	// logHandler.InfoLogger.Printf("Cache Stats - Created: %v, Updated: %v, Tables: %v, Entries: %v", created.Format(time.RFC3339Nano), updated.Format(time.RFC3339Nano), noTables, noCacheEntries)

	return uKey
}

// pagingTest creates extra records, then verifies GetAllPaged, GetAllPagedWhere, PageCount
// and PageCountWhere against Count/CountWhere/GetAll. It returns the number of failures.
func pagingTest(ctx context.Context, phase string, extraRecords, pageSize int) int {
	logHandler.Banner("INFO", "PAGING", fmt.Sprintf("Phase %v Paging Test - %d extra records, page size %d", phase, extraRecords, pageSize))
	failures := 0
	fail := func(format string, args ...any) {
		failures++
		logHandler.Error.Printf("Phase %v Paging FAIL: "+format, append([]any{phase}, args...)...)
	}

	logHandler.Info.Printf("Phase %v Creating %d Paging Records", phase, extraRecords)
	premiumKey := "premium"
	for i := 0; i < extraRecords; i++ {
		record, err := logic.Add(ctx, fmt.Sprintf("%v%03d", phase, i+1))
		if err != nil {
			fail("error creating paging record %d: %v", i+1, err)
			continue
		}
		if i%3 == 0 {
			record.ProfileKey = premiumKey
			if err := record.UpdateWithAction(ctx, audit.UPDATE, "Paging Test set ProfileKey="+premiumKey); err != nil {
				fail("error updating paging record %v: %v", record.Key, err)
			}
		}
	}

	total, err := templateStoreV7.Count()
	if err != nil {
		fail("error counting records: %v", err)
		return failures
	}
	allRecords, err := templateStoreV7.GetAll()
	if err != nil {
		fail("error getting all records: %v", err)
		return failures
	}
	if len(allRecords) != total {
		fail("GetAll returned %d records, Count returned %d", len(allRecords), total)
	}
	logHandler.Info.Printf("Phase %v Total Records for Paging: %d", phase, total)

	// Unfiltered paging
	expectedPages := database.TotalPages(total, pageSize)
	pageCount, err := templateStoreV7.PageCount(pageSize)
	if err != nil {
		fail("PageCount error: %v", err)
	} else if pageCount != expectedPages {
		fail("PageCount=%d expected %d", pageCount, expectedPages)
	}
	failures += verifyPages(phase, "ALL", total, expectedPages, pageSize, recordIDs(allRecords), func(page int) ([]*templateStoreV7.TemplateStoreV7, database.PageInfo, error) {
		return templateStoreV7.GetAllPaged(page, pageSize)
	})

	// Filtered paging
	whereTotal, err := templateStoreV7.CountWhere(templateStoreV7.Fields.ProfileKey, premiumKey)
	if err != nil {
		fail("CountWhere error: %v", err)
		return failures
	}
	whereRecords, err := templateStoreV7.GetAllWhere(templateStoreV7.Fields.ProfileKey, premiumKey)
	if err != nil {
		fail("GetAllWhere error: %v", err)
		return failures
	}
	if whereTotal == 0 || whereTotal == total {
		fail("filtered test needs a subset of records, got %d of %d", whereTotal, total)
	}
	expectedWherePages := database.TotalPages(whereTotal, pageSize)
	wherePageCount, err := templateStoreV7.PageCountWhere(templateStoreV7.Fields.ProfileKey, premiumKey, pageSize)
	if err != nil {
		fail("PageCountWhere error: %v", err)
	} else if wherePageCount != expectedWherePages {
		fail("PageCountWhere=%d expected %d", wherePageCount, expectedWherePages)
	}
	failures += verifyPages(phase, "WHERE "+premiumKey, whereTotal, expectedWherePages, pageSize, recordIDs(whereRecords), func(page int) ([]*templateStoreV7.TemplateStoreV7, database.PageInfo, error) {
		return templateStoreV7.GetAllPagedWhere(templateStoreV7.Fields.ProfileKey, premiumKey, page, pageSize)
	})

	// Edge cases
	beyond, _, err := templateStoreV7.GetAllPaged(expectedPages+5, pageSize)
	if err != nil || len(beyond) != 0 {
		fail("page beyond range returned %d records, err=%v (expected 0, nil)", len(beyond), err)
	}

	logHandler.Banner("INFO", "PAGING", fmt.Sprintf("Phase %v Paging Test Completed - %d failure(s)", phase, failures))
	return failures
}

// verifyPages walks every page returned by fetch and checks the page metadata and that
// the combined pages contain exactly the expected record IDs with no duplicates.
func verifyPages(phase, label string, total, expectedPages, pageSize int, expectedIDs map[int]bool, fetch func(page int) ([]*templateStoreV7.TemplateStoreV7, database.PageInfo, error)) int {
	failures := 0
	fail := func(format string, args ...any) {
		failures++
		logHandler.Error.Printf("Phase %v Paging %v FAIL: "+format, append([]any{phase, label}, args...)...)
	}

	seen := map[int]bool{}
	for page := 1; page <= expectedPages; page++ {
		records, info, err := fetch(page)
		if err != nil {
			fail("page %d error: %v", page, err)
			continue
		}
		logHandler.Info.Printf("Phase %v Paging %v %v - %d records (Total=%d)", phase, label, info.String(), len(records), info.TotalRecords)
		if info.Page != page || info.TotalPages != expectedPages || info.TotalRecords != total || info.PageSize != pageSize {
			fail("page %d info mismatch: %+v (expected pages=%d total=%d)", page, info, expectedPages, total)
		}
		if info.HasPrevious() != (page > 1) || info.HasNext() != (page < expectedPages) {
			fail("page %d HasPrevious=%v HasNext=%v", page, info.HasPrevious(), info.HasNext())
		}
		expectedLen := pageSize
		if page == expectedPages {
			expectedLen = total - (expectedPages-1)*pageSize
		}
		if len(records) != expectedLen {
			fail("page %d returned %d records, expected %d", page, len(records), expectedLen)
		}
		for _, r := range records {
			if seen[r.ID] {
				fail("record %d returned on more than one page", r.ID)
			}
			if !expectedIDs[r.ID] {
				fail("record %d not expected in result set", r.ID)
			}
			seen[r.ID] = true
		}
	}
	if len(seen) != len(expectedIDs) {
		fail("pages returned %d distinct records, expected %d", len(seen), len(expectedIDs))
	}
	return failures
}

func recordIDs(records []*templateStoreV7.TemplateStoreV7) map[int]bool {
	ids := make(map[int]bool, len(records))
	for _, r := range records {
		ids[r.ID] = true
	}
	return ids
}

func entityTypeSample() *templateStoreV7.TemplateStoreV7 {
	record := templateStoreV7.New()
	record.EntityInt.Set(-42)
	record.EntityInt64 = entities.Int64{Value: "-9223372036854775808"}
	record.EntityUInt64 = entities.UInt64{Value: "18446744073709551615"}
	record.EntityUInt = entities.UInt{Value: "42"}
	record.EntityInt32 = entities.Int32{Value: "-2147483648"}
	record.EntityUInt32 = entities.UInt32{Value: "4294967295"}
	record.EntityFloat.Set(-123.125)
	record.EntityFloat32 = entities.Float32{Value: "1.25"}
	record.EntityFloat64 = entities.Float64{Value: "-123456.125"}
	record.EntityDecimal = entities.Decimal{Value: "12345.6789"}
	record.EntityMoney = entities.Money{Value: "19.95"}
	record.EntityPercentage = entities.Percentage{Value: "25.5"}
	record.EntityRate = entities.Rate{Value: "0.035"}
	record.EntityCurrency.Set("GBP", 123.45)
	record.EntityBool.Set(true)
	record.EntityStormBool.Set(true)
	stamp := time.Date(2026, 10, 30, 12, 34, 56, 123456789, time.FixedZone("sample", 2*60*60))
	record.EntityDate.Set(stamp)
	record.EntityTime.Set(stamp)
	record.EntityDateTime.Set(stamp)
	record.EntityField = entities.Field("SampleField")
	record.EntityURI.Set("/entity-tests/sample")
	record.EntityAction.Set("TEST")
	record.EntityKey.Set("entity-test-key")
	record.EntityTable = entities.Table("EntityTestTable")
	return record
}

func createEntityTypeRecords(ctx context.Context) ([]*templateStoreV7.TemplateStoreV7, error) {
	var records []*templateStoreV7.TemplateStoreV7
	for i := 0; i < 7; i++ {
		record := templateStoreV7.New()
		record.Name = fmt.Sprintf("EntityTest_%02d", i)
		record.Destination = "Entity test destination"
		record, err := templateStoreV7.Create(ctx, record)
		if err != nil {
			return nil, fmt.Errorf("creating entity test record %d: %w", i, err)
		}
		// Set fields after creation so application-specific creators cannot discard them.
		if i < 5 {
			sample := reflect.ValueOf(entityTypeSample()).Elem()
			target := reflect.ValueOf(record).Elem()
			for _, field := range reflect.VisibleFields(sample.Type()) {
				if strings.HasPrefix(field.Name, "Entity") {
					target.FieldByIndex(field.Index).Set(sample.FieldByIndex(field.Index))
				}
			}
			if err := record.Update(ctx, "Populate entity type test values"); err != nil {
				return nil, fmt.Errorf("updating entity test record %d: %w", i, err)
			}
		}
		records = append(records, record)
	}
	return records, nil
}

func verifyEntityType(field entities.Field, value any, expected []*templateStoreV7.TemplateStoreV7) error {
	const pageSize = 2
	for _, record := range expected {
		loaded, err := templateStoreV7.GetBy(templateStoreV7.Fields.Key, record.Key)
		if err != nil {
			return err
		}
		got := reflect.ValueOf(loaded).Elem().FieldByName(field.String()).Interface()
		if !reflect.DeepEqual(got, value) {
			return fmt.Errorf("%s round-trip: got %#v, want %#v", field, got, value)
		}
	}
	count, err := templateStoreV7.CountWhere(field, value)
	if err != nil {
		return err
	}
	if count != len(expected) {
		return fmt.Errorf("%s CountWhere=%d, want %d", field, count, len(expected))
	}
	pages, err := templateStoreV7.PageCountWhere(field, value, pageSize)
	if err != nil {
		return err
	}
	expectedPages := (len(expected) + pageSize - 1) / pageSize
	if pages != expectedPages {
		return fmt.Errorf("%s PageCountWhere=%d, want %d", field, pages, expectedPages)
	}
	all, err := templateStoreV7.GetAllWhere(field, value)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(recordIDs(all), recordIDs(expected)) {
		return fmt.Errorf("%s GetAllWhere returned unexpected record IDs", field)
	}
	failures := verifyPages("Entity", field.String(), len(expected), pages, pageSize, recordIDs(expected), func(page int) ([]*templateStoreV7.TemplateStoreV7, database.PageInfo, error) {
		records, info, err := templateStoreV7.GetAllPagedWhere(field, value, page, pageSize)
		if err != nil {
			return nil, info, err
		}
		for _, record := range records {
			got := reflect.ValueOf(record).Elem().FieldByName(field.String()).Interface()
			if !reflect.DeepEqual(got, value) {
				return nil, info, fmt.Errorf("%s paged value: got %#v, want %#v", field, got, value)
			}
		}
		return records, info, nil
	})
	if failures != 0 {
		return fmt.Errorf("%s paging: %d failure(s)", field, failures)
	}
	beyond, info, err := templateStoreV7.GetAllPagedWhere(field, value, pages+1, pageSize)
	if err != nil {
		return err
	}
	if len(beyond) != 0 || info.TotalRecords != len(expected) || info.TotalPages != pages {
		return fmt.Errorf("%s beyond-last-page result: %d records, info=%+v", field, len(beyond), info)
	}
	return nil
}

func entityTypesTest(ctx context.Context) error {
	records, err := createEntityTypeRecords(ctx)
	if err != nil {
		return err
	}
	sample := reflect.ValueOf(entityTypeSample()).Elem()
	var failures int
	for _, field := range reflect.VisibleFields(sample.Type()) {
		if !strings.HasPrefix(field.Name, "Entity") {
			continue
		}
		if err := verifyEntityType(entities.Field(field.Name), sample.FieldByIndex(field.Index).Interface(), records[:5]); err != nil {
			logHandler.Error.Printf("Entity type %s FAILED: %v", field.Type, err)
			failures++
		} else {
			logHandler.Info.Printf("Entity type %s PASSED", field.Type)
		}
	}
	if failures != 0 {
		return fmt.Errorf("%d entity type(s) failed", failures)
	}
	return nil
}
