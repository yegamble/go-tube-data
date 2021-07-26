package models

import (
	"time"
)

type Views struct {
	ID        int64
	User      User
	Video     Video
	CreatedAt time.Time
}
