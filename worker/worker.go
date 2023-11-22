package worker

import (
	"encoding/base64"
	"fmt"
	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/model"
	"github.com/conductor-sdk/conductor-go/sdk/settings"
	"github.com/conductor-sdk/conductor-go/sdk/worker"
	log "github.com/sirupsen/logrus"
	"runtime/debug"
	"time"
	"worker-sample/config"
)

type Worker interface {
	SetServiceContext(ctx *config.ServiceContext)
	GetTaskDefName() string
	Execute(t *model.Task) (*model.TaskResult, error)
}

func Wrap(worker Worker) model.ExecuteTaskFunction {
	return func(t *model.Task) (taskResult interface{}, err error) {
		defer func() {
			if e := recover(); e != nil {
				errStack := string(debug.Stack())
				log.Errorf("Uncaught panic at WorkflowInstanceId %s TaskId %s, error: %+v, stack: %s", t.WorkflowInstanceId, t.TaskId, e, errStack)
				tr := model.NewTaskResultFromTask(t)
				tr.Status = model.FailedWithTerminalErrorTask
				tr.ReasonForIncompletion = e.(error).Error()
				tr.Logs = append(
					tr.Logs,
					model.TaskExecLog{
						Log:         errStack,
						TaskId:      t.TaskId,
						CreatedTime: time.Now().UnixMilli(),
					},
				)
				taskResult = tr
			}
		}()

		return worker.Execute(t)
	}
}

var workerMap = make(map[string]Worker)

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
	for _, workerTask := range workerMap {
		workerTask.SetServiceContext(ctx)
		taskRunner.StartWorkerWithDomain(workerTask.GetTaskDefName(), Wrap(workerTask), ctx.Config.Worker.BatchSize, ctx.Config.Worker.PollingInterval, ctx.Config.Worker.Domain)
	}

	return conductorWorker
}

func newAPIClient(ctx *config.ServiceContext) *client.APIClient {
	httpSettings := settings.NewHttpSettings(
		ctx.Config.Worker.BaseUrl,
	)
	if ctx.Config.Worker.Username != "" && ctx.Config.Worker.Password != "" {
		httpSettings.Headers["Authorization"] = getAuthentication(ctx.Config.Worker.Username, ctx.Config.Worker.Password)
	}
	return client.NewAPIClient(
		nil,
		httpSettings)
}

func getAuthentication(username, password string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
}
