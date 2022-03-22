package main

import (
	"context"
	telebot "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	cache2 "telegram_bot/internal/cache"
	"telegram_bot/internal/config"
	"telegram_bot/internal/http/rest"
	"telegram_bot/internal/logger"
	"telegram_bot/internal/message_broker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	bot, err := telebot.NewBotAPI(config.TelebotToken())
	if err != nil {
		logger.Logger.Panic(err.Error())
	}

	cache := cache2.NewRedisCache(config.RedisConn())
	err = cache.Ping(ctx)
	if err != nil {
		logger.Logger.Panic(err.Error())
	}

	brokerAddr := config.KafkaConn()
	broker := message_broker.NewContestBroker(cache, "peer1")
	err = broker.Connect(ctx, []string{brokerAddr})
	if err != nil {
		logger.Logger.Panic(err.Error())
	}

	defer broker.Close()

	catchTerminationFunc := func(cancel context.CancelFunc) {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		logger.Logger.Warn("caught termination signal")
		cancel()
	}

	go catchTerminationFunc(cancel)

	server := rest.NewServer(
		rest.ServerWithBroker(broker),
		rest.ServerWithClient(&http.Client{}),
		rest.ServerWithCache(cache),
		rest.ServerWithTelebot(bot),
	)

	server.Run()
}
