package states

import (
	"fin-manager-bot/internal/models"
	"sync"
)

var (
	userStates = make(map[int64]*models.UserState)
	statesMu   sync.Mutex
)

// ResetUserState удаляет состояние пользователя.
func ResetUserState(chatID int64) {
	statesMu.Lock()
	defer statesMu.Unlock()
	delete(userStates, chatID)
}

func GetUserState(chatID int64) *models.UserState {
	statesMu.Lock()
	defer statesMu.Unlock()
	if state, ok := userStates[chatID]; ok {
		return state
	}
	return nil
}

// SetUserState устанавливает состояние для пользователя с указанным действием.
func SetUserState(chatID int64, action string) {
	statesMu.Lock()
	defer statesMu.Unlock()
	userStates[chatID] = &models.UserState{Action: action}
}
func UpdateUserState(chatID int64, stateData *models.UserState) {
	statesMu.Lock()
	defer statesMu.Unlock()
	userStates[chatID] = stateData
}
