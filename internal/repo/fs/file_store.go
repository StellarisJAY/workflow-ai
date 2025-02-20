package fs

import (
	"context"
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"io"
)

type FileStore interface {
	Upload(ctx context.Context, path string, data io.Reader) error
	Download(ctx context.Context, path string) ([]byte, error)
	Delete(ctx context.Context, path string) error
}

const (
	FileSystem          string = "file_system"
	FileStoreTencentCOS string = "tencent_cos"
)

func NewFileStore(config *config.Config) FileStore {
	switch config.FileStoreType {
	case FileStoreTencentCOS:
		return MakeTencentCOS(*config)
	default:
		panic("Unsupported file store type: " + config.FileStoreType)
	}
}
