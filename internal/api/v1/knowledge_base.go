package v1

import (
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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

func (k *KnowledgeBaseHandler) BatchUpload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		panic(err)
	}
	kbIdParam := form.Value["kbId"]
	kbId, err := strconv.ParseInt(kbIdParam[0], 10, 64)
	if err != nil {
		panic(err)
	}
	files := make([]*model.KnowledgeBaseFile, 0, len(form.File))
	data := make([]io.Reader, 0, len(files))
	for _, headers := range form.File {
		header := headers[0]
		file := &model.KnowledgeBaseFile{
			Name:   header.Filename,
			Length: header.Size,
			KbId:   kbId,
		}
		reader, err := header.Open()
		if err != nil {
			panic(err)
		}
		data = append(data, reader)
		files = append(files, file)
	}
	n, err := k.service.BatchUploadFile(c, files, data)
	if err != nil {
		c.JSON(200, common.NewErrorResponseWithData(err.Error(), n))
		return
	}
	c.JSON(200, common.NewSuccessResponse(n))
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

func (k *KnowledgeBaseHandler) Delete(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(err)
	}
	if err := k.service.Delete(c, id); err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(nil))
}

func (k *KnowledgeBaseHandler) DownloadFile(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(err)
	}
	data, name, err := k.service.DownloadFile(c, id)
	if err != nil {
		panic(err)
	}
	c.Header("Filename", name)
	c.Header("Content-Disposition", "attachment; filename=\""+name+"\"")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition,Filename")
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func (k *KnowledgeBaseHandler) StartFileProcessing(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(err)
	}
	if err := k.service.ProcessFile(c, id); err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(nil))
}

func (k *KnowledgeBaseHandler) SimilaritySearch(c *gin.Context) {
	var request model.KbSearchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(err)
	}
	result, err := k.service.SimilaritySearch(c, &request)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(result))
}

func (k *KnowledgeBaseHandler) FulltextSearch(c *gin.Context) {
	var request model.KbSearchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(err)
	}
	result, err := k.service.FulltextSearch(c, &request)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(result))
}

func (k *KnowledgeBaseHandler) HybridSearch(c *gin.Context) {
	var request model.KbSearchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(err)
	}
	result, err := k.service.HybridSearch(c, &request)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(result))
}

func (k *KnowledgeBaseHandler) ListChunks(c *gin.Context) {
	var request model.ListChunksRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		panic(err)
	}
	result, total, err := k.service.ListChunks(c, &request)
	if err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponseWithTotal(result, total))
}
