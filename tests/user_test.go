package tests

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/yegamble/go-tube-api/modules/api/auth"
	"github.com/yegamble/go-tube-api/modules/api/models"
	"testing"
	"time"
)

var users []models.User

func TestUserSignUp(t *testing.T) {
	t1 := time.Now()
	for i := 0; i < 10000; i++ {
		username := gofakeit.Username()
		email := gofakeit.Email()
		dob := gofakeit.Date()
		user := models.User{
			FirstName:   gofakeit.FirstName(),
			LastName:    gofakeit.LastName(),
			Email:       &email,
			DateOfBirth: &dob,
			Password:    Password,
			Username:    &username,
		}

		err := auth.EncodeToArgon(&user.Password)
		if err != nil {
			t.Fatal(err.Error())
		}

		models.CreateUser(&user)

		users = append(users, user)
	}

	t2 := t1.Add(time.Second * 341)
	diff := t2.Sub(t1)
	fmt.Println(diff)
}

func TestUploadProfilePicture(t *testing.T) {
	//app := fiber.New()
	//app.Test()
	//
	//http.NewRequest("Post", "localhost:3000")
	////for k, v := range users {
	////	body, err := app.Test("GET /demo HTTP/1.1\r\nHost: google.com\r\n\r\n")
	////	models.UploadUserPhoto()
	////}

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
