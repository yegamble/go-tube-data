package tests

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"github.com/yegamble/go-tube-api/modules/api/models"
	"testing"
	"time"
)

var (
	Username    = "test"
	FirstName   = "Thomas"
	LastName    = "Lok"
	Email       = gofakeit.Email()
	DateOfBirth = time.Date(2001, time.November, 10, 23, 0, 0, 0, time.UTC)
	Password    = "xvZDr5AR/-EM"
)

var u models.User

func TestUserValidation(t *testing.T) {

	assert.New(t)

	var res []handler.ErrorResponse

	u.Username = &Username
	u.FirstName = FirstName
	u.LastName = LastName
	u.Email = &Email
	u.DateOfBirth = &DateOfBirth
	u.Password = Password
	result := models.ValidateUserStruct(&u)
	if result != nil {
		for k := range result {
			res = append(res, *result[k])
		}
		assert.Fail(t, "Check if User Validation is Correct", res)
		t.Fail()
	}
	assert.Nil(t, result, "Failed Validation Results Are Empty")
}

func TestFirstNameFieldMissingValidation(t *testing.T) {
	assert.New(t)
	var u models.User
	u.FirstName = FirstName
	result := models.ValidateUserStruct(&u)
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}

func TestLastNameFieldMissingValidation(t *testing.T) {
	assert.New(t)
	u.LastName = ""
	result := models.ValidateUserStruct(&u)
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}

func TestDateOfBirthFieldMissingValidation(t *testing.T) {
	assert.New(t)
	u.DateOfBirth = nil
	result := models.ValidateUserStruct(&u)
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}

func TestPasswordFieldMissingValidation(t *testing.T) {
	assert.New(t)
	u.Password = ""
	result := models.ValidateUserStruct(&u)
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}

func TestEmailFieldMissingValidation(t *testing.T) {
	assert.New(t)
	print(u.Password)
	u.Email = nil
	result := models.ValidateUserStruct(&u)
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}
