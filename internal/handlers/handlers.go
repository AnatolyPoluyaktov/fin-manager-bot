package handlers

import (
	"fin-manager-bot/internal/auth"
	"fin-manager-bot/internal/config"
	"fin-manager-bot/internal/models"
	"fin-manager-bot/internal/states"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FinManagerClient interface {
	CreateCategory(category *models.Category) (*http.Response, error)
	GetCategories() ([]models.Category, error)
	CreateExpense(expense *models.RawExpense) (*http.Response, error)
}

func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, config *config.Config, client FinManagerClient) {
	chatID := message.Chat.ID
	if message.From == nil || !auth.IsAdmin(message.From.ID, config) {
		SendMessage(bot, chatID, "У вас нет прав доступа!")
		return
	}
	text := message.Text

	switch text {
	case "/start":
		reply := "Привет! Доступные команды:\n" +
			"/create_category - создать категорию\n" +
			"/create_expense - создать трату\n" +
			"/list_categories - список категорий"
		SendMessage(bot, chatID, reply)
		states.ResetUserState(chatID)
	case "/create_category":
		SendMessage(bot, chatID, "Введите название категории:")
		states.SetUserState(chatID, "create_category")
	case "/create_expense":
		SendMessage(bot, chatID, "Введите сумму траты:")
		states.SetUserState(chatID, "create_expense_amount")
	case "/list_categories":
		GetCategoriesHandler(bot, chatID, client)
	default:
		userState := states.GetUserState(chatID)
		if userState == nil {
			SendMessage(bot, chatID, "Неизвестная команда. Используйте /start для списка команд.")
			return
		}
		switch userState.Action {
		case "create_category":
			createCategoryHandler(bot, chatID, client, text)

		case "create_expense_amount":
			CreateExpenseAmountHandler(bot, chatID, client, text)

		case "create_expense_action_date":
			CreateExpenseActionDateHandler(bot, chatID, text)
		case "create_expense_description":
			CreateExpenseDescriptionHandler(bot, chatID, client, text)
		default:

			SendMessage(bot, chatID, "Неверная команда или состояние.")
		}
	}
}

func createCategoryHandler(bot *tgbotapi.BotAPI, chatID int64, client FinManagerClient, text string) {
	if text == "" {
		SendMessage(bot, chatID, "Название категории не может быть пустым")
	}
	category := &models.Category{Name: text}
	resp, err := client.CreateCategory(category)
	if err != nil {
		states.ResetUserState(chatID)
		SendMessage(bot, chatID, "Ошибка при создании категории "+err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		SendMessage(bot, chatID, "Ошибка при создании категории "+resp.Status)
		states.ResetUserState(chatID)
		return
	}

	SendMessage(bot, chatID, "Категория успешно создана")
	states.ResetUserState(chatID)
}

func CreateExpenseAmountHandler(bot *tgbotapi.BotAPI, chatID int64, client FinManagerClient, text string) {
	if text == "" {
		SendMessage(bot, chatID, "Сумма траты не может быть пустой")
		return
	}

	amount, err := strconv.Atoi(text)
	if err != nil {
		SendMessage(bot, chatID, "Сумма траты должна быть числом")
		states.ResetUserState(chatID)
		return
	}
	if amount <= 0 {
		SendMessage(bot, chatID, "Сумма траты должна быть больше нуля")
		states.ResetUserState(chatID)
		return
	}
	userState := states.GetUserState(chatID)
	userState.ExpenseData.Amount = amount
	userState.ExpenseData.Currency = "RUB"
	userState.Action = "create_expense_select_category"
	states.UpdateUserState(chatID, userState)
	categories, err := client.GetCategories()
	if err != nil {
		SendMessage(bot, chatID, "Ошибка при получении категорий "+err.Error())
		states.ResetUserState(chatID)
		return
	}
	if len(categories) == 0 {
		SendMessage(bot, chatID, "Нет доступных категорий")
		states.ResetUserState(chatID)
		return
	}
	var buttons [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	for _, c := range categories {
		btn := tgbotapi.NewInlineKeyboardButtonData(c.Name, strconv.Itoa(c.ID))
		row = append(row, btn)
		if len(row) >= 2 {
			buttons = append(buttons, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(row) > 0 {
		buttons = append(buttons, row)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	replyMsg := tgbotapi.NewMessage(chatID, "Выберите категорию для траты:")
	replyMsg.ReplyMarkup = keyboard
	_, err = bot.Send(replyMsg)
	if err != nil {
		return
	}
}

func CreateExpenseActionDateHandler(bot *tgbotapi.BotAPI, chatID int64, text string) {
	userState := states.GetUserState(chatID)
	date := text
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		SendMessage(bot, chatID, "Неверный формат даты. Используйте YYYY-MM-DD")
		states.ResetUserState(chatID)
		return
	}
	userState.ExpenseData.ActionDate = date
	userState.Action = "create_expense_description"
	states.UpdateUserState(chatID, userState)
	SendMessage(bot, chatID, "Введите описание траты")

}

func HandleCallbackQuery(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery, config *config.Config, client FinManagerClient) {
	chatID := cq.Message.Chat.ID

	// Проверяем авторизацию
	if cq.From == nil || !auth.IsAdmin(cq.From.ID, config) {
		SendMessage(bot, chatID, "У вас нет прав доступа!")
		return
	}

	userState := states.GetUserState(chatID)
	if userState == nil || userState.Action != "create_expense_select_category" {
		return
	}
	catID, err := strconv.Atoi(cq.Data)
	if err != nil {
		SendMessage(bot, chatID, "Неверные данные категории.")
		return
	}
	userState.ExpenseData.CategoryID = catID
	userState.Action = "create_expense_action_date"
	states.UpdateUserState(chatID, userState)
	SendMessage(bot, chatID, "Введите дату траты в формате YYYY-MM-DD")

}
func CreateExpenseDescriptionHandler(bot *tgbotapi.BotAPI, chatID int64, client FinManagerClient, text string) {
	userState := states.GetUserState(chatID)
	userState.ExpenseData.Note = text
	resp, err := client.CreateExpense(&userState.ExpenseData)
	if err != nil {
		SendMessage(bot, chatID, "Ошибка при создании траты "+err.Error())
		states.ResetUserState(chatID)
		return
	}
	if resp.StatusCode != http.StatusOK {
		SendMessage(bot, chatID, "Ошибка при создании траты "+resp.Status)
		states.ResetUserState(chatID)
		return
	}
	SendMessage(bot, chatID, "Трата успешно создана")
	states.ResetUserState(chatID)
}

func GetCategoriesHandler(bot *tgbotapi.BotAPI, chatID int64, client FinManagerClient) {
	categories, err := client.GetCategories()
	if err != nil {
		SendMessage(bot, chatID, "Ошибка получения категорий: "+err.Error())
		states.ResetUserState(chatID)
		return
	}
	if len(categories) == 0 {
		SendMessage(bot, chatID, "Категорий не найдено.")
		states.ResetUserState(chatID)
		return
	}
	reply := fmt.Sprintf("Список категорий (%d):\n", len(categories))
	for _, c := range categories {
		reply += strconv.Itoa(c.ID) + " - " + c.Name + "\n"
	}
	SendMessage(bot, chatID, reply)
	states.ResetUserState(chatID)

}
func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}
