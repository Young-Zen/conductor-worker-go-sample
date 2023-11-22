package worker

import (
	"context"
	"github.com/conductor-sdk/conductor-go/sdk/model"
	log "github.com/sirupsen/logrus"
	"worker-sample/config"
)

type StartConductorWorkflowWorker struct {
	ServiceContext *config.ServiceContext
}

func init() {
	worker := NewStartConductorWorkflowWorker()
	workerMap[worker.GetTaskDefName()] = worker
}

func NewStartConductorWorkflowWorker() *StartConductorWorkflowWorker {
	return &StartConductorWorkflowWorker{}
}

func (w *StartConductorWorkflowWorker) SetServiceContext(ctx *config.ServiceContext) {
	w.ServiceContext = ctx
}

func (w *StartConductorWorkflowWorker) GetTaskDefName() string {
	return "DEMO_START_CONDUCTOR_WORKFLOW"
}

func (w *StartConductorWorkflowWorker) Execute(t *model.Task) (*model.TaskResult, error) {
	log.Infof("WorkflowInstanceId: %s, TaskId: %s, Type: %s, TDN: %s", t.WorkflowInstanceId, t.TaskId, t.TaskType, t.TaskDefName)

	taskResult := model.NewTaskResultFromTask(t)
	workflowName, ok := t.InputData["workflowName"].(string)
	if !ok {
		taskResult.Status = model.FailedWithTerminalErrorTask
		taskResult.ReasonForIncompletion = "Input param 'workflowName' error"
		return taskResult, nil
	}

	workflowId, _, err := w.ServiceContext.Worker.WorkflowClient.StartWorkflow(context.Background(), make(map[string]interface{}), workflowName, nil)
	if err != nil {
		return nil, err
	}

	taskResult.Status = model.CompletedTask
	taskResult.OutputData = map[string]interface{}{
		"workflowExecutionId": workflowId,
	}
	return taskResult, nil
}
