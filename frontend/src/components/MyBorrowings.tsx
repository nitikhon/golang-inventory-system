// MyBorrowings.tsx (Skeleton)
import { useQuery } from '@tanstack/react-query'
import borrowingService from '../services/borrowing'
import { useAuth } from '../hooks/useAuth'
import type React from 'react'
import Loading from './Loading'
import Error from './Error'
import formatDate from '../utils/formatDate'

const MyBorrowings: React.FC = () => {
  const { user, token, setIsLoginModalOpen } = useAuth()
  
  const { data: borrowings, isLoading, isError } = useQuery({
    queryKey: ['my-borrowings', user?.id],
    queryFn: () => borrowingService.getBorrowingsByUserID(token?.access_token),
    enabled: !!user?.id,
  })
  
  if (!user) { 
    setIsLoginModalOpen(true)
    return <p>Please login</p> 
  }

  if (isLoading) {
    return (
        <Loading message={"Loading borrowings"}/>
    )
  }

  if (isError) {
    return (
        <Error message={"Failed to load borrowings. Please check if your backend is running."}/>
    )
  }

  if (borrowings?.length === 0) {
    return (
        <div>
            No borrowings found. You can borrow items at home page.
        </div>
    )
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-slate-800">My Borrowing History</h1>
      
      <div className="bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead className="bg-slate-50 border-b border-slate-100">
             <tr>
                <th scope="col">Item Name</th>
                <th scope="col">Date</th>
                <th scope="col">Amount</th>
                <th scope="col">Status</th>
            </tr>
          </thead>
          <tbody>
            {borrowings?.map((borrow) => (
              <tr key={borrow.id} className="border-b border-slate-50 last:border-0">
                <th scope="row">{borrow.item?.name ?? 'Unknown Item'}</th>
                <td>{formatDate(borrow.borrowed_at)}</td>
                <td>{borrow.borrowing_amount}</td>
                <td>{borrow.borrowing_status}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default MyBorrowings