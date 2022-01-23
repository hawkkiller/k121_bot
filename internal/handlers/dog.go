package handlers

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hawkkiller/k121_bot/internal/model"
	"io/ioutil"
	"net/http"
)

func HandleDog(bot *tgbotapi.BotAPI, id int64) {
	dogUrl, dogUrlErr := http.Get("https://dog.ceo/api/breeds/image/random")
	if dogUrlErr != nil {
		return
	}
	defer dogUrl.Body.Close()
	var dog = new(model.DogPhoto)
	jsonError := json.NewDecoder(dogUrl.Body).Decode(dog)

	if jsonError != nil {
		return
	}
	dogImg, dogImgErr := http.Get(dog.Message)
	defer dogImg.Body.Close()
	if dogImgErr != nil {
		return
	}
	bytes, _ := ioutil.ReadAll(dogImg.Body)
	file := tgbotapi.FileBytes{
		Name:  "image.jpg",
		Bytes: bytes,
	}
	conf := tgbotapi.NewPhoto(id, file)
	if _, err := bot.Send(conf); err != nil {
		fmt.Println(err)
	}
}

//if _, err := bot.Send(msg); err != nil {
//fmt.Println(err)
//}
