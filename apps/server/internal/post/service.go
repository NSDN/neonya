package post

import (
	"strconv"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func GetTopics(db *gorm.DB, plateID string) ([]TopicListItem, error) {
	var topics []Topic

	result := db.
		Where("plate_id = ?", plateID).
		Select("id", "title", "thumbnail_link", "updated_at").
		Order("updated_at DESC").
		Find(&topics)

	if result.Error != nil {
		return nil, result.Error
	}

	items := make([]TopicListItem, 0, len(topics))

	for _, topic := range topics {
		items = append(items, TopicListItem{
			ID:            strconv.FormatInt(topic.ID, 10),
			Title:         topic.Title,
			ThumbnailLink: topic.ThumbnailLink,
			UpdatedAt:     topic.UpdatedAt,
		})
	}

	return items, nil
}

func CreateTopic(db *gorm.DB, request *NewTopicRequest) (*Topic, error) {
	now := time.Now()

	topic := Topic{
		AuthorID:      request.Author.UID,
		PlateID:       request.Plate,
		Title:         request.Title,
		TopicType:     request.BodyType,
		ThumbnailLink: "",
		Tag:           pq.StringArray(request.Tag),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	result := db.Create(&topic)

	if result.Error != nil {
		return nil, result.Error
	}

	return &topic, nil
}
