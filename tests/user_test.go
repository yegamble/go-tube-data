package tests

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/yegamble/go-tube-api/modules/api/handler"
	"github.com/yegamble/go-tube-api/modules/api/user"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserSignUp(t *testing.T) {

	assert.New(t)

	var u user.User
	var res []handler.ErrorResponse
}

// go test -run -v Test_Handler
func Test_Handler(t *testing.T) {
	app := New(Config{
		ErrorHandler: func(c *Ctx, err error) error {
			utils.AssertEqual(t, "1: USE error", err.Error())
			return DefaultErrorHandler(c, err)
		},
	})

	app.Post("/user", func(c *Ctx) {
		c.SendStatus(400)
	})

	resp, err := app.Test(httptest.NewRequest("POST", "/user", nil))

	utils.AssertEqual(t, nil, err, "app.Test")
	utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
}
