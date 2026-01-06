export interface PaginatedResponse<T> {
  data: T[]
  total_items: number
  total_pages: number
  page: number
  limit: number
}
