package models

import (
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

func (user *User) UnsubscribeFromChannel(channelId uuid.UUID) error {

	tx := db.Begin()

	channel := User{}
	err := tx.First(&channel, "id = ?", channelId).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&user).Association("Subscriptions").Delete(&channel)
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

func (user *User) GetSubscriptions() error {
	subscriptions := []User{}
	err := db.Model(&user).Association("Subscriptions").Find(&subscriptions)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) GetSubscribers() error {
	err := db.Model(&user).Where("channel_id = ?", user.ID).Association("Subscriptions").Find(&user.Subscriptions)
	if err != nil {
		return err
	}

	return nil
}
