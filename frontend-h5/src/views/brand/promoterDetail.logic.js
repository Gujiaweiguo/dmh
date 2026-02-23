/**
 * PromoterDetail page logic functions
 */

/**
 * Map promoter status to display text
 * @param {string} status - Status key
 * @returns {string} Display text
 */
export function getPromoterStatusText(status) {
  const map = {
    active: '活跃',
    inactive: '不活跃',
    blocked: '已禁用'
  }
  return map[status] || status || '-'
}

/**
 * Format promoter detail data for display
 * @param {Object} promoter - Raw promoter data
 * @returns {Object} Formatted promoter data
 */
export function formatPromoterDetail(promoter) {
  if (!promoter) return {}

  return {
    ...promoter,
    displayStatus: getPromoterStatusText(promoter.status),
    displayLevel: promoter.level || '普通',
    displayTotalOrders: promoter.totalOrders || 0,
    displayTotalRewards: promoter.totalRewards || 0,
    displayConversionRate: promoter.conversionRate || 0,
    displayCampaignCount: promoter.campaignCount || 0
  }
}

/**
 * Build reward records route path
 * @param {string|number} promoterId - Promoter ID
 * @returns {string} Route path
 */
export function buildRewardRecordsPath(promoterId) {
  if (!promoterId) return '/brand/reward-records'
  return `/brand/reward-records/${promoterId}`
}

/**
 * Build promoter link generation path
 * @param {string|number} promoterId - Promoter ID
 * @returns {string} Route path with query
 */
export function buildGenerateLinkPath(promoterId) {
  if (!promoterId) return '/brand/promoters'
  return `/brand/promoters?generateLink=${promoterId}`
}
