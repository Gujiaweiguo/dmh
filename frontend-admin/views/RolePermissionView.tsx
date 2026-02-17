import { defineComponent, h, ref, onMounted, reactive, computed } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import { PermissionGuard, usePermission } from '../components/PermissionGuard';
import { roleApi } from '../services/roleApi';

// 角色权限管理视图
export const RolePermissionView = defineComponent({
  setup() {
    const { hasPermission } = usePermission();
    const loading = ref(false);
    const activeTab = ref<'roles' | 'permissions' | 'audit'>('roles');
    
    // 角色数据
    const roles = ref<any[]>([]);
    const permissions = ref<any[]>([]);
    const auditLogs = ref<any[]>([]);
    
    // 编辑角色权限弹窗
    const showRolePermissionDialog = ref(false);
    const editingRole = ref<any>(null);
    const selectedPermissions = ref<number[]>([]);
    
    // 权限分组
    const permissionGroups = computed(() => {
      const groups: Record<string, any[]> = {};
      permissions.value.forEach(permission => {
        const resource = permission.resource;
        if (!groups[resource]) {
          groups[resource] = [];
        }
        groups[resource].push(permission);
      });
      return groups;
    });

    // 加载角色列表
    const loadRoles = async () => {
      loading.value = true;
      try {
        roles.value = await roleApi.getRoles();
      } catch (error) {
        console.error('加载角色列表失败', error);
      } finally {
        loading.value = false;
      }
    };

    // 加载权限列表
    const loadPermissions = async () => {
      try {
        permissions.value = await roleApi.getPermissions();
      } catch (error) {
        console.error('加载权限列表失败', error);
      }
    };

    // 加载审计日志
    const loadAuditLogs = async () => {
      try {
        auditLogs.value = await roleApi.getAuditLogs();
      } catch (error) {
        console.error('加载审计日志失败', error);
        auditLogs.value = [];
      }
    };

    // 打开角色权限配置对话框
    const openRolePermissionDialog = (role: any) => {
      editingRole.value = role;
      selectedPermissions.value = role.permissions.includes('*') 
        ? permissions.value.map(p => p.id)
        : permissions.value.filter(p => role.permissions.includes(p.code)).map(p => p.id);
      showRolePermissionDialog.value = true;
    };

    // 保存角色权限
    const saveRolePermissions = async () => {
      try {
        const permissionCodes = selectedPermissions.value.map(id => 
          permissions.value.find(p => p.id === id)?.code
        ).filter(Boolean);
        
        await roleApi.configRolePermissions(editingRole.value.id, selectedPermissions.value);
        
        // 更新本地数据
        editingRole.value.permissions = permissionCodes;
        showRolePermissionDialog.value = false;
        editingRole.value = null;
        
        alert('权限配置保存成功');
      } catch (error) {
        console.error('保存角色权限失败', error);
        alert('保存失败，请重试');
      }
    };

    // 获取资源图标
    const getResourceIcon = (resource: string) => {
      const icons: Record<string, any> = {
        user: LucideIcons.Users,
        brand: LucideIcons.Building,
        campaign: LucideIcons.Target,
        order: LucideIcons.ShoppingCart,
        reward: LucideIcons.Gift,
        withdrawal: LucideIcons.Wallet,
      };
      return icons[resource] || LucideIcons.Shield;
    };

    // 获取资源名称
    const getResourceName = (resource: string) => {
      const names: Record<string, string> = {
        user: '用户管理',
        brand: '品牌管理',
        campaign: '活动管理',
        order: '订单管理',
        reward: '奖励管理',
        withdrawal: '提现管理',
      };
      return names[resource] || resource;
    };

    onMounted(() => {
      loadRoles();
      loadPermissions();
      loadAuditLogs();
    });

    return () => h(PermissionGuard, { permission: 'role:config' }, () => [
      h('div', { class: 'space-y-8 animate-in fade-in' }, [
        // 页面标题
        h('div', [
          h('h2', { class: 'text-4xl font-black text-slate-900' }, '角色权限管理'),
          h('p', { class: 'text-slate-400 mt-2' }, '配置系统角色和权限，管理用户访问控制')
        ]),

        // 标签页导航
        h('div', { class: 'bg-white rounded-3xl p-2 border border-slate-100' }, [
          h('div', { class: 'flex gap-2' }, [
            h('button', {
              onClick: () => activeTab.value = 'roles',
              class: `flex-1 py-3 px-6 rounded-2xl font-bold text-sm transition-all ${
                activeTab.value === 'roles'
                  ? 'bg-indigo-600 text-white shadow-lg'
                  : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
              }`
            }, [
              h(LucideIcons.Shield, { size: 18, class: 'inline mr-2' }),
              '角色管理'
            ]),
            h('button', {
              onClick: () => activeTab.value = 'permissions',
              class: `flex-1 py-3 px-6 rounded-2xl font-bold text-sm transition-all ${
                activeTab.value === 'permissions'
                  ? 'bg-indigo-600 text-white shadow-lg'
                  : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
              }`
            }, [
              h(LucideIcons.Key, { size: 18, class: 'inline mr-2' }),
              '权限列表'
            ]),
            h('button', {
              onClick: () => activeTab.value = 'audit',
              class: `flex-1 py-3 px-6 rounded-2xl font-bold text-sm transition-all ${
                activeTab.value === 'audit'
                  ? 'bg-indigo-600 text-white shadow-lg'
                  : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
              }`
            }, [
              h(LucideIcons.FileText, { size: 18, class: 'inline mr-2' }),
              '操作日志'
            ])
          ])
        ]),

        // 角色管理标签页
        activeTab.value === 'roles' && h('div', { class: 'space-y-6' }, [
          loading.value
            ? h('div', { class: 'text-center py-20 text-slate-400' }, '加载中...')
            : h('div', { class: 'grid grid-cols-1 md:grid-cols-2 gap-6' }, 
                roles.value.map(role => h('div', { class: 'bg-white rounded-3xl border border-slate-100 p-8 hover:shadow-lg transition-all' }, [
                  h('div', { class: 'flex items-start justify-between mb-6' }, [
                    h('div', { class: 'flex items-center gap-4' }, [
                      h('div', { class: `w-12 h-12 rounded-2xl flex items-center justify-center ${
                        role.code === 'platform_admin' ? 'bg-purple-100 text-purple-600' :
                        role.code === 'participant' ? 'bg-green-100 text-green-600' :
                        'bg-gray-100 text-gray-600'
                      }` }, [
                        h(LucideIcons.Shield, { size: 24 })
                      ]),
                      h('div', [
                        h('h3', { class: 'text-xl font-black text-slate-900' }, role.name),
                        h('p', { class: 'text-sm text-slate-500 mt-1' }, role.code)
                      ])
                    ]),
                    role.code !== 'platform_admin' && h('button', {
                      onClick: () => openRolePermissionDialog(role),
                      class: 'p-2 hover:bg-slate-100 rounded-xl transition-colors',
                      title: '配置权限'
                    }, h(LucideIcons.Settings, { size: 20, class: 'text-slate-600' }))
                  ]),
                  
                  h('p', { class: 'text-slate-600 mb-6' }, role.description),
                  
                  h('div', { class: 'space-y-3' }, [
                    h('div', { class: 'flex items-center justify-between' }, [
                      h('span', { class: 'text-sm font-bold text-slate-700' }, '权限数量'),
                      h('span', { class: 'text-sm text-slate-500' }, 
                        role.permissions.includes('*') ? '全部权限' : `${role.permissions.length} 个权限`
                      )
                    ]),
                    h('div', { class: 'flex items-center justify-between' }, [
                      h('span', { class: 'text-sm font-bold text-slate-700' }, '创建时间'),
                      h('span', { class: 'text-sm text-slate-500' }, role.createdAt)
                    ])
                  ])
                ]))
              )
        ]),

        // 权限列表标签页
        activeTab.value === 'permissions' && h('div', { class: 'space-y-6' }, [
          Object.entries(permissionGroups.value).map(([resource, perms]) => 
            h('div', { class: 'bg-white rounded-3xl border border-slate-100 overflow-hidden' }, [
              h('div', { class: 'bg-slate-50 px-6 py-4 border-b border-slate-100' }, [
                h('div', { class: 'flex items-center gap-3' }, [
                  h(getResourceIcon(resource), { size: 24, class: 'text-slate-600' }),
                  h('h3', { class: 'text-lg font-black text-slate-900' }, getResourceName(resource)),
                  h('span', { class: 'text-sm text-slate-500' }, `${perms.length} 个权限`)
                ])
              ]),
              h('div', { class: 'p-6' }, [
                h('div', { class: 'grid grid-cols-1 md:grid-cols-2 gap-4' }, 
                  perms.map(permission => 
                    h('div', { class: 'flex items-center justify-between p-4 rounded-2xl border border-slate-100 hover:bg-slate-50 transition-colors' }, [
                      h('div', [
                        h('div', { class: 'font-bold text-slate-900' }, permission.name),
                        h('div', { class: 'text-sm text-slate-500 mt-1' }, permission.code),
                        h('div', { class: 'text-xs text-slate-400 mt-1' }, permission.description)
                      ]),
                      h('div', { class: `px-3 py-1 rounded-full text-xs font-bold ${
                        permission.action === 'read' ? 'bg-blue-100 text-blue-700' :
                        permission.action === 'create' ? 'bg-green-100 text-green-700' :
                        permission.action === 'update' ? 'bg-amber-100 text-amber-700' :
                        permission.action === 'delete' ? 'bg-red-100 text-red-700' :
                        'bg-gray-100 text-gray-700'
                      }` }, permission.action.toUpperCase())
                    ])
                  )
                )
              ])
            ])
          )
        ]),

        // 操作日志标签页
        activeTab.value === 'audit' && h('div', { class: 'bg-white rounded-3xl border border-slate-100 overflow-hidden' }, [
          h('div', { class: 'px-6 py-4 border-b border-slate-100 bg-slate-50' }, [
            h('h3', { class: 'text-lg font-black text-slate-900' }, '权限操作日志'),
            h('p', { class: 'text-sm text-slate-500 mt-1' }, '记录所有权限相关的操作和变更')
          ]),
          auditLogs.value.length === 0
            ? h('div', { class: 'text-center py-20 text-slate-400' }, '暂无操作日志')
            : h('div', { class: 'divide-y divide-slate-100' }, 
                auditLogs.value.map(log => 
                  h('div', { class: 'p-6 hover:bg-slate-50 transition-colors' }, [
                    h('div', { class: 'flex items-start justify-between' }, [
                      h('div', { class: 'flex items-start gap-4' }, [
                        h('div', { class: 'w-10 h-10 rounded-full bg-indigo-100 text-indigo-600 flex items-center justify-center flex-shrink-0' }, [
                          h(LucideIcons.Activity, { size: 20 })
                        ]),
                        h('div', [
                          h('div', { class: 'font-bold text-slate-900' }, log.action),
                          h('div', { class: 'text-sm text-slate-600 mt-1' }, [
                            '操作者：',
                            h('span', { class: 'font-medium' }, log.username),
                            ' | 目标：',
                            h('span', { class: 'font-medium' }, log.target)
                          ]),
                          h('div', { class: 'text-sm text-slate-500 mt-2' }, log.details),
                          h('div', { class: 'text-xs text-slate-400 mt-2' }, [
                            `IP: ${log.ip} | 时间: ${log.createdAt}`
                          ])
                        ])
                      ]),
                      h('button', {
                        class: 'p-2 hover:bg-slate-100 rounded-xl transition-colors',
                        title: '查看详情'
                      }, h(LucideIcons.Eye, { size: 16, class: 'text-slate-400' }))
                    ])
                  ])
                )
              )
        ]),

        // 角色权限配置对话框
        showRolePermissionDialog.value && editingRole.value && h('div', { 
          class: 'fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4',
          onClick: () => showRolePermissionDialog.value = false
        }, [
          h('div', { 
            class: 'bg-white rounded-3xl p-8 max-w-4xl w-full max-h-[80vh] overflow-auto',
            onClick: (e: Event) => e.stopPropagation()
          }, [
            h('div', { class: 'flex items-center justify-between mb-6' }, [
              h('div', [
                h('h3', { class: 'text-2xl font-black text-slate-900' }, `配置"${editingRole.value.name}"权限`),
                h('p', { class: 'text-slate-500 mt-1' }, '选择该角色可以访问的功能权限')
              ]),
              h('button', {
                onClick: () => showRolePermissionDialog.value = false,
                class: 'p-2 hover:bg-slate-100 rounded-xl transition-colors'
              }, h(LucideIcons.X, { size: 20 }))
            ]),
            
            h('div', { class: 'space-y-6' }, 
              Object.entries(permissionGroups.value).map(([resource, perms]) => 
                h('div', { class: 'border border-slate-200 rounded-2xl overflow-hidden' }, [
                  h('div', { class: 'bg-slate-50 px-4 py-3 border-b border-slate-200' }, [
                    h('div', { class: 'flex items-center gap-3' }, [
                      h(getResourceIcon(resource), { size: 20, class: 'text-slate-600' }),
                      h('span', { class: 'font-bold text-slate-900' }, getResourceName(resource))
                    ])
                  ]),
                  h('div', { class: 'p-4 space-y-3' }, 
                    perms.map(permission => 
                      h('label', { 
                        class: 'flex items-center gap-3 p-3 rounded-xl hover:bg-slate-50 cursor-pointer transition-colors'
                      }, [
                        h('input', {
                          type: 'checkbox',
                          checked: selectedPermissions.value.includes(permission.id),
                          onChange: (e: any) => {
                            if (e.target.checked) {
                              selectedPermissions.value.push(permission.id);
                            } else {
                              selectedPermissions.value = selectedPermissions.value.filter(id => id !== permission.id);
                            }
                          },
                          class: 'w-5 h-5 rounded border-slate-300'
                        }),
                        h('div', { class: 'flex-1' }, [
                          h('div', { class: 'font-medium text-slate-900' }, permission.name),
                          h('div', { class: 'text-sm text-slate-500' }, permission.description)
                        ]),
                        h('div', { class: `px-2 py-1 rounded text-xs font-bold ${
                          permission.action === 'read' ? 'bg-blue-100 text-blue-700' :
                          permission.action === 'create' ? 'bg-green-100 text-green-700' :
                          permission.action === 'update' ? 'bg-amber-100 text-amber-700' :
                          permission.action === 'delete' ? 'bg-red-100 text-red-700' :
                          'bg-gray-100 text-gray-700'
                        }` }, permission.action)
                      ])
                    )
                  )
                ])
              )
            ),
            
            h('div', { class: 'flex gap-3 mt-8' }, [
              h('button', {
                onClick: () => showRolePermissionDialog.value = false,
                class: 'flex-1 px-6 py-3 rounded-xl border border-slate-200 font-bold hover:bg-slate-50 transition-colors'
              }, '取消'),
              h('button', {
                onClick: saveRolePermissions,
                class: 'flex-1 px-6 py-3 rounded-xl bg-indigo-600 text-white font-bold hover:bg-indigo-500 transition-colors'
              }, '保存权限配置')
            ])
          ])
        ])
      ])
    ]);
  }
});
