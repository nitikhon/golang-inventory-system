// MyBorrowings.tsx (Skeleton)
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import borrowingService from '../services/borrowing'
import { useAuth } from '../hooks/useAuth'
import type React from 'react'
import Loading from './Loading'
import Error from './Error'
import formatDate from '../utils/formatDate'
import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Activity, Clock, PackageCheck, PackageSearch } from 'lucide-react'
import StatCard from './StatCard'
import Swal from 'sweetalert2'
import toast from 'react-hot-toast'

const MyBorrowings: React.FC = () => {
  const { user, token, isSilentLoading, setIsLoginModalOpen } = useAuth()

  const navigate = useNavigate()

  const queryClient = useQueryClient()

  const {
    data: borrowings,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ['my-borrowings', user?.id],
    queryFn: () => borrowingService.getBorrowingsByUserID(token?.access_token),
    enabled: !!user?.id,
  })

  const { data } = useQuery({
    queryKey: ['stats'],
    queryFn: () => borrowingService.getBorrowingStats(token?.access_token),
    enabled: !!token?.access_token,
  })

  const { mutate: cancelMutation } = useMutation({
    mutationFn: (id: number) => borrowingService.cancelBorrowing(id, token?.access_token),
    onSuccess: () => {
      toast.success('Request cancelled successfully!', { duration: 5000 })
      queryClient.invalidateQueries({ queryKey: ['my-borrowings'] })
      queryClient.invalidateQueries({ queryKey: ['stats'] })
    },
  })

  const handleCancel = (id: number) => {
    Swal.fire({
      title: 'Are you sure?',
      text: "You won't be able to revert this!",
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#3085d6',
      cancelButtonColor: '#d33',
      confirmButtonText: 'Yes, cancel it!',
    }).then((result) => {
      if (result.isConfirmed) {
        cancelMutation(id)
      }
    })
  }

  const getStatusStyles = (status: string) => {
    const baseClasses =
      'px-2.5 py-1 rounded-full text-xs font-semibold border inline-flex items-center shadow-sm'

    switch (status.toLowerCase()) {
      case 'pending':
        return `${baseClasses} bg-amber-50 text-amber-700 border-amber-200`
      case 'active':
        return `${baseClasses} bg-blue-50 text-blue-700 border-blue-200`
      case 'returned':
        return `${baseClasses} bg-emerald-50 text-emerald-700 border-emerald-200`
      case 'overdue':
        return `${baseClasses} bg-rose-50 text-rose-700 border-rose-200`
      default:
        return `${baseClasses} bg-slate-50 text-slate-600 border-slate-200`
    }
  }

  useEffect(() => {
    if (!isSilentLoading && !user) {
      setIsLoginModalOpen(true)
    }
  }, [user, setIsLoginModalOpen])

  if (!user) {
    return <p>Please login</p>
  }

  if (isLoading || isSilentLoading) {
    return <Loading message={'Loading borrowings'} />
  }

  if (isError) {
    return <Error message={'Failed to load borrowings. Please check if your backend is running.'} />
  }

  if (borrowings?.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center p-12 bg-white rounded-3xl border border-dashed ...">
        <PackageSearch size={20} />
        <h3 className="mt-4 text-lg font-bold">No borrowings yet</h3>
        <p className="text-slate-500">Items you borrow will appear here.</p>
        <button
          className="mt-6 p-3 rounded-lg bg-blue-500 text-white"
          onClick={() => navigate('/')}
        >
          Back to Home Page
        </button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <StatCard
          icon={<Clock className="text-amber-500" />}
          label="Ongoing Borrows"
          value={String(data?.ongoing_borrows ?? 0)}
        />
        <StatCard
          icon={<PackageCheck className="text-emerald-500" />}
          label="Total Returned"
          value={String(data?.total_returned ?? 0)}
        />
        <StatCard
          icon={<Activity className="text-blue-500" />}
          label="Currently Borrows"
          value={String(data?.currently_borrows ?? 0)}
        />
      </div>

      <h1 className="text-2xl font-bold text-slate-800">My Borrowing History</h1>

      <div className="bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden">
        {/* Desktop / Ipad / Widescreen */}
        <table className="hidden md:table w-full text-left border-collapse">
          <thead className="bg-slate-50 border-b border-slate-100">
            <tr className="hover:bg-slate-50/50">
              <th className="py-4 px-6" scope="col">
                Item Name
              </th>
              <th className="py-4 px-6" scope="col">
                Date
              </th>
              <th className="py-4 px-6" scope="col">
                Amount
              </th>
              <th className="py-4 px-6" scope="col">
                Status
              </th>
              <th className="py-4 px-6" scope="col">
                Action
              </th>
            </tr>
          </thead>
          <tbody>
            {borrowings?.map((borrow) => (
              <tr key={borrow.id} className="border-b border-slate-50 last:border-0">
                <th className="py-4 px-6" scope="row">
                  {borrow.item?.name ?? 'Unknown Item'}
                </th>
                <td>{formatDate(borrow.borrowed_at)}</td>
                <td>{borrow.borrowing_amount}</td>
                <td className="py-4 px-6">
                  <span className={getStatusStyles(borrow.borrowing_status)}>
                    <span className="w-1.5 h-1.5 rounded-full bg-current mr-1.5 opacity-60"></span>
                    {borrow.borrowing_status}
                  </span>
                </td>
                <td className="py-4 px-6">
                  {borrow.borrowing_status === 'pending' && (
                    <button
                      onClick={() => handleCancel(borrow.id)}
                      className="w-full py-2 text-sm text-red-600 font-medium bg-red-50 rounded-lg 
                    active:bg-red-100 transition-colors border border-red-700"
                    >
                      Cancel
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {/* Mobile */}
        <div className="md:hidden divide-y divide-slate-100">
          {borrowings?.map((borrow) => (
            <div key={borrow.id} className="p-4 space-y-3">
              {/* Header: name + status */}
              <div className="flex justify-between items-start">
                <span className="font-semibold text-slate-900">{borrow.item?.name}</span>
                <span className={getStatusStyles(borrow.borrowing_status)}>
                  <span className="w-1.5 h-1.5 rounded-full bg-current mr-1.5 opacity-60"></span>
                  {borrow.borrowing_status}
                </span>
              </div>
              {/* Details*/}
              <div className="grid grid-cols-2 gap-2 text-sm text-slate-500">
                <div>
                  <p className="text-xs text-slate-400">Date</p>
                  <p>{formatDate(borrow.borrowed_at)}</p>
                </div>
                <div>
                  <p className="text-xs text-slate-400">Amount</p>
                  <p>{borrow.borrowing_amount}</p>
                </div>
              </div>
              {/* Action Button Section */}
              {borrow.borrowing_status === 'pending' && (
                <div className="pt-2 border-t border-slate-50">
                  <button
                    onClick={() => handleCancel(borrow.id)}
                    className="w-full py-2 text-sm text-red-600 font-medium bg-red-50 rounded-lg 
                    active:bg-red-100 transition-colors border border-red-700"
                  >
                    Cancel Request
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export default MyBorrowings
