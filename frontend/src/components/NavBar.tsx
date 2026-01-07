import { Box, History, Package, User, X, ShieldUser } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import LoginCard from './LoginCard'
import { useState } from 'react'
import NavItem from './NavItem'
import GhostNavItem from './GhostNavItem'
import MobileMenu from './MobileMenu'
import { useTranslation } from '../contexts/LanguageContext'

const NavBar: React.FC = () => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

  const { user, isLoginModalOpen, setIsLoginModalOpen } = useAuth()

  const { language, setLanguage, t } = useTranslation()

  const toggleLang = () => {
      setLanguage(prev => prev === 'en' ? 'th' : 'en')
  }

  return (
    <>
      <nav className="relative sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-slate-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            {/* left: logo, brand */}
            <div className="flex items-center gap-8">
              <div className="flex items-center gap-2 group cursor-pointer">
                <div className="w-10 h-10 bg-blue-600 rounded-xl flex items-center justify-center text-white shadow-lg shadow-blue-200 group-hover:scale-110 transition-transform">
                  <Box size={24} />
                </div>
                <Link to="/">
                  <span className="text-xl font-bold text-slate-900 tracking-tight">
                    {t.nav.brand}
                  </span>
                </Link>
              </div>

              {/* desktop menu */}
              <div className="hidden md:flex items-center gap-1">
                <NavItem icon={<Package size={18} />} label={t.nav.inventory} to="/" />
                {user ? (
                  <NavItem icon={<History size={18} />} label={t.nav.myBorrowings} to="/my-borrowings" />
                ) : (
                  <GhostNavItem
                    icon={<History size={18} />}
                    label={t.nav.myBorrowings}
                    onClick={() => setIsLoginModalOpen(true)}
                  />
                )}
                {user?.is_admin && (
                  <NavItem icon={<ShieldUser size={18} />} label={t.nav.admin} to="/admin" />
                )}
              </div>
            </div>

            {/* right: languages, profile */}
            <div className="flex items-center gap-4 hidden md:flex">
              <button 
                className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-slate-600 hover:bg-slate-100 
                rounded-lg transition-colors"
                onClick={toggleLang}
              >
                <span>{language === 'en' ? '🇹🇭 TH' : '🇬🇧 EN'}</span>
              </button>

              <div className="h-6 w-[1px] bg-slate-200 mx-1"></div>

              {/* profile */}
              {user ? (
                <div className="flex items-center gap-3 pl-2 cursor-pointer hover:opacity-80 transition-opacity">
                  <div className="hidden sm:block text-right">
                    <p className="text-xs font-bold text-slate-900">{`${user?.first_name} ${user?.last_name}`}</p>
                    <p className="text-[10px] text-slate-500">{t.auth.username}: {`${user?.username}`}</p>
                  </div>

                  <NavItem icon={<User size={32} />} to="/profile" />
                </div>
              ) : (
                <button onClick={() => setIsLoginModalOpen(true)}>{t.auth.loginButton}</button>
              )}
            </div>

            {/* Hamburger Button (Visible ONLY on Mobile) */}
            <div className="md:hidden">
              <button
                onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                className="p-2 text-slate-600 hover:bg-slate-100 rounded-lg"
              >
                {/* Toggler */}
                {isMobileMenuOpen ? (
                  <X size={24} />
                ) : (
                  <div className="space-y-1.5">
                    <span className="block w-6 h-0.5 bg-current"></span>
                    <span className="block w-6 h-0.5 bg-current"></span>
                    <span className="block w-6 h-0.5 bg-current"></span>
                  </div>
                )}
              </button>
            </div>
          </div>
        </div>

        {/* Mobile Menu */}
        {isMobileMenuOpen && (
          <MobileMenu
            isOpen={isMobileMenuOpen}
            onClose={() => setIsMobileMenuOpen(false)}
            user={user}
            onLoginClick={() => setIsLoginModalOpen(true)}
          />
        )}
      </nav>

      {/* Login Modal */}
      {isLoginModalOpen && (
        <div
          className="fixed inset-0 z-[100] flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4"
          onClick={() => setIsLoginModalOpen(false)}
        >
          <div
            className="relative w-full max-w-md bg-white rounded-2xl shadow-2xl shadow-slate-900/20 overflow-hidden"
            onClick={(e) => e.stopPropagation()}
          >
            {/* close X button */}
            <button
              className="absolute top-4 right-4 text-slate-400 hover:text-slate-600 transition-colors"
              onClick={() => setIsLoginModalOpen(false)}
            >
              <X size={20} />
            </button>

            <LoginCard />
          </div>
        </div>
      )}
    </>
  )
}

export default NavBar
