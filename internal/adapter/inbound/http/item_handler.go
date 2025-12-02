package http

import (
	"errors"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/service"
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
		case "record not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "record not found"})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	if err := validateItemInput(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	item.Name = normalizeItemName(item.Name)

	createdItem, err := h.service.Create(&item)
	if err != nil {
		// Handle business logic errors
		if err.Error() == "item with this name already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Item with this name already exists"})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "item ID is required for update"})
	}

	_, err := h.service.GetItemByID(uint(item.ID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "item not found"})
	}

	if err := validateItemInput(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	updatedItem, err := h.service.Update(&item)
	if err != nil {
		// Handle business logic errors
		switch err.Error() {
		case "item not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
		case "item with this name already exists":
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Item with this name already exists"})
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "item not found"})
	}

	var req struct {
		Name            *string `json:"name"`
		Description     *string `json:"description"`
		AvailableAmount *int    `json:"available_amount"`
		TotalAmount     *int    `json:"total_amount"`
		Status          *string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid item payload"})
	}

	if req.Name != nil {
		if *req.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name should not be empty"})
		}

		item.Name = normalizeItemName(*req.Name)
	}

	if req.Description != nil {
		item.Description = *req.Description
	}

	if req.AvailableAmount != nil {
		if *req.AvailableAmount < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "available_amount cannot be negative"})
		}

		if (req.TotalAmount != nil && *req.AvailableAmount > *req.TotalAmount) || (req.TotalAmount == nil && *req.AvailableAmount > item.TotalAmount) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "available_amount cannot exceed total_amount"})
		}

		item.AvailableAmount = *req.AvailableAmount
	}

	if req.TotalAmount != nil {
		if *req.TotalAmount <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "total_amount must be greater than 0"})
		}

		if (req.AvailableAmount != nil && *req.TotalAmount < *req.AvailableAmount) || (req.AvailableAmount == nil && *req.TotalAmount < item.AvailableAmount) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "updated total_amount is less than item's available amount"})
		}

		item.TotalAmount = *req.TotalAmount
	}

	if req.Status != nil {
		if *req.Status == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "status must be specified"})
		}

		if !slices.Contains(validStatuses, *req.Status) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "status must be one of: available, borrowed, maintenance, lost"})
		}

		item.Status = *req.Status
	}

	updatedItem, err := h.service.Update(item)
	if err != nil {
		// Handle business logic errors
		switch err.Error() {
		case "item not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
		case "item with this name already exists":
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Item with this name already exists"})
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
		return errors.New("name is required")
	}

	if item.Description == "" {
		return errors.New("description is required")
	}

	// Data type and range validation
	if item.AvailableAmount < 0 {
		return errors.New("available_amount cannot be negative")
	}

	if item.TotalAmount <= 0 {
		return errors.New("total_amount must be greater than 0")
	}

	if item.AvailableAmount > item.TotalAmount {
		return errors.New("available_amount cannot exceed total_amount")
	}

	// Set default status if empty

	return nil
}
