// 会员管理 API 服务

const API_BASE_URL = '/api/v1';

// 获取 Token
const getToken = () => localStorage.getItem('dmh_token');

// 通用请求方法
async function request(url: string, options: RequestInit = {}) {
  const token = getToken();
  
  // 如果没有 token，提示用户重新登录
  if (!token) {
    // 清除可能存在的无效状态
    localStorage.removeItem('dmh_token');
    throw new Error('请先登录');
  }
  
  const headers = {
    'Content-Type': 'application/json',
    Authorization: `Bearer ${token}`,
    ...options.headers,
  };

  const response = await fetch(`${API_BASE_URL}${url}`, {
    ...options,
    headers,
  });

  // 处理 401 未授权错误
  if (response.status === 401) {
    // Token 无效或过期，清除并提示重新登录
    localStorage.removeItem('dmh_token');
    throw new Error('登录已过期，请重新登录');
  }

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: '请求失败' }));
    throw new Error(error.message || '请求失败');
  }

  return response.json();
}

// 会员相关接口
export const memberApi = {
  // 获取会员列表
  getMembers: (params: {
    page?: number;
    pageSize?: number;
    brandId?: number;
    keyword?: string;
    source?: string;
    status?: string;
    tagIds?: number[];
    startDate?: string;
    endDate?: string;
  }) => {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        if (Array.isArray(value)) {
          value.forEach(v => queryParams.append(key, v.toString()));
        } else {
          queryParams.append(key, value.toString());
        }
      }
    });
    return request(`/members?${queryParams.toString()}`);
  },

  // 获取会员详情
  getMember: (id: number) => {
    return request(`/members/${id}`);
  },

  // 创建会员标签
  createTag: (data: {
    name: string;
    category?: string;
    color?: string;
    description?: string;
  }) => {
    return request('/members/tags', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  // 为会员添加标签
  addMemberTags: (memberId: number, tagIds: number[]) => {
    return request(`/members/${memberId}/tags`, {
      method: 'POST',
      body: JSON.stringify({ memberId, tagIds }),
    });
  },

  // 会员合并预览
  previewMerge: (data: {
    sourceMemberId: number;
    targetMemberId: number;
    reason?: string;
  }) => {
    return request('/members/merge/preview', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  // 执行会员合并
  mergeMember: (data: {
    sourceMemberId: number;
    targetMemberId: number;
    reason?: string;
  }) => {
    return request('/members/merge', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  // 创建导出申请
  createExportRequest: (data: {
    brandId: number;
    reason: string;
    filters?: string;
  }) => {
    return request('/members/export-requests', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  // 获取导出申请列表
  getExportRequests: (params: {
    page?: number;
    pageSize?: number;
    brandId?: number;
    status?: string;
  }) => {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        queryParams.append(key, value.toString());
      }
    });
    return request(`/members/export-requests?${queryParams.toString()}`);
  },

  // 审批导出申请
  approveExportRequest: (id: number, data: {
    approve: boolean;
    reason?: string;
  }) => {
    return request(`/members/export-requests/${id}/approve`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },
};
