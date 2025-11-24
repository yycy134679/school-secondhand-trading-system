export interface Product {
  id: number
  title: string
  description: string
  price: number
  mainImageUrl: string
  status: ProductStatus
  conditionId: number
  sellerId: number
  categoryId: number
  createdAt: string
  updatedAt: string
}

export type ProductStatus = 'ForSale' | 'Sold' | 'Delisted'
