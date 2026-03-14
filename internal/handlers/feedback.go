package handlers

import (
	"github.com/FulgurCode/BlankMap/internal/db/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type FeedbackHandler struct {
	queries *generated.Queries
}

func NewFeedbackHandler(queries *generated.Queries) *FeedbackHandler {
	return &FeedbackHandler{queries: queries}
}

func (h *FeedbackHandler) CreateFeedback(c fiber.Ctx) error {
	type request struct {
		PinID  string  `json:"pin_id"`
		Rating *int32  `json:"rating"`
		Review *string `json:"review"`
	}

	var body request
	if err := c.Bind().Body(&body); err != nil || body.PinID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "pin_id is required",
		})
	}

	pinID, err := uuid.Parse(body.PinID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid pin_id",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	feedback, err := h.queries.CreateFeedback(c.Context(), generated.CreateFeedbackParams{
		PinID:  pinID,
		UserID: userID,
		Rating: body.Rating,
		Review: body.Review,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create feedback",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(feedback)
}

func (h *FeedbackHandler) GetFeedbackByPin(c fiber.Ctx) error {
	pinID, err := uuid.Parse(c.Params("pinID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid pin_id",
		})
	}

	feedback, err := h.queries.GetFeedbackByPinID(c.Context(), pinID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch feedback",
		})
	}

	return c.JSON(feedback)
}

func (h *FeedbackHandler) UpdateFeedback(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	type request struct {
		Rating *int32  `json:"rating"`
		Review *string `json:"review"`
	}

	var body request
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	feedback, err := h.queries.UpdateFeedback(c.Context(), generated.UpdateFeedbackParams{
		ID:     id,
		Rating: body.Rating,
		Review: body.Review,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update feedback",
		})
	}

	return c.JSON(feedback)
}

func (h *FeedbackHandler) DeleteFeedback(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	if err := h.queries.DeleteFeedback(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete feedback",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *FeedbackHandler) GetPinRating(c fiber.Ctx) error {
	pinID, err := uuid.Parse(c.Params("pinID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid pin_id",
		})
	}

	feedback, err := h.queries.GetPinRating(c.Context(), pinID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch feedback",
		})
	}

	return c.JSON(feedback)
}
