package main

import (
	"genreport/Startup/Models"
	"github.com/spf13/viper"
)

var Settings Models.Settings

func ConfigureSettings() error {

	viper.SetConfigFile("../../Config/settings.json")
	err := viper.ReadInConfig() // Read the config file
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&Settings) // Unmarshal config into the struct

	return err
}
