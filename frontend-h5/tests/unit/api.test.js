import { beforeEach, describe, expect, it, vi } from 'vitest';
import { api } from '../../src/services/api.js';

describe('h5 api service', () => {
	beforeEach(() => {
		global.fetch = vi.fn();
		global.localStorage = {
			getItem: vi.fn(() => 'token-123'),
		};
	});

	it('api.get builds query string and auth header', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ code: 0 }) });
		await api.get('/orders/list', { page: '1', pageSize: '10' });
		expect(fetch).toHaveBeenCalledTimes(1);
		const [url, config] = fetch.mock.calls[0];
		expect(url).toContain('/api/v1/orders/list?page=1&pageSize=10');
		expect(config.headers.Authorization).toBe('Bearer token-123');
	});

	it('api.post sends json body', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		await api.post('/feedback', { title: 't1' });
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/feedback');
		expect(config.method).toBe('POST');
		expect(config.body).toBe(JSON.stringify({ title: 't1' }));
	});

	it('api throws on http error', async () => {
		fetch.mockResolvedValue({
			ok: false,
			status: 500,
			text: async () => 'server error',
		});
		await expect(api.delete('/x')).rejects.toThrow('HTTP error! status: 500');
	});
});
