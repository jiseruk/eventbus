package model

import (
	err "errors"
	"time"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/gorm"
	"github.com/jonboulle/clockwork"
	uuid "github.com/satori/go.uuid"
	"github.com/wenance/wequeue-management_api/app/errors"
)

const (
	SUBSCRIBER_PUSH = "push"
	SUBSCRIBER_PULL = "pull"
)

var Clock clockwork.Clock

//Topic Model
type Topic struct {
	//gorm.Model

	ID            uint      `gorm:"primary_key" json:"-"`
	Name          string    `gorm:"not null;unique" json:"name" example:"topic_name"`
	Engine        string    `json:"engine" example:"AWS"`
	ResourceID    string    `json:"resource_id"`
	SecurityToken string    `json:"security_token,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"-"`
	//DeletedAt *time.Time `sql:"index"`
}

func (t Topic) Validate() error {
	return validation.ValidateStruct(
		&t,
		validation.Field(&t.Name, validation.Required.Error(errors.ErrorFieldRequired),
			validation.Length(1, 50)),
		validation.Field(&t.Engine, validation.Required.Error(errors.ErrorFieldRequired),
			validation.In("AWS", "AWSStream").Error(errors.GetInListError("AWS", "AWSStream"))),
	)
}

type Subscriber struct {
	ID                uint      `gorm:"primary_key" json:"-"`
	Name              string    `gorm:"not null;unique" json:"name" example:"subscriber_name"`
	ResourceID        string    `json:"-"`
	Endpoint          *string   `gorm:"unique" json:"endpoint,omitempty" example:"http://subscriber.wequeue.com/subscriber"`
	Topic             string    `json:"topic" example:"topic_name"`
	Type              string    `json:"type"`
	DeadLetterQueue   string    `json:"dead_letter_queue,omitempty"`
	PullingQueue      string    `json:"pulling_queue,omitempty"`
	VisibilityTimeout *int      `json:"visibility_timeout,omitempty"`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
}

func (s Subscriber) Validate() error {
	rules := []*validation.FieldRules{
		validation.Field(&s.Name, validation.Required.Error(errors.ErrorFieldRequired)),
		validation.Field(&s.Topic, validation.Required.Error(errors.ErrorFieldRequired)),
	}
	typeRule := validation.Field(&s.Type,
		validation.Required.Error(errors.ErrorFieldRequired),
		validation.In(SUBSCRIBER_PULL, SUBSCRIBER_PUSH).Error(errors.GetInListError(SUBSCRIBER_PULL, SUBSCRIBER_PUSH)),
	)
	if err := validation.ValidateStruct(&s, typeRule); err == nil {
		if s.Type == SUBSCRIBER_PULL {
			rules = append(rules, validation.Field(&s.VisibilityTimeout,
				validation.Required.Error(errors.ErrorFieldRequired),
				validation.Min(0),
				validation.Max(43200),
			))
		} else {
			rules = append(rules, validation.Field(&s.Endpoint,
				validation.Required.Error(errors.ErrorFieldRequired),
				is.URL,
			))
		}
	}
	rules = append(rules, typeRule)

	return validation.ValidateStruct(&s, rules...)
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
	Payload        interface{} `json:"payload"`
	Timestamp      *int64      `json:"timestamp"`
	SequenceNumber *string     `json:"sequence_number,omitempty"`
}

func (p PublishMessage) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.Topic,
			validation.Required.Error(errors.ErrorFieldRequired)),
		validation.Field(&p.Payload,
			validation.Required.Error(errors.ErrorFieldRequired),
			validation.By(func(obj interface{}) error {

				if _, ok := obj.(map[string]interface{}); !ok {
					if _, ok := obj.([]interface{}); !ok {
						return err.New("it should be a valid json object")
					}
				}
				return nil
			})),
	)
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
	MaxMessages int64  `form:"max_messages"`
	Subscriber  string `form:"subscriber"`
}

func (c ConsumerRequest) Validate() error {
	return validation.Errors{
		"max_messages": validation.Validate(
			&c.MaxMessages,
			validation.Required.Error(errors.ErrorFieldRequired),
			validation.Min(1), validation.Max(10),
		),
		"subscriber": validation.Validate(
			&c.Subscriber,
			validation.Required.Error(errors.ErrorFieldRequired),
		),
	}.Filter()
	/*return validation.ValidateStruct(
		&c,
		validation.Field(&c.MaxMessages,
			validation.Required.Error(errors.ErrorFieldRequired),
			validation.Min(1), validation.Max(10),
		),
		validation.Field(&c.Subscriber,
			validation.Required.Error(errors.ErrorFieldRequired),
		),
	)*/
}

type DeleteDeadLetterQueueMessagesRequest struct {
	Messages   []Message `json:"messages" binding:"required"`
	Subscriber string    `json:"subscriber" binding:"required"`
}

type DeleteDeadLetterQueueMessagesResponse struct {
	Failed []Message `json:"failed"`
}

type UUID interface {
	GetUUID() string
}

type UUIDImpl struct {
}

func (u UUIDImpl) GetUUID() string {
	return uuid.NewV4().String()
}
