export const ProductStatus = {
  FOR_SALE: 'ForSale',
  SOLD: 'Sold',
  DELISTED: 'Delisted',
} as const

export type ProductStatusType = (typeof ProductStatus)[keyof typeof ProductStatus]
