import { beforeEach, describe, expect, it, vi } from 'vitest';
import { authApi } from '../../services/authApi';

describe('authApi', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    authApi.logout();
    global.fetch = vi.fn();
  });

  it('login with platform_admin stores auth data', async () => {
    const mockResp = {
      token: 'token-abc',
      username: 'admin',
      roles: ['platform_admin'],
    };
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => mockResp,
    } as Response);

    const result = await authApi.login({ username: 'admin', password: '123456' });

    expect(result).toEqual(mockResp);
    expect(localStorage.getItem('dmh_token')).toBe('token-abc');
    expect(localStorage.getItem('dmh_user_role')).toBe('platform_admin');
    expect(localStorage.getItem('dmh_username')).toBe('admin');
    expect(authApi.isLoggedIn()).toBe(true);
  });

  it('login rejects non-platform-admin user', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => ({
        token: 'token-user',
        username: 'user001',
        roles: ['participant'],
      }),
    } as Response);

    await expect(authApi.login({ username: 'user001', password: '123456' })).rejects.toThrow(
      '管理后台仅限平台管理员访问',
    );
    expect(localStorage.getItem('dmh_token')).toBeNull();
    expect(authApi.isLoggedIn()).toBe(false);
  });

  it('login throws when response is not ok', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 400,
      text: async () => '用户名或密码错误',
    } as Response);

    await expect(authApi.login({ username: 'admin', password: 'wrong' })).rejects.toThrow(
      'HTTP error! status: 400 - 用户名或密码错误',
    );
  });

  it('login state manager subscribe and unsubscribe works', () => {
    const listener = vi.fn();
    const manager = authApi.getLoginStateManager();
    const unsubscribe = manager.subscribe(listener);

    manager.value = true;
    manager.value = true;
    manager.value = false;

    expect(listener).toHaveBeenCalledTimes(2);
    expect(listener).toHaveBeenNthCalledWith(1, true);
    expect(listener).toHaveBeenNthCalledWith(2, false);

    unsubscribe();
    manager.value = true;
    expect(listener).toHaveBeenCalledTimes(2);
  });

  it('isPlatformAdmin returns true when user is platform_admin', () => {
    localStorage.setItem('dmh_user_role', 'platform_admin');
    expect(authApi.isPlatformAdmin()).toBe(true);
  });

  it('isPlatformAdmin returns false when user is not platform_admin', () => {
    localStorage.setItem('dmh_user_role', 'brand_admin');
    expect(authApi.isPlatformAdmin()).toBe(false);
  });

  it('getUsername returns stored username', () => {
    localStorage.setItem('dmh_username', 'testuser');
    expect(authApi.getUsername()).toBe('testuser');
  });

  it('getUsername returns null when no username stored', () => {
    expect(authApi.getUsername()).toBeNull();
  });
});
