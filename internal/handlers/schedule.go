package handlers

import (
	"errors"
	"fmt"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
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
	//defer func() {
	//	if r := recover(); r != nil {
	//		utils.SendMessage(bot, update.FromChat().ID, fmt.Sprintf("Encountered an error %s", r))
	//	}
	//}()
	var s = new(model.Schedule)
	db := internal.DB.Preload(clause.Associations).Preload("Times").Preload("Days.Pairs").Where("chat_id=?", update.FromChat().ID).Last(&s)
	if db.Error != nil {
		err := db.Error.Error()
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			err = "Расписание не **найдено**. Вы можете загрузить его воспользовавшись командой __/uploadSchedule__ к прикреплённому объекту JSON."
		}

		utils.SendMessage(bot, update.FromChat().ID, err)
		return
	}
	days := s.Days
	ds := strings.Split(update.Message.Text, " ")

	if len(ds) == 2 {
		if strings.ToLower(ds[1]) == "все" {
			for _, d := range days {
				printDay(d, bot, update)
			}
			return
		}
		for _, d := range days {
			if strings.ToLower(d.Caption) == strings.ToLower(ds[1]) {
				printDay(d, bot, update)
				return
			}
		}
		/// try to find string that is very similar
		args := make([]dayIdentifier, 0)
		for _, d := range days {
			args = append(args, dayIdentifier{
				Title:    d.Caption,
				Distance: strutil.Similarity(d.Caption, update.Message.Text, metrics.NewJaro()),
			})
		}

		var biggest = dayIdentifier{}
		for _, i := range args {
			if i.Distance > biggest.Distance {
				biggest = i
			}
		}
		for _, d := range days {
			if d.Caption == biggest.Title {
				printDay(d, bot, update)
				return
			}
		}

		utils.SendMessage(bot, update.FromChat().ID, "Нет такого дня в распорядке :(")
	} else {
		now := int(update.Message.Time().Weekday()) - 1

		if len(days) >= now && now >= 0 {
			today := days[now]
			printDay(today, bot, update)
			return
		} else {
			printDay(days[0], bot, update)
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
				additionalInfo += fmt.Sprintf("%s", c)
				continue
			}
			strFormatted := fmt.Sprintf("`%s` : `%s`", strSplit[0], strSplit[1])
			newCode := strings.Replace(c, str, strFormatted, -1)
			additionalInfo += fmt.Sprintf("\n%s", newCode)
		}
		message += fmt.Sprintf("\n\n*%s*\n_Код_: %s\n__Аудитория__: %s", p.Title, additionalInfo, p.Auditory)
	}
	utils.SendMessage(bot, update.FromChat().ID, message)
}

type dayIdentifier struct {
	Title    string
	Distance float64
}
