package handlers

import (
	"github.com/FulgurCode/BlankMap/internal/db/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type BlankMapHandler struct {
	queries *generated.Queries
}

func NewBlankMapHandler(queries *generated.Queries) *BlankMapHandler {
	return &BlankMapHandler{queries: queries}
}

func (h *BlankMapHandler) CreateBlankMap(c fiber.Ctx) error {
	type request struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Icon        *string `json:"icon"`
	}

	var body request
	if err := c.Bind().Body(&body); err != nil || body.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is required",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	bm, err := h.queries.CreateBlankMap(c.Context(), generated.CreateBlankMapParams{
		Name:        body.Name,
		Description: body.Description,
		Icon:        body.Icon,
		CreatedBy:   userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create blank map",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(bm)
}

func (h *BlankMapHandler) GetBlankMaps(c fiber.Ctx) error {
	maps, err := h.queries.GetAllBlankMaps(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch blank maps",
		})
	}
	return c.JSON(maps)
}

func (h *BlankMapHandler) GetBlankMapByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	bm, err := h.queries.GetBlankMapByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "blank map not found",
		})
	}

	return c.JSON(bm)
}

func (h *BlankMapHandler) UpdateBlankMap(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	type request struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Icon        *string `json:"icon"`
	}

	var body request
	if err := c.Bind().Body(&body); err != nil || body.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is required",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	bm, err := h.queries.UpdateBlankMap(c.Context(), generated.UpdateBlankMapParams{
		ID:          id,
		Name:        body.Name,
		Description: body.Description,
		Icon:        body.Icon,
		UpdatedBy:   userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update blank map",
		})
	}

	return c.JSON(bm)
}

func (h *BlankMapHandler) DeleteBlankMap(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.queries.DeleteBlankMap(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete blank map",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *BlankMapHandler) GetPinCount(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	count, err := h.queries.GetNoOfPins(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "blank map not found",
		})
	}

	return c.JSON(count)
}
