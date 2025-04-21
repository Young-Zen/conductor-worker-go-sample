package conductor

import (
	"github.com/conductor-sdk/conductor-go/sdk/model"
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

func GetAllRegisteredWorker() map[string]Worker {
	return workerMap
}

func Register(worker Worker) {
	workerMap[worker.GetTaskDefName()] = worker
	log.Infof("Worker '%s' Registered.", worker.GetTaskDefName())
}
