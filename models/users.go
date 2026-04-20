package models

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email" gorm:"uniqueIndex"`
	Posts []Post `json:"posts,omitempty" gorm:"foreignKey:UserID"` // Relacionamento 1:N
}
