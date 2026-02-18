import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import PlatformRewardView, { type PlatformRewardViewInstance } from '../../views/PlatformRewardView';

vi.mock('../../services/axios', () => ({
  default: {
    get: vi.fn().mockResolvedValue({
      data: { code: 200, data: { rewards: [], total: 0, brands: [], campaigns: [] } }
    }),
  }
}));

describe('PlatformRewardView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should mount without errors', () => {
    const wrapper = mount(PlatformRewardView);
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(PlatformRewardView);
    expect(wrapper.html()).toBeDefined();
    expect(wrapper.find('.platform-reward-view').exists()).toBe(true);
  });

  it('should have component instance with required properties', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(vm.rewards).toBeDefined();
    expect(vm.loading).toBeDefined();
    expect(vm.total).toBeDefined();
    expect(vm.currentPage).toBeDefined();
    expect(vm.pageSize).toBeDefined();
    expect(vm.filters).toBeDefined();
    expect(vm.brands).toBeDefined();
    expect(vm.campaigns).toBeDefined();
  });

  it('should have filter properties initialized', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(vm.filters.keyword).toBe('');
    expect(vm.filters.brandId).toBeNull();
    expect(vm.filters.campaignId).toBeNull();
    expect(vm.filters.level).toBeNull();
    expect(vm.filters.status).toBe('');
  });

  it('should have formatAmount method', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(vm.formatAmount(100)).toBe('¥100.00');
    expect(vm.formatAmount(0)).toBe('¥0.00');
    expect(vm.formatAmount(1234.5)).toBe('¥1234.50');
  });

  it('should have formatDate method', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(vm.formatDate('')).toBe('-');
    expect(vm.formatDate(null)).toBe('-');
  });

  it('should have rewardStatusText method', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(vm.rewardStatusText('pending')).toBe('待结算');
    expect(vm.rewardStatusText('settled')).toBe('已结算');
    expect(vm.rewardStatusText('cancelled')).toBe('已取消');
    expect(vm.rewardStatusText('unknown')).toBe('unknown');
  });

  it('should have rewardStatusColor method', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(vm.rewardStatusColor('pending')).toBe('orange');
    expect(vm.rewardStatusColor('settled')).toBe('green');
    expect(vm.rewardStatusColor('cancelled')).toBe('gray');
    expect(vm.rewardStatusColor('unknown')).toBe('gray');
  });

  it('should have levelText method', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(vm.levelText(1)).toBe('一级');
    expect(vm.levelText(2)).toBe('二级');
    expect(vm.levelText(3)).toBe('三级');
    expect(vm.levelText(99)).toBe('99级');
  });

  it('should have handleSearch method', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(typeof vm.handleSearch).toBe('function');
  });

  it('should have handleReset method', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(typeof vm.handleReset).toBe('function');
  });

  it('should have handlePageChange method', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(typeof vm.handlePageChange).toBe('function');
  });

  it('should have viewDetail method', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(typeof vm.viewDetail).toBe('function');
  });

  it('should have handleExport method', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(typeof vm.handleExport).toBe('function');
  });

  it('should have loadRewards method', () => {
    const wrapper = mount(PlatformRewardView);
    const vm = wrapper.vm as unknown as PlatformRewardViewInstance;
    expect(typeof vm.loadRewards).toBe('function');
  });
});
