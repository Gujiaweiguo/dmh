import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import MemberListView, { type MemberListViewInstance } from '../../views/MemberListView';

vi.mock('../../services/memberApi', () => ({
  memberApi: {
    getMembers: vi.fn().mockResolvedValue({ members: [], total: 0 }),
  },
}));

describe('MemberListView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should mount without errors', () => {
    const wrapper = mount(MemberListView);
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(MemberListView);
    expect(wrapper.html()).toBeDefined();
    expect(wrapper.find('.member-list-view').exists()).toBe(true);
  });

  it('should have component instance with required properties', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(vm.members).toBeDefined();
    expect(vm.loading).toBeDefined();
    expect(vm.total).toBeDefined();
    expect(vm.currentPage).toBeDefined();
    expect(vm.pageSize).toBeDefined();
    expect(vm.filters).toBeDefined();
    expect(vm.selectedMembers).toBeDefined();
  });

  it('should have filter properties initialized', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(vm.filters.keyword).toBe('');
    expect(vm.filters.brandId).toBeNull();
    expect(vm.filters.source).toBe('');
    expect(vm.filters.status).toBe('');
  });

  it('should have selectedMembers initialized as empty array', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(Array.isArray(vm.selectedMembers)).toBe(true);
    expect(vm.selectedMembers.length).toBe(0);
  });

  it('should have formatAmount method', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(vm.formatAmount(100)).toBe('¥100.00');
    expect(vm.formatAmount(0)).toBe('¥0.00');
    expect(vm.formatAmount(1234.5)).toBe('¥1234.50');
  });

  it('should have formatDate method', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(vm.formatDate('')).toBe('-');
    expect(vm.formatDate(null)).toBe('-');
  });

  it('should have genderText method', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(vm.genderText(0)).toBe('未知');
    expect(vm.genderText(1)).toBe('男');
    expect(vm.genderText(2)).toBe('女');
    expect(vm.genderText(99)).toBe('未知');
  });

  it('should have statusColor method', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(vm.statusColor('active')).toBe('green');
    expect(vm.statusColor('disabled')).toBe('red');
    expect(vm.statusColor('unknown')).toBe('gray');
  });

  it('should have handleSearch method', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(typeof vm.handleSearch).toBe('function');
    vm.currentPage = 5;
    vm.handleSearch();
    expect(vm.currentPage).toBe(1);
  });

  it('should have handleReset method', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    vm.filters.keyword = 'test';
    vm.filters.status = 'active';
    vm.handleReset();
    expect(vm.filters.keyword).toBe('');
    expect(vm.filters.status).toBe('');
  });

  it('should have handlePageChange method', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(typeof vm.handlePageChange).toBe('function');
  });

  it('should have viewDetail method', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(typeof vm.viewDetail).toBe('function');
  });

  it('should have handleExport method', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(typeof vm.handleExport).toBe('function');
  });

  it('should have handleMerge method', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(typeof vm.handleMerge).toBe('function');
  });

  it('should have loadMembers method', () => {
    const wrapper = mount(MemberListView);
    const vm = wrapper.vm as unknown as MemberListViewInstance;
    expect(typeof vm.loadMembers).toBe('function');
  });
});
