package jobs

import (
	"github.com/mt1976/frantic-core/commonConfig"

	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
	"github.com/robfig/cron/v3"
)

var (
	domain         = "Schedule"
	scheduledTasks *cron.Cron
	appName        string
)

func Initialise(cfg *commonConfig.Settings) error {
	clock := timing.Start(domain, "Initialise", "")

	logHandler.Info.Println("Initialise - Started")

	scheduledTasks = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))

	appName = cfg.GetApplication_Name()
	logHandler.Info.Println("Initialise - Complete")
	clock.Stop(1)
	return nil
}
