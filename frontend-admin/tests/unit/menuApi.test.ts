import { beforeEach, describe, expect, it, vi } from 'vitest';
import { menuApi } from '../../services/menuApi';

describe('menuApi service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    Object.defineProperty(globalThis, 'fetch', {
      value: vi.fn(),
      writable: true,
    });
  });

  it('getUserMenus uses platform query and auth header', async () => {
    localStorage.setItem('dmh_token', 'token-menu');
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ userId: 1, platform: 'admin', menus: [{ id: 1, name: 'Dashboard', code: 'dash', path: '/dashboard', sort: 1, type: 'menu', platform: 'admin', status: 'active' }] }),
    } as unknown as Response);

    const menus = await menuApi.getUserMenus('admin');
    const [url, options] = vi.mocked(fetch).mock.calls[0];
    const headers = (options as RequestInit).headers as Record<string, string>;

    expect(url).toBe('/api/v1/users/menus?platform=admin');
    expect(headers.Authorization).toBe('Bearer token-menu');
    expect(menus).toHaveLength(1);
  });

  it('returns empty array when menus is not an array', async () => {
    localStorage.setItem('dmh_token', 'token-menu');
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ userId: 1, platform: 'admin', menus: null }),
    } as unknown as Response);

    const menus = await menuApi.getUserMenus();
    expect(menus).toEqual([]);
  });

  it('throws login error when token missing', async () => {
    await expect(menuApi.getUserMenus()).rejects.toThrow('请先登录');
  });

  it('clears token and throws on 401 response', async () => {
    localStorage.setItem('dmh_token', 'expired-menu-token');
    const removeSpy = vi.spyOn(Storage.prototype, 'removeItem');
    vi.mocked(fetch).mockResolvedValue({ ok: false, status: 401, json: async () => ({ message: 'expired' }) } as unknown as Response);

    await expect(menuApi.getUserMenus()).rejects.toThrow('登录已过期，请重新登录');
    expect(removeSpy).toHaveBeenCalledWith('dmh_token');
  });

  it('throws backend message on non-OK response', async () => {
    localStorage.setItem('dmh_token', 'token-menu');
    vi.mocked(fetch).mockResolvedValue({ ok: false, status: 500, json: async () => ({ message: 'menu failed' }) } as unknown as Response);

    await expect(menuApi.getUserMenus()).rejects.toThrow('menu failed');
  });

  it('falls back to generic message when json parse fails', async () => {
    localStorage.setItem('dmh_token', 'token-menu');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error('invalid json');
      },
    } as unknown as Response);

    await expect(menuApi.getUserMenus()).rejects.toThrow('请求失败');
  });
});
