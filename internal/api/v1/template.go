package v1

import (
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type TemplateHandler struct {
	service *service.TemplateService
}

func NewTemplateHandler(service *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{service: service}
}

func (t *TemplateHandler) Create(c *gin.Context) {
	var template model.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		panic(err)
	}
	id, err := t.service.Insert(c, &template)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(struct {
		Id int64 `json:"id,string"`
	}{id}))
}

func (t *TemplateHandler) GetDetail(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(err)
	}
	detail, err := t.service.Get(c, id)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(detail))
}

func (t *TemplateHandler) List(c *gin.Context) {
	var query model.TemplateQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(err)
	}
	list, err := t.service.List(c, &query)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponseWithTotal(list, len(list)))
}

func (t *TemplateHandler) Delete(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(err)
	}
	if err := t.service.Delete(c, id); err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(nil))
}

func (t *TemplateHandler) Update(c *gin.Context) {
	var template model.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		panic(err)
	}
	if err := t.service.Update(c, &template); err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(nil))
}
