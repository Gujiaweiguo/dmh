import { beforeEach, describe, expect, it, vi } from 'vitest';
import { roleApi } from '../../services/roleApi';

describe('roleApi service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    Object.defineProperty(globalThis, 'fetch', {
      value: vi.fn(),
      writable: true,
    });
  });

  it('getRoles fetches role list', async () => {
    localStorage.setItem('dmh_token', 'token-role');
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => [{ id: 1, name: 'admin', code: 'platform_admin', description: '', permissions: [], createdAt: '' }],
    } as unknown as Response);

    const roles = await roleApi.getRoles();
    expect(roles).toHaveLength(1);
    expect(roles[0].code).toBe('platform_admin');
  });

  it('configRolePermissions sends POST payload', async () => {
    localStorage.setItem('dmh_token', 'token-role');
    vi.mocked(fetch).mockResolvedValue({ ok: true, status: 200, json: async () => ({ message: 'ok' }) } as unknown as Response);

    await roleApi.configRolePermissions(7, [1, 2, 3]);

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    expect(url).toBe('/api/v1/roles/permissions');
    expect((options as RequestInit).method).toBe('POST');
    expect((options as RequestInit).body).toBe(JSON.stringify({ roleId: 7, permissionIds: [1, 2, 3] }));
  });

  it('getAuditLogs returns empty array when logs missing', async () => {
    localStorage.setItem('dmh_token', 'token-role');
    vi.mocked(fetch).mockResolvedValue({ ok: true, status: 200, json: async () => ({ total: 9 }) } as unknown as Response);

    const logs = await roleApi.getAuditLogs();
    expect(logs).toEqual([]);
  });

  it('throws login error when token missing', async () => {
    await expect(roleApi.getPermissions()).rejects.toThrow('请先登录');
  });

  it('clears token and throws on 401 response', async () => {
    localStorage.setItem('dmh_token', 'expired-role-token');
    const removeSpy = vi.spyOn(Storage.prototype, 'removeItem');
    vi.mocked(fetch).mockResolvedValue({ ok: false, status: 401, json: async () => ({ message: 'expired' }) } as unknown as Response);

    await expect(roleApi.getPermissions()).rejects.toThrow('登录已过期，请重新登录');
    expect(removeSpy).toHaveBeenCalledWith('dmh_token');
  });

  it('falls back to generic message when json parse fails', async () => {
    localStorage.setItem('dmh_token', 'token-role');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error('invalid json');
      },
    } as unknown as Response);

    await expect(roleApi.getPermissions()).rejects.toThrow('请求失败');
  });
});
