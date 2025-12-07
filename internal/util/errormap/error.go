package errormap

const (
	ErrUserNotExist              = "user does not exist or is not active"
	ErrItemNotAvailable          = "item is not available for borrowing"
	ErrItemNotEnough             = "item is not available enough to borrow"
	ErrAlreadyBorrowed           = "user has already borrowed this item and it is still pending"
	ErrBorrowingNotExist         = "borrowing does not exist"
	ErrApproverNotExist          = "approver does not exist or is not active"
	ErrRejecterNotExist          = "rejecter does not exist or is not active"
	ErrItemNotExistOrActive      = "item does not exist or is not active"
	ErrNotEnoughItemsForApproval = "not enough items available to approve borrowing"
	ErrBorrowingNotPending       = "borrowing request is not pending"
)
