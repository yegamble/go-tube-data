package models

import (
	"github.com/google/uuid"
)

type Comment struct {
	ID    uuid.UUID `json:"id" gorm:"primary_key"`
	User  User
	Video Video
	Vote  []CommentVotes
	Text  string `json:"text" gorm:"type:text"`
}

type CommentVotes struct {
	ID       uuid.UUID `json:"id" gorm:"primary_key"`
	User     User
	Reaction bool `json:"reaction" gorm:"type:boolean"`
}
