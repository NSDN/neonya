package board

import "github.com/NSDN/neonya/apps/server/internal/config"

type Board struct {
	ID          string         `json:"id" gorm:"size:100;primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"not null"`
	Background  string         `json:"background" gorm:"not null"`
	PageType    config.PageType `json:"pageType" gorm:"size:20;not null"`
	SortOrder   int            `json:"sortOrder" gorm:"unique;not null;check:(sort_order >= 0)"`
}


