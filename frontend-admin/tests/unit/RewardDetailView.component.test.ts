import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import RewardDetailView, { type RewardDetailViewInstance } from '../../views/RewardDetailView';

vi.mock('../../services/axios', () => ({
  default: {
    get: vi.fn().mockResolvedValue({
      data: { code: 200, data: { rewards: [], total: 0, brands: [], campaigns: [] } }
    }),
  }
}));

describe('RewardDetailView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should mount without errors', () => {
    const wrapper = mount(RewardDetailView);
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(RewardDetailView);
    expect(wrapper.html()).toBeDefined();
    expect(wrapper.find('.reward-detail-view').exists()).toBe(true);
  });

  it('should have component instance with required properties', () => {
    const wrapper = mount(RewardDetailView);
    const vm = wrapper.vm as unknown as RewardDetailViewInstance;
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
    const wrapper = mount(RewardDetailView);
    const vm = wrapper.vm as unknown as RewardDetailViewInstance;
    expect(vm.filters.keyword).toBe('');
    expect(vm.filters.brandId).toBeNull();
    expect(vm.filters.campaignId).toBeNull();
    expect(vm.filters.status).toBe('');
  });

  it('should have formatAmount method', () => {
    const wrapper = mount(RewardDetailView);
    const vm = wrapper.vm as unknown as RewardDetailViewInstance;
    expect(vm.formatAmount(100)).toBe('¥100.00');
    expect(vm.formatAmount(0)).toBe('¥0.00');
    expect(vm.formatAmount(1234.5)).toBe('¥1234.50');
  });

  it('should have formatDate method', () => {
    const wrapper = mount(RewardDetailView);
    const vm = wrapper.vm as unknown as RewardDetailViewInstance;
    expect(vm.formatDate('')).toBe('-');
    expect(vm.formatDate(null)).toBe('-');
  });

  it('should have rewardStatusText method', () => {
    const wrapper = mount(RewardDetailView);
    const vm = wrapper.vm as unknown as RewardDetailViewInstance;
    expect(vm.rewardStatusText('pending')).toBe('待结算');
    expect(vm.rewardStatusText('settled')).toBe('已结算');
    expect(vm.rewardStatusText('cancelled')).toBe('已取消');
    expect(vm.rewardStatusText('unknown')).toBe('unknown');
  });

  it('should have rewardStatusColor method', () => {
    const wrapper = mount(RewardDetailView);
    const vm = wrapper.vm as unknown as RewardDetailViewInstance;
    expect(vm.rewardStatusColor('pending')).toBe('orange');
    expect(vm.rewardStatusColor('settled')).toBe('green');
    expect(vm.rewardStatusColor('cancelled')).toBe('gray');
    expect(vm.rewardStatusColor('unknown')).toBe('gray');
  });

  it('should have handleSearch method', () => {
    const wrapper = mount(RewardDetailView);
    const vm = wrapper.vm as unknown as RewardDetailViewInstance;
    expect(typeof vm.handleSearch).toBe('function');
  });

  it('should have handleReset method', () => {
    const wrapper = mount(RewardDetailView);
    const vm = wrapper.vm as unknown as RewardDetailViewInstance;
    expect(typeof vm.handleReset).toBe('function');
  });

  it('should have handlePageChange method', () => {
    const wrapper = mount(RewardDetailView);
    const vm = wrapper.vm as unknown as RewardDetailViewInstance;
    expect(typeof vm.handlePageChange).toBe('function');
  });

  it('should have viewDetail method', () => {
    const wrapper = mount(RewardDetailView);
    const vm = wrapper.vm as unknown as RewardDetailViewInstance;
    expect(typeof vm.viewDetail).toBe('function');
  });

  it('should have handleExport method', () => {
    const wrapper = mount(RewardDetailView);
    const vm = wrapper.vm as unknown as RewardDetailViewInstance;
    expect(typeof vm.handleExport).toBe('function');
  });

  it('should have loadRewards method', () => {
    const wrapper = mount(RewardDetailView);
    const vm = wrapper.vm as unknown as RewardDetailViewInstance;
    expect(typeof vm.loadRewards).toBe('function');
  });
});
