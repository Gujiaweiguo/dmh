// 分销商 API 服务
const API_BASE = '/api/v1';

// 获取 token
const getToken = () => localStorage.getItem('dmh_token');

// 通用请求封装
async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`,
      ...options?.headers,
    },
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`HTTP error! status: ${response.status} - ${errorText}`);
  }

  return response.json();
}

export const distributorApi = {
  // 获取分销商申请列表
  getApplications: async (brandId: number, params?: any) => {
    const query = params ? '?' + new URLSearchParams(params).toString() : '';
    return request<any>(`${API_BASE}/brands/${brandId}/distributor/applications${query}`);
  },

  // 获取申请详情
  getApplication: async (brandId: number, applicationId: number) => {
    return request<any>(`${API_BASE}/brands/${brandId}/distributor/applications/${applicationId}`);
  },

  // 审批申请
  approveApplication: async (brandId: number, applicationId: number, data: any) => {
    return request<any>(`${API_BASE}/brands/${brandId}/distributor/approve/${applicationId}`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  // 获取分销商列表
  getDistributors: async (brandId: number, params?: any) => {
    const query = params ? '?' + new URLSearchParams(params).toString() : '';
    return request<any>(`${API_BASE}/brands/${brandId}/distributors${query}`);
  },

  // 获取分销商详情
  getDistributor: async (brandId: number, distributorId: number) => {
    return request<any>(`${API_BASE}/brands/${brandId}/distributors/${distributorId}`);
  },

  // 更新分销商级别
  updateDistributorLevel: async (distributorId: number, level: number) => {
    return request<any>(`${API_BASE}/brands/distributors/${distributorId}/level`, {
      method: 'PUT',
      body: JSON.stringify({ level }),
    });
  },

  // 更新分销商状态
  updateDistributorStatus: async (distributorId: number, status: string, reason?: string) => {
    return request<any>(`${API_BASE}/brands/distributors/${distributorId}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status, reason }),
    });
  },

  // 获取级别奖励配置
  getLevelRewards: async (brandId: number) => {
    return request<any>(`${API_BASE}/brands/${brandId}/distributor/level-rewards`);
  },

  // 设置级别奖励配置
  setLevelRewards: async (brandId: number, rewards: any[]) => {
    return request<any>(`${API_BASE}/brands/${brandId}/distributor/level-rewards`, {
      method: 'PUT',
      body: JSON.stringify({ rewards }),
    });
  },

  // 获取提现记录列表（平台管理员）
  getWithdrawals: async (brandId: number, status?: string, page?: number, pageSize?: number) => {
    const params: any = {};
    if (brandId > 0) params.brandId = brandId;
    if (status) params.status = status;
    if (page) params.page = page;
    if (pageSize) params.pageSize = pageSize;
    const query = Object.keys(params).length > 0 ? '?' + new URLSearchParams(params).toString() : '';
    return request<any>(`${API_BASE}/platform/withdrawals${query}`);
  },

  // 批准提现
  approveWithdrawal: async (withdrawalId: number, notes?: string) => {
    return request<any>(`${API_BASE}/platform/withdrawals/${withdrawalId}/approve`, {
      method: 'PUT',
      body: JSON.stringify({ notes }),
    });
  },

  // 拒绝提现
  rejectWithdrawal: async (withdrawalId: number, reason: string) => {
    return request<any>(`${API_BASE}/platform/withdrawals/${withdrawalId}/reject`, {
      method: 'PUT',
      body: JSON.stringify({ reason }),
    });
  },

  // 获取全局分销商列表（平台管理员）
  getGlobalDistributors: async (brandId?: number, status?: string, page?: number, pageSize?: number) => {
    const params: any = {};
    if (brandId) params.brandId = brandId;
    if (status) params.status = status;
    if (page) params.page = page;
    if (pageSize) params.pageSize = pageSize;
    const query = Object.keys(params).length > 0 ? '?' + new URLSearchParams(params).toString() : '';
    return request<any>(`${API_BASE}/platform/distributors${query}`);
  },

  // 获取全局奖励列表（平台管理员）
  getGlobalRewards: async (brandId?: number, page?: number, pageSize?: number) => {
    const params: any = {};
    if (brandId) params.brandId = brandId;
    if (page) params.page = page;
    if (pageSize) params.pageSize = pageSize;
    const query = Object.keys(params).length > 0 ? '?' + new URLSearchParams(params).toString() : '';
    return request<any>(`${API_BASE}/platform/rewards${query}`);
  },
 };
