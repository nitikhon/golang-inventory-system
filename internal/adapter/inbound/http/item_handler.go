package http

import (
	"errors"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/service"
	"github.com/nitikhon/golang-inventory-system/internal/util/errormap"
)

// for item status enum validation
var validStatuses = []string{"available", "borrowed", "maintenance", "lost"}

// ItemHandler handles HTTP requests for items.
type ItemHandler struct {
	service *service.ItemService
}

// NewItemHandler creates a new ItemHandler.
func NewItemHandler(service *service.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

// GetAllItems retrieves all items.
func (h *ItemHandler) GetAllItems(c *fiber.Ctx) error {
	items, err := h.service.GetAllItems()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

// GetItemByID retrieves an item by its ID.
func (h *ItemHandler) GetItemByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	item, err := h.service.GetItemByID(uint(id))
	if err != nil {
		switch err.Error() {
		case errormap.ErrRecordNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": errormap.ErrRecordNotFound})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	return c.JSON(item)
}

// Create adds a new item.
func (h *ItemHandler) Create(c *fiber.Ctx) error {
	var item entity.Item

	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidJSONFormat})
	}

	if err := validateItemInput(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	item.Name = normalizeItemName(item.Name)

	createdItem, err := h.service.Create(&item)
	if err != nil {
		// Handle business logic errors
		if err.Error() == errormap.ErrItemNameAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": errormap.ErrItemNameAlreadyExists})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(createdItem)
}

// Update modifies an existing item.
func (h *ItemHandler) PutUpdate(c *fiber.Ctx) error {
	var item entity.Item
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	item.Name = normalizeItemName(item.Name)

	if item.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrItemIDRequired})
	}

	_, err := h.service.GetItemByID(uint(item.ID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": errormap.ErrItemNotFound})
	}

	if err := validateItemInput(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	updatedItem, err := h.service.Update(&item)
	if err != nil {
		// Handle business logic errors
		switch err.Error() {
		case errormap.ErrItemNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": errormap.ErrItemNotFound})
		case errormap.ErrItemNameAlreadyExists:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": errormap.ErrItemNameAlreadyExists})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(updatedItem)
}

func (h *ItemHandler) PatchUpdate(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	item, err := h.service.GetItemByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": errormap.ErrItemNotFound})
	}

	var req struct {
		Name            *string `json:"name"`
		Description     *string `json:"description"`
		AvailableAmount *int    `json:"available_amount"`
		TotalAmount     *int    `json:"total_amount"`
		Status          *string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidItemPayload})
	}

	if req.Name != nil {
		if *req.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrNameNotEmpty})
		}

		item.Name = normalizeItemName(*req.Name)
	}

	if req.Description != nil {
		item.Description = *req.Description
	}

	if req.AvailableAmount != nil {
		if *req.AvailableAmount < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrAvailableAmountNegative})
		}

		if (req.TotalAmount != nil && *req.AvailableAmount > *req.TotalAmount) || (req.TotalAmount == nil && *req.AvailableAmount > item.TotalAmount) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrAvailableExceedsTotal})
		}

		item.AvailableAmount = *req.AvailableAmount
	}

	if req.TotalAmount != nil {
		if *req.TotalAmount <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrTotalAmountPositive})
		}

		if (req.AvailableAmount != nil && *req.TotalAmount < *req.AvailableAmount) || (req.AvailableAmount == nil && *req.TotalAmount < item.AvailableAmount) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrTotalLessThanAvailable})
		}

		item.TotalAmount = *req.TotalAmount
	}

	if req.Status != nil {
		if *req.Status == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrStatusRequired})
		}

		if !slices.Contains(validStatuses, *req.Status) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidStatus})
		}

		item.Status = *req.Status
	}

	updatedItem, err := h.service.Update(item)
	if err != nil {
		// Handle business logic errors
		switch err.Error() {
		case errormap.ErrItemNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": errormap.ErrItemNotFound})
		case errormap.ErrItemNameAlreadyExists:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": errormap.ErrItemNameAlreadyExists})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(updatedItem)
}

// Delete removes an item by its ID.
func (h *ItemHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func normalizeItemName(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

// validateItemInput validates input data at handler level
func validateItemInput(item *entity.Item) error {
	// Required field validation
	if item.Name == "" {
		return errors.New(errormap.ErrNameRequired)
	}

	if item.Description == "" {
		return errors.New(errormap.ErrDescriptionRequired)
	}

	// Data type and range validation
	if item.AvailableAmount < 0 {
		return errors.New(errormap.ErrAvailableAmountNegative)
	}

	if item.TotalAmount <= 0 {
		return errors.New(errormap.ErrTotalAmountPositive)
	}

	if item.AvailableAmount > item.TotalAmount {
		return errors.New(errormap.ErrAvailableExceedsTotal)
	}

	// Set default status if empty

	return nil
}
