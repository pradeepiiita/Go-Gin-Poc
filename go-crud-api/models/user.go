package models

type User struct {
	ID            string `json:"id" gorm:"primarykey"`
	Name          string `json:"name"`
	LastName      string `json:"lastname"`
	Email         string `json:"email"`
	Password      string `json:"password"`
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
