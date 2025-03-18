package main

import (
	"flag"
	v1 "github.com/StellrisJAY/workflow-ai/internal/api/v1"
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"github.com/StellrisJAY/workflow-ai/internal/rag"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/StellrisJAY/workflow-ai/internal/repo/fs"
	"github.com/StellrisJAY/workflow-ai/internal/repo/vector"
	"github.com/StellrisJAY/workflow-ai/internal/service"
	"github.com/StellrisJAY/workflow-ai/internal/workflow"
	"github.com/bwmarrin/snowflake"
	"strconv"
)

var (
	configFile = flag.String("config", "config.yaml", "Path to config file")
)

func main() {
	flag.Parse()
	conf, err := config.ParseConfig(*configFile)
	if err != nil {
		panic(err)
	}
	router := v1.NewRouter(conf)

	repository, err := repo.NewRepository(conf)
	if err != nil {
		panic(err)
	}
	repository.MigrateDB()
	nodeId, err := strconv.ParseInt(conf.Server.Id, 10, 64)
	if err != nil {
		panic(err)
	}
	snowflakeNode, err := snowflake.NewNode(nodeId)
	if err != nil {
		panic(err)
	}

	store := fs.NewFileStore(conf)
	llmRepo := repo.NewProviderRepo(repository)
	templateRepo := repo.NewTemplateRepo(repository)
	instanceRepo := repo.NewInstanceRepo(repository)
	kbRepo := repo.NewKnowledgeBaseRepo(repository)
	fileRepo := repo.NewFileRepo(repository)
	providerRepo := repo.NewProviderRepo(repository)
	tm := repo.NewTransactionManager(repository)
	vectorstoreFactory := vector.MakeFactory(*conf)
	documentProcessor := rag.NewDocumentProcessor(8, kbRepo, store, llmRepo, vectorstoreFactory)
	engine := workflow.NewEngine(instanceRepo, llmRepo, snowflakeNode, tm, kbRepo, documentProcessor, conf, fileRepo, store)

	templateService := service.NewTemplateService(templateRepo, snowflakeNode)
	workflowService := service.NewWorkflowService(templateRepo, engine, instanceRepo)
	kbService := service.NewKnowledgeBaseService(kbRepo, snowflakeNode, tm, store, documentProcessor, vectorstoreFactory)
	fileService := service.NewFileService(fileRepo, store, tm, snowflakeNode)
	providerService := service.NewProviderService(providerRepo, snowflakeNode)

	templateHandler := v1.NewTemplateHandler(templateService)
	workflowHandler := v1.NewWorkflowHandler(workflowService)
	kbHandler := v1.NewKnowledgeBaseHandler(kbService)
	fileHandler := v1.NewFSHandler(fileService)
	providerHandler := v1.NewProviderHandler(providerService)

	if err := router.Init(templateHandler, workflowHandler, kbHandler, fileHandler, providerHandler); err != nil {
		panic(err)
	}
	if err = router.Start(); err != nil {
		panic(err)
	}
}
