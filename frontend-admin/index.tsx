import { createApp, ref, onMounted, computed, defineComponent, h, reactive, Transition, watch } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import { authApi } from './services/authApi';

// 简单的Badge组件
const Badge = defineComponent({
  props: ['status', 'label'],
  setup(props) {
    const styles = computed(() => ({
      ACTIVE: 'bg-emerald-100 text-emerald-700 border-emerald-200',
      PAUSED: 'bg-slate-100 text-slate-700 border-slate-200',
      PAID: 'bg-indigo-100 text-indigo-700 border-indigo-200',
      PENDING: 'bg-amber-100 text-amber-700 border-amber-200',
      APPROVED: 'bg-emerald-100 text-emerald-700 border-emerald-200',
      REJECTED: 'bg-rose-100 text-rose-700 border-rose-200',
    }[props.status as string] || 'bg-blue-100 text-blue-700 border-blue-200'));

    return () => h('span', { class: `px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider border ${styles.value}` }, props.label || props.status);
  }
});

// 控制面板视图
const DashboardView = defineComponent({
  setup() {
    return () => h('div', { class: 'grid grid-cols-1 md:grid-cols-4 gap-8 animate-in fade-in duration-500' }, [
      { label: '累计报名', value: '1,284', icon: 'Users', color: 'bg-indigo-600' },
      { label: '推广总收益', value: '¥42,050', icon: 'Wallet', color: 'bg-emerald-600' },
      { label: '活跃活动', value: '12', icon: 'Target', color: 'bg-amber-600' },
      { label: '待处理提现', value: '5', icon: 'Clock', color: 'bg-rose-600' }
    ].map(stat => h('div', { class: 'bg-white p-8 rounded-[2.5rem] border border-slate-100 shadow-sm' }, [
      h('div', { class: `w-12 h-12 ${stat.color} text-white rounded-2xl flex items-center justify-center mb-6` }, h((LucideIcons as any)[stat.icon], { size: 24 })),
      h('p', { class: 'text-[10px] font-black text-slate-400 uppercase tracking-widest' }, stat.label),
      h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, stat.value)
    ])));
  }
});

// 用户管理视图
const UserManagementView = defineComponent({
  setup() {
    const users = ref([
      { id: 1, username: 'admin', realName: '系统管理员', role: '平台管理员', status: 'ACTIVE', phone: '138****8888' },
      { id: 3, username: 'user001', realName: '张三', role: '参与者', status: 'ACTIVE', phone: '136****6666' }
    ]);

    return () => h('div', { class: 'space-y-6' }, [
      h('div', { class: 'flex justify-between items-center' }, [
        h('div', [
          h('h2', { class: 'text-2xl font-black text-slate-900' }, '用户管理'),
          h('p', { class: 'text-slate-400 text-sm mt-1' }, '管理系统用户账号和权限')
        ]),
        h('button', { class: 'bg-indigo-600 text-white px-6 py-3 rounded-2xl font-bold hover:bg-indigo-700 transition-colors flex items-center gap-2' }, [
          h(LucideIcons.Plus, { size: 18 }),
          '新增用户'
        ])
      ]),
      h('div', { class: 'bg-white rounded-3xl border border-slate-100 overflow-hidden shadow-sm' }, [
        h('table', { class: 'w-full text-left' }, [
          h('thead', { class: 'bg-slate-50' }, [
            h('tr', [
              'ID', '用户名', '真实姓名', '角色', '手机号', '状态', '操作'
            ].map(th => h('th', { class: 'px-6 py-4 text-xs font-black text-slate-400 uppercase tracking-widest' }, th)))
          ]),
          h('tbody', users.value.map(user => h('tr', { class: 'border-b border-slate-50 last:border-0 hover:bg-slate-50/40' }, [
            h('td', { class: 'px-6 py-4 text-sm text-slate-400 font-mono' }, String(user.id)),
            h('td', { class: 'px-6 py-4 text-sm font-bold text-slate-900' }, user.username),
            h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, user.realName),
            h('td', { class: 'px-6 py-4 text-sm' }, [
              h('span', { class: 'px-2 py-1 bg-blue-100 text-blue-800 rounded-lg text-xs font-bold' }, user.role)
            ]),
            h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, user.phone),
            h('td', { class: 'px-6 py-4' }, [h(Badge, { status: user.status, label: user.status === 'ACTIVE' ? '正常' : '禁用' })]),
            h('td', { class: 'px-6 py-4' }, [
              h('div', { class: 'flex gap-2' }, [
                h('button', { class: 'px-3 py-1 text-xs bg-indigo-50 text-indigo-600 rounded-lg hover:bg-indigo-100' }, '编辑'),
                h('button', { class: 'px-3 py-1 text-xs bg-red-50 text-red-600 rounded-lg hover:bg-red-100' }, '删除')
              ])
            ])
          ])))
        ])
      ])
    ]);
  }
});

// 品牌管理视图
const BrandManagementView = defineComponent({
  setup() {
    const brands = ref([
      { id: 1, name: '科技公司A', logo: 'https://api.dicebear.com/7.x/initials/svg?seed=TechA', status: 'ACTIVE' },
      { id: 2, name: '教育机构B', logo: 'https://api.dicebear.com/7.x/initials/svg?seed=EduB', status: 'ACTIVE' },
      { id: 3, name: '电商平台C', logo: 'https://api.dicebear.com/7.x/initials/svg?seed=EcomC', status: 'PAUSED' }
    ]);

    return () => h('div', { class: 'space-y-6' }, [
      h('div', { class: 'flex justify-between items-center' }, [
        h('div', [
          h('h2', { class: 'text-2xl font-black text-slate-900' }, '品牌管理'),
          h('p', { class: 'text-slate-400 text-sm mt-1' }, '管理入驻平台的合作品牌')
        ]),
        h('button', { class: 'bg-purple-600 text-white px-6 py-3 rounded-2xl font-bold hover:bg-purple-700 transition-colors flex items-center gap-2' }, [
          h(LucideIcons.Plus, { size: 18 }),
          '新增品牌'
        ])
      ]),
      h('div', { class: 'grid grid-cols-1 md:grid-cols-3 gap-6' }, brands.value.map(brand => 
        h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm hover:shadow-lg transition-all' }, [
          h('div', { class: 'flex items-center gap-4 mb-4' }, [
            h('img', { src: brand.logo, class: 'w-12 h-12 rounded-2xl border-2 border-slate-100' }),
            h('div', { class: 'flex-1' }, [
              h('h3', { class: 'text-lg font-black text-slate-900' }, brand.name)
            ])
          ]),
          h('div', { class: 'flex items-center justify-between' }, [
            h(Badge, { status: brand.status, label: brand.status === 'ACTIVE' ? '运营中' : '已暂停' }),
            h('div', { class: 'flex gap-2' }, [
              h('button', { class: 'p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors' }, 
                h(LucideIcons.Edit3, { size: 16 })
              ),
              h('button', { class: 'p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors' }, 
                h(LucideIcons.Trash2, { size: 16 })
              )
            ])
          ])
        ])
      ))
    ]);
  }
});

// 活动管理视图
const CampaignManagementView = defineComponent({
  setup() {
    const campaigns = ref([]);
    const loading = ref(true);

    // 获取活动列表
    const fetchCampaigns = async () => {
      try {
        const token = localStorage.getItem('dmh_token');
        const response = await fetch('/api/v1/campaigns?page=1&pageSize=100', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });
        if (response.ok) {
          const data = await response.json();
          campaigns.value = (data.campaigns || data.list || []).map(c => ({
            id: c.id,
            name: c.name,
            startTime: c.startTime?.substring(0, 10) || '',
            endTime: c.endTime?.substring(0, 10) || '',
            status: c.status?.toUpperCase() || 'ACTIVE',
            participants: c.orderCount || 0,
            description: c.description,
            rewardRule: c.rewardRule,
            brandId: c.brandId
          }));
        }
      } catch (error) {
        console.error('获取活动列表失败', error);
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
      fetchCampaigns();
    });

    const showCreateModal = ref(false);
    const showViewModal = ref(false);
    const editingCampaign = ref(null);
    
    const campaignForm = reactive({
      name: '',
      description: '',
      startTime: '',
      endTime: '',
      rewardRule: 0,
      brandId: 1,
      formFields: [
        { type: 'text', name: 'name', label: '姓名', required: true, placeholder: '请输入姓名' },
        { type: 'phone', name: 'phone', label: '手机号', required: true, placeholder: '请输入手机号' }
      ]
    });





    const resetForm = () => {
      campaignForm.name = '';
      campaignForm.description = '';
      campaignForm.startTime = '';
      campaignForm.endTime = '';
      campaignForm.rewardRule = 0;
      campaignForm.brandId = 1;
      campaignForm.formFields = [
        { type: 'text', name: 'name', label: '姓名', required: true, placeholder: '请输入姓名' },
        { type: 'phone', name: 'phone', label: '手机号', required: true, placeholder: '请输入手机号' }
      ];
    };

    // 打开创建模态框
    const openCreateModal = () => {
      resetForm();
      showCreateModal.value = true;
    };

    // 打开查看模态框
    const openViewModal = (campaign) => {
      editingCampaign.value = campaign;
      showViewModal.value = true;
    };

    // 关闭模态框
    const closeModals = () => {
      showCreateModal.value = false;
      showViewModal.value = false;
      editingCampaign.value = null;
      resetForm();
    };

    // 添加表单字段
    const addFormField = () => {
      campaignForm.formFields.push({
        type: 'text',
        name: `field_${Date.now()}`,
        label: '新字段',
        required: false,
        placeholder: '请输入内容'
      });
    };

    // 删除表单字段
    const removeFormField = (index) => {
      if (campaignForm.formFields.length > 1) {
        campaignForm.formFields.splice(index, 1);
      }
    };

    // 创建活动
    const createCampaign = async () => {
      try {
        const newCampaign = {
          id: campaigns.value.length + 1,
          name: campaignForm.name,
          startTime: campaignForm.startTime,
          endTime: campaignForm.endTime,
          status: 'ACTIVE',
          participants: 0
        };
        
        campaigns.value.unshift(newCampaign);
        closeModals();
        alert('活动创建成功！');
      } catch (error) {
        alert('创建失败：' + error.message);
      }
    };





    // 删除活动
    const deleteCampaign = (campaign) => {
      if (confirm(`确定要删除活动"${campaign.name}"吗？`)) {
        const index = campaigns.value.findIndex(c => c.id === campaign.id);
        if (index !== -1) {
          campaigns.value.splice(index, 1);
          alert('活动删除成功！');
        }
      }
    };

    // 模态框组件
    const Modal = defineComponent({
      props: ['show', 'title', 'size'],
      emits: ['close'],
      setup(props, { emit, slots }) {
        return () => props.show ? h('div', { 
          class: 'fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50',
          onClick: (e) => e.target === e.currentTarget && emit('close')
        }, [
          h('div', { 
            class: `bg-white rounded-3xl p-8 mx-4 max-h-[90vh] overflow-y-auto ${
              props.size === 'large' ? 'w-full max-w-6xl' : 'w-full max-w-2xl'
            }`
          }, [
            h('div', { class: 'flex justify-between items-center mb-6' }, [
              h('h3', { class: 'text-2xl font-black text-slate-900' }, props.title),
              h('button', { 
                onClick: () => emit('close'),
                class: 'p-2 hover:bg-slate-100 rounded-xl transition-colors'
              }, h(LucideIcons.X, { size: 20 }))
            ]),
            slots.default?.()
          ])
        ]) : null;
      }
    });

    // 表单组件
    const FormField = defineComponent({
      props: ['label', 'type', 'value', 'placeholder'],
      emits: ['update:value'],
      setup(props, { emit }) {
        return () => h('div', { class: 'mb-4' }, [
          h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, props.label),
          h('input', {
            type: props.type || 'text',
            value: props.value,
            placeholder: props.placeholder,
            onInput: (e) => emit('update:value', e.target.value),
            class: 'w-full px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-transparent'
          })
        ]);
      }
    });

    return () => h('div', { class: 'space-y-6' }, [
      h('div', { class: 'flex justify-between items-center' }, [
        h('div', [
          h('h2', { class: 'text-2xl font-black text-slate-900' }, '活动管理'),
          h('p', { class: 'text-slate-400 text-sm mt-1' }, '创建和管理营销活动')
        ]),
        h('button', { 
          onClick: openCreateModal,
          class: 'bg-emerald-600 text-white px-6 py-3 rounded-2xl font-bold hover:bg-emerald-700 transition-colors flex items-center gap-2' 
        }, [
          h(LucideIcons.Plus, { size: 18 }),
          '创建活动'
        ])
      ]),
      h('div', { class: 'bg-white rounded-3xl border border-slate-100 overflow-hidden shadow-sm' }, [
        h('table', { class: 'w-full text-left' }, [
          h('thead', { class: 'bg-slate-50' }, [
            h('tr', [
              'ID', '活动名称', '时间范围', '参与数据', '转化数据', '状态', '操作'
            ].map(th => h('th', { class: 'px-6 py-4 text-xs font-black text-slate-400 uppercase tracking-widest' }, th)))
          ]),
          h('tbody', campaigns.value.map(campaign => h('tr', { class: 'border-b border-slate-50 last:border-0 hover:bg-slate-50/40' }, [
            h('td', { class: 'px-6 py-4 text-sm text-slate-400 font-mono' }, String(campaign.id)),
            h('td', { class: 'px-6 py-4' }, [
              h('div', { class: 'text-sm font-bold text-slate-900' }, campaign.name),
              h('div', { class: 'text-xs text-slate-500 mt-1' }, campaign.description || '暂无描述')
            ]),
            h('td', { class: 'px-6 py-4' }, [
              h('div', { class: 'text-sm text-slate-600' }, `${campaign.startTime} 至`),
              h('div', { class: 'text-sm text-slate-600' }, campaign.endTime)
            ]),
            h('td', { class: 'px-6 py-4' }, [
              h('div', { class: 'flex items-center gap-4' }, [
                h('div', { class: 'text-center' }, [
                  h('div', { class: 'text-lg font-black text-indigo-600' }, String(campaign.participants || 0)),
                  h('div', { class: 'text-xs text-slate-500' }, '总参与')
                ]),
                h('div', { class: 'text-center' }, [
                  h('div', { class: 'text-lg font-black text-emerald-600' }, String(Math.floor((campaign.participants || 0) * 0.8))),
                  h('div', { class: 'text-xs text-slate-500' }, '有效报名')
                ])
              ])
            ]),
            h('td', { class: 'px-6 py-4' }, [
              h('div', { class: 'flex items-center gap-4' }, [
                h('div', { class: 'text-center' }, [
                  h('div', { class: 'text-lg font-black text-amber-600' }, String(Math.floor((campaign.participants || 0) * 0.15))),
                  h('div', { class: 'text-xs text-slate-500' }, '转化成功')
                ]),
                h('div', { class: 'text-center' }, [
                  h('div', { class: 'text-lg font-black text-rose-600' }, `${Math.floor((campaign.participants || 0) * 0.15 / Math.max(campaign.participants || 1, 1) * 100)}%`),
                  h('div', { class: 'text-xs text-slate-500' }, '转化率')
                ])
              ])
            ]),
            h('td', { class: 'px-6 py-4' }, [h(Badge, { status: campaign.status, label: campaign.status === 'ACTIVE' ? '进行中' : '已暂停' })]),
            h('td', { class: 'px-6 py-4' }, [
              h('div', { class: 'flex gap-2' }, [
                h('button', { 
                  onClick: () => deleteCampaign(campaign),
                  class: 'px-3 py-1 text-xs bg-red-50 text-red-600 rounded-lg hover:bg-red-100' 
                }, '删除'),
                h('button', { 
                  onClick: () => openViewModal(campaign),
                  class: 'px-3 py-1 text-xs bg-emerald-50 text-emerald-600 rounded-lg hover:bg-emerald-100' 
                }, '详情')
              ])
            ])
          ])))
        ])
      ]),

      // 创建活动模态框
      h(Modal, { 
        show: showCreateModal.value, 
        title: '创建新活动',
        onClose: closeModals
      }, {
        default: () => h('form', { 
          onSubmit: (e) => { e.preventDefault(); createCampaign(); },
          class: 'space-y-4'
        }, [
          h(FormField, {
            label: '活动名称',
            value: campaignForm.name,
            placeholder: '请输入活动名称',
            'onUpdate:value': (val) => campaignForm.name = val
          }),
          h('div', { class: 'mb-4' }, [
            h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '活动描述'),
            h('textarea', {
              value: campaignForm.description,
              placeholder: '请输入活动描述',
              onInput: (e) => campaignForm.description = e.target.value,
              class: 'w-full px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-transparent h-24 resize-none'
            })
          ]),
          h('div', { class: 'grid grid-cols-2 gap-4' }, [
            h(FormField, {
              label: '开始时间',
              type: 'date',
              value: campaignForm.startTime,
              'onUpdate:value': (val) => campaignForm.startTime = val
            }),
            h(FormField, {
              label: '结束时间',
              type: 'date',
              value: campaignForm.endTime,
              'onUpdate:value': (val) => campaignForm.endTime = val
            })
          ]),
          h(FormField, {
            label: '奖励金额（元）',
            type: 'number',
            value: campaignForm.rewardRule,
            placeholder: '请输入奖励金额',
            'onUpdate:value': (val) => campaignForm.rewardRule = Number(val)
          }),
          
          // 动态表单字段配置
          h('div', { class: 'border-t pt-4' }, [
            h('div', { class: 'flex justify-between items-center mb-4' }, [
              h('h4', { class: 'text-lg font-bold text-slate-900' }, '报名表单字段'),
              h('button', {
                type: 'button',
                onClick: addFormField,
                class: 'px-3 py-1 bg-blue-50 text-blue-600 rounded-lg text-sm font-bold hover:bg-blue-100'
              }, '+ 添加字段')
            ]),
            h('div', { class: 'space-y-3' }, campaignForm.formFields.map((field, index) => 
              h('div', { class: 'p-4 border border-slate-200 rounded-xl bg-slate-50' }, [
                h('div', { class: 'grid grid-cols-4 gap-3 mb-2' }, [
                  h('select', {
                    value: field.type,
                    onChange: (e) => field.type = e.target.value,
                    class: 'px-3 py-2 border border-slate-200 rounded-lg text-sm'
                  }, [
                    h('option', { value: 'text' }, '文本'),
                    h('option', { value: 'phone' }, '手机号'),
                    h('option', { value: 'email' }, '邮箱'),
                    h('option', { value: 'select' }, '选择')
                  ]),
                  h('input', {
                    type: 'text',
                    value: field.name,
                    placeholder: '字段名称',
                    onInput: (e) => field.name = e.target.value,
                    class: 'px-3 py-2 border border-slate-200 rounded-lg text-sm'
                  }),
                  h('input', {
                    type: 'text',
                    value: field.label,
                    placeholder: '显示标签',
                    onInput: (e) => field.label = e.target.value,
                    class: 'px-3 py-2 border border-slate-200 rounded-lg text-sm'
                  }),
                  h('div', { class: 'flex items-center gap-2' }, [
                    h('label', { class: 'flex items-center gap-1 text-sm' }, [
                      h('input', {
                        type: 'checkbox',
                        checked: field.required,
                        onChange: (e) => field.required = e.target.checked,
                        class: 'rounded'
                      }),
                      '必填'
                    ]),
                    h('button', {
                      type: 'button',
                      onClick: () => removeFormField(index),
                      class: 'p-1 text-red-500 hover:bg-red-50 rounded'
                    }, h(LucideIcons.Trash2, { size: 14 }))
                  ])
                ]),
                h('input', {
                  type: 'text',
                  value: field.placeholder,
                  placeholder: '占位符文本',
                  onInput: (e) => field.placeholder = e.target.value,
                  class: 'w-full px-3 py-2 border border-slate-200 rounded-lg text-sm'
                })
              ])
            ))
          ]),

          h('div', { class: 'flex gap-4 pt-4' }, [
            h('button', {
              type: 'button',
              onClick: closeModals,
              class: 'flex-1 px-6 py-3 border border-slate-200 text-slate-600 rounded-xl font-bold hover:bg-slate-50'
            }, '取消'),
            h('button', {
              type: 'submit',
              class: 'flex-1 px-6 py-3 bg-emerald-600 text-white rounded-xl font-bold hover:bg-emerald-700'
            }, '创建活动')
          ])
        ])
      }),

      // 活动查看模态框
      h(Modal, { 
        show: showViewModal.value, 
        title: '查看活动详情',
        onClose: closeModals
      }, {
        default: () => editingCampaign.value ? h('div', { class: 'space-y-6' }, [
          // 活动描述
          h('div', { class: 'bg-white border border-slate-200 p-6 rounded-2xl' }, [
            h('h4', { class: 'text-lg font-bold text-slate-900 mb-4 flex items-center gap-2' }, [
              h(LucideIcons.FileText, { size: 20 }),
              '活动描述'
            ]),
            h('p', { class: 'text-slate-700 leading-relaxed' }, 
              editingCampaign.value.description || '暂无活动描述'
            )
          ]),

          // 表单字段预览
          h('div', { class: 'bg-white border border-slate-200 p-6 rounded-2xl' }, [
            h('h4', { class: 'text-lg font-bold text-slate-900 mb-4 flex items-center gap-2' }, [
              h(LucideIcons.FormInput, { size: 20 }),
              '报名表单字段'
            ]),
            h('div', { class: 'grid grid-cols-1 md:grid-cols-2 gap-4' }, campaignForm.formFields.map((field, index) => 
              h('div', { 
                key: `view-field-${index}`,
                class: 'p-4 bg-slate-50 rounded-xl border border-slate-200' 
              }, [
                h('div', { class: 'flex items-center justify-between mb-2' }, [
                  h('span', { class: 'font-bold text-slate-900' }, field.label),
                  field.required && h('span', { class: 'text-xs bg-red-100 text-red-600 px-2 py-1 rounded-full' }, '必填')
                ]),
                h('div', { class: 'text-sm text-slate-600' }, [
                  h('span', { class: 'inline-block bg-blue-100 text-blue-600 px-2 py-1 rounded mr-2 text-xs' }, 
                    field.type === 'text' ? '文本' : 
                    field.type === 'phone' ? '手机号' : 
                    field.type === 'email' ? '邮箱' : '选择'
                  ),
                  field.placeholder && h('span', { class: 'text-slate-500' }, `占位符: ${field.placeholder}`)
                ])
              ])
            ))
          ]),

          // 操作记录
          h('div', { class: 'bg-white border border-slate-200 p-6 rounded-2xl' }, [
            h('h4', { class: 'text-lg font-bold text-slate-900 mb-4 flex items-center gap-2' }, [
              h(LucideIcons.Clock, { size: 20 }),
              '操作记录'
            ]),
            h('div', { class: 'space-y-3' }, [
              h('div', { class: 'flex items-center gap-3 p-3 bg-slate-50 rounded-lg' }, [
                h('div', { class: 'w-2 h-2 bg-green-500 rounded-full' }),
                h('div', { class: 'flex-1' }, [
                  h('p', { class: 'text-sm font-medium text-slate-900' }, '活动创建'),
                  h('p', { class: 'text-xs text-slate-500' }, '2024-01-15 10:30:00')
                ])
              ]),
              h('div', { class: 'flex items-center gap-3 p-3 bg-slate-50 rounded-lg' }, [
                h('div', { class: 'w-2 h-2 bg-blue-500 rounded-full' }),
                h('div', { class: 'flex-1' }, [
                  h('p', { class: 'text-sm font-medium text-slate-900' }, '活动启动'),
                  h('p', { class: 'text-xs text-slate-500' }, '2024-01-15 12:00:00')
                ])
              ]),
              h('div', { class: 'flex items-center gap-3 p-3 bg-slate-50 rounded-lg' }, [
                h('div', { class: 'w-2 h-2 bg-amber-500 rounded-full' }),
                h('div', { class: 'flex-1' }, [
                  h('p', { class: 'text-sm font-medium text-slate-900' }, '最后更新'),
                  h('p', { class: 'text-xs text-slate-500' }, '2024-01-20 15:45:00')
                ])
              ])
            ])
          ]),

          // 快速操作
          h('div', { class: 'flex gap-4 pt-4 border-t' }, [
            h('button', {
              type: 'button',
              onClick: closeModals,
              class: 'flex-1 px-6 py-3 border border-slate-200 text-slate-600 rounded-xl font-bold hover:bg-slate-50'
            }, '关闭')
          ])
        ]) : null
      })
    ]);
  }
});

// 系统设置视图
const SystemSettingsView = defineComponent({
  setup() {
    return () => h('div', { class: 'space-y-6' }, [
      h('div', [
        h('h2', { class: 'text-2xl font-black text-slate-900' }, '系统设置'),
        h('p', { class: 'text-slate-400 text-sm mt-1' }, '系统配置和管理')
      ]),
      h('div', { class: 'grid grid-cols-1 md:grid-cols-2 gap-6' }, [
        h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm' }, [
          h('div', { class: 'flex items-center gap-3 mb-4' }, [
            h('div', { class: 'w-10 h-10 bg-blue-100 rounded-2xl flex items-center justify-center' }, 
              h(LucideIcons.Database, { size: 20, class: 'text-blue-600' })
            ),
            h('h3', { class: 'text-lg font-black text-slate-900' }, '数据库配置')
          ]),
          h('p', { class: 'text-slate-600 text-sm mb-4' }, '管理数据库连接和备份设置'),
          h('button', { class: 'px-4 py-2 bg-blue-50 text-blue-600 rounded-xl text-sm font-bold hover:bg-blue-100' }, '配置')
        ]),
        h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm' }, [
          h('div', { class: 'flex items-center gap-3 mb-4' }, [
            h('div', { class: 'w-10 h-10 bg-green-100 rounded-2xl flex items-center justify-center' }, 
              h(LucideIcons.Shield, { size: 20, class: 'text-green-600' })
            ),
            h('h3', { class: 'text-lg font-black text-slate-900' }, '安全设置')
          ]),
          h('p', { class: 'text-slate-600 text-sm mb-4' }, '管理系统安全策略和权限'),
          h('button', { class: 'px-4 py-2 bg-green-50 text-green-600 rounded-xl text-sm font-bold hover:bg-green-100' }, '配置')
        ])
      ])
    ]);
  }
});

// 主应用
const AdminApp = defineComponent({
  setup() {
    const isLoggedIn = ref(false);
    const loginForm = reactive({
      username: '',
      password: ''
    });
    const loginError = ref('');
    const loginLoading = ref(false);
    const activeTab = ref('dashboard');

    // 检查登录状态
    const checkLoginStatus = () => {
      isLoggedIn.value = authApi.isLoggedIn();
    };

    // 监听登录状态变化
    onMounted(() => {
      checkLoginStatus();
      setInterval(checkLoginStatus, 1000);
    });

    // 快速填充功能
    const quickFillAdmin = () => {
      loginForm.username = 'admin';
      loginForm.password = '123456';
      loginError.value = '';
    };



    // 登录处理
    const handleLogin = async () => {
      if (!loginForm.username || !loginForm.password) {
        loginError.value = '请输入用户名和密码';
        return;
      }

      loginLoading.value = true;
      loginError.value = '';
      
      try {
        await authApi.login(loginForm);
        isLoggedIn.value = true;
        loginForm.username = '';
        loginForm.password = '';
      } catch (error) {
        loginError.value = error instanceof Error ? error.message : '登录失败';
      } finally {
        loginLoading.value = false;
      }
    };

    // 退出登录
    const handleLogout = () => {
      if (confirm('确定要退出登录吗？')) {
        authApi.logout();
        isLoggedIn.value = false;
        activeTab.value = 'dashboard';
      }
    };

    // 侧边栏菜单项
    const sidebarItems = [
      { id: 'dashboard', label: '控制面板', icon: 'LayoutDashboard' },
      { id: 'users', label: '用户管理', icon: 'Users' },
      { id: 'brands', label: '品牌管理', icon: 'Shield' },
      { id: 'campaigns', label: '活动管理', icon: 'Flag' },
      { id: 'system', label: '系统设置', icon: 'Settings' },
    ];

    return () => {
      // 未登录状态 - 显示登录界面
      if (!isLoggedIn.value) {
        return h('div', { class: 'min-h-screen bg-gradient-to-br from-indigo-50 to-slate-100 flex items-center justify-center p-4' }, [
          h('div', { class: 'bg-white rounded-3xl shadow-2xl p-8 w-full max-w-md' }, [
            h('div', { class: 'text-center mb-8' }, [
              h('div', { class: 'w-16 h-16 bg-indigo-600 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-lg shadow-indigo-600/30' }, 
                h(LucideIcons.Zap, { class: 'text-white', size: 32 })
              ),
              h('h1', { class: 'text-2xl font-black text-slate-900 mb-2' }, 'DMH 管理后台'),
              h('p', { class: 'text-slate-500 text-sm' }, '数字营销中台管理系统')
            ]),
            h('form', { 
              onSubmit: (e: Event) => { 
                e.preventDefault(); 
                handleLogin(); 
              },
              class: 'space-y-4' 
            }, [
              h('div', [
                h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '用户名'),
                h('input', {
                  type: 'text',
                  value: loginForm.username,
                  onInput: (e: any) => loginForm.username = e.target.value,
                  class: 'w-full px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-transparent',
                  placeholder: '请输入用户名'
                })
              ]),
              h('div', [
                h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '密码'),
                h('input', {
                  type: 'password',
                  value: loginForm.password,
                  onInput: (e: any) => loginForm.password = e.target.value,
                  class: 'w-full px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-transparent',
                  placeholder: '请输入密码'
                })
              ]),
              loginError.value && h('div', { class: 'text-red-600 text-sm text-center p-3 bg-red-50 rounded-xl border border-red-200' }, loginError.value),
              h('button', {
                type: 'submit',
                disabled: loginLoading.value,
                class: 'w-full bg-indigo-600 text-white py-3 rounded-xl font-medium hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors'
              }, loginLoading.value ? '登录中...' : '登录')
            ]),
            h('div', { class: 'mt-6 text-center text-sm text-slate-500' }, [
              h('div', { class: 'mb-4 p-4 bg-amber-50 border border-amber-200 rounded-2xl' }, [
                h('p', { class: 'text-amber-800 font-bold mb-2' }, '⚠️ 测试账号'),
                h('div', { class: 'text-amber-700 text-xs space-y-1' }, [
                  h('p', '管理员: admin / 123456')
                ]),
                h('div', { class: 'flex gap-2 mt-3' }, [
                  h('button', {
                    type: 'button',
                    onClick: quickFillAdmin,
                    class: 'w-full px-3 py-2 bg-amber-100 text-amber-800 rounded-xl text-xs font-bold hover:bg-amber-200 transition-colors'
                  }, '填充管理员')
                ])
              ])
            ])
          ])
        ]);
      }

      // 已登录状态 - 显示完整的管理界面
      return h('div', { class: 'flex h-screen overflow-hidden' }, [
        // 侧边栏
        h('aside', { class: 'w-72 bg-slate-900 h-full flex flex-col shadow-2xl z-20 shrink-0' }, [
          h('div', { class: 'p-10 flex items-center gap-4' }, [
            h('div', { class: 'w-12 h-12 bg-indigo-600 rounded-2xl flex items-center justify-center shadow-lg shadow-indigo-600/30' }, h(LucideIcons.Zap, { class: 'text-white', size: 28 })),
            h('div', [h('h2', { class: 'text-white font-black text-xl leading-none tracking-tighter' }, 'DMH HUB'), h('p', { class: 'text-slate-500 text-[9px] font-black uppercase tracking-widest mt-1' }, 'CORE PLATFORM')])
          ]),
          h('nav', { class: 'flex-1 mt-6 px-6 space-y-1' }, sidebarItems.map(item => 
            h('button', {
              onClick: () => activeTab.value = item.id,
              class: `w-full flex items-center gap-4 px-6 py-4 rounded-2xl text-left transition-all ${
                activeTab.value === item.id 
                  ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30' 
                  : 'text-slate-400 hover:text-white hover:bg-slate-800'
              }`
            }, [
              h((LucideIcons as any)[item.icon], { size: 20 }),
              h('span', { class: 'font-bold text-sm' }, item.label)
            ])
          )),
          h('div', { class: 'p-6 border-t border-slate-800 mt-auto' }, [
            h('button', {
              onClick: handleLogout,
              class: 'w-full flex items-center gap-3 px-4 py-3 rounded-2xl text-slate-400 hover:text-white hover:bg-slate-800 transition-all'
            }, [
              h(LucideIcons.LogOut, { size: 18 }),
              h('span', { class: 'font-bold text-sm' }, '退出登录')
            ])
          ])
        ]),
        
        // 主内容区域
        h('main', { class: 'flex-1 flex flex-col overflow-hidden bg-slate-50' }, [
          h('header', { class: 'bg-white border-b border-slate-100 px-10 py-6 flex items-center justify-between shadow-sm' }, [
            h('div', [
              h('h1', { class: 'text-2xl font-black text-slate-900' }, 
                sidebarItems.find(item => item.id === activeTab.value)?.label || '控制面板'
              ),
              h('p', { class: 'text-slate-400 text-sm mt-1' }, '数字营销中台管理系统')
            ]),
            h('div', { class: 'flex items-center gap-6' }, [
              h('div', { class: 'flex items-center gap-3 border-l pl-6 border-slate-100' }, [
                h('div', { class: 'text-right' }, [h('p', { class: 'text-[10px] font-black text-slate-900' }, '管理员'), h('p', { class: 'text-[9px] font-bold text-slate-400 uppercase' }, 'Super Admin')]),
                h('img', { src: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Admin', class: 'w-10 h-10 rounded-2xl border-2 border-white shadow-sm hover:scale-105 transition-all' })
              ])
            ])
          ]),
          h('div', { class: 'p-10 flex-1 overflow-auto' }, [
            h(Transition, { name: 'fade', mode: 'out-in' }, {
              default: () => {
                if (activeTab.value === 'dashboard') return h(DashboardView);
                if (activeTab.value === 'users') return h(UserManagementView);
                if (activeTab.value === 'brands') return h(BrandManagementView);
                if (activeTab.value === 'campaigns') return h(CampaignManagementView);
                if (activeTab.value === 'system') return h(SystemSettingsView);
                return h(DashboardView);
              }
            })
          ])
        ])
      ]);
    };
  }
});

createApp(AdminApp).mount('#root');