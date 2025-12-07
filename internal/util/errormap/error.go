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

	// Item errors
	ErrItemNotFound            = "item not found"
	ErrItemNameAlreadyExists   = "item with this name already exists"
	ErrRecordNotFound          = "record not found"
	ErrInvalidJSONFormat       = "Invalid JSON format"
	ErrInvalidItemPayload      = "invalid item payload"
	ErrItemIDRequired          = "item ID is required for update"
	ErrNameRequired            = "name is required"
	ErrDescriptionRequired     = "description is required"
	ErrNameNotEmpty            = "name should not be empty"
	ErrAvailableAmountNegative = "available_amount cannot be negative"
	ErrTotalAmountPositive     = "total_amount must be greater than 0"
	ErrAvailableExceedsTotal   = "available_amount cannot exceed total_amount"
	ErrTotalLessThanAvailable  = "updated total_amount is less than item's available amount"
	ErrStatusRequired          = "status must be specified"
	ErrInvalidStatus           = "status must be one of: available, borrowed, maintenance, lost"
)
