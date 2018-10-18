package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type EngineEnum string

//Topic Model
type Topic struct {
	//gorm.Model
	ID        uint `gorm:"primary_key" json:"-"`
	Name   string `gorm:"not null;unique" json:"name"`
	Engine string `json:"engine"`
	ResourceID string `json:"resource_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
	//DeletedAt *time.Time `sql:"index"`
}

type Subscriber struct {
	ID        uint `gorm:"primary_key" json:"-"`
	Name   string `gorm:"not null;unique" json:"name"`
	ResourceID string `json:"resource_id"`
	Endpoint string `json:"endpoint"`
	Topic string `json:"topic"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

type Publisher struct {
	ID        uint `gorm:"primary_key" json:"-"`
	Name   string `gorm:"not null;unique" json:"name"`
	ResourceID string `json:"resource_id"`
	Topic string `json:"topic"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`

}

type Engine struct {
	gorm.Model
	ID   int64
	Name string
}

type PublishMessage struct {
	Topic string
	Payload interface{}
	MessageID string
}