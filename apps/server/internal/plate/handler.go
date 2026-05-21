package plate

import (
	"net/http"

	"github.com/NSDN/neonya/apps/server/internal/shared"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleGetPlates(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		list, err := GetPlates(db)

		if err != nil {
			shared.HandleRequestError(context, http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, list)
	}
}

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	InitPlateList(db)
	api.GET("/plates", HandleGetPlates(db))
}
