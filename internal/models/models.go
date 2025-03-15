package models

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RawCategory struct {
	Name string `json:"name"`
}

type RawExpense struct {
	CategoryID int    `json:"category_id"`
	Currency   string `json:"currency"`
	Amount     int    `json:"amount"`
	Note       string `json:"note"`
	ActionDate string `json:"action_date"`
}

// UserState хранит текущее состояние диалога с пользователем.
type UserState struct {
	// Возможные значения Action:
	// "create_category",
	// "create_expense_amount",
	// "create_expense_description",
	// "create_expense_select_category"
	Action      string
	ExpenseData RawExpense
}
