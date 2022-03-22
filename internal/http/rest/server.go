package rest

import (
	"context"
	"encoding/json"
	"fmt"
	telebot "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"net/http"
	cache2 "telegram_bot/internal/cache"
	"telegram_bot/internal/config"
	"telegram_bot/internal/consts"
	"telegram_bot/internal/datastruct"
	"telegram_bot/internal/logger"
	"telegram_bot/internal/message_broker"
	"telegram_bot/internal/tools"
	"time"
)

type Server struct {
	ctx    context.Context
	client *http.Client
	broker *message_broker.ContestBroker
	cache  *cache2.RedisCache
	bot    *telebot.BotAPI
}

func NewServer(options ...ServerOption) *Server {
	srv := &Server{}
	for _, opt := range options {
		opt(srv)
	}
	return srv
}

func (s *Server) Run() {
	contests, err := s.GetContestsByTimeInterval()
	if err != nil {
		logger.Logger.Panic(err.Error())
	}

	for _, contest := range contests {
		err = s.cache.Set(s.ctx, message_broker.WithPrefix(contest.Id), contest)
		if err != nil {
			logger.Logger.Panic(err.Error())
		}
	}

	for {
		time.Sleep(10 * time.Second)
		for _, contest := range contests {
			if time.Until(contest.StartDate.Add(-1*time.Hour)) < 0 {
				msg := telebot.NewMessageToChannel(
					config.TelebotChannelName(),
					fmt.Sprintf("%s starts in ~ 1 hour", contest.Name),
				)
				str, err := s.bot.Send(msg)
				if err != nil {
					logger.Logger.Error(err.Error())
				}
				logger.Logger.Sugar().Debugf("%v", str)
			}
		}
	}
}

func (s *Server) GetContestsByTimeInterval() ([]*datastruct.Contest, error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, consts.GET_CONTESTS_BY_TIME_INTERVAL.Addr(), nil)
	if err != nil {
		return nil, err
	}

	req.Close = true

	q := req.URL.Query()
	q.Add("time_from", tools.Now().Format(time.RFC3339))
	q.Add("time_to", tools.EndOfTheDay(time.Now()).Format(time.RFC3339))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var contests []*datastruct.Contest
	err = json.Unmarshal(body, &contests)
	if err != nil {
		return nil, err
	}
	return contests, nil
}
