import { beforeEach, describe, expect, it, vi } from 'vitest';
import { campaignApi } from '../../services/campaignApi';

describe('campaignApi service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
    Object.defineProperty(globalThis, 'fetch', {
      value: vi.fn(),
      writable: true,
    });
  });

  it('createCampaign sends auth header and payload', async () => {
    localStorage.setItem('dmh_token', 'token-abc');
    const requestData = {
      name: 'Campaign A',
      description: 'desc',
      formFields: ['name', 'phone'],
      rewardRule: 10,
      startTime: '2026-02-13T10:00:00Z',
      endTime: '2026-02-20T10:00:00Z',
      paymentConfig: 'wechat',
      enableDistribution: true,
      distributionLevel: 2,
      distributionRewards: '5,2',
      posterTemplateId: 1,
    };

    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => ({ id: 1, ...requestData }),
    } as Response);

    await campaignApi.createCampaign(requestData);

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    const headers = (options as RequestInit).headers as Record<string, string>;
    expect(url).toBe('/api/v1/campaigns');
    expect((options as RequestInit).method).toBe('POST');
    expect(headers.Authorization).toBe('Bearer token-abc');
    expect(JSON.parse((options as RequestInit).body as string)).toEqual(requestData);
  });

  it('getCampaigns builds expected query parameters', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => ({ total: 0, campaigns: [] }),
    } as Response);

    await campaignApi.getCampaigns(2, 15, 'active', 'spring');

    const [url] = vi.mocked(fetch).mock.calls[0];
    expect(String(url)).toContain('/api/v1/campaigns?page=2&pageSize=15&status=active&keyword=spring');
  });

  it('getCampaigns throws readable error on failure', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      text: async () => 'db down',
    } as Response);

    await expect(campaignApi.getCampaigns()).rejects.toThrow('Failed to fetch campaigns: 500 db down');
  });

  it('updateCampaign sends PUT request without auth when token missing', async () => {
    const updateData = {
      name: 'Campaign B',
      description: 'updated',
      formFields: ['name'],
      rewardRule: 20,
      startTime: '2026-02-13T10:00:00Z',
      endTime: '2026-03-01T10:00:00Z',
      status: 'active' as const,
    };

    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => ({ id: 2, ...updateData }),
    } as Response);

    await campaignApi.updateCampaign(2, updateData);

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    const headers = (options as RequestInit).headers as Record<string, string>;
    expect(url).toBe('/api/v1/campaigns/2');
    expect((options as RequestInit).method).toBe('PUT');
    expect(headers.Authorization).toBeUndefined();
  });

  it('deleteCampaign throws readable error on failure', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 403,
      text: async () => 'forbidden',
    } as Response);

    await expect(campaignApi.deleteCampaign(9)).rejects.toThrow('Failed to delete campaign: 403 forbidden');
  });

  it('createCampaign throws readable error on failure', async () => {
    localStorage.setItem('dmh_token', 'token-abc');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 400,
      text: async () => 'invalid data',
    } as Response);

    await expect(campaignApi.createCampaign({
      name: 'Test',
      description: 'test',
      formFields: [],
      rewardRule: 0,
      startTime: '2026-01-01',
      endTime: '2026-01-02',
    })).rejects.toThrow('Failed to create campaign: 400 invalid data');
  });

  it('getCampaign fetches single campaign by id', async () => {
    localStorage.setItem('dmh_token', 'token-abc');
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => ({ id: 1, name: 'Campaign 1' }),
    } as Response);

    const result = await campaignApi.getCampaign(1);

    const [url, options] = vi.mocked(fetch).mock.calls[0];
    expect(url).toBe('/api/v1/campaigns/1');
    expect((options as RequestInit).headers).toHaveProperty('Authorization');
    expect(result).toEqual({ id: 1, name: 'Campaign 1' });
  });

  it('getCampaign throws readable error on failure', async () => {
    localStorage.setItem('dmh_token', 'token-abc');
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 404,
      text: async () => 'not found',
    } as Response);

    await expect(campaignApi.getCampaign(999)).rejects.toThrow('Failed to fetch campaign: 404 not found');
  });

  it('updateCampaign throws readable error on failure', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 400,
      text: async () => 'validation failed',
    } as Response);

    await expect(campaignApi.updateCampaign(1, {
      name: 'Test',
      description: 'test',
      formFields: [],
      rewardRule: 0,
      startTime: '2026-01-01',
      endTime: '2026-01-02',
      status: 'active',
    })).rejects.toThrow('Failed to update campaign: 400 validation failed');
  });
});
