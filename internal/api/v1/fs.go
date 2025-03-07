package v1

import (
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type FSHandler struct {
	service *service.FileService
}

func NewFSHandler(service *service.FileService) *FSHandler {
	return &FSHandler{service: service}
}

func (f *FSHandler) Upload(c *gin.Context) {
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
	fileModel := model.File{
		Name: filename,
		Size: file.Size,
	}
	if err := f.service.Upload(c, &fileModel, reader); err != nil {
		panic(err)
	}
	c.JSON(200, common.NewSuccessResponse(struct {
		FileID int64 `json:"fileId,string"`
	}{fileModel.Id}))
}

func (f *FSHandler) BatchUpload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		panic(err)
	}
	files := make([]*model.File, 0, len(form.File))
	data := make([]io.Reader, 0, len(files))
	for _, headers := range form.File {
		header := headers[0]
		file := &model.File{
			Name: header.Filename,
			Size: header.Size,
			Type: strings.Replace(path.Ext(header.Filename), ".", "", 1),
		}
		reader, err := header.Open()
		if err != nil {
			panic(err)
		}
		data = append(data, reader)
		files = append(files, file)
	}
	if err := f.service.BatchUpload(c, files, data); err != nil {
		panic(err)
	}
	fileIDs := make([]string, len(files))
	for i, file := range files {
		fileIDs[i] = strconv.FormatInt(file.Id, 10)
	}
	c.JSON(200, common.NewSuccessResponse(struct {
		FileIDs []string `json:"fileIds"`
	}{fileIDs}))
}

func (f *FSHandler) Download(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(err)
	}
	data, name, err := f.service.Download(c, id)
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
