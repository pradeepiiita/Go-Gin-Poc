package models

type User struct {
	ID            string `json:"id" gorm:"primarykey"`
	Name          string `json:"name" validate:"required,min=2,max=100"`
	LastName      string `json:"lastname"`
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=8"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	Occupation    string `json:"occupation"`
	Age           int    `json:"age"`
	Qualification string `json:"qualification"`
	Username      string `json:"username"`
	Gender        string `json:"gender"`
	Pincode       int    `json:"pincode"`
	LanguagePref  string `json:"language_pref"`
}
