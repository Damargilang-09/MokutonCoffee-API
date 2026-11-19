package common

import "gorm.io/gorm"

type PaginationQuery struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`
}

func (p *PaginationQuery) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// minimal 1 data
		if p.Page < 1 {
			p.Page = 1
		}

		// membatasi limit
		switch {
		case p.Limit > 100:
			p.Limit = 100
		case p.Limit <= 0:
			p.Limit = 10
		}

		offset := (p.Page - 1) * p.Limit
		return db.Offset(offset).Limit(p.Limit)
	}
}
