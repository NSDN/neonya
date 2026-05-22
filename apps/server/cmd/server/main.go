package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NSDN/neonya/apps/server/internal/auth"
	"github.com/NSDN/neonya/apps/server/internal/config"
	"github.com/NSDN/neonya/apps/server/internal/plate"
	"github.com/NSDN/neonya/apps/server/internal/post"
	"github.com/NSDN/neonya/apps/server/internal/shared"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	err := godotenv.Load(
		config.ENV_FILE,
		config.ENV_FILE_DATABASE,
	)

	if err != nil {
		log.Fatal(shared.Messages.EnvironmentErrorNotFound)
	}
}

func mustGetenv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatal(shared.Messages.EnvironmentErrorNeedOne(key))
	}

	return value
}

func setupRouter(db *gorm.DB, tokenKey string, frontendDist string) *gin.Engine {
	router := gin.Default()

	router.Use(shared.CORSMiddleware())
	router.Use(shared.ValidateAuth(tokenKey))

	router.NoRoute(shared.SPANoRoute(frontendDist))

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})

	api := router.Group("/api")

	auth.RegisterRoutes(api, db, tokenKey)
	plate.RegisterRoutes(api, db)
	post.RegisterRoutes(api, db)

	return router
}

func openDatabase() *gorm.DB {
	db, err := shared.OpenDatabaseSimple(
		mustGetenv("POSTGRES_DB"),
		mustGetenv("POSTGRES_USER"),
		mustGetenv("POSTGRES_PASSWORD"),
	)

	if err != nil {
		log.Fatalf("Database open error: %v", err)
	}

	return db
}

func main() {
	db := openDatabase()
	tokenKey := mustGetenv(config.ENV_TOKEN_KEY)
	port := mustGetenv(config.ENV_APPLICATION_PORT)
	frontendDist := os.Getenv(config.ENV_FRONTEND_DIST)

	router := setupRouter(db, tokenKey, frontendDist)

	address := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", address)
	router.Run(address)
}
