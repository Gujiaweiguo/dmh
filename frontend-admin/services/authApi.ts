import { LoginRequest, LoginResponse } from '../types';

const API_BASE_URL = '/api/v1';

// 用户信息接口
interface UserInfo {
  id: number;
  username: string;
  role: string;
  realName?: string;
}

// 简单的响应式状态管理
class LoginStateManager {
  private _isLoggedIn: boolean;
  private listeners: Array<(state: boolean) => void> = [];

  constructor() {
    this._isLoggedIn = !!localStorage.getItem('dmh_token');
  }

  get value() {
    return this._isLoggedIn;
  }

  set value(newValue: boolean) {
    if (this._isLoggedIn !== newValue) {
      this._isLoggedIn = newValue;
      this.listeners.forEach(listener => {
        try {
          listener(newValue);
        } catch (e) {
          console.error('Error in login state listener:', e);
        }
      });
    }
  }

  subscribe(listener: (state: boolean) => void) {
    this.listeners.push(listener);
    return () => {
      const index = this.listeners.indexOf(listener);
      if (index > -1) this.listeners.splice(index, 1);
    };
  }
}

const loginStateManager = new LoginStateManager();

class AuthApiService {
  // 获取登录状态管理器
  getLoginStateManager() {
    return loginStateManager;
  }

  // 登录 - 只允许平台管理员登录
  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP error! status: ${response.status} - ${errorText}`);
    }

    const result = await response.json() as { roles?: string[]; role?: string; token?: string; username?: string };

    const roles: string[] = Array.isArray(result.roles)
      ? result.roles
      : (typeof result.role === 'string' ? [result.role] : []);
    // 检查用户角色 - 只允许平台管理员登录管理后台
    if (!roles.includes('platform_admin')) {
      throw new Error('管理后台仅限平台管理员访问，请使用 H5 端登录');
    }

    // 保存token和用户信息到localStorage
    if (result.token) {
      localStorage.setItem('dmh_token', result.token);
      localStorage.setItem('dmh_user_role', 'platform_admin');
      if (result.username) {
        localStorage.setItem('dmh_username', result.username);
      }
      loginStateManager.value = true; // 更新状态
    }

    return result;
  }

  // 注册
  async register(data: {
    username: string;
    password: string;
    phone: string;
    email?: string;
    realName?: string;
  }): Promise<LoginResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Registration failed: ${response.status} ${errorText}`);
    }

    const result = await response.json() as { roles?: string[]; role?: string; token?: string };

    // 保存token到localStorage
    if (result.token) {
      localStorage.setItem('dmh_token', result.token);
      const roles: string[] = Array.isArray(result.roles)
        ? result.roles
        : (typeof result.role === 'string' ? [result.role] : []);
      if (roles.length > 0) {
        localStorage.setItem('dmh_user_role', roles[0]);
      }
      loginStateManager.value = true; // 更新状态
    }

    return result;
  }

  // 登出
  logout() {
    localStorage.removeItem('dmh_token');
    localStorage.removeItem('dmh_user_role');
    localStorage.removeItem('dmh_username');
    loginStateManager.value = false; // 更新状态
  }

  // 检查是否已登录
  isLoggedIn(): boolean {
    return !!localStorage.getItem('dmh_token');
  }

  // 获取token
  getToken(): string | null {
    return localStorage.getItem('dmh_token');
  }

  // 获取用户角色
  getUserRole(): string | null {
    return localStorage.getItem('dmh_user_role');
  }

  // 检查是否是平台管理员
  isPlatformAdmin(): boolean {
    return this.getUserRole() === 'platform_admin';
  }

  // 获取用户名
  getUsername(): string | null {
    return localStorage.getItem('dmh_username');
  }
}

export const authApi = new AuthApiService();
