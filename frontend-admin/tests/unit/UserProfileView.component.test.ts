import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { UserProfileView } from '../../views/UserProfileView';

describe('UserProfileView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should mount without errors', () => {
    const wrapper = mount(UserProfileView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('should render component structure', () => {
    const wrapper = mount(UserProfileView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.html()).toBeDefined();
  });

  it('should have component instance', () => {
    const wrapper = mount(UserProfileView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.vm).toBeDefined();
  });
});
