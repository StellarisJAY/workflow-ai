package v1

import (
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type LLMHandler struct {
	ls *service.LLMService
}

func NewLLMHandler(ls *service.LLMService) *LLMHandler {
	return &LLMHandler{ls: ls}
}

func (lh *LLMHandler) CreateLLM(c *gin.Context) {
	var data model.LLM
	if err := c.ShouldBindJSON(&data); err != nil {
		panic(err)
	}
	if err := lh.ls.Create(c, &data); err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(nil))
}

func (lh *LLMHandler) ListLLM(c *gin.Context) {
	var query model.LLMQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(err)
	}
	list, total, err := lh.ls.List(c, &query)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponseWithTotal(list, total))
}

func (lh *LLMHandler) GetLLMDetail(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(err)
	}
	detail, err := lh.ls.Get(c, id)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(detail))
}
