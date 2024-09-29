package db

import (
	"github.com/gocql/gocql"
)

type Article struct {
	ID          gocql.UUID `cql:"id"`
	Title       string     `cql:"title"`
	Content     string     `cql:"content"`
	Description string     `cql:"description"`
	CreatedAt   int64      `cql:"created_at"`
	UpdatedAt   int64      `cql:"updated_at"`
	Tags        []string   `cql:"tags"`
	Status      int        `cql:"status"`
}

type ArticleDescription struct {
	ID          gocql.UUID `cql:"id"`
	Title       string     `cql:"title"`
	Description string     `cql:"description"`
	CreatedAt   int64      `cql:"created_at"`
	UpdatedAt   int64      `cql:"updated_at"`
	Tags        []string   `cql:"tags"`
	Status      int        `cql:"status"`
}

type ArticleTag struct {
	Tag       string     `cql:"tag"`
	ArticleID gocql.UUID `cql:"article_id"`
	Title     string     `cql:"title"`
	CreatedAt int64      `cql:"created_at"`
}

type ArticleInfo struct {
	ArticleID    gocql.UUID `cql:"article_id"`
	CommentCount int        `cql:"comment_count"`
	ShareCount   int        `cql:"share_count"`
	WatchCount   int        `cql:"watch_count"`
	UpdatedAt    int64      `cql:"updated_at"`
}

type ArticleTagInfo struct {
	Tag       string `cql:"tag"`
	Count     int    `cql:"count"`
	Status    int    `cql:"status"`
	UpdatedAt int64  `cql:"updated_at"`
}

type Comment struct {
	ArticleID gocql.UUID  `cql:"article_id"`
	CommentID gocql.UUID  `cql:"comment_id"`
	ParentID  *gocql.UUID `cql:"parent_id"`
	Name      string      `cql:"name"`
	Content   string      `cql:"content"`
	Status    int         `cql:"status"`
	CreatedAt int64       `cql:"created_at"`
	UpdatedAt int64       `cql:"updated_at"`
}
