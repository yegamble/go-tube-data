package models

import (
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/database"
)

type Subscription struct {
	ID             uint64 `json:"id" gorm:"primary_key"`
	UID            uuid.UUID
	UserID         uint64 `json:"user_id" form:"user_id"`
	User           User   `gorm:"foreignKey:UserID;references:ID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SubscribedToID uint64 `json:"subscribed_to_id" form:"subscribed_to_id"`
	SubscribedTo   User   `gorm:"foreignKey:SubscribedToID;references:ID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func Subscribe(u *User, subbedUser *User) error {
	db := database.DBConn
	var sub Subscription
	sub.UserID = u.ID
	sub.SubscribedToID = subbedUser.ID
	err := db.Create(&sub).Error
	if err != nil {
		return err
	}

	return nil
}

func Unsubscribe(u *User, subbedUser *User) error {
	db := database.DBConn
	var sub Subscription
	sub.UserID = u.ID
	sub.SubscribedToID = subbedUser.ID
	err := db.Create(&sub).Error
	if err != nil {
		return err
	}

	return nil
}
