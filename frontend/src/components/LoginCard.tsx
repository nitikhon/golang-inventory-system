import React, { useState } from 'react'
import { useAuth } from '../hooks/useAuth'
import { Link } from 'react-router-dom'
import toast from 'react-hot-toast'
import { useTranslation } from '../hooks/useTranslation'

const LoginCard = () => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')

  const { login, isPending, setIsLoginModalOpen } = useAuth()
  const { t } = useTranslation()

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (username === '' || password === '') {
      toast.error(t.auth.validation.empty, { duration: 5000 })
    }

    login({ username, password })
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { id, value } = e.target

    if (id === 'username') setUsername(value)
    if (id === 'password') setPassword(value)
  }

  return (
    <div className="p-8">
      <div className="text-center mb-8">
        <h2 className="text-2xl font-bold text-slate-900">{t.auth.signIn}</h2>
        <p className="text-slate-500 text-sm mt-1">{t.auth.subtitle}</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-5">
        <div className="space-y-2">
          <label className="text-sm font-semibold text-slate-700 ml-1">{t.auth.username}</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={handleChange}
            placeholder={t.auth.usernamePlaceholder}
            className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 outline-none transition-all"
          />
        </div>

        <div className="space-y-2">
          <label className="text-sm font-semibold text-slate-700 ml-1">{t.auth.password}</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={handleChange}
            placeholder="••••••••"
            className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 outline-none transition-all"
          />
        </div>

        <button
          disabled={isPending}
          type="submit"
          className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-300 text-white font-bold py-3.5 rounded-xl shadow-lg shadow-blue-200 transition-all transform active:scale-[0.98] mt-4"
        >
          {isPending ? t.auth.authenticating : t.auth.loginButton}
        </button>
      </form>

      <div className="mt-6 text-center text-sm text-slate-500">
        {t.register.signupLink.text}{' '}
        <Link 
          to="/register" 
          className="text-blue-600 font-semibold hover:underline"
          onClick={() => setIsLoginModalOpen(false)}
        >
          {t.register.signupLink.action}
        </Link>
      </div>
    </div>
  )
}

export default LoginCard
