package services

import (
	"encoding/csv"
	"errors"
	"example.com/go-crud-api/db"
	"example.com/go-crud-api/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Client struct {
	conn *websocket.Conn
}

var clients = make(map[*Client]bool)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default (configure this as per your needs)
		return true
	},
}

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	client := &Client{conn: conn}
	clients[client] = true

	// Listen to incoming WebSocket messages (if needed)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			delete(clients, client)
			break
		}
		log.Println("Received message:", string(msg))
	}
}

func CreateUser(user *models.User) (*models.User, error) {
	user.ID = uuid.New().String()
	res := db.Db.Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

// Function to insert batch into PostgreSQL
func insertBatch(batch [][]string, wg *sync.WaitGroup, semaphore chan struct{}) error {
	// Start a transaction
	defer wg.Done()
	semaphore <- struct{}{}

	tx := db.Db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %v", tx.Error)
	}

	// Insert each record in the batch
	for _, record := range batch {
		if record[0] == "Name" {
			continue
		}
		age, err := strconv.Atoi(record[8]) // Convert Age
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to convert age: %v", err)
		}

		pincode, err := strconv.Atoi(record[12])
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to convert pincode: %v", err)
		}

		// Create User object
		user := models.User{
			ID:            uuid.New().String(),
			Name:          record[0],
			LastName:      record[1],
			Email:         record[2],
			Password:      record[3],
			City:          record[4],
			State:         record[5],
			Country:       record[6],
			Occupation:    record[7],
			Age:           age,
			Qualification: record[9],
			Username:      record[10],
			Gender:        record[11],
			Pincode:       pincode,
			LanguagePref:  record[13],
		}

		// Insert the user record
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert record: %v", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	<-semaphore
	return nil
}

func CreateBulkUser(reader *csv.Reader, requestId string) {
	startTime := time.Now() // Capture the start time
	batches, err := readCSVInBatches(reader, 10000)
	if err != nil {
		fmt.Printf("Error reading CSV in batches: %v\n", err)
		return
	}
	var wg sync.WaitGroup
	concurrencyLimit := 20
	semaphore := make(chan struct{}, concurrencyLimit)
	// Insert each batch into the database
	for i, batch := range batches {
		wg.Add(1)
		go insertBatch(batch, &wg, semaphore)
		fmt.Printf("Batch %d inserted successfully.\n", i+1)
	}
	wg.Wait()
	duration := time.Since(startTime)

	result := db.Db.Exec("UPDATE async_process_statuses SET status = ? WHERE id = ?", "Completed", requestId)
	if result.Error != nil {
		return
	}

	log.Printf("Execution time: %v\n", duration)
	fmt.Println("All data inserted successfully.")
	SendUpdateToClients(fmt.Sprintf("Request %s is complete", requestId))
	createResponseFile(requestId)
}

func UpdateAsyncStatus(requestId string, status string) {
	asynprstatus := models.AsyncProcessStatus{
		ID:     requestId,
		Status: status,
	}
	tx := db.Db.Begin()
	if err := tx.Create(&asynprstatus).Error; err != nil {
		tx.Rollback()
	}
	if err := tx.Commit().Error; err != nil {
	}
}

func SendUpdateToClients(message string) {
	for client := range clients {
		err := client.conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Write error:", err)
			client.conn.Close()
			delete(clients, client)
		}
	}
}

func createResponseFile(requestId string) {
	filePath := "/Users/pk/Desktop/AsyncResponse/response_" + requestId + "_.csv"
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Example data for the CSV
	header := []string{"ID", "Name", "Email", "CreatedAt"}
	data := [][]string{
		{"1", "John Doe", "john@example.com", time.Now().Format(time.RFC3339)},
		{"2", "Jane Smith", "jane@example.com", time.Now().Format(time.RFC3339)},
		{"3", "Mark Johnson", "mark@example.com", time.Now().Format(time.RFC3339)},
	}

	// Write the header to the CSV file
	if err := writer.Write(header); err != nil {
		log.Println("Error writing header:", err)
		return
	}

	// Write data rows
	for _, record := range data {
		if err := writer.Write(record); err != nil {
			log.Println("Error writing record:", err)
			return
		}
	}

	// Flush the data to the file and check for errors
	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Println("Error flushing CSV:", err)
		return
	}

}

// Function to read CSV in batches
func readCSVInBatches(r *csv.Reader, batchSize int) ([][][]string, error) {
	var allBatches [][][]string

	for {
		var batch [][]string
		for i := 0; i < batchSize; i++ {
			record, err := r.Read()
			if err == io.EOF {
				if len(batch) > 0 {
					allBatches = append(allBatches, batch) // Append the last non-full batch
				}
				return allBatches, nil // End of file
			}
			if err != nil {
				return nil, fmt.Errorf("error reading CSV: %v", err)
			}
			batch = append(batch, record)
		}
		allBatches = append(allBatches, batch) // Add batch to the collection of batches
	}
}

func GetUser(id string) (*models.User, error) {
	var user models.User
	res := db.Db.First(&user, "id = ?", id)
	if res.RowsAffected == 0 {
		return nil, errors.New(fmt.Sprintf("user of id %s not found", id))
	}
	return &user, nil
}

func GetAsyncProcessStatus(id string) (*models.AsyncProcessStatus, error) {
	var AsyncProcessStatus models.AsyncProcessStatus
	res := db.Db.First(&AsyncProcessStatus, "id = ?", id)
	if res.RowsAffected == 0 {
		return nil, errors.New(fmt.Sprintf("AsyncProcessStatus of id %s not found", id))
	}
	return &AsyncProcessStatus, nil
}

func GetUsers() ([]*models.User, error) {
	var users []*models.User
	res := db.Db.Find(&users)
	if res.Error != nil {
		return nil, errors.New("no users found")
	}
	return users, nil
}

func UpdateUser(user *models.User) (*models.User, error) {
	var userToUpdate models.User
	result := db.Db.Model(&userToUpdate).Where("id = ?", user.ID).Updates(user)
	if result.RowsAffected == 0 {
		return &userToUpdate, errors.New("user not updated")
	}
	return user, nil
}

func DeleteUser(id string) error {
	var deletedUser models.User
	result := db.Db.Where("id = ?", id).Delete(&deletedUser)
	if result.RowsAffected == 0 {
		return errors.New("user not deleted")
	}
	return nil
}
