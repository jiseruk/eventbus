package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type EngineEnum string

//Topic Model
type Topic struct {
	//gorm.Model
	ID         uint      `gorm:"primary_key" json:"-"`
	Name       string    `gorm:"not null;unique" json:"name" binding:"required"`
	Engine     string    `json:"engine" binding:"required"`
	ResourceID string    `json:"resource_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"-"`
	//DeletedAt *time.Time `sql:"index"`
}

type Subscriber struct {
	ID             uint      `gorm:"primary_key" json:"-"`
	Name           string    `gorm:"not null;unique" json:"name"`
	ResourceID     string    `json:"resource_id"`
	Endpoint       string    `gorm:"not null;unique" json:"endpoint" binding:"url"`
	Topic          string    `json:"topic"`
	PullResourceID string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"-"`
}

type Publisher struct {
	ID         uint      `gorm:"primary_key" json:"-"`
	Name       string    `gorm:"not null;unique" json:"name"`
	ResourceID string    `json:"resource_id"`
	Topic      string    `json:"topic"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"-"`
}

type Engine struct {
	gorm.Model
	ID   int64
	Name string
}

type PublishMessage struct {
	Topic     string      `json:"topic"`
	Payload   interface{} `json:"payload" validate:"required"`
	MessageID string      `json:"message_id"`
}

type Messages struct {
	Topic    string
	Messages []Message
}

type Message struct {
	Payload   interface{} `json:"payload"`
	MessageID string      `json:"message_id"`
}

type ConsumerRequest struct {
	MaxMessages int64  `form:"max_messages"`
	Subscriber  string `form:"subscriber"`
}
