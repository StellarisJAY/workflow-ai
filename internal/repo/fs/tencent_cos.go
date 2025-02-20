package fs

import (
	"context"
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"io"
	"net/http"
	"net/url"
)
import "github.com/tencentyun/cos-go-sdk-v5"

type TencentCOS struct {
	cli *cos.Client
}

func MakeTencentCOS(config config.Config) *TencentCOS {
	bu, _ := url.Parse(config.Cos.BucketUrl)
	su, _ := url.Parse(config.Cos.ServiceUrl)
	baseURL := cos.BaseURL{
		BucketURL:  bu,
		ServiceURL: su,
	}
	client := cos.NewClient(&baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.Cos.SecretId,
			SecretKey: config.Cos.SecretKey,
		},
	})
	return &TencentCOS{client}
}

func (t *TencentCOS) Upload(ctx context.Context, path string, reader io.Reader) error {
	_, err := t.cli.Object.Put(ctx, path, reader, nil)
	if err != nil {
		return err
	}
	return nil
}

func (t *TencentCOS) Download(ctx context.Context, path string) ([]byte, error) {
	response, err := t.cli.Object.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t *TencentCOS) Delete(ctx context.Context, path string) error {
	_, err := t.cli.Object.Delete(ctx, path)
	if err != nil {
		return err
	}
	return nil
}
