package handlers

import (
	"strconv"

	"github.com/FulgurCode/BlankMap/internal/db/generated"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type PinHandler struct {
	queries *generated.Queries
}

func NewPinHandler(queries *generated.Queries) *PinHandler {
	return &PinHandler{queries: queries}
}

func (h *PinHandler) CreatePin(c fiber.Ctx) error {
	type request struct {
		Name       string     `json:"name"`
		BlankMapID *uuid.UUID `json:"blank_map_id"`
		Latitude   float64    `json:"latitude"`
		Longitude  float64    `json:"longitude"`
		Address    *string    `json:"address"`
		Contact    *string    `json:"contact"`
	}

	var body request
	if err := c.Bind().Body(&body); err != nil || body.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name, latitude, and longitude are required",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	var blankMapID uuid.UUID
	if body.BlankMapID != nil {
		blankMapID = *body.BlankMapID
	}

	pin, err := h.queries.CreatePin(c.Context(), generated.CreatePinParams{
		Name:       body.Name,
		BlankMapID: blankMapID,
		Latitude:   body.Latitude,
		Longitude:  body.Longitude,
		Address:    body.Address,
		Contact:    body.Contact,
		CreatedBy:  userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create pin",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(pin)
}

func (h *PinHandler) GetPins(c fiber.Ctx) error {
	// Support filtering by blank_map_id
	if bmID := c.Query("blank_map_id"); bmID != "" {
		id, err := uuid.Parse(bmID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid blank_map_id"})
		}
		pins, err := h.queries.GetPinsByBlankMapID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to fetch pins",
			})
		}
		return c.JSON(pins)
	}

	pins, err := h.queries.GetAllPins(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch pins",
		})
	}
	return c.JSON(pins)
}

func (h *PinHandler) GetPinByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	pin, err := h.queries.GetPinByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "pin not found",
		})
	}

	return c.JSON(pin)
}

func (h *PinHandler) GetPinsNearby(c fiber.Ctx) error {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid lat"})
	}

	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid lng"})
	}

	radius, err := strconv.ParseFloat(c.Query("radius", "5"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid radius"})
	}

	pins, err := h.queries.GetPinsNearLocation(c.Context(), generated.GetPinsNearLocationParams{
		Radians:   lat,    // $1 — user latitude
		Radians_2: lng,    // $2 — user longitude
		Latitude:  radius, // $3 — radius in km
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch nearby pins",
		})
	}

	return c.JSON(pins)
}

func (h *PinHandler) UpdatePin(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	type request struct {
		Name    string  `json:"name"`
		Address *string `json:"address"`
		Contact *string `json:"contact"`
	}

	var body request
	if err := c.Bind().Body(&body); err != nil || body.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is required",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	pin, err := h.queries.UpdatePin(c.Context(), generated.UpdatePinParams{
		ID:        id,
		Name:      body.Name,
		Address:   body.Address,
		Contact:   body.Contact,
		UpdatedBy: userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update pin",
		})
	}

	return c.JSON(pin)
}

func (h *PinHandler) DeletePin(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.queries.DeletePin(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete pin",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
