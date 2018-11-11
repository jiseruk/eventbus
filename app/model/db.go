package model

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wenance/wequeue-management_api/app/config"
)

type DB struct {
	//Mock with https://github.com/Selvatico/go-mocket/wiki/Documentation
	*gorm.DB
}

func NewDB() (*DB, error) {
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True",
			config.Get("database.mysql.user"),
			config.Get("database.mysql.password"),
			config.Get("database.mysql.host"),
			config.Get("database.mysql.database")))

	if err != nil {
		return nil, err
	}
	db.DropTableIfExists(&Topic{}, &Subscriber{})
	db.AutoMigrate(&Topic{})
	db.AutoMigrate(&Subscriber{})
	return &DB{DB: db}, nil
}
