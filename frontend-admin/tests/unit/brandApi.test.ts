import { beforeEach, describe, expect, it, vi } from 'vitest';
import { brandApi } from '../../services/brandApi';

describe('brandApi service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    Object.defineProperty(globalThis, 'fetch', {
      value: vi.fn(),
      writable: true,
    });
  });

  it('getBrands returns empty array when payload misses brands', async () => {
    localStorage.setItem('dmh_token', 'token-brand');
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ total: 1 }),
    } as Response);

    const result = await brandApi.getBrands();
    expect(result).toEqual([]);
  });

  it('createBrand sends authorized POST request', async () => {
    localStorage.setItem('dmh_token', 'token-brand');
    const payload = { name: 'Brand A', logo: 'logo', description: 'desc' };
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ id: 1, ...payload, status: 'active', createdAt: '2026-01-01 00:00:00' }),
    } as Response);

    await brandApi.createBrand(payload);

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    const headers = (options as RequestInit).headers as Record<string, string>;
    expect(url).toBe('/api/v1/brands');
    expect((options as RequestInit).method).toBe('POST');
    expect(headers.Authorization).toBe('Bearer token-brand');
    expect((options as RequestInit).body).toBe(JSON.stringify(payload));
  });

  it('throws login error when token missing', async () => {
    await expect(brandApi.getBrandStats(1)).rejects.toThrow('请先登录');
  });

  it('clears token and throws on 401 response', async () => {
    localStorage.setItem('dmh_token', 'expired-token');
    const removeSpy = vi.spyOn(Storage.prototype, 'removeItem');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 401,
      json: async () => ({ message: 'expired' }),
    } as Response);

    await expect(brandApi.getBrandAssets(1)).rejects.toThrow('登录已过期，请重新登录');
    expect(removeSpy).toHaveBeenCalledWith('dmh_token');
  });

  it('throws backend json message when request fails', async () => {
    localStorage.setItem('dmh_token', 'token-brand');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => ({ message: 'brand server error' }),
    } as Response);

    await expect(brandApi.updateBrand(3, { name: 'next' })).rejects.toThrow('brand server error');
  });

  it('falls back to generic message when json parse fails', async () => {
    localStorage.setItem('dmh_token', 'token-brand');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error('invalid json');
      },
    } as Response);

    await expect(brandApi.deleteBrandAsset(1, 2)).rejects.toThrow('请求失败');
  });
});
