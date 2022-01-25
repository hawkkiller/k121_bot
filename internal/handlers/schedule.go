package handlers

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hawkkiller/k121_bot/internal"
	"github.com/hawkkiller/k121_bot/internal/model"
	"github.com/hawkkiller/k121_bot/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"regexp"
	"strings"
)

func HandleSchedule(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var s = new(model.Schedule)
	db := internal.DB.Preload(clause.Associations).Preload("Times").Preload("Days.Pairs").Where("chat_id=?", update.FromChat().ID).Last(&s)
	if db.Error != nil {
		err := db.Error.Error()
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			err = "Расписание не **найдено**. Вы можете загрузить его воспользовавшись командой __/uploadSchedule__ к прикреплённому объекту JSON."
		}

		utils.SendMessage(bot, update.FromChat().ID, err)

	}
	days := s.Days
	ds := strings.Split(update.Message.Text, " ")

	if len(ds) == 2 {
		for _, d := range days {
			if strings.ToLower(d.Caption) == strings.ToLower(ds[1]) {
				printDay(d, bot, update)
				return
			}
		}
		utils.SendMessage(bot, update.FromChat().ID, "Нет такого дня в распорядке :(")
	} else {
		for _, d := range days {
			printDay(d, bot, update)
		}
	}
}

func printDay(d model.Day, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := fmt.Sprintf("*%s*", d.Caption)
	regex := regexp.MustCompile("\\([^)]*\\)")
	for _, p := range d.Pairs {
		codes := strings.Split(p.AdditionalInfo, " ")
		additionalInfo := ""
		for _, c := range codes {
			found := regex.Find([]byte(c))
			if found == nil {
				found = []byte(c)
			}
			str := strings.Replace(strings.Replace(string(found[:]), "(", "", -1), ")", "", -1)
			strSplit := strings.Split(str, ":")
			if len(strSplit) != 2 {
				additionalInfo += fmt.Sprintf("%s\n", c)
				continue
			}
			strFormatted := fmt.Sprintf("`%s` : `%s`", strSplit[0], strSplit[1])
			newCode := strings.Replace(c, str, strFormatted, -1)
			additionalInfo += fmt.Sprintf("%s\n", newCode)
		}
		message += fmt.Sprintf("\n*%s*\n_Код_: %s\n__Аудитория__: %s", p.Title, additionalInfo, p.Auditory)
	}
	utils.SendMessage(bot, update.FromChat().ID, message)
}