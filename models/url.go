package models

import (
	"time"

	"gorm.io/gorm"
)

type URL struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	LongURL   string         `gorm:"not null" json:"long_url"`
	ShortCode string         `gorm:"uniqueIndex;not null" json:"short_code"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	ViewCount uint           `gorm:"default:0" json:"view_count"`
	Name  string        	 `gorm:"index" json:"Name"`
	Campaign  string         `gorm:"index" json:"Campaign"`
	Medium  string           `gorm:"index" json:"Medium"`
	Source  string           `gorm:"index" json:"Source"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}