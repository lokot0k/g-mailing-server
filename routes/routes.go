package routes

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lokot0k/mservice/config"
	"github.com/lokot0k/mservice/controllers"
	"github.com/lokot0k/mservice/database"
	"github.com/lokot0k/mservice/middleware"
	"github.com/lokot0k/mservice/queue"
)

func SetupRouter(cfg *config.Config) (*gin.Engine, error) {
	db, err := database.Connect(cfg)
	if err != nil {
		return nil, err
	}
	if err := database.Migrate(cfg); err != nil {
		return nil, err
	}

	rmq, err := queue.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq connect: %w", err)
	}

	auth := &controllers.AuthController{DB: db, Config: cfg}
	msg := &controllers.MessageController{DB: db, RMQ: rmq}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/register", auth.Register)
		v1.POST("/auth/login", auth.Login)

		sec := v1.Group("/")
		sec.Use(middleware.AuthMiddleware(cfg.JwtSecret))
		sec.GET("/users/me", func(c *gin.Context) {
			c.JSON(200, gin.H{"user_id": c.MustGet("user_id")})
		})

		sec.POST("/messages/send", msg.Send)
		sec.GET("/messages/inbox", msg.Inbox)
		sec.GET("/messages/sent", msg.Sent)
	}

	return r, nil
}
