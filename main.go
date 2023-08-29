package main

import (
	"log"
	"os"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	log.Println("Hopa bot started")
	defer log.Println("Hopa bot stopped")

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN env variable is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// Set up updates channel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	// Listen for messages in group chats
	for update := range updates {
		if update.Message == nil {
			continue
		}

		message := update.Message.Text

		if message != "" && containsStupidQuestion(message) {
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, "На рынке Хопа!")
			_, err := bot.Send(reply)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func containsStupidQuestion(message string) bool {
	var re = regexp.MustCompile(`где.*купить.*\?|где.*найти.*\?|где.*прода[её]тся.*\?`)
	return re.MatchString(strings.ToLower(message))
}
