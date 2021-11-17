package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type MessageThread struct {
	UUID    uuid.UUID `json:"id" gorm:"primary_key"`
	UserID  uuid.UUID `gorm:"type:uuid;not null"`
	Message []Message `gorm:"foreignkey:ThreadID"`
}

type MessageThreadParticipant struct {
	UUID              uuid.UUID     `json:"id" gorm:"primary_key"`
	MessageThreadUUID uuid.UUID     `json:"message_thread_id"`
	MessageThread     MessageThread `gorm:"foreignkey:MessageThreadUUID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserUUID          uuid.UUID     `json:"user_id"`
	User              User          `gorm:"foreignkey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Nickname          *string       `json:"nickname"`
	CreatedAt         time.Time     `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt         time.Time     `json:"updated_at"`
	DeletedAt         gorm.DeletedAt
}

type Message struct {
	UUID          uuid.UUID `json:"id" gorm:"primary_key"`
	ThreadID      uuid.UUID `json:"message_thread_id"`
	SenderUUID    uuid.UUID `json:"sender_id"`
	Sender        User      `gorm:"foreignkey:UserID;references:uuid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RecipientUUID uuid.UUID `json:"recipient_id"`
	Recipient     User      `gorm:"foreignkey:RecipientUUID;references:uuid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Body          *string   `json:"body" gorm:"type:varchar(max)"`
	CreatedAt     time.Time `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     gorm.DeletedAt
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

func (threadParticipant *MessageThreadParticipant) BeforeCreate() (err error) {
	threadParticipant.UUID = uuid.New()
	return
}

func (user *User) CreateMessageThread(users []User) error {

	tx := db.Begin()
	thread.UUID = uuid.New()
	thread.UserID = user.ID

	err := db.Create(&thread).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, u := range users {
		threadParticipant.MessageThreadUUID = thread.UUID
		threadParticipant.UserUUID = u.ID
		messageThreadParticipants = append(messageThreadParticipants, threadParticipant)
	}

	err = db.Create(&messageThreadParticipants).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
