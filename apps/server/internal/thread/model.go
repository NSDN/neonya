package thread

import (
	"time"

	"github.com/NSDN/neonya/apps/server/internal/auth"
	"github.com/lib/pq"
)

type NewThreadRequest struct {
	Author       auth.UserPublicInfo `json:"author"`
	Board        string              `json:"board"`
	Title        string              `json:"title"`
	Tag          []string            `json:"tag"`
	CreationDate time.Time           `json:"creationDate"`
	BodyType     string              `json:"bodyType"`
	Body         string              `json:"body"`
}

type Thread struct {
	ID            int64  `gorm:"primaryKey;autoIncrement"`
	AuthorID      string `gorm:"not null"`
	BoardID       string `gorm:"not null"`
	Title         string `gorm:"not null"`
	ThreadType    string `gorm:"size:20;not null"`
	ThumbnailLink string `gorm:"not null"`
	Tag           pq.StringArray `gorm:"type:varchar(20)[]"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}

type ThreadListItem struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	ThumbnailLink string    `json:"thumbnailLink"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
