package auth

import "fin-manager-bot/internal/config"

// IsAdmin возвращает true, если userID совпадает с adminUserID из переменных окружения.
func IsAdmin(userID int, config *config.Config) bool {
	return userID == config.AdminUserId
}
