export type BorrowingStatus = 'pending' | 'approved' | 'rejected'

export interface Borrowing {
  id: number
  user_id: number
  item_id: number
  description: string
  borrow_at: Date
  returned_at: Date
  due_date: Date
  borrowing_amount: number
  borrowing_status: BorrowingStatus
  approval_at: Date
  approved_by: number
  rejected_at: Date
  rejected_by: number
}
