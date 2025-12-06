export function formatRelativeTime(dateString: string): string {
  if (!dateString) {
    return ''
  }
  const date = new Date(dateString)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  const now = new Date()
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)

  if (diffInSeconds < 60) {
    return '刚刚'
  }

  const diffInMinutes = Math.floor(diffInSeconds / 60)
  if (diffInMinutes < 60) {
    return `${diffInMinutes}分钟前`
  }

  const diffInHours = Math.floor(diffInMinutes / 60)
  if (diffInHours < 24) {
    return `${diffInHours}小时前`
  }

  const diffInDays = Math.floor(diffInHours / 24)
  if (diffInDays < 30) {
    return `${diffInDays}天前`
  }

  const diffInMonths = Math.floor(diffInDays / 30)
  if (diffInMonths < 12) {
    return `${diffInMonths}个月前`
  }

  return `${Math.floor(diffInMonths / 12)}年前`
}

export function formatPrice(price: number): string {
  return `¥${price.toFixed(2)}`
}

/**
 * 将数字数组转换为逗号分隔的字符串
 * 用于 tagIds 和 conditionIds 参数
 */
export function arrayToCommaSeparated(arr: number[]): string {
  return arr.join(',')
}

/**
 * 将逗号分隔的字符串转换为数字数组
 */
export function commaSeparatedToArray(str: string): number[] {
  if (!str || str.trim() === '') {
    return []
  }
  return str
    .split(',')
    .map((s) => parseInt(s.trim(), 10))
    .filter((n) => !isNaN(n))
}
