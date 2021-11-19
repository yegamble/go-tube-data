package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type MessageThread struct {
	ID         uuid.UUID             `json:"id" gorm:"primary_key"`
	UserID     uuid.UUID             `gorm:"type:uuid;not null"`
	Title      *string               `json:"body" gorm:"type:varchar(255)"`
	Messages   []*Message            `json:"messages"`
	MessageLog []*MessageActivityLog `json:"messageLog" gorm:"type:varchar(255);not null;"`
	CreatedAt  time.Time             `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt  time.Time             `json:"updated_at"`
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
	ID              uuid.UUID `json:"id" gorm:"primary_key"`
	MessageThreadID uuid.UUID `json:"message_thread_id"`
	UserID          uuid.UUID `json:"sender_id"`
	User            User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RecipientID     uuid.UUID `json:"recipient_id"`
	Recipient       User      `gorm:"foreignkey:RecipientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Body            *string   `json:"body" gorm:"type:varchar(max)"`
	ReplyTo         *Message  `json:"reply_to" gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt       time.Time `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt
}

type Attachment struct {
	UUID      uuid.UUID `json:"id" gorm:"primary_key"`
	MessageID uuid.UUID `json:"message_id"`
	Message   Message   `gorm:"foreignkey:MessageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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

func (user *User) CreateMessageThread(users []User, ipAddress string) error {

	tx := db.Begin()
	thread.UserID = user.ID
	activityString := "new message thread created"
	log := &MessageActivityLog{
		Log: &Log{
			UserID:    user.ID,
			Activity:  &activityString,
			IPAddress: &ipAddress,
			CreatedAt: time.Now(),
		},
	}
	thread.MessageLog = append(thread.MessageLog, log)

	err := db.Create(&thread).Error
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
