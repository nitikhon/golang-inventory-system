import { Box, History, Package, User, Globe, X } from 'lucide-react'
import { Link, NavLink } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import LoginCard from './LoginCard'

const NavBar: React.FC = () => {
  const { user, isLoginModalOpen, setIsLoginModalOpen } = useAuth()

  return (
    <>
      <nav className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-slate-200">
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
                    INV.SYSTEM
                  </span>
                </Link>
              </div>

              {/* menu */}
              <div className="hidden md:flex items-center gap-1">
                <NavItem icon={<Package size={18} />} label="Inventory" to="/" />
                {user 
                  ? <NavItem icon={<History size={18} />} label="My Borrowings" to="/my-borrowings" />
                  : <GhostNavItem icon={<History size={18} />} label="My Borrowings" onClick={() => setIsLoginModalOpen(true)} />
                }
              </div>
            </div>

            {/* right: languages, profile */}
            <div className="flex items-center gap-4">
              <button className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-slate-600 hover:bg-slate-100 rounded-lg transition-colors">
                <Globe size={16} />
                <span>EN</span>
              </button>

              <div className="h-6 w-[1px] bg-slate-200 mx-1"></div>

              {/* profile */}

              {user ? (
                <div className="flex items-center gap-3 pl-2 cursor-pointer hover:opacity-80 transition-opacity">
                  <div className="hidden sm:block text-right">
                    <p className="text-xs font-bold text-slate-900">{`${user?.first_name} ${user?.last_name}`}</p>
                    <p className="text-[10px] text-slate-500">username: {`${user?.username}`}</p>
                  </div>

                  <NavItem icon={<User size={32} />} to="/profile" />
                </div>
              ) : (
                <button onClick={() => setIsLoginModalOpen(true)}>Login</button>
              )}
            </div>
          </div>
        </div>
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

// sub-component for navbar menu
const NavItem = ({ icon, label, to }: { icon: React.ReactNode, label?: string, to: string, active?: boolean }) => (
  <NavLink
    to={to}
    className={({ isActive }) => `
    flex items-center ${label ? 'gap-2' : ''} px-4 py-2 rounded-xl text-sm font-medium transition-all
    ${isActive ? 'bg-blue-50 text-blue-600' : 'text-slate-500 hover:bg-slate-50 hover:text-slate-900'}
  `}
  >
    {icon}
    {label}
  </NavLink>
)

const GhostNavItem = ({ icon, label, onClick }: { icon: React.ReactNode, label: string, onClick: () => void }) => (
  <button 
    onClick={onClick}
    className={`
    flex items-center ${label ? 'gap-2' : ''} px-4 py-2 rounded-xl text-sm font-medium transition-all`}
  >
    {icon}
    <span className="text-sm font-medium">{label}</span>
  </button>
)

export default NavBar
