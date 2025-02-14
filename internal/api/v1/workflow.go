package v1

import (
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type WorkflowHandler struct {
	service *service.WorkflowService
}

func NewWorkflowHandler(service *service.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{
		service: service,
	}
}

func (w *WorkflowHandler) Start(c *gin.Context) {
	var request model.StartWorkflowRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(err)
	}
	workflowId, err := w.service.Start(c, &request)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(struct {
		WorkflowId int64 `json:"workflowId,string"`
	}{WorkflowId: workflowId}))
}

func (w *WorkflowHandler) Outputs(c *gin.Context) {
	param := c.Param("id")
	workflowId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(err)
	}
	outputs, err := w.service.Outputs(c, workflowId)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(outputs))
}
