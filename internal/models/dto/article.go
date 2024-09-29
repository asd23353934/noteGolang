package dto

type CreateArticleRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type CreateCommentRequest struct {
	ArticleID string `json:"article_id"`
	ParentID  string `json:"parent_id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
}
