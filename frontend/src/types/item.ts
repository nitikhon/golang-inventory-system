export type ItemStatus = 'available' | 'borrowed' | 'maintenance' | 'lost'

export interface Item {
  id: number
  available_amount: number
  total_amount: number
  name: string
  description: string
  status: ItemStatus
}

export interface ItemCreateRequest {
  available_amount: number
  total_amount: number
  name: string
  description: string
  status: ItemStatus
}

export interface ItemPatchRequest {
  available_amount?: number
  total_amount?: number
  name?: string
  description?: string
  status?: ItemStatus
}

