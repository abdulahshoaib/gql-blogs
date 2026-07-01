package model

import "time"

type Post struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	Body      string    `gorm:"not null"`
	AuthorID  uint      `gorm:"not null;index"`
	Author    User      `gorm:"foreignKey:AuthorID"`
	Comments  []Comment `gorm:"foreignKey:PostID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
