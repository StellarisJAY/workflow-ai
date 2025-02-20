package repo

import (
	"context"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"gorm.io/gorm"
)

type LLMRepo struct {
	*Repository
}

func NewLLMRepo(repo *Repository) *LLMRepo {
	return &LLMRepo{repo}
}

func (lr *LLMRepo) Insert(ctx context.Context, llm *model.LLM) error {
	return lr.DB(ctx).Table(llm.TableName()).WithContext(ctx).Create(llm).Error
}

func (lr *LLMRepo) GetDetail(ctx context.Context, id int64) (*model.LLMDetailDTO, error) {
	var llm *model.LLMDetailDTO
	err := lr.DB(ctx).Table(model.LLM{}.TableName()+" llm").
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

func (lr *LLMRepo) List(ctx context.Context, query *model.LLMQuery) ([]*model.LLMListDTO, int, error) {
	p := common.Pagination{Page: query.Page, PageSize: query.PageSize, Paged: query.Paged}
	tx := lr.DB(ctx).Table(model.LLM{}.TableName() + " llm").
		Joins("INNER JOIN wf_user u ON u.user_id = llm.add_user").
		Select("llm.*, u.username AS add_username").
		Scopes(common.WithPagination(&p)).
		WithContext(ctx)
	if query.Code != "" {
		tx = tx.Where("llm.code =?", query.Code)
	}
	if query.Name != "" {
		tx = tx.Where("llm.name like?", "%"+query.Name+"%")
	}
	if query.ModelType != "" {
		tx = tx.Where("llm.model_type = ?", query.ModelType)
	}
	var res []*model.LLMListDTO
	err := tx.Scan(&res).Error
	if err != nil {
		return nil, 0, err
	}
	return res, p.Total, nil
}
