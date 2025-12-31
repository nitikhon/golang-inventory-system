export interface User {
  id: number
  username: string
  email: string
  phone: string
  first_name: string
  last_name: string
  is_admin: boolean
  refresh_token: string
  created_at: string
  updated_at: string
  deleted_at: string
}

export interface LoginPayload {
  username: string
  password: string
}

export interface Token {
  access_token: string
}
