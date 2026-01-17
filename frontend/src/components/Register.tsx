import React, { useState } from 'react'
import userService from '../services/user'
import type { RegisterPayload } from '../types/user'
import toast from 'react-hot-toast'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import { useTranslation } from '../hooks/useTranslation'

const Register = () => {
  const navigate = useNavigate()
  const { setIsLoginModalOpen } = useAuth()
  const { t } = useTranslation() // Add this
  
  const [formData, setFormData] = useState<RegisterPayload>({
    username: '',
    password: '',
    email: '',
    first_name: '',
    last_name: '',
    phone: ''
  })
  const [isPending, setIsPending] = useState(false)
  const [errors, setErrors] = useState<Partial<RegisterPayload>>({})

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({ ...prev, [name]: value }))
    // Clear error when user types
    if (errors[name as keyof RegisterPayload]) {
      setErrors(prev => ({ ...prev, [name]: undefined }))
    }
  }

  const validate = (): boolean => {
    const newErrors: Partial<RegisterPayload> = {}
    let isValid = true

    // Email validation
    const emailRegex = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/
    if (!formData.email) {
      newErrors.email = t.register.validation.emailRequired
      isValid = false
    } else if (formData.email.length < 6 || formData.email.length > 128) {
      newErrors.email = t.register.validation.emailLength
      isValid = false
    } else if (!emailRegex.test(formData.email)) {
      newErrors.email = t.register.validation.emailFormat
      isValid = false
    }

    // Phone validation
    const phoneCleaned = formData.phone.replace(/[^0-9]/g, '')
    if (!formData.phone) {
      newErrors.phone = t.register.validation.phoneRequired
      isValid = false
    } else if (phoneCleaned.length !== 10) {
      newErrors.phone = t.register.validation.phoneFormat
      isValid = false
    }

    // Username validation
    const usernameRegex = /^[a-zA-Z0-9._]+$/
    if (!formData.username) {
      newErrors.username = t.register.validation.usernameRequired
      isValid = false
    } else if (formData.username.length < 3 || formData.username.length > 20) {
      newErrors.username = t.register.validation.usernameLength
      isValid = false
    } else if (!usernameRegex.test(formData.username)) {
      newErrors.username = t.register.validation.usernameFormat
      isValid = false
    } else if (formData.username.startsWith('.') || formData.username.startsWith('_') || 
               formData.username.endsWith('.') || formData.username.endsWith('_')) {
      newErrors.username = t.register.validation.usernameStartEnd
      isValid = false
    } else if (formData.username.includes('..') || formData.username.includes('__') || 
               formData.username.includes('._') || formData.username.includes('_.')) {
      newErrors.username = t.register.validation.usernameConsecutive
      isValid = false
    }

    // Password validation
    if (!formData.password) {
      newErrors.password = t.register.validation.passwordRequired
      isValid = false
    } else if (formData.password.length < 8 || formData.password.length > 128) {
      newErrors.password = t.register.validation.passwordLength
      isValid = false
    } else {
      if (!/[a-z]/.test(formData.password)) {
        newErrors.password = t.register.validation.passwordLower
        isValid = false
      } else if (!/[A-Z]/.test(formData.password)) {
        newErrors.password = t.register.validation.passwordUpper
        isValid = false
      } else if (!/[0-9]/.test(formData.password)) {
        newErrors.password = t.register.validation.passwordDigit
        isValid = false
      } else if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(formData.password)) {
        newErrors.password = t.register.validation.passwordSpecial
        isValid = false
      }
    }

    // Name validation helper
    const validateName = (name: string, field: 'first_name' | 'last_name') => {
      const nameRegex = /^[a-zA-ZÀ-ÿ\s\-']+$/
      if (!name) {
        newErrors[field] = field === 'first_name' ? t.register.validation.firstNameRequired : t.register.validation.lastNameRequired
        isValid = false
        return
      }
      if (name.length > 200) {
        newErrors[field] = t.register.validation.nameLength
        isValid = false
        return
      }
      if (!nameRegex.test(name)) {
        newErrors[field] = t.register.validation.nameFormat
        isValid = false
        return
      }
      if (name.includes('--') || name.includes("''") || name.includes('  ')) {
        newErrors[field] = t.register.validation.nameConsecutive
        isValid = false
        return
      }
      if (/^[\s\-']|[\s\-']$/.test(name)) {
        newErrors[field] = t.register.validation.nameStartEnd
        isValid = false
        return
      }
    }

    validateName(formData.first_name, 'first_name')
    validateName(formData.last_name, 'last_name')

    setErrors(newErrors)
    return isValid
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validate()) {
      toast.error(t.register.validation.error)
      return
    }

    setIsPending(true)
    try {
      await userService.register(formData)
      toast.success(t.register.success)
      navigate('/')
      setIsLoginModalOpen(true)
    } catch (error: any) {
      const msg = error.response?.data?.error || 'Registration failed'
      toast.error(msg)
    } finally {
      setIsPending(false)
    }
  }

  return (
    <div className="flex justify-center items-center min-h-[80vh]">
      <div className="w-full max-w-lg bg-white rounded-2xl shadow-2xl shadow-slate-900/10 p-8 m-4">
        <div className="text-center mb-8">
          <h2 className="text-3xl font-bold text-slate-900">{t.register.title}</h2>
          <p className="text-slate-500 mt-2">{t.register.subtitle}</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-semibold text-slate-700 ml-1">{t.register.labels.firstName}</label>
              <input
                name="first_name"
                value={formData.first_name}
                onChange={handleChange}
                className={`w-full px-4 py-2 mt-1 rounded-xl border ${errors.first_name ? 'border-red-500 ring-4 ring-red-500/10' : 'border-slate-200 focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10'} outline-none transition-all`}
                placeholder={t.register.placeholders.firstName}
              />
              {errors.first_name && <p className="text-red-500 text-xs mt-1 ml-1">{errors.first_name}</p>}
            </div>
            <div>
              <label className="text-sm font-semibold text-slate-700 ml-1">{t.register.labels.lastName}</label>
              <input
                name="last_name"
                value={formData.last_name}
                onChange={handleChange}
                className={`w-full px-4 py-2 mt-1 rounded-xl border ${errors.last_name ? 'border-red-500 ring-4 ring-red-500/10' : 'border-slate-200 focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10'} outline-none transition-all`}
                placeholder={t.register.placeholders.lastName}
              />
              {errors.last_name && <p className="text-red-500 text-xs mt-1 ml-1">{errors.last_name}</p>}
            </div>
          </div>

          <div>
            <label className="text-sm font-semibold text-slate-700 ml-1">{t.register.labels.username}</label>
            <input
              name="username"
              value={formData.username}
              onChange={handleChange}
              className={`w-full px-4 py-2 mt-1 rounded-xl border ${errors.username ? 'border-red-500 ring-4 ring-red-500/10' : 'border-slate-200 focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10'} outline-none transition-all`}
              placeholder={t.register.placeholders.username}
            />
            {errors.username && <p className="text-red-500 text-xs mt-1 ml-1">{errors.username}</p>}
          </div>

          <div>
            <label className="text-sm font-semibold text-slate-700 ml-1">{t.register.labels.email}</label>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              className={`w-full px-4 py-2 mt-1 rounded-xl border ${errors.email ? 'border-red-500 ring-4 ring-red-500/10' : 'border-slate-200 focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10'} outline-none transition-all`}
              placeholder={t.register.placeholders.email}
            />
             {errors.email && <p className="text-red-500 text-xs mt-1 ml-1">{errors.email}</p>}
          </div>

          <div>
            <label className="text-sm font-semibold text-slate-700 ml-1">{t.register.labels.phone}</label>
            <input
              name="phone"
              value={formData.phone}
              onChange={handleChange}
              className={`w-full px-4 py-2 mt-1 rounded-xl border ${errors.phone ? 'border-red-500 ring-4 ring-red-500/10' : 'border-slate-200 focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10'} outline-none transition-all`}
              placeholder={t.register.placeholders.phone}
            />
            {errors.phone && <p className="text-red-500 text-xs mt-1 ml-1">{errors.phone}</p>}
          </div>

          <div>
            <label className="text-sm font-semibold text-slate-700 ml-1">{t.register.labels.password}</label>
            <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              className={`w-full px-4 py-2 mt-1 rounded-xl border ${errors.password ? 'border-red-500 ring-4 ring-red-500/10' : 'border-slate-200 focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10'} outline-none transition-all`}
              placeholder={t.register.placeholders.password}
            />
            {errors.password && <p className="text-red-500 text-xs mt-1 ml-1">{errors.password}</p>}
          </div>

          <button
            disabled={isPending}
            type="submit"
            className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-300 text-white font-bold py-3.5 rounded-xl shadow-lg shadow-blue-200 transition-all transform active:scale-[0.98] mt-6"
          >
            {isPending ? t.register.submitting : t.register.submit}
          </button>
        </form>

        <div className="mt-6 text-center text-sm text-slate-500">
          {t.register.loginLink.text}{' '}
          <button onClick={() => { navigate('/'); setIsLoginModalOpen(true); }} className="text-blue-600 font-semibold hover:underline">
            {t.register.loginLink.action}
          </button>
        </div>
      </div>
    </div>
  )
}

export default Register
