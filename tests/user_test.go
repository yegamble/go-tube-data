package tests

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/yegamble/go-tube-api/modules/api/auth"
	"github.com/yegamble/go-tube-api/modules/api/models"
	"gorm.io/gorm"
	"testing"
)

var users []*models.User

func init() {
	err := models.SyncModels()
	if err != nil {
		return
	}
}

func SeedUsers() error {
	users = nil
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

		err = user.New(gofakeit.IPv4Address())
		if err != nil {
			return err
		}
		users = append(users, &user)
	}
	return nil
}

func SeedMillionUsers(t *testing.T) {
	for i := 0; i < 100; i++ {
		users = nil
		size := 100000
		fmt.Println("Running for loop…")
		for i := 0; i < size; i++ {
			username := gofakeit.Date().String() + gofakeit.Username()
			email := gofakeit.Date().String() + gofakeit.Email()
			dob := gofakeit.Date()
			user := &models.User{
				FirstName:   gofakeit.FirstName(),
				LastName:    gofakeit.LastName(),
				Email:       &email,
				DateOfBirth: &dob,
				Password:    Password,
				Username:    &username,
			}
			users = append(users, user)
		}

		err := models.CreateUsers(users)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
			return
		}
	}

	return
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

	for _, user := range users {
		user.GetLogs()
		assert.Equal(t, 1, len(user.Logs), "user log is created")
		assert.Equal(t, "registered", *user.Logs[0].Activity, "user log activity is registered")
		assert.NotEmpty(t, user.Logs[0].IPAddress, "user log ip address is not empty")
	}

	t.Log("Deleting Test Users")
	DeleteTestUsers(t)
}

func DeleteTestUsers(t *testing.T) {
	for _, user := range users {
		err := user.Delete()
		user.GetLogs()
		assert.Equal(t, 1, len(user.Logs), "user log still exists")
		if err != nil {
			t.Fatal(err.Error())
		}

		u, err := models.GetUserByUUID(user.ID)
		if assert.Error(t, err) {
			assert.Equal(t, err, gorm.ErrRecordNotFound)
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

//func TestUserCreateTags(t *testing.T) {
//	err := SeedUsers()
//	if err != nil {
//		t.Log(err.Error())
//		t.Fail()
//		return
//	}
//
//	for _, user := range users {
//		userTags := seedTags()
//		err = user.CreateTags(userTags)
//		if err != nil {
//			t.Log(err.Error())
//			t.Fail()
//			return
//		}
//		if err != nil {
//			t.Log(err.Error())
//			t.Fail()
//			return
//		}
//	}
//
//	t.Log("Deleting Test Users")
//	DeleteTestUsers(t)
//}

func TestUserSubscriptions(t *testing.T) {
	err := SeedUsers()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
		return
	}

	for i, user := range users {
		if i > 0 {
			fmt.Println(users[0].ID)
			err = user.SubscribeToChannel(users[0].ID)
			if err != nil {
				t.Log(err.Error())
				DeleteTestUsers(t)
				t.Fail()
				return
			}

			assert.Equal(t, 1, len(user.Subscriptions), "user has one subscription")
		}
	}

	t.Log("Deleting Test Users")
	DeleteTestUsers(t)

}

func TestUserDeleteSubscriptions(t *testing.T) {
	err := SeedUsers()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
		return
	}

	for i, user := range users {
		if i > 0 {
			err = user.SubscribeToChannel(users[0].ID)
			if err != nil {
				t.Log(err.Error())
				DeleteTestUsers(t)
				t.Fail()
				return
			}

			assert.Equal(t, 1, len(user.Subscriptions), "user has one channel subscription")

			err = user.UnsubscribeFromChannel(users[0].ID)
			if err != nil {
				t.Log(err.Error())
				DeleteTestUsers(t)
				t.Fail()
				return
			}

			assert.Equal(t, 0, len(user.Subscriptions), "user is not subscribed to a channel")

			err = user.GetSubscriptions()
			if err != nil {
				t.Log(err.Error())
				DeleteTestUsers(t)
				t.Fail()
				return
			}

			fmt.Println(user.Subscriptions)

		}
	}

	t.Log("Deleting Test Users")
	DeleteTestUsers(t)
}

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
