import axios from 'axios'
import type { PaginatedResponse } from './pagination'
import type { Item } from '../types/item'

const isDev = import.meta.env.DEV
const baseUrl = isDev ? 'http://localhost:8080/api/items' : '/api/items'

const getAll = async (page = 1, limit = 12, search = ''): Promise<PaginatedResponse<Item>> => {
  const request = await axios.get(`${baseUrl}?page=${page}&limit=${limit}&search=${search}`)
  return request.data
}

export default { getAll }
