package repository

import (
	"context"
	"database/sql"
	"notes-service/database/models"
)

type NoteRepository interface {
	CreateNote(ctx context.Context, note *models.Notes) error
	GetNoteByID(ctx context.Context, id int) (*models.Notes, error)
	GetAllNotes(ctx context.Context) ([]models.Notes, error)
	UpdateNote(ctx context.Context, note *models.Notes) error
	DeleteNote(ctx context.Context, id int) error
}

type noteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) NoteRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) CreateNote(ctx context.Context, note *models.Notes) error {
	query := `INSERT INTO notes (user_publilc_id,content,category,created_at)
				VALUES ($1,$2,$3,$4)
				RETURNING id`
	return r.db.QueryRowContext(ctx, query).Scan(&note.Id)
}

func (r *noteRepository) GetNoteByID(ctx context.Context, id int) (*models.Notes, error) {
	var note models.Notes
	query := `SELECT id,title, content, category, color, created_at 
				FROM user 
				WHERE id=$1`
	err := r.db.QueryRowContext(ctx, query).Scan(
		&note.Id,
		&note.Title,
		&note.Content,
		&note.Category,
		&note.Color,
		&note.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &note, nil
}

func (r *noteRepository) GetAllNotes(ctx context.Context) ([]models.Notes, error) {
	query := `SELECT id,title, content, category, color, created_at 
				FROM notes 
				ORDER BY created_at`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Notes
	for rows.Next() {
		var note models.Notes
		if err := rows.Scan(
			&note.Id,
			&note.Title,
			&note.Content,
			&note.Category,
			&note.Color,
			&note.CreatedAt,
		); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (r *noteRepository) UpdateNote(ctx context.Context, note *models.Notes) error {
	query := `UPDATE notes
				SET title = $1,
					content=$2,
					category=$3,
					color=$4,
					created_at=$5
				WHERE id=$1`
	result, err := r.db.ExecContext(ctx, query,
		note.Id,
		note.Title,
		note.Content,
		note.Category,
		note.Color,
		note.CreatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *noteRepository) DeleteNote(ctx context.Context, id int) error {
	query := `DELETE FROM notes WHERE id=$1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
