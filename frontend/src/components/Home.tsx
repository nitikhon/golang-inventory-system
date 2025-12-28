import itemService from '../services/item'
import { useQuery } from '@tanstack/react-query'
import type { Item } from '../types/item'
import ItemCard from './ItemCard'

const HomePage: React.FC = () => {
  const { data, isLoading, isError } = useQuery({
    queryKey: ['items'],
    queryFn: itemService.getAll,
  })

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span className="ml-3 text-slate-500">Loading inventory...</span>
      </div>
    )
  }

  if (isError) {
    return (
      <div className="p-8 bg-red-50 text-red-600 rounded-2xl border border-red-100 text-center">
        Failed to load items. Please check if your backend is running.
      </div>
    )
  }

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Inventory System</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {data?.map((item: Item) => (
          <ItemCard key={item.id} item={item} />
        ))}
      </div>
    </div>
  )
}

export default HomePage
