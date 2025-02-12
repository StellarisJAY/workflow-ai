package v1

import (
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/service"
	"github.com/gin-gonic/gin"
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
	err := w.service.Start(c, &request)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(nil))
}
