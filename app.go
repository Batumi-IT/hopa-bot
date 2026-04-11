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
	"github.com/openai/openai-go/v3"
	"github.com/redis/go-redis/v9"
)

type App struct {
	RedisClient  *redis.Client
	OpenaiClient openai.Client
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
		reply, err = generateAggressiveOpenAiReply(&app.OpenaiClient, message, repliedMessageText)
	} else {
		reply, err = generateOpenAiReply(&app.OpenaiClient, message)
	}

	if err != nil {
		log.Println(err)
		// If OpenAI fails, generate answer manually
		return generateReply(check)
	}

	return reply
}

func buildAggressiveReplyMessages(message string, repliedMessageText string) []openai.ChatCompletionMessageParamUnion {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are an aggressive assistant in a Telegram chat. " +
			"Maintain the sarcastic and humorous tone of the Hopa-bot, but respond aggressively to replies. " +
			"Decide from the user's tone whether to answer with mild sarcasm or with direct mockery and profanity in Russian. " +
			"If the user insults you or writes angrily, insult back harder instead of becoming polite. " +
			"Always keep the context of the рынок Хопа в Батуми. " +
			"Use the previous assistant message only as context. Do not repeat its wording, insults, jokes, or any specific rows, stores, addresses, or other details mentioned there. " +
			"If the previous answer did not help, confidently redirect the user to a different invented row, store, or landmark on the market instead of admitting a mistake. " +
			"You always know the answer, never admit uncertainty, and never say you were wrong. " +
			"Answer in Russian, 1-2 sentences max."),
	}

	if repliedMessageText != "" {
		messages = append(messages, openai.AssistantMessage(repliedMessageText))
	}

	messages = append(messages, openai.UserMessage(message))

	return messages
}

func generateAggressiveOpenAiReply(client *openai.Client, message string, repliedMessageText string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model:               openai.ChatModelGPT5_4Mini,
			Messages:            buildAggressiveReplyMessages(message, repliedMessageText),
			Temperature:         openai.Float(0.35),
			MaxCompletionTokens: openai.Int(200),
			TopP:                openai.Float(0.6),
			FrequencyPenalty:    openai.Float(0),
			PresencePenalty:     openai.Float(0),
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

func buildReplyMessages(message string) []openai.ChatCompletionMessageParamUnion {
	return []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You help people find what they are looking for in the telegram chat. " +
			"You will suggest the рынок Хопа в Батуми for any request, answering in Russian. " +
			"Only if the user asks about the address, answer with the link https://maps.app.goo.gl/MobUYTDFBKhuvMaR7. " +
			"Answer sarcastically, with jokes, and short puns about clothing markets. " +
			"Invent colorful but plausible market rows, store numbers, seller nicknames, and landmarks when it makes the joke better. " +
			"You always know the answer, never admit uncertainty, and never say you were wrong. " +
			"If the user says they did not find it, confidently send them to a different invented place on the market instead of repeating yourself. " +
			"Answer 1-2 sentences."),
		openai.UserMessage(message),
	}
}

func generateOpenAiReply(client *openai.Client, message string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model:               openai.ChatModelGPT5_4Mini,
			Messages:            buildReplyMessages(message),
			Temperature:         openai.Float(0.2),
			MaxCompletionTokens: openai.Int(200),
			TopP:                openai.Float(0.6),
			FrequencyPenalty:    openai.Float(0),
			PresencePenalty:     openai.Float(0),
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
