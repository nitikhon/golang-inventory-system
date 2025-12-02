package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/service"
)

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
func (h *ItemHandler) Update(c *fiber.Ctx) error {
	var item entity.Item
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updatedItem, err := h.service.Update(&item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
