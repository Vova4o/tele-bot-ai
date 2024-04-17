package storage

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"github.com/vova4o/tele-bot-ai/internal/model"
)

type SourceSQLiteStorage struct {
	db *sqlx.DB
}

func NewSourceStorage(db *sqlx.DB) *SourceSQLiteStorage {
	s := &SourceSQLiteStorage{db: db}
	s.setup()
	return s
}

func (s *SourceSQLiteStorage) setup() {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS sources (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			feed_url TEXT NOT NULL,
			priority INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Printf("[ERROR] failed to create sources table: %v", err)
		return
	}
}

func (s *SourceSQLiteStorage) Sources(ctx context.Context) ([]model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var sources []dbSource
	if err := conn.SelectContext(ctx, &sources, `SELECT * FROM sources`); err != nil {
		return nil, err
	}

	return lo.Map(sources, func(source dbSource, _ int) model.Source { return model.Source(source) }), nil
}

func (s *SourceSQLiteStorage) SourceByID(ctx context.Context, id int64) (*model.Source, error) {
    conn, err := s.db.Connx(ctx)
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    var source dbSource
    if err := conn.GetContext(ctx, &source, `SELECT * FROM sources WHERE id = ?`, id); err != nil {
        return nil, err
    }

    return (*model.Source)(&source), nil
}

func (s *SourceSQLiteStorage) Add(ctx context.Context, source model.Source) (int64, error) {
    conn, err := s.db.Connx(ctx)
    if err != nil {
        return 0, err
    }
    defer conn.Close()

    res, err := conn.ExecContext(
        ctx,
        `INSERT INTO sources (name, feed_url, priority)
                    VALUES (?, ?, ?);`,
        source.Name, source.FeedURL, source.Priority,
    )
    if err != nil {
        return 0, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }

    return id, nil
}

func (s *SourceSQLiteStorage) SetPriority(ctx context.Context, id int64, priority int) error {
    conn, err := s.db.Connx(ctx)
    if err != nil {
        return err
    }
    defer conn.Close()

    _, err = conn.ExecContext(ctx, `UPDATE sources SET priority = ? WHERE id = ?`, priority, id)

    return err
}

func (s *SourceSQLiteStorage) Delete(ctx context.Context, id int64) error {
    conn, err := s.db.Connx(ctx)
    if err != nil {
        return err
    }
    defer conn.Close()

    if _, err := conn.ExecContext(ctx, `DELETE FROM sources WHERE id = ?`, id); err != nil {
        return err
    }

    return nil
}

type dbSource struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	Priority  int       `db:"priority"`
	CreatedAt time.Time `db:"created_at"`
}
