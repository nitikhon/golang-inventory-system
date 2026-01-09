import type React from 'react'
import type { Item } from '../types/item'
import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import borrowService from '../services/borrowing'
import { useAuth } from '../hooks/useAuth'
import type { AxiosError } from 'axios'
import type { BorrowRequest } from '../types/borrowing'
import { useNavigate } from 'react-router-dom'
import addDays from '../utils/addDays'
import toast from 'react-hot-toast'
import capitalizeSentence from '../utils/capitalizeSentence'
import { useTranslation } from '../hooks/useTranslation'
import Swal, { type SweetAlertIcon } from 'sweetalert2'

interface BorrowCardProps {
  item: Item
  setIsBorrowModalOpen: React.Dispatch<React.SetStateAction<boolean>>
}

const BorrowCard: React.FC<BorrowCardProps> = ({ item, setIsBorrowModalOpen }) => {
  const [borrowForm, setBorrowForm] = useState({
    amount: item.available_amount > 0 ? 1 : 0,
    description: '',
    due_date: addDays(new Date(), 7).toISOString().split('T')[0],
  })

  const { t } = useTranslation()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { user, token } = useAuth()

  const { mutate: borrowsMutation, isPending } = useMutation({
    mutationFn: (data: BorrowRequest) => borrowService.createBorrowing(data, token?.access_token),
    onSuccess: async () => {
      queryClient.invalidateQueries({ queryKey: ['items'] })
      queryClient.invalidateQueries({ queryKey: ['my-borrowings'] })
      queryClient.invalidateQueries({ queryKey: ['stats'] })
      toast.success(t.myBorrowings.messages.borrowSuccess, { duration: 5000 })
      navigate('/my-borrowings')
      setIsBorrowModalOpen(false)
    },
    onError: (error: AxiosError<{ error: string }>) => {
      let message = error.response?.data?.error || t.common.error
      const status = error.response?.status
      if (status == 500) {
        message = t.common.error
      }
      toast.error(capitalizeSentence(message), { duration: 5000 })
    },
  })

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target

    setBorrowForm((prev) => ({
      ...prev,
      [name]: value,
    }))
  }

  const swalFire = (title: string, text: string, icon: SweetAlertIcon) => {
    Swal.fire({
      title: title,
      text: text,
      icon: icon,
    })
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (borrowForm.amount <= 0) {
      swalFire('Error', t.inventory.borrowErr.negativeAmount, 'error')
      return
    }

    if (borrowForm.amount > item.available_amount) {
      swalFire('Error', t.inventory.borrowErr.invalidAmoumt, 'error')
      return
    }

    if (new Date(borrowForm.due_date) <= new Date()) {
      swalFire('Error', t.inventory.borrowErr.invalidDate, 'error')
      return
    }

    if (user && token) {
      const payload: BorrowRequest = {
        item_id: item?.id,
        borrowing_amount: Number(borrowForm.amount),
        description: borrowForm.description,
        due_date: new Date(borrowForm.due_date).toISOString(),
      }

      borrowsMutation(payload)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="p-6 space-y-4">
      {/* Header: item's name, amount*/}
      <div>
        <h2 className="text-xl font-bold">{item.name}</h2>
        <p className="text-sm text-slate-500">
          {t.inventory.labels.available}: {item.available_amount} {t.inventory.labels.items}
        </p>
      </div>
      {/* Input: borrow amount */}
      <div>
        <label className="block text-sm font-medium">{t.admin.table.qty}</label>
        <input
          name="amount"
          type="number"
          min={1}
          max={item.available_amount}
          value={borrowForm.amount}
          onChange={handleChange}
          className="w-full p-2 border rounded-lg"
        />
      </div>
      {/* Input: due date */}
      <div>
        <label className="block text-sm font-medium">{t.admin.table.due}</label>
        <input
          name="due_date"
          type="date"
          value={borrowForm.due_date}
          onChange={handleChange}
          className="w-full p-2 border rounded-lg"
        />
      </div>
      {/* Input: description */}
      <div>
        <label className="block text-sm font-medium">{t.inventory.labels.description}</label>
        <textarea
          name="description"
          value={borrowForm.description}
          onChange={handleChange}
          className="w-full p-2 border rounded-lg"
          placeholder={t.inventory.placeholders.borrowReason}
        />
      </div>
      {/* submit */}
      <button
        type="submit"
        disabled={isPending}
        className="w-full py-3 bg-blue-600 text-white rounded-xl font-bold hover:bg-blue-700"
      >
        {isPending ? t.common.processing : t.inventory.actions.borrow}
      </button>
    </form>
  )
}

export default BorrowCard
