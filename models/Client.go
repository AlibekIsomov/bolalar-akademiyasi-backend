package models

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name           string `json:"name"`
	PhoneNumber    string `json:"phone_number"`
	Age            int    `json:"age"`
	ClientsComment string `json:"clients_comment"`
	Status         string `json:"status"`
}
