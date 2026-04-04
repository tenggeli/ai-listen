package mysql

import (
	"context"
	"database/sql"
	"strings"

	domain "listen/backend/internal/domain/audio"
)

type SoundContentRepository struct {
	db *sql.DB
}

func NewSoundContentRepository(db *sql.DB) SoundContentRepository {
	return SoundContentRepository{db: db}
}

func (r SoundContentRepository) ListCategories(ctx context.Context) ([]domain.Category, error) {
	const query = `
SELECT category_key, label
FROM audio_sound_categories
WHERE status = 'active'
ORDER BY sort_order ASC, category_key ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Category, 0)
	for rows.Next() {
		var item domain.Category
		if err := rows.Scan(&item.Key, &item.Label); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r SoundContentRepository) ListTracks(ctx context.Context, query domain.TrackQuery) ([]domain.Track, int, error) {
	where := "WHERE t.status = 'active'"
	args := make([]any, 0, 4)
	if query.CategoryKey != "" && query.CategoryKey != "all" {
		where += " AND t.category_key = ?"
		args = append(args, query.CategoryKey)
	}

	countSQL := "SELECT COUNT(1) FROM audio_sound_tracks t " + where
	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (query.PageNo - 1) * query.PageSize
	listSQL := `
SELECT
  t.track_id,
  t.title,
  COALESCE(c.label, ''),
  t.play_count_text,
  t.duration_text,
  t.emoji,
  t.author
FROM audio_sound_tracks t
LEFT JOIN audio_sound_categories c ON c.category_key = t.category_key
` + where + `
ORDER BY t.sort_order ASC, t.track_id ASC
LIMIT ? OFFSET ?`

	args = append(args, query.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.Track, 0)
	for rows.Next() {
		var item domain.Track
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Category,
			&item.PlayCountText,
			&item.DurationText,
			&item.Emoji,
			&item.Author,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r SoundContentRepository) ListAdminTracks(ctx context.Context, query domain.AdminTrackQuery) ([]domain.AdminTrack, int, error) {
	where := "WHERE 1=1"
	args := make([]any, 0, 8)
	if query.CategoryKey != "" && query.CategoryKey != "all" {
		where += " AND category_key = ?"
		args = append(args, query.CategoryKey)
	}
	if query.Status != "" {
		where += " AND status = ?"
		args = append(args, query.Status)
	}
	if strings.TrimSpace(query.Keyword) != "" {
		k := "%" + strings.TrimSpace(query.Keyword) + "%"
		where += " AND (track_id LIKE ? OR title LIKE ? OR author LIKE ?)"
		args = append(args, k, k, k)
	}

	countSQL := "SELECT COUNT(1) FROM audio_sound_tracks " + where
	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (query.PageNo - 1) * query.PageSize
	listSQL := `
SELECT track_id, category_key, title, play_count_text, duration_text, emoji, author, sort_order, status
FROM audio_sound_tracks
` + where + `
ORDER BY sort_order ASC, track_id ASC
LIMIT ? OFFSET ?`
	args = append(args, query.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.AdminTrack, 0)
	for rows.Next() {
		var item domain.AdminTrack
		if err := rows.Scan(
			&item.ID,
			&item.CategoryKey,
			&item.Title,
			&item.PlayCountText,
			&item.DurationText,
			&item.Emoji,
			&item.Author,
			&item.SortOrder,
			&item.Status,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r SoundContentRepository) GetAdminTrackByID(ctx context.Context, trackID string) (domain.AdminTrack, error) {
	const query = `
SELECT track_id, category_key, title, play_count_text, duration_text, emoji, author, sort_order, status
FROM audio_sound_tracks
WHERE track_id = ?
LIMIT 1`
	var item domain.AdminTrack
	err := r.db.QueryRowContext(ctx, query, trackID).Scan(
		&item.ID,
		&item.CategoryKey,
		&item.Title,
		&item.PlayCountText,
		&item.DurationText,
		&item.Emoji,
		&item.Author,
		&item.SortOrder,
		&item.Status,
	)
	if err == sql.ErrNoRows {
		return domain.AdminTrack{}, domain.ErrSoundNotFound
	}
	if err != nil {
		return domain.AdminTrack{}, err
	}
	return item, nil
}

func (r SoundContentRepository) CreateAdminTrack(ctx context.Context, input domain.AdminTrack) error {
	const query = `
INSERT INTO audio_sound_tracks(track_id, category_key, title, play_count_text, duration_text, emoji, author, sort_order, status, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
	_, err := r.db.ExecContext(
		ctx,
		query,
		input.ID,
		input.CategoryKey,
		input.Title,
		input.PlayCountText,
		input.DurationText,
		input.Emoji,
		input.Author,
		input.SortOrder,
		input.Status,
	)
	if isMySQLDuplicateError(err) {
		return domain.ErrInvalidInput
	}
	return err
}

func (r SoundContentRepository) UpdateAdminTrack(ctx context.Context, input domain.AdminTrack) error {
	const query = `
UPDATE audio_sound_tracks
SET category_key = ?, title = ?, play_count_text = ?, duration_text = ?, emoji = ?, author = ?, sort_order = ?, updated_at = NOW()
WHERE track_id = ?`
	result, err := r.db.ExecContext(
		ctx,
		query,
		input.CategoryKey,
		input.Title,
		input.PlayCountText,
		input.DurationText,
		input.Emoji,
		input.Author,
		input.SortOrder,
		input.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrSoundNotFound
	}
	return nil
}

func (r SoundContentRepository) UpdateAdminTrackStatus(ctx context.Context, trackID string, status string) error {
	const query = `UPDATE audio_sound_tracks SET status = ?, updated_at = NOW() WHERE track_id = ?`
	result, err := r.db.ExecContext(ctx, query, status, trackID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrSoundNotFound
	}
	return nil
}

var _ domain.Repository = SoundContentRepository{}
