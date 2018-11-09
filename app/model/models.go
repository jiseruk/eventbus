package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	SUBSCRIBER_PUSH = "push"
	SUBSCRIBER_PULL = "pull"
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
	Endpoint        *string   `gorm:"unique" json:"endpoint,omitempty" binding:"omitempty,url" example:"http://subscriber.wequeue.com/subscriber"`
	Topic           string    `json:"topic" binding:"required" example:"topic_name"`
	Type            string    `json:"type" binding:"required,oneof=pull push"`
	DeadLetterQueue string    `json:"dead_letter_queue,omitempty"`
	PullingQueue    string    `json:"pulling_queue,omitempty"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

func (s Subscriber) GetQueueURL() string {
	if s.Type == SUBSCRIBER_PUSH {
		return s.DeadLetterQueue
	}
	return s.PullingQueue
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
	Topic          string      `json:"topic"`
	Payload        interface{} `json:"payload" validate:"required"`
	Timestamp      *int64      `json:"timestamp"`
	SequenceNumber *string     `json:"sequence_number,omitempty"`
}

type Messages struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Message     PublishMessage `json:"message"`
	MessageID   string         `json:"message_id"`
	DeleteToken *string        `json:"delete_token"`
	DeleteError *deleteError   `json:"delete_error,omitempty"`
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
}
