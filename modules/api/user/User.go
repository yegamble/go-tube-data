package user

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/database"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"github.com/yegamble/go-tube-api/modules/api/video"
	"gorm.io/gorm"
	"os/user"
	"time"
)

type User struct {
	gorm.Model
	ID           uuid.UUID     `json:"id" json:"id" form:"id" gorm:"primary_key"`
	FirstName    string        `json:"first_name" form:"first_name" gorm:"type:text" validate:"required,min=1,max=30"`
	LastName     string        `json:"last_name" form:"last_name" gorm:"type:text" validate:"required,min=1,max=30"`
	Email        string        `json:"email" form:"email" gorm:unique",type:text" validate:"required,min=6,max=32"`
	Username     string        `json:"username" form:"username" gorm:"unique" validate:"required,alphanum,min=1,max=32"`
	Password     string        `json:"password" form:"password" gorm:"type:text" validate:"required,min=8,max=120"`
	DisplayName  string        `json:"display_name" form:"display_name" gorm:"unique" validate:"max=50"`
	DateOfBirth  time.Time     `json:"date_of_birth" form:"date_of_birth" gorm:"type:datetime" validate:"required,datetime"`
	Gender       string        `json:"gender" form:"gender" gorm:"type:text"`
	CurrentCity  string        `json:"current_city" form:"current_city" gorm:"type:text"`
	HomeTown     string        `json:"hometown" form:"hometown" gorm:"type:text"`
	Bio          string        `json:"bio" form:"bio" gorm:"type:varchar"`
	ProfilePhoto string        `json:"profile_photo" form:"profile_photo" gorm:"type:text"`
	HeaderPhoto  string        `json:"header_photo" form:"header_photo" gorm:"type:text"`
	Videos       []video.Video `json:"videos" form:"videos" gorm:"type:array"`
	Subscribers  []*user.User  `json:"subscribers" form:"subscribers" gorm:"type:array"`
	PGPKey       string        `json:"pgp_key" form:"pgp_key" gorm:"type:text"`
	IsBanned     bool          `json:"is_Banned" form:"is_banned" gorm:"type:bool"`
	LastActive   string        `json:"last_active" gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

type WatchLater struct {
	User      User
	Video     video.Video
	CreatedAt time.Time
}

type Block struct {
	User      User
	BlockedID User
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func RegisterUserFormParser(c *fiber.Ctx) ([]*handler.ErrorResponse, error) {

	var body User

	err := c.BodyParser(&body)
	if err != nil {
		return nil, err
	}

	formErr := ValidateStruct(&body)
	if formErr != nil {
		return formErr, nil
	}

	return nil, nil
}

func createUser(u *User) error {
	db := database.DBConn

	//var body User
	u.ID = uuid.New()
	db.Create(u)

	return nil
}

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
