package utils

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// add the middleware function
func GuidMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New()
		c.Set("uuid", uuid.String())
		fmt.Printf("The request with uuid %s is started \n", uuid)
		fmt.Printf("Request body %s  \n", c.Request)
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w
		c.Next()
		fmt.Println("Response body: " + w.body.String())
		fmt.Printf("The request with uuid %s is served \n", uuid)
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
