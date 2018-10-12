package model

import "github.com/jinzhu/gorm"

type EngineEnum string

//Topic Model
type Topic struct {
	gorm.Model
	//ID     int64
	Name   string `gorm:"not null;unique" json:"name"`
	Engine string `json:"engine"`
	ResourceID string
}

type Engine struct {
	gorm.Model
	ID   int64
	Name string
}
