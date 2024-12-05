package router

import (
	"context"
	"example.com/go-crud-api/controllers"
	_ "example.com/go-crud-api/docs"
	pb "example.com/go-crud-api/go-crud-api"
	"example.com/go-crud-api/middleware"
	"example.com/go-crud-api/models"
	"example.com/go-crud-api/services"
	"example.com/go-crud-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"net/http"
	"os"
	"time"
)

// @title My API
// @version 1.0
// @description This is a sample API.
// @host localhost:8080
// @BasePath /
func InitRouter() *gin.Engine {
	logger := logrus.New()
	file, err := os.OpenFile("/Users/pk/Desktop/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}
	logger.SetFormatter(&logrus.JSONFormatter{})

	router := gin.Default()
	router.Use(gin.LoggerWithWriter(logger.Writer()))
	router.Use(gin.Recovery())

	router.Use(utils.GuidMiddleware())
	// Start the Datadog tracer
	tracer.Start(
		tracer.WithAgentAddr("localhost:8126"), // Address of the Datadog Agent
		tracer.WithServiceName("go-crud-api"),
	)
	defer tracer.Stop()
	router.Use(gintrace.Middleware("go-crud-api"))
	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Simple route
	router.POST("/login", controllers.Login)

	// Example endpoint to call gRPC from Gin
	router.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to gRPC server"})
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		res, err := client.GetUser(ctx, &pb.UserRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"name": res.Name, "age": res.Age})
	})

	protected := router.Group("/api/v1")
	protected.Use(gin.Recovery())
	protected.Use(gintrace.Middleware(os.Getenv("DD_SERVICE")))
	protected.Use(utils.GuidMiddleware())
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/users", controllers.GetUsers)
	protected.POST("/users", middleware.ValidationMiddleware(&models.User{}), controllers.PostUser)
	protected.GET("/users/:id", controllers.GetUser)
	protected.POST("/limitUpdate", controllers.LimitUpdate)
	protected.POST("/bulkusers", controllers.PostBulkUser)
	protected.GET("/asyncstatus/:id", controllers.Asyncstatus)
	protected.PUT("/users/:id", controllers.UpdateUser)
	protected.DELETE("/users/:id", controllers.DeleteUser)
	router.GET("/ws", func(c *gin.Context) {
		services.HandleWebSocket(c)
	})
	router.GET("/ws-connect", func(c *gin.Context) {
		services.CallWebSocket(c)
	})
	return router
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.POST("/testusers", controllers.PostUserNew(db))
	return r
}
