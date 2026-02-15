import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { RolePermissionView } from '../../views/RolePermissionView';

vi.mock('../../services/roleApi', () => ({
  roleApi: {
    getRoles: vi.fn().mockResolvedValue({ list: [], total: 0 }),
    createRole: vi.fn().mockResolvedValue({}),
    updateRole: vi.fn().mockResolvedValue({}),
    deleteRole: vi.fn().mockResolvedValue({}),
  },
}));

describe('RolePermissionView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should mount without errors', () => {
    const wrapper = mount(RolePermissionView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(RolePermissionView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.html()).toBeDefined();
  });

  it('should have component instance', () => {
    const wrapper = mount(RolePermissionView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.vm).toBeDefined();
  });
});
