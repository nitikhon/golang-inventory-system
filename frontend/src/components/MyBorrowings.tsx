import { useMutation, useQuery, useQueryClient, keepPreviousData } from '@tanstack/react-query'
import borrowingService from '../services/borrowing'
import { useAuth } from '../hooks/useAuth'
import type React from 'react'
import Loading from './Loading'
import Error from './Error'
import formatDate from '../utils/formatDate'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Activity, Clock, PackageCheck, PackageSearch } from 'lucide-react'
import StatCard from './StatCard'
import Swal from 'sweetalert2'
import toast from 'react-hot-toast'
import Pagination from './Pagination'
import SearchBar from './SearchBar'
import useDebounce from '../hooks/useDebounce'
import { useTranslation } from '../hooks/useTranslation'
import StatusBadge from './StatusBadge'

const MyBorrowings: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState('')
  const [page, setPage] = useState(1)

  const { t } = useTranslation()
  const { user, token, isSilentLoading, setIsLoginModalOpen } = useAuth()

  const navigate = useNavigate()

  const queryClient = useQueryClient()

  const debouncedSearch = useDebounce(searchTerm, 500)

  useEffect(() => {
    setPage(1)
  }, [debouncedSearch])

  const {
    data: borrowings,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ['my-borrowings', debouncedSearch, page, user?.id],
    queryFn: () =>
      borrowingService.getBorrowingsByUserID(page, 12, debouncedSearch, token?.access_token),
    enabled: !!user?.id,
    placeholderData: keepPreviousData,
  })

  useEffect(() => {
    if (borrowings?.total_pages && page > borrowings.total_pages) {
      setPage(borrowings.total_pages)
    }
  }, [borrowings?.data, page])

  const { data } = useQuery({
    queryKey: ['stats'],
    queryFn: () => borrowingService.getBorrowingStats(token?.access_token),
    enabled: !!token?.access_token,
  })

  const { mutate: cancelMutation } = useMutation({
    mutationFn: (id: number) => borrowingService.cancelBorrowing(id, token?.access_token),
    onSuccess: () => {
      toast.success(t.myBorrowings.messages.cancelSuccess, { duration: 5000 })
      queryClient.invalidateQueries({ queryKey: ['my-borrowings'] })
      queryClient.invalidateQueries({ queryKey: ['stats'] })
    },
  })

  const handleCancel = (id: number) => {
    Swal.fire({
      title: t.admin.dialog.title,
      text: t.admin.dialog.text,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#3085d6',
      cancelButtonColor: '#d33',
      confirmButtonText: t.admin.dialog.confirmReject,
    }).then((result) => {
      if (result.isConfirmed) {
        cancelMutation(id)
      }
    })
  }

  useEffect(() => {
    if (!isSilentLoading && !user) {
      setIsLoginModalOpen(true)
    }
  }, [user, setIsLoginModalOpen])

  if (!user) {
    return <p>{t.nav.login}</p>
  }

  if (isLoading || isSilentLoading) {
    return <Loading message={t.inventory.loading} />
  }

  if (isError) {
    return <Error message={t.common.error} />
  }

  if (borrowings?.total_items === 0 && !searchTerm) {
    return (
      <div className="flex flex-col items-center justify-center p-12 bg-white rounded-3xl border border-dashed ...">
        <PackageSearch size={20} />
        <h3 className="mt-4 text-lg font-bold">{t.myBorrowings.empty.title}</h3>
        <p className="text-slate-500">{t.myBorrowings.empty.description}</p>
        <button
          className="mt-6 p-3 rounded-lg bg-blue-500 text-white"
          onClick={() => navigate('/')}
        >
          {t.myBorrowings.empty.backHome}
        </button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <StatCard
          icon={<Clock className="text-amber-500" />}
          label={t.myBorrowings.stats.ongoing}
          value={String(data?.ongoing_borrows ?? 0)}
        />
        <StatCard
          icon={<PackageCheck className="text-emerald-500" />}
          label={t.myBorrowings.stats.returned}
          value={String(data?.total_returned ?? 0)}
        />
        <StatCard
          icon={<Activity className="text-blue-500" />}
          label={t.myBorrowings.stats.current}
          value={String(data?.currently_borrows ?? 0)}
        />
      </div>

      <h1 className="text-2xl font-bold text-slate-800">{t.myBorrowings.title}</h1>

      <SearchBar
        value={searchTerm}
        onChange={setSearchTerm}
        placeholder={t.myBorrowings.searchPlaceholder}
      />

      <div className="bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden">
        {/* Desktop / Ipad / Widescreen */}
        <table className="hidden md:table w-full text-left border-collapse">
          <thead className="bg-slate-50 border-b border-slate-100">
            <tr className="hover:bg-slate-50/50">
              <th className="py-4 px-6" scope="col">
                {t.myBorrowings.table.itemName}
              </th>
              <th className="py-4 px-6" scope="col">
                {t.myBorrowings.table.date}
              </th>
              <th className="py-4 px-6" scope="col">
                {t.myBorrowings.table.amount}
              </th>
              <th className="py-4 px-6" scope="col">
                {t.myBorrowings.table.status}
              </th>
              <th className="py-4 px-6" scope="col">
                {t.myBorrowings.table.action}
              </th>
            </tr>
          </thead>
          <tbody>
            {borrowings?.data?.map((borrow) => (
              <tr key={borrow.id} className="border-b border-slate-50 last:border-0">
                <th className="py-4 px-6" scope="row">
                  {borrow.item?.name ?? t.common.unknownItem}
                </th>
                <td>{formatDate(borrow.borrowed_at)}</td>
                <td>{borrow.borrowing_amount}</td>
                <td className="py-4 px-6">
                  <StatusBadge borrowing_status={borrow.borrowing_status} />
                </td>
                <td className="py-4 px-6">
                  {borrow.borrowing_status === 'pending' && (
                    <button
                      onClick={() => handleCancel(borrow.id)}
                      className="w-full py-2 text-sm text-red-600 font-medium bg-red-50 rounded-lg 
                    active:bg-red-100 transition-colors border border-red-700"
                    >
                      {t.myBorrowings.actions.cancel}
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {/* Mobile */}
        <div className="md:hidden divide-y divide-slate-100">
          {borrowings?.data?.map((borrow) => (
            <div key={borrow.id} className="p-4 space-y-3">
              {/* Header: name + status */}
              <div className="flex justify-between items-start">
                <span className="font-semibold text-slate-900">{borrow.item?.name}</span>
                <StatusBadge borrowing_status={borrow.borrowing_status} />
              </div>
              {/* Details*/}
              <div className="grid grid-cols-2 gap-2 text-sm text-slate-500">
                <div>
                  <p className="text-xs text-slate-400">{t.myBorrowings.table.date}</p>
                  <p>{formatDate(borrow.borrowed_at)}</p>
                </div>
                <div>
                  <p className="text-xs text-slate-400">{t.myBorrowings.table.amount}</p>
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
                    {t.myBorrowings.actions.cancelRequest}
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>

        {borrowings?.total_items === 0 && (
          <div className="flex flex-col items-center justify-center p-12 text-slate-500">
            <PackageSearch size={48} className="text-slate-200 mb-4" strokeWidth={1.5} />
            <p>
              {t.common.noResults} "{searchTerm}"
            </p>
          </div>
        )}
      </div>

      {(borrowings?.total_items || 0) > 0 && (
        <Pagination
          page={page}
          totalPages={borrowings?.total_pages || 1}
          isLoading={isLoading}
          handlePageChange={setPage}
        />
      )}
    </div>
  )
}

export default MyBorrowings
