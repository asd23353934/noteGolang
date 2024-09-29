package transformers

import (
	"NoteGolang/internal/models/db"
	"NoteGolang/internal/models/dto"
	"NoteGolang/internal/models/view"
	"NoteGolang/internal/utils"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gocql/gocql"
	"github.com/microcosm-cc/bluemonday"
)

type ArticleTransformer struct{}

func NewArticleTransformer() *ArticleTransformer {
	return &ArticleTransformer{}
}

func (t *ArticleTransformer) ArticleRequestToDB(req *dto.CreateArticleRequest) (*db.Article, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	tags := t.processTags(req.Tags)

	now := time.Now().Unix()

	return &db.Article{
		ID:          gocql.TimeUUID(),
		Title:       req.Title,
		Content:     utils.EncodeHTML(req.Content),
		Description: utils.EncodeHTML(t.generateDescription(req.Content)),
		CreatedAt:   now,
		UpdatedAt:   now,
		Tags:        tags,
		Status:      1,
	}, nil
}

func (t *ArticleTransformer) ArticleInfoRequestToDB(id gocql.UUID, commentCount int, shareCount int, watchCount int, updatedAt int64) db.ArticleInfo {
	return db.ArticleInfo{
		ArticleID:    id,
		CommentCount: commentCount,
		ShareCount:   shareCount,
		WatchCount:   watchCount,
		UpdatedAt:    updatedAt,
	}
}

func (t *ArticleTransformer) ArticleTagsRequestToDB(tag string, articleID gocql.UUID, title string, createdAt int64) (*db.ArticleTag, error) {
	return &db.ArticleTag{
		Tag:       tag,
		ArticleID: articleID,
		Title:     title,
		CreatedAt: createdAt,
	}, nil
}

func (t *ArticleTransformer) CommentRequestToDB(req dto.CreateCommentRequest) (*db.Comment, error) {
	now := time.Now().Unix()
	articleID, err := gocql.ParseUUID(req.ArticleID)
	if err != nil {
		return nil, err
	}
	var parentID *gocql.UUID
	if req.ParentID != "" {
		id, err := gocql.ParseUUID(req.ParentID)
		if err != nil {
			return nil, err
		}
		parentID = &id
	}

	return &db.Comment{
		ArticleID: articleID,
		CommentID: gocql.TimeUUID(),
		ParentID:  parentID,
		Name:      req.Name,
		Content:   utils.EncodeHTML(req.Content),
		Status:    1,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (t *ArticleTransformer) CommentsDBToResponse(dbList []db.Comment) ([]view.CommentResponse, error) {
	var list []view.CommentResponse
	var parentList []view.CommentResponse
	for _, dbItem := range dbList {
		if dbItem.ParentID == nil {
			data, _ := t.CommentDBToResponse(&dbItem)
			parentList = append(parentList, data)
		}
	}
	for _, parent := range parentList {
		data := parent
		for _, dbItem := range dbList {
			if dbItem.ParentID != nil && parent.CommentID == *dbItem.ParentID {
				child, _ := t.CommentDBToResponse(&dbItem)
				data.Replies = append(data.Replies, &child)
			}
		}
		list = append(list, data)
	}

	return list, nil
}

func (t *ArticleTransformer) CommentDBToResponse(data *db.Comment) (view.CommentResponse, error) {
	content, _ := utils.DecodeHTML(data.Content)
	return view.CommentResponse{
		CommentID: data.CommentID,
		ArticleID: data.ArticleID,
		ParentID:  data.ParentID,
		Name:      data.Name,
		Content:   content,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Status:    data.Status,
		Replies:   nil,
	}, nil
}

func (t *ArticleTransformer) TagInfosDBToResponse(dbList []db.ArticleTagInfo) ([]view.ArticleTagInfo, error) {
	var list []view.ArticleTagInfo

	for _, dbItem := range dbList {
		data := view.ArticleTagInfo{
			Tag:       dbItem.Tag,
			Count:     dbItem.Count,
			Status:    dbItem.Status,
			UpdatedAt: dbItem.UpdatedAt,
		}

		list = append(list, data)
	}

	return list, nil
}

func (t *ArticleTransformer) ArticleDescriptionsDBToResponse(article db.ArticleDescription, info db.ArticleInfo) view.ArticleDescription {
	des, _ := utils.DecodeHTML(article.Description)
	return view.ArticleDescription{
		ID:           article.ID,
		Tags:         article.Tags,
		Title:        article.Title,
		ShareCount:   info.ShareCount,
		CommentCount: info.CommentCount,
		WatchCount:   info.WatchCount,
		Description:  des,
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
	}
}

func (t *ArticleTransformer) ArticleDBToResponse(article db.Article, info db.ArticleInfo) view.Article {
	content, _ := utils.DecodeHTML(article.Content)
	des, _ := utils.DecodeHTML(article.Description)
	return view.Article{
		ID:           article.ID,
		Tags:         article.Tags,
		Title:        article.Title,
		Content:      content,
		ShareCount:   info.ShareCount,
		CommentCount: info.CommentCount,
		WatchCount:   info.WatchCount,
		Description:  des,
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
	}
}

func (t *ArticleTransformer) generateDescription(content string) string {
	// 移除 HTML 标签
	plainText := stripHTML(content)

	// 清理和截断文本
	description := truncateText(plainText, 200)

	return description
}

func stripHTML(content string) string {
	// 使用 bluemonday 来安全地删除所有 HTML 标签
	p := bluemonday.StripTagsPolicy()
	return p.Sanitize(content)
}

func truncateText(text string, maxLength int) string {
	// 移除多余的空白字符
	text = strings.Join(strings.Fields(text), " ")

	// 确保文本是有效的 UTF-8
	text = strings.Map(func(r rune) rune {
		if r == utf8.RuneError {
			return -1
		}
		return r
	}, text)

	// 截断文本
	if utf8.RuneCountInString(text) > maxLength {
		runes := []rune(text)
		text = string(runes[:maxLength])
		// 尝试在单词边界截断
		lastSpace := strings.LastIndex(text, " ")
		if lastSpace > maxLength/2 {
			text = text[:lastSpace]
		}
		text += "..."
	}

	return strings.TrimSpace(text)
}

func (t *ArticleTransformer) processTags(tags []string) []string {
	// 处理标签，例如去重、转小写等
	uniqueTags := make(map[string]bool)
	var result []string
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag != "" && !uniqueTags[tag] {
			uniqueTags[tag] = true
			result = append(result, tag)
		}
	}
	return result
}
