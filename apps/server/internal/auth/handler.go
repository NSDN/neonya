package auth

import (
	"errors"
	"net/http"

	"github.com/NSDN/neonya/apps/server/internal/config"
	"github.com/NSDN/neonya/apps/server/internal/shared"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func HandleRegister(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		var info RegisterInfo

		if err := context.ShouldBindJSON(&info); err != nil {
			shared.HandleRequestError(
				context,
				http.StatusBadRequest,
				errors.New(shared.Messages.AuthorizeFailedBadParameter),
			)
			return
		}

		if info.Username == "" {
			shared.HandleRequestError(
				context,
				http.StatusBadRequest,
				errors.New(shared.Messages.AuthorizeFailedMissingParameter("用户名")),
			)
			return
		}

		if info.Password == "" {
			shared.HandleRequestError(
				context,
				http.StatusBadRequest,
				errors.New(shared.Messages.AuthorizeFailedMissingParameter("密码")),
			)
			return
		}

		if info.Password != info.ConfirmPassword {
			shared.HandleRequestError(
				context,
				http.StatusBadRequest,
				errors.New(shared.Messages.AuthorizeFailedWrongPassword),
			)
			return
		}

		result, err := Register(db, &info)

		if err != nil {
			shared.HandleRequestError(context, http.StatusForbidden, err)
			return
		}

		context.JSON(http.StatusOK, result)
	}
}

func HandleLogin(db *gorm.DB, tokenKey string) gin.HandlerFunc {
	return func(context *gin.Context) {
		var info LoginInfo

		if err := context.ShouldBindJSON(&info); err != nil {
			shared.HandleRequestError(
				context,
				http.StatusBadRequest,
				errors.New(shared.Messages.AuthorizeFailedBadParameter),
			)
			return
		}

		if info.Username == "" {
			shared.HandleRequestError(
				context,
				http.StatusBadRequest,
				errors.New(shared.Messages.AuthorizeFailedMissingParameter("用户名")),
			)
			return
		}

		if info.Password == "" {
			shared.HandleRequestError(
				context,
				http.StatusBadRequest,
				errors.New(shared.Messages.AuthorizeFailedMissingParameter("密码")),
			)
			return
		}

		token, err := Login(db, tokenKey, info)

		if err != nil {
			shared.HandleRequestError(context, http.StatusForbidden, err)
			return
		}

		context.JSON(http.StatusOK, Token{
			AccessToken: token,
		})
	}
}

func HandleGetUserInfo(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, exist := context.Get(config.CONTEXT_KEY_CLAIMS)

		if !exist {
			shared.HandleWrongTokenError(context)
			return
		}

		uid := claims.(jwt.MapClaims)["uid"].(string)

		user, err := FindUserByUID(db, uid)

		if err != nil {
			shared.HandleRequestError(context, http.StatusForbidden, err)
			return
		}

		context.JSON(http.StatusOK, user)
	}
}

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB, tokenKey string) {
	api.POST(config.API_REGISTER, HandleRegister(db))
	api.POST(config.API_LOGIN, HandleLogin(db, tokenKey))
	api.GET(config.API_GET_USER_INFO, HandleGetUserInfo(db))
}
