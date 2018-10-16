package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jonboulle/clockwork"
)


type DB struct {
	//Mock with https://github.com/Selvatico/go-mocket/wiki/Documentation
	*gorm.DB
}

var Clock clockwork.Clock

func NewDB() (*DB, error) {
	Clock = clockwork.NewRealClock()
	db, err := gorm.Open("mysql", "root:root@tcp(mysql:3306)/wequeue?charset=utf8&parseTime=True")
	if err != nil {
		return nil, err
	}
	db.CreateTable(&Topic{})
	return &DB{DB: db}, nil
}
