package shared

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/NSDN/neonya/apps/server/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NotFound(context *gin.Context) {
	HandleRequestError(
		context,
		http.StatusNotFound,
		fmt.Errorf("请求的路径 %s 不存在", context.Request.URL.Path),
	)
}

func SPANoRoute(indexPath string) gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path

		if strings.HasPrefix(path, "/api") || indexPath == "" {
			NotFound(context)
			return
		}

		context.File(indexPath)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		context.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
			return
		}

		context.Next()
	}
}

func ValidateAuth(tokenKey string) gin.HandlerFunc {
	return func(context *gin.Context) {
		authorization := context.GetHeader(config.HTTP_HEADER_AUTHORIZATION)

		if authorization == "" {
			context.Next()
			return
		}

		tokenString := extractBearerToken(authorization)

		if tokenString == "" {
			context.Next()
			return
		}

		token, err := parseToken(tokenString, tokenKey)

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			context.Set(config.CONTEXT_KEY_CLAIMS, claims)
			context.Next()
		} else {
			log.Println(err)
			HandleWrongTokenError(context)
		}
	}
}

func extractBearerToken(authorization string) string {
	result, found := strings.CutPrefix(authorization, config.AUTHENTICATION_TYPE)

	if !found {
		return ""
	}

	return strings.TrimSpace(result)
}

func parseToken(tokenString string, tokenKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenKey), nil
	})
}
