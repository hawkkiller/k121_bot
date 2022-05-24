package schedule

import (
	"bytes"
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
	h.Bot.Handle("скачать расписание", h.DownloadSchedule)
	h.Bot.Handle("удалить расписание", h.DeleteSchedule)
	h.Bot.Handle(telebot.OnText, h.GetSchedule)
	h.Bot.Handle(telebot.OnEdited, h.GetSchedule)
	h.Bot.Handle(telebot.OnDocument, h.UploadSchedule)
}

func (h *Handler) GetSchedule(ctx telebot.Context) error {
	// defer recover func
	defer func() {
		if r := recover(); r != nil {
			h.Logger.Error(fmt.Sprintf("%v", r))
			err := ctx.Send(fmt.Sprintf("%v", r))
			if err != nil {
				return
			}
		}
	}()
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

	if schedule.ID == "" {
		h.Logger.Info("Schedule not found")
		err = ctx.Reply("Расписание не найдено")
		return nil
	}

	segments := strings.Split(text, " ")

	day := new(Day)

	today := int(time.Now().Weekday()) - 1

	if len(schedule.Days)-1 < today || today < 0 {
		today = 0
	}

	if len(schedule.Days)-1 >= today {
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
	doc := ctx.Message().Document
	// check if file is a json mime type
	if doc.MIME != "application/json" {
		return nil
	}
	admin := !ctx.Message().FromGroup()

	if ctx.Message().FromGroup() {
		of, err := ctx.Bot().AdminsOf(ctx.Message().Chat)
		if err != nil {
			return err
		}

		for _, a := range of {
			if a.User.ID == ctx.Message().Sender.ID {
				admin = true
			}
		}
	}

	if !admin {
		return ctx.Reply("Вы не администратор")
	}
	model := Schedule{}

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

	err = h.Service.CreateSchedule(context.Background(), model)
	if err != nil {
		return err
	}
	err = ctx.Reply("Запись успешно добавлена")
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

func (h *Handler) DownloadSchedule(ctx telebot.Context) error {
	// send schedule from database
	schedule, err := h.Service.GetSchedule(context.Background(), ctx.Message().Chat.ID)
	if err != nil {
		h.Logger.Error(err)
		return err
	}
	b, err := json.Marshal(schedule)
	if err != nil {
		h.Logger.Error(err)
		return err
	}
	file := telebot.FromReader(bytes.NewReader(b))
	doc := &telebot.Document{File: file, Caption: "расписание", MIME: "application/json", FileName: "schedule.json"}
	err = ctx.Reply(doc)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) DeleteSchedule(ctx telebot.Context) error {
	admin := !ctx.Message().FromGroup()

	if ctx.Message().FromGroup() {
		of, err := ctx.Bot().AdminsOf(ctx.Message().Chat)
		if err != nil {
			return err
		}

		for _, a := range of {
			if a.User.ID == ctx.Message().Sender.ID {
				admin = true
			}
		}
	}

	if !admin {
		return ctx.Reply("Вы не администратор")
	}
	// delete schedule from database using service
	err := h.Service.DeleteSchedule(context.Background(), ctx.Message().Chat.ID)
	if err != nil {
		h.Logger.Error(err)
		return err
	}
	err = ctx.Reply("Расписание удалено")
	if err != nil {
		return err
	}
	return nil
}
