package models

import (
	"github.com/google/uuid"
)

type Subscription struct {
	UUID        uuid.UUID `gorm:"type:varchar(255);index;primaryKey;"`
	ChannelUUID uuid.UUID `json:"subscribed_to_uuid" form:"subscribed_to_uuid"`
	Channel     User      `gorm:"foreignKey:ChannelUUID;references:uuid;OnUpdate:CASCADE,OnDelete:CASCADE;varchar(255)"`
}
