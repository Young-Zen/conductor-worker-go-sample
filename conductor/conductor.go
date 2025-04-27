package conductor

import (
	"encoding/base64"
	"fmt"

	_ "worker-sample/conductor/demo"
	cw "worker-sample/conductor/worker"
	"worker-sample/config"

	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/settings"
	"github.com/conductor-sdk/conductor-go/sdk/worker"
)

func InitWorker(ctx *config.ServiceContext) *config.ConductorWorker {
	apiClient := newAPIClient(ctx)
	workflowClient := &client.WorkflowResourceApiService{
		APIClient: apiClient,
	}
	taskRunner := worker.NewTaskRunnerWithApiClient(apiClient)
	conductorWorker := &config.ConductorWorker{
		APIClient:      apiClient,
		WorkflowClient: workflowClient,
		TaskRunner:     taskRunner,
	}
	ctx.Worker = conductorWorker

	// Between polls if there are no tasks available to execute
	workerMap := cw.GetAllRegisteredWorker()
	for _, workerTask := range workerMap {
		workerTask.SetServiceContext(ctx)
		err := taskRunner.StartWorkerWithDomain(
			workerTask.GetTaskDefName(),
			cw.Wrap(workerTask),
			ctx.Config.Worker.BatchSize,
			ctx.Config.Worker.PollingInterval,
			ctx.Config.Worker.Domain,
		)
		if err != nil {
			panic(err)
		}
	}

	return conductorWorker
}

func newAPIClient(ctx *config.ServiceContext) *client.APIClient {
	httpSettings := settings.NewHttpSettings(
		ctx.Config.Worker.BaseUrl,
	)
	if ctx.Config.Worker.Username != "" && ctx.Config.Worker.Password != "" {
		httpSettings.Headers["Authorization"] = getAuthentication(
			ctx.Config.Worker.Username,
			ctx.Config.Worker.Password,
		)
	}
	return client.NewAPIClient(
		nil,
		httpSettings)
}

func getAuthentication(username, password string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
}
