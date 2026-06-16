package thread

import (
	"errors"
	"net/http"

	"github.com/NSDN/neonya/apps/server/internal/shared"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleGetThreads(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		boardID := context.Query("boardId")

		list, err := GetThreads(db, boardID)

		if err != nil {
			shared.HandleRequestError(context, http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, list)
	}
}

func HandleCreateThread(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		var request NewThreadRequest

		if err := context.ShouldBindJSON(&request); err != nil {
			shared.HandleRequestError(
				context,
				http.StatusBadRequest,
				errors.New(shared.Messages.ThreadFailedBadContent),
			)
			return
		}

		thread, err := CreateThread(db, &request)

		if err != nil {
			shared.HandleRequestError(context, http.StatusForbidden, err)
			return
		}

		context.JSON(http.StatusOK, thread)
	}
}

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	api.GET("/threads", HandleGetThreads(db))
	api.POST("/threads", HandleCreateThread(db))
}
