package models

type Session struct {
	ID          uint64
	SessionID   string `json:"session_id"`
	UserID      int64
	User        User   `json:"user_id" form:"user_id" gorm:"foreignKey:UserID;references:ID"`
	Fingerprint string `json:"fingerprint"`
}
