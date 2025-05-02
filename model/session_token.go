package model

import (
	"time"

	"gorm.io/gorm"
)

type UserSession struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64    `gorm:"column:user_id;not null;" json:"user_id"`
	SessionToken  string    `gorm:"column:session_token" json:"session_token"`
	SessionExpiry int64     `gorm:"column:session_expiry" json:"session_expiry"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}

func MigrateUserSession(db *gorm.DB) error {
	return db.AutoMigrate(&UserSession{})
}
