import { beforeEach, describe, expect, it, vi } from 'vitest';
import { securityApi } from '../../services/securityApi';

describe('securityApi service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    Object.defineProperty(globalThis, 'fetch', {
      value: vi.fn(),
      writable: true,
    });
  });

  it('getPasswordPolicy sends GET request', async () => {
    localStorage.setItem('dmh_token', 'token-security');
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ id: 1, minLength: 8 }),
    } as Response);

    await securityApi.getPasswordPolicy();

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    expect(url).toBe('/api/v1/security/password-policy');
    expect((options as RequestInit).method).toBe('GET');
  });

  it('updatePasswordPolicy sends PUT payload', async () => {
    localStorage.setItem('dmh_token', 'token-security');
    const payload = {
      id: 1,
      minLength: 10,
      requireUppercase: true,
      requireLowercase: true,
      requireNumbers: true,
      requireSpecialChars: false,
      maxAge: 90,
      historyCount: 5,
      maxLoginAttempts: 5,
      lockoutDuration: 30,
      sessionTimeout: 120,
      maxConcurrentSessions: 3,
    };
    vi.mocked(fetch).mockResolvedValue({ ok: true, status: 200, json: async () => payload } as Response);

    await securityApi.updatePasswordPolicy(payload);

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    expect(url).toBe('/api/v1/security/password-policy');
    expect((options as RequestInit).method).toBe('PUT');
    expect((options as RequestInit).body).toBe(JSON.stringify(payload));
  });

  it('getUserSessions and getSecurityEvents append pagination query', async () => {
    localStorage.setItem('dmh_token', 'token-security');
    vi.mocked(fetch)
      .mockResolvedValueOnce({ ok: true, status: 200, json: async () => ({ total: 0, sessions: [] }) } as Response)
      .mockResolvedValueOnce({ ok: true, status: 200, json: async () => ({ total: 0, events: [] }) } as Response);

    await securityApi.getUserSessions(2, 30);
    await securityApi.getSecurityEvents(3, 40);

    expect(String(vi.mocked(fetch).mock.calls[0][0])).toContain('/api/v1/security/sessions?page=2&pageSize=30');
    expect(String(vi.mocked(fetch).mock.calls[1][0])).toContain('/api/v1/security/events?page=3&pageSize=40');
  });

  it('forceLogoutUser and handleSecurityEvent send default payload values', async () => {
    localStorage.setItem('dmh_token', 'token-security');
    vi.mocked(fetch)
      .mockResolvedValueOnce({ ok: true, status: 200, json: async () => ({ message: 'ok' }) } as Response)
      .mockResolvedValueOnce({ ok: true, status: 200, json: async () => ({ message: 'ok' }) } as Response);

    await securityApi.forceLogoutUser(9);
    await securityApi.handleSecurityEvent(11);

    expect(vi.mocked(fetch).mock.calls[0][0]).toBe('/api/v1/security/force-logout/9');
    expect((vi.mocked(fetch).mock.calls[0][1] as RequestInit).body).toBe(JSON.stringify({ reason: '管理员操作' }));
    expect(vi.mocked(fetch).mock.calls[1][0]).toBe('/api/v1/security/events/11/handle');
    expect((vi.mocked(fetch).mock.calls[1][1] as RequestInit).body).toBe(JSON.stringify({ note: '' }));
  });

  it('throws login error when token missing', async () => {
    await expect(securityApi.revokeSession('s1')).rejects.toThrow('请先登录');
  });

  it('clears token and throws on 401 response', async () => {
    localStorage.setItem('dmh_token', 'expired-security-token');
    const removeSpy = vi.spyOn(Storage.prototype, 'removeItem');
    vi.mocked(fetch).mockResolvedValue({ ok: false, status: 401, json: async () => ({ message: 'expired' }) } as Response);

    await expect(securityApi.revokeSession('s1')).rejects.toThrow('登录已过期，请重新登录');
    expect(removeSpy).toHaveBeenCalledWith('dmh_token');
  });

  it('falls back to generic message when json parse fails', async () => {
    localStorage.setItem('dmh_token', 'token-security');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error('invalid json');
      },
    } as Response);

    await expect(securityApi.getPasswordPolicy()).rejects.toThrow('请求失败');
  });
});
