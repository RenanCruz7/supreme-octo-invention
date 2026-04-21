package models

import "time"

type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" binding:"required,min=3"`
	Body      string    `json:"body" binding:"required,min=10"`
	UserID    uint      `json:"user_id" binding:"required"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
