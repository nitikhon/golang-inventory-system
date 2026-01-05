import type React from 'react'
import capitalizeSentence from '../utils/capitalizeSentence'

interface ErrorProps {
  message: string
}

const Error: React.FC<ErrorProps> = ({ message }) => {
  return (
    <div className="p-8 bg-red-50 text-red-600 rounded-2xl border border-red-100 text-center">
      {capitalizeSentence(message)}
    </div>
  )
}

export default Error
