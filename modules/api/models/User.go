package models

import (
	"errors"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/modules/api/config"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"gorm.io/gorm"
	"log"
	"time"
)

type User struct {
	gorm.Model
	ID            uint64              `json:"id" json:"id" form:"id" gorm:"primary_key"`
	UUID          uuid.UUID           `json:"uuid" form:"uuid" gorm:"->;<-:create;unique;type:varchar(255);not null"`
	FirstName     string              `json:"first_name,omitempty" form:"first_name" gorm:"type:varchar(100);not null" validate:"min=1,max=30"`
	LastName      string              `json:"last_name,omitempty" form:"last_name" gorm:"type:varchar(100);not null" validate:"min=1,max=30"`
	Email         *string             `json:"email,omitempty" form:"email" gorm:"unique;not null;type:varchar(100)" validate:"email,required,min=6,max=32"`
	Username      *string             `json:"username" form:"username" gorm:"unique;type:varchar(100);not null" validate:"required,alphanum,min=1,max=32"`
	Password      string              `json:"-" form:"password" gorm:"type:varchar(100)" validate:"required,min=8,max=120"`
	DisplayName   *string             `json:"display_name,omitempty" form:"display_name" gorm:"type:varchar(100)" validate:"max=100"`
	DateOfBirth   *time.Time          `json:"date_of_birth,omitempty" form:"date_of_birth" gorm:"type:datetime;not null" validate:"required"`
	Gender        *string             `json:"gender,omitempty" form:"gender" gorm:"type:varchar(100)"`
	CurrentCity   *string             `json:"current_city,omitempty" form:"current_city" gorm:"type:varchar(255)"`
	Hometown      *string             `json:"hometown,omitempty" form:"hometown" gorm:"type:varchar(255)"`
	Bio           *string             `json:"bio,omitempty" form:"bio" gorm:"type:varchar(255)"`
	ProfilePhoto  *string             `json:"profile_photo,omitempty" form:"profile_photo" gorm:"type:varchar(255)"`
	HeaderPhoto   *string             `json:"header_photo,omitempty" form:"header_photo" gorm:"type:varchar(255)"`
	PGPKey        *string             `json:"pgp_key,omitempty" form:"pgp_key" gorm:"type:text"`
	Settings      UserSettings        `json:"settings" gorm:"foreignKey:UserUUID;references:UUID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Videos        []Video             `json:"videos,omitempty" gorm:"foreignKey:UserUUID;references:UUID;OnUpdate:CASCADE,OnDelete:SET NULL;"`
	WatchLater    []WatchLaterQueue   `json:"watch_later,omitempty" gorm:"foreignKey:UserUUID;references:UUID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Subscriptions []Subscription      `json:"subscriptions,omitempty" gorm:"foreignKey:UserUUID;references:UUID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserPlaylist  []UserPlaylist      `json:"playlist,omitempty" gorm:"foreignKey:UserUUID;references:UUID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags          []*Tag              `json:"tags,omitempty" gorm:"many2many:user_tags;foreignKey:UUID;joinForeignKey:UserUUID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BlockedUsers  []BlockedUserRecord `json:"blocked_users,omitempty" gorm:"foreignKey:UserUUID;references:UUID;OnUpdate:CASCADE,OnDelete:CASCADE;type:varchar(255);"`
	Logs          []IPLog             `json:"logs,omitempty" gorm:"foreignKey:UserUUID;references:UUID;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Admin         bool                `json:"is_admin" form:"is_admin" gorm:"type:bool;default:0"`
	Moderator     bool                `json:"is_moderator" form:"is_banned" gorm:"type:bool;default:0"`
	Banned        bool                `json:"is_banned" form:"is_banned" gorm:"type:bool;default:0"`
	Private       bool                `json:"is_private" form:"is_private" gorm:"type:bool;default:0"`
	LastActive    time.Time           `json:"last_active"  gorm:"autoCreateTime"`
	CreatedAt     time.Time           `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

type UserSettings struct {
	ID                  uint64
	UserUUID            uuid.UUID `json:"user_id" form:"user_id" gorm:"varchar(255);"`
	EmailVisible        bool      `json:"email_visible" form:"email_visible" gorm:"type:bool"`
	DateOfBirthVisible  bool      `json:"date_of_birth_visible" form:"date_of_birth_visible" gorm:"type:bool"`
	GenderVisible       bool      `json:"gender_visible" form:"gender_visible" gorm:"type:bool"`
	CurrentCityVisible  bool      `json:"current_city_visible" form:"current_city_visible" gorm:"type:bool"`
	LastActiveVisible   bool      `json:"last_active_visible" form:"last_active_visible" gorm:"type:bool"`
	OnlineStatusVisible bool      `json:"online_status_visible" form:"online_status_visible" gorm:"type:bool"`
	CreatedAt           time.Time `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type BlockedUserRecord struct {
	ID              uint64    `json:"id" json:"id" form:"id" gorm:"primary_key"`
	UserUUID        uuid.UUID `json:"user_uuid" form:"user_uuid"`
	BlockedUserUUID uuid.UUID `json:"blocked_user_uuid" form:"blocked_user_uuid"`
	BlockedUserID   User      `json:"blocked_user_id" form:"blocked_user_id" gorm:"foreignKey:BlockedUserUUID;references:uuid; not null;OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}

var (
	user  User
	users []User
	page  int
)

func CreateUsers(users *[]User) error {
	tx := db.Begin()

	err := db.Create(&users)
	if err != nil {
		return err.Error
	}

	tx.Commit()

	return nil
}

func (user *User) BeforeCreate(*gorm.DB) (err error) {
	user.UUID = uuid.New()
	user.LastActive = time.Now()
	user.Settings = UserSettings{
		UserUUID: user.UUID,
	}

	return
}

func (user *User) Create(ipAddress string) error {

	tx := db.Begin()

	err := tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	newUserLog := user.CreateUserLog("registered", ipAddress)
	err = tx.Create(&newUserLog).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (user *User) isBlocked(u User) (bool, error) {
	var blockedUser User
	err := db.Where("blocked_user_uuid = ? AND user_uuid = ?", user.UUID, u.UUID).First(&blockedUser).Error
	if err == gorm.ErrRecordNotFound || db.Row() == nil {
		return false, err
	}
	if err != nil {
		return true, err
	}

	return true, nil
}

func (user *User) CreateTags(tagsArray []*Tag) error {

	tx := db.Begin()

	for i, tag := range tagsArray {
		log.Println(len(tagsArray))
		err := tag.findTag(tag.Value)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		} else if err != gorm.ErrRecordNotFound {
			tagsArray = append(tagsArray[:i], tagsArray[i+1:]...)
		}
	}

	err := tx.Model(&user).Association("Tags").Append(tagsArray)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}

func (user *User) CreateWatchLaterQueue() error {
	watchQueue := WatchLaterQueue{
		UserUUID: user.UUID,
	}

	err := db.Create(&watchQueue).Error
	if err != nil {
		return err
	}

	return nil
}

func (user *User) Delete() error {
	err := db.Delete(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserByID(uuid uuid.UUID) error {
	return db.Where("id = ?", uuid).Delete(&User{}).Error
}

func ValidateUserStruct(user *User) []*handler.ErrorResponse {
	var errorResponses []*handler.ErrorResponse
	var element handler.ErrorResponse
	validate := validator.New()

	results := db.Where("username = ?", user.Username).First(&user)
	if results.Row() != nil {
		element.FailedField = "username"
		element.Tag = "unique"
		element.Value = *user.Username
		errorResponses = append(errorResponses, &element)
	}

	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errorResponses = append(errorResponses, &element)
		}
	}

	return errorResponses
}

/**
	CheckAuthorisationIsValid Check
**/

func (user User) isBanned() bool {
	err := db.First(&user).Error
	if err != nil {
		return false
	}

	return user.Banned
}

func (user *User) isAdmin() bool {
	db.First(&user)
	return user.Admin == true
}

func (user *User) isMod() bool {
	db.First(&user)
	return user.Admin == true
}

/**
Search for User
**/

func GetUserByUUID(uid uuid.UUID) (*User, error) {
	tx := db.Begin()
	err := tx.First(&user, "uuid = ?", uid).Error
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

func SearchUsersByName(searchTerm string, page int) (*[]User, error) {

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

func (user *User) HidePrivateFields() error {

	userSettings := UserSettings{}
	err := db.First(&userSettings, "user_uid = ?", user.UUID).Error
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
