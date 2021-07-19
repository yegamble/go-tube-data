package video

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Video struct {
	gorm.Model
	ID 	  				uuid.UUID  		`json:"id" gorm:"primary_key"`
	ShortID 			string     		`json:"short_id" gorm:"unique"`
	Title  				string     		`json:"title" gorm:"required"`
	UserID				string	   		`json:"user_id`
	Description 		string			`json:"description`
	Thumbnail   		string			`json:"thumbnail""`
	Resolutions			[]VidRes 		`json:"title" gorm:"required,type:array"`
	MaxResolution       string 			`json:"max_resolution"`
	PublishedAt 		time.Time  		`json:"published_at"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt
}

type VidRes struct {
	Resolution string
}