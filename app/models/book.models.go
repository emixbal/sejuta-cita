package models

import (
	"fmt"
	"net/http"
	"sejuta-cita/config"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Author string `json:"Author"`
	Name   string `json:"Name"`
	NoISBN string `json:"NoISBN"`
}

func FethAllBooks() (Response, error) {
	var books []Book
	var res Response

	db := config.GetDBInstance()

	if result := db.Find(&books); result.Error != nil {
		fmt.Print("error FethAllBooks")
		fmt.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error fetchin records"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = books

	return res, nil
}

func CreateABook(book *Book) (Response, error) {
	var res Response
	db := config.GetDBInstance()

	if result := db.Create(&book); result.Error != nil {
		fmt.Print("error CreateABook")
		fmt.Print(result.Error)

		res.Status = http.StatusInternalServerError
		res.Message = "error save new record"
		return res, result.Error
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = book

	return res, nil
}
