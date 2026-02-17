const API_BASE_URL = '/api/v1';

export interface MenuItem {
  id: number;
  name: string;
  code: string;
  path: string;
  icon?: string;
  parentId?: number;
  sort: number;
  type: 'menu' | 'button';
  platform: 'admin' | 'h5';
  status: 'active' | 'disabled';
  permission?: string;
  role?: string;
  children?: MenuItem[];
}

interface UserMenuResponse {
  userId: number;
  platform: string;
  menus: MenuItem[];
}

class MenuApiService {
  private getToken() {
    return localStorage.getItem('dmh_token');
  }

  async getUserMenus(platform: 'admin' | 'h5' = 'admin'): Promise<MenuItem[]> {
    const token = this.getToken();
    if (!token) {
      throw new Error('请先登录');
    }

    const response = await fetch(`${API_BASE_URL}/users/menus?platform=${platform}`, {
      headers: {
        Authorization: `Bearer ${token}`,
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

    const data = (await response.json()) as UserMenuResponse;
    return Array.isArray(data.menus) ? data.menus : [];
  }
}

export const menuApi = new MenuApiService();
