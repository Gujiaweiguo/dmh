const API_BASE_URL = '/api/v1';

export interface RoleItem {
  id: number;
  name: string;
  code: string;
  description: string;
  permissions: string[];
  createdAt: string;
}

export interface PermissionItem {
  id: number;
  name: string;
  code: string;
  resource: string;
  action: string;
  description: string;
}

export interface AuditLogItem {
  id: number;
  userId?: number;
  username?: string;
  action: string;
  details?: string;
  clientIp?: string;
  userAgent?: string;
  createdAt: string;
}

class RoleApiService {
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

  async getRoles(): Promise<RoleItem[]> {
    return this.request<RoleItem[]>('/roles');
  }

  async getPermissions(): Promise<PermissionItem[]> {
    return this.request<PermissionItem[]>('/permissions');
  }

  async configRolePermissions(roleId: number, permissionIds: number[]): Promise<void> {
    await this.request<{ message: string }>('/roles/permissions', {
      method: 'POST',
      body: JSON.stringify({ roleId, permissionIds }),
    });
  }

  async getAuditLogs(): Promise<AuditLogItem[]> {
    const data = await this.request<{ logs?: AuditLogItem[]; total?: number }>('/security/audit-logs');
    return Array.isArray(data.logs) ? data.logs : [];
  }
}

export const roleApi = new RoleApiService();
