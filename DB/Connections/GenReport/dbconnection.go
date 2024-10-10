package GenReport

import (
	"database/sql"
	"genreport/Startup/Config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type SelfDbConnection struct {
	dbConnection *gorm.DB
	db           *sql.DB
}

func (s *SelfDbConnection) InitDBConnection() error {

	dbConnection, err := gorm.Open(postgres.Open(Config.GetSettings().ConnectionString), &gorm.Config{})
	if err != nil {
		Config.GetLogger().Error("error initiating connection to the DB", zap.Error(err))
		return err
	}
	s.db, err = dbConnection.DB()
	s.db.SetConnMaxLifetime(time.Duration(Config.GetSettings().MaxConnectionTime) * time.Second)
	s.db.SetMaxIdleConns(Config.GetSettings().MaxIdleConnection)
	s.db.SetMaxOpenConns(Config.GetSettings().MaxAllowedConnections)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	return err
}

func (s *SelfDbConnection) CheckDBConnection() error {
	for {
		err := s.db.Ping()
		if err != nil {
			Config.GetLogger().Error("error pinging the DB", zap.Error(err))
			Config.GetLogger().Info("Reconnecting to the database")
			innerErr := s.InitDBConnection()
			if innerErr != nil {
				Config.GetLogger().Error("error pinging the DB", zap.Error(err))
				return innerErr
			}
		}
	}
}
