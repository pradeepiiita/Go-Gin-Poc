package controllers

import (
	"bufio"
	"encoding/csv"
	"example.com/go-crud-api/middleware"
	"example.com/go-crud-api/models"
	"example.com/go-crud-api/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func GetUsers(ctx *gin.Context) {
	res, err := services.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"users": res,
	})
}

func Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	// Verify username and password (pseudo-code)
	if username == "user" && password == "password" { // Replace with real authentication
		token, _ := middleware.GenerateToken(username)
		ctx.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

func LimitUpdate(ctx *gin.Context) {
	var limitUpdateJson models.LimitUpdateJson
	err := ctx.Bind(&limitUpdateJson)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	res, err := services.CreateLimitUpdate(&limitUpdateJson)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"limitUpdateJson": res,
	})
}

func PostUserNew(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser models.User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Create the user in the database
		if err := db.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Respond with the created user details, excluding the password
		c.JSON(http.StatusCreated, gin.H{"id": newUser.ID, "name": newUser.Name, "email": newUser.Email})
	}
}

func PostUser(ctx *gin.Context) {
	var user models.User
	err := ctx.Bind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	res, err := services.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"user": res,
	})
}

func PostBulkUser(c *gin.Context) {
	// Parse the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from form"})
		return
	}

	// Open the uploaded file
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open the uploaded file"})
		return
	}
	//defer f.Close()

	// Read and process the CSV file
	reader := csv.NewReader(bufio.NewReader(f))
	value, _ := c.Get("uuid")
	requestId, ok := value.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
	services.UpdateAsyncStatus(requestId, "In Progress")
	go services.CreateBulkUser(reader, requestId)

	message := fmt.Sprintf("Request has been taken, you can track status with request id %s", requestId)
	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func Asyncstatus(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := services.GetAsyncProcessStatus(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	filePath := "/Users/pk/Desktop/AsyncResponse/response_" + id + "_.csv"
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		// File does not exist
		log.Println("File does not exist:", filePath)
		ctx.JSON(http.StatusOK, gin.H{
			"GetAsyncProcessStatus": res,
		})
		return
	} else if err != nil {
		// Other errors, such as permission issues
		log.Println("Error checking file:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error checking file",
			"error":   err.Error(),
		})
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to open file",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	ctx.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Transfer-Encoding", "binary")

	if _, err := io.Copy(ctx.Writer, file); err != nil {
		log.Println("Error sending file:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error sending file",
			"error":   err.Error(),
		})
	}

	//ctx.JSON(http.StatusOK, gin.H{
	//	"GetAsyncProcessStatus": res,
	//})
}

func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := services.GetUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"user": res,
	})
}

func UpdateUser(ctx *gin.Context) {
	var updatedUser models.User
	err := ctx.Bind(&updatedUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	id := ctx.Param("id")
	dbUser, err := services.GetUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	dbUser.Name = updatedUser.Name
	dbUser.Email = updatedUser.Email
	dbUser.Password = updatedUser.Password

	res, err := services.UpdateUser(dbUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"task": res,
	})
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := services.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "task deleted successfully",
	})
}
