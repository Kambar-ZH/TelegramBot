package config

import (
	"os"
)

func TelebotToken() string {
	token := os.Getenv("TELEBOT_TOKEN")
	if token == "" {
		token = "1752339795:AAFomsaP4I3hr2Xh6QYi3s09Yyps6nEGGM4"
	}
	return token
}

func TelebotChannelName() string {
	channelName := os.Getenv("TELEBOT_CHANNEL_NAME")
	if channelName == "" {
		channelName = "@CompetitiveProgrammingPlatform"
	}
	return channelName
}

func RedisConn() string {
	conn := os.Getenv("REDIS_CONN")
	if conn == "" {
		conn = "localhost:6379"
	}
	return conn
}

func KafkaConn() string {
	conn := os.Getenv("KAFKA_CONN")
	if conn == "" {
		conn = "localhost:29092"
	}
	return conn
}
