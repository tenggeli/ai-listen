package mysql

import (
	"context"
	"database/sql"
	"time"

	domain "listen/backend/internal/domain/ai"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return SessionRepository{db: db}
}

func (r SessionRepository) Create(ctx context.Context, session domain.Session) error {
	const insert = `
INSERT INTO ai_sessions(session_uid, user_id, scene_type, summary, status, last_message_at, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, NULLIF(?, ''), NOW(), NOW())`
	last := ""
	if !session.LastMessageAt.IsZero() {
		last = session.LastMessageAt.Format("2006-01-02 15:04:05")
	}
	_, err := r.db.ExecContext(ctx, insert, session.ID, session.UserID, session.SceneType, session.Summary, session.Status, last)
	return err
}

func (r SessionRepository) GetByID(ctx context.Context, id string) (domain.Session, error) {
	const query = `SELECT session_uid, user_id, scene_type, summary, status, last_message_at FROM ai_sessions WHERE session_uid = ? LIMIT 1`
	var session domain.Session
	var last sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(&session.ID, &session.UserID, &session.SceneType, &session.Summary, &session.Status, &last)
	if err == sql.ErrNoRows {
		return domain.Session{}, domain.ErrSessionNotFound
	}
	if err != nil {
		return domain.Session{}, err
	}
	if last.Valid {
		session.LastMessageAt = last.Time
	}

	const msgQuery = `SELECT sender_type, content, created_at FROM ai_messages WHERE session_uid = ? ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, msgQuery, id)
	if err != nil {
		return domain.Session{}, err
	}
	defer rows.Close()

	session.Messages = make([]domain.Message, 0)
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.SenderType, &m.Content, &m.CreatedAt); err != nil {
			return domain.Session{}, err
		}
		session.Messages = append(session.Messages, m)
	}
	if err := rows.Err(); err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (r SessionRepository) Save(ctx context.Context, session domain.Session) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const update = `UPDATE ai_sessions SET status = ?, summary = ?, last_message_at = ?, updated_at = NOW() WHERE session_uid = ?`
	last := time.Now()
	if !session.LastMessageAt.IsZero() {
		last = session.LastMessageAt
	}
	if _, err = tx.ExecContext(ctx, update, session.Status, session.Summary, last, session.ID); err != nil {
		return err
	}

	const deleteMsg = `DELETE FROM ai_messages WHERE session_uid = ?`
	if _, err = tx.ExecContext(ctx, deleteMsg, session.ID); err != nil {
		return err
	}

	const insertMsg = `INSERT INTO ai_messages(session_uid, sender_type, message_type, content, intent_json, safety_level, created_at, updated_at) VALUES(?, ?, 'text', ?, NULL, 'normal', ?, NOW())`
	for _, m := range session.Messages {
		if _, err = tx.ExecContext(ctx, insertMsg, session.ID, m.SenderType, m.Content, m.CreatedAt); err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}
