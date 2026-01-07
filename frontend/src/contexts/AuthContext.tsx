import { useState, useEffect, type ReactNode } from 'react'
import type { User, LoginPayload, Token } from '../types/user'
import userService from '../services/user'
import { useMutation } from '@tanstack/react-query'
import type { AxiosError } from 'axios'
import { AuthContext } from './AuthContextType'
import { useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import toast from 'react-hot-toast'
import capitalizeSentence from '../utils/capitalizeSentence'
import Swal from 'sweetalert2'
import { useTranslation } from '../hooks/useTranslation'

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const queryClient = useQueryClient()
  const navigator = useNavigate()

  const { t } = useTranslation()

  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<Token | null>(null)
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false)
  const [isSilentLoading, setIsSilentLoading] = useState(true)

  useEffect(() => {
    const silentRefresh = async () => {
      try {
        const token = await userService.refreshToken()
        if (token.access_token !== undefined) {
          setToken(token)
          try {
            const user = await userService.getProfile(token.access_token)
            setUser(user)
          } catch {
            console.log('Silent refresh: Failed to retrieve data')
          }
        }
      } catch {
        console.log('Silent refresh: No session found')
      } finally {
        setIsSilentLoading(false)
      }
    }

    silentRefresh()
  }, [])

  const { mutate: loginMutation, isPending } = useMutation({
    mutationFn: userService.login,
    onSuccess: async (data: Token) => {
      setToken(data)
      try {
        const userProfile = await userService.getProfile(data?.access_token)
        setUser(userProfile)
        queryClient.setQueryData(['user'], userProfile)
      } catch {
        toast.error(t.common.error, { duration: 5000 })
      }
      setIsLoginModalOpen(false)
    },
    onError: (error: AxiosError<{ error: string }>) => {
      let message = error.response?.data?.error || t.common.error
      const status = error.response?.status
      if (status == 500) {
        message = t.common.error
      }
      toast.error(capitalizeSentence(message), { duration: 5000 })
    },
  })

  const { mutate: logoutMutation } = useMutation({
    mutationFn: userService.logout,
    onSuccess: () => {
      toast.success(t.auth.logoutSuccess)
      setUser(null)
      setToken(null)
      navigator('/')
      queryClient.removeQueries({ queryKey: ['user'] })
    },
    onError: () => {
      setUser(null)
      setToken(null)
      navigator('/')
    },
  })

  const login = (payload: LoginPayload) => {
    loginMutation(payload)
  }

  const logout = () => {
    Swal.fire({
      title: t.auth.logoutTitle,
      text: t.auth.logoutText,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#3085d6',
      cancelButtonColor: '#d33',
      confirmButtonText: t.profile.actions.signOut,
      cancelButtonText: t.profile.actions.signOutCancel
    }).then((result) => {
      if (result.isConfirmed) {
        logoutMutation()
      }
    })
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        login,
        logout,
        isPending,
        isLoginModalOpen,
        setIsLoginModalOpen,
        isSilentLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}
