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
