package shared

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleRequestError(context *gin.Context, httpStatus int, err error) {
	log.Println(err)
	context.String(httpStatus, err.Error())
}

func HandleWrongTokenError(context *gin.Context) {
	HandleRequestError(
		context,
		http.StatusUnauthorized,
		fmt.Errorf(Messages.AuthorizeFailedWrongToken),
	)
	context.Abort()
}

func FilterSlice[T any](slice *[]T, predicate func(item *T, index int) bool) *[]T {
	var result []T

	for index, item := range *slice {
		if predicate(&item, index) {
			result = append(result, item)
		}
	}

	return &result
}
