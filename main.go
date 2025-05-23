package main

import (
	"log"
	"os"

	"github.com/go-redis/redis_rate/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
)

const (
	AiMessageMaxLength       = 250 // Max number of symbols to process (OpenAI)
	AiUserRateLimitPerMinute = 5   // Max number of messages per minute per user (OpenAI)
	AiUserRateLimitPerDay    = 40  // Max number of messages per day per user (OpenAI)
	AiTotalRateLimitPerDay   = 300 // Total max number of messages per day (OpenAI)
)

func main() {
	// TODO: add validation for all env variables with viper package
	log.Println("Hopa bot is started")
	defer log.Println("Hopa bot stopped")

	tgToken := os.Getenv("TELEGRAM_TOKEN")
	if tgToken == "" {
		log.Fatal("TELEGRAM_TOKEN env variable is not set")
	}

	openaiToken := os.Getenv("OPENAI_TOKEN")
	if openaiToken == "" {
		log.Fatal("OPENAI_TOKEN env variable is not set")
	}

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	rdb := connectToRedis(RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		User:     os.Getenv("REDIS_USER"),
	})

	app := App{
		RedisClient:  rdb,
		OpenaiClient: openai.NewClient(openaiToken),
		TelegramBot:  bot,
		RedisLimiter: redis_rate.NewLimiter(rdb),
	}

	log.Println("Hopa bot is running")
	app.run()
}
