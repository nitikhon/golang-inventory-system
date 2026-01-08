import { useTranslation } from '../hooks/useTranslation'

interface StatusBadgeProps {
  borrowing_status: string
}

const StatusBadge: React.FC<StatusBadgeProps> = ({ borrowing_status }) => {
  const { t } = useTranslation()

  const getStatusText = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return t.inventory.borrowingStatus.pending
      case 'active':
        return t.inventory.borrowingStatus.active
      case 'returned':
        return t.inventory.borrowingStatus.returned
      case 'overdue':
        return t.inventory.borrowingStatus.overdue
      case 'rejected':
        return t.inventory.borrowingStatus.rejected
      case 'cancelled':
        return t.inventory.borrowingStatus.cancelled
      case 'lost':
        return t.inventory.borrowingStatus.lost
      default:
        return status
    }
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
      case 'rejected':
        return `${baseClasses} bg-red-50 text-red-700 border-red-200`
      case 'cancelled':
        return `${baseClasses} bg-slate-100 text-slate-500 border-slate-200`
      case 'lost':
        return `${baseClasses} bg-purple-50 text-purple-700 border-purple-200`
      default:
        return `${baseClasses} bg-slate-50 text-slate-600 border-slate-200`
    }
  }

  return (
    <span className={getStatusStyles(borrowing_status)}>
      <span className="w-1.5 h-1.5 rounded-full bg-current mr-1.5 opacity-60"></span>
      {getStatusText(borrowing_status)}
    </span>
  )
}

export default StatusBadge
