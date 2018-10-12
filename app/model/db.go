package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)


type DB struct {
	*gorm.DB
}

func NewDB() (*DB, error) {
	db, err := gorm.Open("mysql", "root:root@tcp(mysql:3306)/wequeue?charset=utf8&parseTime=True")
	if err != nil {
		return nil, err
	}
	db.CreateTable(&Topic{})
	return &DB{db}, nil
}
