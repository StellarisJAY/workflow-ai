package repo

import (
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type txKey struct{}

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

func (r *Repository) DB(ctx context.Context) *gorm.DB {
	val := ctx.Value(txKey{})
	if val == nil {
		return r.db
	}
	return val.(*gorm.DB)
}

type TransactionManager struct {
	repo *Repository
}

func NewTransactionManager(repo *Repository) *TransactionManager {
	return &TransactionManager{repo: repo}
}

func (tm *TransactionManager) Tx(ctx context.Context, fn func(c context.Context) error) error {
	return tm.repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txKey{}, tx))
	})
}
