package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	SUBSCRIBER_PUSH     = "push"
	SUBSCRIBER_PUSH_DLT = "push-dlq"
	SUBSCRIBER_PULL     = "pull"
)

//Topic Model
type Topic struct {
	//gorm.Model
	ID         uint      `gorm:"primary_key" json:"-"`
	Name       string    `gorm:"not null;unique" json:"name" binding:"required" example:"topic_name"`
	Engine     string    `json:"engine" binding:"required,oneof=AWSStream AWS" example:"AWS"`
	ResourceID string    `json:"resource_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"-"`
	//DeletedAt *time.Time `sql:"index"`
}

type Subscriber struct {
	ID              uint      `gorm:"primary_key" json:"-"`
	Name            string    `gorm:"not null;unique" json:"name" binding:"required" example:"subscriber_name"`
	ResourceID      string    `json:"-"`
	Endpoint        string    `gorm:"not null;unique" json:"endpoint" binding:"required,url" example:"http://subscriber.wequeue.com/subscriber"`
	Topic           string    `json:"topic" binding:"required" example:"topic_name"`
	DeadLetterQueue string    `json:"-"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
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
	Topic    string    `json:"topic"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Payload     interface{}  `json:"payload"`
	MessageID   string       `json:"message_id"`
	DeleteToken *string      `json:"delete_token"`
	DeleteError *deleteError `json:"delete_error,omitempty"`
}

type deleteError struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

type ConsumerRequest struct {
	MaxMessages int64  `form:"max_messages" binding:"required,max=10,min=1"`
	Subscriber  string `form:"subscriber" binding:"required"`
}

type DeleteDeadLetterQueueMessagesRequest struct {
	Messages   []Message `json:"messages" binding:"required"`
	Subscriber string    `json:"subscriber" binding:"required"`
}

type DeleteDeadLetterQueueMessagesResponse struct {
	Failed []Message `json:"failed"`
	Topic  string    `json:"topic"`
}
