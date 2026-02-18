import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import MemberDetailView, { type MemberDetailViewInstance } from '../../views/MemberDetailView';

vi.mock('../../services/memberApi', () => ({
  memberApi: {
    getMember: vi.fn().mockResolvedValue({ id: 1, nickname: 'test', totalPayment: 100, totalReward: 50 }),
    addMemberTags: vi.fn().mockResolvedValue({}),
  },
}));

describe('MemberDetailView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    window.location.hash = '#/members/1';
  });

  it('should mount without errors', () => {
    const wrapper = mount(MemberDetailView);
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(MemberDetailView);
    expect(wrapper.html()).toBeDefined();
    expect(wrapper.find('.member-detail-view').exists()).toBe(true);
  });

  it('should have formatAmount method', () => {
    const wrapper = mount(MemberDetailView);
    const vm = wrapper.vm as unknown as MemberDetailViewInstance;
    expect(vm.formatAmount(100)).toBe('¥100.00');
    expect(vm.formatAmount(0)).toBe('¥0.00');
    expect(vm.formatAmount(1234.5)).toBe('¥1234.50');
  });

  it('should have formatDate method', () => {
    const wrapper = mount(MemberDetailView);
    const vm = wrapper.vm as unknown as MemberDetailViewInstance;
    expect(vm.formatDate('')).toBe('-');
    expect(vm.formatDate(null)).toBe('-');
  });

  it('should have genderText method', () => {
    const wrapper = mount(MemberDetailView);
    const vm = wrapper.vm as unknown as MemberDetailViewInstance;
    expect(vm.genderText(0)).toBe('未知');
    expect(vm.genderText(1)).toBe('男');
    expect(vm.genderText(2)).toBe('女');
  });

  it('should have formatFormData method', () => {
    const wrapper = mount(MemberDetailView);
    const vm = wrapper.vm as unknown as MemberDetailViewInstance;
    expect(vm.formatFormData(null)).toBe('-');
    expect(vm.formatFormData({})).toBe('-');
    expect(vm.formatFormData({ name: 'test', phone: '123' })).toBe('name:test、phone:123');
  });

  it('should have orderStatusText method', () => {
    const wrapper = mount(MemberDetailView);
    const vm = wrapper.vm as unknown as MemberDetailViewInstance;
    expect(vm.orderStatusText('pending')).toBe('待支付');
    expect(vm.orderStatusText('paid')).toBe('已支付');
    expect(vm.orderStatusText('cancelled')).toBe('已取消');
    expect(vm.orderStatusText('')).toBe('-');
    expect(vm.orderStatusText(null)).toBe('-');
  });

  it('should have handleAddTags method', () => {
    const wrapper = mount(MemberDetailView);
    const vm = wrapper.vm as unknown as MemberDetailViewInstance;
    expect(typeof vm.handleAddTags).toBe('function');
  });

  it('should have goBack method', () => {
    const wrapper = mount(MemberDetailView);
    const vm = wrapper.vm as unknown as MemberDetailViewInstance;
    expect(typeof vm.goBack).toBe('function');
  });
});
