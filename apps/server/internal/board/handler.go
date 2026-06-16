package board

import (
	"net/http"

	"github.com/NSDN/neonya/apps/server/internal/shared"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleGetBoards(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		list, err := GetBoards(db)

		if err != nil {
			shared.HandleRequestError(context, http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, list)
	}
}

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	api.GET("/boards", HandleGetBoards(db))
}
