package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Adminname string `json:"adminname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}
