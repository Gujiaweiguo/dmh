import { api } from './api.js'

// 认证相关API
export const authApi = {
  // 品牌管理员登录
  login: (username, password) => {
    return api.post('/auth/login', { username, password })
  },

  // 获取当前用户信息
  getUserInfo: () => {
    return api.get('/auth/userinfo')
  },

  // 登出
  logout: () => {
    return api.post('/auth/logout')
  }
}

// 活动管理API
export const campaignApi = {
  // 获取活动列表
  getCampaigns: (params = {}) => {
    return api.get('/campaigns', params)
  },

  // 获取活动详情
  getCampaign: (id) => {
    return api.get(`/campaign/detail/${id}`)
  },

  // 创建活动
  createCampaign: (data) => {
    return api.post('/campaign/create', data)
  },

  // 更新活动
  updateCampaign: (id, data) => {
    return api.put(`/campaign/update/${id}`, data)
  },

  // 删除活动
  deleteCampaign: (id) => {
    return api.delete(`/campaign/delete/${id}`)
  },

  // 保存页面配置
  savePageConfig: (campaignId, config) => {
    return api.post(`/campaign/page-config/${campaignId}`, config)
  },

  // 获取页面配置
  getPageConfig: (campaignId) => {
    return api.get(`/campaign/page-config/${campaignId}`)
  },

  // 获取支付二维码
  getPaymentQrcode: (campaignId) => {
    return api.get(`/campaigns/${campaignId}/payment-qrcode`)
  }
}

// 订单管理API
export const orderApi = {
  // 获取订单列表
  getOrders: (params = {}) => {
    return api.get('/order/list', params)
  },

  // 获取订单详情
  getOrder: (id) => {
    return api.get(`/order/detail/${id}`)
  },

  // 更新订单状态
  updateOrderStatus: (id, status) => {
    return api.put(`/order/status/${id}`, { status })
  },

  // 扫码获取订单信息
  scanOrderCode: (code) => {
    return api.get('/orders/scan', { code })
  },

  verifyOrder: (code, data = {}) => {
    const payload = typeof code === 'object' && code !== null
      ? { ...code }
      : { ...data }

    if (typeof code === 'string') {
      if (code.trim() !== '') {
        payload.code = code
      }
    } else if (code !== undefined && code !== null && typeof code !== 'object') {
      payload.code = code
    }

    return api.post('/orders/verify', payload)
  },

  unverifyOrder: (code, data = {}) => {
    const payload = typeof code === 'object' && code !== null
      ? { ...code }
      : { ...data }

    if (typeof code === 'string') {
      if (code.trim() !== '') {
        payload.code = code
      }
    } else if (code !== undefined && code !== null && typeof code !== 'object') {
      payload.code = code
    }

    return api.post('/orders/unverify', payload)
  },

  getVerificationRecords: () => {
    return api.get('/order/verification-records')
  }
}

// 推广员管理API
export const promoterApi = {
  // 获取推广员列表
  getPromoters: (params = {}) => {
    return api.get('/promoter/list', params)
  },

  // 获取推广员详情
  getPromoter: (id) => {
    return api.get(`/promoter/detail/${id}`)
  },

  // 生成推广链接
  generateLink: (promoterId, campaignId) => {
    return api.post('/promoter/generate-link', { promoterId, campaignId })
  },

  // 获取推广员奖励记录
  getRewards: (promoterId, params = {}) => {
    return api.get(`/promoter/rewards/${promoterId}`, params)
  }
}

// 数据分析API
export const analyticsApi = {
  // 获取核心指标
  getMetrics: (period = 'month') => {
    return api.get('/analytics/metrics', { period })
  },

  // 获取趋势数据
  getTrends: (period = 'month') => {
    return api.get('/analytics/trends', { period })
  },

  // 获取活动排行
  getCampaignRanking: (period = 'month') => {
    return api.get('/analytics/campaign-ranking', { period })
  },

  // 获取推广员排行
  getPromoterRanking: (period = 'month') => {
    return api.get('/analytics/promoter-ranking', { period })
  }
}

// 海报生成相关API
export const posterApi = {
  // 获取海报模板列表
  getPosterTemplates: () => {
    return api.get('/poster/templates')
  },

  // 生成海报（活动海报或分销商海报）
  generatePoster: (campaignId, distributorId) => {
    return api.post(`/campaigns/${campaignId}/poster`, {
      distributorId: distributorId
    })
  },

  // 获取海报文件
  getPosterFile: (filename) => {
    return api.get(`/poster/${filename}`)
  },

  // 获取海报生成记录
  getPosterRecords: () => {
    return api.get('/poster/records')
  }
}

// 素材管理API
export const materialApi = {
  // 上传素材
  uploadMaterial: (file, type = 'image') => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('type', type)
    return api.post('/material/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  // 获取素材列表
  getMaterials: (params = {}) => {
    return api.get('/material/list', params)
  },

  // 删除素材
  deleteMaterial: (id) => {
    return api.delete(`/material/delete/${id}`)
  }
}

// 品牌设置API
export const settingsApi = {
  getBrandInfo: (brandId) => {
    return api.get(`/brands/${brandId}`)
  },

  updateBrandInfo: (brandId, data) => {
    return api.put(`/brands/${brandId}`, data)
  },

  getRewardSettings: () => {
    const raw = localStorage.getItem('dmh_brand_reward_settings')
    if (!raw) return Promise.resolve({})

    try {
      return Promise.resolve(JSON.parse(raw))
    } catch (error) {
      console.error('奖励设置解析失败:', error)
      return Promise.resolve({})
    }
  },

  updateRewardSettings: (data) => {
    localStorage.setItem('dmh_brand_reward_settings', JSON.stringify(data || {}))
    return Promise.resolve({ success: true })
  },

  getNotificationSettings: () => {
    const raw = localStorage.getItem('dmh_brand_notification_settings')
    if (!raw) return Promise.resolve({})

    try {
      return Promise.resolve(JSON.parse(raw))
    } catch (error) {
      console.error('通知设置解析失败:', error)
      return Promise.resolve({})
    }
  },

  updateNotificationSettings: (data) => {
    localStorage.setItem('dmh_brand_notification_settings', JSON.stringify(data || {}))
    return Promise.resolve({ success: true })
  },

  getSyncSettings: () => {
    const raw = localStorage.getItem('dmh_brand_sync_settings')
    if (!raw) return Promise.resolve({})

    try {
      return Promise.resolve(JSON.parse(raw))
    } catch (error) {
      console.error('同步设置解析失败:', error)
      return Promise.resolve({})
    }
  },

  updateSyncSettings: (data) => {
    localStorage.setItem('dmh_brand_sync_settings', JSON.stringify(data || {}))
    return Promise.resolve({ success: true })
  },

  changePassword: (data) => {
    return api.post('/users/change-password', data)
  },

  getSyncHealth: () => {
    return api.get('/sync/health')
  }
}

// AI 文案生成 API
export const aiApi = {
  // 生成营销文案
  generateCopywriting: (params) => {
    return api.post('/ai/generate-copywriting', params)
  }
}
