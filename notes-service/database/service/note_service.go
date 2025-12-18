package service

import (
	"context"
	"notes-service/database/models"
	"notes-service/database/repository"
)

type NoteService interface {
	CreateNote(ctx context.Context, note *models.Notes) error
	GetNoteByID(ctx context.Context, id int) (*models.Notes, error)
	GetAllNotes(ctx context.Context) ([]models.Notes, error)
	UpdateNote(ctx context.Context, note *models.Notes) error
	DeleteNote(ctx context.Context, id int) error
}

type noteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) NoteService {
	return &noteService{repo: repo}
}

func (s *noteService) CreateNote(ctx context.Context, user *models.Notes) error {
	return s.repo.CreateNote(ctx, user)
}

func (s *noteService) GetNoteByID(ctx context.Context, id int) (*models.Notes, error) {
	return s.repo.GetNoteByID(ctx, id)
}

func (s *noteService) GetAllNotes(ctx context.Context) ([]models.Notes, error) {
	return s.repo.GetAllNotes(ctx)
}

func (s *noteService) UpdateNote(ctx context.Context, note *models.Notes) error {
	return s.repo.UpdateNote(ctx, note)
}

func (s *noteService) DeleteNote(ctx context.Context, id int) error {
	return s.repo.DeleteNote(ctx, id)
}
