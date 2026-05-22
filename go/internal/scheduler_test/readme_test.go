package schedule_test

import (
	"testing"

	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestForREADME(t *testing.T) {
	// create a client for kk-scheduler
	conn, err := grpc.NewClient("127.0.0.1:8666",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer conn.Close() //nolint
	if err != nil {
		t.Fatal(err)
	}
	client := kk_scheduler.NewKKScheduleClient(conn)

	myServiceName := "my-service"
	testJob := func() *kk_scheduler.PBRegisterJob {
		j := &kk_scheduler.PBRegisterJob{}
		j.SetDescription("test job")
		j.SetServiceName(myServiceName)
		j.SetFuncName("Func1")
		return j
	}()
	testService := func() *kk_scheduler.PBRegisterService {
		s := &kk_scheduler.PBRegisterService{}
		s.SetServiceName(myServiceName)
		s.SetTarget("127.0.0.1:8000")
		return s
	}()
	{
		// put the running service info to kk-scheduler
		input := &kk_scheduler.ServicePut_Input{}
		input.SetService(testService)
		resp, err := client.ServicePut(t.Context(), input)
		assert.NoError(t, err)
		t.Log(resp)
	}
	{
		// put a job to kk-scheduler with the service name
		input := &kk_scheduler.JobPut_Input{}
		input.SetJob(testJob)
		resp, err := client.JobPut(t.Context(), input)
		assert.NoError(t, err)
		t.Log(resp)
	}
	{
		// set job spec
		input := &kk_scheduler.JobSetSpec_Input{}
		input.SetId(jobId)
		input.SetSpec("* * * * *")
		resp, err := client.JobSetSpec(t.Context(), input)
		assert.NoError(t, err)
		t.Log(resp)
	}
	{
		// enable job to be triggered with the spec
		input := &kk_scheduler.JobEnable_Input{}
		input.SetId(jobId)
		resp, err := client.JobEnable(t.Context(), input)
		assert.NoError(t, err)
		t.Log(resp)
	}
}
