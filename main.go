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

		replyMessage := generateReply(message)
		if replyMessage == "" {
			continue
		}

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, replyMessage)
		reply.ReplyToMessageID = update.Message.MessageID

		_, err := bot.Send(reply)
		if err != nil {
			log.Println(err)
		}
	}
}

func generateReply(message string) string {
	message = strings.ToLower(message)

	// Note: add if-else statement in case there will be more checks in the future
	if stupidHopaBazarQuestion(message) && smartHopaBazarQuestion(message) {
		return "Хопа на рынке Хопа! Вот, ну:\nhttps://goo.gl/maps/aqN4rzapdDXvRJNW9"
	} else if smartHopaBazarQuestion(message) {
		return "Держи ссылку с адресом рынка Хопа, раз в гугле забанили:\nhttps://goo.gl/maps/aqN4rzapdDXvRJNW9"
	} else if stupidHopaBazarQuestion(message) {
		return "На рынке Хопа!"
	} else if toxicPositivity(message) {
		return "Идите на хуй со своей токсичной позитивностью!"
	}
	return ""
}

func stupidHopaBazarQuestion(message string) bool {
	var re = regexp.MustCompile(
		`где.*купить.*\?|где.*найти.*\?|где.*прода[её]тся.*\?|где.*починить.*\?|где.*посмотреть.*\?`,
	)
	return re.MatchString(message)
}

func smartHopaBazarQuestion(message string) bool {
	var re = regexp.MustCompile(
		`где.*хоп[ау].*\?|как.*хоп[ауы].*\?`,
	)
	return re.MatchString(message)
}

func toxicPositivity(message string) bool {
	var re = regexp.MustCompile(
		`бы(?:ть|л|ла|ли).*позитивн`,
	)
	return re.MatchString(message)
}
