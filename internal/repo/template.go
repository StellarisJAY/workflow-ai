package repo

import (
	"context"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"gorm.io/gorm"
)

type TemplateRepo struct {
	*Repository
}

func NewTemplateRepo(repo *Repository) *TemplateRepo {
	return &TemplateRepo{repo}
}

func (tr *TemplateRepo) Insert(ctx context.Context, template *model.Template) error {
	return tr.db.Table(template.TableName()).WithContext(ctx).Create(template).Error
}

func (tr *TemplateRepo) GetDetail(ctx context.Context, id int64) (*model.TemplateDetailDTO, error) {
	var template *model.TemplateDetailDTO
	err := tr.DB(ctx).Table(model.Template{}.TableName()+" t").
		Joins("INNER JOIN wf_user u ON u.user_id = t.add_user").
		Select("t.*, u.username AS add_username").
		Where("t.id =?", id).
		WithContext(ctx).
		Scan(&template).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return template, err
}

func (tr *TemplateRepo) List(ctx context.Context, query *model.TemplateQuery) ([]*model.TemplateListDTO, int, error) {
	p := common.Pagination{Page: query.Page, PageSize: query.PageSize, Paged: query.Paged}
	var templates []*model.TemplateListDTO
	tx := tr.DB(ctx).Table(model.Template{}.TableName() + " t").
		Joins("INNER JOIN wf_user u ON u.user_id = t.add_user").
		Select("t.*, u.username AS add_username").
		Scopes(common.WithPagination(&p)).
		WithContext(ctx)
	if query.Name != "" {
		tx.Where("t.name like ?", "%"+query.Name+"%")
	}
	err := tx.Scan(&templates).Error
	if err != nil {
		return nil, 0, err
	}
	return templates, p.Total, nil
}

func (tr *TemplateRepo) Delete(ctx context.Context, id int64) error {
	return tr.db.Table(model.Template{}.TableName()).
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Template{}).
		Error
}

func (tr *TemplateRepo) Update(ctx context.Context, template *model.Template) error {
	return tr.db.Table(template.TableName()).
		WithContext(ctx).
		Where("id = ?", template.Id).
		UpdateColumns(map[string]interface{}{
			"name":        template.Name,
			"description": template.Description,
			"data":        template.Data,
		}).Error
}
