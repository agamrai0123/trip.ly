package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

// kafkaEvent is the envelope produced by pkg/kafka.Producer.
type kafkaEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// EventConsumer subscribes to multiple Kafka topics and persists notifications.
type EventConsumer struct {
	group    sarama.ConsumerGroup
	topics   []string
	repo     *NotificationRepo
	hub      *Hub
	ctx      context.Context
	cancel   context.CancelFunc
	lagGauge *prometheus.GaugeVec
}

// NewEventConsumer creates and starts a Kafka consumer group.
// It registers a kafka_consumer_lag gauge into reg.
func NewEventConsumer(brokers []string, groupID string, topics []string, repo *NotificationRepo, hub *Hub, reg *prometheus.Registry) (*EventConsumer, error) {
	scfg := sarama.NewConfig()
	scfg.Version = sarama.V2_1_0_0
	scfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	scfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(brokers, groupID, scfg)
	if err != nil {
		return nil, err
	}

	lagGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_consumer_lag",
		Help: "Kafka consumer lag per topic and partition.",
	}, []string{"topic", "partition"})
	reg.MustRegister(lagGauge)

	ctx, cancel := context.WithCancel(context.Background())
	return &EventConsumer{
		group:    group,
		topics:   topics,
		repo:     repo,
		hub:      hub,
		ctx:      ctx,
		cancel:   cancel,
		lagGauge: lagGauge,
	}, nil
}

// Start launches the consumer loop in a goroutine.
func (ec *EventConsumer) Start() {
	go func() {
		for {
			if err := ec.group.Consume(ec.ctx, ec.topics, ec); err != nil {
				if ec.ctx.Err() != nil {
					return
				}
				log.Error().Err(err).Msg("kafka consumer error")
			}
		}
	}()
}

// Close shuts down the consumer group.
func (ec *EventConsumer) Close() error {
	ec.cancel()
	return ec.group.Close()
}

// Setup implements sarama.ConsumerGroupHandler.
func (ec *EventConsumer) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup implements sarama.ConsumerGroupHandler.
func (ec *EventConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from a single partition.
func (ec *EventConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		lag := claim.HighWaterMarkOffset() - msg.Offset - 1
		if lag < 0 {
			lag = 0
		}
		ec.lagGauge.WithLabelValues(msg.Topic, fmt.Sprintf("%d", msg.Partition)).Set(float64(lag))
		ec.handle(msg)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (ec *EventConsumer) handle(msg *sarama.ConsumerMessage) {
	var env kafkaEvent
	if err := json.Unmarshal(msg.Value, &env); err != nil {
		log.Warn().Err(err).Msg("unmarshal kafka event")
		return
	}

	// Map event types to notification records.
	var n Notification
	switch env.Type {
	case "trip.created":
		var p struct {
			UserID string `json:"user_id"`
		}
		_ = json.Unmarshal(env.Payload, &p)
		n = Notification{UserID: p.UserID, Type: "trip.created", Title: "Trip created", Body: "Your new trip has been created."}
	case "trip.updated":
		var p struct {
			UserID string `json:"user_id"`
		}
		_ = json.Unmarshal(env.Payload, &p)
		n = Notification{UserID: p.UserID, Type: "trip.updated", Title: "Trip updated", Body: "A trip you follow was updated."}
	case "collaboration.invited":
		var p struct {
			UserID string `json:"user_id"`
		}
		_ = json.Unmarshal(env.Payload, &p)
		n = Notification{UserID: p.UserID, Type: "collaboration.invited", Title: "Invitation", Body: "You have been invited to collaborate on a trip."}
	case "user.registered":
		var p struct {
			UserID string `json:"user_id"`
		}
		_ = json.Unmarshal(env.Payload, &p)
		n = Notification{UserID: p.UserID, Type: "user.registered", Title: "Welcome to WanderPlan!", Body: "Your account is ready."}
	default:
		return
	}

	if n.UserID == "" {
		return
	}
	created, err := ec.repo.Create(context.Background(), &n)
	if err != nil {
		log.Error().Err(err).Msg("persist notification")
		return
	}
	ec.hub.Broadcast(created.UserID, created)
}
