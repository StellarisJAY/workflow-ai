package v1

import (
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/service"
	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	service *service.ProviderService
}

func NewProviderHandler(service *service.ProviderService) *ProviderHandler {
	return &ProviderHandler{service: service}
}

func (p *ProviderHandler) ListProviderSchemas(c *gin.Context) {
	c.JSON(200, common.NewSuccessResponse(model.ProviderSchemas))
}

func (p *ProviderHandler) CreateProvider(c *gin.Context) {
	var provider model.Provider
	if err := c.ShouldBindJSON(&provider); err != nil {
		panic(err)
	}
	if err := p.service.CreateProvider(c.Request.Context(), &provider); err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(nil))
}

func (p *ProviderHandler) CreateProviderModel(c *gin.Context) {
	var pm model.ProviderModel
	if err := c.ShouldBindJSON(&pm); err != nil {
		panic(err)
	}
	if err := p.service.CreateProviderModel(c.Request.Context(), &pm); err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(nil))
}

func (p *ProviderHandler) ListProviderModel(c *gin.Context) {
	var query model.ProviderModelQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(err)
	}
	list, total, err := p.service.ListProviderModels(c, &query)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponseWithTotal(list, total))
}

func (p *ProviderHandler) ListProviders(c *gin.Context) {
	list, err := p.service.ListProviders(c)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(list))
}
