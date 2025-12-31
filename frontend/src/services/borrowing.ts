import axios from 'axios'
import type { Borrowing, BorrowingStats } from '../types/borrowing'

const isDev = import.meta.env.DEV
const baseUrl = isDev ? 'http://localhost:8080/api/borrows' : '/api/borrows'

const getBorrowingsByUserID = async (userId: number | undefined): Promise<Borrowing[]> => {
  const request = await axios.get(`${baseUrl}/user/${userId}`)
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

export default { getBorrowingsByUserID, getBorrowingStats }
