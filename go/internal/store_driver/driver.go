package store_driver

import (
	"github.com/cruvie/kk-scheduler/go/internal/g_config"
	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
)

type StoreDriver interface {
	JobList(serviceName string) ([]*kk_scheduler.PBJob, error)
	JobGet(jobId string) (*kk_scheduler.PBJob, error)
	JobGetByServiceFuncName(serviceName, funcName string) (*kk_scheduler.PBJob, error)
	JobDelete(jobId string) error
	JobPut(job *kk_scheduler.PBJob) error

	ServiceList() ([]*kk_scheduler.PBRegisterService, error)
	ServicePut(service *kk_scheduler.PBRegisterService) error
	ServiceGet(serviceName string) (*kk_scheduler.PBRegisterService, error)
	ServiceDelete(serviceName string) error

	// TaskCreate creates a new task execution record
	TaskCreate(in *kk_scheduler.TaskCreate_Input) error
	// TaskUpdateStatus updates the status
	TaskUpdateStatus(id string, status kk_scheduler.TaskExecutionStatus) error
	// TaskAppendLog append log to the execution record
	TaskAppendLog(id string, log string) error
	// TaskExecutionList lists task execution records, optionally filtered by jobId
	TaskExecutionList(jobId string) ([]*kk_scheduler.PBTaskExecution, error)
	// TaskExecutionGet returns a specific task execution record
	TaskExecutionGet(id string) (*kk_scheduler.PBTaskExecution, error)
	// TaskExecutionDelete deletes a task execution record
	TaskExecutionDelete(id string) error
}

func NewStoreDriver() StoreDriver {
	switch g_config.Config.Store.Choose {
	case "PG":
		return NewPostgresStore()
	default:
		panic("store choose error")
	}
}
