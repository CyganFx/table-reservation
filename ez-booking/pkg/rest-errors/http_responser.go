package rest_errors

import (
	"fmt"
	"github.com/CyganFx/table-reservation/ez-booking/internal/delivery/http-v1"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"runtime/debug"
)

type httpResponser struct {
	errorLog *log.Logger
}

func NewHttpResponser(errorLog *log.Logger) http_v1.Responser {
	return &httpResponser{errorLog: errorLog}
}

func (h *httpResponser) ServerError(c *gin.Context, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.errorLog.Output(2, trace)

	http.Error(c.Writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *httpResponser) ClientError(c *gin.Context, status int) {
	http.Error(c.Writer, http.StatusText(status), status)
}

func (h *httpResponser) NotFound(c *gin.Context) {
	h.ClientError(c, http.StatusNotFound)
}
