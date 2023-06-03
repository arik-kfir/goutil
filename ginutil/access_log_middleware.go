package ginutil

import (
	"bytes"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
	"time"
)

type customReadCloser struct {
	r io.Reader
	c io.Closer
}

func (rc *customReadCloser) Read(p []byte) (int, error) {
	return rc.r.Read(p)
}

func (rc *customReadCloser) Close() error {
	return rc.c.Close()
}

func AccessLogMiddleware(c *gin.Context) {
	if c.Request.RequestURI == "/healthz" {
		c.Next()
		return
	}

	//
	// Set up a request logger, which:
	// - adds simple metadata fields
	// - replaces request body reader with a Tee reader which also logs the request body to a side copy
	// - adds all transfer-encoding headers
	// - adds all headers
	// - adds all trailers
	//
	requestBody := bytes.Buffer{}
	c.Request.Body = &customReadCloser{r: io.TeeReader(c.Request.Body, &requestBody), c: c.Request.Body}

	start := time.Now()
	c.Next()
	duration := time.Since(start)

	status := c.Writer.Status()
	var event *zerolog.Event
	if len(c.Errors) == 0 {
		if status >= 200 && status <= 399 {
			event = log.Ctx(c).Info()
		} else if status >= 400 && status <= 499 {
			event = log.Ctx(c).Warn()
		} else {
			event = log.Ctx(c).Error().Stack()
		}
	} else {
		event = log.Ctx(c).Error().Stack()
	}

	event = event.
		Str("request:id", requestid.Get(c)).
		Dur("http:process:duration", duration).
		Int("http:res:status", c.Writer.Status()).
		Int("http:res:size", c.Writer.Size())
	for name, values := range c.Writer.Header() {
		if strings.HasPrefix(name, "sec-") {
			continue
		}
		arr := zerolog.Arr()
		for _, value := range values {
			arr.Str(value)
		}
		event = event.Array("http:res:header:"+strings.ToLower(name), arr)
	}

	var errorsArr []error
	for _, err := range c.Errors {
		errorsArr = append(errorsArr, err.Err)
	}
	if len(errorsArr) > 0 {
		event = event.Err(errorsArr[0])
		if len(errorsArr) > 1 {
			event = event.Errs("http:res:errors", errorsArr)
		}
	}

	event.Msg("HTTP Request processed")
}
