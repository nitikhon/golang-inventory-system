import axios from 'axios'
import type { Borrowing, BorrowingStats, BorrowRequest } from '../types/borrowing'

const isDev = import.meta.env.DEV
const baseUrl = isDev ? 'http://localhost:8080/api/borrows' : '/api/borrows'

const getBorrowingsByUserID = async (access_token: string | undefined): Promise<Borrowing[]> => {
  const request = await axios.get(`${baseUrl}/user/`, {
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

export default { getBorrowingsByUserID, getBorrowingStats, createBorrowing }
