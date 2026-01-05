import itemService from '../services/item'
import { useQuery } from '@tanstack/react-query'
import type { Item } from '../types/item'
import ItemCard from './ItemCard'
import Loading from './Loading'
import Error from './Error'
import { useAuth } from '../hooks/useAuth'
import { useState } from 'react'
import BorrowCard from './BorrowCard'
import { X } from 'lucide-react'
import SearchBar from './SearchBar'

const HomePage: React.FC = () => {
  const [selectedItem, setSelectedItem] = useState<Item | null>(null)
  const [isBorrowModalOpen, setIsBorrowModalOpen] = useState(false)
  const [searchTerm, setSearchTerm] = useState('')

  const { user, setIsLoginModalOpen } = useAuth()

  const handleBorrowClick = (item: Item) => {
    if (!user) {
      setIsLoginModalOpen(true)
      return
    }
    setSelectedItem(item)
    setIsBorrowModalOpen(true)
  }

  const { data, isLoading, isError } = useQuery({
    queryKey: ['items'],
    queryFn: itemService.getAll,
  })

  const filteredItems = data?.filter(
    (item: Item) =>
      item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      item.description.toLowerCase().includes(searchTerm.toLowerCase())
  )

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

        <SearchBar value={searchTerm} onChange={setSearchTerm} />

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {searchTerm == ''
            ? data?.map((item: Item) => (
                <ItemCard key={item.id} item={item} handleBorrow={() => handleBorrowClick(item)} />
              ))
            : filteredItems?.map((item: Item) => (
                <ItemCard key={item.id} item={item} handleBorrow={() => handleBorrowClick(item)} />
              ))}
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
