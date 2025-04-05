package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nitikhon/golang-inventory-system/internal/core/service"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
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

	item, err := h.service.GetItemByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(item)
}

// Create adds a new item.
func (h *ItemHandler) Create(c *fiber.Ctx) error {
	var item entity.Item
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	createdItem, err := h.service.Create(item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(createdItem)
}

// Update modifies an existing item.
func (h *ItemHandler) Update(c *fiber.Ctx) error {
	var item entity.Item
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updatedItem, err := h.service.Update(item)
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

	if err := h.service.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}