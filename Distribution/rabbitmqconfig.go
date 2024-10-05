package Distribution

import (
	"genreport/Startup/Config"
	"genreport/Startup/Models"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var Broker *amqp.Connection
var Settings *Models.Settings
var Logger *zap.Logger

func InitBroker() {
	Logger = Config.GetLogger()
	settings, err := Config.GetSettings()
	if err != nil {
		Logger.Error("error running the application cannot get settings")
		return
	}
	Settings = settings

	Broker, err = amqp.Dial(Settings.AMQPServerURL)
	if err != nil {
		Logger.Error("error running the application cannot connect to RabbitMQ")
		return
	}

}

// CheckConnection pings RabbitMQ server to check if connection is still active
func CheckConnection(connection *amqp.Connection) bool {
	return !(connection == nil || connection.IsClosed())
}
