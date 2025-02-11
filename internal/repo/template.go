package repo

import (
	"context"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"gorm.io/gorm"
)

type TemplateRepo struct {
	*Repository
}

var templateTableName = "wf_template"

func NewTemplateRepo(repo *Repository) *TemplateRepo {
	return &TemplateRepo{repo}
}

func (tr *TemplateRepo) Insert(ctx context.Context, template *model.Template) error {
	return tr.db.Table(templateTableName).WithContext(ctx).Create(template).Error
}

func (tr *TemplateRepo) GetDetail(ctx context.Context, id int64) (*model.TemplateDetailDTO, error) {
	var template *model.TemplateDetailDTO
	err := tr.db.Table(templateTableName+" t").
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

func (tr *TemplateRepo) List(ctx context.Context, query *model.TemplateQuery) ([]*model.TemplateListDTO, error) {
	var templates []*model.TemplateListDTO
	tx := tr.db.Table(templateTableName + " t").
		Joins("INNER JOIN wf_user u ON u.user_id = t.add_user").
		Select("t.*, u.username AS add_username").
		WithContext(ctx)
	if query.Name != "" {
		tx.Where("t.name like ?", "%"+query.Name+"%")
	}
	err := tx.Scan(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}
