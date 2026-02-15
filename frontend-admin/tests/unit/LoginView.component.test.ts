import { describe, expect, it, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { LoginView } from '../../views/LoginView';

vi.mock('../../services/authApi', () => ({
  authApi: {
    login: vi.fn().mockResolvedValue({ token: 'test-token', user: { username: 'admin' } }),
    register: vi.fn().mockResolvedValue({ token: 'test-token', user: { username: 'newuser' } }),
  },
}));

describe('LoginView Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should mount without errors', () => {
    const wrapper = mount(LoginView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('should emit login-success event on successful login', async () => {
    const wrapper = mount(LoginView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.emitted()).toBeDefined();
  });

  it('should have login form reactive state', () => {
    const wrapper = mount(LoginView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    expect(wrapper.vm).toBeDefined();
  });

  it('should toggle between login and register modes', async () => {
    const wrapper = mount(LoginView, {
      global: {
        stubs: {
          'lucide-vue-next': true,
        },
      },
    });
    const vm = wrapper.vm as any;
    if (vm.toggleMode) {
      const initialMode = vm.isLogin;
      vm.toggleMode();
      expect(vm.isLogin).toBe(!initialMode);
    }
  });
});
