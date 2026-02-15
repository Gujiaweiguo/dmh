import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { DistributorManagementView } from '../../views/DistributorManagementView';

vi.mock('../../services/distributorApi', () => ({
  distributorApi: {
    getDistributors: vi.fn().mockResolvedValue({ list: [], total: 0 }),
    updateDistributorStatus: vi.fn().mockResolvedValue({}),
  },
}));

describe('DistributorManagementView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should mount without errors', () => {
    const wrapper = mount(DistributorManagementView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(DistributorManagementView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.html()).toBeDefined();
  });

  it('should have component instance', () => {
    const wrapper = mount(DistributorManagementView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.vm).toBeDefined();
  });
});
