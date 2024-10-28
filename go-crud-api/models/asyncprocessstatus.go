package models

type AsyncProcessStatus struct {
	ID     string `json:"id" gorm:"primarykey"`
	Status string `json:"status"`
}
