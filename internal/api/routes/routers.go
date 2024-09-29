package routes

import (
	"NoteGolang/internal/middleware"
	"NoteGolang/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	Router     *gin.Engine
	Services   *service.Services
	Middleware *middleware.Middleware
	// middleware *middleware.Middleware
}

type RoutesGroup struct {
	PublicRouter   *gin.RouterGroup
	InternalRouter *gin.RouterGroup
}

func SetupRoutes(r *Routes) {
	// r.Router.Use(cors.Default())
	// r.Router.Use(r.Middleware.CORSMiddleware())
	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Router.Use(cors.New(config))

	// 定义公共路由组
	publicRouter := r.Router.Group("/v1")
	// 定义內部路由组
	internalRouter := r.Router.Group("/v2", r.Middleware.InternalMiddleware)

	r.SetupArticleRoutes(&RoutesGroup{PublicRouter: publicRouter, InternalRouter: internalRouter})
}
