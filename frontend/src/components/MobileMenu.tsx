import { Globe, User, Package, History, ShieldUser } from 'lucide-react'
import type { User as UserType } from '../types/user'
import { Link } from 'react-router-dom'
import GhostNavItem from './GhostNavItem'
import NavItem from './NavItem'
import { useTranslation } from '../hooks/useTranslation'

interface MobileMenuProps {
  isOpen: boolean
  onClose: () => void
  user: UserType | null
  onLoginClick: () => void
}

const MobileMenu: React.FC<MobileMenuProps> = ({ isOpen, onClose, user, onLoginClick }) => {
  const { t, language, setLanguage } = useTranslation()

  const toggleLanguage = () => {
    setLanguage((prev) => (prev === 'en' ? 'th' : 'en'))
  }

  if (!isOpen) return null

  return (
    <div className="md:hidden absolute top-16 left-0 w-full bg-slate-50 border-b border-slate-200 shadow-xl z-40 p-4 flex flex-col gap-4">
      {/* Navigation Menus */}
      <div className="flex flex-col gap-2 border-b border-slate-200 pb-4">
        <NavItem icon={<Package size={18} />} label={t.nav.inventory} to="/" onClick={onClose} />
        {user ? (
          <NavItem
            icon={<History size={18} />}
            label={t.nav.myBorrowings}
            to="/my-borrowings"
            onClick={onClose}
          />
        ) : (
          <GhostNavItem
            icon={<History size={18} />}
            label={t.nav.myBorrowings}
            onClick={() => {
              onLoginClick
              onClose
            }}
          />
        )}
        {user?.is_admin && (
          <NavItem icon={<ShieldUser size={18} />} label={t.nav.admin} to="/admin" onClick={onClose} />
        )}
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
                {t.nav.profile}
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
            {t.auth.signIn}
          </button>
        )}

        {/* Language */}
        <button 
          onClick={toggleLanguage}
          className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-slate-600 hover:bg-slate-100 rounded-lg transition-colors">
          <Globe size={16} />
          <span>{language === 'en' ? 'EN' : 'TH'}</span>
        </button>
      </div>
    </div>
  )
}

export default MobileMenu
