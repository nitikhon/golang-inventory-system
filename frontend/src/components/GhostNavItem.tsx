const GhostNavItem = ({
  icon,
  label,
  onClick,
}: {
  icon: React.ReactNode
  label: string
  onClick: () => void
}) => (
  <button
    onClick={onClick}
    className={`
    flex items-center ${label ? 'gap-2' : ''} px-4 py-2 rounded-xl text-sm font-medium transition-all 
    text-slate-500 hover:bg-slate-50 hover:text-slate-900`}
  >
    {icon}
    <span className="text-sm font-medium">{label}</span>
  </button>
)

export default GhostNavItem
