import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { BrandManagementView } from '../../views/BrandManagementView';

vi.mock('../../components/PermissionGuard', () => ({
  PermissionGuard: {
    setup: () => ({ slots: { default: () => null } }),
  },
  usePermission: () => ({
    hasPermission: vi.fn().mockReturnValue(true),
    canAccessBrand: vi.fn().mockReturnValue(true),
    user: { value: { roles: ['platform_admin'] } },
  }),
}));

describe('BrandManagementView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should mount without errors', () => {
    const wrapper = mount(BrandManagementView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
          PermissionGuard: true,
        },
      },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(BrandManagementView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
          PermissionGuard: true,
        },
      },
    });
    expect(wrapper.html()).toBeDefined();
  });

  it('should have component instance', () => {
    const wrapper = mount(BrandManagementView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
          PermissionGuard: true,
        },
      },
    });
    expect(wrapper.vm).toBeDefined();
  });
});
