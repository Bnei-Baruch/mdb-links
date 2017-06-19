package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/stvp/rollbar"
	"gopkg.in/gin-gonic/gin.v1"
)

// Set required resources for handlers in context
func EnvironmentMiddleware(mbdDB *sql.DB, backendUrls []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("MDB_DB", mbdDB)
		c.Set("BACKEND_URLS", backendUrls)
		c.Next()
	}
}

// Recover with error
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rval := recover(); rval != nil {
				debug.PrintStack()
				err, ok := rval.(error)
				if !ok {
					err = errors.Errorf("panic: %s", rval)
				}
				c.AbortWithError(http.StatusInternalServerError, err).SetType(gin.ErrorTypePrivate)
			}
		}()

		c.Next()
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func errorsToRollbarStack(st stackTracer) rollbar.Stack {
	t := st.StackTrace()
	rs := make(rollbar.Stack, len(t))
	for i, f := range t {
		// Program counter as it's computed internally in errors.Frame
		pc := uintptr(f) - 1
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			rs[i] = rollbar.Frame{
				Filename: "unknown",
				Method:   "?",
				Line:     0,
			}
			continue
		}

		// symtab info
		file, line := fn.FileLine(pc)
		name := fn.Name()

		// trim compile time GOPATH from file name
		fileWImportPath := trimGOPATH(name, file)

		// Strip only method name from FQN
		idx := strings.LastIndex(name, "/")
		name = name[idx+1:]
		idx = strings.Index(name, ".")
		name = name[idx+1:]

		rs[i] = rollbar.Frame{
			Filename: fileWImportPath,
			Method:   name,
			Line:     line,
		}
	}

	return rs
}

// Taken AS IS from errors pkg since it's not exported there.
// Check out the source code with good comments on https://github.com/pkg/errors/blob/master/stack.go
func trimGOPATH(name, file string) string {
	const sep = "/"
	goal := strings.Count(name, sep) + 2
	i := len(file)
	for n := 0; n < goal; n++ {
		i = strings.LastIndex(file[:i], sep)
		if i == -1 {
			i = -len(sep)
			break
		}
	}
	file = file[i+len(sep):]
	return file
}

// Handle all errors
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				switch e.Type {
				case gin.ErrorTypePublic:
					if e.Err != nil {
						log.Warnf("Public error: %s", e.Error())
						c.JSON(c.Writer.Status(), gin.H{"status": "error", "error": e.Error()})
					}

				case gin.ErrorTypeBind:
					// Keep the preset response status
					status := http.StatusBadRequest
					if c.Writer.Status() != http.StatusOK {
						status = c.Writer.Status()
					}
					log.WithFields(log.Fields{
						"error": e.Error(),
					}).Warn("Bind error")
					c.JSON(status, gin.H{
						"status": "error",
						"error":  e.Error(),
					})

				default:
					// Log all other errors
					log.Error(e.Err)
					st, ok := e.Err.(stackTracer)
					if ok {
						fmt.Printf("%s: %+v\n", st, st.StackTrace())
					}

					// Log to rollbar if we have a token setup
					if len(rollbar.Token) != 0 {
						if ok {
							rollbar.RequestErrorWithStack(rollbar.ERR, c.Request, e.Err,
								errorsToRollbarStack(st))
						} else {
							rollbar.RequestError(rollbar.ERR, c.Request, e.Err)
						}
					}
				}
			}

			// If there was no public or bind error, display default 500 message
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError,
					gin.H{"status": "error", "error": "Internal Server Error"})
			}
		}
	}
}
