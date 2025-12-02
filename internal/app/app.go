package app

import (
	"os"
	"log"
	_ "sync"

	"tyto/internal/service"
	"tyto/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/gzip"
)

type Application struct {
	service service.ServiceInterface
	routes *gin.Engine
	WebhookChan chan struct{}
	WebhookSecret string
}

func NewApplication() *Application {
	s := store.NewStore()
	// 加载 .env 文件
	_ = godotenv.Load()

	service := service.NewService(
		s,
		mustGetEnv("GIT_REPO_URL"),
		mustGetEnv("REPOSITORY_DIR"),
	)

	app := &Application{
		service: service,
		routes: gin.Default(),
		WebhookChan: make(chan struct{}),
		WebhookSecret: mustGetEnv("WEBHOOK_SECRET"),
	}

	// 启动后台任务
	go app.updateCache()
	return app
}

func (app *Application) RegisterRoutes() {
	app.routes.Use(cors.Default())
	app.routes.Use(gzip.Gzip(gzip.DefaultCompression))
	app.routes.GET("/api/healthcheck", app.healthcheckHandler)
	app.routes.GET("/api/categories", app.GetCategoriesHandler)
	app.routes.GET("/api/categoryTree", app.GetCategoryTreeHandler)
	app.routes.GET("/api/file", app.GetContentHandler)
	app.routes.POST("/api/webhook", app.webhookHandler)
}

func (app *Application) Run(addr string) {
	app.routes.Run(addr)
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Fatal Error: Environment variable %s is not set.", key)
	}
	return value
}