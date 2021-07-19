package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/yegamble/go-tube-api/modules/api/errorhandler"
	"github.com/yegamble/go-tube-api/modules/api/user"
	"testing"
	"time"
)

var (
	Username    = "test"
	FirstName   = "Thomas"
	LastName    = "Lok"
	Email       = "tomLok@tube.com"
	DateOfBirth = time.Date(2001, time.November, 10, 23, 0, 0, 0, time.UTC)
	Password    = "xvZDr5AR/-EM"
)

func TestUserValidation(t *testing.T) {
	assert.New(t)

	var u user.User
	var res []errorhandler.ErrorResponse

	u.Username = Username
	u.FirstName = FirstName
	u.LastName = LastName
	u.Email = Email
	u.DateOfBirth = DateOfBirth
	u.Password = Password
	result := user.ValidateStruct(&u)
	if result != nil {
		for k := range result {
			res = append(res, *result[k])
		}
		assert.Fail(t, "Check if User Validation is Correct", res)
		t.Fail()
	}
	assert.Nil(t, result, "Failed Validation Results Are Empty")
}

func TestFirstNameFieldMissing(t *testing.T) {
	assert.New(t)
	var u user.User
	u.FirstName = FirstName
	result := user.ValidateStruct(&u)
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}

func TestLastNameFieldMissing(t *testing.T) {
	assert.New(t)
	var u user.User
	u.LastName = LastName
	result := user.ValidateStruct(&u)
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}

func TestDateOfBirthFieldMissing(t *testing.T) {
	assert.New(t)
	var u user.User
	u.DateOfBirth = DateOfBirth
	result := user.ValidateStruct(&u)
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}

func TestPasswordFieldMissing(t *testing.T) {
	assert.New(t)
	var u user.User
	u.Password = Password
	result := user.ValidateStruct(&u)
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}

func TestEmailFieldMissing(t *testing.T) {
	assert.New(t)
	var u user.User
	u.Email = Email
	result := user.ValidateStruct(&u)
	assert.NotEmpty(t, result, "Failed Validation Results Are Empty")
}
