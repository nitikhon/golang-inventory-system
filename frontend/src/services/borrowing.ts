import axios from 'axios'
import type { Borrowing, BorrowingStats, BorrowRequest } from '../types/borrowing'
import type { PaginatedResponse } from '../types/pagination'

const isDev = import.meta.env.DEV
const baseUrl = isDev ? 'http://localhost:8080/api/borrows' : '/api/borrows'

const getBorrowingsByUserID = async (
  page = 1,
  limit = 12,
  search = '',
  access_token: string | undefined
): Promise<PaginatedResponse<Borrowing>> => {
  const request = await axios.get(`${baseUrl}/user/?page=${page}&limit=${limit}&search=${search}`, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${access_token}`,
    },
    withCredentials: true,
  })
  return request.data
}

const getBorrowingStats = async (access_token: string | undefined): Promise<BorrowingStats> => {
  const request = await axios.get(`${baseUrl}/stats/`, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${access_token}`,
    },
    withCredentials: true,
  })
  return request.data
}

const createBorrowing = async (
  data: BorrowRequest | undefined,
  access_token: string | undefined
): Promise<Borrowing> => {
  const request = await axios.post(`${baseUrl}`, data, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${access_token}`,
    },
    withCredentials: true,
  })
  return request.data
}

const cancelBorrowing = async (
  borrowing_id: number | undefined,
  access_token: string | undefined
): Promise<Borrowing> => {
  const request = await axios.post(
    `${baseUrl}/reject/${borrowing_id}`,
    {},
    {
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${access_token}`,
      },
      withCredentials: true,
    }
  )
  return request.data
}

const getBorrowingsByStatus = async (
  page = 1,
  limit = 12,
  search = '',
  status: string | undefined,
  access_token: string | undefined
): Promise<PaginatedResponse<Borrowing>> => {
  const request = await axios.get(
    `${baseUrl}/status/${status}?page=${page}&limit=${limit}&search=${search}`,
    {
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${access_token}`,
      },
      withCredentials: true,
    }
  )
  return request.data
}

const approveBorrowing = async (
  borrowing_id: number | undefined,
  access_token: string | undefined
) => {
  const request = await axios.post(
    `${baseUrl}/approve/${borrowing_id}`,
    {},
    {
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${access_token}`,
      },
      withCredentials: true,
    }
  )
  return request.data
}

export default {
  getBorrowingsByUserID,
  getBorrowingStats,
  createBorrowing,
  cancelBorrowing,
  getBorrowingsByStatus,
  approveBorrowing,
}
