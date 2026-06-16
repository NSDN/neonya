package thread

import (
	"strconv"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func GetThreads(db *gorm.DB, boardID string) ([]ThreadListItem, error) {
	var threads []Thread

	result := db.
		Where("board_id = ?", boardID).
		Select("id", "title", "thumbnail_link", "updated_at").
		Order("updated_at DESC").
		Find(&threads)

	if result.Error != nil {
		return nil, result.Error
	}

	items := make([]ThreadListItem, 0, len(threads))

	for _, thread := range threads {
		items = append(items, ThreadListItem{
			ID:            strconv.FormatInt(thread.ID, 10),
			Title:         thread.Title,
			ThumbnailLink: thread.ThumbnailLink,
			UpdatedAt:     thread.UpdatedAt,
		})
	}

	return items, nil
}

func CreateThread(db *gorm.DB, request *NewThreadRequest) (*Thread, error) {
	now := time.Now()

	thread := Thread{
		AuthorID:      request.Author.UID,
		BoardID:       request.Board,
		Title:         request.Title,
		ThreadType:    request.BodyType,
		ThumbnailLink: "",
		Tag:           pq.StringArray(request.Tag),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	result := db.Create(&thread)

	if result.Error != nil {
		return nil, result.Error
	}

	return &thread, nil
}
