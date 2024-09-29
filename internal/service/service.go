package service

import (
	"NoteGolang/internal/repository"
)

type Services struct {
	Article *ArticleService
}

type ServiceDependencies struct {
	Repos *repository.Repositories
}

func NewServices(deps *ServiceDependencies) (*Services, error) {
	articleService := NewArticleService(*deps.Repos.ArticleRepo, *deps.Repos.ArticleCache)
	return &Services{
		Article: articleService,
	}, nil
}
