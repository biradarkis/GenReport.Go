package main

import (
	"genreport/DB/Connections/GenReportDB"
	"go.uber.org/dig"
)

func main() {
	//Dependecy Injection
	container := dig.New()
	err := container.Provide(GenReportDB.SelfDbConnection{}, GenReportDB.NewDBConnection())
	if err != nil {
		return
	}
}
