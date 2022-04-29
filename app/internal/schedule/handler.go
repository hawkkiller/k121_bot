package schedule

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hawkkiller/k121_bot/pkg/logging"
	"gopkg.in/telebot.v3"
	"strings"
	"time"
)

type Handler struct {
	Logger  logging.Logger
	Bot     *telebot.Bot
	Service Service
}

func (h *Handler) Register() {
	h.Bot.Handle("хелп", h.Help)
	h.Bot.Handle(telebot.OnText, h.GetSchedule)
	h.Bot.Handle(telebot.OnEdited, h.GetSchedule)
	h.Bot.Handle(telebot.OnVoice, h.AnswerAudio)
	h.Bot.Handle(telebot.OnDocument, h.UploadSchedule)
}

func (h *Handler) GetSchedule(ctx telebot.Context) error {
	text := ctx.Message().Text
	// if message starts with расписание
	if !strings.HasPrefix(strings.ToLower(text), "расписание") {
		return nil
	}

	h.Logger.Info("Get schedule")
	schedule, err := h.Service.GetSchedule(context.Background(), ctx.Message().Chat.ID)
	if err != nil {
		h.Logger.Error(err)
		return err
	}

	segments := strings.Split(text, " ")

	day := new(Day)

	today := int(time.Now().Weekday()) - 1

	if len(schedule.Days) < today || today < 0 {
		today = 0
	}

	if len(schedule.Days) >= today {
		day = &schedule.Days[today]
	}

	if len(segments) > 1 {
		for _, d := range schedule.Days {
			if strings.ToLower(d.Caption) == strings.ToLower(segments[1]) {
				day = &d
				break
			}
		}
	}
	str := fmt.Sprintf("`%s`\n", day.Caption)
	for _, p := range day.Pairs {
		str += fmt.Sprintf("\n**%s**\n", p.Title)
		codes := strings.Split(strings.ReplaceAll(p.AdditionalInfo, "\n", " "), " ")
		for _, code := range codes {
			c := strings.Split(code, "-")
			if len(c) == 3 {
				str += fmt.Sprintf("**%s** `%s`:`%s`\n", c[0], c[1], c[2])
			} else {
				str += fmt.Sprintf("**%s**\n", code)
			}
		}
	}
	err = ctx.Reply(str)
	if err != nil {
		return err
	}
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
	model.ChatId = int(ctx.Message().Chat.ID)

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

func (h *Handler) Help(ctx telebot.Context) error {
	err := ctx.Reply("" +
		"Здесь вы можете создать [расписание](https://site.michaeldeveloper.com)\n" +
		"Формат конфы: Учитель-код-пароль\n" +
		"Чтобы получить расписание напишите расписание или выберите день недели\n" +
		"Например: `расписание понедельник` или `расписание`\n" +
		"Исходный код - [github](https://github.com/hawkkiller/k121_bot/)\n" +
		"")
	if err != nil {
		return err
	}
	return nil
}
