package jobs

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"

	"github.com/mt1976/frantic-core/dateHelpers"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/stringHelpers"
	"github.com/mt1976/frantic-core/timing"
)

func StartOfDay(t time.Time) time.Time {
	// Purpose: To remove the time from a date
	return dateHelpers.StartOfDay(t)
}

func DateBeforeOrEqualTo(t1, t2 time.Time) bool {
	return dateHelpers.IsBeforeOrEqualTo(dateHelpers.StartOfDay(t1), dateHelpers.StartOfDay(t2))
}

func DateAfterOrEqualTo(t1, t2 time.Time) bool {
	return dateHelpers.IsAfterOrEqualTo(dateHelpers.StartOfDay(t1), dateHelpers.StartOfDay(t2))
}

func Before(t1, t2 time.Time) bool {
	return dateHelpers.IsBefore(dateHelpers.StartOfDay(t1), dateHelpers.StartOfDay(t2))
}

func DateAfter(t1, t2 time.Time) bool {
	return dateHelpers.IsAfter(dateHelpers.StartOfDay(t1), dateHelpers.StartOfDay(t2))
}

func Equal(t1, t2 time.Time) bool {
	return dateHelpers.StartOfDay(t1).Equal(dateHelpers.StartOfDay(t2))
}

func NextRun(name, schedule string) string {
	// Purpose: To determine the next run time of a job
	rtn := fmt.Sprintf("[%v] NextRun: %v", name, GetHumanReadableCronFreq(schedule))
	logHandler.Service.Println(rtn)
	return rtn
}

// Announce - Announce the start of a job to the log
// Deprecated: Use PreRun instead
func Announce(name, inAction string) {
	// logHandler.ServiceBanner(domain, name, inAction)
}

func GetHumanReadableCronFreq(freq string) string {
	// bkHuman1, _ := crondescriptor.NewCronDescriptor(freq)
	// bkHuman, _ := bkHuman1.GetDescription(crondescriptor.Full)
	nextTime := cronexpr.MustParse(freq).Next(time.Now())
	bkHuman := nextTime.Format("02 Jan 2006 (Mon) 15:04:05")
	return bkHuman
}

func PreRun(job Job) {
	// Purpose: To log the start of a job
	logHandler.Service.Printf("[%v] Started", job.Name())
}

func PostRun(job Job) {
	// Purpose: To log the completion of a job
	nextRun := GetHumanReadableCronFreq(job.Schedule())
	// logHandler.Service.Printf("[%v] Completed", job.Name())
	logHandler.Service.Printf("[%v] Completed - Scheduled (%v) Next Run: %v", job.Name(), job.Schedule(), nextRun)
}

func CodedName(job Job) string {
	// Purpose: To return the coded name of a job
	name := job.Name()
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	name = stringHelpers.RemoveSpecialChars(name)
	return name
}

func AddJobToScheduler(job Job) {
	// logHandler.ServiceLogger.Printf("[%v] Scheduling Job [%v] [%v]", domain, j.Name(), j.Schedule())
	clock := timing.Start(domain, "Schedule", job.Name())
	// Start the job
	jobID, err := scheduledTasks.AddFunc(job.Schedule(), job.Service())
	if err != nil {
		logHandler.Error.Printf("[%v] Scheduling Error=[%v]", job.Name(), err.Error())
		return
	}
	nextRun := GetHumanReadableCronFreq(job.Schedule())
	logHandler.Service.Printf("[%v] Scheduled (%v) Next Run: %v (id=%v)", job.Name(), job.Schedule(), nextRun, jobID)
	clock.Stop(1)
}

func AddJobsToScheduler(jobs []Job) {
	clock := timing.Start(domain, "Schedule", "Jobs")
	// Schedule a list of jobs
	for _, job := range jobs {
		AddJobToScheduler(job)
	}
	clock.Stop(len(jobs))
}

func StartScheduler() {
	clock := timing.Start(domain, "Start", "Scheduler")
	logHandler.Service.Println("Scheduler - Starting")
	// Start the scheduler
	scheduledTasks.Start()

	noEntries := len(scheduledTasks.Entries())
	// Log the scheduled tasks
	// for x, entry := range scheduledTasks.Entries() {
	// 	logHandler.ServiceLogger.Printf("(%v/%v) [%v] [%v] [%v]", x+1, noEntries, entry.ID, entry.Next, entry.Job)
	// }
	logHandler.Service.Println("Scheduler - Started")
	clock.Stop(noEntries + 1)
}
