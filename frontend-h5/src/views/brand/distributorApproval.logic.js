/**
 * DistributorApproval business logic
 * Extracted from DistributorApproval.vue for testability
 */

/**
 * Format time string to localized display format
 * @param {string} timeString - ISO time string
 * @returns {string} Formatted time or '-' if invalid
 */
export const formatTime = (timeString) => {
  if (!timeString) return '-'
  const date = new Date(timeString)
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

/**
 * Get display text for application status
 * @param {string} status - Application status key
 * @returns {string} Human-readable status text
 */
export const getStatusText = (status) => {
  const map = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return map[status] || status
}

/**
 * Get default approval form values
 * @returns {Object} Default form object
 */
export const getDefaultApprovalForm = () => ({
  level: 1,
  notes: ''
})

/**
 * Get default pagination state
 * @returns {Object} Default pagination object
 */
export const getDefaultPagination = () => ({
  currentPage: 1,
  pageSize: 20,
  hasMore: false
})

/**
 * Tab configuration builder
 * @param {number} pendingCount - Number of pending applications
 * @param {number} processedCount - Number of processed applications
 * @returns {Array} Tab configuration array
 */
export const buildTabs = (pendingCount = 0, processedCount = 0) => [
  { key: 'pending', label: '待审批', count: pendingCount },
  { key: 'processed', label: '已处理', count: processedCount }
]

/**
 * Validate rejection form (must have notes)
 * @param {Object} form - Approval form object
 * @returns {boolean} True if valid for rejection
 */
export const validateRejection = (form) => {
  return !!(form && form.notes && form.notes.trim())
}

/**
 * Build approve request payload
 * @param {number} applicationId - Application ID
 * @param {Object} form - Approval form with level and notes
 * @returns {Object} Request payload
 */
export const buildApprovePayload = (applicationId, form) => ({
  applicationId,
  action: 'approve',
  level: form?.level || 1,
  reason: form?.notes || ''
})

/**
 * Build reject request payload
 * @param {number} applicationId - Application ID
 * @param {string} reason - Rejection reason
 * @returns {Object} Request payload
 */
export const buildRejectPayload = (applicationId, reason) => ({
  applicationId,
  action: 'reject',
  reason: reason || ''
})

/**
 * Get mock pending applications for development
 * @returns {Array} Mock application array
 */
export const getMockPendingApplications = () => [
  {
    id: 1,
    userId: 4,
    username: 'user004',
    phone: '138****4444',
    brandId: 1,
    brandName: '品牌A',
    status: 'pending',
    reason: '我想成为分销商，推广品牌产品',
    createdAt: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 2,
    userId: 5,
    username: 'user005',
    phone: '139****5555',
    brandId: 1,
    brandName: '品牌A',
    status: 'pending',
    reason: '有多年的销售经验，希望能加入',
    createdAt: new Date(Date.now() - 5 * 60 * 60 * 1000).toISOString()
  }
]

/**
 * Get avatar text (first character of username)
 * @param {string} username - Username
 * @param {string} fallback - Fallback text
 * @returns {string} Avatar text
 */
export const getAvatarText = (username, fallback = '申') => {
  return username?.charAt(0) || fallback
}

export const buildPaginationParams = (page, pageSize, status) => {
  const params = new URLSearchParams()
  params.set('page', String(page || 1))
  params.set('pageSize', String(pageSize || 20))
  if (status) {
    params.set('status', status)
  }
  return params.toString()
}

export const calculateHasMore = (currentListLength, pageSize, total) => {
  if (typeof total === 'number') {
    return currentListLength < total
  }
  return currentListLength >= pageSize
}

export const mergePaginatedResults = (existingList, newList, append = true) => {
  if (!append) {
    return newList || []
  }
  if (!Array.isArray(newList) || newList.length === 0) {
    return existingList || []
  }
  const existingIds = new Set((existingList || []).map(item => item.id))
  const uniqueNewItems = newList.filter(item => !existingIds.has(item.id))
  return [...(existingList || []), ...uniqueNewItems]
}
