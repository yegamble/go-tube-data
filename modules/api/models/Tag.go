package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
)

type Tag struct {
	UUID     uuid.UUID `gorm:"type:varchar(255);index;primaryKey;"`
	Value    *string   `json:"value" gorm:"unique"`
	Disabled *bool     `json:"disabled" gorm:"type:boolean;default:0"`
}

func (tag *Tag) BeforeCreate(*gorm.DB) error {
	tag.UUID = uuid.New()
	return nil
}

func (tag *Tag) findTag(tagSearchString *string) error {

	err := db.First(&tag, "value = ?", *tagSearchString)
	if err != nil {
		return err.Error
	}

	return nil
}

func (user *User) findTags() error {

	err := db.Model(&user).Where("user_id = ?", user.ID).Association("Tags").Find(&user.Tags)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	return nil
}
