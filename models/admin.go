package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	UserName string `json:"username" gorm:"unique_index"`
	Email    string `json:"email" gorm:"unique_index"`
	Password string `json:"password"`
}
