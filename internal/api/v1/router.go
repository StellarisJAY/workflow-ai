package v1

import (
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"github.com/StellrisJAY/workflow-ai/internal/middleware"
	"github.com/StellrisJAY/workflow-ai/internal/rag"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/StellrisJAY/workflow-ai/internal/repo/fs"
	"github.com/StellrisJAY/workflow-ai/internal/repo/vector"
	"github.com/StellrisJAY/workflow-ai/internal/service"
	"github.com/StellrisJAY/workflow-ai/internal/workflow"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-contrib/cors"
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
	repository.MigrateDB()
	nodeId, err := strconv.ParseInt(r.conf.Server.Id, 10, 64)
	if err != nil {
		return err
	}
	snowflakeNode, err := snowflake.NewNode(nodeId)
	if err != nil {
		return err
	}

	store := fs.NewFileStore(r.conf)
	llmRepo := repo.NewLLMRepo(repository)
	templateRepo := repo.NewTemplateRepo(repository)
	instanceRepo := repo.NewInstanceRepo(repository)
	kbRepo := repo.NewKnowledgeBaseRepo(repository)
	tm := repo.NewTransactionManager(repository)
	vectorstoreFactory := vector.MakeFactory(*r.conf)
	documentProcessor := rag.NewDocumentProcessor(8, kbRepo, store, llmRepo, vectorstoreFactory)
	engine := workflow.NewEngine(instanceRepo, llmRepo, snowflakeNode, tm, kbRepo, documentProcessor)

	llmService := service.NewLLMService(llmRepo, snowflakeNode)
	templateService := service.NewTemplateService(templateRepo, snowflakeNode)
	workflowService := service.NewWorkflowService(templateRepo, engine, instanceRepo)
	kbService := service.NewKnowledgeBaseService(kbRepo, snowflakeNode, tm, store, documentProcessor)

	llmHandler := NewLLMHandler(llmService)
	templateHandler := NewTemplateHandler(templateService)
	workflowHandler := NewWorkflowHandler(workflowService)
	kbHandler := NewKnowledgeBaseHandler(kbService)

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
			template.DELETE("/:id", templateHandler.Delete)
			template.PUT("/update", templateHandler.Update)
			template.GET("/detail/:id", templateHandler.GetDetail)
			template.GET("/list", templateHandler.List)
			template.GET("/start-variables/:id", templateHandler.GetStartInputVariables)
		}
		wf := v1.Group("/workflow")
		{
			wf.POST("/start", workflowHandler.Start)
			wf.GET("/detail/:id", workflowHandler.GetDetail)
			wf.GET("/outputs/:id", workflowHandler.Outputs)
			wf.GET("/list", workflowHandler.List)
			wf.GET("/node/detail", workflowHandler.GetNodeInstanceDetail)
		}
		kb := v1.Group("/knowledgeBase")
		{
			kb.POST("/create", kbHandler.Create)
			kb.PUT("/update", kbHandler.Update)
			kb.GET("/detail/:id", kbHandler.Detail)
			kb.GET("/files/:kbId", kbHandler.ListFiles)
			kb.GET("/list", kbHandler.List)
			kb.POST("/upload", kbHandler.Upload)
			kb.DELETE("/file/:id", kbHandler.Delete)
			kb.GET("/process/options/:id", kbHandler.GetFileProcessOptions)
			kb.GET("/download/:id", kbHandler.DownloadFile)
			kb.PUT("/process/options", kbHandler.UpdateFileProcessOptions)
			kb.POST("/process/start/:id", kbHandler.StartFileProcessing)
			kb.POST("/similarity-search", kbHandler.SimilaritySearch)
			kb.POST("/fulltext-search", kbHandler.FulltextSearch)
			kb.GET("/chunks", kbHandler.ListChunks)
		}
	}
	return nil
}

func (r *Router) Start() error {
	return r.e.Run(r.conf.Server.Host + ":" + r.conf.Server.Port)
}
