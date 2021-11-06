package models

import (
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/database"
)

type Subscription struct {
	ID             uint64    `json:"id" gorm:"primary_key"`
	UID            uuid.UUID `json:"uid"`
	UserID         uuid.UUID `json:"user_id" form:"user_id" gorm:"varchar(255);size:255"`
	User           User      `gorm:"foreignKey:UserID;references:ID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SubscribedToID uuid.UUID `json:"subscribed_to_id" form:"subscribed_to_id" gorm:"varchar(255);size:255;"`
	SubscribedTo   User      `gorm:"foreignKey:SubscribedToID;references:ID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *User) Subscribe(channel *User) error {
	db := database.DBConn
	var sub Subscription
	sub.UserID = u.UID
	sub.SubscribedToID = channel.UID
	err := db.Create(&sub).Error
	if err != nil {
		return err
	}

	return nil
}

func Unsubscribe(u *User, subbedUser *User) error {
	db := database.DBConn
	var sub Subscription
	sub.UserID = u.UID
	sub.SubscribedToID = subbedUser.UID
	err := db.Create(&sub).Error
	if err != nil {
		return err
	}

	return nil
}
