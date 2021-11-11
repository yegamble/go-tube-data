package models

import (
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/database"
)

type Subscription struct {
	ID               uint64    `json:"id" gorm:"primary_key"`
	UID              uuid.UUID `json:"uid"`
	UserUUID         uuid.UUID `json:"user_uuid" form:"user_uuid" gorm:"type:varchar(255);"`
	SubscribedToUUID uuid.UUID `json:"subscribed_to_uuid" form:"subscribed_to_uuid"`
	SubscribedTo     User      `gorm:"foreignKey:SubscribedToUUID;references:uuid;OnUpdate:CASCADE,OnDelete:CASCADE;varchar(255)"`
}

func (sub *Subscription) Subscribe(user *User, channel *User) error {
	db := database.DBConn
	sub.UserUUID = user.UUID
	sub.SubscribedToUUID = channel.UUID
	err := db.Create(&sub).Error
	if err != nil {
		return err
	}

	return nil
}

func (sub *Subscription) Unsubscribe() error {
	db := database.DBConn
	err := db.Delete(&sub).Error
	if err != nil {
		return err
	}

	return nil
}
