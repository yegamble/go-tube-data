package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

type IPLog struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	UserUUID  uuid.UUID `json:"user_id" form:"user_id" gorm:"type:varchar(255);"'`
	IPAddress string    `json:"ip_address" gorm:"type:text"`
	Activity  string    `json:"activity" gorm:"type:text"`
	CreatedAt time.Time
}

type BannedIPLog struct {
	IPAddress string
}

var (
	dsn string
)

func init() {
	godotenv.Load(".env")
	dsn = "root@tcp(127.0.0.1:3306)/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func CreateUserLog(activity string, userUID uuid.UUID, ctx *fiber.Ctx) (uint64, error) {

	var log IPLog
	log.UserUUID = userUID
	log.IPAddress = ctx.IP()
	log.Activity = activity

	result := db.Create(&log)

	return log.ID, result.Error
}

func ScheduleCleanup() error {
	log.Info("Create new cron")

	c := cron.New()
	c.AddFunc("*/1 * * * *", func() { ClearIPLogs(dsn) })

	// Start cron with one scheduled job
	log.Info("Start cron")
	c.Start()
	printCronEntries(c.Entries())
	time.Sleep(2 * time.Minute)

	// Funcs may also be added to a running Cron
	log.Info("Add new job to a running cron")
	entryID2, _ := c.AddFunc("*/2 * * * *", func() { log.Info("[Job 2]Every two minutes job\n") })
	printCronEntries(c.Entries())
	time.Sleep(5 * time.Minute)

	//Remove Job2 and add new Job2 that run every 1 minute
	log.Info("Remove Job2 and add new Job2 with schedule run every minute")
	c.Remove(entryID2)
	c.AddFunc("*/1 * * * *", func() { log.Info("[Job 2]Every one minute job\n") })
	time.Sleep(5 * time.Minute)
	return nil
}

func printCronEntries(cronEntries []cron.Entry) {
	log.Infof("Cron Info: %+v\n", cronEntries)
}

func ClearIPLogs(dsn string) error {

	var iplogs IPLog
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	res := db.Delete(&iplogs, "created_at < NOW() - INTERVAL 26 WEEK").Error
	log.Println("Logs Deleted: " + string(db.RowsAffected))
	return res
}
