package helper

import (
	// "api-mobile/general"
	"bytes"
	"clean_arch_v2/models"
	"fmt"
	"io"
	"regexp"
	"runtime"
	"strings"
	"time"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct{}
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// CORS will handle the CORS middleware
func (w *GoMiddleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func identifyPanic() string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(3, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	switch {
	case name != "":
		return fmt.Sprintf("%v:%v", name, line)
	case file != "":
		return fmt.Sprintf("%v:%v", file, line)
	}

	return fmt.Sprintf("pc:%x", pc)
}

// PanicCatcher is use for collecting panic that happened in endpoint
func (w *GoMiddleware) PanicCatcher(mw io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			_ = c.Request.ParseForm()
			reqBodyStr := ""
			for key, val := range c.Request.Form {
				reqBodyStr += key + ":"
				for i, v := range val {
					reqBodyStr += v
					if i+1 != len(val) {
						reqBodyStr += ","
					} else {
						reqBodyStr += " "
					}
				}
			}
			rec := recover()
			if rec != nil {
				user := "unknown user"
				device := "unknown device"
				userdata, isExist := c.Get("user")
				if isExist {
					user = userdata.(string)
				}
				devicedata, isExist2 := c.Get("device")
				if isExist2 {
					device = devicedata.(string)
				}
				fmt.Fprintf(mw, `level=error datetime="%s" ip=%s method=%s url="%s" user="%s" device="%s" panic=%v trace=%v`+"\n",
					time.Now().Format(time.RFC1123),
					c.ClientIP(),
					c.Request.Method,
					c.Request.URL.String(),
					user,
					device,
					rec,
					identifyPanic(),
				)
				// fmt.Println(string(debug.Stack()))
				// errPanic := fmt.Sprintf("Endpoint: %s - panic: %v", c.Request.RequestURI, rec)
				// errorcollector.WritePanic(errPanic, debug.Stack())
				var message []string
				message = append(message, viper.GetString("default_unhandled_error"))
				c.JSON(500, gin.H{
					"status":  false,
					"messages": message,
					"data":    new(struct{}),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CustomLogger provide custom log
func (w *GoMiddleware) CustomLogger(mw io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		latency := time.Since(t)
		_ = c.Request.ParseForm()
		reqBodyStr := ""
		for key, val := range c.Request.Form {
			reqBodyStr += key + ":"
			for i, v := range val {
				reqBodyStr += "'" + v + "'"
				if i+1 != len(val) {
					reqBodyStr += ","
				} else {
					reqBodyStr += " "
				}
			}
		}
		reqHeaderStr := ""
		for key, val := range c.Request.Header {
			reqHeaderStr += key + ":"
			for i, v := range val {
				reqHeaderStr += v
				if i+1 != len(val) {
					reqHeaderStr += ","
				} else {
					reqHeaderStr += " "
				}
			}
		}
		res := blw.body.String()
		contentType := blw.ResponseWriter.Header().Get("Content-Type")
		contentTypeSplit := strings.Split(contentType, ";")
		if len(contentTypeSplit) > 0 {
			if contentTypeSplit[0] == "text/html" {
				var re = regexp.MustCompile(`/\s+|\n+|\r/`) //make html string to one line
				res = re.ReplaceAllString(res, "")
				res += " \n"
			}
		}
		if res == "" {
			res = " \n"
		}
		user := "unknown user"
		device := "unknown device"
		userdata, isExist := c.Get("user")
		if isExist {
			user = userdata.(string)
		}
		devicedata, isExist2 := c.Get("device")
		if isExist2 {
			device = devicedata.(string)
		}
		fmt.Fprintf(mw, `
		level=info
		datetime="%s" 
		ip=%s method=%s 
		url="%s" 
		proto=%s 
		status=%d 
		latency=%s 
		user="%s" 
		device="%s" 
		req_header:"%s" 
		req_body="%s" 
		response=%s`,
		t.Format(time.RFC1123), c.ClientIP(), c.Request.Method, c.Request.URL.String(),
		c.Request.Proto, c.Writer.Status(), latency, user,
		device, reqHeaderStr, reqBodyStr, res,
		)
	}
}
var TimeoutContext time.Duration
var Conn	[]models.WebSocketConnection
// // InitMiddleware intialize the middleware
func InitMiddleware() *GoMiddleware {
	TimeoutContext = time.Duration(viper.GetInt("context.timeout")) * time.Second
	Conn = make([]models.WebSocketConnection, 0)
	return &GoMiddleware{}
}

