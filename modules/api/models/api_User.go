package models

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/twinj/uuid"
	"github.com/yegamble/go-tube-api/modules/api/auth"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func FetchUserByUID(c *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(c.Params("uid"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	user, err := GetUserByUID(*parsedUUID)
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

	var body *User

	body.UID = uuid.NewV4()

	body.LastActive = time.Now()
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	auth.EncodeToArgon(&body.Password)

	formErr := ValidateUserStruct(body)
	if formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(formErr)
	}

	_, err = CreateUser(body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	CreateUserLog("registered", body.ID, c)

	return c.Status(fiber.StatusCreated).JSON(body.UID.String())
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

	CreateAuth(user.ID, token)

	AccessToken := reflect.ValueOf((*token).AccessToken).String()
	RefreshToken := reflect.ValueOf((*token).RefreshToken).String()

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + AccessToken
	c.Set("Authorization", bearer)

	err = SaveSession(user.ID, reflect.ValueOf(AccessToken).String(), reflect.ValueOf(RefreshToken).String(), c)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"access_token": token.AccessToken, "refresh_token": token.RefreshToken})
}

func EditUserRequest(c *fiber.Ctx) error {

	var editUser User

	authUser, err := Auth(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}

	err = db.First(&editUser, c.Params("id")).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	if authUser.ID != editUser.ID && !authUser.IsAdmin {
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

func DeleteUserPhoto(c *fiber.Ctx, photoKey string) error {

	err := db.First(&user, c.Params("id")).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	if photoKey == "profile_photo" {
		err = os.Remove(user.ProfilePhoto)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		user.ProfilePhoto = ""
	} else if photoKey == "header_photo" {
		err = os.Remove(user.HeaderPhoto)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
		user.HeaderPhoto = ""
	}

	err = db.Save(user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON("photo deleted")
}

func UploadUserPhoto(c *fiber.Ctx, photoKey string) error {

	id := c.Params("id")
	user, err := GetUserByID(id)
	if err != nil {
		return err
	}

	dir := "uploads/photos/user/" + user.Username + "/"

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

	filename := uuid.NewV4()

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