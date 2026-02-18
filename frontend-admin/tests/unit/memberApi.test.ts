import { beforeEach, describe, expect, it, vi } from 'vitest';
import { memberApi } from '../../services/memberApi';

describe('memberApi service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    Object.defineProperty(globalThis, 'fetch', {
      value: vi.fn(),
      writable: true,
    });
  });

  it('getMembers builds query with array params', async () => {
    localStorage.setItem('dmh_token', 'token-member');
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ total: 0, members: [] }),
    } as unknown as Response);

    await memberApi.getMembers({
      page: 2,
      pageSize: 20,
      keyword: 'alice',
      tagIds: [1, 3],
      status: 'active',
    });

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    const urlStr = String(url);
    const headers = (options as RequestInit).headers as Record<string, string>;
    expect(urlStr).toContain('/api/v1/members?');
    expect(urlStr).toContain('page=2');
    expect(urlStr).toContain('pageSize=20');
    expect(urlStr).toContain('keyword=alice');
    expect(urlStr).toContain('status=active');
    expect(urlStr).toContain('tagIds=1');
    expect(urlStr).toContain('tagIds=3');
    expect(headers.Authorization).toBe('Bearer token-member');
  });

  it('throws login error when token missing', async () => {
    const removeSpy = vi.spyOn(Storage.prototype, 'removeItem');

    await expect(memberApi.getMember(1)).rejects.toThrow('请先登录');
    expect(removeSpy).toHaveBeenCalledWith('dmh_token');
  });

  it('clears token and throws on 401 response', async () => {
    localStorage.setItem('dmh_token', 'expired-token');
    const removeSpy = vi.spyOn(Storage.prototype, 'removeItem');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 401,
      json: async () => ({ message: 'expired' }),
    } as unknown as Response);

    await expect(memberApi.getMember(11)).rejects.toThrow('登录已过期，请重新登录');
    expect(removeSpy).toHaveBeenCalledWith('dmh_token');
  });

  it('throws backend json message when request fails', async () => {
    localStorage.setItem('dmh_token', 'token-member');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => ({ message: 'server error' }),
    } as unknown as Response);

    await expect(memberApi.getMember(2)).rejects.toThrow('server error');
  });

  it('falls back to generic message when error json parse fails', async () => {
    localStorage.setItem('dmh_token', 'token-member');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error('invalid json');
      },
    } as unknown as Response);

    await expect(memberApi.getMember(3)).rejects.toThrow('请求失败');
  });

  it('createTag sends POST with body', async () => {
    localStorage.setItem('dmh_token', 'token-member');
    const payload = { name: 'vip', category: 'value', color: '#ff0000' };
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ id: 1, ...payload }),
    } as unknown as Response);

    await memberApi.createTag(payload);

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    expect(url).toBe('/api/v1/members/tags');
    expect((options as RequestInit).method).toBe('POST');
    expect((options as RequestInit).body).toBe(JSON.stringify(payload));
  });

  it('getExportRequests appends query params', async () => {
    localStorage.setItem('dmh_token', 'token-member');
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ total: 0, items: [] }),
    } as unknown as Response);

    await memberApi.getExportRequests({ page: 1, pageSize: 10, brandId: 99, status: 'pending' });

    const [url] = vi.mocked(fetch).mock.calls[0];
    const urlStr = String(url);
    expect(urlStr).toContain('/api/v1/members/export-requests?');
    expect(urlStr).toContain('page=1');
    expect(urlStr).toContain('pageSize=10');
    expect(urlStr).toContain('brandId=99');
    expect(urlStr).toContain('status=pending');
  });
});
