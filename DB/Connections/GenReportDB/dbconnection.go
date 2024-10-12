package GenReportDB

import (
	"context"
	"database/sql"
	"genreport/Startup/Config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type SelfDbConnection struct {
	dbConnection *gorm.DB
	db           *sql.DB
}

func NewDBConnection() (*SelfDbConnection, error) {

	dbConnection, err := gorm.Open(postgres.Open(Config.GetSettings().ConnectionString), &gorm.Config{})
	if err != nil {
		Config.GetLogger().Error("error initiating connection to the DB", zap.Error(err))
		return nil, err
	}
	db, err := dbConnection.DB()
	if err != nil {
		Config.GetLogger().Error("error getting the database connection", zap.Error(err))
		return nil, err
	}
	db.SetConnMaxLifetime(time.Duration(Config.GetSettings().MaxConnectionTime) * time.Second)
	db.SetMaxIdleConns(Config.GetSettings().MaxIdleConnection)
	db.SetMaxOpenConns(Config.GetSettings().MaxAllowedConnections)

	return &SelfDbConnection{
		dbConnection: dbConnection,
		db:           db,
	}, nil
}

func (s *SelfDbConnection) CheckDBConnection(ctx context.Context, connectionTime int64) {
	// Run the connection check logic in a goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Exit the goroutine when the context is canceled
				Config.GetLogger().Info("Stopping DB connection check goroutine")
				return
			default:
				err := s.db.Ping()
				if err != nil {
					Config.GetLogger().Error("error pinging the DB", zap.Error(err))
					Config.GetLogger().Info("Reconnecting to the database")

					newConnection, innerErr := NewDBConnection()

					if innerErr != nil {
						Config.GetLogger().Error("error initializing the DB connection", zap.Error(innerErr))
						return
					}
					s.db = newConnection.db
					s.dbConnection = newConnection.dbConnection
				}
				time.Sleep(time.Second * time.Duration(connectionTime))
			}
		}
	}()

}
