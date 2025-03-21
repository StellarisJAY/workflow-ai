package v1

import (
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"github.com/StellrisJAY/workflow-ai/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

func (r *Router) Init(templateHandler *TemplateHandler, workflowHandler *WorkflowHandler, kbHandler *KnowledgeBaseHandler,
	fileHandler *FSHandler, providerHandler *ProviderHandler) error {
	r.e.Use(middleware.Recovery)
	r.e.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
	}))
	v1 := r.e.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, "pong")
		})
		provider := v1.Group("/provider")
		{
			provider.GET("/schemas", providerHandler.ListProviderSchemas)
			provider.POST("/create", providerHandler.CreateProvider)
			provider.POST("/model/create", providerHandler.CreateProviderModel)
			provider.GET("/model/list", providerHandler.ListProviderModel)
			provider.GET("/list", providerHandler.ListProviders)
		}
		template := v1.Group("/template")
		{
			template.POST("/create", templateHandler.Create)
			template.DELETE("/:id", templateHandler.Delete)
			template.PUT("/update", templateHandler.Update)
			template.GET("/detail/:id", templateHandler.GetDetail)
			template.GET("/list", templateHandler.List)
			template.GET("/start-variables/:id", templateHandler.GetStartInputVariables)
			template.GET("/prototype", templateHandler.GetNodePrototype)
		}
		wf := v1.Group("/workflow")
		{
			wf.POST("/start", workflowHandler.Start)
			wf.GET("/detail/:id", workflowHandler.GetDetail)
			wf.GET("/list", workflowHandler.List)
			wf.GET("/node/detail", workflowHandler.GetNodeInstanceDetail)
			wf.POST("/start-and-listen", workflowHandler.StartAndListen)
		}
		kb := v1.Group("/knowledgeBase")
		{
			kb.POST("/create", kbHandler.Create)
			kb.PUT("/update", kbHandler.Update)
			kb.GET("/detail/:id", kbHandler.Detail)
			kb.GET("/files/:kbId", kbHandler.ListFiles)
			kb.GET("/list", kbHandler.List)
			kb.POST("/upload", kbHandler.Upload)
			kb.POST("/upload-batch", kbHandler.BatchUpload)
			kb.DELETE("/file/:id", kbHandler.Delete)
			kb.GET("/download/:id", kbHandler.DownloadFile)
			kb.POST("/process/start/:id", kbHandler.StartFileProcessing)
			kb.POST("/similarity-search", kbHandler.SimilaritySearch)
			kb.POST("/fulltext-search", kbHandler.FulltextSearch)
			kb.POST("/hybrid-search", kbHandler.HybridSearch)
			kb.GET("/chunks", kbHandler.ListChunks)
		}
		file := v1.Group("/fs")
		{
			file.POST("/upload", fileHandler.Upload)
			file.POST("/upload-batch", fileHandler.BatchUpload)
			file.GET("/download/:id", fileHandler.Download)
		}
	}
	return nil
}

func (r *Router) Start() error {
	return r.e.Run(r.conf.Server.Host + ":" + r.conf.Server.Port)
}
