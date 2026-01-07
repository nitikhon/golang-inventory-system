import itemService from '../services/item'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import type { Item } from '../types/item'
import ItemCard from './ItemCard'
import Loading from './Loading'
import Error from './Error'
import { useAuth } from '../hooks/useAuth'
import { useEffect, useState } from 'react'
import BorrowCard from './BorrowCard'
import { PackageSearch, X } from 'lucide-react'
import SearchBar from './SearchBar'
import type { PaginatedResponse } from '../types/pagination'
import Pagination from './Pagination'
import useDebounce from '../hooks/useDebounce'
import { useTranslation } from '../contexts/LanguageContext'

const HomePage: React.FC = () => {
  const [selectedItem, setSelectedItem] = useState<Item | null>(null)
  const [isBorrowModalOpen, setIsBorrowModalOpen] = useState(false)
  const [searchTerm, setSearchTerm] = useState('')
  const [page, setPage] = useState(1)

  const debouncedSearch = useDebounce(searchTerm, 500)
  const { t } = useTranslation()

  useEffect(() => {
    setPage(1)
  }, [debouncedSearch])

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
    queryKey: ['items', page, debouncedSearch],
    queryFn: () => itemService.getAll(page, 9, debouncedSearch),
    placeholderData: keepPreviousData,
  })

  if (isLoading) {
    return <Loading message={t.inventory.loading} />
  }

  if (isError) {
    return <Error message={t.common.error} />
  }

  if (data?.total_items === 0 && !searchTerm) {
    return (
      <div className="flex flex-col items-center justify-center p-12 bg-white rounded-3xl border border-dashed ...">
        <PackageSearch size={20} />
        <h3 className="mt-4 text-lg font-bold">{t.inventory.noItems}</h3>
        <p className="text-slate-500">{t.inventory.contactStaff}</p>
      </div>
    )
  }

  return (
    <>
      <div className="p-4">
        <h1 className="text-2xl font-bold mb-4">{t.inventory.title}</h1>

        <SearchBar value={searchTerm} onChange={setSearchTerm} placeholder={t.search.placeholder} />

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {data?.data?.map((item: Item) => (
            <ItemCard key={item.id} item={item} handleBorrow={() => handleBorrowClick(item)} />
          ))}
        </div>

        {data?.total_items === 0 && (
          <div className="flex flex-col items-center justify-center p-12 text-slate-500">
            <PackageSearch size={48} className="text-slate-200 mb-4" strokeWidth={1.5} />
            <p>{t.common.noResults} "{searchTerm}"</p>
          </div>
        )}

        {(data?.total_items || 0) > 0 && (
          <Pagination
            page={page}
            totalPages={data?.total_pages || 1}
            isLoading={isLoading}
            handlePageChange={setPage}
          />
        )}
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
