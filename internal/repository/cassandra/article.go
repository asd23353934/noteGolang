package cassandra

import (
	"NoteGolang/internal/models/db"
	"context"
	"fmt"

	"github.com/gocql/gocql"
)

const (
	// Table names
	articlesTable      = "articles.articles"
	articlesByTagTable = "articles.articles_by_tag"
	tagInfoTable       = "articles.tag_info"
	commentsTable      = "articles.comments"
	articlesInfoTable  = "articles.articles_info"

	// TTL for tag info (90 days)
	tagInfoTTL = 90 * 24 * 60 * 60
)

// Article related queries
const (
	createArticle = `
		INSERT INTO ` + articlesTable + ` (
			id, title, content, description, created_at, updated_at, tags, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	// getArticleDescription = `
	// 	SELECT id, title, description, created_at, updated_at, tags, status
	// 	FROM ` + articlesTable

	getPaginatedArticleDescriptions = `
		SELECT id, title, description, created_at, updated_at, tags, status
		FROM ` + articlesTable + `
		LIMIT ?`

	getArticleByID = `
		SELECT id, title, content, description, created_at, updated_at, tags, status
		FROM ` + articlesTable + `
		 WHERE id = ?
	 `
	getArticleDescriptionsByID = `
		SELECT id, title, description, created_at, updated_at, tags, status
		FROM ` + articlesTable + `
		 WHERE id = ?
	 `
	createArticleInfo = `INSERT INTO ` + articlesInfoTable + ` (article_id, comment_count, share_count, watch_count, updated_at) VALUES (?, ?, ?, ?, ?) IF NOT EXISTS`

	updateArticleInfoByCommentCount = `
		UPDATE ` + articlesInfoTable + ` 
		SET comment_count = ?, updated_at = ? 
		WHERE article_id = ? 
		IF comment_count = ?`

	updateArticleInfoByShareCount = `
		UPDATE ` + articlesInfoTable + ` 
		SET share_count = ?, updated_at = ? 
		WHERE article_id = ? 
		IF share_count = ?`
	updateArticleInfoByWatchCount = `
		UPDATE ` + articlesInfoTable + ` 
		SET watch_count = ?, updated_at = ? 
		WHERE article_id = ? 
		IF watch_count = ?`

	getArticleInfoByArticleID = `
		SELECT comment_count, share_count, watch_count, updated_at 
		FROM ` + articlesInfoTable + `
		 WHERE article_id = ?`
)

// Tag related queries
const (
	createArticleTag = `
		INSERT INTO ` + articlesByTagTable + ` (
			tag, article_id, title, created_at
		) VALUES (?, ?, ?, ?) 
		IF NOT EXISTS`
	getArticleTagByTag = `
		SELECT article_id, title, created_at FROM
	` + articlesByTagTable + ` WHERE tag = ?`

	createTagInfo = `
		INSERT INTO ` + tagInfoTable + ` (
			tag, status, count, updated_at
		) VALUES (?, ?, ?, ?) 
 	`

	updateTagInfo = `
		UPDATE ` + tagInfoTable + ` 
		SET count = ?, status = ?, updated_at = ? 
		WHERE tag = ? 
		IF count = ?`

	// getArticleTag = `
	// 	SELECT tag, article_id, title, created_at
	// 	FROM ` + articlesByTagTable

	getTagInfo = `
		SELECT tag, status, count, updated_at 
		FROM ` + tagInfoTable

	getTagInfoByTag = getTagInfo + ` WHERE tag = ?`
)

// Comment related queries
const (
	createComment = `
		INSERT INTO ` + commentsTable + ` (
			article_id, comment_id, parent_id, name, content, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	getCommentByArticleID = `
		SELECT article_id, comment_id, parent_id, name, content, status, created_at, updated_at 
		FROM ` + commentsTable + `
		WHERE article_id = ?`
)

type ArticleRepository struct {
	db *CassandraDB
}

func NewArticleRepository(db *CassandraDB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// articlesTable
func (r *ArticleRepository) CreateArticle(ctx context.Context, data *db.Article) error {
	return r.db.Session.Query(createArticle, data.ID, data.Title, data.Content, data.Description, data.CreatedAt, data.UpdatedAt, data.Tags, data.Status).WithContext(ctx).Exec()
}
func (r *ArticleRepository) GetPaginatedArticleDescriptions(ctx context.Context, page int, pageSize int) ([]db.ArticleDescription, error) {
	// offset := (page - 1) * pageSize
	iter := r.db.Session.Query(getPaginatedArticleDescriptions, pageSize).WithContext(ctx).Iter()

	var articles []db.ArticleDescription
	var article db.ArticleDescription

	for iter.Scan(&article.ID, &article.Title, &article.Description, &article.CreatedAt, &article.UpdatedAt, &article.Tags, &article.Status) {
		articles = append(articles, article)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return articles, nil
}
func (r *ArticleRepository) GetArticleByID(ctx context.Context, id gocql.UUID) (db.Article, error) {
	var data db.Article
	err := r.db.Session.Query(getArticleByID, id).WithContext(ctx).Scan(&data.ID, &data.Title, &data.Content, &data.Description, &data.CreatedAt, &data.UpdatedAt, &data.Tags, &data.Status)
	return data, err
}
func (r *ArticleRepository) GetArticleDescriptionsByID(ctx context.Context, id gocql.UUID) (db.ArticleDescription, error) {
	var data db.ArticleDescription
	err := r.db.Session.Query(getArticleDescriptionsByID, id).WithContext(ctx).Scan(&data.ID, &data.Title, &data.Description, &data.CreatedAt, &data.UpdatedAt, &data.Tags, &data.Status)
	return data, err
}

// articlesByTagTable
func (r *ArticleRepository) CreateArticleTags(ctx context.Context, data *db.ArticleTag) error {
	return r.db.Session.Query(createArticleTag, data.Tag, data.ArticleID, data.Title, data.CreatedAt).WithContext(ctx).Exec()
}
func (r *ArticleRepository) GetArticleTagsByTag(ctx context.Context, tag string) ([]db.ArticleTag, error) {
	var list []db.ArticleTag
	var data db.ArticleTag
	iter := r.db.Session.Query(getArticleTagByTag, tag).Iter()
	for iter.Scan(&data.ArticleID, &data.Title, &data.CreatedAt) {
		list = append(list, data)
	}
	if err := iter.Close(); err != nil {
		return list, fmt.Errorf("failed to GetArticleTags")
	}
	return list, nil
}

// tagInfoTable
func (r *ArticleRepository) CreateTagInfo(ctx context.Context, tag string, status int, count int, time int64) error {
	return r.db.Session.Query(createTagInfo,
		tag, status, count, time).WithContext(ctx).Exec()
}
func (r *ArticleRepository) IncrementTagInfo(ctx context.Context, tag string, status int, count int, time int64) error {
	tagInfo, err := r.GetTagInfoByTag(ctx, tag)
	if err != nil {
		if err == gocql.ErrNotFound {
			return r.CreateTagInfo(ctx, tag, status, count, time)
		}
		return err
	}
	newCount := tagInfo.Count + 1

	return r.UpdateTagInfo(ctx, tag, status, tagInfo.Count, newCount, time)
}
func (r *ArticleRepository) UpdateTagInfo(ctx context.Context, tag string, status int, oldCount, newCount int, time int64) error {
	return r.db.Session.Query(updateTagInfo,
		newCount, status, time, tag, oldCount).WithContext(ctx).Exec()
}
func (r *ArticleRepository) GetTagInfoByTag(ctx context.Context, tag string) (db.ArticleTagInfo, error) {
	var data db.ArticleTagInfo
	err := r.db.Session.Query(getTagInfoByTag,
		tag).WithContext(ctx).Scan(&data.Tag, &data.Status, &data.Count, &data.UpdatedAt)
	return data, err
}
func (r *ArticleRepository) GetTagInfo(ctx context.Context) ([]db.ArticleTagInfo, error) {
	var list []db.ArticleTagInfo
	var data db.ArticleTagInfo
	iter := r.db.Session.Query(getTagInfo).WithContext(ctx).Iter()
	for iter.Scan(&data.Tag, &data.Status, &data.Count, &data.UpdatedAt) {
		list = append(list, data)
	}
	if err := iter.Close(); err != nil {
		return list, fmt.Errorf("failed to GetTagInfoByTag")
	}
	return list, nil
}

// commentsTable
func (r *ArticleRepository) CreateComment(ctx context.Context, data *db.Comment) error {
	return r.db.Session.Query(createComment, data.ArticleID, data.CommentID, data.ParentID, data.Name, data.Content, data.Status, data.CreatedAt, data.UpdatedAt).WithContext(ctx).Exec()
}
func (r *ArticleRepository) CreateChildComment(ctx context.Context, data *db.Comment) error {
	return r.db.Session.Query(createComment, data.ArticleID, data.CommentID, data.ParentID, data.Name, data.Content, data.Status, data.CreatedAt, data.UpdatedAt).WithContext(ctx).Exec()
}
func (r *ArticleRepository) GetCommentsByArticleID(ctx context.Context, articleID gocql.UUID) ([]db.Comment, error) {

	iter := r.db.Session.Query(getCommentByArticleID, articleID).WithContext(ctx).Iter()

	var comments []db.Comment
	var comment db.Comment
	for iter.Scan(
		&comment.ArticleID,
		&comment.CommentID,
		&comment.ParentID,
		&comment.Name,
		&comment.Content,
		&comment.Status,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	) {
		comments = append(comments, comment)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return comments, nil
}

// articlesInfoTable
func (r *ArticleRepository) CreateArticleInfo(ctx context.Context, data db.ArticleInfo) error {
	return r.db.Session.Query(createArticleInfo, data.ArticleID, data.CommentCount, data.ShareCount, data.WatchCount, data.UpdatedAt).WithContext(ctx).Exec()
}

func (r *ArticleRepository) IncrementArticleInfoByCommentCount(ctx context.Context, data db.ArticleInfo) error {
	articleInfo, err := r.GetArticleInfoByArticleID(ctx, data.ArticleID)
	if err != nil {
		if err == gocql.ErrNotFound {
			_, err := r.GetArticleByID(ctx, data.ArticleID)
			if err != nil {
				return r.CreateArticleInfo(ctx, data)
			}
		}
	}
	newCount := articleInfo.CommentCount + 1

	return r.UpdateArticleInfoByCommentCount(ctx, newCount, data.ArticleID, data.UpdatedAt, articleInfo.CommentCount)
}
func (r *ArticleRepository) IncrementArticleInfoByShareCount(ctx context.Context, data db.ArticleInfo) error {
	articleInfo, err := r.GetArticleInfoByArticleID(ctx, data.ArticleID)
	if err != nil {
		if err == gocql.ErrNotFound {
			_, err := r.GetArticleByID(ctx, data.ArticleID)
			if err != nil {
				return r.CreateArticleInfo(ctx, data)
			}
		}
	}
	newCount := articleInfo.ShareCount + 1

	return r.UpdateArticleInfoByShareCount(ctx, newCount, data.ArticleID, data.UpdatedAt, articleInfo.ShareCount)
}
func (r *ArticleRepository) IncrementArticleInfoByWatchCount(ctx context.Context, data db.ArticleInfo) error {
	articleInfo, err := r.GetArticleInfoByArticleID(ctx, data.ArticleID)
	if err != nil {
		if err == gocql.ErrNotFound {
			_, err := r.GetArticleByID(ctx, data.ArticleID)
			if err == nil {
				return r.CreateArticleInfo(ctx, data)
			}
		}
		return err
	}
	newCount := articleInfo.WatchCount + 1

	return r.UpdateArticleInfoByWatchCount(ctx, newCount, data.ArticleID, data.UpdatedAt, articleInfo.WatchCount)
}
func (r *ArticleRepository) UpdateArticleInfoByCommentCount(ctx context.Context, newCount int, articleID gocql.UUID, updatedAt int64, oldCount int) error {
	return r.db.Session.Query(updateArticleInfoByCommentCount, newCount, updatedAt, articleID, oldCount).WithContext(ctx).Exec()
}
func (r *ArticleRepository) UpdateArticleInfoByShareCount(ctx context.Context, newCount int, articleID gocql.UUID, updatedAt int64, oldCount int) error {
	return r.db.Session.Query(updateArticleInfoByShareCount, newCount, updatedAt, articleID, oldCount).WithContext(ctx).Exec()
}
func (r *ArticleRepository) UpdateArticleInfoByWatchCount(ctx context.Context, newCount int, articleID gocql.UUID, updatedAt int64, oldCount int) error {
	return r.db.Session.Query(updateArticleInfoByWatchCount, newCount, updatedAt, articleID, oldCount).WithContext(ctx).Exec()
}
func (r *ArticleRepository) GetArticleInfoByArticleID(ctx context.Context, articleID gocql.UUID) (db.ArticleInfo, error) {
	var data db.ArticleInfo
	err := r.db.Session.Query(getArticleInfoByArticleID, articleID).WithContext(ctx).Scan(&data.CommentCount, &data.ShareCount, &data.WatchCount, &data.UpdatedAt)

	return data, err
}
