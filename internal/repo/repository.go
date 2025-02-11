package repo

import (
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(conf *config.Config) (*Repository, error) {
	db, err := gorm.Open(mysql.Open(conf.Database.Url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}
