import { beforeEach, describe, expect, it, vi } from 'vitest';
import { userApi } from '../../services/userApi';

describe('userApi service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    Object.defineProperty(globalThis, 'fetch', {
      value: vi.fn(),
      writable: true,
    });
  });

  it('getUsers builds query and skips all filters', async () => {
    localStorage.setItem('dmh_token', 'token-user');
    vi.mocked(fetch).mockResolvedValue({ ok: true, status: 200, json: async () => ({ total: 0, users: [] }) } as Response);

    await userApi.getUsers({ page: 2, pageSize: 10, role: 'all', status: 'active', keyword: 'alice' });

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    const urlStr = String(url);
    const headers = (options as RequestInit).headers as Record<string, string>;
    expect(urlStr).toContain('/api/v1/admin/users?');
    expect(urlStr).toContain('page=2');
    expect(urlStr).toContain('pageSize=10');
    expect(urlStr).toContain('status=active');
    expect(urlStr).toContain('keyword=alice');
    expect(urlStr).not.toContain('role=all');
    expect(headers.Authorization).toBe('Bearer token-user');
  });

  it('updateUser sends PUT payload', async () => {
    localStorage.setItem('dmh_token', 'token-user');
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ id: 8, username: 'u', phone: '1', status: 'active', roles: [], createdAt: '2026-01-01' }),
    } as Response);

    await userApi.updateUser(8, { realName: 'Alice', brandIds: [1, 2] });

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    expect(url).toBe('/api/v1/admin/users/8');
    expect((options as RequestInit).method).toBe('PUT');
    expect((options as RequestInit).body).toBe(JSON.stringify({ realName: 'Alice', brandIds: [1, 2] }));
  });

  it('reset and delete user send expected methods', async () => {
    localStorage.setItem('dmh_token', 'token-user');
    vi.mocked(fetch).mockResolvedValue({ ok: true, status: 200, json: async () => ({ message: 'ok' }) } as Response);

    await userApi.resetUserPassword(3, 'NewPassword123!');
    await userApi.deleteUser(3);

    expect(vi.mocked(fetch).mock.calls[0][0]).toBe('/api/v1/admin/users/3/reset-password');
    expect((vi.mocked(fetch).mock.calls[0][1] as RequestInit).method).toBe('POST');
    expect(vi.mocked(fetch).mock.calls[1][0]).toBe('/api/v1/admin/users/3');
    expect((vi.mocked(fetch).mock.calls[1][1] as RequestInit).method).toBe('DELETE');
  });

  it('throws login error when token missing', async () => {
    await expect(userApi.getUsers()).rejects.toThrow('请先登录');
  });

  it('clears token and throws on 401 response', async () => {
    localStorage.setItem('dmh_token', 'expired-user-token');
    const removeSpy = vi.spyOn(Storage.prototype, 'removeItem');
    vi.mocked(fetch).mockResolvedValue({ ok: false, status: 401, json: async () => ({ message: 'expired' }) } as Response);

    await expect(userApi.getUsers()).rejects.toThrow('登录已过期，请重新登录');
    expect(removeSpy).toHaveBeenCalledWith('dmh_token');
  });

  it('falls back to generic message when json parse fails', async () => {
    localStorage.setItem('dmh_token', 'token-user');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error('invalid json');
      },
    } as Response);

    await expect(userApi.getUsers()).rejects.toThrow('请求失败');
  });
});
