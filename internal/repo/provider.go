package repo

import (
	"context"
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
)

type ProviderRepo struct {
	*Repository
}

func NewProviderRepo(repo *Repository) *ProviderRepo {
	return &ProviderRepo{repo}
}

func (m *ProviderRepo) InsertProvider(ctx context.Context, provider *model.Provider) error {
	return m.DB(ctx).Table(provider.TableName()).WithContext(ctx).Create(provider).Error
}

func (m *ProviderRepo) InsertProviderModel(ctx context.Context, pm *model.ProviderModel) error {
	return m.DB(ctx).Table(pm.TableName()).WithContext(ctx).Create(pm).Error
}

func (m *ProviderRepo) GetProviderModelList(ctx context.Context, query *model.ProviderModelQuery) ([]*model.ProviderModelListDTO,
	int, error) {
	var result []*model.ProviderModelListDTO
	page := common.Pagination{Page: query.Page, PageSize: query.PageSize, Paged: query.Paged}
	d := m.DB(ctx).Table(model.Provider{}.TableName() + " p").
		Joins("INNER JOIN wf_provider_model pm ON pm.provider_id = p.id").
		Select("pm.*, p.name AS provider_name, p.code AS provider_code, p.credentials AS provider_credentials").
		WithContext(ctx).
		Scopes(common.WithPagination(&page))
	if query.ProviderId != 0 {
		d = d.Where("pm.provider_id =?", query.ProviderId)
	}
	if query.ModelType != "" {
		d = d.Where("pm.model_type =?", query.ModelType)
	}
	if err := d.Scan(&result).Error; err != nil {
		return nil, 0, err
	}
	return result, page.Total, nil
}

func (m *ProviderRepo) GetProviderList(ctx context.Context) ([]*model.ProviderListDTO, error) {
	var result []*model.ProviderListDTO
	err := m.DB(ctx).Table(model.Provider{}.TableName()).Select("id, name, code, add_time").Find(&result).Error
	return result, err
}

func (m *ProviderRepo) GetProviderModelDetail(ctx context.Context, id int64) (*model.ProviderModelDetail, error) {
	var result *model.ProviderModelDetail
	err := m.DB(ctx).Table(model.Provider{}.TableName()+" p").
		Joins("INNER JOIN wf_provider_model pm ON pm.provider_id = p.id").
		Select("pm.*, p.name AS provider_name, p.code AS provider_code, p.credentials AS provider_credentials").
		WithContext(ctx).Where("pm.id=?", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *ProviderRepo) GetProviderModel(ctx context.Context, id int64) (*model.ProviderModel, error) {
	var result *model.ProviderModel
	err := m.DB(ctx).Table(model.ProviderModel{}.TableName()).Where("id=?", id).First(&result).Error
	return result, err
}

func (m *ProviderRepo) GetProvider(ctx context.Context, id int64) (*model.Provider, error) {
	var result *model.Provider
	err := m.DB(ctx).Table(model.Provider{}.TableName()).Where("id=?", id).First(&result).Error
	return result, err
}
