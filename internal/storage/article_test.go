package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vova4o/tele-bot-ai/internal/model"
	"github.com/vova4o/tele-bot-ai/internal/storage"
)

func TestArticlePostgresStorage_Store(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	storage := storage.NewArticleStorage(sqlx.NewDb(db, "sqlmock"))

	article := model.Article{
		SourceID:    1,
		Title:       "Test Article",
		Link:        "https://example.com/article",
		Summary:     "This is a test article",
		PublishedAt: time.Now(),
	}

	mock.ExpectExec("INSERT OR IGNORE INTO articles").
		WithArgs(article.SourceID, article.Title, article.Link, article.Summary, article.PublishedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.Store(context.Background(), article)
	require.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestArticlePostgresStorage_AllNotPosted(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	storage := storage.NewArticleStorage(sqlx.NewDb(db, "sqlmock"))

	since := time.Now().Add(-24 * time.Hour)
	limit := uint64(10)

	rows := sqlmock.NewRows([]string{
		"a_id", "s_priority", "s_id", "a_title", "a_link", "a_summary", "a_published_at", "a_posted_at", "a_created_at",
	}).AddRow(1, 10, 1, "Article 1", "https://example.com/article1", "Summary 1", time.Now(), nil, time.Now())

	mock.ExpectQuery("SELECT").
		WithArgs(since.UTC().Format(time.RFC3339), limit).
		WillReturnRows(rows)

	articles, err := storage.AllNotPosted(context.Background(), since, limit)
	require.NoError(t, err)

	expectedArticle := model.Article{
		ID:          1,
		SourceID:    1,
		Title:       "Article 1",
		Link:        "https://example.com/article1",
		Summary:     "Summary 1",
		PublishedAt: articles[0].PublishedAt,
		CreatedAt:   articles[0].CreatedAt,
	}

	assert.Len(t, articles, 1)
	assert.Equal(t, expectedArticle, articles[0])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestArticlePostgresStorage_MarkAsPosted(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	storage := storage.NewArticleStorage(sqlx.NewDb(db, "sqlmock"))

	article := model.Article{
		ID: 1,
	}

	mock.ExpectExec("UPDATE articles").
		WithArgs(sqlmock.AnyArg(), article.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.MarkAsPosted(context.Background(), article)
	require.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
