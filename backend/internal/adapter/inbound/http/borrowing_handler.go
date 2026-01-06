package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/service"
	"github.com/nitikhon/golang-inventory-system/internal/util/errormap"
	"gorm.io/gorm"
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
	var borrowing entity.BorrowRequest
	if err := c.BodyParser(&borrowing); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	borrowing.UserID = c.Locals("user_id").(uint)

	// Input validation
	if borrowing.UserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user ID is required"})
	}
	if borrowing.ItemID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "item ID is required"})
	}
	if borrowing.BorrowingAmount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "borrowing amount must be greater than zero"})
	}

	borrowedItem, err := h.service.BorrowItem(borrowing)
	if err != nil {
		switch err.Error() {
		case errormap.ErrUserNotExist, errormap.ErrItemNotAvailable, gorm.ErrRecordNotFound.Error():
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case errormap.ErrItemNotEnough, errormap.ErrAlreadyBorrowed:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		case errormap.ErrInvalidDueDateValue, errormap.ErrInvalidDueDateFormat:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(borrowedItem)
}

// ApproveBorrowing handles approving a borrowing request.
func (h *BorrowingHandler) ApproveBorrowing(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	borrowingID, err := c.ParamsInt("id")
	if userID == 0 || err != nil || borrowingID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidRequestBody})
	}

	approvedBorrowing, err := h.service.ApproveBorrowing(uint(borrowingID), userID)
	if err != nil {
		switch err.Error() {
		case gorm.ErrRecordNotFound.Error():
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case errormap.ErrBorrowingNotExist, errormap.ErrApproverNotExist, errormap.ErrItemNotExistOrActive:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case errormap.ErrBorrowingNotPending, errormap.ErrNotEnoughItemsForApproval:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	return c.JSON(approvedBorrowing)
}

// RejectBorrowing handles rejecting a borrowing request.
func (h *BorrowingHandler) RejectBorrowing(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidRequestBody})
	}

	borrowingID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidRequestBody})
	}
	if borrowingID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidRequestBody})
	}

	rejectedBorrowing, err := h.service.RejectBorrowing(uint(borrowingID), userID)
	if err != nil {
		switch err.Error() {
		case errormap.ErrBorrowingNotExist, errormap.ErrRejecterNotExist, gorm.ErrRecordNotFound.Error():
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case errormap.ErrBorrowingNotPending:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		case errormap.ErrUnauthorizedToRejectAndCancel:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	return c.JSON(rejectedBorrowing)
}

// GetBorrowingsByStatus handles fetching borrowings by status.
func (h *BorrowingHandler) GetBorrowingsByBorrowingStatus(c *fiber.Ctx) error {
	status := c.Params("status")

	page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 12)
    search := c.Query("search")

	// Input validation
	if !isValidBorrowingStatus(status) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid borrowing status"})
	}

	borrowings, err := h.service.GetBorrowingsByBorrowingStatus(status, search, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(borrowings)
}

// GetBorrowingsByApprovalStatus handles fetching borrowings by approval status.
func (h *BorrowingHandler) GetBorrowingsByApprovalStatus(c *fiber.Ctx) error {
	status := c.Params("status")

	// Input validation
	if !isValidApprovalStatus(status) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid approval status"})
	}

	borrowings, err := h.service.GetBorrowingsByApprovalStatus(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(borrowings)
}

func (h *BorrowingHandler) GetBorrowingByUserID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 12)
    search := c.Query("search")

	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidRequestBody})
	}

	borrowings, err := h.service.GetBorrowingsByUserID(uint(userID), page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(borrowings)
}

func (h *BorrowingHandler) UserStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidRequestBody})
	}

	borrowingStat, err := h.service.GetUserBorrowingStats(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(borrowingStat)
}

func isValidBorrowingStatus(status string) bool {
	return status == entity.BORROWING_PENDING ||
		status == entity.BORROWING_ACTIVE ||
		status == entity.BORROWING_RETURNED ||
		status == entity.BORROWING_OVERDUE ||
		status == entity.BORROWING_CANCELLED ||
		status == entity.BORROWING_LOST
}

func isValidApprovalStatus(status string) bool {
	return status == entity.APPROVAL_PENDING ||
		status == entity.APPROVAL_APPROVED ||
		status == entity.APPROVAL_REJECTED
}
