package repo

import (
	"context"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"gorm.io/gorm"
)

type LLMRepo struct {
	*Repository
}

var llmTableName = "wf_llm"

func NewLLMRepo(repo *Repository) *LLMRepo {
	return &LLMRepo{repo}
}

func (lr *LLMRepo) Insert(ctx context.Context, llm *model.LLM) error {
	return lr.db.Table(llmTableName).WithContext(ctx).Create(llm).Error
}

func (lr *LLMRepo) GetDetail(ctx context.Context, id int64) (*model.LLMDetailDTO, error) {
	var llm *model.LLMDetailDTO
	err := lr.db.Table(llmTableName+" llm").
		Joins("INNER JOIN wf_user u ON u.user_id = llm.add_user").
		Select("llm.*, u.username AS add_username").
		Where("llm.id =?", id).
		WithContext(ctx).
		Scan(&llm).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return llm, err
}

func (lr *LLMRepo) List(ctx context.Context, query *model.LLMQuery) ([]*model.LLMListDTO, error) {
	tx := lr.db.Table(llmTableName + " llm").
		Joins("INNER JOIN wf_user u ON u.user_id = llm.add_user").
		Select("llm.*, u.username AS add_username").
		WithContext(ctx)
	if query.Code != "" {
		tx.Where("llm.code =?", query.Code)
	}
	if query.Name != "" {
		tx.Where("llm.name like?", "%"+query.Name+"%")
	}
	var res []*model.LLMListDTO
	err := tx.Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
