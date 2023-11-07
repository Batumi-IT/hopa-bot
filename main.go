package main

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Check struct {
	Stupid bool
	Smart  bool
}

func main() {
	log.Println("Hopa bot started")
	defer log.Println("Hopa bot stopped")

	tgToken := os.Getenv("TELEGRAM_TOKEN")
	if tgToken == "" {
		log.Fatal("TELEGRAM_TOKEN env variable is not set")
	}

	openaiToken := os.Getenv("OPENAI_TOKEN")
	if openaiToken == "" {
		log.Fatal("OPENAI_TOKEN env variable is not set")
	}

	openaiClient := openai.NewClient(openaiToken)

	bot, err := tgbotapi.NewBotAPI(tgToken)
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

		message := strings.ToLower(update.Message.Text)
		if message == "" {
			continue
		}

		check := generateCheck(message)
		if !(check.Stupid || check.Smart) {
			continue
		}

		// TODO: Add rate limit by user id
		// First try to create answer with OpenAI
		replyMessage, err := generateOpenAiReply(openaiClient, message)
		if err != nil {
			log.Println(err)
			// If OpenAI fails, generate answer manually
			replyMessage = generateReply(message, check)
		}

		if replyMessage == "" {
			continue
		}

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, replyMessage)
		reply.ReplyToMessageID = update.Message.MessageID

		_, err = bot.Send(reply)
		if err != nil {
			log.Println(err)
		}
	}
}

func generateCheck(message string) Check {
	c := Check{
		Stupid: containsStupidQuestion(message),
		Smart:  containsSmartQuestion(message),
	}

	return c
}

func generateReply(message string, check Check) string {
	switch check {
	case Check{Stupid: true, Smart: false}:
		return "На рынке Хопа!"
	case Check{Stupid: false, Smart: true}:
		return "Держи ссылку с адресом рынка Хопа, раз в гугле забанили:\nhttps://goo.gl/maps/aqN4rzapdDXvRJNW9"
	case Check{Stupid: true, Smart: true}:
		return "Хопа на рынке Хопа! Вот, ну:\nhttps://goo.gl/maps/aqN4rzapdDXvRJNW9"
	default:
		return ""
	}
}

func containsStupidQuestion(message string) bool {
	var re = regexp.MustCompile(
		`(\s|^)(?:где|в)\s.*(?:купи(ть|л|ли|ла)|на(йти|шла|ш[её]л)|прода[её]тся|починить|посмотреть|продаже|доста(ть|л|ли|ла)|взя(ть|л|ли|ла)|покупа(л|ли|ла)).*\?`,
	)
	return re.MatchString(message)
}

func containsSmartQuestion(message string) bool {
	var re = regexp.MustCompile(
		`(\s|^)(?:где|как)\s.*(?:хоп[ау]|хоп[ауы]).*\?`,
	)
	return re.MatchString(message)
}

func generateOpenAiReply(client *openai.Client, message string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo1106,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Help people find what they are looking for at the рынок Хопа in Batumi, answering in Russian. Only if user ask about the address, answer with the link https://maps.app.goo.gl/MobUYTDFBKhuvMaR7. Answer sarcastically with jokes, puns, prejudices about clothing markets. Answer 1-2 sentences.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
			Temperature:      1.1,
			MaxTokens:        128,
			TopP:             1,
			Stop:             []string{"\n"},
			FrequencyPenalty: 0,
			PresencePenalty:  0.5,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
