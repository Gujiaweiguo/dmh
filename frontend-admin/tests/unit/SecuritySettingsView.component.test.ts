import { beforeEach, describe, expect, it, vi } from 'vitest';
import { flushPromises, mount } from '@vue/test-utils';
import { SecuritySettingsView } from '../../views/SecuritySettingsView';
import { securityApi } from '../../services/securityApi';

vi.mock('../../services/securityApi', () => ({
  securityApi: {
    getPasswordPolicy: vi.fn(),
    updatePasswordPolicy: vi.fn(),
    getUserSessions: vi.fn(),
    revokeSession: vi.fn(),
    forceLogoutUser: vi.fn(),
    getSecurityEvents: vi.fn(),
    handleSecurityEvent: vi.fn(),
  },
}));

const defaultPolicy = {
  id: 1,
  minLength: 8,
  requireUppercase: true,
  requireLowercase: true,
  requireNumbers: true,
  requireSpecialChars: true,
  maxAge: 90,
  historyCount: 5,
  maxLoginAttempts: 5,
  lockoutDuration: 30,
  sessionTimeout: 480,
  maxConcurrentSessions: 3,
};

const findButtonByText = (wrapper: ReturnType<typeof mount>, text: string) =>
  wrapper.findAll('button').find((button) => button.text().includes(text));

describe('SecuritySettingsView component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.stubGlobal('confirm', vi.fn(() => true));
    vi.stubGlobal('prompt', vi.fn(() => 'ops note'));

    vi.mocked(securityApi.getPasswordPolicy).mockResolvedValue(defaultPolicy);
    vi.mocked(securityApi.updatePasswordPolicy).mockResolvedValue(defaultPolicy);
    vi.mocked(securityApi.getUserSessions).mockResolvedValue({
      total: 1,
      sessions: [
        {
          id: 'session-1',
          userId: 101,
          clientIp: '127.0.0.1',
          userAgent: 'ua',
          loginAt: '2026-02-17 12:00:00',
          lastActiveAt: '2026-02-17 12:10:00',
          expiresAt: '2026-02-17 14:00:00',
          status: 'active',
          createdAt: '2026-02-17 12:00:00',
        },
      ],
    });
    vi.mocked(securityApi.getSecurityEvents).mockResolvedValue({
      total: 1,
      events: [
        {
          id: 9,
          eventType: 'login_failed',
          severity: 'high',
          username: 'admin',
          userId: 101,
          clientIp: '127.0.0.1',
          userAgent: 'ua',
          description: 'failed login',
          details: '{}',
          handled: false,
          createdAt: '2026-02-17 12:11:00',
        },
      ],
    });
    vi.mocked(securityApi.revokeSession).mockResolvedValue(undefined);
    vi.mocked(securityApi.forceLogoutUser).mockResolvedValue(undefined);
    vi.mocked(securityApi.handleSecurityEvent).mockResolvedValue(undefined);
  });

  it('loads policy, sessions and events on mount', async () => {
    const wrapper = mount(SecuritySettingsView);
    await flushPromises();

    expect(securityApi.getPasswordPolicy).toHaveBeenCalledTimes(1);
    expect(securityApi.getUserSessions).toHaveBeenCalledWith(1, 20);
    expect(securityApi.getSecurityEvents).toHaveBeenCalledWith(1, 20);
    expect(wrapper.text()).toContain('活跃会话 (1)');
    expect(wrapper.text()).toContain('安全事件 (1)');
    expect(wrapper.text()).toContain('session-1');
  });

  it('shows error message when initial load fails', async () => {
    vi.mocked(securityApi.getPasswordPolicy).mockRejectedValueOnce(new Error('load failed'));
    const wrapper = mount(SecuritySettingsView);
    await flushPromises();

    expect(wrapper.text()).toContain('load failed');
  });

  it('saves policy and displays success message', async () => {
    const wrapper = mount(SecuritySettingsView);
    await flushPromises();

    const saveButton = findButtonByText(wrapper, '保存策略');
    expect(saveButton).toBeDefined();
    await saveButton!.trigger('click');
    await flushPromises();

    expect(securityApi.updatePasswordPolicy).toHaveBeenCalledTimes(1);
    expect(wrapper.text()).toContain('密码策略已更新');
  });

  it('revokeSession respects cancel confirmation', async () => {
    vi.stubGlobal('confirm', vi.fn(() => false));
    const wrapper = mount(SecuritySettingsView);
    await flushPromises();

    const revokeButton = findButtonByText(wrapper, '撤销会话');
    expect(revokeButton).toBeDefined();
    await revokeButton!.trigger('click');
    await flushPromises();

    expect(securityApi.revokeSession).not.toHaveBeenCalled();
  });

  it('force logout uses default reason when prompt is empty', async () => {
    vi.stubGlobal('prompt', vi.fn(() => ''));
    const wrapper = mount(SecuritySettingsView);
    await flushPromises();

    const forceButton = findButtonByText(wrapper, '强制下线');
    expect(forceButton).toBeDefined();
    await forceButton!.trigger('click');
    await flushPromises();

    expect(securityApi.forceLogoutUser).toHaveBeenCalledWith(101, '管理员操作');
    expect(wrapper.text()).toContain('用户已强制下线');
  });

  it('handles security event and shows success message', async () => {
    const wrapper = mount(SecuritySettingsView);
    await flushPromises();

    const handleButton = findButtonByText(wrapper, '标记已处理');
    expect(handleButton).toBeDefined();
    await handleButton!.trigger('click');
    await flushPromises();

    expect(securityApi.handleSecurityEvent).toHaveBeenCalledWith(9, 'ops note');
    expect(wrapper.text()).toContain('安全事件已处理');
  });
});
