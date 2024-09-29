package handlers

import (
	"NoteGolang/internal/service"
	"NoteGolang/pkg/apperrors"
)

type Handlers struct {
	ArticleHandler  *ArticleHandler
	ResponseHandler *apperrors.ResponseHandler
	// 可以添加更多處理器...
}

func NewHandlers(services *service.Services) *Handlers {
	responseHandler := apperrors.NewResponseHandler()
	return &Handlers{
		ArticleHandler: NewArticleHandler(services.Article, responseHandler),
	}
}
