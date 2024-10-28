package router

import (
	"example.com/go-crud-api/controllers"
	"example.com/go-crud-api/middleware"
	"example.com/go-crud-api/services"
	"example.com/go-crud-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(utils.GuidMiddleware())
	r.POST("/login", controllers.Login)
	protected := r.Group("/api")
	protected.Use(utils.GuidMiddleware())
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/users", controllers.GetUsers)
	protected.POST("/users", controllers.PostUser)
	protected.GET("/users/:id", controllers.GetUser)
	protected.POST("/limitUpdate", controllers.LimitUpdate)
	protected.POST("/bulkusers", controllers.PostBulkUser)
	protected.GET("/asyncstatus/:id", controllers.Asyncstatus)
	protected.PUT("/users/:id", controllers.UpdateUser)
	protected.DELETE("/users/:id", controllers.DeleteUser)
	r.GET("/ws", func(c *gin.Context) {
		services.HandleWebSocket(c)
	})
	return r
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.POST("/testusers", controllers.PostUserNew(db))
	return r
}
