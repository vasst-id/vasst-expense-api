package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	logs "github.com/vasst-id/vasst-expense-api/internal/utils/logger"
)

// log is a unexported package-level global variable that holds julo-go-library logger instance
var log *logs.Logger

// Log returns logs.Logger instance
func Log() *logs.Logger {
	return log
}

// SetLogger set log variable. Must be set when init, otherwise every call to log will panic
func SetLogger(logger *logs.Logger) {
	if log == nil {
		log = logger
	}
}

// LogRequestBody is an alternative to RequestBodyLoggerHandler without logger args
func LogRequestBody(c *gin.Context, funcName string) {
	body, err := io.ReadAll(c.Request.Body)

	// always close the body when surrounding function returns
	defer func() {
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
	}()

	if err != nil {
		log.Err(err).
			Msg("RequestBodyLogger.io.ReadAll")
		return
	}

	// minify json
	dst := &bytes.Buffer{}
	err = json.Compact(dst, body)
	if err != nil {
		log.Err(err).
			Msg("RequestBodyLogger.json.Compact")
		return
	}

	// log request body
	log.Info().
		Str("funcName", funcName).
		Str("body", dst.String()).
		Send()
}

func RequestBodyLoggerHandler(c *gin.Context, funcName string) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	log.Printf("Request Body of %s: \n %s", funcName, body)
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
}
