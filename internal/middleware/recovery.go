package middleware

import (
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/gin-gonic/gin"
	"net"
	"os"
	"strings"
)

func Recovery(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			err, _ := e.(error)
			var brokenPipe bool
			var ne *net.OpError
			if errors.As(err, &ne) {
				var se *os.SyscallError
				if errors.As(ne, &se) {
					seStr := strings.ToLower(se.Error())
					if strings.Contains(seStr, "broken pipe") ||
						strings.Contains(seStr, "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			if brokenPipe {
				c.JSON(200, common.NewErrorResponse(err.Error()))
				c.Abort()
			} else {
				c.JSON(200, common.NewErrorResponse(err.Error()))
			}
		}
	}()
	c.Next()
}
