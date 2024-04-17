package storage

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/vova4o/tele-bot-ai/internal/model"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	s := &ArticlePostgresStorage{db: db}
	s.setup()
	return s
}

func (s *ArticlePostgresStorage) setup() {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS articles (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            source_id INTEGER NOT NULL,
            title TEXT NOT NULL,
            link TEXT NOT NULL UNIQUE,
			summary TEXT,
            published_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            posted_at DATETIME,
            FOREIGN KEY (source_id) REFERENCES sources (id) ON DELETE CASCADE
        );
    `)
	if err != nil {
		log.Printf("[ERROR] failed to create table: %v", err)
		return
	}
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT OR IGNORE INTO articles (source_id, title, link, summary, published_at)
        				VALUES (?, ?, ?, ?, ?);`,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt.Format("2006-01-02 15:04:05"),
	); err != nil {
		return err
	}

	return nil
}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var articles []dbArticleWithPriority

	sinceStr := since.UTC().Format(time.RFC3339)
	sinceStr = strings.Replace(sinceStr, "T", " ", 1)
	sinceStr = strings.TrimSuffix(sinceStr, "Z")

	if err := conn.SelectContext(
		ctx,
		&articles,
		`SELECT 
                a.id AS a_id, 
                s.priority AS s_priority,
                s.id AS s_id,
                a.title AS a_title,
                a.link AS a_link,
                a.summary AS a_summary,
                a.published_at AS a_published_at,
                a.posted_at AS a_posted_at,
                a.created_at AS a_created_at
            FROM articles a JOIN sources s ON s.id = a.source_id
            WHERE a.posted_at IS NULL 
                AND a.published_at >= ?
            ORDER BY a.created_at DESC, s_priority DESC LIMIT ?;`,
		sinceStr,
		limit,
	); err != nil {
		return nil, err
	}

	return lo.Map(articles, func(article dbArticleWithPriority, _ int) model.Article {
		return model.Article{
			ID:          article.ID,
			SourceID:    article.SourceID,
			Title:       article.Title,
			Link:        article.Link,
			Summary:     article.Summary.String,
			PublishedAt: article.PublishedAt,
			CreatedAt:   article.CreatedAt,
		}
	}), nil
}

func (s *ArticlePostgresStorage) MarkAsPosted(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	nowStr := time.Now().UTC().Format(time.RFC3339)
	nowStr = strings.Replace(nowStr, "T", " ", 1)
	nowStr = strings.TrimSuffix(nowStr, "Z")

	if _, err := conn.ExecContext(
		ctx,
		`UPDATE articles SET posted_at = ? WHERE id = ?;`,
		nowStr,
		article.ID,
	); err != nil {
		return err
	}

	return nil
}

type dbArticleWithPriority struct {
	ID             int64          `db:"a_id"`
	SourcePriority int64          `db:"s_priority"`
	SourceID       int64          `db:"s_id"`
	Title          string         `db:"a_title"`
	Link           string         `db:"a_link"`
	Summary        sql.NullString `db:"a_summary"`
	PublishedAt    time.Time      `db:"a_published_at"`
	PostedAt       sql.NullTime   `db:"a_posted_at"`
	CreatedAt      time.Time      `db:"a_created_at"`
}
