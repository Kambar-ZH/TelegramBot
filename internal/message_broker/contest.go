package message_broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	cache2 "telegram_bot/internal/cache"
	"telegram_bot/internal/datastruct"
)

const (
	topic = "contest"
)

type (
	ContestBroker struct {
		consumerGroup sarama.ConsumerGroup

		consumeHandler contestConsumeHandler
		clientId       string
	}

	contestConsumeHandler struct {
		cache *cache2.RedisCache
		ready chan bool
	}
)

func NewContestBroker(cache *cache2.RedisCache, clientId string) *ContestBroker {
	return &ContestBroker{
		clientId: clientId,
		consumeHandler: contestConsumeHandler{
			cache: cache,
			ready: make(chan bool),
		},
	}
}

func (c *ContestBroker) Connect(ctx context.Context, brokers []string) error {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true
	consumerGroup, err := sarama.NewConsumerGroup(brokers, c.clientId, consumerConfig)
	if err != nil {
		panic(err)
	}
	c.consumerGroup = consumerGroup

	go func() {
		for {
			err := c.consumerGroup.Consume(ctx, []string{topic}, c.consumeHandler)
			if err != nil {
				log.Println(err)
			}
			if ctx.Err() != nil {
				return
			}
			c.consumeHandler.ready = make(chan bool)
		}
	}()
	<-c.consumeHandler.ready

	return nil
}

func (c *ContestBroker) Close() error {
	if err := c.consumerGroup.Close(); err != nil {
		return err
	}

	return nil
}

func (c contestConsumeHandler) Setup(session sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c contestConsumeHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c contestConsumeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		contestMsg := new(datastruct.ContestMessage)
		if err := json.Unmarshal(msg.Value, &contestMsg); err != nil {
			return err
		}
		switch contestMsg.Command {
		case datastruct.ContestCommandCreate:
			if contestMsg.Contest == nil {
				return errors.New("wtff")
			}
			err := c.cache.Set(context.Background(), WithPrefix(contestMsg.Contest.Id), contestMsg.Contest)
			return err
		default:
			fmt.Println("Undefined command")
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

func WithPrefix(id int32) string {
	return fmt.Sprintf("contests-%d", id)
}
