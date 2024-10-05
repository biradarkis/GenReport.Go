package Distribution

import (
	"genreport/Startup/Config"
	"genreport/Startup/Models"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"math"
	"sync"
	"time"
)

var Broker *amqp.Connection
var Settings *Models.Settings
var Logger *zap.Logger

type RabbitMQBroker struct {
	Settings *Models.Settings
	Logger   *zap.Logger
}

// InitBroker connects to the RabbitMQ server and sets the connection to the server
func (r *RabbitMQBroker) InitBroker() {
	Logger = Config.GetLogger()
	settings, err := Config.GetSettings()
	if err != nil {
		Logger.Error("error running the application cannot get settings", zap.Error(err))
		return
	}
	Settings = settings

	Broker, err = amqp.Dial(Settings.AMQPServerURL)
	if err != nil {
		Logger.Error("error running the application cannot connect to RabbitMQ", zap.Error(err))
		return
	}

}

// CheckConnection pings RabbitMQ server to check if connection is still active
func (r *RabbitMQBroker) CheckConnection(connection *amqp.Connection) bool {
	return !(connection == nil || connection.IsClosed())
}

func (r *RabbitMQBroker) GetRabbitMQConnection() *amqp.Connection {
	once := sync.Once{}
	once.Do(func() {
		r.InitBroker()
	})
	return Broker
}

func (r *RabbitMQBroker) PublishMessageTextImmediate(exchangeName string, routingKey string, data []byte) error {
	err := r.publishMessageImmediate(exchangeName, routingKey, data, "text/plain")
	if err != nil {
		Config.GetLogger().Error("error publishing message immediate", zap.Error(err))
	}
	return err
}

func (r *RabbitMQBroker) publishMessageImmediateAsync(exchangeName string, routingKey string, data []byte, encoding string) {
	go func() {
		err := r.publishMessageImmediate(exchangeName, routingKey, data, encoding)
		if err != nil {
			Config.GetLogger().Error("error executing publish message", zap.Error(err))
		}
	}()

}
func (r *RabbitMQBroker) publishMessageImmediate(exchangeName string, routingKey string, data []byte, encoding string) error {
	channel, err := r.GetRabbitMQConnection().Channel()
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

func (r *RabbitMQBroker) CreateExchange(name string, exchangeType string) error {
	channel, err := r.GetRabbitMQConnection().Channel()
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

func (r *RabbitMQBroker) GetOrCreateQueue(name string) (amqp.Queue, error) {
	channel, err := r.GetRabbitMQConnection().Channel()
	if err != nil {
		Config.GetLogger().Error("error creating the channel", zap.Error(err))
	}

	return channel.QueueDeclarePassive(name, true, false, false, false, nil)
}
