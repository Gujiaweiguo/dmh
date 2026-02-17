const API_BASE_URL = '/api/v1';

export interface AdminUser {
  id: number;
  username: string;
  phone: string;
  email?: string;
  realName?: string;
  status: 'active' | 'disabled' | 'locked' | string;
  roles: string[];
  brandIds?: number[];
  createdAt: string;
}

interface AdminUserListResponse {
  total: number;
  users: AdminUser[];
}

class UserApiService {
  private getToken() {
    return localStorage.getItem('dmh_token');
  }

  private async request<T>(url: string, options: RequestInit = {}): Promise<T> {
    const token = this.getToken();
    if (!token) {
      throw new Error('请先登录');
    }

    const response = await fetch(`${API_BASE_URL}${url}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
        ...options.headers,
      },
    });

    if (response.status === 401) {
      localStorage.removeItem('dmh_token');
      throw new Error('登录已过期，请重新登录');
    }

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: '请求失败' }));
      throw new Error(error.message || '请求失败');
    }

    return response.json() as Promise<T>;
  }

  async getUsers(params: {
    page?: number;
    pageSize?: number;
    role?: string;
    status?: string;
    keyword?: string;
  } = {}): Promise<AdminUserListResponse> {
    const query = new URLSearchParams();
    if (params.page) query.append('page', String(params.page));
    if (params.pageSize) query.append('pageSize', String(params.pageSize));
    if (params.role && params.role !== 'all') query.append('role', params.role);
    if (params.status && params.status !== 'all') query.append('status', params.status);
    if (params.keyword) query.append('keyword', params.keyword);

    const queryString = query.toString();
    return this.request<AdminUserListResponse>(`/admin/users${queryString ? `?${queryString}` : ''}`);
  }

  async updateUser(id: number, payload: {
    realName?: string;
    email?: string;
    role?: string;
    status?: string;
    brandIds?: number[];
  }): Promise<AdminUser> {
    return this.request<AdminUser>(`/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    });
  }

  async resetUserPassword(id: number, newPassword: string): Promise<void> {
    await this.request<{ message: string }>(`/admin/users/${id}/reset-password`, {
      method: 'POST',
      body: JSON.stringify({ newPassword }),
    });
  }

  async deleteUser(id: number): Promise<void> {
    await this.request<{ message: string }>(`/admin/users/${id}`, {
      method: 'DELETE',
    });
  }
}

export const userApi = new UserApiService();
