package tests

import (
	"github.com/yegamble/go-tube-api/modules/api/user"
	"testing"
)

func TestUserSignUp(t *testing.T) {

	var u user.User
	u.Username = Username
	u.FirstName = FirstName
	u.LastName = LastName
	u.Email = Email
	u.DateOfBirth = DateOfBirth
	u.Password = Password

	user.CreateUser(u)
}

//// go test -run -v Test_Handler
//func Test_Handler(t *testing.T) {
//	app := New(Config{
//		ErrorHandler: func(c *Ctx, err error) error {
//			utils.AssertEqual(t, "1: USE error", err.Error())
//			return DefaultErrorHandler(c, err)
//		},
//	})
//
//	app.Post("/user", func(c *Ctx) {
//		c.SendStatus(400)
//	})
//
//	resp, err := app.Test(httptest.NewRequest("POST", "/user", nil))
//
//	utils.AssertEqual(t, nil, err, "app.Test")
//	utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
//}
