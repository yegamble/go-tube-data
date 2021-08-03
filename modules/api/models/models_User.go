package models

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/twinj/uuid"
	"github.com/yegamble/go-tube-api/modules/api/auth"
	"github.com/yegamble/go-tube-api/modules/api/config"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"gorm.io/gorm"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

type User struct {
	ID           uint64            `json:"id" json:"id" form:"id" gorm:"primary_key"`
	UID          uuid.UUID         `json:"uid" form:"uid" gorm:"->;<-:create;unique;type:varchar(255);not null"`
	FirstName    string            `json:"first_name,omitempty" form:"first_name" gorm:"type:varchar(100);not null" validate:"min=1,max=30"`
	LastName     string            `json:"last_name,omitempty" form:"last_name" gorm:"type:varchar(100);not null" validate:"min=1,max=30"`
	Email        string            `json:"email,omitempty" form:"email" gorm:"unique;not null;type:varchar(100)" validate:"email,required,min=6,max=32"`
	Username     string            `json:"username" form:"username" gorm:"unique;type:varchar(100);not null" validate:"required,alphanum,min=1,max=32"`
	Password     string            `json:"-" form:"password" gorm:"type:varchar(100)" validate:"required,min=8,max=120"`
	DisplayName  string            `json:"display_name,omitempty" form:"display_name" gorm:"type:varchar(100)" validate:"max=50"`
	DateOfBirth  time.Time         `json:"date_of_birth" form:"date_of_birth" gorm:"type:datetime;not null" validate:"required"`
	Gender       string            `json:"gender,omitempty" form:"gender" gorm:"type:varchar(100)"`
	CurrentCity  string            `json:"current_city,omitempty" form:"current_city" gorm:"type:varchar(255)"`
	HomeTown     string            `json:"hometown,omitempty" form:"hometown" gorm:"type:varchar(255)"`
	Bio          string            `json:"bio,omitempty" form:"bio" gorm:"type:varchar(255)"`
	ProfilePhoto string            `json:"profile_photo,omitempty" form:"profile_photo" gorm:"type:varchar(255)"`
	HeaderPhoto  string            `json:"header_photo,omitempty" form:"header_photo" gorm:"type:varchar(255)"`
	PGPKey       string            `json:"pgp_key,omitempty" form:"pgp_key" gorm:"type:text"`
	Videos       []Video           `json:"videos,omitempty"`
	WatchLater   []WatchLaterQueue `json:"watch_later,omitempty"`
	IsAdmin      bool              `json:"is_admin" form:"is_banned" gorm:"type:bool"`
	IsModerator  bool              `json:"is_moderator" form:"is_banned" gorm:"type:bool"`
	IsBanned     bool              `json:"is_banned" form:"is_banned" gorm:"type:bool"`
	LastActive   time.Time         `json:"last_active"  gorm:"autoCreateTime"`
	CreatedAt    time.Time         `json:"created_at" gorm:"<-:create;autoCreateTime"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type WatchLaterQueue struct {
	ID        uint64
	UserID    uint64
	User      User      `json:"user_id" form:"user_id" gorm:"foreignKey:UserID;references:ID"`
	VideoID   uuid.UUID `json:"video_id" form:"video_id"`
	Video     Video     `gorm:"foreignKey:VideoID;references:ID;not null"`
	CreatedAt time.Time
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

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
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

func Auth(c *fiber.Ctx) (*User, error) {
	userToken, err := auth.VerifyToken(c)

	tokenAuth, err := auth.ExtractTokenMetadata(c)
	if err != nil {
		return nil, c.Status(fiber.StatusUnauthorized).JSON("unauthorized")
	}

	userId, err := FetchAuth(tokenAuth)
	if err != nil {
		return nil, c.Status(fiber.StatusUnauthorized).JSON("unauthorized")
	}

	user.ID = userId

	claims := userToken.Claims.(jwt.MapClaims)
	log.Println(claims["user_id"])

	return &user, err
}

func isAdmin(u User) bool {
	db.First(&u)
	return u.IsAdmin == true
}

func CreateUser(u *User) (*User, error) {
	err := db.Create(&u).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateAuth(userid uint64, td *auth.TokenDetails) error {

	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(reflect.ValueOf((*td).AccessUuid).String(), strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(reflect.ValueOf((*td).RefreshUuid).String(), strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func RefreshAuthorisation(c *fiber.Ctx) error {

	//verify the token
	token, err := jwt.Parse(c.Cookies("refresh_token"), func(token *jwt.Token) (interface{}, error) {

		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errors.New("refresh token expired"))
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(errors.New("token invalid"))
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string

		if !ok {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
		}

		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
		}
		//Delete the previous RefreshAuthorisation Token
		deleted, delErr := DeleteAuth(refreshUuid)

		if delErr != nil || deleted == 0 { //if any goes wrong
			log.Println(claims)
			return c.Status(fiber.StatusUnauthorized).JSON(errors.New("refresh token expired").Error())
		}

		//Create new pairs of refresh and access tokens
		ts, createErr := auth.CreateJWTToken(userId)
		if createErr != nil {
			return c.Status(fiber.StatusCreated).JSON(createErr.Error())
		}
		//save the tokens metadata to redis
		saveErr := CreateAuth(userId, ts)
		if saveErr != nil {
			return c.Status(fiber.StatusForbidden).JSON(err.Error())
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}

		SaveSession(userId, ts.AccessToken, ts.RefreshToken, c)

		return c.Status(fiber.StatusCreated).JSON(tokens)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(errors.New("refresh expired"))
	}
}

func DeleteAuth(refreshUuid string) (int64, error) {

	deleted := client.Del(refreshUuid)
	if deleted.Err() != nil {
		return 0, deleted.Err()
	}

	return deleted.Val(), nil
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie()
	return nil
}

func DeleteUser(c *fiber.Ctx) error {
	return db.Where("uid = ?", c.Query("uid")).Delete(&User{}).Error
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
