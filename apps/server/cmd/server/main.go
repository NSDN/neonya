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

func openDatabase() *gorm.DB {
	dbname := mustGetenv("POSTGRES_DB")
	user := mustGetenv("POSTGRES_USER")
	password := mustGetenv("POSTGRES_PASSWORD")

	db, err := shared.OpenDatabaseSimple(dbname, user, password)

	if err != nil {
		log.Fatalf("Database open error: %v", err)
	}

	return db
}

func setupRouter(db *gorm.DB, tokenKey string) *gin.Engine {
	router := gin.Default()

	router.Use(shared.CORSMiddleware())
	router.Use(shared.ValidateAuth(tokenKey))

	router.NoRoute(shared.NotFound)

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})

	api := router.Group("/api")

	auth.RegisterRoutes(api, db, tokenKey)
	plate.RegisterRoutes(api, db)
	post.RegisterRoutes(api, db)

	return router
}

func main() {
	db := openDatabase()
	tokenKey := mustGetenv(config.ENV_TOKEN_KEY)
	port := mustGetenv(config.ENV_APPLICATION_PORT)

	router := setupRouter(db, tokenKey)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	router.Run(addr)
}
