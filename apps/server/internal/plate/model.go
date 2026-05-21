package plate

import "github.com/NSDN/neonya/apps/server/internal/config"

type Plate struct {
	ID          string         `json:"id" gorm:"size:100;primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"not null"`
	Background  string         `json:"background" gorm:"not null"`
	PageType    config.PageType `json:"pageType" gorm:"size:20;not null"`
	SortOrder   int            `json:"sortOrder" gorm:"unique;not null;check:(sort_order >= 0)"`
}

var DefaultPlates = []Plate{
	{
		ID:          "localization",
		Name:        "喵玉汉化馆",
		Description: "",
		Background:  "https://i.imgur.com/ohQuzivl.jpg",
		PageType:    config.COMIC,
		SortOrder:   0,
	},
	{
		ID:          "music",
		Name:        "喵玉咏唱组",
		Description: "",
		Background:  "https://i.imgur.com/IHo7tTyl.jpg",
		PageType:    config.ARTICLE,
		SortOrder:   1,
	},
	{
		ID:          "chat",
		Name:        "魔女的茶会",
		Description: "",
		Background:  "https://i.imgur.com/JsWkJ4jl.jpg",
		PageType:    config.ARTICLE,
		SortOrder:   2,
	},
}
