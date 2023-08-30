package main

import (
	"log"
	"os"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Checks struct {
	Stupid bool
	Smart  bool
}

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
		if message == "" {
			continue
		}

		check := Checks{
			Stupid: containsStupidQuestion(message),
			Smart:  containsSmartQuestion(message),
		}

		var replyMessage string

		// Note: add switch in case where will be more checks in the future
		switch check {
		case Checks{Stupid: false, Smart: false}:
			replyMessage = "На рынке Хопа!"
		case Checks{Stupid: true, Smart: false}:
			replyMessage = "Держи ссылку с адресом рынка Хопа, раз в гугле забанили: https://goo.gl/maps/aqN4rzapdDXvRJNW9"
		}

		if replyMessage != "" {
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, replyMessage)
			reply.ReplyToMessageID = update.Message.MessageID

			_, err := bot.Send(reply)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func containsStupidQuestion(message string) bool {
	var re = regexp.MustCompile(
		`где.*купить.*\?|где.*найти.*\?|где.*прода[её]тся.*\?|где.*починить.*\?|где.*посмотреть.*\?`,
	)
	return re.MatchString(strings.ToLower(message))
}

func containsSmartQuestion(message string) bool {
	var re = regexp.MustCompile(
		`где.*хопа.*\?|как.*хопа.*\?|где.*хопу.*\?|как.*хопу.*\?`,
	)
	return re.MatchString(strings.ToLower(message))
}
