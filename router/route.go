package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/letenk/use_deal_user/handler"
	"github.com/letenk/use_deal_user/middleware"
	"github.com/letenk/use_deal_user/repository"
	"github.com/letenk/use_deal_user/service"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Database) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	repositoryUser := repository.NewUserRepository(db)
	serviceUser := service.NewServiceUser(repositoryUser)
	handlerUser := handler.NewHandlerUser(serviceUser)

	// Route home
	router.GET("/", func(c *gin.Context) {
		resp := gin.H{"say": "Server is healthy 💪"}
		c.JSON(http.StatusOK, resp)
	})

	// Group api version 1
	v1 := router.Group("/api/v1")

	// Login
	v1.POST("/login", handlerUser.Login)

	// Group users
	users := v1.Group("/users")
	// Endpoint get all user with middleware
	users.GET("", middleware.AuthMiddleware(serviceUser), handlerUser.GetAll)
	// Endpoint create user with middleware
	users.POST("", middleware.AuthMiddleware(serviceUser), handlerUser.Create)
	// Endpoint get one user with middleware
	users.GET("/:id", middleware.AuthMiddleware(serviceUser), handlerUser.GetOne)
	// Endpoint update user with middleware
	users.PATCH("/:id", middleware.AuthMiddleware(serviceUser), handlerUser.Update)
	// Endpoint delete user with middleware
	users.DELETE("/:id", middleware.AuthMiddleware(serviceUser), handlerUser.Delete)

	return router
}
