package redis

import (
	"NoteGolang/internal/models/view"
	"context"
	"encoding/json"
	"time"
)

type ArticleRepository struct {
	db *RedisDB
}

func NewArticleRepository(db *RedisDB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) GetTagInfosCache(ctx context.Context) ([]view.ArticleTagInfo, error) {
	dataStr, err := r.db.Get(ctx, "tag_infos")
	if err != nil {
		return nil, err
	}
	var tagInfos []view.ArticleTagInfo
	err = json.Unmarshal([]byte(dataStr), &tagInfos)
	if err != nil {
		return nil, err
	}
	return tagInfos, nil
}

func (r *ArticleRepository) SetTagInfosCache(ctx context.Context, data []view.ArticleTagInfo) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = r.db.Set(ctx, "tag_infos", dataByte, time.Minute*5)
	if err != nil {
		return err
	}

	return nil
}
