import axios from 'axios'
import type { PaginatedResponse } from '../types/pagination'
import type { Item, ItemCreateRequest, ItemPatchRequest } from '../types/item'

const isDev = import.meta.env.DEV
const baseUrl = isDev ? 'http://localhost:8080/api/items' : '/api/items'

const getAll = async (page = 1, limit = 12, search = ''): Promise<PaginatedResponse<Item>> => {
  const request = await axios.get(`${baseUrl}?page=${page}&limit=${limit}&search=${search}`)
  return request.data
}

const create = async (data: ItemCreateRequest, access_token: string | undefined): Promise<Item> => {
  const request = await axios.post(`${baseUrl}?`, data, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${access_token}`,
    },
    withCredentials: true,
  })
  return request.data
}

const patchUpdate = async (id: number, data: ItemPatchRequest, access_token: string | undefined): Promise<Item> => {
  const request = await axios.patch(`${baseUrl}/${id}`, data, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${access_token}`,
    },
    withCredentials: true,
  })
  return request.data
}

export default { getAll, create, patchUpdate }
