import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Home from './components/Home'
import NavBar from './components/NavBar'

const queryClient = new QueryClient()

function App() {

    return (
        <QueryClientProvider client={queryClient}>
            <div className='min-h-screen bg-slate-50'>
                <NavBar />
                <main className='max-w-7xl mx-auto px-4 py-8'>
                    <Home />
                </main>
            </div>
        </QueryClientProvider>
    )
}

export default App
