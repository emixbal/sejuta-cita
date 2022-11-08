package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name string `json:"Name" gorm:"index:idx_name,unique"`
}
