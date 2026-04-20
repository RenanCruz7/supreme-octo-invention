package models

type Post struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body" binding:"required"`
	UserID uint   `json:"user_id" binding:"required"` // Chave estrangeira
	User   *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
