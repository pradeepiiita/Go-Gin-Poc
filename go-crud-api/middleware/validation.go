package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationMiddleware applies validation to all API inputs
func ValidationMiddleware(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON input to the struct
		if err := c.ShouldBindJSON(obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			c.Abort()
			return
		}

		// Validate the struct
		if err := validate.Struct(obj); err != nil {
			// Format validation errors
			var validationErrors []string
			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, err.Field()+" is invalid: "+err.ActualTag())
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			c.Abort()
			return
		}

		// Proceed to the next middleware/handler if valid
		c.Next()
	}
}
