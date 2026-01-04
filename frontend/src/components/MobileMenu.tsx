import { Globe, User, Package, History, ShieldUser } from 'lucide-react'
import type { User as UserType } from '../types/user'
import { Link } from 'react-router-dom'
import GhostNavItem from './GhostNavItem'
import NavItem from './NavItem'

interface MobileMenuProps {
  isOpen: boolean
  onClose: () => void
  user: UserType | null
  onLoginClick: () => void
}

const MobileMenu: React.FC<MobileMenuProps> = ({ isOpen, onClose, user, onLoginClick }) => {
  if (!isOpen) return null

  return (
    <div className="md:hidden absolute top-16 left-0 w-full bg-slate-50 border-b border-slate-200 shadow-xl z-40 p-4 flex flex-col gap-4">
      {/* Navigation Menus */}
      <div className="flex flex-col gap-2 border-b border-slate-200 pb-4">
        <NavItem icon={<Package size={18} />} label="Inventory" to="/" onClick={onClose} />
        {user ? (
          <NavItem
            icon={<History size={18} />}
            label="My Borrowings"
            to="/my-borrowings"
            onClick={onClose}
          />
        ) : (
          <GhostNavItem
            icon={<History size={18} />}
            label="My Borrowings"
            onClick={() => {
              onLoginClick
              onClose
            }}
          />
        )}
        {user?.is_admin && <NavItem icon={<ShieldUser size={18} />} label="Admin" to="/admin" />}
      </div>

      {/* Profile */}
      <div className="flex items-center justify-between">
        {user ? (
          <div className="flex items-center gap-3">
            <div className="bg-slate-200 p-2 rounded-full">
              <User size={20} />
            </div>
            <div>
              <p className="font-bold text-sm">{user.first_name}</p>
              <Link to="/profile" className="text-xs text-blue-600" onClick={onClose}>
                View Profile
              </Link>
            </div>
          </div>
        ) : (
          <button
            onClick={() => {
              onLoginClick
              onClose
            }}
          >
            Sign In
          </button>
        )}

        {/* Language */}
        <button className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-slate-600 hover:bg-slate-100 rounded-lg transition-colors">
          <Globe size={16} />
          <span>EN</span>
        </button>
      </div>
    </div>
  )
}

export default MobileMenu
