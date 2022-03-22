package rest

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	cache2 "telegram_bot/internal/cache"
	"telegram_bot/internal/message_broker"
)

type ServerOption func(srv *Server)

func ServerWithBroker(broker *message_broker.ContestBroker) ServerOption {
	return func(srv *Server) {
		srv.broker = broker
	}
}

func ServerWithClient(client *http.Client) ServerOption {
	return func(srv *Server) {
		srv.client = client
	}
}

func ServerWithCache(cache *cache2.RedisCache) ServerOption {
	return func(srv *Server) {
		srv.cache = cache
	}
}

func ServerWithTelebot(bot *tgbotapi.BotAPI) ServerOption {
	return func(srv *Server) {
		srv.bot = bot
	}
}
