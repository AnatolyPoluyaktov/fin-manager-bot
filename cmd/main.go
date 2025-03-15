package main

import (
	"fin-manager-bot/internal/api"
	"fin-manager-bot/internal/config"
	"fin-manager-bot/internal/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()
	log := setupLogger(cfg.Env)
	log.Info("Starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("Debug mode is enabled")
	client := api.NewClient(cfg.FinManagerAPI.BaseUrl, cfg.FinManagerAPI.AuthToken)
	tgBot, err := tgbotapi.NewBotAPI(cfg.BotAuthToken)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("Authorized on account", slog.String("bot_name", tgBot.Self.UserName))
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := tgBot.GetUpdatesChan(u)
	if err != nil {
		log.Error(err.Error())
	}
	for update := range updates {
		if update.Message != nil {
			handlers.HandleMessage(tgBot, update.Message, cfg, client)
		}

		if update.CallbackQuery != nil {
			handlers.HandleCallbackQuery(tgBot, update.CallbackQuery, cfg, client)
		}
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
