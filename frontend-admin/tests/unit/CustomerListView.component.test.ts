import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import CustomerListView, { type CustomerListViewInstance } from '../../views/CustomerListView';

vi.mock('../../services/axios', () => ({
  default: {
    get: vi.fn().mockResolvedValue({
      data: { code: 200, data: { customers: [], total: 0, brands: [], campaigns: [] } }
    })
  }
}));

describe('CustomerListView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should mount without errors', () => {
    const wrapper = mount(CustomerListView);
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(CustomerListView);
    expect(wrapper.html()).toBeDefined();
    expect(wrapper.find('.customer-list-view').exists()).toBe(true);
  });

  it('should have component instance with required properties', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    expect(vm.customers).toBeDefined();
    expect(vm.loading).toBeDefined();
    expect(vm.total).toBeDefined();
    expect(vm.currentPage).toBeDefined();
    expect(vm.filters).toBeDefined();
    expect(vm.brands).toBeDefined();
    expect(vm.campaigns).toBeDefined();
  });

  it('should have filter properties initialized', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    expect(vm.filters.keyword).toBe('');
    expect(vm.filters.brandId).toBeNull();
    expect(vm.filters.campaignId).toBeNull();
    expect(vm.filters.status).toBe('');
  });

  it('should have formatAmount method', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    expect(vm.formatAmount(100)).toBe('¥100.00');
    expect(vm.formatAmount(0)).toBe('¥0.00');
    expect(vm.formatAmount(1234.5)).toBe('¥1234.50');
  });

  it('should have formatDate method', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    expect(vm.formatDate('')).toBe('-');
    expect(vm.formatDate(null)).toBe('-');
  });

  it('should have paymentStatusText method', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    expect(vm.paymentStatusText('unpaid')).toBe('未支付');
    expect(vm.paymentStatusText('paid')).toBe('已支付');
    expect(vm.paymentStatusText('refunded')).toBe('已退款');
    expect(vm.paymentStatusText('unknown')).toBe('unknown');
  });

  it('should have paymentStatusColor method', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    expect(vm.paymentStatusColor('unpaid')).toBe('orange');
    expect(vm.paymentStatusColor('paid')).toBe('green');
    expect(vm.paymentStatusColor('refunded')).toBe('gray');
    expect(vm.paymentStatusColor('unknown')).toBe('gray');
  });

  it('should have handleSearch method', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    expect(typeof vm.handleSearch).toBe('function');
    vm.currentPage = 5;
    vm.handleSearch();
    expect(vm.currentPage).toBe(1);
  });

  it('should have handleReset method', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    vm.filters.keyword = 'test';
    vm.filters.status = 'paid';
    vm.handleReset();
    expect(vm.filters.keyword).toBe('');
    expect(vm.filters.status).toBe('');
  });

  it('should have handlePageChange method', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    expect(typeof vm.handlePageChange).toBe('function');
  });

  it('should have viewDetail method', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    expect(typeof vm.viewDetail).toBe('function');
  });

  it('should have loadCustomers method', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    expect(typeof vm.loadCustomers).toBe('function');
  });

  it('should have loadBrands method', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    expect(typeof vm.loadBrands).toBe('function');
  });

  it('should have loadCampaigns method', () => {
    const wrapper = mount(CustomerListView);
    const vm = wrapper.vm as unknown as CustomerListViewInstance;
    expect(typeof vm.loadCampaigns).toBe('function');
  });
});
