/**
 * RewardRecords page logic functions
 */

/**
 * Map reward type to display text
 * @param {string} type - Type key
 * @returns {string} Display text
 */
export function getRewardTypeText(type) {
  const map = {
    commission: '佣金',
    bonus: '奖金',
    penalty: '扣款'
  }
  return map[type] || type
}

/**
 * Map reward status to display text
 * @param {string} status - Status key
 * @returns {string} Display text
 */
export function getRewardStatusText(status) {
  const map = {
    pending: '待发放',
    paid: '已发放',
    cancelled: '已取消'
  }
  return map[status] || status
}

/**
 * Calculate reward records summary
 * @param {Array} records - Reward records array
 * @returns {Object} Summary object with totalAmount, pendingAmount, totalCount
 */
export function calculateRewardSummary(records) {
  if (!Array.isArray(records) || records.length === 0) {
    return {
      totalAmount: '0.00',
      pendingAmount: '0.00',
      totalCount: 0
    }
  }

  const totalAmount = records.reduce((sum, r) => sum + (r.amount || 0), 0)
  const pendingAmount = records
    .filter(r => r.status === 'pending')
    .reduce((sum, r) => sum + (r.amount || 0), 0)

  return {
    totalAmount: totalAmount.toFixed(2),
    pendingAmount: pendingAmount.toFixed(2),
    totalCount: records.length
  }
}

/**
 * Filter reward records by type and status
 * @param {Array} records - Reward records array
 * @param {string} type - Type filter ('all' or specific type)
 * @param {string} status - Status filter ('all' or specific status)
 * @returns {Array} Filtered records
 */
export function filterRewardRecords(records, type, status) {
  if (!Array.isArray(records)) return []

  let result = records

  if (type && type !== 'all') {
    result = result.filter(r => r.type === type)
  }

  if (status && status !== 'all') {
    result = result.filter(r => r.status === status)
  }

  return result
}

/**
 * Format amount with sign
 * @param {number} amount - Amount value
 * @returns {string} Formatted amount with sign
 */
export function formatRewardAmount(amount) {
  if (amount === null || amount === undefined) return '¥0'
  const sign = amount >= 0 ? '+' : ''
  return `${sign}¥${amount}`
}

/**
 * Determine amount CSS class
 * @param {number} amount - Amount value
 * @returns {string} CSS class name
 */
export function getAmountClass(amount) {
  return amount >= 0 ? 'positive' : 'negative'
}
