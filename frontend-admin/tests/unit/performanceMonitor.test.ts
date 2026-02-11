import { beforeEach, describe, expect, it, vi } from 'vitest';
import PerformanceMonitor from '../../services/performanceMonitor';

type MockWindow = {
	navigator: { sendBeacon: ReturnType<typeof vi.fn> };
	addEventListener: ReturnType<typeof vi.fn>;
};

type MockPerformance = {
	now: ReturnType<typeof vi.fn>;
	timing: { navigationStart: number; loadEventEnd: number };
};

describe('performanceMonitor', () => {
	beforeEach(() => {
		vi.useRealTimers();
		(globalThis as unknown as { window: MockWindow }).window = {
			navigator: { sendBeacon: vi.fn() },
			addEventListener: vi.fn(),
		};
		(globalThis as unknown as { performance: MockPerformance }).performance = {
			now: vi.fn(() => 100),
			timing: { navigationStart: 0, loadEventEnd: 2000 },
		};
	});

	it('measureRender returns fn result', () => {
		(globalThis as unknown as { performance: MockPerformance }).performance.now = vi.fn()
			.mockReturnValueOnce(10)
			.mockReturnValueOnce(35);
		const result = PerformanceMonitor.measureRender('Comp', () => 'ok');
		expect(result).toBe('ok');
	});

	it('logApiRequest sends beacon', () => {
		PerformanceMonitor.logApiRequest('/x', 120, true);
		expect(window.navigator.sendBeacon).toHaveBeenCalledTimes(1);
	});

	it('debounce delays call', async () => {
		vi.useFakeTimers();
		const fn = vi.fn();
		const debounced = PerformanceMonitor.debounce(fn, 100);
		debounced('a');
		debounced('b');
		expect(fn).not.toHaveBeenCalled();
		vi.advanceTimersByTime(110);
		expect(fn).toHaveBeenCalledTimes(1);
		expect(fn).toHaveBeenCalledWith('b');
	});

	it('throttle limits calls', () => {
		vi.useFakeTimers();
		const fn = vi.fn();
		const throttled = PerformanceMonitor.throttle(fn, 100);
		throttled();
		throttled();
		expect(fn).toHaveBeenCalledTimes(1);
		vi.advanceTimersByTime(100);
		throttled();
		expect(fn).toHaveBeenCalledTimes(2);
	});
});
