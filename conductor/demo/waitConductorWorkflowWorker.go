package demo

import (
	"context"
	"fmt"
	"github.com/conductor-sdk/conductor-go/sdk/model"
	log "github.com/sirupsen/logrus"
	cw "worker-sample/conductor/worker"
	"worker-sample/config"
)

type WaitConductorWorkflowWorker struct {
	ServiceContext *config.ServiceContext
}

func init() {
	waitConductorWorkflowWorker := NewWaitConductorWorkflowWorker()
	cw.Register(waitConductorWorkflowWorker)
}

func NewWaitConductorWorkflowWorker() *WaitConductorWorkflowWorker {
	return &WaitConductorWorkflowWorker{}
}

func (w *WaitConductorWorkflowWorker) SetServiceContext(ctx *config.ServiceContext) {
	w.ServiceContext = ctx
}

func (w *WaitConductorWorkflowWorker) GetTaskDefName() string {
	return "DEMO_WAIT_CONDUCTOR_WORKFLOW"
}

func (w *WaitConductorWorkflowWorker) Execute(t *model.Task) (*model.TaskResult, error) {
	log.Infof("WorkflowInstanceId: %s, TaskId: %s, Type: %s, TDN: %s", t.WorkflowInstanceId, t.TaskId, t.TaskType, t.TaskDefName)

	taskResult := model.NewTaskResultFromTask(t)
	workflowExecutionId, ok := t.InputData["workflowExecutionId"].(string)
	if !ok {
		taskResult.Status = model.FailedWithTerminalErrorTask
		taskResult.ReasonForIncompletion = "Input param 'workflowExecutionId' error"
		return taskResult, nil
	}

	workflow, _, err := w.ServiceContext.Worker.WorkflowClient.GetExecutionStatus(context.Background(), workflowExecutionId, nil)
	if err != nil {
		return model.NewTaskResultFromTaskWithError(t, err), err
	}

	taskResult.Status = model.InProgressTask
	switch workflow.Status {
	case model.RunningWorkflow, model.PausedWorkflow:
		taskResult.CallbackAfterSeconds = 10
	case model.CompletedWorkflow:
		taskResult.Status = model.CompletedTask
	case model.TerminatedWorkflow:
		taskResult.Status = model.FailedWithTerminalErrorTask
		taskResult.ReasonForIncompletion = fmt.Sprintf("Workflow %s execute failed: %s", workflowExecutionId, workflow.ReasonForIncompletion)
	default:
		taskResult.Status = model.FailedTask
		taskResult.ReasonForIncompletion = fmt.Sprintf("Workflow %s execute failed: %s", workflowExecutionId, workflow.ReasonForIncompletion)
	}
	return taskResult, nil
}
