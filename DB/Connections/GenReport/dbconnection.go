package GenReport

import (
	"database/sql"
	"genreport/Startup/Models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var dbConnection *gorm.DB
var db *sql.DB

func InitDBConnection() error {

	dbConnection, err := gorm.Open(postgres.Open(Models.Settings{}.ConnectionString), &gorm.Config{})
	db, err = dbConnection.DB()
	db.SetConnMaxLifetime(time.Duration(Models.Settings{}.MaxConnectionTime) * time.Second)
	db.SetMaxIdleConns(Models.Settings{}.MaxIdleConnection)
	db.SetMaxOpenConns(Models.Settings{}.MaxAllowedConnections)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	return err
}

func CheckDBConnection() error {
	for {
		err := db.Ping()
		if err != nil {
			innerErr := InitDBConnection()
			if innerErr != nil {
				return innerErr
			}
		}
	}
}
