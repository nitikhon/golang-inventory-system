import { NavLink } from 'react-router-dom'

const NavItem = ({
  icon,
  label,
  to,
  onClick,
}: {
  icon: React.ReactNode
  label?: string
  to: string
  onClick?: () => void
}) => (
  <NavLink
    to={to}
    className={({ isActive }) => `
    flex items-center ${label ? 'gap-2' : ''} px-4 py-2 rounded-xl text-sm font-medium transition-all
    ${isActive ? 'bg-blue-50 text-blue-600' : 'text-slate-500 hover:bg-slate-50 hover:text-slate-900'}
  `}
    onClick={onClick}
  >
    {icon}
    {label}
  </NavLink>
)

export default NavItem
