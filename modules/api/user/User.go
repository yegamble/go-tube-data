package user

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/database"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"github.com/yegamble/go-tube-api/modules/api/video"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID           int64     `json:"id" json:"id" form:"id" gorm:"primary_key"`
	UID          uuid.UUID `json:"user_id" form:"user_id" gorm:unique;type:text"`
	FirstName    string    `json:"first_name" form:"first_name" gorm:"type:text" validate:"required,min=1,max=30"`
	LastName     string    `json:"last_name" form:"last_name" gorm:"type:text" validate:"required,min=1,max=30"`
	Email        string    `json:"email" form:"email" gorm:unique",type:text" validate:"required,min=6,max=32"`
	Username     string    `json:"username" form:"username" gorm:"unique;type:varchar" validate:"required,alphanum,min=1,max=32"`
	Password     string    `json:"password" form:"password" gorm:"type:text" validate:"required,min=8,max=120"`
	DisplayName  string    `json:"display_name" form:"display_name" gorm:"type:varchar" validate:"max=50"`
	DateOfBirth  time.Time `json:"date_of_birth" form:"date_of_birth" gorm:"type:datetime" validate:"required"`
	Gender       string    `json:"gender" form:"gender" gorm:"type:varchar"`
	CurrentCity  string    `json:"current_city" form:"current_city" gorm:"type:varchar"`
	HomeTown     string    `json:"hometown" form:"hometown" gorm:"type:varchar"`
	Bio          string    `json:"bio" form:"bio" gorm:"type:varchar"`
	ProfilePhoto string    `json:"profile_photo" form:"profile_photo" gorm:"type:varchar"`
	HeaderPhoto  string    `json:"header_photo" form:"header_photo" gorm:"type:varchar"`
	PGPKey       string    `json:"pgp_key" form:"pgp_key" gorm:"type:text"`
	IsBanned     bool      `json:"is_Banned" form:"is_banned" gorm:"type:bool"`
	LastActive   time.Time `json:"last_active"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

type WatchLaterQueue struct {
	UserID    int64
	User      User        `json:"user_id" form:"user_id" gorm:"foreignKey:UserID;references:ID"`
	VideoID   uuid.UUID   `json:"video_id" form:"video_id"`
	Video     video.Video `json:"video_id" form:"video_id" gorm:"foreignKey:VideoID;references:ID"`
	CreatedAt time.Time
}

type Block struct {
	ID            int `json:"id" json:"id" form:"id" gorm:"primary_key"`
	UserID        int64
	User          User `json:"user_id" form:"video_id" gorm:"foreignKey:UserID;references:ID"`
	BlockedUserID User `json:"blocked_user_id" form:"blocked_user_id" gorm:"foreignKey:UserID;references:ID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func RegisterUser(c *fiber.Ctx) (string, []*handler.ErrorResponse, error) {

	db := database.DBConn
	var body User

	body.UID = uuid.New()
	err := c.BodyParser(&body)
	if err != nil {
		return "", nil, err
	}

	formErr := ValidateStruct(&body)
	if formErr != nil {
		return "", formErr, nil
	}

	err = db.Create(&body).Error
	if err != nil {
		return "", nil, err
	}

	return body.UID.String(), formErr, nil
}

//func EditUser(c *fiber.Ctx) (string, error){
//
//
//}
//
//func AddProfilePicture(c *fiber.Ctx) (string, error){
//
//
//}

func ValidateStruct(user *User) []*handler.ErrorResponse {
	var errors []*handler.ErrorResponse
	validate := validator.New()

	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element handler.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
