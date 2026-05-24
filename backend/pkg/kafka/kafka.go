// Package kafka provides Sarama-based producer and consumer wrappers for WanderPlan services.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

// Topic constants shared by all services.
const (
	TopicAuthEvents   = "auth-events"
	TopicTripEvents   = "trip-events"
	TopicCollabEvents = "collab-events"
)

// Event is the standard envelope for all Kafka messages.
type Event struct {
	Type      string          `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// Producer wraps a Sarama sync producer.
type Producer struct {
	client sarama.SyncProducer
}

// NewProducer creates a connected Sarama sync producer.
func NewProducer(brokers []string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_6_0_0
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Retry.Backoff = 250 * time.Millisecond

	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("create sarama producer: %w", err)
	}
	log.Info().Strs("brokers", brokers).Msg("kafka producer connected")
	return &Producer{client: p}, nil
}

// Publish encodes payload as JSON and sends an event to the given topic.
func (p *Producer) Publish(ctx context.Context, topic, eventType string, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	env, err := json.Marshal(Event{
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Payload:   b,
	})
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(env),
	}
	partition, offset, err := p.client.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("send message to %s: %w", topic, err)
	}

	log.Debug().
		Str("topic", topic).
		Str("event_type", eventType).
		Int32("partition", partition).
		Int64("offset", offset).
		Msg("kafka message published")
	return nil
}

// Close closes the underlying Sarama producer gracefully.
func (p *Producer) Close() error { return p.client.Close() }

// Handler is a callback invoked for each consumed message.
type Handler func(ctx context.Context, event Event) error

// Consumer wraps a Sarama consumer group.
type Consumer struct {
	group    sarama.ConsumerGroup
	topics   []string
	handlers map[string]Handler
	groupID  string
}

// NewConsumer creates a consumer group for the given topics.
func NewConsumer(brokers []string, groupID string, topics []string) (*Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_6_0_0
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, fmt.Errorf("create consumer group: %w", err)
	}
	log.Info().Strs("brokers", brokers).Str("group", groupID).Strs("topics", topics).Msg("kafka consumer group created")
	return &Consumer{group: group, topics: topics, handlers: make(map[string]Handler), groupID: groupID}, nil
}

// Register associates an event type with a handler function.
func (c *Consumer) Register(eventType string, h Handler) {
	c.handlers[eventType] = h
}

// Start begins consuming messages and dispatches to registered handlers.
// It blocks until ctx is cancelled.
func (c *Consumer) Start(ctx context.Context) error {
	ch := &consumerGroupHandler{handlers: c.handlers}
	for {
		if err := c.group.Consume(ctx, c.topics, ch); err != nil {
			return fmt.Errorf("consumer group error: %w", err)
		}
		if ctx.Err() != nil {
			return nil
		}
	}
}

// Close shuts down the consumer group.
func (c *Consumer) Close() error { return c.group.Close() }

// consumerGroupHandler implements sarama.ConsumerGroupHandler.
type consumerGroupHandler struct {
	handlers map[string]Handler
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Error().Err(err).Bytes("raw", msg.Value).Msg("failed to unmarshal kafka event")
			session.MarkMessage(msg, "")
			continue
		}

		if handler, ok := h.handlers[event.Type]; ok {
			if err := handler(session.Context(), event); err != nil {
				log.Error().Err(err).Str("event_type", event.Type).Msg("event handler error")
			}
		} else {
			log.Debug().Str("event_type", event.Type).Msg("no handler registered for event type")
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
