package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis_rate/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
)

type App struct {
	RedisClient  *redis.Client
	OpenaiClient *openai.Client
	TelegramBot  *tgbotapi.BotAPI
	RedisLimiter *redis_rate.Limiter
}

type Check struct {
	Stupid bool
	Smart  bool
}

func (app *App) run() {
	ctx := context.Background()

	// Set up updates channel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := app.TelegramBot.GetUpdatesChan(u)
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

		var isReplyToBot bool
		var repliedMessageText string
		if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From != nil {
			isReplyToBot = update.Message.ReplyToMessage.From.ID == app.TelegramBot.Self.ID
			repliedMessageText = update.Message.ReplyToMessage.Text
		}

		check := generateCheck(message)
		if !(check.Stupid || check.Smart) && !isReplyToBot {
			continue
		}

		replyMessage := app.generateReplyMessage(ctx, message, check, update.Message.From.ID, isReplyToBot, repliedMessageText)
		if replyMessage == "" {
			continue
		}

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, replyMessage)
		reply.ReplyToMessageID = update.Message.MessageID

		_, err = app.TelegramBot.Send(reply)
		if err != nil {
			log.Println(err)
		}
	}
}

func (app *App) generateReplyMessage(ctx context.Context, message string, check Check, userID int, isReplyToBot bool, repliedMessageText string) string {
	if len(message) > AiMessageMaxLength {
		return generateReply(check)
	}

	// If global limit is reached for today or user limit is reached, generate answer manually
	globalLimit, err := app.RedisLimiter.Allow(
		ctx,
		"global:ai_total_rate_limit",
		perDay(AiTotalRateLimitPerDay),
	)
	if err != nil {
		log.Fatal(err)
	}
	userDailyLimit, err := app.RedisLimiter.Allow(
		ctx,
		fmt.Sprintf("user:%d:ai_daily_limit", userID),
		perDay(AiUserRateLimitPerDay),
	)
	if err != nil {
		log.Fatal(err)
	}
	if globalLimit.Allowed == 0 || userDailyLimit.Allowed == 0 {
		return generateReply(check)
	}

	// If users minute limit is reached - escort him to the nearest hyi :)
	userMinuteLimit, err := app.RedisLimiter.Allow(
		ctx,
		fmt.Sprintf("user:%d:ai_minute_limit", userID),
		redis_rate.PerMinute(AiUserRateLimitPerMinute),
	)
	if err != nil {
		log.Fatal(err)
	}
	if userMinuteLimit.Allowed == 0 {
		return "Error 429: Пошёл на хуй, пёс!"
	}

	var reply string
	if isReplyToBot {
		reply, err = generateAggressiveOpenAiReply(app.OpenaiClient, message+"\n\nContext: "+repliedMessageText)
	} else {
		reply, err = generateOpenAiReply(app.OpenaiClient, message)
	}

	if err != nil {
		log.Println(err)
		// If OpenAI fails, generate answer manually
		return generateReply(check)
	}

	return reply
}

func generateAggressiveOpenAiReply(client *openai.Client, message string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: "You are an aggressive assistant in a Telegram chat. " +
						"Maintain the sarcastic and humorous tone of the Hopa-bot, but respond aggressively to replies. " +
						"If the user is unkind or rude, swear at them in Russian, using sharp humor and insults. " +
						"Always keep the context of the рынок Хопа в Батуми. " +
						"Mock the user, make jokes, and use sharp humor, but stay within the context of the market. " +
						"Answer in Russian, 1-2 sentences max.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
			Temperature:      0.9,
			MaxTokens:        200,
			TopP:             1,
			Stop:             []string{"\n"},
			FrequencyPenalty: 0.4,
			PresencePenalty:  0.4,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func generateCheck(message string) Check {
	m := strings.ToLower(message)
	return Check{
		Stupid: containsStupidQuestion(m),
		Smart:  containsSmartQuestion(m),
	}
}

func generateReply(check Check) string {
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: "You help people find what they are looking for in the telegram chat. " +
						"You will suggest the рынок Хопа в Батуми for any request, answering in Russian. " +
						"Only if the user asks about the address, answer with the link https://maps.app.goo.gl/MobUYTDFBKhuvMaR7. " +
						"Answer sarcastically, with jokes, puns, prejudices about clothing markets. " +
						"If someone asks where to find something on Хопа - you can randomly generate the market row number and the store. " +
						"Answer 1-2 sentences.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
			Temperature:      0.85,
			MaxTokens:        200,
			TopP:             1,
			Stop:             []string{"\n"},
			FrequencyPenalty: 0.4,
			PresencePenalty:  0.4,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// perDay is a missing function in redis_rate package
func perDay(rate int) redis_rate.Limit {
	return redis_rate.Limit{
		Rate:   rate,
		Period: 24 * time.Hour,
		Burst:  rate,
	}
}
