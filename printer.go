package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/labstack/echo"
)

type httpWriter struct {
	b *bytes.Buffer
	http.ResponseWriter
}

func (w *httpWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	if err == nil {
		w.b.Write(b)
	}

	return n, err
}

// Fully log the request and response
func FullRequestLog() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			reqStr := ""

			bys, err := httputil.DumpRequest(c.Request(), true)
			if err != nil {
				msg := fmt.Sprintf("Couldnt dump incoming request (%s)", err)
				panic(msg)
			}
			reqStr = string(bys)

			// Have to mock out our own writer to be able to dump
			origWriter := c.Response().Writer

			buff := &bytes.Buffer{}

			c.Response().Writer = &httpWriter{
				b:              buff,
				ResponseWriter: origWriter,
			}

			if err := next(c); err != nil {
				c.Error(err)
			}

			logBuf := &bytes.Buffer{}

			logBuf.WriteString("****NEW REQUEST************************************************************************\n\n")
			logBuf.WriteString(reqStr)
			logBuf.WriteString("\n****RESPONSE***************************************************************************\n")

			resp := c.Response()
			stdHeader := resp.Header()

			for headerName := range stdHeader {
				for _, val := range stdHeader[headerName] {
					logBuf.WriteString(headerName)
					logBuf.WriteString(": ")
					logBuf.WriteString(val)
					logBuf.WriteString("\n")
				}
			}
			logBuf.WriteString("\n")

			jsonMap := map[string]interface{}{}
			err = json.Unmarshal(buff.Bytes(), &jsonMap)

			if err == nil { // Pretty print json
				prettyJsonBytes, _ := json.MarshalIndent(jsonMap, "", "  ")
				logBuf.Write(prettyJsonBytes)
			} else {
				logBuf.WriteString("<output omitted, not json>")
			}
			logBuf.WriteString("\n***************************************************************************************")
			fmt.Println(logBuf.String())

			return nil
		}
	}
}
