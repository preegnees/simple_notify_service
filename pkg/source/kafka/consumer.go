package kafka

import (
	"context"
	"fmt"

	models "notifier/pkg/models"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
)

type kafka_ struct {
	EventCh    chan models.Message
	conn       *kafka.Conn
	ctx        context.Context
	topic      string
	partition  int
	host       string
	port       string
	minByte    int
	maxByte    int
	bufferSize int
}

type ConfKafkaConsumer struct {
	Ctx        context.Context
	Topic      string
	Partition  int
	Host       string
	Port       string
	MinByte    int
	MaxByte    int
	BufferSize int
}

func New(cnf ConfKafkaConsumer) (*kafka_, error) {

	conn, err := kafka.DialLeader(context.Background(), "tcp", cnf.Host+":"+cnf.Port, cnf.Topic, cnf.Partition) // tcp, localhost:29092, my-topic, 0
	if err != nil {
		return nil, fmt.Errorf("failed to dial leader: %w", err)
	}
	return &kafka_{
		EventCh:    make(chan models.Message),
		ctx:        cnf.Ctx,
		conn:       conn,
		topic:      cnf.Topic,
		partition:  cnf.Partition,
		host:       cnf.Host,
		port:       cnf.Port,
		minByte:    cnf.MinByte,
		maxByte:    cnf.MaxByte,
		bufferSize: cnf.BufferSize,
	}, nil
}

func (c *kafka_) Run() error {

	batch := c.conn.ReadBatch(c.minByte, c.maxByte) // 1 1e6
	defer c.conn.Close()
	defer batch.Close()

	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			mes, err := batch.ReadMessage()
			if err != nil {
				break
			}
			hs := mes.Headers
			message, err := c.convertToMessage(hs)
			if err != nil {
				return err
			}
			c.EventCh <- message
		}
	}

	// b := make([]byte, c.bufferSize)                 // 10e3
	// for {
	// 	n, err := batch.ReadMessage()
	// 	if err != nil {
	// 		break
	// 	}
	// 	c.EventCh <- b[:n]
	// }
	// return nil
}

func (c *kafka_) convertToMessage(hs []protocol.Header) (models.Message, error) {
	message := models.Message{}
	for _, h := range hs {
		v := h.Key
		if v == "Destination" {
			message.Destination = string(h.Value)
		} else if v == "Email" {
			message.Email = string(h.Value)
		} else if v == "Username" {
			message.Username = string(h.Value)
		} else if v == "MessageSubject" {
			message.MessageSubject = string(h.Value)
		} else if v == "Message" {
			message.Message = string(h.Value)
		} else {
			return models.Message{}, fmt.Errorf("Err, invalid message")
		}
	}
	return message, nil
}
