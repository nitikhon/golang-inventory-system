package http

import (
	"strings"

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

	borrowedItem, err := h.service.BorrowItem(c.UserContext(), borrowing)
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

	approvedBorrowing, err := h.service.ApproveBorrowing(c.UserContext(), uint(borrowingID), userID)
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

	rejectedBorrowing, err := h.service.RejectBorrowing(c.UserContext(), uint(borrowingID), userID)
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

func (h *BorrowingHandler) ReturnBorrowing(c *fiber.Ctx) error {
	borrowingID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidRequestBody})
	}
	if borrowingID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errormap.ErrInvalidRequestBody})
	}

	returnedBorrowing, err := h.service.ReturnBorrowing(c.UserContext(), uint(borrowingID))
	if err != nil {
		switch err.Error() {
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

	borrowings, err := h.service.GetBorrowingsByBorrowingStatus(c.UserContext(), statuses, search, page, limit)
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

	borrowings, err := h.service.GetBorrowingsByUserID(c.UserContext(), uint(userID), page, limit, search)
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

	borrowingStat, err := h.service.GetUserBorrowingStats(c.UserContext(), userID)
	if err != nil {
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
