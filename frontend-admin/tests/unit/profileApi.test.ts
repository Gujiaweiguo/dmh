import { beforeEach, describe, expect, it, vi } from 'vitest';
import { profileApi } from '../../services/profileApi';

describe('profileApi service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    Object.defineProperty(globalThis, 'fetch', {
      value: vi.fn(),
      writable: true,
    });
  });

  it('getUserInfo sends auth GET request', async () => {
    localStorage.setItem('dmh_token', 'token-profile');
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ id: 1, username: 'u', phone: '1', email: '', realName: '', avatar: '', status: 'active', roles: [], createdAt: '' }),
    } as unknown as Response);

    await profileApi.getUserInfo();

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    const headers = (options as RequestInit).headers as Record<string, string>;
    expect(url).toBe('/api/v1/auth/userinfo');
    expect((options as RequestInit).method).toBe('GET');
    expect(headers.Authorization).toBe('Bearer token-profile');
  });

  it('writes expected payloads for profile operations', async () => {
    localStorage.setItem('dmh_token', 'token-profile');
    vi.mocked(fetch).mockResolvedValue({ ok: true, status: 200, json: async () => ({ message: 'ok' }) } as unknown as Response);

    await profileApi.updateProfile('Alice');
    await profileApi.sendPhoneCode('13800138000');
    await profileApi.sendEmailCode('a@b.com');
    await profileApi.bindPhone('13800138000', '1234');
    await profileApi.bindEmail('a@b.com', '5678');
    await profileApi.changePassword('oldpass', 'newpass');

    expect(vi.mocked(fetch).mock.calls[0][0]).toBe('/api/v1/users/profile');
    expect((vi.mocked(fetch).mock.calls[0][1] as RequestInit).body).toBe(JSON.stringify({ realName: 'Alice' }));

    expect(vi.mocked(fetch).mock.calls[1][0]).toBe('/api/v1/auth/send-phone-code');
    expect((vi.mocked(fetch).mock.calls[1][1] as RequestInit).body).toBe(
      JSON.stringify({ target: '13800138000', type: 'phone' }),
    );

    expect(vi.mocked(fetch).mock.calls[2][0]).toBe('/api/v1/auth/send-email-code');
    expect((vi.mocked(fetch).mock.calls[2][1] as RequestInit).body).toBe(
      JSON.stringify({ target: 'a@b.com', type: 'email' }),
    );

    expect(vi.mocked(fetch).mock.calls[3][0]).toBe('/api/v1/users/bind-phone');
    expect((vi.mocked(fetch).mock.calls[3][1] as RequestInit).body).toBe(
      JSON.stringify({ phone: '13800138000', code: '1234' }),
    );

    expect(vi.mocked(fetch).mock.calls[4][0]).toBe('/api/v1/users/bind-email');
    expect((vi.mocked(fetch).mock.calls[4][1] as RequestInit).body).toBe(
      JSON.stringify({ email: 'a@b.com', code: '5678' }),
    );

    expect(vi.mocked(fetch).mock.calls[5][0]).toBe('/api/v1/users/change-password');
    expect((vi.mocked(fetch).mock.calls[5][1] as RequestInit).body).toBe(
      JSON.stringify({ oldPassword: 'oldpass', newPassword: 'newpass' }),
    );
  });

  it('throws login error when token missing', async () => {
    await expect(profileApi.getUserInfo()).rejects.toThrow('请先登录');
  });

  it('clears token and throws on 401', async () => {
    localStorage.setItem('dmh_token', 'expired-profile-token');
    const removeSpy = vi.spyOn(Storage.prototype, 'removeItem');
    vi.mocked(fetch).mockResolvedValue({ ok: false, status: 401, json: async () => ({ message: 'expired' }) } as unknown as Response);

    await expect(profileApi.getUserInfo()).rejects.toThrow('登录已过期，请重新登录');
    expect(removeSpy).toHaveBeenCalledWith('dmh_token');
  });

  it('falls back to generic message when error json parse fails', async () => {
    localStorage.setItem('dmh_token', 'token-profile');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error('invalid json');
      },
    } as unknown as Response);

    await expect(profileApi.getUserInfo()).rejects.toThrow('请求失败');
  });
});
