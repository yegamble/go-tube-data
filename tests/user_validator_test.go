package tests

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"github.com/yegamble/go-tube-api/modules/api/models"
	"testing"
)

var (
	Username    = "test"
	FirstName   = gofakeit.FirstName()
	LastName    = gofakeit.LastName()
	Email       = gofakeit.Email()
	DateOfBirth = gofakeit.Date()
	Password    = gofakeit.Password(true, false, false, false, false, 32)
	res         []handler.ErrorResponse
)

var u models.User

func TestUserValidation(t *testing.T) {

	assert.New(t)

	u.DisplayName = &Username
	u.Username = &Username
	u.FirstName = FirstName
	u.LastName = LastName
	u.Email = &Email
	u.DateOfBirth = &DateOfBirth
	u.Password = Password
	result := u.ValidateUserStruct()
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
	result := u.ValidateUserStruct()
	if assert.NotEmpty(t, result) {
		for k := range result {
			res = append(res, *result[k])
		}
		assert.NotEmpty(t, "Validation Field is Not Empty", res)
	}
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}

func TestLastNameFieldMissingValidation(t *testing.T) {
	assert.New(t)
	u.LastName = ""
	result := u.ValidateUserStruct()
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}

func TestDateOfBirthFieldMissingValidation(t *testing.T) {
	assert.New(t)
	u.DateOfBirth = nil
	result := u.ValidateUserStruct()
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}

func TestPasswordFieldMissingValidation(t *testing.T) {
	assert.New(t)
	u.Password = ""
	result := u.ValidateUserStruct()
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}

func TestEmailFieldMissingValidation(t *testing.T) {
	assert.New(t)
	print(u.Password)
	u.Email = nil
	result := u.ValidateUserStruct()
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}
