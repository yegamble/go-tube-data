package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
)

type Tag struct {
	UUID  uuid.UUID `gorm:"type:varchar(255);index;primaryKey;"`
	Value *string   `json:"value" gorm:"unique"`
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
