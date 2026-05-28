package store_driver

import (
	"errors"
	"fmt"
	"log/slog"

	"gitee.com/cruvie/kk_kit/go/kk_protobuf"
	"gitee.com/cruvie/kk_kit/go/kk_time"
	"gitee.com/cruvie/kk_kit/go/multi_lang/kk_id"
	"github.com/cruvie/kk_scheduler/go/internal/models"
	"github.com/cruvie/kk_scheduler/go/internal/models/query"
	"github.com/cruvie/kk_scheduler/go/kk_scheduler"
	"gorm.io/gorm"
)

// PostgresStore implements StoreDriver using PostgreSQL
type PostgresStore struct{}

// NewPostgresStore creates a new PostgreSQL store
func NewPostgresStore() *PostgresStore {
	return &PostgresStore{}
}

// TaskCreate creates a new task execution record
func (s *PostgresStore) TaskCreate(in *kk_scheduler.TaskCreate_Input) error {
	{
		// check
		_, err := s.JobGet(in.GetJobId())
		if err != nil {
			return err
		}
	}
	execution := &models.TaskExecution{
		Id:         in.GetId(),
		JobId:      in.GetJobId(),
		Status:     kk_scheduler.TaskExecutionStatus_TASK_EXECUTION_STATUS_Init,
		StartedAt:  kk_time.NowUTCTime(),
		FinishedAt: kk_time.DefaultTime,
		Log:        "",
	}

	if err := query.TaskExecution.Create(execution); err != nil {
		slog.Error("failed to create task execution", "err", err)
		return err
	}
	return nil
}

// TaskUpdateStatus updates status
func (s *PostgresStore) TaskUpdateStatus(id string, status kk_scheduler.TaskExecutionStatus) error {
	if status == kk_scheduler.TaskExecutionStatus_TASK_EXECUTION_STATUS_COMPLETED {
		_, err := query.TaskExecution.
			Where(query.TaskExecution.Id.Eq(id)).
			UpdateSimple(
				query.TaskExecution.Status.Value(int32(status)),
				query.TaskExecution.FinishedAt.Value(kk_time.NowUTCTime()),
			)
		if err != nil {
			return err
		}
	}
	_, err := query.TaskExecution.
		Where(query.TaskExecution.Id.Eq(id)).
		UpdateSimple(query.TaskExecution.Status.Value(int32(status)))
	if err != nil {
		return err
	}
	return err
}

// TaskAppendLog appends log to the execution record
func (s *PostgresStore) TaskAppendLog(id string, log string) error {
	err := query.Q.Transaction(func(tx *query.Query) error {
		execution, err := tx.TaskExecution.
			Where(tx.TaskExecution.Id.Eq(id)).
			Take()
		if err != nil {
			return err
		}

		_, err = tx.TaskExecution.
			Where(tx.TaskExecution.Id.Eq(id)).
			UpdateSimple(tx.TaskExecution.Log.Value(fmt.Sprintf("%s%s\n", execution.Log, log)))
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

// TaskExecutionList lists task execution records
func (s *PostgresStore) TaskExecutionList(jobId string) ([]*kk_scheduler.PBTaskExecution, error) {
	q := query.TaskExecution.
		Omit(query.TaskExecution.Log).
		Where()
	if jobId != "" {
		q = query.TaskExecution.Where(query.TaskExecution.JobId.Eq(jobId))
	}

	executions, err := q.Order(query.TaskExecution.StartedAt.Desc()).Find()
	if err != nil {
		slog.Error("failed to list task executions", "err", err, "jobId", jobId)
		return nil, err
	}
	list := kk_protobuf.ToPBList(executions)
	return list, nil
}

// TaskExecutionGet returns a specific task execution record
func (s *PostgresStore) TaskExecutionGet(id string) (*kk_scheduler.PBTaskExecution, error) {
	execution, err := query.TaskExecution.
		Where(query.TaskExecution.Id.Eq(id)).
		Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("failed to get task execution", "err", err, "id", id)
		return nil, err
	}
	return execution.ToPB(), nil
}

// TaskExecutionDelete deletes a task execution record
func (s *PostgresStore) TaskExecutionDelete(id string) error {
	_, err := query.TaskExecution.
		Where(query.TaskExecution.Id.Eq(id)).
		Delete()
	if err != nil {
		slog.Error("failed to delete task execution", "err", err, "id", id)
		return err
	}
	return nil
}

// JobList returns all jobs for a service
func (s *PostgresStore) JobList(serviceName string) ([]*kk_scheduler.PBJob, error) {
	q := query.Job.Where()
	if serviceName != "" {
		q = query.Job.Where(query.Job.ServiceName.Eq(serviceName))
	}

	jobs, err := q.Find()
	if err != nil {
		slog.Error("failed to list jobs", "err", err, "serviceName", serviceName)
		return nil, err
	}
	return kk_protobuf.ToPBList(jobs), nil
}

// JobGet returns a specific job
func (s *PostgresStore) JobGet(jobId string) (*kk_scheduler.PBJob, error) {
	job, err := query.Job.
		Where(query.Job.Id.Eq(jobId)).
		Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("failed to get job", "err", err, "jobId", jobId)
		return nil, err
	}
	return job.ToPB(), nil
}

func (s *PostgresStore) JobGetByServiceFuncName(serviceName, funcName string) (*kk_scheduler.PBJob, error) {
	job, err := query.Job.
		Where(query.Job.ServiceName.Eq(serviceName)).
		Where(query.Job.FuncName.Eq(funcName)).
		Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("failed to get job", "err", err, "serviceName", serviceName, "funcName", funcName)
		return nil, err
	}
	return job.ToPB(), nil
}

// JobDelete deletes a job
func (s *PostgresStore) JobDelete(jobId string) error {
	_, err := query.Job.
		Where(query.Job.Id.Eq(jobId)).
		Delete()
	if err != nil {
		slog.Error("failed to delete job", "err", err, "jobId", jobId)
		return err
	}
	return nil
}

// JobPut creates or updates a job (upsert)
func (s *PostgresStore) JobPut(entry *kk_scheduler.PBJob) error {
	// Check if job exists
	existing, err := query.Job.
		Where(query.Job.ServiceName.Eq(entry.GetServiceName())).
		Where(query.Job.FuncName.Eq(entry.GetFuncName())).
		Take()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new job
		job := &models.Job{}
		job.FromPB(entry)
		job.Id = kk_id.GenUUID7()
		if err = query.Job.Create(job); err != nil {
			slog.Error("failed to create job", "err", err)
			return err
		}
		return nil
	}

	if err != nil {
		slog.Error("failed to check existing job", "err", err)
		return err
	}

	// Update existing job
	existing.FromPB(entry)
	if err = query.Job.Save(existing); err != nil {
		slog.Error("failed to update job", "err", err)
		return err
	}
	return nil
}

// ServiceList returns all registered services
func (s *PostgresStore) ServiceList() ([]*kk_scheduler.PBRegisterService, error) {
	services, err := query.Service.Find()
	if err != nil {
		slog.Error("failed to list services", "err", err)
		return nil, err
	}

	return kk_protobuf.ToPBList(services), nil
}

// ServiceGet returns a specific service
func (s *PostgresStore) ServiceGet(serviceName string) (*kk_scheduler.PBRegisterService, error) {
	service, err := query.Service.
		Where(query.Service.ServiceName.Eq(serviceName)).
		Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("failed to get service", "err", err, "serviceName", serviceName)
		return nil, err
	}
	return service.ToPB(), nil
}

// ServicePut creates or updates a service
func (s *PostgresStore) ServicePut(service *kk_scheduler.PBRegisterService) error {
	existing, err := query.Service.
		Where(query.Service.ServiceName.Eq(service.GetServiceName())).
		Take()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		svc := &models.Service{}
		svc.FromPB(service)
		if err = query.Service.Create(svc); err != nil {
			slog.Error("failed to create service", "err", err)
			return err
		}
		return nil
	}

	if err != nil {
		slog.Error("failed to check existing service", "err", err)
		return err
	}

	existing.FromPB(service)
	if err = query.Service.Save(existing); err != nil {
		slog.Error("failed to update service", "err", err)
		return err
	}
	return nil
}

// ServiceDelete deletes a service
func (s *PostgresStore) ServiceDelete(serviceName string) error {
	_, err := query.Service.
		Where(query.Service.ServiceName.Eq(serviceName)).
		Delete()
	if err != nil {
		slog.Error("failed to delete service", "err", err, "serviceName", serviceName)
		return err
	}
	return nil
}
