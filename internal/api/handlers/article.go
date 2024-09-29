package handlers

import (
	"NoteGolang/internal/models/dto"
	"NoteGolang/internal/service"
	"NoteGolang/pkg/apperrors"
	"context"
	"strconv"
	"time"

	"NoteGolang/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type ArticleHandler struct {
	ArticleService  *service.ArticleService
	ResponseHandler *apperrors.ResponseHandler
}

func NewArticleHandler(articleService *service.ArticleService, responseHandler *apperrors.ResponseHandler) *ArticleHandler {
	return &ArticleHandler{ArticleService: articleService, ResponseHandler: responseHandler}
}

// PublicRouter
func (h *ArticleHandler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, "Hellow world")
}
func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	articleID, _ := gocql.ParseUUID(c.Param("id"))
	rep, err := h.ArticleService.GetArticleByID(ctx, articleID)
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, rep)
}
func (h *ArticleHandler) GetArticleDescriptions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	pageSizeQuery := c.DefaultQuery("pageSize", "10")
	pageQuery := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	rep, err := h.ArticleService.GetPaginatedArticleDescriptions(ctx, page, pageSize)
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, rep)
}
func (h *ArticleHandler) GetArticleDescriptionsByTag(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	pageSizeQuery := c.DefaultQuery("pageSize", "10")
	pageQuery := c.DefaultQuery("page", "1")
	tag := c.Param("tag")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	rep, err := h.ArticleService.GetPaginatedArticleDescriptionsByTag(ctx, tag, page, pageSize)
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, rep)
}
func (h *ArticleHandler) GetArticleTagInfos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	view, err := h.ArticleService.GetTagInfos(ctx)
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, view)
}
func (h *ArticleHandler) CreateComment(c *gin.Context) {
	var req dto.CreateCommentRequest
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	if !utils.BindJSON(c, &req) {
		return
	}
	err := h.ArticleService.CreateComment(ctx, req)
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	h.ResponseHandler.NoContent(c)
}
func (h *ArticleHandler) GetCommentByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	articleID, err := gocql.ParseUUID(c.Param("id"))
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	view, err := h.ArticleService.GetCommentByArticleID(ctx, articleID)
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, view)
}
func (h *ArticleHandler) IncrementArticleInfoShareCountByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	articleID, err := gocql.ParseUUID(c.Param("id"))

	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	err = h.ArticleService.IncrementArticleInfoShareCountByArticleID(ctx, articleID)
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	h.ResponseHandler.NoContent(c)
}
func (h *ArticleHandler) IncrementArticleInfoWatchCountByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	articleID, err := gocql.ParseUUID(c.Param("id"))

	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	err = h.ArticleService.IncrementArticleInfoWatchCountByArticleID(ctx, articleID)
	if err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	h.ResponseHandler.NoContent(c)
}

// InternalRouter
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req *dto.CreateArticleRequest
	ctx := c.Request.Context()

	if !utils.BindJSON(c, &req) {
		return
	}

	if err := h.ArticleService.CreateArticleWithTags(ctx, req); err != nil {
		h.ResponseHandler.Error(c, err)
		return
	}
	h.ResponseHandler.NoContent(c)
}
