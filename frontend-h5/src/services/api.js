// API 基础配置
const API_BASE_URL = '/api/v1'

// 获取认证头
const getAuthHeaders = () => {
  const token = localStorage.getItem('dmh_token')
  return token ? { 'Authorization': `Bearer ${token}` } : {}
}

// 通用请求方法
const request = async (url, options = {}) => {
  const config = {
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
      ...options.headers
    },
    ...options
  }

  try {
    const fullUrl = `${API_BASE_URL}${url}`
    console.log('API请求:', fullUrl, config)

    const response = await fetch(fullUrl, config)

    console.log('API响应:', response.status, response.statusText)

    if (!response.ok) {
      // 尝试获取错误详情
      let errorMessage = `HTTP error! status: ${response.status}`
      try {
        const errorData = await response.text()
        if (errorData) {
          errorMessage += ` - ${errorData}`
        }
      } catch (e) {
        // 忽略解析错误
      }
      throw new Error(errorMessage)
    }

    const data = await response.json()
    console.log('API数据:', data)
    return data
  } catch (error) {
    console.error('API request failed:', error)
    throw error
  }
}

export const api = {
  // GET 请求
  get: (url, params = {}) => {
    const queryString = new URLSearchParams(params).toString()
    const fullUrl = queryString ? `${url}?${queryString}` : url
    return request(fullUrl, { method: 'GET' })
  },

  // POST 请求
  post: (url, data = {}) => {
    return request(url, {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // PUT 请求
  put: (url, data = {}) => {
    return request(url, {
      method: 'PUT',
      body: JSON.stringify(data)
    })
  },

  // DELETE 请求
  delete: (url) => {
    return request(url, { method: 'DELETE' })
  }
}

// 导出默认对象（兼容性）
export default api

// 订单管理API
export const orderApi = {
  // 获取订单列表
  getOrders: (params = {}) => {
    return api.get('/orders/list', params)
  },

  // 获取订单详情
  getOrder: (id) => {
    return api.get(`/orders/${id}`)
  },

  // 更新订单状态
  updateOrderStatus: (id, status) => {
    return api.put(`/orders/${id}`, { status })
  },

  // 扫码获取订单信息（后端路由: GET /orders/scan）
  scanOrderCode: (code) => {
    return api.get('/orders/scan', { code })
  },

  // 核销订单（后端路由: POST /orders/verify）
  verifyOrder: (code, notes) => {
    return api.post('/orders/verify', { code, notes })
  },

  // 取消核销（后端路由: POST /orders/unverify）
  unverifyOrder: (code, reason) => {
    return api.post('/orders/unverify', { code, reason })
  },

  // 获取核销记录（后端路由: GET /orders/verification-records）
  getVerificationRecords: () => {
    return api.get('/orders/verification-records')
  }
}

// 品牌管理员专用API
export const brandApi = {
  // 扫码核销（品牌管理员）
  scanOrderCode: (code) => {
    return api.get('/orders/scan', { code })
  },

  // 验证订单码（品牌管理员）- 后端路由: POST /orders/verify
  verifyOrderCode: (code, notes) => {
    return api.post('/orders/verify', { code, notes })
  },

  // 取消订单核销（品牌管理员）- 后端路由: POST /orders/unverify
  unverifyOrderCode: (code, reason) => {
    return api.post('/orders/unverify', { code, reason })
  },

  // 核销记录查询
  getVerificationRecords: () => {
    return api.get('/orders/verification-records')
  }
}

// 反馈系统 API
export const feedbackApi = {
  createFeedback: (data) => {
    return api.post('/feedback', data)
  },

  listFeedback: (params = {}) => {
    return api.get('/feedback/list', params)
  },

  getFeedbackDetail: (id) => {
    return api.get('/feedback/detail', { id })
  },

  submitSatisfactionSurvey: (data) => {
    return api.post('/feedback/satisfaction-survey', data)
  },

  listFaq: (params = {}) => {
    return api.get('/feedback/faq', params)
  },

  markFaqHelpful: (id, type = 'helpful') => {
    return api.post('/feedback/faq/helpful', { id, type })
  },

  recordFeatureUsage: (data) => {
    return api.post('/feedback/feature-usage', data)
  },

  getFeedbackStatistics: (params = {}) => {
    return api.get('/feedback/statistics', params)
  }
}
