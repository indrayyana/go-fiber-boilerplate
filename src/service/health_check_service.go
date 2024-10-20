package service

import (
	"app/src/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HealthCheckService interface {
	GormCheck() error
}

type healthCheckService struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewHealthCheckService(db *gorm.DB) HealthCheckService {
	return &healthCheckService{
		Log: utils.Log,
		DB:  db,
	}
}

func (s *healthCheckService) GormCheck() error {
	sqlDB, errDB := s.DB.DB()
	if errDB != nil {
		s.Log.Errorf("failed to access the database connection pool: %v", errDB)
		return errDB
	}

	if err := sqlDB.Ping(); err != nil {
		s.Log.Errorf("failed to ping the database: %v", err)
		return err
	}

	return nil
}
