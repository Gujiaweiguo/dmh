import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import MemberMergeView, { type MemberMergeViewInstance } from '../../views/MemberMergeView';

vi.mock('../../services/memberApi', () => ({
  memberApi: {
    previewMerge: vi.fn().mockResolvedValue({
      canMerge: true,
      sourceMember: { id: 1, nickname: 'source' },
      targetMember: { id: 2, nickname: 'target' },
      conflicts: [],
    }),
    mergeMember: vi.fn().mockResolvedValue({}),
  },
}));

describe('MemberMergeView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    window.location.hash = '#/members/merge?source=1&target=2';
  });

  it('should mount without errors', () => {
    const wrapper = mount(MemberMergeView);
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(MemberMergeView);
    expect(wrapper.html()).toBeDefined();
    expect(wrapper.find('.member-merge-view').exists()).toBe(true);
  });

  it('should have component instance with required properties', () => {
    const wrapper = mount(MemberMergeView);
    const vm = wrapper.vm as unknown as MemberMergeViewInstance;
    expect(vm.loading).toBeDefined();
    expect(vm.sourceMemberId).toBeDefined();
    expect(vm.targetMemberId).toBeDefined();
    expect(vm.reason).toBeDefined();
    expect(vm.preview).toBeDefined();
    expect(vm.showPreview).toBeDefined();
  });

  it('should have reason initialized as empty string', () => {
    const wrapper = mount(MemberMergeView);
    const vm = wrapper.vm as unknown as MemberMergeViewInstance;
    expect(vm.reason).toBe('');
  });

  it('should have showPreview initialized as false', () => {
    const wrapper = mount(MemberMergeView);
    const vm = wrapper.vm as unknown as MemberMergeViewInstance;
    expect(vm.showPreview).toBe(false);
  });

  it('should have preview initialized as null', () => {
    const wrapper = mount(MemberMergeView);
    const vm = wrapper.vm as unknown as MemberMergeViewInstance;
    expect(vm.preview).toBeNull();
  });

  it('should have formatAmount method', () => {
    const wrapper = mount(MemberMergeView);
    const vm = wrapper.vm as unknown as MemberMergeViewInstance;
    expect(vm.formatAmount(100)).toBe('¥100.00');
    expect(vm.formatAmount(0)).toBe('¥0.00');
    expect(vm.formatAmount(1234.5)).toBe('¥1234.50');
  });

  it('should have handlePreview method', () => {
    const wrapper = mount(MemberMergeView);
    const vm = wrapper.vm as unknown as MemberMergeViewInstance;
    expect(typeof vm.handlePreview).toBe('function');
  });

  it('should have handleMerge method', () => {
    const wrapper = mount(MemberMergeView);
    const vm = wrapper.vm as unknown as MemberMergeViewInstance;
    expect(typeof vm.handleMerge).toBe('function');
  });

  it('should have goBack method', () => {
    const wrapper = mount(MemberMergeView);
    const vm = wrapper.vm as unknown as MemberMergeViewInstance;
    expect(typeof vm.goBack).toBe('function');
  });
});
