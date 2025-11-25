export const ProductConditionCode = {
  BRAND_NEW: 'BRAND_NEW',
  NINE_TENTHS: 'NINE_TENTHS',
  EIGHT_TENTHS: 'EIGHT_TENTHS',
  SEVEN_TENTHS: 'SEVEN_TENTHS',
} as const

export type ProductConditionCodeType = (typeof ProductConditionCode)[keyof typeof ProductConditionCode]
