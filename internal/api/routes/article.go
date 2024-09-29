package routes

import (
	"NoteGolang/internal/api/handlers"
	"NoteGolang/internal/middleware"
	"time"
)

func (r *Routes) SetupArticleRoutes(rg *RoutesGroup) {
	handler := handlers.NewHandlers(r.Services)

	{
		rg.PublicRouter.GET("/test", handler.ArticleHandler.Test)
		rg.PublicRouter.GET("/article/descriptions", r.Middleware.DiskCache(middleware.CacheConfig{Duration: time.Hour * 1}), handler.ArticleHandler.GetArticleDescriptions)
		rg.PublicRouter.GET("/article/descriptions/:tag", r.Middleware.DiskCache(middleware.CacheConfig{Duration: time.Hour * 1}), handler.ArticleHandler.GetArticleDescriptionsByTag)
		rg.PublicRouter.GET("/article/tags", r.Middleware.DiskCache(middleware.CacheConfig{Duration: time.Hour * 1}), handler.ArticleHandler.GetArticleTagInfos)
		rg.PublicRouter.GET("/article/:id", r.Middleware.DiskCache(middleware.CacheConfig{Duration: time.Hour * 1}), handler.ArticleHandler.GetArticleByID)
		rg.PublicRouter.GET("/article/:id/comment", handler.ArticleHandler.GetCommentByID)

		rg.PublicRouter.PUT("/article/:id/info/share/count", handler.ArticleHandler.IncrementArticleInfoShareCountByID)
		rg.PublicRouter.PUT("/article/:id/info/watch/count", handler.ArticleHandler.IncrementArticleInfoWatchCountByID)
		rg.PublicRouter.POST("/article/comment", handler.ArticleHandler.CreateComment)
	}
	{
		rg.InternalRouter.POST("/article", handler.ArticleHandler.CreateArticle)
	}
}
