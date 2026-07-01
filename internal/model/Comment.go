package model

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	Body      string    `gorm:"not null"`
	AuthorID  uint      `gorm:"not null;index"`
	Author    User      `gorm:"foreignKey:AuthorID"`
	PostID    uint      `gorm:"not null;index"`
	Post      Post      `gorm:"foreignKey:PostID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
