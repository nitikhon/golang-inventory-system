import React from 'react'
import { Box, History, Package, User, Globe } from 'lucide-react'

const NavBar: React.FC = () => {
    return (
        <nav className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-slate-200">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex justify-between h-16 items-center">
          
                    {/* left: logo, brand */}
                    <div className="flex items-center gap-8">
                        <div className="flex items-center gap-2 group cursor-pointer">
                            <div className="w-10 h-10 bg-blue-600 rounded-xl flex items-center justify-center text-white shadow-lg shadow-blue-200 group-hover:scale-110 transition-transform">
                                <Box size={24} />
                            </div>
                            <span className="text-xl font-bold text-slate-900 tracking-tight">INV.SYSTEM</span>
                        </div>

                        {/* menu */}
                        <div className="hidden md:flex items-center gap-1">
                            <NavItem icon={<Package size={18} />} label="Inventory" active />
                            <NavItem icon={<History size={18} />} label="My Borrowings" />
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
                        <div className="flex items-center gap-3 pl-2 cursor-pointer hover:opacity-80 transition-opacity">
                            <div className="hidden sm:block text-right">
                                <p className="text-xs font-bold text-slate-900">John Doe</p>
                                <p className="text-[10px] text-slate-500">username: johndoe</p>
                            </div>
                            <User size={32}/>
                        </div>
                    </div>
                </div>
            </div>
        </nav>
    )
}

// sub-component for navbar menu
const NavItem = ({ icon, label, active = false }: { icon: React.ReactNode, label: string, active?: boolean }) => (
    <a className={`
    flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium transition-all
    ${active 
        ? 'bg-blue-50 text-blue-600' 
        : 'text-slate-500 hover:bg-slate-50 hover:text-slate-900'}
  `}>
        {icon}
        {label}
    </a>
)

export default NavBar