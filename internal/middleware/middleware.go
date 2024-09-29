package middleware

import (
	"NoteGolang/internal/repository"
	"NoteGolang/internal/service"
)

type Middleware struct {
	service *service.Services
	repos   *repository.Repositories
	// 可以添加其他依賴
}

func NewMiddleware(service *service.Services, repos *repository.Repositories) *Middleware {
	return &Middleware{
		service: service,
		repos:   repos,
	}
}
