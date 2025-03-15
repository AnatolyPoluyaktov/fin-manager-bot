package api

import (
	"bytes"
	"encoding/json"
	"fin-manager-bot/internal/models"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(baseURL string, token string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
		Token:      token,
	}
}
func (c *Client) CreateCategory(category *models.Category) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/categories", c.BaseURL)
	jsonData, err := json.Marshal(category)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil

}

func (c *Client) GetCategories() ([]models.Category, error) {
	url := fmt.Sprintf("%s/api/v1/categories", c.BaseURL)

	// Создаем новый запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Устанавливаем заголовки (например, указываем, что ожидаем JSON)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	// Создаем HTTP-клиент и выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var categories []models.Category
	err = json.NewDecoder(resp.Body).Decode(&categories)
	if err != nil {
		return nil, err
	}
	return categories, nil

}

func (c *Client) CreateExpense(expense *models.RawExpense) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/expenses", c.BaseURL)
	jsonData, err := json.Marshal(expense)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil

}
