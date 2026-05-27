package schedule_test

import (
	"testing"

	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn

func getClient(t *testing.T) kk_scheduler.KKScheduleClient {
	var err error
	conn, err = grpc.NewClient(
		"127.0.0.1:8666",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	return kk_scheduler.NewKKScheduleClient(conn)
}

func down() {
	conn.Close() // nolint
}

const testAuthToken = "sdgoisdglodshlghlshlghdlskg"

var testJob = func() *kk_scheduler.PBRegisterJob {
	j := &kk_scheduler.PBRegisterJob{}
	j.SetDescription("test job")
	j.SetServiceName("test-service")
	j.SetFuncName("test-func")
	return j
}()

var jobId = "029d4317-039b-7b89-88b1-689d1e0f24c0"

var testService = func() *kk_scheduler.PBRegisterService {
	s := &kk_scheduler.PBRegisterService{}
	s.SetServiceName("test-service")
	s.SetTarget("127.0.0.1:8000")
	s.SetAuthToken(testAuthToken)
	return s
}()

func TestJobList(t *testing.T) {
	defer down()
	input := &kk_scheduler.JobList_Input{}
	jobs, err := getClient(t).JobList(t.Context(), input)
	assert.NoError(t, err)
	for _, job := range jobs.GetJobList() {
		t.Log(job)
		t.Log(job.GetPrev().AsTime(), job.GetNext().AsTime())
	}
}

func TestJobGet(t *testing.T) {
	defer down()
	input := &kk_scheduler.JobGet_Input{}
	input.SetId(jobId)
	job, err := getClient(t).JobGet(t.Context(), input)
	assert.NoError(t, err)
	t.Log(job.GetJob())
	t.Log(job.GetJob().GetPrev().AsTime(), job.GetJob().GetNext().AsTime())
}

func TestJobSetSpec(t *testing.T) {
	defer down()
	input := &kk_scheduler.JobSetSpec_Input{}
	input.SetId(jobId)
	input.SetSpec("* * * * *")
	resp, err := getClient(t).JobSetSpec(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestJobEnable(t *testing.T) {
	defer down()
	input := &kk_scheduler.JobEnable_Input{}
	input.SetId(jobId)
	resp, err := getClient(t).JobEnable(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestJobDisable(t *testing.T) {
	defer down()
	input := &kk_scheduler.JobDisable_Input{}
	input.SetId(jobId)
	resp, err := getClient(t).JobDisable(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestJobPut(t *testing.T) {
	defer down()
	input := &kk_scheduler.JobPut_Input{}
	input.SetJob(testJob)
	resp, err := getClient(t).JobPut(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestServiceList(t *testing.T) {
	defer down()
	input := &kk_scheduler.ServiceList_Input{}
	resp, err := getClient(t).ServiceList(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestServicePut(t *testing.T) {
	defer down()
	input := &kk_scheduler.ServicePut_Input{}
	input.SetService(testService)
	resp, err := getClient(t).ServicePut(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestServiceGet(t *testing.T) {
	defer down()
	input := &kk_scheduler.ServiceGet_Input{}
	input.SetServiceName(testService.GetServiceName())
	resp, err := getClient(t).ServiceGet(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestServiceDelete(t *testing.T) {
	defer down()
	input := &kk_scheduler.ServiceDelete_Input{}
	input.SetServiceName(testService.GetServiceName())
	resp, err := getClient(t).ServiceDelete(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}
