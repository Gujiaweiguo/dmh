import { beforeEach, describe, expect, it, vi } from 'vitest';
import { api, orderApi, brandApi, feedbackApi } from '../../src/services/api.js';

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

	it('api.get without params', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ code: 0 }) });
		await api.get('/orders/list');
		const [url] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/orders/list');
	});

	it('api.post sends json body', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		await api.post('/feedback', { title: 't1' });
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/feedback');
		expect(config.method).toBe('POST');
		expect(config.body).toBe(JSON.stringify({ title: 't1' }));
	});

	it('api.put sends json body', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		await api.put('/orders/1', { status: 'paid' });
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/orders/1');
		expect(config.method).toBe('PUT');
		expect(config.body).toBe(JSON.stringify({ status: 'paid' }));
	});

	it('api.post keeps FormData body without forcing json content-type', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		const formData = new FormData();
		formData.append('file', new Blob(['a'], { type: 'text/plain' }), 'a.txt');
		await api.post('/material/upload', formData);
		const [, config] = fetch.mock.calls[0];
		expect(config.method).toBe('POST');
		expect(config.body).toBe(formData);
		expect(config.headers['Content-Type']).toBeUndefined();
	});

	it('api.delete sends request', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		await api.delete('/orders/1');
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/orders/1');
		expect(config.method).toBe('DELETE');
	});

	it('api throws on http error', async () => {
		fetch.mockResolvedValue({
			ok: false,
			status: 500,
			text: async () => 'server error',
		});
		await expect(api.delete('/x')).rejects.toThrow('HTTP error! status: 500');
	});

	it('api throws on http error with empty text', async () => {
		fetch.mockResolvedValue({
			ok: false,
			status: 404,
			text: async () => '',
		});
		await expect(api.get('/not-found')).rejects.toThrow('HTTP error! status: 404');
	});

	it('api throws on http error when text() fails', async () => {
		fetch.mockResolvedValue({
			ok: false,
			status: 500,
			text: async () => { throw new Error('text parse error'); },
		});
		await expect(api.get('/error')).rejects.toThrow('HTTP error! status: 500');
	});

	it('api works without token', async () => {
		global.localStorage.getItem = vi.fn(() => null);
		fetch.mockResolvedValue({ ok: true, json: async () => ({ code: 0 }) });
		await api.get('/public');
		const [, config] = fetch.mock.calls[0];
		expect(config.headers.Authorization).toBeUndefined();
	});
});

describe('h5 api service - orderApi', () => {
	beforeEach(() => {
		global.fetch = vi.fn();
		global.localStorage = {
			getItem: vi.fn(() => 'token-123'),
		};
	});

	it('orderApi.getOrders calls correct endpoint', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ orders: [] }) });
		await orderApi.getOrders({ page: 1, pageSize: 10 });
		const [url] = fetch.mock.calls[0];
		expect(url).toContain('/api/v1/orders/list');
	});

	it('orderApi.getOrder calls correct endpoint', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ id: 1 }) });
		await orderApi.getOrder(123);
		const [url] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/orders/123');
	});

	it('orderApi.updateOrderStatus sends put request', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		await orderApi.updateOrderStatus(123, 'verified');
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/orders/123');
		expect(config.method).toBe('PUT');
		expect(config.body).toBe(JSON.stringify({ status: 'verified' }));
	});

	it('orderApi.scanOrderCode sends get with code param', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ order: {} }) });
		await orderApi.scanOrderCode('CODE123');
		const [url] = fetch.mock.calls[0];
		expect(url).toContain('/api/v1/orders/scan?code=CODE123');
	});

	it('orderApi.verifyOrder sends post to /orders/verify', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		await orderApi.verifyOrder('CODE123', 'ok');
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/orders/verify');
		expect(config.method).toBe('POST');
		expect(config.body).toBe(JSON.stringify({ code: 'CODE123', notes: 'ok' }));
	});

	it('orderApi.unverifyOrder sends post to /orders/unverify', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		await orderApi.unverifyOrder('CODE123', 'mistake');
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/orders/unverify');
		expect(config.method).toBe('POST');
		expect(config.body).toBe(JSON.stringify({ code: 'CODE123', reason: 'mistake' }));
	});

	it('orderApi.getVerificationRecords sends get request', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ records: [] }) });
		await orderApi.getVerificationRecords();
		const [url] = fetch.mock.calls[0];
		expect(url).toContain('/api/v1/orders/verification-records');
	});
});

describe('h5 api service - brandApi', () => {
	beforeEach(() => {
		global.fetch = vi.fn();
		global.localStorage = {
			getItem: vi.fn(() => 'token-123'),
		};
	});

	it('brandApi.scanOrderCode calls /orders/scan', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ order: {} }) });
		await brandApi.scanOrderCode('CODE456');
		const [url] = fetch.mock.calls[0];
		expect(url).toContain('/api/v1/orders/scan?code=CODE456');
	});

	it('brandApi.verifyOrderCode sends post to /orders/verify', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		await brandApi.verifyOrderCode('CODE456', 'approved');
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/orders/verify');
		expect(config.body).toBe(JSON.stringify({ code: 'CODE456', notes: 'approved' }));
	});

	it('brandApi.unverifyOrderCode sends post to /orders/unverify', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		await brandApi.unverifyOrderCode('CODE456', 'error');
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/orders/unverify');
		expect(config.body).toBe(JSON.stringify({ code: 'CODE456', reason: 'error' }));
	});

	it('brandApi.getVerificationRecords calls correct endpoint', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ records: [] }) });
		await brandApi.getVerificationRecords();
		const [url] = fetch.mock.calls[0];
		expect(url).toContain('/api/v1/orders/verification-records');
	});
});

describe('h5 api service - feedbackApi', () => {
	beforeEach(() => {
		global.fetch = vi.fn();
		global.localStorage = {
			getItem: vi.fn(() => 'token-123'),
		};
	});

	it('feedbackApi.createFeedback sends post request', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ id: 1 }) });
		await feedbackApi.createFeedback({ title: 'test', content: 'content' });
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/feedback');
		expect(config.method).toBe('POST');
	});

	it('feedbackApi.listFeedback sends get request', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ list: [] }) });
		await feedbackApi.listFeedback({ page: 1 });
		const [url] = fetch.mock.calls[0];
		expect(url).toContain('/api/v1/feedback/list');
	});

	it('feedbackApi.getFeedbackDetail sends get with id param', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ id: 1 }) });
		await feedbackApi.getFeedbackDetail(123);
		const [url] = fetch.mock.calls[0];
		expect(url).toContain('id=123');
	});

	it('feedbackApi.submitSatisfactionSurvey sends post request', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		await feedbackApi.submitSatisfactionSurvey({ rating: 5 });
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/feedback/satisfaction-survey');
		expect(config.method).toBe('POST');
	});

	it('feedbackApi.listFaq sends get request', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ faqs: [] }) });
		await feedbackApi.listFaq({ category: 'general' });
		const [url] = fetch.mock.calls[0];
		expect(url).toContain('/api/v1/feedback/faq');
	});

	it('feedbackApi.markFaqHelpful sends post request', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		await feedbackApi.markFaqHelpful(1, 'helpful');
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/feedback/faq/helpful');
		expect(config.body).toBe(JSON.stringify({ id: 1, type: 'helpful' }));
	});

	it('feedbackApi.recordFeatureUsage sends post request', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ ok: true }) });
		await feedbackApi.recordFeatureUsage({ feature: 'scan', duration: 500 });
		const [url, config] = fetch.mock.calls[0];
		expect(url).toBe('/api/v1/feedback/feature-usage');
		expect(config.method).toBe('POST');
		expect(config.body).toBe(JSON.stringify({ feature: 'scan', duration: 500 }));
	});

	it('feedbackApi.getFeedbackStatistics sends get request', async () => {
		fetch.mockResolvedValue({ ok: true, json: async () => ({ stats: {} }) });
		await feedbackApi.getFeedbackStatistics({ period: 'month' });
		const [url] = fetch.mock.calls[0];
		expect(url).toContain('/api/v1/feedback/statistics');
	});
});
