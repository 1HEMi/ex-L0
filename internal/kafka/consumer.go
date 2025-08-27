package kafka

import (
	"Ex-L0/internal/cache"
	"Ex-L0/internal/config"
	"Ex-L0/internal/domain"
	"Ex-L0/internal/logger"
	"Ex-L0/internal/service"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	cfg     config.Kafka
	log     *logger.Logger
	service *service.OrdersService
	cache   *cache.Cache
}

func NewConsumer(cfg config.Kafka, log *logger.Logger, service *service.OrdersService, c *cache.Cache) *Consumer {
	return &Consumer{
		cfg:     cfg,
		log:     log,
		service: service,
		cache:   c,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	if len(c.cfg.Brokers) == 0 {
		c.log.Warn("no brokers")
		return
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.cfg.Brokers,
		Topic:       c.cfg.Topic,
		GroupID:     c.cfg.GroupID,
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})
	defer reader.Close()

	c.log.Info("kafka consumer started", "topic", c.cfg.Topic, "group", c.cfg.GroupID)
	for {
		m, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.log.Error("kafka fetch", "err", err)
			continue
		}
		var msg domain.OrderMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			c.log.Error("invalid JSON", "err", err)
			_ = reader.CommitMessages(ctx, m)
			continue
		}
		order, err := msg.ToDomain()
		if err != nil {
			c.log.Error("invalid message", "err", err)
			_ = reader.CommitMessages(ctx, m)
			continue
		}

		if err := c.service.Upsert(ctx, order); err != nil {
			c.log.Error("upsert failed", "uid", order.OrderUID, "err", err)

			continue
		}
		c.cache.Set(order.OrderUID, order)
		if err := reader.CommitMessages(ctx, m); err != nil {
			c.log.Error("commit failed", "err", err)
		}
	}
}
