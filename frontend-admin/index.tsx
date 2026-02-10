import { createApp, ref, onMounted, computed, defineComponent, h, reactive, Transition, watch } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import { authApi } from './services/authApi';
import MemberListView from './views/MemberListView';
import MemberDetailView from './views/MemberDetailView';
import MemberMergeView from './views/MemberMergeView';
import MemberExportView from './views/MemberExportView';
import { DistributorManagementView } from './views/DistributorManagementView';
import FeedbackManagementView from './views/FeedbackManagementView.vue';
import VerificationRecordsView from './views/VerificationRecordsView.vue';
import PosterRecordsView from './views/PosterRecordsView.vue';
import { resolveAdminHashRoute, type MemberRoute } from './utils/adminHashRoute';
import { UserManagementView } from './views/UserManagementView';
import './src/index.css';
import './styles/member.css';

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

// 控制面板视图
const BrandManagementView = defineComponent({
  setup() {
    const brands = ref([]);
    const loading = ref(true);

    // 获取品牌列表
    const fetchBrands = async () => {
      try {
        loading.value = true;
        const token = localStorage.getItem('dmh_token');
        const response = await fetch('/api/v1/brands', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });
        if (response.ok) {
          const data = await response.json();
          brands.value = data.brands || [];
        }
      } catch (error) {
        console.error('获取品牌列表失败', error);
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
      fetchBrands();
    });

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
      loading.value ? 
        h('div', { class: 'p-12 text-center' }, [
          h('div', { class: 'inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-purple-600' }),
          h('p', { class: 'mt-4 text-slate-500' }, '加载中...')
        ]) :
        h('div', { class: 'grid grid-cols-1 md:grid-cols-3 gap-6' }, brands.value.map(brand => 
          h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm hover:shadow-lg transition-all' }, [
            h('div', { class: 'flex items-center gap-4 mb-4' }, [
              h('img', { 
                src: brand.logo || `https://api.dicebear.com/7.x/initials/svg?seed=${brand.name}`, 
                class: 'w-12 h-12 rounded-2xl border-2 border-slate-100' 
              }),
              h('div', { class: 'flex-1' }, [
                h('h3', { class: 'text-lg font-black text-slate-900' }, brand.name),
                h('p', { class: 'text-xs text-slate-500 mt-1' }, brand.description || '暂无描述')
              ])
            ]),
            h('div', { class: 'flex items-center justify-between' }, [
              h(Badge, { status: brand.status, label: brand.status === 'active' ? '运营中' : '已暂停' }),
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

// 活动管理视图（平台监控版本 - 只查询不创建）
const CampaignManagementView = defineComponent({
  setup() {
    const campaigns = ref([]);
    const loading = ref(true);
    
    // 筛选条件
    const filters = reactive({
      brandId: '',
      status: '',
      startDate: '',
      endDate: '',
      keyword: ''
    });

    // 动态品牌列表
    const brands = computed(() => {
      const uniqueBrands = new Map();
      allCampaigns.value.forEach(campaign => {
        if (!uniqueBrands.has(campaign.brandId)) {
          uniqueBrands.set(campaign.brandId, {
            id: campaign.brandId,
            name: campaign.brandName || `品牌${campaign.brandId}`
          });
        }
      });
      const brandList = Array.from(uniqueBrands.values());
      console.log('动态品牌列表:', brandList);
      return brandList;
    });

    // 原始活动数据
    const allCampaigns = ref([]);
    
    // 获取活动列表
    const fetchCampaigns = async () => {
      try {
        loading.value = true;
        const token = localStorage.getItem('dmh_token');
        
        const response = await fetch('/api/v1/campaigns?page=1&pageSize=100', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });
        
        if (response.ok) {
          const data = await response.json();
          allCampaigns.value = (data.campaigns || data.list || []).map(c => ({
            id: c.id,
            name: c.name,
            startTime: c.startTime?.substring(0, 10) || '',
            endTime: c.endTime?.substring(0, 10) || '',
            status: c.status?.toLowerCase() || 'active',
            participants: c.orderCount || 0,
            description: c.description,
            rewardRule: c.rewardRule,
            brandId: c.brandId,
            brandName: c.brandName || `品牌${c.brandId}`
          }));
          
          console.log('获取到的活动数据:', allCampaigns.value);
          
          // 应用筛选
          applyFilters();
        }
      } catch (error) {
        console.error('获取活动列表失败', error);
        campaigns.value = [];
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
      fetchCampaigns();
    });

    // 监听筛选条件变化，自动应用筛选
    watch([() => filters.keyword, () => filters.brandId, () => filters.status, () => filters.startDate, () => filters.endDate], () => {
      if (allCampaigns.value.length > 0) {
        applyFilters();
      }
    });

    const showViewModal = ref(false);
    const viewingCampaign = ref(null);

    
    // 应用筛选
    const applyFilters = () => {
      let filtered = [...allCampaigns.value];
      
      // 关键词筛选
      if (filters.keyword.trim()) {
        const keyword = filters.keyword.trim().toLowerCase();
        filtered = filtered.filter(campaign => 
          campaign.name.toLowerCase().includes(keyword) ||
          (campaign.description && campaign.description.toLowerCase().includes(keyword))
        );
      }
      
      // 品牌筛选
      if (filters.brandId) {
        filtered = filtered.filter(campaign => 
          campaign.brandId.toString() === filters.brandId
        );
      }
      
      // 状态筛选
      if (filters.status) {
        filtered = filtered.filter(campaign => 
          campaign.status === filters.status
        );
      }
      
      // 开始时间筛选
      if (filters.startDate) {
        filtered = filtered.filter(campaign => 
          campaign.startTime >= filters.startDate
        );
      }
      
      // 结束时间筛选
      if (filters.endDate) {
        filtered = filtered.filter(campaign => 
          campaign.endTime <= filters.endDate
        );
      }
      
      campaigns.value = filtered;
    };

    // 重置筛选条件
    const resetFilters = () => {
      filters.brandId = '';
      filters.status = '';
      filters.startDate = '';
      filters.endDate = '';
      filters.keyword = '';
      campaigns.value = [...allCampaigns.value];
    };

    // 暂停/恢复活动
    const toggleCampaignStatus = async (campaign, newStatus) => {
      const action = newStatus === 'paused' ? '暂停' : '恢复';
      if (!confirm(`确定要${action}活动"${campaign.name}"吗？`)) return;
      
      try {
        const token = localStorage.getItem('dmh_token');
        const response = await fetch(`/api/v1/campaigns/${campaign.id}/status`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify({ status: newStatus })
        });
        
        if (response.ok) {
          alert(`活动${action}成功！`);
          // 重新加载活动列表以确保数据同步
          await fetchCampaigns();
        } else {
          const data = await response.json();
          alert(`活动${action}失败: ${data.message || '未知错误'}`);
        }
      } catch (error) {
        console.error('操作失败:', error);
        alert(`活动${action}失败！`);
      }
    };

    // 打开查看模态框
    const openViewModal = (campaign) => {
      viewingCampaign.value = campaign;
      showViewModal.value = true;
    };

    // 关闭模态框
    const closeModal = () => {
      showViewModal.value = false;
      viewingCampaign.value = null;
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
          h('h2', { class: 'text-2xl font-black text-slate-900' }, '活动监控'),
          h('p', { class: 'text-slate-400 text-sm mt-1' }, '查看和监控所有营销活动数据')
        ]),
        h('div', { class: 'flex items-center gap-3' }, [
          h('button', { 
            onClick: resetFilters,
            class: 'px-4 py-2 text-slate-600 border border-slate-200 rounded-xl hover:bg-slate-50 transition-colors flex items-center gap-2' 
          }, [
            h(LucideIcons.RotateCcw, { size: 16 }),
            '重置筛选'
          ]),
          h('div', { class: 'text-sm text-slate-500' }, [
            `共找到 ${campaigns.value.length} 个活动`
          ])
        ])
      ]),
      
      // 筛选条件面板
      h('div', { class: 'bg-white rounded-3xl border border-slate-100 p-6 shadow-sm' }, [
        h('h3', { class: 'text-lg font-bold text-slate-900 mb-4 flex items-center gap-2' }, [
          h(LucideIcons.Filter, { size: 20 }),
          '实时筛选',
          h('span', { class: 'text-xs bg-blue-100 text-blue-600 px-2 py-1 rounded-full font-normal' }, '自动应用')
        ]),
        h('div', { class: 'grid grid-cols-1 md:grid-cols-5 gap-4' }, [
          h('div', [
            h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '关键词'),
            h('input', {
              type: 'text',
              value: filters.keyword,
              placeholder: '搜索活动名称',
              onInput: (e: Event) => {
                const target = e.target as HTMLInputElement | null;
                filters.keyword = target?.value ?? '';
              },
              class: 'w-full px-3 py-2 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm'
            })
          ]),
          h('div', [
            h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '品牌'),
            h('select', {
              value: filters.brandId,
              onChange: (e: Event) => {
                const target = e.target as HTMLSelectElement | null;
                filters.brandId = target?.value ?? '';
              },
              class: 'w-full px-3 py-2 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm'
            }, [
              h('option', { value: '' }, '全部品牌'),
              ...brands.value.map(brand => 
                h('option', { value: brand.id.toString() }, brand.name)
              )
            ])
          ]),
          h('div', [
            h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '状态'),
            h('select', {
              value: filters.status,
              onChange: (e: Event) => {
                const target = e.target as HTMLSelectElement | null;
                filters.status = target?.value ?? '';
              },
              class: 'w-full px-3 py-2 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm'
            }, [
              h('option', { value: '' }, '全部状态'),
              h('option', { value: 'active' }, '进行中'),
              h('option', { value: 'paused' }, '已暂停'),
              h('option', { value: 'ended' }, '已结束')
            ])
          ]),
          h('div', [
            h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '开始时间'),
            h('input', {
              type: 'date',
              value: filters.startDate,
              onInput: (e: Event) => {
                const target = e.target as HTMLInputElement | null;
                filters.startDate = target?.value ?? '';
              },
              class: 'w-full px-3 py-2 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm'
            })
          ]),
          h('div', [
            h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '结束时间'),
            h('input', {
              type: 'date',
              value: filters.endDate,
              onInput: (e: Event) => {
                const target = e.target as HTMLInputElement | null;
                filters.endDate = target?.value ?? '';
              },
              class: 'w-full px-3 py-2 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm'
            })
          ])
        ])
      ]),

      // 活动统计卡片
      h('div', { class: 'grid grid-cols-1 md:grid-cols-4 gap-6' }, [
        { label: '总活动数', value: campaigns.value.length, icon: 'Target', color: 'bg-blue-600' },
        { label: '进行中', value: campaigns.value.filter(c => c.status === 'active').length, icon: 'Play', color: 'bg-emerald-600' },
        { label: '已暂停', value: campaigns.value.filter(c => c.status === 'paused').length, icon: 'Pause', color: 'bg-amber-600' },
        { label: '总参与数', value: campaigns.value.reduce((sum, c) => sum + (c.participants || 0), 0), icon: 'Users', color: 'bg-purple-600' }
      ].map(stat => h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm' }, [
        h('div', { class: `w-12 h-12 ${stat.color} text-white rounded-2xl flex items-center justify-center mb-4` }, h((LucideIcons as any)[stat.icon], { size: 24 })),
        h('p', { class: 'text-[10px] font-black text-slate-400 uppercase tracking-widest' }, stat.label),
        h('p', { class: 'text-2xl font-black text-slate-900 mt-1' }, String(stat.value))
      ]))),

      h('div', { class: 'bg-white rounded-3xl border border-slate-100 overflow-hidden shadow-sm' }, [
        loading.value ? 
          h('div', { class: 'p-12 text-center' }, [
            h('div', { class: 'inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600' }),
            h('p', { class: 'mt-4 text-slate-500' }, '加载中...')
          ]) :
          campaigns.value.length === 0 ?
            h('div', { class: 'p-12 text-center' }, [
              h(LucideIcons.Search, { size: 48, class: 'mx-auto text-slate-300 mb-4' }),
              h('p', { class: 'text-slate-500 text-lg' }, '暂无活动数据'),
              h('p', { class: 'text-slate-400 text-sm mt-2' }, '请调整筛选条件或联系品牌方创建活动')
            ]) :
            h('table', { class: 'w-full text-left' }, [
              h('thead', { class: 'bg-slate-50' }, [
                h('tr', [
                  'ID', '活动名称', '品牌', '时间范围', '参与数据', '转化数据', '状态', '操作'
                ].map(th => h('th', { class: 'px-6 py-4 text-xs font-black text-slate-400 uppercase tracking-widest' }, th)))
              ]),
              h('tbody', campaigns.value.map(campaign => h('tr', { class: 'border-b border-slate-50 last:border-0 hover:bg-slate-50/40' }, [
                h('td', { class: 'px-6 py-4 text-sm text-slate-400 font-mono' }, String(campaign.id)),
                h('td', { class: 'px-6 py-4' }, [
                  h('div', { class: 'text-sm font-bold text-slate-900' }, campaign.name),
                  h('div', { class: 'text-xs text-slate-500 mt-1' }, campaign.description || '暂无描述')
                ]),
                h('td', { class: 'px-6 py-4' }, [
                  h('div', { class: 'text-sm font-medium text-slate-700' }, campaign.brandName || `品牌${campaign.brandId}`)
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
                h('td', { class: 'px-6 py-4' }, [h(Badge, { status: campaign.status, label: campaign.status === 'active' ? '进行中' : campaign.status === 'paused' ? '已暂停' : '已结束' })]),
                h('td', { class: 'px-6 py-4' }, [
                  h('div', { class: 'flex gap-2' }, [
                    // 状态控制按钮
                    campaign.status === 'active' ? 
                      h('button', { 
                        onClick: () => toggleCampaignStatus(campaign, 'paused'),
                        class: 'px-3 py-1 text-xs bg-amber-50 text-amber-600 rounded-lg hover:bg-amber-100 flex items-center gap-1' 
                      }, [
                        h(LucideIcons.Pause, { size: 12 }),
                        '暂停'
                      ]) :
                      h('button', { 
                        onClick: () => toggleCampaignStatus(campaign, 'active'),
                        class: 'px-3 py-1 text-xs bg-emerald-50 text-emerald-600 rounded-lg hover:bg-emerald-100 flex items-center gap-1' 
                      }, [
                        h(LucideIcons.Play, { size: 12 }),
                        '恢复'
                      ]),
                    // 查看详情按钮
                    h('button', { 
                      onClick: () => openViewModal(campaign),
                      class: 'px-3 py-1 text-xs bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 flex items-center gap-1' 
                    }, [
                      h(LucideIcons.Eye, { size: 12 }),
                      '查看详情'
                    ])
                  ])
                ])
              ])))
            ])
      ]),

      // 活动查看模态框
      h(Modal, { 
        show: showViewModal.value, 
        title: '活动详情监控',
        size: 'large',
        onClose: closeModal
      }, {
        default: () => viewingCampaign.value ? h('div', { class: 'space-y-6' }, [
          // 活动基本信息
          h('div', { class: 'grid grid-cols-1 md:grid-cols-2 gap-6' }, [
            h('div', { class: 'bg-gradient-to-br from-blue-50 to-indigo-50 border border-blue-200 p-6 rounded-2xl' }, [
              h('h4', { class: 'text-lg font-bold text-blue-900 mb-4 flex items-center gap-2' }, [
                h(LucideIcons.Info, { size: 20 }),
                '基本信息'
              ]),
              h('div', { class: 'space-y-3' }, [
                h('div', { class: 'flex justify-between' }, [
                  h('span', { class: 'text-blue-700 font-medium' }, '活动名称:'),
                  h('span', { class: 'text-blue-900 font-bold' }, viewingCampaign.value.name)
                ]),
                h('div', { class: 'flex justify-between' }, [
                  h('span', { class: 'text-blue-700 font-medium' }, '所属品牌:'),
                  h('span', { class: 'text-blue-900 font-bold' }, viewingCampaign.value.brandName || `品牌${viewingCampaign.value.brandId}`)
                ]),
                h('div', { class: 'flex justify-between' }, [
                  h('span', { class: 'text-blue-700 font-medium' }, '活动状态:'),
                  h(Badge, { status: viewingCampaign.value.status, label: viewingCampaign.value.status === 'active' ? '进行中' : viewingCampaign.value.status === 'paused' ? '已暂停' : '已结束' })
                ]),
                h('div', { class: 'flex justify-between' }, [
                  h('span', { class: 'text-blue-700 font-medium' }, '开始时间:'),
                  h('span', { class: 'text-blue-900 font-bold' }, viewingCampaign.value.startTime)
                ]),
                h('div', { class: 'flex justify-between' }, [
                  h('span', { class: 'text-blue-700 font-medium' }, '结束时间:'),
                  h('span', { class: 'text-blue-900 font-bold' }, viewingCampaign.value.endTime)
                ])
              ])
            ]),
            h('div', { class: 'bg-gradient-to-br from-emerald-50 to-green-50 border border-emerald-200 p-6 rounded-2xl' }, [
              h('h4', { class: 'text-lg font-bold text-emerald-900 mb-4 flex items-center gap-2' }, [
                h(LucideIcons.TrendingUp, { size: 20 }),
                '数据概览'
              ]),
              h('div', { class: 'grid grid-cols-2 gap-4' }, [
                h('div', { class: 'text-center' }, [
                  h('div', { class: 'text-2xl font-black text-emerald-600' }, String(viewingCampaign.value.participants || 0)),
                  h('div', { class: 'text-xs text-emerald-700 font-medium' }, '总参与数')
                ]),
                h('div', { class: 'text-center' }, [
                  h('div', { class: 'text-2xl font-black text-emerald-600' }, String(Math.floor((viewingCampaign.value.participants || 0) * 0.8))),
                  h('div', { class: 'text-xs text-emerald-700 font-medium' }, '有效报名')
                ]),
                h('div', { class: 'text-center' }, [
                  h('div', { class: 'text-2xl font-black text-emerald-600' }, String(Math.floor((viewingCampaign.value.participants || 0) * 0.15))),
                  h('div', { class: 'text-xs text-emerald-700 font-medium' }, '转化成功')
                ]),
                h('div', { class: 'text-center' }, [
                  h('div', { class: 'text-2xl font-black text-emerald-600' }, `${Math.floor((viewingCampaign.value.participants || 0) * 0.15 / Math.max(viewingCampaign.value.participants || 1, 1) * 100)}%`),
                  h('div', { class: 'text-xs text-emerald-700 font-medium' }, '转化率')
                ])
              ])
            ])
          ]),

          // 活动描述
          h('div', { class: 'bg-white border border-slate-200 p-6 rounded-2xl' }, [
            h('h4', { class: 'text-lg font-bold text-slate-900 mb-4 flex items-center gap-2' }, [
              h(LucideIcons.FileText, { size: 20 }),
              '活动描述'
            ]),
            h('p', { class: 'text-slate-700 leading-relaxed' }, 
              viewingCampaign.value.description || '暂无活动描述'
            )
          ]),

          // 奖励规则
          h('div', { class: 'bg-amber-50 border border-amber-200 p-6 rounded-2xl' }, [
            h('h4', { class: 'text-lg font-bold text-amber-900 mb-4 flex items-center gap-2' }, [
              h(LucideIcons.Gift, { size: 20 }),
              '奖励规则'
            ]),
            h('div', { class: 'flex items-center gap-3' }, [
              h('div', { class: 'text-2xl font-black text-amber-600' }, `¥${viewingCampaign.value.rewardRule || 0}`),
              h('div', { class: 'text-amber-700' }, '每次成功转化奖励')
            ])
          ]),

          // 监控指标
          h('div', { class: 'bg-white border border-slate-200 p-6 rounded-2xl' }, [
            h('h4', { class: 'text-lg font-bold text-slate-900 mb-4 flex items-center gap-2' }, [
              h(LucideIcons.BarChart3, { size: 20 }),
              '关键指标监控'
            ]),
            h('div', { class: 'grid grid-cols-1 md:grid-cols-3 gap-4' }, [
              h('div', { class: 'p-4 bg-blue-50 rounded-xl border border-blue-200' }, [
                h('div', { class: 'text-center' }, [
                  h('div', { class: 'text-xl font-black text-blue-600' }, '85%'),
                  h('div', { class: 'text-xs text-blue-700 font-medium mt-1' }, '报名完成率')
                ])
              ]),
              h('div', { class: 'p-4 bg-purple-50 rounded-xl border border-purple-200' }, [
                h('div', { class: 'text-center' }, [
                  h('div', { class: 'text-xl font-black text-purple-600' }, '12.5%'),
                  h('div', { class: 'text-xs text-purple-700 font-medium mt-1' }, '平均转化率')
                ])
              ]),
              h('div', { class: 'p-4 bg-rose-50 rounded-xl border border-rose-200' }, [
                h('div', { class: 'text-center' }, [
                  h('div', { class: 'text-xl font-black text-rose-600' }, '¥' + String((viewingCampaign.value.rewardRule || 0) * Math.floor((viewingCampaign.value.participants || 0) * 0.15))),
                  h('div', { class: 'text-xs text-rose-700 font-medium mt-1' }, '总奖励支出')
                ])
              ])
            ])
          ]),

          // 操作记录
          h('div', { class: 'bg-white border border-slate-200 p-6 rounded-2xl' }, [
            h('h4', { class: 'text-lg font-bold text-slate-900 mb-4 flex items-center gap-2' }, [
              h(LucideIcons.Clock, { size: 20 }),
              '活动时间线'
            ]),
            h('div', { class: 'space-y-3' }, [
              h('div', { class: 'flex items-center gap-3 p-3 bg-slate-50 rounded-lg' }, [
                h('div', { class: 'w-2 h-2 bg-green-500 rounded-full' }),
                h('div', { class: 'flex-1' }, [
                  h('p', { class: 'text-sm font-medium text-slate-900' }, '活动创建'),
                  h('p', { class: 'text-xs text-slate-500' }, viewingCampaign.value.startTime + ' 00:00:00')
                ])
              ]),
              h('div', { class: 'flex items-center gap-3 p-3 bg-slate-50 rounded-lg' }, [
                h('div', { class: 'w-2 h-2 bg-blue-500 rounded-full' }),
                h('div', { class: 'flex-1' }, [
                  h('p', { class: 'text-sm font-medium text-slate-900' }, '活动启动'),
                  h('p', { class: 'text-xs text-slate-500' }, viewingCampaign.value.startTime + ' 09:00:00')
                ])
              ]),
              h('div', { class: 'flex items-center gap-3 p-3 bg-slate-50 rounded-lg' }, [
                h('div', { class: 'w-2 h-2 bg-amber-500 rounded-full' }),
                h('div', { class: 'flex-1' }, [
                  h('p', { class: 'text-sm font-medium text-slate-900' }, '最后数据更新'),
                  h('p', { class: 'text-xs text-slate-500' }, new Date().toLocaleString())
                ])
              ])
            ])
          ]),

          // 关闭按钮
          h('div', { class: 'flex justify-end pt-4 border-t' }, [
            h('button', {
              type: 'button',
              onClick: closeModal,
              class: 'px-6 py-3 bg-slate-600 text-white rounded-xl font-bold hover:bg-slate-700 transition-colors'
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
    const memberRoute = ref<MemberRoute>('list');

	    const syncFromHash = () => {
	      const route = resolveAdminHashRoute(window.location.hash || '');
	      if (route.activeTab) {
	        activeTab.value = route.activeTab;
	      }
	      if (route.memberRoute) {
	        memberRoute.value = route.memberRoute;
	      }
	      if (route.normalizedHash && window.location.hash !== route.normalizedHash) {
	        window.location.hash = route.normalizedHash;
	      }
	    };

	    // 检查登录状态
	    const checkLoginStatus = () => {
	      const loggedIn = authApi.isLoggedIn();
	      if (loggedIn && !authApi.isPlatformAdmin()) {
	        authApi.logout();
	        isLoggedIn.value = false;
	        loginError.value = '管理后台仅限平台管理员访问，请使用 H5 端登录';
	        return;
	      }
	      isLoggedIn.value = loggedIn;
	    };

    // 监听登录状态变化
    onMounted(() => {
      checkLoginStatus();
      setInterval(checkLoginStatus, 1000);
      syncFromHash();
      window.addEventListener('hashchange', syncFromHash);
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
        memberRoute.value = 'list';
        window.location.hash = '#/dashboard';
      }
    };

    const handleNavigate = (tabId: string) => {
      activeTab.value = tabId;
      if (tabId === 'members') {
        memberRoute.value = 'list';
        window.location.hash = '#/members';
        return;
      }
      window.location.hash = `#/${tabId}`;
    };

	    // 侧边栏菜单项
	    const sidebarItems = [
	      { id: 'dashboard', label: '控制面板', icon: 'LayoutDashboard' },
	      { id: 'users', label: '用户管理', icon: 'Users' },
	      { id: 'brands', label: '品牌管理', icon: 'Shield' },
	      { id: 'campaigns', label: '活动监控', icon: 'Monitor' },
	      { id: 'members', label: '会员管理', icon: 'Users' },
	      { id: 'distributor-management', label: '分销监控', icon: 'TrendingUp' },
	      { id: 'feedback', label: '反馈管理', icon: 'MessageSquare' },
	      { id: 'verification-records', label: '核销记录', icon: 'CheckCircle' },
	      { id: 'poster-records', label: '海报记录', icon: 'Image' },
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
              h('h1', { class: 'text-2xl font-black text-slate-900 mb-2' }, 'DMH 平台管理后台'),
              h('p', { class: 'text-slate-500 text-sm' }, '仅限平台管理员访问')
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
                  onInput: (e: Event) => {
                    const target = e.target as HTMLInputElement | null;
                    loginForm.username = target?.value ?? '';
                  },
                  class: 'w-full px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-transparent',
                  placeholder: '请输入用户名'
                })
              ]),
              h('div', [
                h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '密码'),
                h('input', {
                  type: 'password',
                  value: loginForm.password,
                  onInput: (e: Event) => {
                    const target = e.target as HTMLInputElement | null;
                    loginForm.password = target?.value ?? '';
                  },
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
              h('div', { class: 'mb-4 p-4 bg-indigo-50 border border-indigo-200 rounded-2xl' }, [
                h('p', { class: 'text-indigo-800 font-bold mb-2' }, '🔒 平台管理员登录'),
                h('div', { class: 'text-indigo-700 text-xs space-y-1' }, [
                  h('p', '账号: admin / 123456'),
                  h('p', { class: 'text-indigo-600 mt-2' }, '品牌管理员请使用 H5 端访问')
                ]),
                h('div', { class: 'flex gap-2 mt-3' }, [
                  h('button', {
                    type: 'button',
                    onClick: quickFillAdmin,
                    class: 'w-full px-3 py-2 bg-indigo-100 text-indigo-800 rounded-xl text-xs font-bold hover:bg-indigo-200 transition-colors'
                  }, '填充账号')
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
              onClick: () => handleNavigate(item.id),
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
                h('div', { class: 'text-right' }, [h('p', { class: 'text-[10px] font-black text-slate-900' }, authApi.getUsername() || 'Admin'), h('p', { class: 'text-[9px] font-bold text-slate-400 uppercase' }, 'Platform Admin')]),
                h('img', { src: `https://api.dicebear.com/7.x/avataaars/svg?seed=${authApi.getUsername() || 'Admin'}`, class: 'w-10 h-10 rounded-2xl border-2 border-white shadow-sm hover:scale-105 transition-all' })
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
	                if (activeTab.value === 'members') {
	                  if (memberRoute.value === 'detail') return h(MemberDetailView);
	                  if (memberRoute.value === 'merge') return h(MemberMergeView);
	                  if (memberRoute.value === 'export') return h(MemberExportView);
	                  return h(MemberListView);
	                }
                if (activeTab.value === 'distributor-management') {
                  return h(DistributorManagementView, { readOnly: true, isPlatformAdmin: true });
                }
                if (activeTab.value === 'feedback') return h(FeedbackManagementView);
                if (activeTab.value === 'verification-records') return h(VerificationRecordsView);
                if (activeTab.value === 'poster-records') return h(PosterRecordsView);
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
