package handlers

import (
	"notes-service/database/models"
	"notes-service/database/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type NoteHandler struct {
	noteService service.NoteService
}

func NewNoteHandler(noteService service.NoteService) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

func (h *NoteHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	notes := api.Group("/notes")
	notes.Post("/", h.CreateNote)
	notes.Get("/:id", h.GetNoteByID)
	notes.Get("/", h.GetAllNotes)
	notes.Put("/:id", h.UpdateNote)
	notes.Delete("/:id", h.DeleteNote)
}

func (h *NoteHandler) CreateNote(c *fiber.Ctx) error {
	var note models.Notes
	if err := c.BodyParser(&note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid request data",
			"details": err.Error(),
		})
	}

	if err := h.noteService.CreateNote(c.Context(), &note); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to create note",
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(note)
}
func (h *NoteHandler) GetNoteByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid note id",
			"details": err.Error(),
		})
	}
	note, err := h.noteService.GetNoteByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to get note",
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(note)
}
func (h *NoteHandler) GetAllNotes(c *fiber.Ctx) error {
	notes, err := h.noteService.GetAllNotes(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to get notes",
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(notes)
}
func (h *NoteHandler) UpdateNote(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid note id",
			"details": err.Error(),
		})
	}
	var note models.Notes
	if err := c.BodyParser(&note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid request data",
			"details": err.Error(),
		})
	}
	note.Id = id

	if err := h.noteService.UpdateNote(c.Context(), &note); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to update note",
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(note)
}

func (h *NoteHandler) DeleteNote(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid note id",
			"details": err.Error(),
		})
	}

	if err := h.noteService.DeleteNote(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to delete note",
			"details": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
