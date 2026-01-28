package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mt1976/frantic-amphora/dao/audit"
	"github.com/mt1976/frantic-amphora/dao/entities"
	"github.com/mt1976/frantic-amphora/dao/test/templateStoreV3"
	tmpllogic "github.com/mt1976/frantic-amphora/dao/test/tmplLogic"
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

var SampleKey entities.Field
var Sample2Key entities.Field

func main() {

	ctx := context.Background()
	// SampleKey = entities.Field("Field1")
	// Sample2Key = entities.Field("FieldA")

	// //Cache.Initialise()
	// //Cache.Spew()
	// //Cache.Activate(Sausage{})
	// //Cache.RegisterExpiry(Sausage{}, 30*time.Minute+15*time.Second)
	// //Cache.RegisterKey(Sausage{}, SampleKey)
	// //Cache.Spew()
	// xx := Sausage{}
	// xx.Field1 = "Bum"
	// xx.Field2.Set(123)
	// xx.When = time.Now()
	// //Cache.AddEntry(xx)
	// xx.Field1 = "Bum2"
	// xx.Field2.Set(456)
	// zz := xx
	// //Cache.AddEntry(xx)
	// xx.Field1 = "Bum3"
	// xx.Field2.Set(789)
	// xx.When = time.Now().Add(10 * time.Minute)
	// //Cache.AddEntry(xx)
	// xx.Field1 = "Bum4"
	// xx.Field2.Set(101112)
	// //Cache.AddEntry(xx)
	// //Cache.Spew()
	// yy := Supper{}
	// yy.FieldA = "Foo"
	// yy.FieldB = 789
	// yy.FieldC = 12.34
	// //Cache.Activate(Supper{})
	// //Cache.RegisterKey(Supper{}, Sample2Key)

	// //Cache.AddEntry(yy)
	// rtn, err := //Cache.FindByKey(Sausage{}, "Bum2")
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Error Finding Entry: %v", err)
	// } else {
	// 	logHandler.InfoLogger.Printf("Found Entry: %+v", rtn)
	// }
	// //Cache.RemoveEntry(zz)
	// //Cache.RemoveByKey(Sausage{}, "Bum3")
	// ////Cache.Spew()
	// created, updated, noTables, noCacheEntries := //Cache.Stats()
	// logHandler.InfoLogger.Printf("Cache Stats - Created: %v, Updated: %v, Tables: %v, Entries: %v", created.Format(time.RFC3339Nano), updated.Format(time.RFC3339Nano), noTables, noCacheEntries)
	// //os.Exit(0)
	// Placeholder main function
	logHandler.InfoBanner("INFO", "START", "Starting DAO Test Application - Phase 1")

	logHandler.InfoLogger.Println("Initialize User Store")
	templateStoreV3.Initialise(ctx, false)
	templateStoreV3.RegisterCreator(tmpllogic.Creator)
	templateStoreV3.RegisterDuplicateCheck(tmpllogic.DuplicateCheck)
	templateStoreV3.RegisterWorker(tmpllogic.JobProcessor)
	templateStoreV3.RegisterPostCreate(tmpllogic.PostCreate)
	templateStoreV3.RegisterPostUpdate(tmpllogic.PostUpdate)

	logHandler.InfoLogger.Println("Clear Down User Store")
	templateStoreV3.ClearDown(ctx)

	totalElapsed := time.Duration(0)
	//start := time.Now()
	for i := 0; i < 2; i++ {
		in_start := time.Now()

		msg2 := test(ctx, "ONE", i+1)
		logHandler.ServiceLogger.Printf("Phase 2 Test Message: %v", msg2)

		stop := time.Now()
		in_elapsed := stop.Sub(in_start)
		totalElapsed += in_elapsed
		//		logHandler.ErrorLogger.Printf("P_%v Test Duration: %v Start: %v Stop: %v", i+1, in_elapsed, start.Format(time.RFC3339), stop.Format(time.RFC3339))
		//Cache.PurgeExpiredEntries()

	}
	// stop := time.Now()
	// recs, _ := templateStoreV2.GetAll()
	// for _, r := range recs {
	// 	logHandler.InfoLogger.Printf("P1 User: %v - %v", r.ID, r.RealName)
	// }

	//templateStoreV2.Close()

	// //Cache.Spew()

	// //Cache.SynchroniseAll()
	// //Cache.Disable(templateStoreV2.TemplateStore{})
	logHandler.InfoLogger.Printf("Total Test Duration: %v", totalElapsed)

	logHandler.InfoBanner("INFO", "STOP", "Stopping DAO Test Application - Phase 1")
	// templateStoreV2.ExportAllAsCSV("AllUsers")

	// templateStoreV2.ExportAllAsJSON("AllUsers")

}

func test(ctx context.Context, phase string, baselineUsers int) string {

	logHandler.InfoLogger.Printf("Phase %v Creating %v Baseline Users", phase, baselineUsers)
	// //Cache.Activate(templateStoreV2.TemplateStore{})
	// //Cache.RegisterExpiry(templateStoreV2.TemplateStore{}, time.Duration(baselineUsers)*time.Second)
	// //Cache.RegisterKey(templateStoreV2.TemplateStore{}, templateStoreV2.Fields.Key)
	// //Cache.RegisterSynchroniser(templateStoreV2.TemplateStore{}, templateStoreV2.CacheSynchroniser(ctx))
	// //Cache.RegisterHydrator(templateStoreV2.TemplateStore{}, templateStoreV2.CacheHydrator(ctx))

	logHandler.InfoLogger.Printf("Phase %v Adding %d Baseline Users to Store", phase, baselineUsers)

	for i := 0; i < baselineUsers; i++ {
		logHandler.WarningLogger.Printf("Phase %v Creating Baseline User %v", phase, i+1)
		_, info := tmpllogic.Login(ctx, fmt.Sprintf("%04v", i))
		if info != nil {
			logHandler.ErrorLogger.Printf("Phase %v Error creating Baseline User %v: %v", phase, i+1, info)
		}
		logHandler.WarningLogger.Printf("Phase %v Created Baseline User %v", phase, i+1)
		//fmt.Print(".")
		//	//Cache.AddEntry(usr)
	}

	logHandler.InfoLogger.Printf("Phase %v Baseline %v Users Added to Store", phase, baselineUsers)
	logHandler.InfoLogger.Printf("Phase %v Hydrating Cache for Users", phase)

	////Cache.HydrateForType(templateStoreV2.TemplateStore{})

	setupTemplates, err := templateStoreV3.GetAll()
	if err != nil {
		logHandler.ErrorLogger.Printf("Phase %v Error getting all users: %v", phase, err)
	}
	logHandler.InfoLogger.Printf("Phase %v Setup Templates Loaded: %v", phase, len(setupTemplates))
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Phase %v Error getting all users: %v", phase, err)
	// }

	// //Cache.SpewForType(templateStoreV2.TemplateStore{})

	// logHandler.InfoLogger.Printf("Phase %v Setup Templates Loaded: %v", phase, len(setupTemplates))
	// if len(setupTemplates) != baselineUsers {
	// 	logHandler.ErrorLogger.Printf("Phase %v Setup Template count mismatch: expected %v, got %v", phase, baselineUsers, len(setupTemplates))
	// 	logHandler.ErrorLogger.Printf("Phase %v Setup Template count mismatch: expected %v, got %v", phase, baselineUsers, len(setupTemplates))
	// 	os.Exit(0)
	// 	return fmt.Sprintf("Phase %v Setup Template count mismatch: expected %v, got %v", phase, baselineUsers, len(setupTemplates))
	// }
	uKey := ""
	for x, u := range setupTemplates {
		if mathHelpers.CoinToss() {
			logHandler.InfoLogger.Printf("Phase %v User: (%v/%v) %v %v", phase, x, baselineUsers, u.RealName, u.Key)
			uKey = u.Key
			break
		}
		logHandler.ErrorLogger.Printf("PostTest: %+v", u.PostTest)
	}
	if uKey == "" && len(setupTemplates) > 0 {
		uKey = setupTemplates[0].Key
	}
	logHandler.InfoLogger.Printf("Phase %v Selected User Key: %v", phase, uKey)

	// Read back in the created user with key uKey, and update its LastHost field to "orion+datetime"
	userRec, getErr := templateStoreV3.GetBy(templateStoreV3.Fields.Key, uKey)
	if getErr != nil {
		logHandler.ErrorLogger.Printf("Phase %v Error getting user by key %v: %v", phase, uKey, getErr)
	} else {
		logHandler.InfoLogger.Printf("Phase %v Retrieved User by Key %v: %v", phase, uKey, userRec.RealName)
		userRec.LastHost = fmt.Sprintf("orion-%v", time.Now().Format("150405"))
		updateErr := userRec.UpdateWithAction(ctx, audit.UPDATE, "Test Update of LastHost")
		if updateErr != nil {
			logHandler.ErrorLogger.Printf("Phase %v Error updating user %v: %v", phase, uKey, updateErr)
		} else {
			logHandler.InfoLogger.Printf("Phase %v Updated User %v LastHost to %v", phase, uKey, userRec.LastHost)
		}
	}

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

	// users, err := templateStoreV2.GetAll()
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Phase %v Error getting all users: %v", phase, err)
	// }
	// logHandler.InfoLogger.Printf("Phase %v Total Users: %v", phase, len(users))
	// // for _, u := range users {
	// // 	logHandler.InfoLogger.Printf("Phase %v User: %v", phase, u.RealName)
	// // }
	// logHandler.InfoLogger.Printf("Phase %v Counting all users", phase)
	// count, err := templateStoreV2.Count()
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Phase %v Error counting users: %v", phase, err)
	// } else {
	// 	logHandler.InfoLogger.Printf("Phase %v User count: %d", phase, count)
	// }
	// if count != baselineUsers {
	// 	logHandler.ErrorLogger.Printf("Phase %v User count mismatch: expected %d, got %d", phase, baselineUsers, count)
	// }

	// logHandler.InfoLogger.Printf("Phase %v Counting active users with LastHost='orion'", phase)
	// countw, err := templateStoreV2.CountWhere(templateStoreV2.Fields.LastHost, "orion")
	// if err != nil {
	// 	logHandler.ErrorLogger.Printf("Phase %v Error counting active users: %v", phase, err)
	// } else {
	// 	logHandler.InfoLogger.Printf("Phase %v Active user count: %d", phase, countw)
	// }
	// if countw > baselineUsers {
	// 	logHandler.ErrorLogger.Printf("Phase %v Active user count exceeds baseline: %d > %d", phase, countw, baselineUsers)
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
	// logHandler.InfoLogger.Printf("Phase %v Repeated user load %v", phase, len(users))
	// for _, u := range users {
	// 	_, err := templateStoreV2.GetAllWhere(templateStoreV2.Fields.UserName, u.UserName)
	// 	if err != nil {
	// 		logHandler.ErrorLogger.Printf("Phase %v Error getting users by LastHost: %v", phase, err)
	// 	} else {
	// 		//			logHandler.InfoLogger.Printf("Phase %v Users by LastHost: %v/%v (%v) %v", phase, x+1, len(users), len(dum), u.UserName)
	// 	}
	// }

	// // logHandler.InfoLogger.Printf("Phase %v Flushing Cache", phase)
	// // err = //Cache.SynchroniseForType(templateStoreV2.TemplateStore{})
	// // if err != nil {
	// // 	logHandler.ErrorLogger.Printf("Phase %v Error flushing cache: %v", phase, err)
	// // } else {
	// // 	logHandler.InfoLogger.Printf("Phase %v Cache flushed successfully", phase)
	// // }
	logHandler.InfoLogger.Printf("Phase %v Completed", phase)

	// created, updated, noTables, noCacheEntries := //Cache.Stats()
	// logHandler.InfoLogger.Printf("Cache Stats - Created: %v, Updated: %v, Tables: %v, Entries: %v", created.Format(time.RFC3339Nano), updated.Format(time.RFC3339Nano), noTables, noCacheEntries)

	return uKey

}
