import React from 'react'
import {
  User as UserIcon,
  Mail,
  Building,
  Calendar,
  PackageCheck,
  Clock,
  ShieldCheck,
  LogOut,
} from 'lucide-react'
import { useAuth } from '../hooks/useAuth'
import borrowingService from '../services/borrowing'
import { useQuery } from '@tanstack/react-query'
import formatDate from '../utils/formatDate'

const Profile: React.FC = () => {
  const { user, token, logout } = useAuth()

  const { data } = useQuery({
    queryKey: ['stats'],
    queryFn: () => borrowingService.getBorrowingStats(token?.access_token),
    enabled: !!token?.access_token,
  })

  if (!user) return <p>Please login</p>

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header Section */}
      <div className="bg-white rounded-3xl p-8 border border-slate-200 shadow-sm flex flex-col md:flex-row items-center gap-6">
        <div className="w-24 h-24 bg-blue-100 rounded-full flex items-center justify-center text-blue-600">
          <UserIcon size={48} strokeWidth={1.5} />
        </div>
        <div className="text-center md:text-left space-y-1">
          <h1 className="text-2xl font-bold text-slate-900">
            {user?.first_name} {user?.last_name}
          </h1>
          <p className="text-slate-500 font-medium">@{user?.username}</p>
          <div className="flex flex-wrap justify-center md:justify-start gap-2 mt-2">
            <span className="px-3 py-1 bg-blue-50 text-blue-600 text-xs font-bold rounded-full uppercase tracking-wider">
              {user?.is_admin ? 'Admin' : 'Standard User'}
            </span>
          </div>
        </div>
      </div>

      {/* stats grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
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
          icon={<ShieldCheck className="text-blue-500" />}
          label="Account Status"
          value={user.deleted_at === null ? 'Active' : 'Unactive'}
        />
      </div>

      {/* detailed info section */}
      <div className="bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="px-8 py-4 border-b border-slate-100 bg-slate-50/50">
          <h2 className="font-semibold text-slate-800">Account Details</h2>
        </div>
        <div className="p-8 space-y-6">
          <InfoRow icon={<Mail size={18} />} label="Email Address" value={user?.email} />
          <InfoRow icon={<Building size={18} />} label="Phone" value={user?.phone} />
          <InfoRow
            icon={<Calendar size={18} />}
            label="Member Since"
            value={formatDate(user?.created_at)}
          />
        </div>
      </div>

      {/* action (logout) */}
      <div className="pt-4">
        <button
          onClick={logout}
          className="w-full flex items-center justify-center gap-2 px-6 py-4 bg-rose-50 text-rose-600 
                    font-bold rounded-2xl border border-rose-100 hover:bg-rose-100 hover:shadow-md hover:shadow-rose-100 
                    transition-all active:scale-[0.98]"
        >
          <LogOut size={20} />
          Sign Out from System
        </button>
        <p className="text-center text-slate-400 text-xs mt-4">
          Your session will be cleared and you will need to login again.
        </p>
      </div>
    </div>
  )
}

// Sub-components
const StatCard = ({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode
  label: string
  value: string
}) => (
  <div className="bg-white p-6 rounded-3xl border border-slate-200 shadow-sm flex items-center gap-4">
    <div className="p-3 bg-slate-50 rounded-2xl">{icon}</div>
    <div>
      <p className="text-sm font-medium text-slate-500">{label}</p>
      <p className="text-xl font-bold text-slate-900">{value}</p>
    </div>
  </div>
)

const InfoRow = ({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode
  label: string
  value: string
}) => (
  <div className="flex items-center gap-4">
    <div className="text-slate-400">{icon}</div>
    <div className="flex-1">
      <p className="text-xs font-medium text-slate-400 uppercase tracking-wider">{label}</p>
      <p className="text-slate-900 font-medium">{value}</p>
    </div>
  </div>
)

export default Profile
