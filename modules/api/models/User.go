package models

import (
	"errors"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/twinj/uuid"
	"github.com/yegamble/go-tube-api/modules/api/config"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type User struct {
	ID           uint64            `json:"id" json:"id" form:"id" gorm:"primary_key"`
	UID          uuid.UUID         `json:"uid" form:"uid" gorm:"->;<-:create;unique;type:varchar(255);not null"`
	FirstName    string            `json:"first_name,omitempty" form:"first_name" gorm:"type:varchar(100);not null" validate:"min=1,max=30"`
	LastName     string            `json:"last_name,omitempty" form:"last_name" gorm:"type:varchar(100);not null" validate:"min=1,max=30"`
	Email        *string           `json:"email,omitempty" form:"email" gorm:"unique;not null;type:varchar(100)" validate:"email,required,min=6,max=32"`
	Username     *string           `json:"username" form:"username" gorm:"unique;type:varchar(100);not null" validate:"required,alphanum,min=1,max=32"`
	Password     string            `json:"-" form:"password" gorm:"type:varchar(100)" validate:"required,min=8,max=120"`
	DisplayName  *string           `json:"display_name,omitempty" form:"display_name" gorm:"type:varchar(100)" validate:"max=50"`
	DateOfBirth  *time.Time        `json:"date_of_birth,omitempty" form:"date_of_birth" gorm:"type:datetime;not null" validate:"required"`
	Gender       *string           `json:"gender,omitempty" form:"gender" gorm:"type:varchar(100)"`
	CurrentCity  *string           `json:"current_city,omitempty" form:"current_city" gorm:"type:varchar(255)"`
	HomeTown     *string           `json:"hometown,omitempty" form:"hometown" gorm:"type:varchar(255)"`
	Bio          string            `json:"bio,omitempty" form:"bio" gorm:"type:varchar(255)"`
	ProfilePhoto *string           `json:"profile_photo,omitempty" form:"profile_photo" gorm:"type:varchar(255)"`
	HeaderPhoto  *string           `json:"header_photo,omitempty" form:"header_photo" gorm:"type:varchar(255)"`
	PGPKey       *string           `json:"pgp_key,omitempty" form:"pgp_key" gorm:"type:text"`
	Videos       []Video           `json:"videos,omitempty"`
	WatchLater   []WatchLaterQueue `json:"watch_later,omitempty"`
	IsAdmin      bool              `json:"is_admin" form:"is_banned" gorm:"type:bool"`
	IsModerator  bool              `json:"is_moderator" form:"is_banned" gorm:"type:bool"`
	IsBanned     bool              `json:"is_banned" form:"is_banned" gorm:"type:bool"`
	LastActive   time.Time         `json:"last_active"  gorm:"autoCreateTime"`
	CreatedAt    time.Time         `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type ChannelProfile struct {
	ID           uint64
	UID          uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	Username     string
	DisplayName  string
	DateOfBirth  time.Time
	Gender       string
	CurrentCity  string
	HomeTown     string
	Bio          string
	ProfilePhoto string
	HeaderPhoto  string
}

type UserSettings struct {
	ID                  uint64
	UserID              uint64
	User                User      `json:"user_id" form:"user_id" gorm:"foreignKey:UserID;references:ID"`
	EmailVisible        bool      `json:"email_visible" form:"email_visible" gorm:"type:bool"`
	DateOfBirthVisible  bool      `json:"date_of_birth_visible" form:"date_of_birth_visible" gorm:"type:bool"`
	GenderVisible       bool      `json:"gender_visible" form:"gender_visible" gorm:"type:bool"`
	CurrentCityVisible  bool      `json:"current_city_visible" form:"current_city_visible" gorm:"type:bool"`
	LastActiveVisible   bool      `json:"last_active_visible" form:"last_active_visible" gorm:"type:bool"`
	OnlineStatusVisible bool      `json:"online_status_visible" form:"online_status_visible" gorm:"type:bool"`
	CreatedAt           time.Time `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type UserBlock struct {
	ID            uint64 `json:"id" json:"id" form:"id" gorm:"primary_key"`
	UserID        uint64 `json:"user_id" form:"user_id" gorm:"not null"`
	User          User   `gorm:"foreignKey:UserID;references:ID;not null"`
	BlockedUserID User   `json:"blocked_user_id" form:"blocked_user_id" gorm:"foreignKey:UserID;references:ID; not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

var (
	authUser User
	user     User
	users    []User
	limit    int
	page     int
)

func init() {
	//TODO: add session initializer here
}

func isAdmin(u User) bool {
	db.First(&u)
	return u.IsAdmin == true
}

func CreateUsers(users *[]User) error {
	tx := db.CreateInBatches(&users, len(*users))
	if tx.Error != nil {
		return tx.Error
	}

	tx.Commit()

	return nil
}

func CreateUser(u *User) error {

	if uuid.IsNil(u.UID) {
		u.UID = uuid.NewV4()
	}

	err := db.Create(&u).Error
	if err != nil {
		return err
	}

	err = CreateUserSettings(u)
	if err != nil {
		return err
	}

	return nil
}

//func getPGPFingerprint(u User) (string, error){
//
//	openpgp.
//	if u.PGPKey == nil {
//		return "", errors.New("PGP key not found")
//	}
//	return u.PGPKey.PublicKey.KeyIdString(), nil
//}

func CreateUserSettings(u *User) error {

	var userSettings UserSettings
	userSettings.User = *u

	err := db.Create(&userSettings).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserByID(userID uint64) error {
	return db.Where("uid = ?", userID).Delete(&User{}).Error
}

func DeleteUserByUID(uuid *uuid.UUID) error {
	return db.Where("uid = ?", uuid).Delete(&User{}).Error
}

func ValidateUserStruct(user *User) []*handler.ErrorResponse {
	var errors []*handler.ErrorResponse
	var element handler.ErrorResponse
	validate := validator.New()

	results := db.Where("username = ?", user.Username).First(&user)
	if results.Row() != nil {
		element.FailedField = "username"
		element.Tag = "unique"
		element.Value = *user.Username
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
	CheckAuthorisationIsValid Check
**/

func isUserBanned(u *User) (*bool, error) {
	err := db.First(&u, "id = ?", u.ID).Error
	if err != nil {
		return nil, err
	}

	return &u.IsBanned, nil
}

func isUserAdmin(u *User) (*bool, error) {
	err := db.First(&u, "id = ?", u.ID).Error
	if err != nil {
		return nil, err
	}

	return &u.IsBanned, nil
}

/**
Search for User
**/

func GetUserByID(id string) (*User, error) {
	err := db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUID(uid uuid.UUID) (*User, error) {
	tx := db.Begin()
	err := tx.First(&user, "uid = ?", uid).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {

	err := db.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func SearchUsersByName(searchTerm string, limit int, page int) (*[]User, error) {

	if page < 0 {
		return nil, errors.New("page cannot be negative")
	}

	offset := (page - 1) * config.GetResultsLimit()

	err := db.Select("username,id,uid,display_name,first_name,last_name").Offset(offset).Limit(config.UserResultsLimit).Where("username LIKE ?", "%"+searchTerm+"%").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func SearchUsersByUsername(c *fiber.Ctx) error {

	var users *[]User

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

	users, err = SearchUsersByName(username, limit, page)
	if err != nil {
		return err
	}

	if len(*users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "false",
			"message": "Profile Not Found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

func HidePrivateFields(user *User) error {

	userSettings := UserSettings{}
	err := db.First(&userSettings, "user_id = ?", user.ID).Error
	if err != nil {
		return err
	}

	if !userSettings.DateOfBirthVisible {
		user.DateOfBirth = nil
	}

	if !userSettings.EmailVisible {
		user.Email = nil
	}

	if !userSettings.GenderVisible {
		user.Gender = nil
	}

	if !userSettings.CurrentCityVisible {
		user.CurrentCity = nil
	}

	return nil
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
