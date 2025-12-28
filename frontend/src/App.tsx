import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Home from './components/Home'
import NavBar from './components/NavBar'
import Profile from './components/Profile'
import { AuthProvider } from './contexts/AuthContext'

const queryClient = new QueryClient()

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AuthProvider>
          <div className="min-h-screen bg-slate-50">
            <NavBar />
            <main className="max-w-7xl mx-auto px-4 py-8">
              <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/profile" element={<Profile />} />
              </Routes>
            </main>
          </div>
        </AuthProvider>
      </BrowserRouter>
    </QueryClientProvider>
  )
}

export default App
