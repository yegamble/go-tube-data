package user

import (
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/modules/api/video"
)

type Comment struct {
	ID    uuid.UUID `json:"id" gorm:"primary_key"`
	User  User
	Video video.Video
	Vote  []CommentVotes
	Text  string `json:"text" gorm:"type:text"`
}

type CommentVotes struct {
	ID       uuid.UUID `json:"id" gorm:"primary_key"`
	User     User
	Reaction bool `json:"reaction" gorm:"type:boolean"`
}
