package repo

import (
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"github.com/StellrisJAY/workflow-ai/internal/model"
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

func (r *Repository) MigrateDB() {
	migrator := r.db.Migrator()
	tables := []any{&model.Provider{}, &model.ProviderModel{}, &model.Template{}, &model.WorkflowInstance{}, &model.NodeInstance{},
		&model.KnowledgeBase{}, &model.KnowledgeBaseFile{}, &model.User{}, &model.KbFileProcessTask{},
		&model.KbFileProcessOptions{}, &model.KbFileChunk{}, &model.File{}}
	for _, table := range tables {
		if !migrator.HasTable(table) {
			if err := migrator.CreateTable(table); err != nil {
				panic(err)
			}
		}
	}
}
