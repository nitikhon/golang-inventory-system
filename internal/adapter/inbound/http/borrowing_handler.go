package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/service"
)

// BorrowingHandler handles HTTP requests for borrowing operations.
type BorrowingHandler struct {
	service *service.BorrowingService
}

// NewBorrowingHandler creates a new BorrowingHandler.
func NewBorrowingHandler(service *service.BorrowingService) *BorrowingHandler {
	return &BorrowingHandler{service: service}
}

// BorrowItem handles borrowing an item.
func (h *BorrowingHandler) BorrowItem(c *fiber.Ctx) error {
	var borrowing entity.Borrowing
	if err := c.BodyParser(&borrowing); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	borrowedItem, err := h.service.BorrowItem(borrowing)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(borrowedItem)
}

// ApproveBorrowing handles approving a borrowing request.
func (h *BorrowingHandler) ApproveBorrowing(c *fiber.Ctx) error {
	var borrowing entity.Borrowing
	if err := c.BodyParser(&borrowing); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	approvedBorrowing, err := h.service.ApproveBorrowing(borrowing.ID, borrowing.ApprovedBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(approvedBorrowing)
}

// RejectBorrowing handles rejecting a borrowing request.
func (h *BorrowingHandler) RejectBorrowing(c *fiber.Ctx) error {
	var borrowing entity.Borrowing
	if err := c.BodyParser(&borrowing); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	rejectedBorrowing, err := h.service.RejectBorrowing(borrowing.ID, borrowing.RejectedBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rejectedBorrowing)
}

// GetBorrowingsByStatus handles fetching borrowings by status.
func (h *BorrowingHandler) GetBorrowingsByBorrowingStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	borrowings, err := h.service.GetBorrowingsByBorrowingStatus(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(borrowings)
}

// GetBorrowingsByStatus handles fetching borrowings by status.
func (h *BorrowingHandler) GetBorrowingsByApprovalStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	borrowings, err := h.service.GetBorrowingsByApprovalStatus(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(borrowings)
}
