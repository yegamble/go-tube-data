package models

//type Subscription struct {
//	UUID        uuid.UUID `gorm:"type:varchar(255);index;primaryKey;"`
//	UserUUID    *uuid.UUID
//	ChannelUUID *uuid.UUID `json:"channel_uuid" form:"channel_uuid"`
//}

//func (subscription *Subscription) BeforeCreate(*gorm.DB) error {
//    subscription.UUID = uuid.New()
//    return nil
//}
//
//func (user *User) SubscribeToChannel(userUUID uuid.UUID) error {
//
//	tx := db.Begin()
//	user.Subscriptions = []*Subscription{
//		{
//            ChannelUUID: &userUUID,
//        },
//    }
//
//	err := tx.Save(&user).Error
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//
//	tx.Commit()
//    return nil
//}
//
//func (subscription *Subscription) UnsubscribeFromChannel() error {
//    tx := db.Begin()
//    err := tx.Delete(&subscription).Error
//    if err != nil {
//        tx.Rollback()
//        return err
//    }
//
//    tx.Commit()
//    return nil
//}
