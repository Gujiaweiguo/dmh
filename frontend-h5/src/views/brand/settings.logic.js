export const getDefaultBrandInfo = () => ({
  name: '示例品牌',
  description: '专业的数字营销服务品牌',
  logo: 'https://api.dicebear.com/7.x/initials/svg?seed=Brand',
  phone: '400-123-4567',
  email: 'contact@brand.com',
})

export const getDefaultRewardSettings = () => ({
  defaultRate: 20,
  minWithdraw: 100,
  settlementType: 'instant',
})

export const getDefaultNotificationSettings = () => ({
  newOrder: true,
  newPromoter: true,
  dailyReport: false,
  email: 'admin@brand.com',
})

export const getDefaultSyncSettings = () => ({
  status: 'connected',
  frequency: 'realtime',
  dataTypes: ['orders', 'users'],
})

export const getDefaultPasswordForm = () => ({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

export const getSyncStatusText = (status) => {
  const statusMap = {
    connected: '已连接',
    disconnected: '未连接',
    error: '连接错误',
  }
  return statusMap[status] || status
}

export const validatePasswordForm = (passwordForm) => {
  if (!passwordForm.oldPassword || !passwordForm.newPassword) {
    return '请填写完整的密码信息'
  }

  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    return '两次输入的新密码不一致'
  }

  return ''
}

export const getCurrentBrandId = (storedBrandId, userInfoRaw) => {
  if (storedBrandId) {
    const parsed = Number.parseInt(storedBrandId, 10)
    if (!Number.isNaN(parsed) && parsed > 0) {
      return parsed
    }
  }

  if (!userInfoRaw) return 0

  try {
    const userInfo = JSON.parse(userInfoRaw)
    const firstBrandId = Array.isArray(userInfo.brandIds) ? Number(userInfo.brandIds[0]) : 0
    return Number.isFinite(firstBrandId) && firstBrandId > 0 ? firstBrandId : 0
  } catch {
    return 0
  }
}

export const unwrapApiResponse = (payload) => {
  if (payload && typeof payload === 'object' && payload.data && typeof payload.data === 'object') {
    return payload.data
  }
  return payload || {}
}

export const resolveSyncStatus = (syncHealthPayload) => {
  const status = String(syncHealthPayload?.status || '').toLowerCase()
  if (status === 'healthy' || status === 'connected' || status === 'ok') {
    return 'connected'
  }
  return 'error'
}

export const validateLogoFile = (file) => {
  if (!file) return '请选择图片文件'

  if (!String(file.type || '').startsWith('image/')) {
    return '仅支持图片文件上传'
  }

  const maxSize = 5 * 1024 * 1024
  if (file.size > maxSize) {
    return '图片大小不能超过 5MB'
  }

  return ''
}

export const resolveUploadedLogoUrl = (payload) => {
  const body = unwrapApiResponse(payload)
  return body.url || body.fileUrl || body.logo || ''
}

const escapeCsvCell = (value) => {
  const text = String(value ?? '')
  if (text.includes(',') || text.includes('"') || text.includes('\n')) {
    return `"${text.replace(/"/g, '""')}"`
  }
  return text
}

export const buildBrandDataCsv = (brandInfo, rewardSettings, notificationSettings) => {
  const rows = [
    { category: '品牌信息', field: '品牌名称', value: brandInfo?.name || '' },
    { category: '品牌信息', field: '品牌描述', value: brandInfo?.description || '' },
    { category: '品牌信息', field: '联系电话', value: brandInfo?.phone || '' },
    { category: '品牌信息', field: '联系邮箱', value: brandInfo?.email || '' },
    { category: '品牌信息', field: 'Logo', value: brandInfo?.logo || '' },
    { category: '奖励设置', field: '默认佣金比例(%)', value: String(rewardSettings?.defaultRate || 0) },
    { category: '奖励设置', field: '最低提现金额', value: String(rewardSettings?.minWithdraw || 0) },
    { category: '奖励设置', field: '结算方式', value: rewardSettings?.settlementType || '' },
    { category: '通知设置', field: '新订单通知', value: notificationSettings?.newOrder ? '是' : '否' },
    { category: '通知设置', field: '新推广员通知', value: notificationSettings?.newPromoter ? '是' : '否' },
    { category: '通知设置', field: '日报通知', value: notificationSettings?.dailyReport ? '是' : '否' },
    { category: '通知设置', field: '通知邮箱', value: notificationSettings?.email || '' },
  ]

  const headers = ['分类', '字段', '值']
  const lines = [headers.map(escapeCsvCell).join(',')]

  rows.forEach((row) => {
    lines.push([row.category, row.field, row.value].map(escapeCsvCell).join(','))
  })

  return lines.join('\n')
}

export const downloadCsv = (csvContent, filename) => {
  if (!csvContent) return false

  const blob = new Blob(['\uFEFF' + csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename || 'export.csv'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)

  return true
}

export const getExportFilename = (prefix = 'brand-data') => {
  const now = new Date()
  const dateStr = now.toISOString().slice(0, 10)
  return `${prefix}-${dateStr}.csv`
}
