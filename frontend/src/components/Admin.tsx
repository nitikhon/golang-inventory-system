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

const Admin: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'pending' | 'active' | 'returned'>('pending')

  const { user, token, isSilentLoading } = useAuth()

  const navigate = useNavigate()

  const queryClient = useQueryClient()

  const { data, isLoading, isError } = useQuery({
    queryKey: ['admin-borrowings', activeTab],
    queryFn: () => borrowingService.getBorrowingsByStatus(activeTab, token?.access_token),
    enabled: !!token?.access_token,
  })

  const { mutate: approveMutation, isPending: isApproving } = useMutation({
    mutationFn: (id: number) => borrowingService.approveBorrowing(id, token?.access_token),
    onSuccess: () => {
      toast.success('Request approved successfully!', { duration: 5000 })
      queryClient.invalidateQueries({ queryKey: ['admin-borrowings'] })
      queryClient.invalidateQueries({ queryKey: ['my-borrowings'] })
      queryClient.invalidateQueries({ queryKey: ['stats'] })
    },
  })

  const { mutate: cancelMutation, isPending: isRejecting } = useMutation({
    mutationFn: (id: number) => borrowingService.cancelBorrowing(id, token?.access_token),
    onSuccess: () => {
      toast.success('Request cancelled successfully!', { duration: 5000 })
      queryClient.invalidateQueries({ queryKey: ['admin-borrowings'] })
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
      confirmButtonText: 'Yes, reject it!',
    }).then((result) => {
      if (result.isConfirmed) {
        cancelMutation(id)
      }
    })
  }

  const handleApprove = (id: number) => {
    Swal.fire({
      title: 'Are you sure?',
      text: "You won't be able to revert this!",
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#3085d6',
      cancelButtonColor: '#d33',
      confirmButtonText: 'Yes, approve it!',
    }).then((result) => {
      if (result.isConfirmed) {
        approveMutation(id)
      }
    })
  }

  useEffect(() => {
    if (!isSilentLoading && (!user || !user.is_admin)) {
      navigate('/')
    }
  }, [user, isSilentLoading, navigate])

  const tabs = [
    { id: 'pending', label: 'Pending Requests' },
    { id: 'active', label: 'Active Borrows' },
    { id: 'returned', label: 'History' },
  ] as const

  if (isSilentLoading || !user?.is_admin) return null

  if (isLoading) {
    return <Loading message={'Loading inventory...'} />
  }

  if (isError) {
    return <Error message={'Failed to load items. Please check if your backend is running.'} />
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-slate-900">Admin Dashboard</h1>

      {/* Tabs UI */}
      <div className="flex gap-4 border-b border-slate-200">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
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
            <h3 className="font-semibold text-slate-700">Requests List</h3>
            <span className="text-xs text-slate-500">Total: {data?.length} items</span>
          </div>

          {/* Table */}
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="bg-slate-50 text-slate-500 font-medium">
                <tr>
                  <th className="px-6 py-3">User</th>
                  <th className="px-6 py-3">Item Detail</th>
                  <th className="px-6 py-3">Dates</th>
                  <th className="px-6 py-3 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {data?.map((req: Borrowing) => (
                  <tr key={req.id} className="hover:bg-slate-50/50 transition-colors">
                    {/* User Info */}
                    <td className="px-6 py-4">
                      <div className="font-medium text-slate-900">{req.user?.username}</div>
                    </td>

                    {/* Item Info */}
                    <td className="px-6 py-4">
                      <div className="font-medium text-slate-900">{req.item?.name}</div>
                      <div className="text-xs text-slate-500">Qty: {req.borrowing_amount}</div>
                    </td>

                    {/* Date */}
                    <td className="px-6 py-4 text-slate-500">
                      <div>Due: {formatDate(req.due_date)}</div>
                    </td>

                    {/* Actions (only Pending Tab) */}
                    <td className="px-6 py-4 text-right space-x-2">
                      {activeTab === 'pending' && (
                        <>
                          <button
                            disabled={isApproving}
                            className="text-green-600 hover:bg-green-50 px-3 py-1 rounded-lg border 
                            border-green-200 disabled:opacity-50 disabled:cursor-not-allowed"
                            onClick={() => handleApprove(req.id)}
                          >
                            {isApproving ? 'Approving...' : 'Approve'}
                          </button>
                          <button
                            disabled={isApproving}
                            className="text-red-600 hover:bg-red-50 px-3 py-1 rounded-lg border 
                            border-red-200 disabled:opacity-50 disabled:cursor-not-allowed"
                            onClick={() => handleCancel(req.id)}
                          >
                            {isRejecting ? 'Rejecting...' : 'Reject'}
                          </button>
                        </>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            {/* Empty State */}
            {(!data || data.length === 0) && (
              <div className="p-12 text-center text-slate-500">No data found in this tab.</div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default Admin
