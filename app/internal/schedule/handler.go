package schedule

import (
	"encoding/json"
	"github.com/hawkkiller/k121_bot/pkg/logging"
	"gopkg.in/telebot.v3"
)

type Handler struct {
	Logger  logging.Logger
	Bot     *telebot.Bot
	Service Service
}

func (h *Handler) Register() {
	h.Bot.Handle(telebot.OnText, h.GetSchedule)
	h.Bot.Handle(telebot.OnVoice, h.AnswerAudio)
	h.Bot.Handle(telebot.OnDocument, h.UploadSchedule)
}

func (h *Handler) GetSchedule(ctx telebot.Context) error {

	return nil
}

func (h *Handler) UploadSchedule(ctx telebot.Context) error {
	model := Schedule{}
	doc := ctx.Message().Document
	file, err := ctx.Bot().File(&doc.File)
	if err != nil {
		h.Logger.Error(err)
		return err
	}
	err = json.NewDecoder(file).Decode(&model)
	if err != nil {
		return err
	}
	model.ChatId = ctx.Message().Chat.ID

	err = h.Service.CreateSchedule(model)
	if err != nil {
		return err
	}
	err = ctx.Reply("Запись успешно добавлена")
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) AnswerAudio(ctx telebot.Context) error {
	err := ctx.Reply("Удалил при мне.")
	if err != nil {
		return err
	}
	return nil
}
