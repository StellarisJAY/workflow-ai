package common

import "gorm.io/gorm"

type Pagination struct {
	Page     int
	PageSize int
	Total    int
	Paged    bool
}

type PageQuery struct {
	Page     int  `json:"page" form:"page"`
	PageSize int  `json:"pageSize" form:"pageSize"`
	Paged    bool `json:"paged" form:"paged"`
}

func WithPagination(p *Pagination) func(db *gorm.DB) *gorm.DB {
	limit := p.PageSize
	offset := (p.Page - 1) * p.PageSize
	return func(db *gorm.DB) *gorm.DB {
		var total int64
		db.Count(&total)
		p.Total = int(total)
		if p.Paged {
			return db.Limit(limit).Offset(offset)
		} else {
			return db
		}
	}
}
