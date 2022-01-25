package handlers

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hawkkiller/k121_bot/internal"
	"github.com/hawkkiller/k121_bot/internal/model"
	"github.com/hawkkiller/k121_bot/pkg/utils"
	"gorm.io/gorm/clause"
	"net/http"
)

func HandleUploadSchedule(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	user := update.SentFrom()
	config := update.FromChat().ChatConfig()
	administratorsConfig := tgbotapi.ChatAdministratorsConfig{ChatConfig: config}
	admins, aErr := bot.GetChatAdministrators(administratorsConfig)
	var trust = false
	if aErr != nil {
		if update.FromChat().IsPrivate() {
			trust = true
		} else {
			return
		}
	}
	for _, admin := range admins {
		if admin.User.ID == user.ID {
			trust = true
			break
		}
	}
	if trust {
		document := update.Message.Document

		if document == nil {
			utils.SendMessage(bot, update.FromChat().ID, "Вы не прикрепили файл")
			return
		}
		fileConfig := tgbotapi.FileConfig{
			FileID: document.FileID,
		}
		utils.SendMessage(bot, update.FromChat().ID, fmt.Sprintf("Файл получен. Размер: %d", document.FileSize))
		if document.MimeType != "application/json" {
			utils.SendMessage(bot, update.FromChat().ID, "Файл должен быть с расширением .json")
			return
		}

		getFile, fileErr := bot.GetFile(fileConfig)
		if fileErr != nil {
			utils.SendMessage(bot, update.FromChat().ID, "Не удаётся найти файл на сервере")
			return
		}
		url, err := bot.GetFileDirectURL(getFile.FileID)
		if err != nil {
			utils.SendMessage(bot, update.FromChat().ID, err.Error())
			return
		}

		utils.SendMessage(bot, update.FromChat().ID, "Скачиваю по этому URL: "+url)

		res, resErr := http.Get(url)

		if resErr != nil {
			utils.SendMessage(bot, update.FromChat().ID, "Ошибка при выполнении запроса на получение файла")
			return
		}
		defer res.Body.Close()

		schedule := new(model.Schedule)

		mErr := json.NewDecoder(res.Body).Decode(&schedule)
		if mErr != nil {
			utils.SendMessage(bot, update.FromChat().ID, "Ошибка при трансформировании json - невалидная структура")
			return
		}
		schedule.ChatId = update.FromChat().ChatConfig().ChatID
		foundSchedule := new(model.Schedule)
		if res := internal.DB.Debug().Where("chat_id=?", schedule.ChatId).Preload(clause.Associations).Preload("Days.Pairs").First(&foundSchedule); res.Error != nil {
			db := internal.DB.Debug().Create(&schedule)
			if db.Error != nil {
				utils.SendMessage(bot, update.FromChat().ID, db.Error.Error())
				return
			}
		} else {
			err := internal.DB.Debug().Model(&foundSchedule).Association("Days").Replace(schedule.Days)
			if err != nil {
				utils.SendMessage(bot, update.FromChat().ID, err.Error())
				return
			}
			err = internal.DB.Debug().Model(&foundSchedule).Association("Times").Replace(schedule.Times)
			if err != nil {
				utils.SendMessage(bot, update.FromChat().ID, err.Error())
				return
			}
			internal.DB.Debug().Model(&foundSchedule).Update("timezone", schedule.Timezone)
		}
	} else {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Изменять могут только администраторы")
		if _, err := bot.Send(msg); err != nil {
			fmt.Println(err)
		}
	}

}
