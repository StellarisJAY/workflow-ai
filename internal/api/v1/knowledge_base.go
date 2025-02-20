package v1

import (
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type KnowledgeBaseHandler struct {
	service *service.KnowledgeBaseService
}

func NewKnowledgeBaseHandler(service *service.KnowledgeBaseService) *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{service: service}
}

func (k *KnowledgeBaseHandler) List(c *gin.Context) {
	var query model.KnowledgeBaseQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(err)
	}
	list, total, err := k.service.List(c, &query)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponseWithTotal(list, total))
}

func (k *KnowledgeBaseHandler) Detail(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(err)
	}
	detail, err := k.service.Detail(c, id)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(detail))
}

func (k *KnowledgeBaseHandler) Create(c *gin.Context) {
	var kb model.KnowledgeBase
	if err := c.ShouldBindJSON(&kb); err != nil {
		panic(err)
	}
	if err := k.service.Create(c, &kb); err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(nil))
}

func (k *KnowledgeBaseHandler) Update(c *gin.Context) {

}

func (k *KnowledgeBaseHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		panic(err)
	}
	filename := file.Filename
	reader, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	val := c.PostForm("kbId")
	kbId, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}
	if err := k.service.UploadFile(c, &model.KnowledgeBaseFile{KbId: kbId, Name: filename, Length: file.Size}, reader); err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(nil))
}

func (k *KnowledgeBaseHandler) ListFiles(c *gin.Context) {
	var query model.KbFileQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(err)
	}
	param := c.Param("kbId")
	kbId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(err)
	}
	list, total, err := k.service.ListFile(c, kbId, &query)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponseWithTotal(list, total))
}
