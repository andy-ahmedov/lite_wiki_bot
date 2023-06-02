package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func telegramBot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			switch update.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, i'm a wikipedia bot, i can search information in a wikipedia, send me something what you want find in Wikipedia.")
				bot.Send(msg)

			case "/number_of_users":

				if os.Getenv("DB_SWITCH") == "on" {
					num, err := getNumberOfUsers()
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
						bot.Send(msg)
					}

					ans := fmt.Sprintf("%d peoples used me for search infromation in Wikipedia", num)

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database not connected, so i can't say you how many peoples used me.")
					bot.Send(msg)
				}

			default:

				language := os.Getenv("LANGUAGE")

				ms, _ := urlEncoded(update.Message.Text)
				url := ms

				request := "https://" + language + ".wikipedia.org/w/api.php?action=opensearch&search=" + url + "&limit=3&origin=*&format=json"
				message := WikipediaAPI(request)

				if os.Getenv("DB_SWITCH") == "on" {
					if err := collectData(update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.Text, message); err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error, but bot still working.")
						bot.Send(msg)
					}
				}

				for _, val := range message {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, val)
					bot.Send(msg)
				}
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}
}

func main() {
	time.Sleep(1 * time.Minute)

	if os.Getenv("CREATE_TABLE") == "yes" {
		if os.Getenv("DB_SWITCH") == "on" {
			if err := createTable(); err != nil {
				panic(err)
			}
		}
	}

	time.Sleep(1 * time.Minute)

	telegramBot()
}
