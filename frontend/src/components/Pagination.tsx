interface PaginationProps {
  page: number
  totalPages: number
  isLoading: boolean
  handlePageChange: React.Dispatch<React.SetStateAction<number>>
}

const Pagination: React.FC<PaginationProps> = ({
  page,
  totalPages,
  isLoading,
  handlePageChange,
}) => {
  return (
    <div className="flex justify-center items-center gap-4 mt-8 pb-8">
      <button
        disabled={page === 1 || isLoading}
        onClick={() => handlePageChange((p: number) => Math.max(p - 1, 1))}
        className="px-4 py-2 border rounded-lg text-sm font-medium hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Previous
      </button>
      <span className="text-sm text-slate-600">
        Page {page} of {totalPages || 1}
      </span>
      <button
        disabled={page >= (totalPages || 1) || isLoading}
        onClick={() => handlePageChange((p: number) => p + 1)}
        className="px-4 py-2 border rounded-lg text-sm font-medium hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Next
      </button>
    </div>
  )
}

export default Pagination
