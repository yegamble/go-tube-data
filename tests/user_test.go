package tests

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/yegamble/go-tube-api/modules/api/auth"
	"github.com/yegamble/go-tube-api/modules/api/models"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

var users []*models.User

func init() {
	models.SyncModels()
}

func SeedUsers() error {
	for i := 0; i < 10; i++ {
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
			return err
		}

		err = user.Create(gofakeit.IPv4Address())
		if err != nil {
			return err
		}
		users = append(users, &user)
	}
	return nil
}

func seedTags() []*models.Tag {
	var userTags []*models.Tag
	for i := 0; i < 10; i++ {
		word := gofakeit.Word()
		tag := models.Tag{
			Value: &word,
		}
		userTags = append(userTags, &tag)
	}

	return userTags
}

func TestCreateUsers(t *testing.T) {
	err := SeedUsers()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
		return
	}

	assert.Equal(t, 10, len(users), "count of users successfully created in the application")

	t.Log("Deleting Test Users")
	DeleteTestUsers(t)
}

func DeleteTestUsers(t *testing.T) {
	for _, user := range users {
		err := user.Delete()
		if err != nil {
			t.Fatal(err.Error())
		}

		u, err := models.GetUserByUUID(user.UUID)
		if assert.Error(t, err) {
			assert.Equal(t, err.Error(), "record not found")
		}

		assert.Empty(t, u, "User ", user.Username, " successfully deleted")
	}
}

func TestComparePasswordAndHash(t *testing.T) {
	err := SeedUsers()
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, u := range users {
		match, err := auth.ComparePasswordAndHash(&Password, u.Password)
		if err != nil {
			t.Fatal(err.Error())
		} else if !match {
			t.Fatal("Password does not match")
		}

		assert.Equal(t, match, true)
	}

	t.Log("Deleting Test Users")
	DeleteTestUsers(t)

}

func TestUserCreateTags(t *testing.T) {
	err := SeedUsers()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
		return
	}

	for _, user := range users {
		userTags := seedTags()
		user.CreateTags(userTags)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
			return
		}
	}

	t.Log("Deleting Test Users")
	DeleteTestUsers(t)
}

func TestUserLogin(t *testing.T) {
	app := fiber.New()

	// http.Request
	req := httptest.NewRequest("GET", "http://google.com", nil)
	req.Header.Set("X-Custom-Header", "hi")

	// http.Response
	resp, _ := app.Test(req)

	// Do something with results:
	if resp.StatusCode == fiber.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body)) // => Hello, World!
	}
}

//func TestUserSubscriptions(t *testing.T) {
//	err := SeedUsers()
//	if err != nil {
//		t.Log(err.Error())
//		t.Fail()
//		return
//	}
//
//	for i, user := range users {
//		if i > 0{
//			err = user.SubscribeToChannel(users[0].UUID)
//			if err != nil {
//				t.Log(err.Error())
//				t.Fail()
//				return
//			}
//
//			fmt.Println(*user.Subscriptions[0])
//		}
//	}
//
//	t.Log("Deleting Test Users")
//	//DeleteTestUsers(t)
//}

//func TestUploadProfilePicture(t *testing.T) {
//		app := fiber.New(fiber.Config{
//			ErrorHandler: func(c *fiber.Ctx, err error) error {
//				utils.AssertEqual(t, "1: USE error", err.Error())
//				t.Fatal(err.Error())
//				return nil
//			},
//		})
//
//		app.Post("/user", func(c *fiber.Ctx) error {
//			err := c.SendStatus(400)
//			if err != nil {
//				t.Fatal(err.Error())
//			}
//			return nil
//		})
//
//		resp, err := app.Test(httptest.NewRequest("POST", "/user", nil))
//
//		utils.AssertEqual(t, nil, err, "app.Test")
//		utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
//}

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
