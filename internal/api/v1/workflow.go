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

func (w *WorkflowHandler) List(c *gin.Context) {
	list, total, err := w.service.ListWorkflowInstance(c)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponseWithTotal(list, total))
}

func (w *WorkflowHandler) GetDetail(c *gin.Context) {
	workflowId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	detail, err := w.service.GetWorkflowInstanceDetail(c, workflowId)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(detail))
}

func (w *WorkflowHandler) GetNodeInstanceDetail(c *gin.Context) {
	var query struct {
		WorkflowId int64  `form:"workflowId"`
		NodeId     string `form:"nodeId"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(err)
	}
	instance, err := w.service.GetNodeInstance(c, query.WorkflowId, query.NodeId)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(instance))
}

func (w *WorkflowHandler) GetWorkflowTimeline(c *gin.Context) {
	workflowId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	timeline, err := w.service.GetWorkflowTimeline(c, workflowId)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(timeline))
}
