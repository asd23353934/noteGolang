package service

import (
	"NoteGolang/internal/models/dto"
	"NoteGolang/internal/models/view"
	"NoteGolang/internal/repository/cassandra"
	"NoteGolang/internal/repository/redis"
	"NoteGolang/internal/transformers"
	"context"
	"time"

	"github.com/gocql/gocql"
)

type ArticleService struct {
	articleRepo  cassandra.ArticleRepository
	articleCache redis.ArticleRepository
	transformer  *transformers.ArticleTransformer
}

func NewArticleService(articleRepo cassandra.ArticleRepository, articleCache redis.ArticleRepository) *ArticleService {
	return &ArticleService{
		articleRepo:  articleRepo,
		articleCache: articleCache,
		transformer:  transformers.NewArticleTransformer(),
	}
}

func (s *ArticleService) CreateArticleWithTags(ctx context.Context, req *dto.CreateArticleRequest) error {
	article, err := s.transformer.ArticleRequestToDB(req)
	if err != nil {
		return err
	}
	// 創建文章
	if err = s.articleRepo.CreateArticle(ctx, article); err != nil {
		return err
	}
	// 創建文章info
	articleInfo := s.transformer.ArticleInfoRequestToDB(article.ID, 0, 0, 0, article.CreatedAt)
	if err = s.articleRepo.CreateArticleInfo(ctx, articleInfo); err != nil {
		return err
	}
	for _, tag := range article.Tags {
		articleByTag, _ := s.transformer.ArticleTagsRequestToDB(tag, article.ID, article.Title, article.CreatedAt)
		// 創建文章tag
		if err := s.articleRepo.CreateArticleTags(ctx, articleByTag); err != nil {
			return err
		}
		// 創建tag info, 已存在則更新
		if err := s.articleRepo.IncrementTagInfo(ctx, tag, 1, 1, article.UpdatedAt); err != nil {
			return err
		}
	}
	return nil
}

func (s *ArticleService) GetTagInfos(ctx context.Context) ([]view.ArticleTagInfo, error) {
	// cache, err := s.articleCache.GetTagInfosCache(ctx)
	// if err == nil {
	// 	return cache, nil
	// }
	list, err := s.articleRepo.GetTagInfo(ctx)
	if err != nil {
		return nil, err
	}
	view, err := s.transformer.TagInfosDBToResponse(list)
	// if err == nil {
	// 	go s.articleCache.SetTagInfosCache(ctx, view)
	// }
	return view, err
}

func (s *ArticleService) GetArticleByID(ctx context.Context, id gocql.UUID) (view.Article, error) {
	var rep view.Article
	articleDB, err := s.articleRepo.GetArticleByID(ctx, id)
	if err != nil {
		return rep, err
	}
	infoDB, _ := s.articleRepo.GetArticleInfoByArticleID(ctx, id)
	rep = s.transformer.ArticleDBToResponse(articleDB, infoDB)
	return rep, nil
}

func (s *ArticleService) GetPaginatedArticleDescriptions(ctx context.Context, page int, pageSize int) ([]view.ArticleDescription, error) {
	list, err := s.articleRepo.GetPaginatedArticleDescriptions(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	var repList []view.ArticleDescription
	for _, item := range list {
		dbData, _ := s.articleRepo.GetArticleInfoByArticleID(ctx, item.ID)
		rep := s.transformer.ArticleDescriptionsDBToResponse(item, dbData)
		repList = append(repList, rep)
	}
	return repList, nil
}

func (s *ArticleService) GetPaginatedArticleDescriptionsByTag(ctx context.Context, tag string, page int, pageSize int) ([]view.ArticleDescription, error) {
	list, err := s.articleRepo.GetArticleTagsByTag(ctx, tag)
	if err != nil {
		return nil, err
	}
	var repList []view.ArticleDescription
	for _, item := range list {
		desDB, _ := s.articleRepo.GetArticleDescriptionsByID(ctx, item.ArticleID)
		infoDB, _ := s.articleRepo.GetArticleInfoByArticleID(ctx, item.ArticleID)
		rep := s.transformer.ArticleDescriptionsDBToResponse(desDB, infoDB)
		repList = append(repList, rep)
	}
	return repList, err
}

func (s *ArticleService) CreateComment(ctx context.Context, req dto.CreateCommentRequest) error {
	dbData, err := s.transformer.CommentRequestToDB(req)
	if err != nil {
		return err
	}
	err = s.articleRepo.CreateComment(ctx, dbData)
	if err != nil {
		return err
	}

	articleInfoData := s.transformer.ArticleInfoRequestToDB(dbData.ArticleID, 1, 0, 0, dbData.UpdatedAt)
	return s.articleRepo.IncrementArticleInfoByCommentCount(ctx, articleInfoData)
}

func (s *ArticleService) GetCommentByArticleID(ctx context.Context, articleID gocql.UUID) ([]view.CommentResponse, error) {
	list, err := s.articleRepo.GetCommentsByArticleID(ctx, articleID)
	if err != nil {
		return nil, err
	}
	return s.transformer.CommentsDBToResponse(list)
}

func (s *ArticleService) IncrementArticleInfoShareCountByArticleID(ctx context.Context, articleID gocql.UUID) error {
	now := time.Now().Unix()
	data := s.transformer.ArticleInfoRequestToDB(articleID, 0, 1, 0, now)
	return s.articleRepo.IncrementArticleInfoByShareCount(ctx, data)
}

func (s *ArticleService) IncrementArticleInfoWatchCountByArticleID(ctx context.Context, articleID gocql.UUID) error {
	now := time.Now().Unix()
	data := s.transformer.ArticleInfoRequestToDB(articleID, 0, 0, 1, now)
	return s.articleRepo.IncrementArticleInfoByWatchCount(ctx, data)
}
