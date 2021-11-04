package models

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/modules/api/auth"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

/**
Authentication
*/

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
		return c.Status(fiber.StatusUnauthorized).JSON("invalid login details")
	}

	if *user.Banned {
		return c.Status(fiber.StatusUnauthorized).JSON("user is banned")
	}

	match, err := auth.ComparePasswordAndHash(&password, user.Password)
	if err != nil {
		return err
	} else if !match {
		return errors.New("invalid login details")
	}

	token, err := CreateJWTToken(user.ID, user.Admin)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	err = CreateAuthRecord(user.ID, token)
	if err != nil {
		return err
	}

	AccessToken := reflect.ValueOf((*token).AccessToken).String()
	RefreshToken := reflect.ValueOf((*token).RefreshToken).String()

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + AccessToken
	c.Set("Authorization", bearer)

	SaveUserCookies(reflect.ValueOf(AccessToken).String(), reflect.ValueOf(RefreshToken).String(), c)
	SaveSession(user.ID, AccessToken, c)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"access_token": token.AccessToken, "refresh_token": token.RefreshToken})
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie()
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

		var isAdmin = false
		fmt.Println(claims["is_admin"])
		if claims["is_admin"] != nil {
			isAdmin, ok = claims["is_admin"].(bool) //convert the interface to string
			if !ok {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
			}
		}

		userId := uuid.MustParse(claims["user_id"].(string))
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
		ts, createErr := CreateJWTToken(userId, isAdmin)
		if createErr != nil {
			return c.Status(fiber.StatusCreated).JSON(createErr.Error())
		}
		//save the tokens metadata to redis
		saveErr := CreateAuthRecord(userId, ts)
		if saveErr != nil {
			return c.Status(fiber.StatusForbidden).JSON(err.Error())
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}

		SaveUserCookies(ts.AccessToken, ts.RefreshToken, c)

		return c.Status(fiber.StatusCreated).JSON(tokens)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(errors.New("refresh expired"))
	}
}

/**
Search Users
*/

func FetchUserByUID(c *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(c.Params("uid"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	user, err := GetUserByUID(parsedUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func FetchUserByID(c *fiber.Ctx) error {
	user, err := GetUserByID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func FetchUserByUsername(c *fiber.Ctx) error {
	user, err := GetUserByUsername(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

/*
	Create and Modify User
*/

func RegisterUser(c *fiber.Ctx) error {

	var body User

	body.ID = uuid.New()
	body.LastActive = time.Now()

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	err = auth.EncodeToArgon(&body.Password)
	if err != nil {
		return err
	}

	formErr := ValidateUserStruct(&body)
	if formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(formErr)
	}

	err = CreateUser(&body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	CreateUserLog("registered", body.ID, c)

	return c.Status(fiber.StatusCreated).JSON(body.ID.String())
}

func EditUserRequest(c *fiber.Ctx) error {

	var editUser User

	authUser, err := CheckAuthorisationIsValid(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}

	err = db.First(&editUser, c.Params("id")).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	if authUser.ID != editUser.ID && !authUser.Admin {
		return c.Status(fiber.StatusUnauthorized).JSON("unauthorised to edit another user")
	}

	err = c.BodyParser(&editUser)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}

	if c.FormValue("password") != "" {
		auth.EncodeToArgon(&editUser.Password)
	}

	if c.FormValue("uid") != "" {
		return errors.New("uid cannot be changed")
	}

	formErr := ValidateUserStruct(&editUser)
	if formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(formErr)
	}

	err = db.Save(editUser).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(editUser)
}

func DeleteUser(c *fiber.Ctx) error {
	uuid, err := uuid.Parse(c.Query("uid"))
	if err != nil {
		return err
	}

	err = DeleteUserByUID(uuid)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON("user deleted")
}

/**
User Photos
**/

func DeleteUserPhoto(c *fiber.Ctx, photoKey string) error {

	err := db.First(&user, c.Params("id")).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	if photoKey == "profile_photo" {

		if user.ProfilePhoto == nil {
			return c.Status(fiber.StatusNotFound).JSON("photo not found")
		}

		err = os.Remove(*user.ProfilePhoto)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		user.ProfilePhoto = nil
	} else if photoKey == "header_photo" {

		if user.HeaderPhoto == nil {
			return c.Status(fiber.StatusNotFound).JSON("photo not found")
		}

		err = os.Remove(*user.HeaderPhoto)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
		user.HeaderPhoto = nil
	}

	err = db.Save(user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON("photo deleted")
}

func UploadUserPhoto(c *fiber.Ctx, photoKey string) error {

	user, err := GetUserByID(c.Params("id"))
	if err != nil {
		return err
	}

	dir := "uploads/photos/user/" + *user.Username + "/"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}

	file, err := c.FormFile("photo")
	if err != nil {
		return err
	}

	contentType := file.Header.Get("content-type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		return errors.New("photo is not jpeg or png")
	}

	filename := uuid.New()

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
