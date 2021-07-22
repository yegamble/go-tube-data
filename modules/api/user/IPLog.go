package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/database"
	"time"
)

type IPLog struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key"`
	UserID    int64     `json:"user_id" form:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BannedLog struct {
	IPAddress IPLog
}

func insertUserIPLog(u *User, ctx *fiber.Ctx) (uuid.UUID, error) {
	db := database.DBConn

	var log IPLog
	log.User.ID = u.ID
	log.IPAddress = ctx.IP()

	result := db.Create(&log)

	return log.ID, result.Error
}

func clearIPLogs(clearAll bool) error {
	db := database.DBConn
	result := db.Where("created_at < NOW() - INTERVAL 1 WEEK").Delete(IPLog{})
	return result.Error
}
