package repository

import (
	"time"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"gorm.io/gorm"
)

// BorrowingRepository is a struct that implements the BorrowingRepository interface.
type BorrowingRepository struct {
	db *gorm.DB
}

// NewBorrowingRepository creates a new instance of BorrowingRepository.
func NewBorrowingRepository(db *gorm.DB) *BorrowingRepository {
	return &BorrowingRepository{db: db}
}

func (r *BorrowingRepository) BorrowItem(borrowing *entity.Borrowing) (*entity.Borrowing, error) {
	err := r.db.Create(&borrowing).Error
	if err != nil {
		return &entity.Borrowing{}, err
	}

	err = r.db.Preload("Item").First(borrowing, borrowing.ID).Error
    if err != nil {
        return nil, err
    }

	return borrowing, nil
}

// ApproveBorrowingWithTx updates the status of a borrowing record to approved and active within the provided transaction.
func (r *BorrowingRepository) ApproveBorrowingWithTx(tx *gorm.DB, borrowingId, approverId uint) (*entity.Borrowing, error) {
	err := tx.Model(&entity.Borrowing{}).Where("id = ?", borrowingId).
		Updates(&entity.Borrowing{
			BorrowingStatus: entity.BORROWING_ACTIVE,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
			ApprovedBy:      approverId,
			ApprovedAt:      time.Now().Format(time.RFC3339),
		}).Error
	if err != nil {
		return &entity.Borrowing{}, err
	}

	var updatedBorrowing entity.Borrowing
	if err := tx.First(&updatedBorrowing, borrowingId).Error; err != nil {
		return &entity.Borrowing{}, err
	}
	return &updatedBorrowing, nil
}

// RejectBorrowing updates the borrowing record to mark it as cancelled and rejected.
func (r *BorrowingRepository) RejectBorrowingWithTx(tx *gorm.DB, borrowingId, rejecterId uint) (*entity.Borrowing, error) {
	err := tx.Model(&entity.Borrowing{}).Where("id = ?", borrowingId).
		Updates(&entity.Borrowing{
			BorrowingStatus: entity.BORROWING_CANCELLED,
			ApprovalStatus:  entity.APPROVAL_REJECTED,
			RejectedBy:      rejecterId,
			RejectedAt:      time.Now().Format(time.RFC3339),
		}).Error
	if err != nil {
		return &entity.Borrowing{}, err
	}

	var updatedBorrowing entity.Borrowing
	if err := tx.First(&updatedBorrowing, borrowingId).Error; err != nil {
		return &entity.Borrowing{}, err
	}
	return &updatedBorrowing, nil
}

// GetAllBorrowings retrieves all borrowing records from the database.
func (r *BorrowingRepository) GetAllBorrowings() ([]*entity.Borrowing, error) {
	var borrowings []*entity.Borrowing
	err := r.db.Find(&borrowings).Error
	if err != nil {
		return nil, err
	}
	return borrowings, nil
}

// GetBorrowingByID retrieves a borrowing record by its ID.
func (r *BorrowingRepository) GetBorrowingByID(borrowingID uint) (*entity.Borrowing, error) {
	var borrowing entity.Borrowing
	err := r.db.First(&borrowing, borrowingID).Error
	if err != nil {
		return nil, err
	}
	return &borrowing, nil
}

// GetBorrowingsByUserID retrieves all borrowings for a specific user.
func (r *BorrowingRepository) GetBorrowingsByUserID(userID uint) ([]*entity.Borrowing, error) {
	var borrowings []*entity.Borrowing
	err := r.db.Preload("Item").Where("user_id = ?", userID).Find(&borrowings).Error
	if err != nil {
		return nil, err
	}
	return borrowings, nil
}

// GetBorrowingsByItemID retrieves all borrowings for a specific item.
func (r *BorrowingRepository) GetBorrowingsByItemID(itemID uint) ([]*entity.Borrowing, error) {
	var borrowings []*entity.Borrowing
	err := r.db.Where("item_id = ?", itemID).Find(&borrowings).Error
	if err != nil {
		return nil, err
	}
	return borrowings, nil
}

// GetBorrowingsByBorrowingStatus retrieves borrowings by their borrowing status.
func (r *BorrowingRepository) GetBorrowingsByBorrowingStatus(status string) ([]*entity.Borrowing, error) {
	var borrowings []*entity.Borrowing
	err := r.db.Where("borrowing_status = ?", status).Find(&borrowings).Error
	if err != nil {
		return nil, err
	}
	return borrowings, nil
}

// GetBorrowingsByApprovalStatus retrieves borrowings by their approval status.
func (r *BorrowingRepository) GetBorrowingsByApprovalStatus(status string) ([]*entity.Borrowing, error) {
	var borrowings []*entity.Borrowing
	err := r.db.Where("approval_status = ?", status).Find(&borrowings).Error
	if err != nil {
		return nil, err
	}
	return borrowings, nil
}

// GetBorrowingsByApproverID retrieves borrowings approved by a specific user.
func (r *BorrowingRepository) GetBorrowingsByApproverID(approverID uint) ([]*entity.Borrowing, error) {
	var borrowings []*entity.Borrowing
	err := r.db.Where("approved_by = ?", approverID).Find(&borrowings).Error
	if err != nil {
		return nil, err
	}
	return borrowings, nil
}

func (r *BorrowingRepository) GetUserBorrowingStats(userID uint) (*entity.BorrowingStats, error) {
	var result entity.BorrowingStats
	var countOngoing int64
	var countReturned int64

	err := r.db.Model(&entity.Borrowing{}).
		Where("user_id = ? AND borrowing_status = ?", userID, "pending").
		Count(&countOngoing).Error
	if err != nil {
		return nil, err
	}
	result.OnGoingBorrows = uint(countOngoing)

	err = r.db.Model(&entity.Borrowing{}).
		Where("user_id = ? AND borrowing_status = ?", userID, "returned").
		Count(&countReturned).Error
	if err != nil {
		return nil, err
	}

	result.TotalReturned = uint(countReturned)

	return &result, nil
}
