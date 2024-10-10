package Config

import (
	"genreport/Startup/Models"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"sync"
)

var once = sync.Once{}

var Settings *Models.Settings
var Logger *zap.Logger

func ConfigureSettings() error {

	viper.SetConfigFile("../../Config/settings.json")
	err := viper.ReadInConfig() // Read the config file
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&Settings) // Unmarshal config into the struct

	return err
}

// InitLogger this function requires ConfigureSettings Function to be called beforehand
func InitLogger() {
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       strings.ToLower(Settings.Environment) == "development",
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",

		OutputPaths:      nil,
		ErrorOutputPaths: nil,
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
		EncoderConfig: zapcore.EncoderConfig{TimeKey: "timestamp"}}
	Logger = zap.Must(config.Build())
	Logger.Info("Initiated logger", zap.String("environment", Settings.Environment))

}

func GetSettings() *Models.Settings {

	once.Do(func() {
		err := ConfigureSettings()
		if err != nil {
			panic(err)
		}
	})

	return Settings
}

func GetLogger() *zap.Logger {
	once.Do(func() {
		InitLogger()
	})
	return Logger
}
