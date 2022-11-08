package database

import (
	"sejuta-cita/app/models"

	"gorm.io/gorm"
)

func InitMigration(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Role{})

	SeederRole(db)
	SeederUser(db)
}

func SeederRole(db *gorm.DB) {
	var dataCount int64

	db.Model(&models.Role{}).Count(&dataCount)

	if dataCount == 0 {
		data := []models.Role{
			{
				Name: "superadmin",
			},
			{
				Name: "admin",
			},
			{
				Name: "oprator",
			},
		}
		db.Create(&data)
	}
}

func SeederUser(db *gorm.DB) {
	var dataCount int64

	db.Model(&models.User{}).Count(&dataCount)

	if dataCount == 0 {
		data := []models.User{
			{
				Name:     "muhammad iqbal",
				Email:    "emixbal@gmail.com",
				Password: "$2a$10$xO0eiq3.64vo1gR1cKkEE.hwn0OvafrzVI0HhsZWeb9UuXsl7bZrq", //aaaaaaaa
				RoleID:   1,
			},
			{
				Name:     "suti",
				Email:    "suti@gmail.com",
				Password: "$2a$10$xO0eiq3.64vo1gR1cKkEE.hwn0OvafrzVI0HhsZWeb9UuXsl7bZrq", //aaaaaaaa
				RoleID:   2,
			},
			{
				Name:     "yoso",
				Email:    "yoso@gmail.com",
				Password: "$2a$10$xO0eiq3.64vo1gR1cKkEE.hwn0OvafrzVI0HhsZWeb9UuXsl7bZrq", //aaaaaaaa
				RoleID:   3,
			},
		}
		db.Create(&data)
	}
}
