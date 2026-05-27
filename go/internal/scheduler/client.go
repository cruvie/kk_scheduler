package scheduler

import (
	"errors"
	"slices"
	"strings"

	"gitee.com/cruvie/kk_kit/go/kk_time"
	"gitee.com/cruvie/kk_kit/go/multi_lang/kk_id"
	"github.com/cruvie/kk-scheduler/go/internal/store_driver"
	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GClient *Client

type Client struct {
	cron   *cron.Cron
	storer store_driver.StoreDriver
}

func InitGClient(cfg *Config) {
	cfg.check()
	c := &Client{
		cron:   cron.New(cfg.Opts...),
		storer: cfg.StoreDriver,
	}
	GClient = c
	GClient.initJob()
}

func (x *Client) initJob() {
	jobList, err := x.JobList("")
	if err != nil {
		panic(err)
	}
	for _, job := range jobList {
		if !job.GetEnabled() {
			continue
		}
		err := x.JobEnable(job.GetId())
		if err != nil {
			panic(err)
		}
	}
}

func (x *Client) Start() {
	x.cron.Start()
}

func (x *Client) Close() {
	x.cron.Stop()
}

func (x *Client) JobPut(jobs ...*kk_scheduler.PBRegisterJob) error {
	for _, job := range jobs {
		err := job.Check()
		if err != nil {
			panic(err)
		}
		{ // check service exist
			_, err = x.storer.ServiceGet(job.GetServiceName())
			if err != nil {
				return err
			}
		}

		newEntry := &kk_scheduler.PBJob{}
		newEntry.SetId(kk_id.GenUUID7())
		newEntry.SetEntryID(0)
		newEntry.SetEnabled(false)
		newEntry.SetNext(nil)
		newEntry.SetPrev(nil)
		newEntry.SetSpec("")
		newEntry.SetServiceName(job.GetServiceName())
		newEntry.SetDescription(job.GetDescription())
		newEntry.SetFuncName(job.GetFuncName())

		entry, err := x.storer.JobGetByServiceFuncName(job.GetServiceName(), job.GetFuncName())
		if err != nil && !errors.Is(err, kk_scheduler.ErrJobNotFount) {
			return err
		}

		if entry != nil {
			newEntry.SetId(entry.GetId())
			newEntry.SetEntryID(entry.GetEntryID())
			newEntry.SetEnabled(entry.GetEnabled())
			newEntry.SetNext(entry.GetNext())
			newEntry.SetPrev(entry.GetPrev())
			newEntry.SetSpec(entry.GetSpec())
		}
		err = x.storer.JobPut(newEntry)
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *Client) JobList(serviceName string) ([]*kk_scheduler.PBJob, error) {
	entries := x.cron.Entries()
	pbJobList, err := x.storer.JobList(serviceName)
	if err != nil {
		return nil, err
	}
	var hasSpecEntryList []int32
	var pbJobs []*kk_scheduler.PBJob
	for _, entry := range entries {
		dbPBJob, b := lo.Find(pbJobList, func(item *kk_scheduler.PBJob) bool {
			return item.GetEntryID() == int32(entry.ID)
		})
		if !b {
			continue
		} else {
			hasSpecEntryList = append(hasSpecEntryList, dbPBJob.GetEntryID())
		}
		job := &kk_scheduler.PBJob{}
		job.SetId(dbPBJob.GetId())
		job.SetEntryID(dbPBJob.GetEntryID())
		job.SetEnabled(dbPBJob.GetEnabled())
		job.SetNext(timestamppb.New(entry.Next))
		if entry.Prev.IsZero() {
			entry.Prev = kk_time.DefaultTime
		}
		job.SetPrev(timestamppb.New(entry.Prev))
		job.SetSpec(dbPBJob.GetSpec())
		job.SetDescription(dbPBJob.GetDescription())
		job.SetFuncName(dbPBJob.GetFuncName())
		job.SetServiceName(dbPBJob.GetServiceName())
		pbJobs = append(pbJobs, job)
	}

	noSpecJobList := lo.Filter(pbJobList, func(item *kk_scheduler.PBJob, index int) bool {
		_, b := lo.Find(hasSpecEntryList, func(id int32) bool {
			return id == item.GetEntryID()
		})
		return !b
	})
	pbJobs = append(pbJobs, noSpecJobList...)
	// sort by serviceName
	slices.SortFunc(pbJobs, func(a, b *kk_scheduler.PBJob) int {
		return strings.Compare(a.GetServiceName(), b.GetServiceName())
	})
	return pbJobs, nil
}

func (x *Client) JobGet(jobId string) (*kk_scheduler.PBJob, error) {
	entry, err := x.storer.JobGet(jobId)
	if err != nil {
		return nil, err
	}
	cEntry := x.cron.Entry(cron.EntryID(entry.GetEntryID()))

	if cEntry.Valid() {
		entry.SetNext(timestamppb.New(cEntry.Next))
		entry.SetPrev(timestamppb.New(cEntry.Prev))
	}

	return entry, nil
}

func (x *Client) JobSetSpec(jobId string, spec string) error {
	job, err := x.storer.JobGet(jobId)
	if err != nil {
		return err
	}
	if job.GetSpec() == spec {
		return nil
	}
	job.SetSpec(spec)

	err = x.storer.JobPut(job)
	if job.GetEnabled() {
		err = x.JobEnable(jobId)
		if err != nil {
			return err
		}
	}

	return err
}

func (x *Client) JobEnable(jobId string) error {
	entry, err := x.storer.JobGet(jobId)
	if err != nil {
		return err
	}
	if entry.GetSpec() == "" {
		return kk_scheduler.ErrSpecIsEmpty
	}
	err = x.JobDisable(jobId)
	if err != nil {
		return err
	}
	service, err := x.storer.ServiceGet(entry.GetServiceName())
	if err != nil {
		return err
	}

	entryID, err := x.cron.AddFunc(entry.GetSpec(), triggerFunc(service, entry))
	if err != nil {
		return err
	}

	entry.SetEntryID(int32(entryID))
	entry.SetEnabled(true)

	err = x.storer.JobPut(entry)
	return err
}

func (x *Client) JobDisable(jobId string) error {
	entry, err := x.storer.JobGet(jobId)
	if err != nil {
		return err
	}
	x.cron.Remove(cron.EntryID(entry.GetEntryID()))

	entry.SetEnabled(false)
	entry.SetEntryID(0)

	err = x.storer.JobPut(entry)
	return err
}

func (x *Client) JobDelete(jobId string) error {
	// disable job
	err := x.JobDisable(jobId)
	if err != nil {
		return err
	}
	return x.storer.JobDelete(jobId)
}

// JobTrigger triggers a job manually
func (x *Client) JobTrigger(jobId string) error {
	entry, err := x.storer.JobGet(jobId)
	if err != nil {
		return err
	}
	service, err := x.storer.ServiceGet(entry.GetServiceName())
	if err != nil {
		return err
	}

	// Trigger the job function directly
	triggerFunc(service, entry)()
	return nil
}

func (x *Client) ServiceList() ([]*kk_scheduler.PBRegisterService, error) {
	return x.storer.ServiceList()
}

func (x *Client) ServicePut(service *kk_scheduler.PBRegisterService) error {
	return x.storer.ServicePut(service)
}

func (x *Client) ServiceGet(serviceName string) (*kk_scheduler.PBRegisterService, error) {
	return x.storer.ServiceGet(serviceName)
}

func (x *Client) ServiceDelete(serviceName string) error {
	// check no job in service
	jobList, err := x.JobList(serviceName)
	if err != nil {
		return err
	}
	if len(jobList) > 0 {
		return kk_scheduler.ErrServiceHasJob
	}
	return x.storer.ServiceDelete(serviceName)
}

// TaskCreate creates a new task execution record
func (x *Client) TaskCreate(in *kk_scheduler.TaskCreate_Input) error {
	return x.storer.TaskCreate(in)
}

// TaskUpdateStatus updates the task execution status
func (x *Client) TaskUpdateStatus(id string, status kk_scheduler.TaskExecutionStatus) error {
	return x.storer.TaskUpdateStatus(id, status)
}

// TaskAppendLog appends log to the task execution record
func (x *Client) TaskAppendLog(id string, log string) error {
	return x.storer.TaskAppendLog(id, log)
}

// TaskExecutionList lists task execution records
func (x *Client) TaskExecutionList(jobId string) ([]*kk_scheduler.PBTaskExecution, error) {
	return x.storer.TaskExecutionList(jobId)
}

// TaskExecutionGet returns a specific task execution record
func (x *Client) TaskExecutionGet(id string) (*kk_scheduler.PBTaskExecution, error) {
	return x.storer.TaskExecutionGet(id)
}

// TaskExecutionDelete deletes a task execution record
func (x *Client) TaskExecutionDelete(id string) error {
	return x.storer.TaskExecutionDelete(id)
}
