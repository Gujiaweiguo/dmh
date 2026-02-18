import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import MemberExportView, { type MemberExportViewInstance } from '../../views/MemberExportView';

vi.mock('../../services/memberApi', () => ({
  memberApi: {
    getExportRequests: vi.fn().mockResolvedValue({ requests: [], total: 0 }),
    createExportRequest: vi.fn().mockResolvedValue({}),
    approveExportRequest: vi.fn().mockResolvedValue({}),
  },
}));

describe('MemberExportView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should mount without errors', () => {
    const wrapper = mount(MemberExportView);
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(MemberExportView);
    expect(wrapper.html()).toBeDefined();
    expect(wrapper.find('.member-export-view').exists()).toBe(true);
  });

  it('should have component instance with required properties', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(vm.loading).toBeDefined();
    expect(vm.exportRequests).toBeDefined();
    expect(vm.total).toBeDefined();
    expect(vm.currentPage).toBeDefined();
    expect(vm.pageSize).toBeDefined();
    expect(vm.showCreateDialog).toBeDefined();
    expect(vm.createForm).toBeDefined();
    expect(vm.showApproveDialog).toBeDefined();
    expect(vm.approveForm).toBeDefined();
  });

  it('should have createForm initialized correctly', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(vm.createForm.brandId).toBeNull();
    expect(vm.createForm.reason).toBe('');
    expect(vm.createForm.filters).toBe('');
  });

  it('should have approveForm initialized correctly', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(vm.approveForm.requestId).toBe(0);
    expect(vm.approveForm.approve).toBe(true);
    expect(vm.approveForm.reason).toBe('');
  });

  it('should have formatDate method', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(vm.formatDate('')).toBe('-');
    expect(vm.formatDate(null)).toBe('-');
  });

  it('should have statusColor method', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(vm.statusColor('pending')).toBe('orange');
    expect(vm.statusColor('approved')).toBe('blue');
    expect(vm.statusColor('rejected')).toBe('red');
    expect(vm.statusColor('completed')).toBe('green');
    expect(vm.statusColor('unknown')).toBe('gray');
  });

  it('should have statusText method', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(vm.statusText('pending')).toBe('待审批');
    expect(vm.statusText('approved')).toBe('已批准');
    expect(vm.statusText('rejected')).toBe('已驳回');
    expect(vm.statusText('completed')).toBe('已完成');
    expect(vm.statusText('unknown')).toBe('unknown');
  });

  it('should have handleCreate method', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(typeof vm.handleCreate).toBe('function');
  });

  it('should have openApproveDialog method', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(typeof vm.openApproveDialog).toBe('function');
  });

  it('should have handleApprove method', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(typeof vm.handleApprove).toBe('function');
  });

  it('should have handleDownload method', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(typeof vm.handleDownload).toBe('function');
  });

  it('should have handlePageChange method', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(typeof vm.handlePageChange).toBe('function');
  });

  it('should have goBack method', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(typeof vm.goBack).toBe('function');
  });

  it('should have loadExportRequests method', () => {
    const wrapper = mount(MemberExportView);
    const vm = wrapper.vm as unknown as MemberExportViewInstance;
    expect(typeof vm.loadExportRequests).toBe('function');
  });
});
