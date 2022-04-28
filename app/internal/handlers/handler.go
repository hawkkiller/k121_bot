package handlers

import "gopkg.in/telebot.v3"

type Handler interface {
	Register(bot *telebot.Bot)
}
