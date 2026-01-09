import { useEffect, useState } from 'react'
import { useAuth } from '../hooks/useAuth'
import { useNavigate } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import borrowingService from '../services/borrowing'
import Loading from './Loading'
import Error from './Error'
import toast from 'react-hot-toast'
import Swal from 'sweetalert2'
import type { Borrowing } from '../types/borrowing'
import formatDate from '../utils/formatDate'
import SearchBar from './SearchBar'
import capitalizeSentence from '../utils/capitalizeSentence'
import type { PaginatedResponse } from '../types/pagination'
import Pagination from './Pagination'
import useDebounce from '../hooks/useDebounce'
import { useTranslation } from '../hooks/useTranslation'
import StatusBadge from './StatusBadge'

const Admin: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'pending' | 'active' | 'history'>('pending')
  const [searchTerm, setSearchTerm] = useState('')
  const [page, setPage] = useState(1)

  const { t } = useTranslation()
  const { user, token, isSilentLoading } = useAuth()

  const navigate = useNavigate()

  const queryClient = useQueryClient()

  const debouncedSearch = useDebounce(searchTerm, 500)

  useEffect(() => {
    setPage(1)
  }, [activeTab, debouncedSearch])

  const { data, isLoading, isError } = useQuery<PaginatedResponse<Borrowing>>({
    queryKey: ['admin-borrowings', page, activeTab, debouncedSearch],
    queryFn: () => {
      let status = String(activeTab)
      if (activeTab === 'history') {
        status = 'returned,cancelled,lost,rejected'
      }
      if (activeTab === 'active') {
        status = 'active,overdue'
      }
      return borrowingService.getBorrowingsByStatus(
        page,
        12,
        debouncedSearch,
        status,
        token?.access_token
      )
    },
    enabled: !!token?.access_token,
  })

  const { mutate: approveMutation, isPending: isApproving } = useMutation({
    mutationFn: (id: number) => borrowingService.approveBorrowing(id, token?.access_token),
    onSuccess: () => {
      toast.success(t.admin.messages.approveSuccess, { duration: 5000 })
      queryClient.invalidateQueries({ queryKey: ['admin-borrowings'] })
      queryClient.invalidateQueries({ queryKey: ['my-borrowings'] })
      queryClient.invalidateQueries({ queryKey: ['stats'] })
    },
  })

  const { mutate: cancelMutation, isPending: isRejecting } = useMutation({
    mutationFn: (id: number) => borrowingService.cancelBorrowing(id, token?.access_token),
    onSuccess: () => {
      toast.success(t.admin.messages.cancelSuccess, { duration: 5000 })
      queryClient.invalidateQueries({ queryKey: ['admin-borrowings'] })
      queryClient.invalidateQueries({ queryKey: ['my-borrowings'] })
      queryClient.invalidateQueries({ queryKey: ['stats'] })
    },
  })

  const { mutate: returnMutation, isPending: isReturning } = useMutation({
    mutationFn: (id: number) => borrowingService.returnBorrowing(id, token?.access_token),
    onSuccess: () => {
      toast.success(t.admin.messages.returnSuccess, { duration: 5000 })
      queryClient.invalidateQueries({ queryKey: ['admin-borrowings'] })
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

  const handleApprove = (id: number) => {
    Swal.fire({
      title: t.admin.dialog.title,
      text: t.admin.dialog.text,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#3085d6',
      cancelButtonColor: '#d33',
      confirmButtonText: t.admin.dialog.confirmApprove,
    }).then((result) => {
      if (result.isConfirmed) {
        approveMutation(id)
      }
    })
  }

  const handleReturn = (id: number) => {
    Swal.fire({
      title: t.admin.dialog.title,
      text: t.admin.dialog.text,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#3085d6',
      cancelButtonColor: '#d33',
      confirmButtonText: t.admin.dialog.confirmReturn,
    }).then((result) => {
      if (result.isConfirmed) {
        returnMutation(id)
      }
    })
  }

  useEffect(() => {
    if (!isSilentLoading && (!user || !user.is_admin)) {
      navigate('/')
    }
  }, [user, isSilentLoading, navigate])

  const tabs = [
    { id: 'pending', label: t.admin.tabs.pending },
    { id: 'active', label: t.admin.tabs.active },
    { id: 'history', label: t.admin.tabs.history },
  ] as const

  if (isSilentLoading || !user?.is_admin) return null

  if (isLoading) {
    return <Loading message={t.inventory.loading} />
  }

  if (isError) {
    return <Error message={t.common.error} />
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-slate-900">{t.admin.title}</h1>

      <SearchBar
        value={searchTerm}
        onChange={setSearchTerm}
        placeholder={t.admin.searchPlaceholder}
      />

      {/* Tabs UI */}
      <div className="flex gap-4 border-b border-slate-200">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id as any)}
            className={`pb-2 px-1 border-b-2 font-medium text-sm transition-colors ${
              activeTab === tab.id
                ? 'border-blue-600 text-blue-600' // Active State
                : 'border-transparent text-slate-500 hover:text-slate-700' // Inactive State
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Content Area */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-6">
        <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
          {/* Header + Filters (TODO) */}
          <div className="p-4 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
            <h3 className="font-semibold text-slate-700">{t.admin.requestsList}</h3>
            <span className="text-xs text-slate-500">
              {t.inventory.labels.total}: {data?.data?.length} {t.inventory.labels.items}
            </span>
          </div>

          {/* Table */}
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="bg-slate-50 text-slate-500 font-medium">
                <tr>
                  <th className="px-6 py-3">{t.admin.table.user}</th>
                  <th className="px-6 py-3">{t.admin.table.itemDetail}</th>
                  <th className="px-6 py-3">{t.admin.table.dates}</th>
                  <th className="px-6 py-3">{t.admin.table.details}</th>

                  {(activeTab === 'history' || activeTab === 'active') && (
                    <th className="px-6 py-3">{t.admin.table.status}</th>
                  )}

                  {(activeTab === 'pending' || activeTab === 'active') && (
                    <th className="px-6 py-3">{t.admin.table.actions}</th>
                  )}
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {data?.data?.map((req: Borrowing) => (
                  <tr key={req.id} className="hover:bg-slate-50/50 transition-colors">
                    {/* User Info */}
                    <td className="px-6 py-4">
                      <div className="font-medium text-slate-900">{req.user?.username}</div>
                      <div className="text-[10px] text-slate-600">{`${capitalizeSentence(req.user?.first_name)} ${capitalizeSentence(req.user?.last_name)}`}</div>
                    </td>

                    {/* Item Info */}
                    <td className="px-6 py-4">
                      <div className="font-medium text-slate-900">{req.item?.name}</div>
                      <div className="text-xs text-slate-500">
                        {t.admin.table.qty}: {req.borrowing_amount}
                      </div>
                    </td>

                    {/* Date */}
                    <td className="px-6 py-4 text-slate-500">
                      <div>
                        {t.admin.table.due}: {formatDate(req.due_date)}
                      </div>
                    </td>

                    <td
                      className="px-6 py-4 text-slate-500 text-sm line-clamp-2 max-w-xs"
                      title={req.description}
                    >
                      <div>{req.description}</div>
                    </td>

                    {/* status */}
                    {(activeTab === 'history' || activeTab === 'active') && (
                      <td className="px-6 py-4 space-x-2">
                        <StatusBadge borrowing_status={req.borrowing_status} />
                      </td>
                    )}

                    {/* Actions only Pending and Active Tab if history turns into status */}
                    {(activeTab === 'pending' || activeTab === 'active') && (
                      <td className="px-6 py-4 space-x-2">
                        {/* approve or reject */}
                        {activeTab === 'pending' && (
                          <>
                            <button
                              disabled={isApproving}
                              className="text-green-600 hover:bg-green-50 px-3 py-1 rounded-lg border 
                            border-green-200 disabled:opacity-50 disabled:cursor-not-allowed"
                              onClick={() => handleApprove(req.id)}
                            >
                              {isApproving ? t.admin.actions.approving : t.admin.actions.approve}
                            </button>
                            <button
                              disabled={isApproving}
                              className="text-red-600 hover:bg-red-50 px-3 py-1 rounded-lg border 
                            border-red-200 disabled:opacity-50 disabled:cursor-not-allowed"
                              onClick={() => handleCancel(req.id)}
                            >
                              {isRejecting ? t.admin.actions.rejecting : t.admin.actions.reject}
                            </button>
                          </>
                        )}

                        {/* mark as returned */}
                        {activeTab === 'active' && (
                          <>
                            <button
                              disabled={isReturning}
                              className="text-green-600 hover:bg-green-50 px-3 py-1 rounded-lg border 
                            border-green-200 disabled:opacity-50 disabled:cursor-not-allowed"
                              onClick={() => handleReturn(req.id)}
                            >
                              {isReturning ? t.admin.actions.returning : t.admin.actions.return}
                            </button>
                          </>
                        )}
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>

            {/* Empty State */}
            {(!data?.data || data?.data.length === 0) && (
              <div className="p-12 text-center text-slate-500">{t.common.noResults}</div>
            )}
          </div>
        </div>
      </div>

      <Pagination
        page={page}
        totalPages={data?.total_pages || 1}
        isLoading={isLoading}
        handlePageChange={setPage}
      />
    </div>
  )
}

export default Admin
