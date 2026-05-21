package post

import (
	"errors"
	"net/http"

	"github.com/NSDN/neonya/apps/server/internal/shared"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleGetTopics(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		plateID := context.Query("plateId")

		list, err := GetTopics(db, plateID)

		if err != nil {
			shared.HandleRequestError(context, http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, list)
	}
}

func HandleCreateTopic(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		var request NewTopicRequest

		if err := context.ShouldBindJSON(&request); err != nil {
			shared.HandleRequestError(
				context,
				http.StatusBadRequest,
				errors.New(shared.Messages.ArticleFailedBadContent),
			)
			return
		}

		topic, err := CreateTopic(db, &request)

		if err != nil {
			shared.HandleRequestError(context, http.StatusForbidden, err)
			return
		}

		context.JSON(http.StatusOK, topic)
	}
}

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	api.GET("/topics", HandleGetTopics(db))
	api.POST("/topic", HandleCreateTopic(db))
}
