import itemService from '../services/item'
import { useQuery } from '@tanstack/react-query'
import type { Item } from '../types/item'
import ItemCard from './ItemCard'
import Loading from './Loading'
import Error from './Error'

const HomePage: React.FC = () => {
  const { data, isLoading, isError } = useQuery({
    queryKey: ['items'],
    queryFn: itemService.getAll,
  })

  if (isLoading) {
    return (
      <Loading message={"Loading inventory..."}/>
    )
  }

  if (isError) {
    return (
      <Error message={"Failed to load items. Please check if your backend is running."}/>
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
