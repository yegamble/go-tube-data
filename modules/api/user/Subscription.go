package user

import (
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/database"
	"os/user"
)

type Subscription struct {
	ID             int64 `json:"id" gorm:"primary_key"`
	UID            uuid.UUID
	UserID         int64
	User           User `json:"user_id" form:"user_id" gorm:"foreignKey:UserID;references:ID"`
	SubscribedToID int64
	SubscribedTo   user.User `json:"subscribed_to" form:"subscribed_to" gorm:"foreignKey:UserID;references:ID"`
}

func Subscribe(u User, subbedUser User) {
	db := database.DBConn
	var sub Subscription
	sub.UserID = u.ID
	sub.SubscribedToID = subbedUser.ID
	db.Create(&sub)
}

func Unsubscribe(u User, subbedUser User) {
	db := database.DBConn
	var sub Subscription
	sub.UserID = u.ID
	sub.SubscribedToID = subbedUser.ID
	db.Create(&sub)
}
