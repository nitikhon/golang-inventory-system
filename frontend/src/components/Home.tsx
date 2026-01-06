import itemService from '../services/item'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import type { Item } from '../types/item'
import ItemCard from './ItemCard'
import Loading from './Loading'
import Error from './Error'
import { useAuth } from '../hooks/useAuth'
import { useEffect, useState } from 'react'
import BorrowCard from './BorrowCard'
import { X } from 'lucide-react'
import SearchBar from './SearchBar'
import type { PaginatedResponse } from '../types/pagination'

const HomePage: React.FC = () => {
  const [selectedItem, setSelectedItem] = useState<Item | null>(null)
  const [isBorrowModalOpen, setIsBorrowModalOpen] = useState(false)
  const [searchTerm, setSearchTerm] = useState('')
  const [page, setPage] = useState(1)

  useEffect(() => {
    setPage(1)
  }, [searchTerm])

  const { user, setIsLoginModalOpen } = useAuth()

  const handleBorrowClick = (item: Item) => {
    if (!user) {
      setIsLoginModalOpen(true)
      return
    }
    setSelectedItem(item)
    setIsBorrowModalOpen(true)
  }

  const { data, isLoading, isError } = useQuery<PaginatedResponse<Item>>({
    queryKey: ['items', page, searchTerm],
    queryFn: () => itemService.getAll(page, 12, searchTerm),
    placeholderData: keepPreviousData,
  })

  if (isLoading) {
    return <Loading message={'Loading inventory...'} />
  }

  if (isError) {
    return <Error message={'Something went wrong.'} />
  }

  return (
    <>
      <div className="p-4">
        <h1 className="text-2xl font-bold mb-4">Inventory System</h1>

        <SearchBar value={searchTerm} onChange={setSearchTerm} placeholder="Search Items..." />

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {data?.data?.map((item: Item) => (
            <ItemCard key={item.id} item={item} handleBorrow={() => handleBorrowClick(item)} />
          ))}
        </div>

        <div className="flex justify-center items-center gap-4 mt-8 pb-8">
          <button
            disabled={page === 1 || isLoading}
            onClick={() => setPage((p) => Math.max(p - 1, 1))}
            className="px-4 py-2 border rounded-lg text-sm font-medium hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <span className="text-sm text-slate-600">
            Page {page} of {data?.total_pages || 1}
          </span>
          <button
            disabled={page >= (data?.total_pages || 1) || isLoading}
            onClick={() => setPage((p) => p + 1)}
            className="px-4 py-2 border rounded-lg text-sm font-medium hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
      </div>

      {/* item model */}
      {isBorrowModalOpen && selectedItem && (
        <div
          className="fixed inset-0 z-[100] flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4"
          onClick={() => setIsBorrowModalOpen(false)}
        >
          <div
            className="relative w-full max-w-md bg-white rounded-2xl shadow-2xl shadow-slate-900/20 overflow-hidden"
            onClick={(e) => e.stopPropagation()}
          >
            {/* close X button */}
            <button
              className="absolute top-4 right-4 text-slate-400 hover:text-slate-600 transition-colors"
              onClick={() => setIsBorrowModalOpen(false)}
            >
              <X size={20} />
            </button>

            <BorrowCard item={selectedItem} setIsBorrowModalOpen={setIsBorrowModalOpen} />
          </div>
        </div>
      )}
    </>
  )
}

export default HomePage
