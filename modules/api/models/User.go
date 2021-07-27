package models

import (
	"errors"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/modules/api/auth"
	"github.com/yegamble/go-tube-api/modules/api/config"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID           uint64            `json:"id" json:"id" form:"id" gorm:"<-:create;primary_key"`
	UID          uuid.UUID         `json:"uid" form:"uid" gorm:"->;unique;type:varchar(255)"`
	FirstName    string            `json:"first_name,omitempty" form:"first_name" gorm:"type:varchar(100)" validate:"min=1,max=30"`
	LastName     string            `json:"last_name,omitempty" form:"last_name" gorm:"type:varchar(100)" validate:"min=1,max=30"`
	Email        string            `json:"email,omitempty" form:"email" gorm:unique",type:varchar(100)" validate:"required,min=6,max=32"`
	Username     string            `json:"username" form:"username" gorm:"unique;type:varchar(100)" validate:"required,alphanum,min=1,max=32"`
	Password     string            `json:"-" form:"password" gorm:"type:varchar(100)" validate:"required,min=8,max=120"`
	DisplayName  string            `json:"display_name,omitempty" form:"display_name" gorm:"type:varchar(100)" validate:"max=50"`
	DateOfBirth  time.Time         `json:"date_of_birth" form:"date_of_birth" gorm:"type:datetime" validate:"required"`
	Gender       string            `json:"gender,omitempty" form:"gender" gorm:"type:varchar(100)"`
	CurrentCity  string            `json:"current_city,omitempty" form:"current_city" gorm:"type:varchar(255)"`
	HomeTown     string            `json:"hometown,omitempty" form:"hometown" gorm:"type:varchar(255)"`
	Bio          string            `json:"bio,omitempty" form:"bio" gorm:"type:varchar(255)"`
	ProfilePhoto string            `json:"profile_photo,omitempty" form:"profile_photo" gorm:"type:varchar(255)"`
	HeaderPhoto  string            `json:"header_photo,omitempty" form:"header_photo" gorm:"type:varchar(255)"`
	PGPKey       string            `json:"pgp_key,omitempty" form:"pgp_key" gorm:"type:text"`
	Videos       []Video           `json:"videos,omitempty"`
	WatchLater   []WatchLaterQueue `json:"watch_later,omitempty"`
	IsBanned     bool              `json:"is_Banned" form:"is_banned" gorm:"type:bool"`
	LastActive   time.Time         `json:"last_active"  gorm:"autoCreateTime"`
	CreatedAt    time.Time         `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type WatchLaterQueue struct {
	ID        uint64
	UserID    uint64
	User      User      `json:"user_id" form:"user_id" gorm:"foreignKey:UserID;references:ID"`
	VideoID   uuid.UUID `json:"video_id" form:"video_id"`
	Video     Video     `gorm:"foreignKey:VideoID;references:ID"`
	CreatedAt time.Time
}

type UserBlock struct {
	ID            uint64 `json:"id" json:"id" form:"id" gorm:"primary_key"`
	UserID        uint64 `json:"user_id" form:"user_id"`
	User          User   `gorm:"foreignKey:UserID;references:ID"`
	BlockedUserID User   `json:"blocked_user_id" form:"blocked_user_id" gorm:"foreignKey:UserID;references:ID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

var (
	user  User
	users []User
	limit int
	page  int
)

func RegisterUser(c *fiber.Ctx) error {

	var body User

	body.UID = uuid.New()
	body.LastActive = time.Now()
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	auth.EncodeToArgon(&body.Password)

	formErr := ValidateUserStruct(&body)
	if formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(formErr)
	}

	err = db.Create(&body).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	CreateUserLog("registered", body, c)

	return c.Status(fiber.StatusOK).JSON(body.UID.String())
}

func Login(c *fiber.Ctx) error {
	var user User
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON("Username field is empty")
	}
	if password == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON("Password field is empty")
	}

	//compare the user from the request, with the one we defined:
	err := db.Where("username = ?", c.FormValue("username")).First(&user).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON("invalid login details")
	}

	match, err := auth.ComparePasswordAndHash(&password, user.Password)
	if err != nil {
		return err
	} else if !match {
		return errors.New("invalid login details")
	}

	token, err := auth.CreateJWTToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	err = SaveSession(user, c)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(token)
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie()
	return nil
}

func DeleteUser(c *fiber.Ctx) error {
	return db.Where("uid = ?", c.Query("uid")).Delete(&User{}).Error
}

func EditUser(c *fiber.Ctx) error {

	err := db.First(&user, c.Params("id")).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	err = c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	if c.FormValue("password") != "" {
		auth.EncodeToArgon(&user.Password)
	}

	if c.FormValue("uid") != "" {
		return errors.New("uid cannot be changed")
	}

	err = db.Save(user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func UploadUserPhoto(c *fiber.Ctx, photoKey string) error {

	dir := "uploads/photos/user/"

	file, err := c.FormFile(photoKey)
	if err != nil {
		return err
	}

	contentType := file.Header.Get("content-type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		return errors.New("photo is not jpeg or png")
	}

	filename, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	dst, err := os.Create(filepath.Join(dir, filepath.Base(strings.Replace(filename.String()+os.Getenv("APP_IMAGE_EXTENSION"), "-", "_", -1))))
	if err != nil {
		return err
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	err = db.Model(&User{}).Where("id = ?", c.Params("id")).Update(photoKey, dst.Name()).Error
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dst.Name())

}

func ValidateUserStruct(user *User) []*handler.ErrorResponse {
	var errors []*handler.ErrorResponse
	var element handler.ErrorResponse
	validate := validator.New()

	results := db.Where("username = ?", user.Username).First(&user)
	if results.Row() != nil {
		element.FailedField = "username"
		element.Tag = "unique"
		element.Value = user.Username
		errors = append(errors, &element)
	}

	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}

/**
Search for User
**/

func GetUserByID(c *fiber.Ctx) error {
	var user User
	db.First(&user, "id = ?", c.Query("id"))

	return c.Status(fiber.StatusOK).JSON(user)
}

func SearchUserByUID(c *fiber.Ctx) error {

	var user User
	db.First(&user, "uid = ?", c.Query("uid"))

	return c.Status(fiber.StatusOK).JSON(user)
}

func GetUserByUsername(c *fiber.Ctx) error {

	var user User
	err := db.First(&user, "username = ?", c.Params("username")).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func SearchUsersByUsername(c *fiber.Ctx) error {

	username := c.Params("*")
	if username == "" {
		return PaginateAllUsers(c)
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil && c.Query("page") != "" {
		return err
	} else if page == 0 {
		page = 1
	}

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			return err
		}
	} else {
		limit = config.GetResultsLimit()
	}

	db.Select("username,id,uid,display_name,first_name,last_name").Limit(config.UserResultsLimit).Where("username LIKE ?", "%"+username+"%").Find(&users)

	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "false",
			"message": "Profile Not Found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(&users)
}

func PaginateAllUsers(c *fiber.Ctx) error {

	offset := (page - 1) * config.GetResultsLimit()

	db.Offset(offset).Limit(config.UserResultsLimit).Find(&users)
	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "false",
			"message": "Profile Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}
