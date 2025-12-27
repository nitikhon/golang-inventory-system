import React from 'react'
import type { Item } from '../types/item'

interface ItemCardProps {
  item: Item;
}

const ItemCard: React.FC<ItemCardProps> = ({ item }) => {
    const isAvailable = item.status === 'available'
    const isBorrowed = item.status === 'borrowed'
    const canBorrow = (item.status === 'available' || item.status === 'borrowed') && item.available_amount > 0
    const getButtonText = () => {
        switch (item.status) {
        case 'available':
            return item.available_amount > 0 ? 'Borrow Item' : 'Out of Stock'
        case 'borrowed':
            return item.available_amount > 0 ? 'Borrow More' : 'Out of Stock'
        case 'maintenance':
            return 'Under Maintenance'
        case 'lost':
            return 'Item Lost'
        default:
            return 'Unavailable'
        }
    }
    const getBadgeColor = () => {
        switch (item.status) {
        case 'available':
            return 'bg-emerald-100 text-emerald-700'
        case 'borrowed':
            return 'bg-rose-100 text-rose-700'
        case 'maintenance':
            return 'bg-amber-100 text-amber-700'
        case 'lost':
            return 'bg-yellow-100 text-black'
        default:
            return 'bg-slate-100 text-slate-600'
        }
    }

    return (
        <div className="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden hover:shadow-md transition-shadow group">
            <div className="p-5">
                <div className="flex justify-between items-start mb-4">
                    {/* badge */}
                    <span className={`text-[10px] font-bold uppercase tracking-wider px-2 py-1 rounded-md ${
                        getBadgeColor()
                    }`}>
                        {item.status}
                    </span>

                    {/* item amount */}
                    {(isAvailable || isBorrowed) && <span className="text-sm font-medium text-slate-400">
                        {item.available_amount}/{item.total_amount} items
                    </span>}
                </div>
                
                {/* item info */}
                <h3 className="text-lg font-semibold text-slate-800 mb-1 group-hover:text-blue-600">
                    {item.name}
                </h3>
                <p className="text-sm text-slate-500 line-clamp-2 min-h-[40px]">
                    {item.description}
                </p>
                
                {/* button */}
                <div className="mt-5 pt-4 border-t border-slate-50">
                    <button 
                        disabled={!canBorrow}
                        className="w-full py-2.5 rounded-xl font-medium transition-all active:scale-95 disabled:opacity-50 
                        disabled:active:scale-100 bg-blue-700 text-white hover:bg-slate-800 disabled:bg-slate-200 
                        disabled:text-slate-400"
                    >
                        {getButtonText()}
                    </button>
                </div>
            </div>
        </div>
    )
}

export default ItemCard