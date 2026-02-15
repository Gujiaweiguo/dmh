import { defineComponent, h, ref, onMounted, reactive, computed } from 'vue';
import * as LucideIcons from 'lucide-vue-next';

// 用户管理视图
export const UserManagementView = defineComponent({
  setup() {
    const users = ref<any[]>([]);
    const loading = ref(false);
    const searchQuery = ref('');
    const filterRole = ref('all');
    const filterStatus = ref('all');
    
    // 编辑用户弹窗
    const showEditDialog = ref(false);
    const editingUser = ref<any>(null);
    
    // 重置密码弹窗
    const showResetPasswordDialog = ref(false);
    const resetPasswordUser = ref<any>(null);
    const newPassword = ref('');
    
    // 角色管理弹窗
    const showRoleDialog = ref(false);
    const roleUser = ref<any>(null);
    const availableRoles = ref([
      { code: 'platform_admin', name: '平台管理员', description: '拥有系统最高权限' },
      { code: 'participant', name: '活动参与者', description: '可参与活动、分享和提现' },
    ]);
    const selectedRoles = ref<string[]>([]);
    
    const availableBrands = ref<any[]>([]);

    // 加载用户列表
    const loadUsers = async () => {
      loading.value = true;
      try {
        // TODO: 调用真实API
        // const response = await fetch('/api/v1/users', {
        //   headers: { 'Authorization': `Bearer ${localStorage.getItem('dmh_token')}` }
        // });
        // users.value = await response.json();
        
        // 模拟数据
        users.value = [
          {
            id: 1,
            username: 'admin',
            realName: '系统管理员',
            phone: '13800000001',
            email: 'admin@dmh.com',
            status: 'active',
            roles: ['platform_admin'],
            brandIds: [],
            createdAt: '2025-01-01 10:00:00',
          },
          {
            id: 3,
            username: 'user001',
            realName: '张三',
            phone: '13800000003',
            email: 'user001@dmh.com',
            status: 'active',
            roles: ['participant'],
            brandIds: [],
            createdAt: '2025-01-03 10:00:00',
          },
          {
            id: 4,
            username: 'user002',
            realName: '李四',
            phone: '13800000004',
            email: 'user002@dmh.com',
            status: 'disabled',
            roles: ['participant'],
            brandIds: [],
            createdAt: '2025-01-04 10:00:00',
          },
        ];
        
        availableBrands.value = [
          { id: 1, name: '品牌A' },
          { id: 2, name: '品牌B' },
          { id: 3, name: '品牌C' },
        ];
      } catch (error) {
        console.error('加载用户列表失败', error);
      } finally {
        loading.value = false;
      }
    };

    // 筛选用户
    const filteredUsers = computed(() => {
      return users.value.filter(user => {
        // 搜索过滤
        if (searchQuery.value) {
          const query = searchQuery.value.toLowerCase();
          const matchSearch = 
            user.username.toLowerCase().includes(query) ||
            user.realName?.toLowerCase().includes(query) ||
            user.phone.includes(query) ||
            user.email?.toLowerCase().includes(query);
          if (!matchSearch) return false;
        }
        
        // 角色过滤
        if (filterRole.value !== 'all') {
          if (!user.roles.includes(filterRole.value)) return false;
        }
        
        // 状态过滤
        if (filterStatus.value !== 'all') {
          if (user.status !== filterStatus.value) return false;
        }
        
        return true;
      });
    });

    // 打开编辑用户对话框
    const openEditDialog = (user: any) => {
      editingUser.value = { ...user };
      showEditDialog.value = true;
    };

    // 保存用户信息
    const saveUser = async () => {
      try {
        // TODO: 调用真实API
        // await fetch(`/api/v1/users/${editingUser.value.id}`, {
        //   method: 'PUT',
        //   headers: {
        //     'Authorization': `Bearer ${localStorage.getItem('dmh_token')}`,
        //     'Content-Type': 'application/json',
        //   },
        //   body: JSON.stringify(editingUser.value),
        // });
        
        const index = users.value.findIndex(u => u.id === editingUser.value.id);
        if (index !== -1) {
          users.value[index] = { ...editingUser.value };
        }
        showEditDialog.value = false;
        editingUser.value = null;
      } catch (error) {
        console.error('保存用户失败', error);
        alert('保存失败，请重试');
      }
    };

    // 切换用户状态
    const toggleUserStatus = async (user: any) => {
      const newStatus = user.status === 'active' ? 'disabled' : 'active';
      const confirmText = newStatus === 'disabled' ? '禁用' : '启用';
      
      if (!confirm(`确定要${confirmText}用户"${user.username}"吗？`)) {
        return;
      }
      
      try {
        // TODO: 调用真实API
        // await fetch(`/api/v1/users/${user.id}/status`, {
        //   method: 'PUT',
        //   headers: {
        //     'Authorization': `Bearer ${localStorage.getItem('dmh_token')}`,
        //     'Content-Type': 'application/json',
        //   },
        //   body: JSON.stringify({ status: newStatus }),
        // });
        
        user.status = newStatus;
      } catch (error) {
        console.error('切换用户状态失败', error);
        alert('操作失败，请重试');
      }
    };

    // 打开重置密码对话框
    const openResetPasswordDialog = (user: any) => {
      resetPasswordUser.value = user;
      newPassword.value = '';
      showResetPasswordDialog.value = true;
    };

    // 重置密码
    const resetPassword = async () => {
      if (!newPassword.value || newPassword.value.length < 6) {
        alert('密码长度至少6位');
        return;
      }
      
      try {
        // TODO: 调用真实API
        // await fetch(`/api/v1/users/${resetPasswordUser.value.id}/reset-password`, {
        //   method: 'POST',
        //   headers: {
        //     'Authorization': `Bearer ${localStorage.getItem('dmh_token')}`,
        //     'Content-Type': 'application/json',
        //   },
        //   body: JSON.stringify({ newPassword: newPassword.value }),
        // });
        
        alert('密码重置成功');
        showResetPasswordDialog.value = false;
        resetPasswordUser.value = null;
        newPassword.value = '';
      } catch (error) {
        console.error('重置密码失败', error);
        alert('重置失败，请重试');
      }
    };

    // 打开角色管理对话框
    const openRoleDialog = (user: any) => {
      roleUser.value = user;
      selectedRoles.value = [...user.roles];
      selectedBrands.value = [...(user.brandIds || [])];
      showRoleDialog.value = true;
    };

    // 保存角色
    const saveRoles = async () => {
      try {
        // TODO: 调用真实API
        // await fetch(`/api/v1/users/${roleUser.value.id}/roles`, {
        //   method: 'PUT',
        //   headers: {
        //     'Authorization': `Bearer ${localStorage.getItem('dmh_token')}`,
        //     'Content-Type': 'application/json',
        //   },
        //   body: JSON.stringify({ 
        //     roles: selectedRoles.value,
        //     brandIds: selectedBrands.value,
        //   }),
        // });
        
        roleUser.value.roles = [...selectedRoles.value];
        roleUser.value.brandIds = [...selectedBrands.value];
        showRoleDialog.value = false;
        roleUser.value = null;
      } catch (error) {
        console.error('保存角色失败', error);
        alert('保存失败，请重试');
      }
    };

    // 删除用户
    const deleteUser = async (user: any) => {
      if (!confirm(`确定要删除用户"${user.username}"吗？此操作不可恢复！`)) {
        return;
      }
      
      try {
        // TODO: 调用真实API
        // await fetch(`/api/v1/users/${user.id}`, {
        //   method: 'DELETE',
        //   headers: {
        //     'Authorization': `Bearer ${localStorage.getItem('dmh_token')}`,
        //   },
        // });
        
        users.value = users.value.filter(u => u.id !== user.id);
      } catch (error) {
        console.error('删除用户失败', error);
        alert('删除失败，请重试');
      }
    };

    // 角色徽章颜色
    const getRoleBadgeColor = (roleCode: string) => {
      const colors: Record<string, string> = {
        platform_admin: 'bg-purple-100 text-purple-700 border-purple-200',
        brand_admin: 'bg-blue-100 text-blue-700 border-blue-200',
        participant: 'bg-green-100 text-green-700 border-green-200',
      };
      return colors[roleCode] || 'bg-gray-100 text-gray-700 border-gray-200';
    };

    // 角色名称映射
    const getRoleName = (roleCode: string) => {
      const names: Record<string, string> = {
        platform_admin: '平台管理员',
        participant: '活动参与者',
      };
      return names[roleCode] || roleCode;
    };

    onMounted(() => {
      loadUsers();
    });

    return () => h('div', { class: 'space-y-8 animate-in fade-in' }, [
      // 页面标题
      h('div', { class: 'flex justify-between items-end' }, [
        h('div', [
          h('h2', { class: 'text-4xl font-black text-slate-900' }, '用户管理'),
          h('p', { class: 'text-slate-400 mt-2' }, '管理系统用户、角色分配、权限控制')
        ]),
        h('button', {
          onClick: loadUsers,
          class: 'bg-slate-900 text-white px-6 py-3 rounded-2xl font-bold shadow-lg flex items-center gap-2 hover:bg-slate-800 transition-colors'
        }, [
          h(LucideIcons.RefreshCw, { size: 18 }),
          '刷新'
        ])
      ]),

      // 筛选和搜索栏
      h('div', { class: 'bg-white rounded-3xl p-6 border border-slate-100' }, [
        h('div', { class: 'grid grid-cols-1 md:grid-cols-4 gap-4' }, [
          // 搜索框
          h('div', { class: 'md:col-span-2' }, [
            h('div', { class: 'relative' }, [
              h('div', { class: 'absolute left-4 top-1/2 -translate-y-1/2 text-slate-400' }, 
                h(LucideIcons.Search, { size: 20 })
              ),
              h('input', {
                type: 'text',
                value: searchQuery.value,
                onInput: (e: any) => { searchQuery.value = e.target.value },
                placeholder: '搜索用户名、姓名、手机号、邮箱...',
                class: 'w-full pl-12 pr-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
              })
            ])
          ]),
          
          // 角色筛选
          h('select', {
            value: filterRole.value,
            onChange: (e: any) => { filterRole.value = e.target.value },
            class: 'px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
          }, [
            h('option', { value: 'all' }, '全部角色'),
            h('option', { value: 'platform_admin' }, '平台管理员'),
            h('option', { value: 'participant' }, '活动参与者'),
          ]),
          
          // 状态筛选
          h('select', {
            value: filterStatus.value,
            onChange: (e: any) => { filterStatus.value = e.target.value },
            class: 'px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
          }, [
            h('option', { value: 'all' }, '全部状态'),
            h('option', { value: 'active' }, '正常'),
            h('option', { value: 'disabled' }, '已禁用'),
          ]),
        ])
      ]),

      // 用户列表
      loading.value
        ? h('div', { class: 'text-center py-20 text-slate-400' }, '加载中...')
        : filteredUsers.value.length === 0
          ? h('div', { class: 'text-center py-20 text-slate-400' }, '暂无用户数据')
          : h('div', { class: 'bg-white rounded-3xl border border-slate-100 overflow-hidden' }, [
              h('table', { class: 'w-full' }, [
                // 表头
                h('thead', { class: 'bg-slate-50 border-b border-slate-100' }, [
                  h('tr', {}, [
                    h('th', { class: 'px-6 py-4 text-left text-xs font-black text-slate-500 uppercase tracking-wider' }, 'ID'),
                    h('th', { class: 'px-6 py-4 text-left text-xs font-black text-slate-500 uppercase tracking-wider' }, '用户名'),
                    h('th', { class: 'px-6 py-4 text-left text-xs font-black text-slate-500 uppercase tracking-wider' }, '姓名'),
                    h('th', { class: 'px-6 py-4 text-left text-xs font-black text-slate-500 uppercase tracking-wider' }, '手机号'),
                    h('th', { class: 'px-6 py-4 text-left text-xs font-black text-slate-500 uppercase tracking-wider' }, '邮箱'),
                    h('th', { class: 'px-6 py-4 text-left text-xs font-black text-slate-500 uppercase tracking-wider' }, '角色'),
                    h('th', { class: 'px-6 py-4 text-left text-xs font-black text-slate-500 uppercase tracking-wider' }, '状态'),
                    h('th', { class: 'px-6 py-4 text-left text-xs font-black text-slate-500 uppercase tracking-wider' }, '注册时间'),
                    h('th', { class: 'px-6 py-4 text-right text-xs font-black text-slate-500 uppercase tracking-wider' }, '操作'),
                  ])
                ]),
                
                // 表体
                h('tbody', { class: 'divide-y divide-slate-100' }, 
                  filteredUsers.value.map(user => h('tr', { class: 'hover:bg-slate-50 transition-colors' }, [
                    h('td', { class: 'px-6 py-4 text-sm font-medium text-slate-900' }, user.id),
                    h('td', { class: 'px-6 py-4' }, [
                      h('div', { class: 'flex items-center gap-3' }, [
                        h('div', { class: 'w-10 h-10 rounded-full bg-indigo-100 text-indigo-600 flex items-center justify-center font-bold' }, 
                          user.username.charAt(0).toUpperCase()
                        ),
                        h('span', { class: 'font-medium text-slate-900' }, user.username)
                      ])
                    ]),
                    h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, user.realName || '-'),
                    h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, user.phone),
                    h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, user.email || '-'),
                    h('td', { class: 'px-6 py-4' }, [
                      h('div', { class: 'flex flex-wrap gap-1' }, 
                        user.roles.map((role: string) => 
                          h('span', { 
                            class: `px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider border ${getRoleBadgeColor(role)}` 
                          }, getRoleName(role))
                        )
                      )
                    ]),
                    h('td', { class: 'px-6 py-4' }, [
                      user.status === 'active'
                        ? h('span', { class: 'px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider border bg-green-100 text-green-700 border-green-200' }, '正常')
                        : h('span', { class: 'px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider border bg-red-100 text-red-700 border-red-200' }, '已禁用')
                    ]),
                    h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, user.createdAt),
                    h('td', { class: 'px-6 py-4 text-right' }, [
                      h('div', { class: 'flex justify-end gap-2' }, [
                        // 编辑按钮
                        h('button', {
                          onClick: () => openEditDialog(user),
                          class: 'p-2 hover:bg-slate-100 rounded-lg transition-colors',
                          title: '编辑用户'
                        }, h(LucideIcons.Edit, { size: 16, class: 'text-blue-600' })),
                        
                        // 角色管理按钮
                        h('button', {
                          onClick: () => openRoleDialog(user),
                          class: 'p-2 hover:bg-slate-100 rounded-lg transition-colors',
                          title: '管理角色'
                        }, h(LucideIcons.Shield, { size: 16, class: 'text-purple-600' })),
                        
                        // 重置密码按钮
                        h('button', {
                          onClick: () => openResetPasswordDialog(user),
                          class: 'p-2 hover:bg-slate-100 rounded-lg transition-colors',
                          title: '重置密码'
                        }, h(LucideIcons.Key, { size: 16, class: 'text-amber-600' })),
                        
                        // 切换状态按钮
                        h('button', {
                          onClick: () => toggleUserStatus(user),
                          class: 'p-2 hover:bg-slate-100 rounded-lg transition-colors',
                          title: user.status === 'active' ? '禁用用户' : '启用用户'
                        }, user.status === 'active'
                          ? h(LucideIcons.UserX, { size: 16, class: 'text-orange-600' })
                          : h(LucideIcons.UserCheck, { size: 16, class: 'text-green-600' })
                        ),
                        
                        // 删除按钮
                        h('button', {
                          onClick: () => deleteUser(user),
                          class: 'p-2 hover:bg-slate-100 rounded-lg transition-colors',
                          title: '删除用户'
                        }, h(LucideIcons.Trash2, { size: 16, class: 'text-red-600' })),
                      ])
                    ]),
                  ]))
                )
              ])
            ]),

      // 编辑用户对话框
      showEditDialog.value && editingUser.value && h('div', { 
        class: 'fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4',
        onClick: () => { showEditDialog.value = false }
      }, [
        h('div', { 
          class: 'bg-white rounded-3xl p-8 max-w-2xl w-full max-h-[80vh] overflow-auto',
          onClick: (e: Event) => e.stopPropagation()
        }, [
          h('div', { class: 'flex items-center justify-between mb-6' }, [
            h('h3', { class: 'text-2xl font-black text-slate-900' }, '编辑用户'),
            h('button', {
              onClick: () => { showEditDialog.value = false },
              class: 'p-2 hover:bg-slate-100 rounded-xl transition-colors'
            }, h(LucideIcons.X, { size: 20 }))
          ]),
          
          h('div', { class: 'space-y-4' }, [
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '用户名'),
              h('input', {
                type: 'text',
                value: editingUser.value.username,
                onInput: (e: any) => { editingUser.value.username = e.target.value },
                class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500',
                disabled: true // 用户名不可修改
              })
            ]),
            
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '真实姓名'),
              h('input', {
                type: 'text',
                value: editingUser.value.realName,
                onInput: (e: any) => { editingUser.value.realName = e.target.value },
                class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500'
              })
            ]),
            
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '手机号'),
              h('input', {
                type: 'tel',
                value: editingUser.value.phone,
                onInput: (e: any) => { editingUser.value.phone = e.target.value },
                class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500'
              })
            ]),
            
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '邮箱'),
              h('input', {
                type: 'email',
                value: editingUser.value.email,
                onInput: (e: any) => { editingUser.value.email = e.target.value },
                class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500'
              })
            ]),
          ]),
          
          h('div', { class: 'flex gap-3 mt-8' }, [
            h('button', {
              onClick: () => { showEditDialog.value = false },
              class: 'flex-1 px-6 py-3 rounded-xl border border-slate-200 font-bold hover:bg-slate-50 transition-colors'
            }, '取消'),
            h('button', {
              onClick: saveUser,
              class: 'flex-1 px-6 py-3 rounded-xl bg-indigo-600 text-white font-bold hover:bg-indigo-500 transition-colors'
            }, '保存')
          ])
        ])
      ]),

      // 重置密码对话框
      showResetPasswordDialog.value && resetPasswordUser.value && h('div', { 
        class: 'fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4',
        onClick: () => { showResetPasswordDialog.value = false }
      }, [
        h('div', { 
          class: 'bg-white rounded-3xl p-8 max-w-md w-full',
          onClick: (e: Event) => e.stopPropagation()
        }, [
          h('div', { class: 'flex items-center justify-between mb-6' }, [
            h('h3', { class: 'text-2xl font-black text-slate-900' }, '重置密码'),
            h('button', {
              onClick: () => { showResetPasswordDialog.value = false },
              class: 'p-2 hover:bg-slate-100 rounded-xl transition-colors'
            }, h(LucideIcons.X, { size: 20 }))
          ]),
          
          h('div', { class: 'space-y-4' }, [
            h('div', { class: 'bg-amber-50 border border-amber-200 rounded-xl p-4' }, [
              h('div', { class: 'flex items-start gap-3' }, [
                h(LucideIcons.AlertTriangle, { size: 20, class: 'text-amber-600 flex-shrink-0 mt-0.5' }),
                h('div', { class: 'text-sm text-amber-800' }, [
                  h('p', { class: 'font-bold' }, '重置用户密码'),
                  h('p', { class: 'mt-1' }, `将为用户"${resetPasswordUser.value.username}"设置新密码`)
                ])
              ])
            ]),
            
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '新密码'),
              h('input', {
                type: 'password',
                value: newPassword.value,
                onInput: (e: any) => { newPassword.value = e.target.value },
                placeholder: '请输入新密码（至少6位）',
                class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500'
              })
            ]),
          ]),
          
          h('div', { class: 'flex gap-3 mt-8' }, [
            h('button', {
              onClick: () => { showResetPasswordDialog.value = false },
              class: 'flex-1 px-6 py-3 rounded-xl border border-slate-200 font-bold hover:bg-slate-50 transition-colors'
            }, '取消'),
            h('button', {
              onClick: resetPassword,
              class: 'flex-1 px-6 py-3 rounded-xl bg-amber-600 text-white font-bold hover:bg-amber-500 transition-colors'
            }, '确认重置')
          ])
        ])
      ]),

      // 角色管理对话框
      showRoleDialog.value && roleUser.value && h('div', { 
        class: 'fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4',
        onClick: () => { showRoleDialog.value = false }
      }, [
        h('div', { 
          class: 'bg-white rounded-3xl p-8 max-w-2xl w-full max-h-[80vh] overflow-auto',
          onClick: (e: Event) => e.stopPropagation()
        }, [
          h('div', { class: 'flex items-center justify-between mb-6' }, [
            h('h3', { class: 'text-2xl font-black text-slate-900' }, '管理用户角色'),
            h('button', {
              onClick: () => { showRoleDialog.value = false },
              class: 'p-2 hover:bg-slate-100 rounded-xl transition-colors'
            }, h(LucideIcons.X, { size: 20 }))
          ]),
          
          h('div', { class: 'space-y-6' }, [
            h('div', { class: 'bg-blue-50 border border-blue-200 rounded-xl p-4' }, [
              h('div', { class: 'flex items-start gap-3' }, [
                h(LucideIcons.Info, { size: 20, class: 'text-blue-600 flex-shrink-0 mt-0.5' }),
                h('div', { class: 'text-sm text-blue-800' }, [
                  h('p', { class: 'font-bold' }, `为用户"${roleUser.value.username}"分配角色`),
                  h('p', { class: 'mt-1' }, '一个用户可以拥有多个角色，系统会合并所有角色的权限')
                ])
              ])
            ]),
            
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-3' }, '选择角色'),
              h('div', { class: 'space-y-3' }, 
                availableRoles.value.map(role => 
                  h('label', { 
                    class: `flex items-start gap-4 p-4 rounded-xl border-2 cursor-pointer transition-all ${
                      selectedRoles.value.includes(role.code)
                        ? 'border-indigo-500 bg-indigo-50'
                        : 'border-slate-200 hover:border-indigo-200'
                    }`
                  }, [
                    h('input', {
                      type: 'checkbox',
                      checked: selectedRoles.value.includes(role.code),
                      onChange: (e: any) => {
                        if (e.target.checked) {
                          selectedRoles.value.push(role.code);
                        } else {
                          selectedRoles.value = selectedRoles.value.filter((r: string) => r !== role.code);
                        }
                      },
                      class: 'mt-1 w-5 h-5 rounded border-slate-300'
                    }),
                    h('div', { class: 'flex-1' }, [
                      h('div', { class: 'font-bold text-slate-900' }, role.name),
                      h('div', { class: 'text-sm text-slate-600 mt-1' }, role.description)
                    ])
                  ])
                )
              )
            ]),
            
          ]),
          
          h('div', { class: 'flex gap-3 mt-8' }, [
            h('button', {
              onClick: () => { showRoleDialog.value = false },
              class: 'flex-1 px-6 py-3 rounded-xl border border-slate-200 font-bold hover:bg-slate-50 transition-colors'
            }, '取消'),
            h('button', {
              onClick: saveRoles,
              class: 'flex-1 px-6 py-3 rounded-xl bg-indigo-600 text-white font-bold hover:bg-indigo-500 transition-colors'
            }, '保存')
          ])
        ])
      ]),
    ]);
  }
});
