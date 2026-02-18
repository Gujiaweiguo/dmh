import { defineComponent, h, onMounted, reactive, ref } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import { securityApi, type PasswordPolicy, type SecurityEvent, type UserSession } from '../services/securityApi';

const defaultPolicy: PasswordPolicy = {
  id: 0,
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

const formatTime = (value?: string) => {
  if (!value) {
    return '-';
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleString();
};

export const SecurityManagementView = defineComponent({
  name: 'SecurityManagementView',
  setup() {
    const loading = ref(false);
    const savingPolicy = ref(false);
    const actionKey = ref('');
    const messageType = ref<'success' | 'error' | ''>('');
    const messageText = ref('');

    const policy = reactive<PasswordPolicy>({ ...defaultPolicy });
    const sessions = ref<UserSession[]>([]);
    const sessionTotal = ref(0);
    const events = ref<SecurityEvent[]>([]);
    const eventTotal = ref(0);

    const setMessage = (type: 'success' | 'error', text: string) => {
      messageType.value = type;
      messageText.value = text;
    };

    const loadPolicy = async () => {
      const policyResp = await securityApi.getPasswordPolicy();
      Object.assign(policy, defaultPolicy, policyResp);
    };

    const loadSessions = async () => {
      const resp = await securityApi.getUserSessions(1, 20);
      sessions.value = resp.sessions || [];
      sessionTotal.value = resp.total || 0;
    };

    const loadEvents = async () => {
      const resp = await securityApi.getSecurityEvents(1, 20);
      events.value = resp.events || [];
      eventTotal.value = resp.total || 0;
    };

    const loadAll = async () => {
      loading.value = true;
      try {
        await Promise.all([loadPolicy(), loadSessions(), loadEvents()]);
      } catch (error) {
        const text = error instanceof Error ? error.message : '加载安全数据失败';
        setMessage('error', text);
      } finally {
        loading.value = false;
      }
    };

    const savePolicy = async () => {
      savingPolicy.value = true;
      try {
        const resp = await securityApi.updatePasswordPolicy({ ...policy });
        Object.assign(policy, defaultPolicy, resp);
        setMessage('success', '密码策略已更新');
      } catch (error) {
        const text = error instanceof Error ? error.message : '更新密码策略失败';
        setMessage('error', text);
      } finally {
        savingPolicy.value = false;
      }
    };

    const revokeSession = async (sessionId: string) => {
      if (!confirm(`确认撤销会话 ${sessionId} 吗？`)) {
        return;
      }
      actionKey.value = `revoke:${sessionId}`;
      try {
        await securityApi.revokeSession(sessionId);
        await loadSessions();
        setMessage('success', '会话已撤销');
      } catch (error) {
        const text = error instanceof Error ? error.message : '撤销会话失败';
        setMessage('error', text);
      } finally {
        actionKey.value = '';
      }
    };

    const forceLogout = async (userId: number) => {
      if (!confirm(`确认强制用户 ${userId} 下线吗？`)) {
        return;
      }
      const reason = window.prompt('请输入强制下线原因（可选）') || '管理员操作';
      actionKey.value = `force:${userId}`;
      try {
        await securityApi.forceLogoutUser(userId, reason);
        await loadSessions();
        setMessage('success', '用户已强制下线');
      } catch (error) {
        const text = error instanceof Error ? error.message : '强制下线失败';
        setMessage('error', text);
      } finally {
        actionKey.value = '';
      }
    };

    const handleEvent = async (eventId: number) => {
      const note = window.prompt('输入处理备注（可选）') || '';
      actionKey.value = `event:${eventId}`;
      try {
        await securityApi.handleSecurityEvent(eventId, note);
        await loadEvents();
        setMessage('success', '安全事件已处理');
      } catch (error) {
        const text = error instanceof Error ? error.message : '处理安全事件失败';
        setMessage('error', text);
      } finally {
        actionKey.value = '';
      }
    };

    onMounted(() => {
      void loadAll();
    });

    const renderNumberInput = (label: string, key: keyof PasswordPolicy, min: number, max: number) =>
      h('label', { class: 'flex flex-col gap-2' }, [
        h('span', { class: 'text-xs font-bold text-slate-500 uppercase tracking-wider' }, label),
        h('input', {
          type: 'number',
          min,
          max,
          value: String(policy[key]),
          onInput: (event: Event) => {
            const target = event.target as HTMLInputElement | null;
            if (!target) {
              return;
            }
            const nextValue = Number.parseInt(target.value, 10);
            if (Number.isNaN(nextValue)) {
              return;
            }
            (policy[key] as number) = nextValue;
          },
          class: 'rounded-xl border border-slate-200 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none',
        }),
      ]);

    const renderSwitch = (label: string, key: keyof PasswordPolicy) =>
      h('label', { class: 'flex items-center gap-3 rounded-xl border border-slate-100 px-3 py-2' }, [
        h('input', {
          type: 'checkbox',
          checked: Boolean(policy[key]),
          onChange: (event: Event) => {
            const target = event.target as HTMLInputElement | null;
            if (!target) {
              return;
            }
            (policy[key] as boolean) = target.checked;
          },
          class: 'h-4 w-4 rounded border-slate-300',
        }),
        h('span', { class: 'text-sm font-medium text-slate-700' }, label),
      ]);

    return () =>
      h('div', { class: 'space-y-6' }, [
        h('div', { class: 'flex items-center justify-between' }, [
          h('div', [
            h('h2', { class: 'text-2xl font-black text-slate-900' }, '安全管理'),
            h('p', { class: 'mt-1 text-sm text-slate-500' }, '管理密码策略、会话与安全事件'),
          ]),
          h(
            'button',
            {
              class: 'inline-flex items-center gap-2 rounded-xl bg-indigo-600 px-4 py-2 text-sm font-bold text-white hover:bg-indigo-700 disabled:cursor-not-allowed disabled:opacity-60',
              disabled: loading.value,
              onClick: () => {
                void loadAll();
              },
            },
            [h(LucideIcons.RefreshCw, { size: 16 }), loading.value ? '刷新中...' : '刷新数据'],
          ),
        ]),

        messageType.value
          ? h(
              'div',
              {
                class:
                  messageType.value === 'success'
                    ? 'rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm font-medium text-emerald-700'
                    : 'rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-medium text-rose-700',
              },
              messageText.value,
            )
          : null,

        h('div', { class: 'rounded-3xl border border-slate-100 bg-white p-6 shadow-sm' }, [
          h('div', { class: 'mb-4 flex items-center justify-between' }, [
            h('h3', { class: 'text-lg font-black text-slate-900' }, '密码策略'),
            h(
              'button',
              {
                class: 'rounded-xl bg-emerald-600 px-4 py-2 text-sm font-bold text-white hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-60',
                disabled: savingPolicy.value,
                onClick: () => {
                  void savePolicy();
                },
              },
              savingPolicy.value ? '保存中...' : '保存策略',
            ),
          ]),
          h('div', { class: 'grid grid-cols-1 gap-4 md:grid-cols-3' }, [
            renderNumberInput('最小长度', 'minLength', 6, 50),
            renderNumberInput('密码有效期(天)', 'maxAge', 0, 365),
            renderNumberInput('历史密码数', 'historyCount', 0, 20),
            renderNumberInput('最大登录尝试', 'maxLoginAttempts', 1, 20),
            renderNumberInput('锁定时长(分钟)', 'lockoutDuration', 1, 1440),
            renderNumberInput('会话超时(分钟)', 'sessionTimeout', 5, 1440),
            renderNumberInput('最大并发会话', 'maxConcurrentSessions', 1, 10),
          ]),
          h('div', { class: 'mt-4 grid grid-cols-1 gap-3 md:grid-cols-2' }, [
            renderSwitch('需要大写字母', 'requireUppercase'),
            renderSwitch('需要小写字母', 'requireLowercase'),
            renderSwitch('需要数字', 'requireNumbers'),
            renderSwitch('需要特殊字符', 'requireSpecialChars'),
          ]),
        ]),

        h('div', { class: 'rounded-3xl border border-slate-100 bg-white p-6 shadow-sm' }, [
          h('div', { class: 'mb-4 flex items-center justify-between' }, [
            h('h3', { class: 'text-lg font-black text-slate-900' }, `活跃会话 (${sessionTotal.value})`),
            h('span', { class: 'text-xs text-slate-500' }, '仅展示最近 20 条'),
          ]),
          h('div', { class: 'overflow-x-auto' }, [
            h('table', { class: 'min-w-full divide-y divide-slate-200 text-sm' }, [
              h('thead', [
                h('tr', { class: 'text-left text-xs uppercase tracking-wider text-slate-500' }, [
                  h('th', { class: 'py-2 pr-4' }, '用户ID'),
                  h('th', { class: 'py-2 pr-4' }, '会话ID'),
                  h('th', { class: 'py-2 pr-4' }, 'IP'),
                  h('th', { class: 'py-2 pr-4' }, '状态'),
                  h('th', { class: 'py-2 pr-4' }, '最后活跃'),
                  h('th', { class: 'py-2' }, '操作'),
                ]),
              ]),
              h(
                'tbody',
                { class: 'divide-y divide-slate-100 text-slate-700' },
                sessions.value.length > 0
                  ? sessions.value.map((session) =>
                      h('tr', { key: session.id }, [
                        h('td', { class: 'py-3 pr-4 font-semibold' }, String(session.userId)),
                        h('td', { class: 'py-3 pr-4 font-mono text-xs' }, session.id),
                        h('td', { class: 'py-3 pr-4' }, session.clientIp || '-'),
                        h('td', { class: 'py-3 pr-4' }, session.status),
                        h('td', { class: 'py-3 pr-4 text-xs text-slate-500' }, formatTime(session.lastActiveAt)),
                        h('td', { class: 'py-3' }, [
                          h(
                            'button',
                            {
                              class: 'mr-2 rounded-lg border border-amber-200 px-3 py-1 text-xs font-bold text-amber-700 hover:bg-amber-50 disabled:cursor-not-allowed disabled:opacity-60',
                              disabled: actionKey.value === `revoke:${session.id}`,
                              onClick: () => {
                                void revokeSession(session.id);
                              },
                            },
                            '撤销会话',
                          ),
                          h(
                            'button',
                            {
                              class: 'rounded-lg border border-rose-200 px-3 py-1 text-xs font-bold text-rose-700 hover:bg-rose-50 disabled:cursor-not-allowed disabled:opacity-60',
                              disabled: actionKey.value === `force:${session.userId}`,
                              onClick: () => {
                                void forceLogout(session.userId);
                              },
                            },
                            '强制下线',
                          ),
                        ]),
                      ]),
                    )
                  : [
                      h('tr', { key: 'empty' }, [
                        h('td', { class: 'py-6 text-center text-slate-400', colSpan: 6 }, '暂无会话数据'),
                      ]),
                    ],
              ),
            ]),
          ]),
        ]),

        h('div', { class: 'rounded-3xl border border-slate-100 bg-white p-6 shadow-sm' }, [
          h('div', { class: 'mb-4 flex items-center justify-between' }, [
            h('h3', { class: 'text-lg font-black text-slate-900' }, `安全事件 (${eventTotal.value})`),
            h('span', { class: 'text-xs text-slate-500' }, '仅展示最近 20 条'),
          ]),
          h('div', { class: 'overflow-x-auto' }, [
            h('table', { class: 'min-w-full divide-y divide-slate-200 text-sm' }, [
              h('thead', [
                h('tr', { class: 'text-left text-xs uppercase tracking-wider text-slate-500' }, [
                  h('th', { class: 'py-2 pr-4' }, '时间'),
                  h('th', { class: 'py-2 pr-4' }, '类型'),
                  h('th', { class: 'py-2 pr-4' }, '严重级别'),
                  h('th', { class: 'py-2 pr-4' }, '描述'),
                  h('th', { class: 'py-2' }, '操作'),
                ]),
              ]),
              h(
                'tbody',
                { class: 'divide-y divide-slate-100 text-slate-700' },
                events.value.length > 0
                  ? events.value.map((event) =>
                      h('tr', { key: event.id }, [
                        h('td', { class: 'py-3 pr-4 text-xs text-slate-500' }, formatTime(event.createdAt)),
                        h('td', { class: 'py-3 pr-4 font-semibold' }, event.eventType),
                        h('td', { class: 'py-3 pr-4' }, event.severity),
                        h('td', { class: 'py-3 pr-4' }, event.description),
                        h(
                          'td',
                          { class: 'py-3' },
                          event.handled
                            ? h('span', { class: 'text-xs font-bold text-emerald-600' }, '已处理')
                            : h(
                                'button',
                                {
                                  class: 'rounded-lg border border-indigo-200 px-3 py-1 text-xs font-bold text-indigo-700 hover:bg-indigo-50 disabled:cursor-not-allowed disabled:opacity-60',
                                  disabled: actionKey.value === `event:${event.id}`,
                                  onClick: () => {
                                    void handleEvent(event.id);
                                  },
                                },
                                '标记已处理',
                              ),
                        ),
                      ]),
                    )
                  : [
                      h('tr', { key: 'empty' }, [
                        h('td', { class: 'py-6 text-center text-slate-400', colSpan: 5 }, '暂无安全事件'),
                      ]),
                    ],
              ),
            ]),
          ]),
        ]),
      ]);
  },
});

export default SecurityManagementView;
