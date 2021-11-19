package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type MessageThread struct {
	ID          uuid.UUID             `json:"id" gorm:"primary_key"`
	UserID      uuid.UUID             `gorm:"type:uuid;not null"`
	Title       *string               `json:"body" gorm:"type:varchar(255)"`
	Messages    []*Message            `json:"messages"`
	MessageLogs []*MessageActivityLog `json:"messageLog" gorm:"type:varchar(255);not null;"`
	CreatedAt   time.Time             `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type MessageActivityLog struct {
	MessageThreadID uuid.UUID
	*Log
}

type MessageThreadParticipant struct {
	ID              uuid.UUID     `json:"id" gorm:"primary_key"`
	MessageThreadID uuid.UUID     `json:"message_thread_id"`
	MessageThread   MessageThread `gorm:"type:varchar(255);constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID          uuid.UUID     `json:"user_id"`
	User            User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Nickname        *string       `json:"nickname"`
	CreatedAt       time.Time     `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       gorm.DeletedAt
}

type Message struct {
	ID              uuid.UUID     `json:"id" gorm:"primary_key"`
	MessageThreadID uuid.UUID     `json:"message_thread_id"`
	UserID          uuid.UUID     `json:"sender_id"`
	User            User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RecipientID     uuid.UUID     `json:"recipient_id"`
	Recipient       User          `gorm:"foreignkey:RecipientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Body            *string       `json:"body" gorm:"type:varchar(max)"`
	ReplyTo         *Message      `json:"reply_to" gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Attachments     []*Attachment `json:"attachments"`
	CreatedAt       time.Time     `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       gorm.DeletedAt
}

type Attachment struct {
	UUID      uuid.UUID `json:"id" gorm:"primary_key"`
	MessageID uuid.UUID `json:"message_id"`
	FileName  *string   `json:"file_name"`
	FileType  *string   `json:"file_type"`
	FileSize  *int64    `json:"file_size"`
	FileURL   *string   `json:"file_url"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	thread                    MessageThread
	threadParticipant         MessageThreadParticipant
	messageThreadParticipants []MessageThreadParticipant
	message                   Message
	attachments               []Attachment
)

func (threadParticipant *MessageThreadParticipant) BeforeCreate(*gorm.DB) (err error) {
	if threadParticipant.ID == uuid.Nil {
		threadParticipant.ID = uuid.New()
	}
	return
}

func (MessageActivityLog *MessageActivityLog) BeforeCreate(*gorm.DB) (err error) {

	return
}

func (thread *MessageThread) BeforeCreate(*gorm.DB) (err error) {
	if thread.ID == uuid.Nil {
		thread.ID = uuid.New()
	}

	return
}

func (user *User) MakeActivityLogStruct(activity string, IPAddress string) (log *MessageActivityLog, err error) {

	log = &MessageActivityLog{
		Log: &Log{
			UserID:    user.ID,
			Activity:  &activity,
			IPAddress: &IPAddress,
			CreatedAt: time.Now(),
		},
	}
	return log, err
}

func (user *User) CreateMessageThread(users []User, ipAddress string) error {

	tx := db.Begin()
	thread.UserID = user.ID
	log, err := user.MakeActivityLogStruct("new message thread created", ipAddress)
	if err != nil {
		tx.Rollback()
		return err
	}

	thread.MessageLogs = append(thread.MessageLogs, log)

	err = db.Create(&thread).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, u := range users {
		threadParticipant.MessageThreadID = thread.ID
		threadParticipant.UserID = u.ID
		messageThreadParticipants = append(messageThreadParticipants, threadParticipant)
	}

	err = db.Create(&messageThreadParticipants).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (user *User) CreateMessage(thread *MessageThread, messageBody string, ipAddress string, replyToMessage *Message) error {

	tx := db.Begin()
	message.MessageThreadID = thread.ID
	message.UserID = user.ID
	message.Body = &messageBody
	message.ReplyTo = replyToMessage
	log, err := user.MakeActivityLogStruct("new message created", ipAddress)
	if err != nil {
		tx.Rollback()
		return err
	}

	thread.MessageLogs = append(thread.MessageLogs, log)

	err = db.Create(&message).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
