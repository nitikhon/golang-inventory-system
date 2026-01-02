package entity

const (
	// BorrowingStatus constants
	BORROWING_PENDING   = "pending"
	BORROWING_ACTIVE    = "active"
	BORROWING_RETURNED  = "returned"
	BORROWING_OVERDUE   = "overdue"
	BORROWING_CANCELLED = "cancelled"
	BORROWING_LOST      = "lost"

	// ApprovalStatus constants
	APPROVAL_PENDING  = "pending"
	APPROVAL_APPROVED = "approved"
	APPROVAL_REJECTED = "rejected"
)

type Borrowing struct {
	GormModel
	UserID          uint   `json:"user_id"`
	ItemID          uint   `json:"item_id" gorm:"index"`
	Item            *Item  `json:"item"`
	Description     string `json:"description"`
	BorrowedAt      string `json:"borrowed_at"`
	ReturnedAt      string `json:"returned_at"`
	DueDate         string `json:"due_date"`
	BorrowingAmount int    `json:"borrowing_amount"`
	BorrowingStatus string `json:"borrowing_status" gorm:"type:VARCHAR(20);default:'pending';index"`
	ApprovalStatus  string `json:"approval_status" gorm:"type:VARCHAR(20);default:'pending';index"` // pending, approved, rejected
	ApprovedAt      string `json:"approved_at"`
	ApprovedBy      uint   `json:"approved_by"`
	RejectedBy      uint   `json:"rejected_by"`
	RejectedAt      string `json:"rejected_at"`
}

type BorrowingStats struct {
	OnGoingBorrows uint `json:"ongoing_borrows"`
	TotalReturned  uint `json:"total_returned"`
}

type BorrowRequest struct {
	UserID          uint
	ItemID          uint   `json:"item_id"`
	BorrowingAmount int    `json:"borrowing_amount"`
	Description     string `json:"description"`
	DueDate         string `json:"due_date"`
}
