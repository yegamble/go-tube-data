package user

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/database"
	"github.com/yegamble/go-tube-api/modules/api/errorhandler"
	"github.com/yegamble/go-tube-api/modules/api/video"
	"gorm.io/gorm"
	"os/user"
	"time"
)

type User struct {
	gorm.Model
	ID           uuid.UUID     `json:"id" gorm:"primary_key"`
	FirstName    string        `json:"first_name" gorm:"type:text" validate:"required,min=1,max=30"`
	LastName     string        `json:"last_name" gorm:"type:text" validate:"required,min=1,max=30"`
	Email        string        `json:"email" gorm:unique",type:text" validate:"required,email,min=6,max=32"`
	Username     string        `json:"username" gorm:"unique" validate:"required,alphanum,min=1,max=32"`
	Password     string        `json:"password" gorm:"type:text" validate:"required,min=8,max=120"`
	DisplayName  string        `json:"username" gorm:"unique" validate:"max=50"`
	DateOfBirth  time.Time     `json:"date_of_birth" gorm:"type:datetime" validate:"required"`
	Gender       string        `json:"gender" gorm:"type:text"`
	CurrentCity  string        `json:"current_city" gorm:"type:text"`
	HomeTown     string        `json:"hometown" gorm:"type:text"`
	Bio          string        `json:"bio" gorm:"type:text"`
	ProfilePhoto string        `json:"profile_photo" gorm:"type:text"`
	HeaderPhoto  string        `json:"header_photo" gorm:"type:text"`
	Videos       []video.Video `json:"videos" gorm:"type:array"`
	Subscribers  []*user.User  `json:"subscribers" gorm:"type:array"`
	PGPKey       string        `json:"pgp_key" gorm:"type:text"`
	IsBanned     bool          `json:"pgp_key" gorm:"type:bool"`
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

func UserFormParser(c *fiber.Ctx) ([]*errorhandler.ErrorResponse, error) {

	var body User

	formErr := ValidateStruct(&body)
	if formErr != nil {
		return formErr, nil
	}

	err := c.BodyParser(&body)
	if err != nil {
		return nil, err
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

func ValidateStruct(user *User) []*errorhandler.ErrorResponse {
	var errors []*errorhandler.ErrorResponse
	validate := validator.New()

	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element errorhandler.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
