package user

import (
	"github.com/yegamble/go-tube-api/modules/api/video"
	"time"
)

type Views struct {
	ID        int64
	User      User
	Video     video.Video
	CreatedAt time.Time
}
