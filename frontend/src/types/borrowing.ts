import type { Item } from './item'

export type ApprovalStatus = 'pending' | 'approved' | 'rejected'
export type BorrowingStatus = 'pending' | 'active' | 'returned' | 'overdue' | 'cancelled' | 'lost'

export interface Borrowing {
  id: number
  user_id: number
  item_id: number
  item: Item
  description: string
  borrowed_at: string
  returned_at: string
  due_date: string
  borrowing_amount: number
  borrowing_status: BorrowingStatus
  approval_status: ApprovalStatus
  approval_at: string
  approved_by: number
  rejected_at: string
  rejected_by: number
}

export interface BorrowingStats {
  ongoing_borrows: number
  total_returned: number
  currently_borrows: number
}

export interface BorrowRequest {
  item_id: number
  borrowing_amount: number
  description: string
  due_date: string
}
