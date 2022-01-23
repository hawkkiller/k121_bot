package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleHulp(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.FromChat().ID, "")
	msg.Text = "Я понимаю команды: \n /hulp \n /dog \n gh: https://github.com/hawkkiller/k121_bot"
	if _, err := bot.Send(msg); err != nil {
		fmt.Println(err)
		return
	}
}
