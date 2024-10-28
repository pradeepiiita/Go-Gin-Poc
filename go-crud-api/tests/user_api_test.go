package tests

import (
	"bytes"
	"encoding/json"
	"example.com/go-crud-api/router"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "postgres://tester:test_password@localhost:5432/testdb?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func BeginTestTransaction(db *gorm.DB) *gorm.DB {
	return db.Begin()
}

func TestCreateUser(t *testing.T) {
	// Initialize test database
	db, err := SetupTestDB()
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	tx := BeginTestTransaction(db)
	defer tx.Rollback()

	// Set up router with test database
	testRouter := router.SetupRouter(tx)

	// Define the JSON payload for the POST request
	userPayload := map[string]string{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "securepassword123",
		"id":       uuid.New().String(),
	}
	requestBody, _ := json.Marshal(userPayload)

	// Create the HTTP POST request and recorder
	req, _ := http.NewRequest("POST", "/testusers", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Serve the request
	testRouter.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, w.Code)

	// Define the expected response structure
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, "John Doe", response["name"])
	assert.Equal(t, "john@example.com", response["email"])
	assert.NotEmpty(t, response["id"])
}
