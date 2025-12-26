import React from 'react';
import type { Item } from '../types/item';

interface ItemCardProps {
  item: Item;
}

const ItemCard: React.FC<ItemCardProps> = ({ item }) => {
    return (
        <div className="border rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow">
            <div>{item?.name}</div>
            <div>{item?.description}</div>
            <div>{item?.status}</div>
        </div>
    );
};

export default ItemCard;