package services

import (
	"encoding/json"
	"example.com/go-crud-api/db"
	"example.com/go-crud-api/models"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

func CreateLimitUpdate(limitUpdateJson *models.LimitUpdateJson) (*models.LimitUpdateJson, error) {
	limitUpdateJson.ID = uuid.New().String()

	detailsJSON, err := json.Marshal(limitUpdateJson.Details)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	//Insert product data into the database
	if err := db.Db.Exec("INSERT INTO limit_update_jsons (id, details, status, started_at) VALUES (?, ?, ?, ?)", limitUpdateJson.ID, string(detailsJSON), "In Progress", time.Now()).Error; err != nil {
		log.Fatalf("Error inserting limit_update_jsons: %v", err)
	}
	go LimitUpdateSvc2(limitUpdateJson.ID)
	return limitUpdateJson, nil
}

func LimitUpdateSvc2(id string) {
	apiURL := "http://127.0.0.1:8081/limitUpdateApi/" + id
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Fatalf("Failed to make GET request")
		return
	}
	defer resp.Body.Close()
}
