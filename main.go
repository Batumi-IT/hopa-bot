package main

import (
	"github.com/go-redis/redis_rate/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
)

const (
	AiMessageMaxLength       = 200 // Max number of symbols to process (OpenAI)
	AiUserRateLimitPerMinute = 5   // Max number of messages per minute per user (OpenAI)
	AiUserRateLimitPerDay    = 30  // Max number of messages per day per user (OpenAI)
	AiTotalRateLimitPerDay   = 300 // Total max number of messages per day (OpenAI)
)

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

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	rdb := connectToRedis(RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	app := App{
		RedisClient:  rdb,
		OpenaiClient: openai.NewClient(openaiToken),
		TelegramBot:  bot,
		RedisLimiter: redis_rate.NewLimiter(rdb),
	}

	app.run()
}
