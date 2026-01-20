package jobs

import "github.com/mt1976/frantic-amphora/dao/database"

type Job interface {
	Run() error
	Service() func()
	Schedule() string
	Name() string
	AddDatabaseAccessFunctions(f func() ([]*database.DB, error))
	Description() string
}
