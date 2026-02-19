const API_BASE_URL = '/api/v1';

export interface UserProfile {
  id: number;
  username: string;
  phone: string;
  email: string;
  realName: string;
  avatar: string;
  status: string;
  roles: string[];
  brandIds?: number[];
  createdAt: string;
}

interface ApiMessageResp {
  message?: string;
}

class ProfileApiService {
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

  async getUserInfo(): Promise<UserProfile> {
    return this.request<UserProfile>('/auth/userinfo', {
      method: 'GET',
    });
  }

  async updateProfile(realName: string): Promise<UserProfile> {
    return this.request<UserProfile>('/users/profile', {
      method: 'PUT',
      body: JSON.stringify({ realName }),
    });
  }

  async sendPhoneCode(phone: string): Promise<ApiMessageResp> {
    return this.request<ApiMessageResp>('/auth/send-phone-code', {
      method: 'POST',
      body: JSON.stringify({
        target: phone,
        type: 'phone',
      }),
    });
  }

  async sendEmailCode(email: string): Promise<ApiMessageResp> {
    return this.request<ApiMessageResp>('/auth/send-email-code', {
      method: 'POST',
      body: JSON.stringify({
        target: email,
        type: 'email',
      }),
    });
  }

  async bindPhone(phone: string, code: string): Promise<UserProfile> {
    return this.request<UserProfile>('/users/bind-phone', {
      method: 'POST',
      body: JSON.stringify({ phone, code }),
    });
  }

  async bindEmail(email: string, code: string): Promise<UserProfile> {
    return this.request<UserProfile>('/users/bind-email', {
      method: 'POST',
      body: JSON.stringify({ email, code }),
    });
  }

  async changePassword(oldPassword: string, newPassword: string): Promise<ApiMessageResp> {
    return this.request<ApiMessageResp>('/users/change-password', {
      method: 'POST',
      body: JSON.stringify({ oldPassword, newPassword }),
    });
  }
}

export const profileApi = new ProfileApiService();
