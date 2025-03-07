package v1

import (
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ModelHandler struct {
	ls *service.ModelService
}

func NewModelHandler(ls *service.ModelService) *ModelHandler {
	return &ModelHandler{ls: ls}
}

func (mh *ModelHandler) CreateModel(c *gin.Context) {
	var data model.Model
	if err := c.ShouldBindJSON(&data); err != nil {
		panic(err)
	}
	if err := mh.ls.Create(c, &data); err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(nil))
}

func (mh *ModelHandler) ListModel(c *gin.Context) {
	var query model.ModelQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(err)
	}
	list, total, err := mh.ls.List(c, &query)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponseWithTotal(list, total))
}

func (mh *ModelHandler) GetModelDetail(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(err)
	}
	detail, err := mh.ls.Get(c, id)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(detail))
}
