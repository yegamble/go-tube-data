package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/database"
	"time"
)

type IpLog struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key"`
	User      User      `json:"user_id"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func insertUserIPLog(u *User, ctx *fiber.Ctx) (uuid.UUID, error) {
	db := database.DBConn

	var log IpLog
	log.User.ID = u.ID
	log.IPAddress = ctx.IP()

	result := db.Create(&log)

	return log.ID, result.Error
}

func clearIPLogs(clearAll bool) error {
	db := database.DBConn
	result := db.Where("created_at < NOW() - INTERVAL 1 WEEK").Delete(IpLog{})
	return result.Error
}
