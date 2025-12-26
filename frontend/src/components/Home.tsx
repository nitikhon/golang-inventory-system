// import React, { useState } from 'react';
// import { InventoryItem } from './types/inventory';

import itemService from '../services/item'
import { useQuery } from '@tanstack/react-query'
import type { Item } from '../types/item';
import ItemCard from './ItemCard';

const HomePage: React.FC = () => {
    const queryItems = useQuery({ queryKey: ['items'], queryFn: itemService.getAll })

    return (
        <div className="p-4">
            <h1 className="text-2xl font-bold mb-4">Inventory System</h1>
      
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">  
                    {queryItems.data?.map((item: Item) => (
                        <ItemCard item={item}/>
                    ))}
            </div>
        </div>
    );
};

export default HomePage;