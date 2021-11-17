package models

import (
	"fmt"
	"github.com/google/uuid"
)

func (user *User) SubscribeToChannel(uuid uuid.UUID) error {

	tx := db.Begin()

	channel := User{}
	err := tx.First(&channel, "id = ?", uuid).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	subID := user.ID.String()
	chanID := channel.ID.String()
	fmt.Println(subID)
	fmt.Println(chanID)

	err = tx.Model(&user).Association("Subscriptions").Append(&channel)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
