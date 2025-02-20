package repo

import (
	"context"
	"github.com/StellrisJAY/workflow-ai/internal/common"
	"github.com/StellrisJAY/workflow-ai/internal/model"
)

type KnowledgeBaseRepo struct {
	*Repository
}

func NewKnowledgeBaseRepo(repo *Repository) *KnowledgeBaseRepo {
	return &KnowledgeBaseRepo{repo}
}

func (k *KnowledgeBaseRepo) Insert(ctx context.Context, kb *model.KnowledgeBase) error {
	return k.DB(ctx).Table(kb.TableName()).WithContext(ctx).Create(kb).Error
}

func (k *KnowledgeBaseRepo) Update(ctx context.Context, kb *model.KnowledgeBase) error {
	return k.DB(ctx).Table(kb.TableName()).Where("id = ?", kb.Id).UpdateColumns(map[string]interface{}{
		"name":   kb.Name,
		"enable": kb.Enable,
	}).Error
}

func (k *KnowledgeBaseRepo) InsertFile(ctx context.Context, kbFile *model.KnowledgeBaseFile) error {
	return k.DB(ctx).WithContext(ctx).Create(kbFile).Error
}

func (k *KnowledgeBaseRepo) UpdateFile(ctx context.Context, kbFile *model.KnowledgeBaseFile) error {
	return k.DB(ctx).Table(kbFile.TableName()).WithContext(ctx).Where("id=?", kbFile.Id).
		UpdateColumns(map[string]interface{}{
			"name":   kbFile.Name,
			"status": kbFile.Status,
		}).Error
}

func (k *KnowledgeBaseRepo) GetFileDetail(ctx context.Context, id int64) (*model.KnowledgeBaseFile, error) {
	var kbFile *model.KnowledgeBaseFile
	err := k.DB(ctx).
		WithContext(ctx).
		Where("id =?", id).
		Find(&kbFile).
		Error
	if err != nil {
		return nil, err
	}
	return kbFile, nil
}

func (k *KnowledgeBaseRepo) List(ctx context.Context, query *model.KnowledgeBaseQuery) ([]*model.KnowledgeBaseListDTO, int, error) {
	p := common.Pagination{Page: query.Page, PageSize: query.PageSize, Paged: query.Paged}
	var result []*model.KnowledgeBaseListDTO
	d := k.DB(ctx).Table(model.KnowledgeBase{}.TableName()).
		Select("id, name, description, enable, add_time, add_user").
		Scopes(common.WithPagination(&p)).
		WithContext(ctx)
	if err := d.Scan(&result).Error; err != nil {
		return nil, 0, err
	}
	return result, p.Total, nil
}

func (k *KnowledgeBaseRepo) Detail(ctx context.Context, id int64) (*model.KnowledgeBaseDetailDTO, error) {
	var detail *model.KnowledgeBaseDetailDTO
	if err := k.DB(ctx).Table(model.KnowledgeBase{}.TableName()).
		Where("id = ?", id).
		WithContext(ctx).
		Find(&detail).
		Error; err != nil {
		return nil, err
	}
	return detail, nil
}

func (k *KnowledgeBaseRepo) CountFile(ctx context.Context, kbId int64) (int64, int64, error) {
	var result *model.KbFileCountSize
	if err := k.DB(ctx).Table(model.KnowledgeBaseFile{}.TableName()).
		Select("COUNT(`id`) AS count, IFNULL(SUM(`length`),0) AS size").
		Where("kb_id = ?", kbId).
		Scan(&result).Error; err != nil {
		return 0, 0, err
	}
	return result.Count, result.Size, nil
}

func (k *KnowledgeBaseRepo) ListFileCount(ctx context.Context, kbIds []int64) ([]*model.KbFileCountSize, error) {
	var result []*model.KbFileCountSize
	if err := k.DB(ctx).Table(model.KnowledgeBaseFile{}.TableName()).
		Select("COUNT(`id`) AS count, kb_id, IFNULL(SUM(`length`),0) AS size").
		Where("kb_id IN (?)", kbIds).
		Group("kb_id").
		Find(&result).
		Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (k *KnowledgeBaseRepo) ListFile(ctx context.Context, kbId int64, query *model.KbFileQuery) ([]*model.KbFileListDTO, int, error) {
	var result []*model.KbFileListDTO
	p := common.Pagination{Page: query.Page, PageSize: query.PageSize, Paged: query.Paged}
	d := k.DB(ctx).Table(model.KnowledgeBaseFile{}.TableName()).
		Scopes(common.WithPagination(&p)).
		Where("kb_id = ?", kbId)
	if err := d.Scan(&result).Error; err != nil {
		return nil, 0, err
	}
	return result, p.Total, nil
}
