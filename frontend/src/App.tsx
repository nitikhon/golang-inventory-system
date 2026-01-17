import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Home from './components/Home'
import NavBar from './components/NavBar'
import Profile from './components/Profile'
import { AuthProvider } from './contexts/AuthContext'
import MyBorrowings from './components/MyBorrowings'
import { Toaster } from 'react-hot-toast'
import Admin from './components/Admin'
import { LanguageProvider } from './contexts/LanguageContext'
import Register from './components/Register'

const queryClient = new QueryClient()

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <LanguageProvider>
          <AuthProvider>
            <div className="min-h-screen bg-slate-50">
              <NavBar />
              <Toaster position="top-right" />
              <main className="max-w-7xl mx-auto px-4 py-8">
                <Routes>
                  <Route path="/register" element={<Register />} />
                  <Route path="/" element={<Home />} />
                  <Route path="/profile" element={<Profile />} />
                  <Route path="/my-borrowings" element={<MyBorrowings />} />
                  <Route path="/admin" element={<Admin />} />
                </Routes>
              </main>
            </div>
          </AuthProvider>
        </LanguageProvider>
      </BrowserRouter>
    </QueryClientProvider>
  )
}

export default App
