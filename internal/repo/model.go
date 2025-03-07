package repo

import (
	"context"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"gorm.io/gorm"
)

type ModelRepo struct {
	*Repository
}

func NewModelRepo(repo *Repository) *ModelRepo {
	return &ModelRepo{repo}
}

func (mr *ModelRepo) Insert(ctx context.Context, m *model.Model) error {
	return mr.DB(ctx).Table(m.TableName()).WithContext(ctx).Create(m).Error
}

func (mr *ModelRepo) GetDetail(ctx context.Context, id int64) (*model.ModelDetailDTO, error) {
	var m *model.ModelDetailDTO
	err := mr.DB(ctx).Table(model.Model{}.TableName()+" m").
		Joins("INNER JOIN wf_user u ON u.user_id = m.add_user").
		Select("m.*, u.username AS add_username").
		Where("m.id =?", id).
		WithContext(ctx).
		Scan(&m).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return m, err
}

func (mr *ModelRepo) List(ctx context.Context, query *model.ModelQuery) ([]*model.ModelListDTO, int, error) {
	p := common.Pagination{Page: query.Page, PageSize: query.PageSize, Paged: query.Paged}
	tx := mr.DB(ctx).Table(model.Model{}.TableName() + " m").
		Joins("INNER JOIN wf_user u ON u.user_id = m.add_user").
		Select("m.*, u.username AS add_username").
		Scopes(common.WithPagination(&p)).
		WithContext(ctx)
	if query.Code != "" {
		tx = tx.Where("m.code =?", query.Code)
	}
	if query.Name != "" {
		tx = tx.Where("m.name like?", "%"+query.Name+"%")
	}
	if query.ModelType != "" {
		tx = tx.Where("m.model_type = ?", query.ModelType)
	}
	var res []*model.ModelListDTO
	err := tx.Scan(&res).Error
	if err != nil {
		return nil, 0, err
	}
	return res, p.Total, nil
}
