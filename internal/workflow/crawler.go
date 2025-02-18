package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (e *Engine) executeCrawlerNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance,
	crawlerData *model.CrawlerNodeData, inputMap map[string]any) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			nodeInstance.Status = model.NodeInstanceStatusFailed
			nodeInstance.CompleteTime = time.Now()
			nodeInstance.Error = err.Error()
			if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
				log.Println("update node instance error:", err)
			}
		}
	}()
	urlStr, ok := inputMap["url"]
	if !ok {
		panic(errors.New("url参数不存在"))
	}
	u, err := url.Parse(urlStr.(string))
	if err != nil {
		panic(err)
	}
	// http请求
	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		panic(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	result := make(map[string]any)
	// 请求失败
	if response.StatusCode != 200 {
		result["code"] = response.StatusCode
		result["message"] = response.Status
		return
	}
	result["code"] = response.StatusCode
	result["message"] = response.Status

	contentType := response.Header.Get("Content-Type")
	contentType = strings.Split(contentType, ";")[0]
	result["contentType"] = contentType
	switch contentType {
	case "application/json", "text/plain":
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		result["data"] = string(bytes)
	case "text/html":
		content, err := parseHTMLContent(response)
		if err != nil {
			panic(err)
		}
		result["data"] = content
	default:
		panic(fmt.Errorf("不支持的网页内容:%s", contentType))
	}
	outputs, _ := json.Marshal(result)
	nodeInstance.Output = string(outputs)
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	nodeInstance.CompleteTime = time.Now()
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		log.Println("update node instance error:", err)
	}
}

func parseHTMLContent(response *http.Response) (string, error) {
	node, err := html.Parse(response.Body)
	if err != nil {
		return "", err
	}
	var pContents []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			var text strings.Builder
			var extractText func(*html.Node)
			extractText = func(child *html.Node) {
				if child.Type == html.TextNode {
					text.WriteString(strings.TrimSpace(child.Data))
				}
				for c := child.FirstChild; c != nil; c = c.NextSibling {
					extractText(c)
				}
			}
			extractText(n)
			pContents = append(pContents, text.String())
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
	content := strings.Join(pContents, "\n")
	return content, nil
}
