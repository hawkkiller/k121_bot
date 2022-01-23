package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessage(bot *tgbotapi.BotAPI, id int64, text string) {
	msg := tgbotapi.NewMessage(id, text)
	msg.ParseMode = "markdown"
	msg.DisableWebPagePreview = true
	msg.DisableNotification = true
	if _, err := bot.Send(msg); err != nil {
		fmt.Println(err)
	}
}
