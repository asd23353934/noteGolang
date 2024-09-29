package view

import (
	"github.com/gocql/gocql"
)

type Article struct {
	ID           gocql.UUID `json:"id"`
	Tags         []string   `json:"tags"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	ShareCount   int        `json:"share_count"`
	CommentCount int        `json:"comment_count"`
	WatchCount   int        `json:"watch_count"`
	Description  string     `json:"description"`
	CreatedAt    int64      `json:"created_at"`
	UpdatedAt    int64      `json:"updated_at"`
}

type ArticleListItem struct {
	ID        gocql.UUID `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	AuthorID  gocql.UUID `json:"author_id"`
	CreatedAt int64      `json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
}

type ArticleTagInfo struct {
	Tag       string `json:"tag"`
	Count     int    `json:"count"`
	Status    int    `json:"status"`
	UpdatedAt int64  `json:"updated_at"`
}

type ArticleDescription struct {
	ID           gocql.UUID `json:"id"`
	Tags         []string   `json:"tags"`
	Title        string     `json:"title"`
	ShareCount   int        `json:"share_count"`
	CommentCount int        `json:"comment_count"`
	WatchCount   int        `json:"watch_count"`
	Description  string     `json:"description"`
	CreatedAt    int64      `json:"created_at"`
	UpdatedAt    int64      `json:"updated_at"`
}

type CommentResponse struct {
	CommentID gocql.UUID         `json:"comment_id"`
	ArticleID gocql.UUID         `json:"article_id"`
	ParentID  *gocql.UUID        `json:"parent_id"`
	Name      string             `json:"name"`
	Content   string             `json:"content"`
	CreatedAt int64              `json:"created_at"`
	UpdatedAt int64              `json:"updated_at"`
	Status    int                `json:"status"`
	Replies   []*CommentResponse `json:"replies"`
}
