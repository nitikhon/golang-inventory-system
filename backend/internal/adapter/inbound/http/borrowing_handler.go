package http

import (
	"context"
	"strings"
	"time"

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

	ctx, cancel := context.WithTimeout(c.UserContext(), 30*time.Second)
	defer cancel()

	borrowedItem, err := h.service.BorrowItem(ctx, borrowing)
	if err != nil {
		switch err.Error() {
		case context.DeadlineExceeded.Error():
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{"error": "The request took too long to process (timeout)"})
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

	ctx, cancel := context.WithTimeout(c.UserContext(), 30*time.Second)
	defer cancel()

	approvedBorrowing, err := h.service.ApproveBorrowing(ctx, uint(borrowingID), userID)
	if err != nil {
		switch err.Error() {
		case context.DeadlineExceeded.Error():
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{"error": "The request took too long to process (timeout)"})
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

	ctx, cancel := context.WithTimeout(c.UserContext(), 30*time.Second)
	defer cancel()

	rejectedBorrowing, err := h.service.RejectBorrowing(ctx, uint(borrowingID), userID)
	if err != nil {
		switch err.Error() {
		case context.DeadlineExceeded.Error():
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{"error": "The request took too long to process (timeout)"})
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

func (h *BorrowingHandler) ReturnBorrowing(c *fiber.Ctx) error {
	borrowingID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidRequestBody})
	}
	if borrowingID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidRequestBody})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 30*time.Second)
	defer cancel()

	returnedBorrowing, err := h.service.ReturnBorrowing(ctx, uint(borrowingID))
	if err != nil {
		switch err.Error() {
		case context.DeadlineExceeded.Error():
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{"error": "The request took too long to process (timeout)"})
		case errormap.ErrBorrowingNotExist, gorm.ErrRecordNotFound.Error():
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case errormap.ErrBorrowingNotActive:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	return c.JSON(returnedBorrowing)
}

// GetBorrowingsByStatus handles fetching borrowings by status.
func (h *BorrowingHandler) GetBorrowingsByBorrowingStatus(c *fiber.Ctx) error {
	statusParam := c.Query("status")

	var statuses []string
	if statusParam != "" {
		statuses = strings.Split(statusParam, ",")
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 12)
	search := c.Query("search")

	// Input validation
	if !isValidBorrowingStatus(statuses) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid borrowing status"})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 30*time.Second)
	defer cancel()

	borrowings, err := h.service.GetBorrowingsByBorrowingStatus(ctx, statuses, search, page, limit)
	if err != nil {
		switch err.Error() {
		case context.DeadlineExceeded.Error():
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{"error": "The request took too long to process (timeout)"})
		}
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

	ctx, cancel := context.WithTimeout(c.UserContext(), 30*time.Second)
	defer cancel()

	borrowings, err := h.service.GetBorrowingsByUserID(ctx, uint(userID), page, limit, search)
	if err != nil {
		switch err.Error() {
		case context.DeadlineExceeded.Error():
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{"error": "The request took too long to process (timeout)"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(borrowings)
}

func (h *BorrowingHandler) UserStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidRequestBody})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 30*time.Second)
	defer cancel()

	borrowingStat, err := h.service.GetUserBorrowingStats(ctx, userID)
	if err != nil {
		switch err.Error() {
		case context.DeadlineExceeded.Error():
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{"error": "The request took too long to process (timeout)"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(borrowingStat)
}

func isValidBorrowingStatus(statuses []string) bool {
	for _, status := range statuses {
		if status != entity.BORROWING_PENDING &&
			status != entity.BORROWING_ACTIVE &&
			status != entity.BORROWING_RETURNED &&
			status != entity.BORROWING_OVERDUE &&
			status != entity.BORROWING_CANCELLED &&
			status != entity.BORROWING_LOST &&
			status != entity.BORROWING_REJECTED {
			return false
		}
	}

	return true
}
