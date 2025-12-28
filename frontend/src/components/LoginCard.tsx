import React, { useState } from 'react'
import { useAuth } from '../hooks/useAuth'

const LoginCard = () => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')

  const { login, isPending } = useAuth()

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (username === '' || password === '') {
      alert('Username and Password must not be empty')
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
        <h2 className="text-2xl font-bold text-slate-900">Sign In</h2>
        <p className="text-slate-500 text-sm mt-1">Access your inventory dashboard</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-5">
        <div className="space-y-2">
          <label className="text-sm font-semibold text-slate-700 ml-1">Username</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={handleChange}
            placeholder="Enter your username"
            className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 outline-none transition-all"
          />
        </div>

        <div className="space-y-2">
          <label className="text-sm font-semibold text-slate-700 ml-1">Password</label>
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
          {isPending ? 'Authenticating...' : 'Login to System'}
        </button>
      </form>
    </div>
  )
}

export default LoginCard
