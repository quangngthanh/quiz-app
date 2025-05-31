package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"quiz-app/internal/config"
	"quiz-app/internal/handler"
	"quiz-app/internal/model"
	"quiz-app/internal/repository"
	"quiz-app/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	env := flag.String("env", "", "Environment: prod, dev, local")
	flag.Parse()
	if *env == "" {
		*env = os.Getenv("APP_ENV")
	}
	cfg, err := config.LoadConfig(*env)
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password,
		cfg.Database.Name, cfg.Database.Port, cfg.Database.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	db.AutoMigrate(&model.User{}, &model.QuizSession{}, &model.Question{}, &model.UserAnswer{})

	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Initialize repositories
	quizRepo := repository.NewQuizRepository(db)
	redisRepo := repository.NewRedisRepository(rdb)

	// Initialize services
	wsService := service.NewWebSocketService(redisRepo)
	quizService := service.NewQuizService(quizRepo, redisRepo, wsService)

	// Initialize handlers
	quizHandler := handler.NewQuizHandler(quizService)
	wsHandler := handler.NewWebSocketHandler(wsService)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API routes
	api := r.Group("/api")
	{
		api.POST("/quiz", quizHandler.CreateQuiz)
		api.GET("/quiz/:quiz_id", quizHandler.GetQuiz)
		api.POST("/quiz/:quiz_id/join", quizHandler.JoinQuiz)
		api.POST("/quiz/:quiz_id/answer", quizHandler.SubmitAnswer)
		api.GET("/quiz/:quiz_id/leaderboard", quizHandler.GetLeaderboard)
	}

	// WebSocket routes
	r.GET("/ws/quiz/:quiz_id/leaderboard", wsHandler.HandleLeaderboardWebSocket)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(r.Run(":" + cfg.Server.Port))
}
