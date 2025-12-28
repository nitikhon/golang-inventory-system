import axios from 'axios'
import type { LoginPayload, Token, User } from '../types/user'

const isDev = import.meta.env.DEV
const baseUrl = isDev ? 'http://localhost:8080/api/users' : '/api/users'

const login = async (data: LoginPayload | undefined): Promise<Token> => {
  const request = await axios.post(`${baseUrl}/login`, data)
  return request.data
}

const logout = async () => {
  const request = await axios.post(`${baseUrl}/logout`, {}, { withCredentials: true })
  return request.data
}

const getProfile = async (access_token: string | undefined): Promise<User> => {
  const request = await axios.get(`${baseUrl}/me`, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${access_token}`,
    },
    withCredentials: true,
  })
  return request.data
}

const refreshToken = async (): Promise<Token> => {
  const request = await axios.get(`${baseUrl}/refresh`, { withCredentials: true })
  return request.data
}

export default { login, logout, getProfile, refreshToken }
