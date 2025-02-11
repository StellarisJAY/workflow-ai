package v1

import (
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"github.com/StellrisJAY/workflow-ai/internal/middleware"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/StellrisJAY/workflow-ai/internal/service"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Router struct {
	e    *gin.Engine
	conf *config.Config
}

func NewRouter(conf *config.Config) *Router {
	e := gin.New()
	r := &Router{e: e, conf: conf}
	return r
}

func (r *Router) Init() error {
	repository, err := repo.NewRepository(r.conf)
	if err != nil {
		return err
	}
	nodeId, err := strconv.ParseInt(r.conf.Server.Id, 10, 64)
	if err != nil {
		return err
	}
	snowflakeNode, err := snowflake.NewNode(nodeId)
	if err != nil {
		return err
	}
	llmRepo := repo.NewLLMRepo(repository)
	templateRepo := repo.NewTemplateRepo(repository)

	llmService := service.NewLLMService(llmRepo, snowflakeNode)
	templateService := service.NewTemplateService(templateRepo, snowflakeNode)
	llmHandler := NewLLMHandler(llmService)
	templateHandler := NewTemplateHandler(templateService)

	r.e.Use(middleware.Recovery)
	v1 := r.e.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, "pong")
		})
		model := v1.Group("/model")
		{
			model.POST("/create", llmHandler.CreateLLM)
			model.DELETE("/delete", nil)
			model.PUT("/update", nil)
			model.GET("/detail/:id", llmHandler.GetLLMDetail)
			model.GET("/list", llmHandler.ListLLM)
		}
		template := v1.Group("/template")
		{
			template.POST("/create", templateHandler.Create)
			template.DELETE("/delete", nil)
			template.PUT("/update", nil)
			template.GET("/detail/:id", templateHandler.GetDetail)
			template.GET("/list", templateHandler.List)
		}
	}
	return nil
}

func (r *Router) Start() error {
	return r.e.Run(r.conf.Server.Host + ":" + r.conf.Server.Port)
}
