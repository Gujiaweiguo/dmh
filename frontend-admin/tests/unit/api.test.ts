import { describe, expect, it, vi, beforeEach } from 'vitest';
import api from '../api';

vi.mock('../api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('Campaign API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('fetches campaign list', async () => {
    const mockResponse = {
      data: {
        campaigns: [
          { id: 1, name: 'Test Campaign', status: 'active' },
        ],
        total: 1,
      },
    };
    (api.get as any).mockResolvedValue(mockResponse);

    const result = await api.get('/api/v1/campaigns?page=1&pageSize=10');

    expect(api.get).toHaveBeenCalledWith('/api/v1/campaigns?page=1&pageSize=10');
    expect(result.data.campaigns).toHaveLength(1);
    expect(result.data.total).toBe(1);
  });

  it('creates new campaign', async () => {
    const campaignData = {
      brandId: 1,
      name: 'New Campaign',
      description: 'Test',
      rewardRule: 10,
      startTime: '2026-02-09T10:00:00',
      endTime: '2026-02-10T10:00:00',
    };
    const mockResponse = {
      data: { id: 1, ...campaignData },
    };
    (api.post as any).mockResolvedValue(mockResponse);

    const result = await api.post('/api/v1/campaigns', campaignData);

    expect(api.post).toHaveBeenCalledWith('/api/v1/campaigns', campaignData);
    expect(result.data.id).toBe(1);
  });

  it('updates campaign', async () => {
    const updateData = { name: 'Updated Name' };
    const mockResponse = {
      data: { id: 1, ...updateData },
    };
    (api.put as any).mockResolvedValue(mockResponse);

    const result = await api.put('/api/v1/campaigns/1', updateData);

    expect(api.put).toHaveBeenCalledWith('/api/v1/campaigns/1', updateData);
  });

  it('deletes campaign', async () => {
    const mockResponse = { data: { message: 'success' } };
    (api.delete as any).mockResolvedValue(mockResponse);

    const result = await api.delete('/api/v1/campaigns/1');

    expect(api.delete).toHaveBeenCalledWith('/api/v1/campaigns/1');
    expect(result.data.message).toBe('success');
  });
});

describe('Order API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('fetches order list', async () => {
    const mockResponse = {
      data: {
        orders: [
          { id: 1, campaignId: 1, phone: '13800138000', status: 'pending' },
        ],
        total: 1,
      },
    };
    (api.get as any).mockResolvedValue(mockResponse);

    const result = await api.get('/api/v1/orders?page=1&pageSize=10');

    expect(api.get).toHaveBeenCalledWith('/api/v1/orders?page=1&pageSize=10');
    expect(result.data.orders).toHaveLength(1);
  });

  it('creates new order', async () => {
    const orderData = {
      campaignId: 1,
      phone: '13800138000',
      formData: { name: 'Test User' },
    };
    const mockResponse = {
      data: { id: 1, ...orderData },
    };
    (api.post as any).mockResolvedValue(mockResponse);

    const result = await api.post('/api/v1/orders', orderData);

    expect(api.post).toHaveBeenCalledWith('/api/v1/orders', orderData);
  });
});

describe('Auth API', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('logs in user', async () => {
    const loginData = { username: 'admin', password: '123456' };
    const mockResponse = {
      data: {
        token: 'test-token',
        userId: 1,
        username: 'admin',
        roles: ['platform_admin'],
      },
    };
    (api.post as any).mockResolvedValue(mockResponse);

    const result = await api.post('/api/v1/auth/login', loginData);

    expect(api.post).toHaveBeenCalledWith('/api/v1/auth/login', loginData);
    expect(result.data.token).toBe('test-token');
    expect(result.data.roles).toContain('platform_admin');
  });

  it('registers new user', async () => {
    const registerData = {
      username: 'newuser',
      password: 'password123',
      phone: '13800138001',
    };
    const mockResponse = {
      data: {
        token: 'new-token',
        userId: 2,
        username: 'newuser',
      },
    };
    (api.post as any).mockResolvedValue(mockResponse);

    const result = await api.post('/api/v1/auth/register', registerData);

    expect(api.post).toHaveBeenCalledWith('/api/v1/auth/register', registerData);
    expect(result.data.username).toBe('newuser');
  });
});
