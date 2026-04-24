package repositories

import (
	"time"

	"awesomeProject/db"
	"awesomeProject/models"
)

func CreateRefreshToken(token models.RefreshToken) error {
	return db.DB.Create(&token).Error
}

// GetRefreshToken busca um refresh token válido (não revogado e não expirado)
func GetRefreshToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	result := db.DB.Where("token = ? AND revoked = false AND expires_at > ?", token, time.Now()).First(&rt)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rt, nil
}

// RevokeRefreshToken revoga um token específico (rotação)
func RevokeRefreshToken(token string) error {
	return db.DB.Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}

// RevokeAllUserTokens revoga todos os tokens ativos de um usuário (logout)
func RevokeAllUserTokens(userID uint) error {
	return db.DB.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}

// DeleteExpiredTokens remove tokens expirados ou revogados do banco
func DeleteExpiredTokens() error {
	return db.DB.
		Where("expires_at < ? OR revoked = true", time.Now()).
		Delete(&models.RefreshToken{}).Error
}
