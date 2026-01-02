import type React from 'react'

interface LoadingProps {
  message: string
}

const Loading: React.FC<LoadingProps> = ({ message }) => {
  return (
    <div className="flex justify-center items-center min-h-[400px]">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      <span className="ml-3 text-slate-500">{message}</span>
    </div>
  )
}

export default Loading
