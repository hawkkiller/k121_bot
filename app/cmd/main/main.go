package main

import (
	"context"
	"github.com/hawkkiller/k121_bot/internal/config"
	"github.com/hawkkiller/k121_bot/internal/schedule"
	"github.com/hawkkiller/k121_bot/internal/schedule/db"
	"github.com/hawkkiller/k121_bot/pkg/client/postgresql"
	"github.com/hawkkiller/k121_bot/pkg/logging"
	"gopkg.in/telebot.v3"
	"time"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()

	logger.Infof("Successfully initialized logger")

	logger.Infof("Get config")
	cfg := config.GetConfig()

	logger.Infof("Connect to postgres client")
	client, err := postgresql.NewClient(context.Background(), cfg.Storage)

	if err != nil {
		logger.Println("failed to connect to postgresql")
		logger.Error(err)
	}
	logger.Infof("Successfully connected to postgresql")

	logger.Infof("Create storage")
	storage := db.NewStorage(client, logger)
	logger.Infof("Storage created")

	logger.Infof("Create service")
	service, err := schedule.NewService(storage, logger)
	if err != nil {
		logger.Println("service creating failed")
		logger.Error(err)
	}
	logger.Infof("Service created")

	logger.Infof("Create bot")
	pref := telebot.Settings{
		Token:   cfg.Telegram.Token,
		Poller:  &telebot.LongPoller{Timeout: 60 * time.Second},
		Verbose: false,
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		logger.Println("bot creating failed")
		logger.Error(err)
	}
	logger.Infof("Bot created")

	logger.Infof("Create handler")
	scheduleHandler := &schedule.Handler{
		Logger:  logger,
		Bot:     bot,
		Service: service,
	}
	logger.Infof("Handler created")

	logger.Infof("Register handlers")
	scheduleHandler.Register()
	logger.Infof("Handlers registered")

	bot.Start()
}
