import type React from 'react'
import { useState, useEffect } from 'react'
import { useAuth } from '../hooks/useAuth'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import itemService from '../services/item'
import type { Item, ItemPatchRequest } from '../types/item'
import toast from 'react-hot-toast'
import { X } from 'lucide-react'
import { useTranslation } from '../hooks/useTranslation'
import type { SweetAlertIcon } from 'sweetalert2'
import Swal from 'sweetalert2'

interface EditItemModalProps {
  isOpen: boolean
  onClose: React.Dispatch<React.SetStateAction<boolean>>
  item: Item
}

const EditItemModal: React.FC<EditItemModalProps> = ({ isOpen, onClose, item }) => {
  const [form, setForm] = useState<ItemPatchRequest>({
    name: item.name,
    description: item.description,
    total_amount: item.total_amount,
    available_amount: item.available_amount,
    status: item.status,
  })

  useEffect(() => {
    if (isOpen) {
      setForm({
        name: item.name,
        description: item.description,
        total_amount: item.total_amount,
        available_amount: item.available_amount,
        status: item.status,
      })
    }
  }, [item, isOpen])

  const { user, token } = useAuth()

  const queryClient = useQueryClient()

  const { t } = useTranslation()

  const { mutate: updateItemMutation, isPending } = useMutation({
    mutationFn: (data: ItemPatchRequest) => itemService.patchUpdate(item.id, data, token?.access_token),
    onSuccess: () => {
      toast.success(t.editItemModal.success)
      onClose(false)
      queryClient.invalidateQueries({ queryKey: ['items'] })
    },
    onError: (error: any) => {
        const errorMsg = error.response?.data?.error || t.editItemModal.error.updateFailed
        swalFire('Error', errorMsg, 'error')
    }
  })

  const swalFire = (title: string, text: string, icon: SweetAlertIcon) => {
    Swal.fire({
      title: title,
      text: text,
      icon: icon,
    })
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (!form.name || form.name === '') {
      swalFire('Error', t.editItemModal.error.emptyName, 'error')
      return
    }

    if (!form.description || form.description === '') {
      swalFire('Error', t.editItemModal.error.emptyDesc, 'error')
      return
    }

    if (!form.total_amount || !form.available_amount) {
      swalFire('Error', t.editItemModal.error.invalidTotalAmount, 'error')
      return
    }

    if (form.total_amount <= 0) {
      swalFire('Error', t.editItemModal.error.invalidTotalAmount, 'error')
      return
    }
    
    if (form.available_amount < 0) {
        swalFire('Error', t.editItemModal.error.negativeAvailable, 'error')
        return
    }

    if (form.total_amount < form.available_amount) {
        swalFire('Error', t.editItemModal.error.invalidAvailable, 'error')
        return
    }

    if (user && token) {
      updateItemMutation(form)
    }
  }

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target

    setForm((prev) => ({
      ...prev,
      [name]: name === 'total_amount' || name === 'available_amount' ? Number(value) : value,
    }))
  }

  const inputClass =
    'w-full p-2.5 bg-slate-50 border border-slate-300 text-slate-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 transition-all font-medium'
  const labelClass = 'block mb-1.5 text-sm font-semibold text-slate-700'

  if (!isOpen) return null

  return (
    <div
      className="fixed inset-0 z-[100] flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4"
      onClick={() => onClose(false)}
    >
      <div
        className="relative w-full max-w-md bg-white rounded-2xl shadow-2xl shadow-slate-900/20 overflow-hidden"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex justify-between items-center p-5 border-b border-slate-100 bg-slate-50/50">
          <h3 className="text-lg font-bold text-slate-800">{t.editItemModal.header}</h3>
          <button
            onClick={() => onClose(false)}
            className="text-slate-400 hover:text-slate-600 p-1 rounded-full hover:bg-slate-100 transition-colors"
          >
            <X size={20} />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-6 space-y-5">
          <div>
            <label className={labelClass}>{t.editItemModal.labels.name}</label>
            <input
              type="text"
              name="name"
              className={inputClass}
              value={form.name}
              onChange={handleChange}
              placeholder={t.editItemModal.placeholder.name}
            />
          </div>
          <div>
            <label className={labelClass}>{t.editItemModal.labels.description}</label>
            <textarea
              name="description"
              rows={3}
              className={inputClass}
              value={form.description}
              onChange={handleChange}
              placeholder={t.editItemModal.placeholder.description}
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className={labelClass}>{t.editItemModal.labels.total_amount}</label>
              <input
                type="number"
                name="total_amount"
                min="1"
                className={inputClass}
                value={form.total_amount}
                onChange={handleChange}
              />
            </div>
             <div>
              <label className={labelClass}>{t.editItemModal.labels.available_amount}</label>
              <input
                type="number"
                name="available_amount"
                min="0"
                className={inputClass}
                value={form.available_amount}
                onChange={handleChange}
              />
            </div>
          </div>
          <div>
             <label className={labelClass}>{t.editItemModal.labels.status}</label>
              <select
                name="status"
                className={inputClass}
                value={form.status}
                onChange={handleChange}
              >
                <option value="available">{t.inventory.status.available}</option>
                <option value="maintenance">{t.inventory.status.maintenance}</option>
                <option value="lost">{t.inventory.status.lost}</option>
              </select>
          </div>

          <button
            type="submit"
            disabled={isPending}
            className="w-full py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-bold rounded-xl shadow-lg shadow-blue-500/30 transition-all active:scale-[0.98]"
          >
            {isPending ? t.editItemModal.actions.updating : t.editItemModal.actions.update}
          </button>
        </form>
      </div>
    </div>
  )
}

export default EditItemModal
