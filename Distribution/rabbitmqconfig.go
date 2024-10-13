package Distribution

import (
	"genreport/Startup/Config"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"math"
	"time"
)

type RabbitMQConnection struct {
	Broker *amqp.Connection
}

// InitBroker connects to the RabbitMQ server and sets the connection to the server

func NewConnection() (*RabbitMQConnection, error) {
	Logger := Config.GetLogger()
	Settings := Config.GetSettings()

	connection, err := amqp.Dial(Settings.AMQPServerURL)
	if err != nil {
		Logger.Error("error running the application cannot connect to RabbitMQ", zap.Error(err))
		return nil, err
	}
	return &RabbitMQConnection{
		Broker: connection,
	}, nil

}

// CheckConnection pings RabbitMQ server to check if connection is still active
func (r *RabbitMQConnection) CheckConnection(connection *amqp.Connection) bool {
	return !(connection == nil || connection.IsClosed())
}

func (r *RabbitMQConnection) PublishMessageTextImmediate(exchangeName string, routingKey string, data []byte) error {
	err := r.publishMessageImmediate(exchangeName, routingKey, data, "text/plain")
	if err != nil {
		Config.GetLogger().Error("error publishing message immediate", zap.Error(err))
	}
	return err
}

func (r *RabbitMQConnection) publishMessageImmediateAsync(exchangeName string, routingKey string, data []byte, encoding string) {
	go func() {
		err := r.publishMessageImmediate(exchangeName, routingKey, data, encoding)
		if err != nil {
			Config.GetLogger().Error("error executing publish message", zap.Error(err))
		}
	}()

}
func (r *RabbitMQConnection) publishMessageImmediate(exchangeName string, routingKey string, data []byte, encoding string) error {
	channel, err := r.Broker.Channel()
	if err != nil {
		Config.GetLogger().Error("error creating the channel", zap.Error(err))
		return err
	}

	err = channel.Publish(exchangeName, routingKey, true, true, amqp.Publishing{
		Headers:         nil,
		ContentType:     "",
		ContentEncoding: encoding,
		DeliveryMode:    amqp.Transient,
		Priority:        math.MaxInt,
		Timestamp:       time.Time{},
		Body:            data,
	})
	if err != nil {
		Config.GetLogger().Error("error publishing json message", zap.Error(err))
		return err
	}
	return nil
}

func (r *RabbitMQConnection) CreateExchange(name string, exchangeType string) error {
	channel, err := r.Broker.Channel()
	if err != nil {
		Config.GetLogger().Error("error creating the channel", zap.Error(err))
		return err
	}
	err = channel.ExchangeDeclare(name, exchangeType, true, false, false, false, nil)
	if err != nil {
		Config.GetLogger().Error("error creating the channel", zap.Error(err))
		return err
	}
	return nil
}

func (r *RabbitMQConnection) GetOrCreateQueue(name string) (amqp.Queue, error) {
	channel, err := r.Broker.Channel()
	if err != nil {
		Config.GetLogger().Error("error creating the channel", zap.Error(err))
	}

	return channel.QueueDeclarePassive(name, true, false, false, false, nil)
}

func (r *RabbitMQConnection) CloseConnection() {
	if !r.Broker.IsClosed() {
		err := r.Broker.Close()
		if err != nil {
			Config.GetLogger().Error("error closing the connection", zap.Error(err))
		}

	}
}
