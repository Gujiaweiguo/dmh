// 分销商审批视图
import { ref, onMounted } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import { distributorApi } from '../services/distributorApi';

export const DistributorApprovalView = (props: { brandId?: number }) => {
  const applications = ref<any[]>([]);
  const loading = ref(false);
  const showApproveModal = ref(false);
  const currentApplication = ref<any>(null);
  const approveForm = ref({
    action: 'approve',
    level: 1,
    reason: ''
  });
  const processing = ref(false);

  // 加载申请列表
  const loadApplications = async () => {
    try {
      loading.value = true;
      const response: any = await distributorApi.getApplications(props.brandId || 1, {
        status: 'pending',
        pageSize: 100
      });
      if (response.code === 200) {
        applications.value = response.data.applications || [];
      }
    } catch (error) {
      console.error('加载申请列表失败:', error);
    } finally {
      loading.value = false;
    }
  };

  // 打开审批弹窗
  const openApproveModal = (application: any) => {
    currentApplication.value = application;
    approveForm.value = {
      action: 'approve',
      level: 1,
      reason: ''
    };
    showApproveModal.value = true;
  };

  // 提交审批
  const submitApproval = async () => {
    try {
      processing.value = true;
      const response: any = await distributorApi.approveApplication(
        props.brandId || 1,
        currentApplication.value.id,
        approveForm.value
      );
      if (response.code === 200) {
        alert(approveForm.value.action === 'approve' ? '已批准申请' : '已拒绝申请');
        showApproveModal.value = false;
        loadApplications();
      } else {
        alert(response.message || '操作失败');
      }
    } catch (error: any) {
      alert(error.response?.data?.message || '操作失败');
    } finally {
      processing.value = false;
    }
  };

  onMounted(() => {
    loadApplications();
  });

  return {
    applications,
    loading,
    showApproveModal,
    currentApplication,
    approveForm,
    processing,
    loadApplications,
    openApproveModal,
    submitApproval
  };
};

// 渲染审批视图
export const renderDistributorApprovalView = (viewModel: ReturnType<typeof DistributorApprovalView>) => {
  const { h } = (window as unknown as { Vue?: { h: typeof h } }).Vue || { h: () => null };

  return h('div', { class: 'space-y-6' }, [
    // 头部
    h('div', { class: 'flex justify-between items-center' }, [
      h('div', [
        h('h2', { class: 'text-2xl font-black text-slate-900' }, '分销商审批'),
        h('p', { class: 'text-slate-400 text-sm mt-1' }, '审批用户成为分销商的申请')
      ])
    ]),

    // 申请列表
    h('div', { class: 'bg-white rounded-3xl border border-slate-100 overflow-hidden shadow-sm' }, [
      viewModel.loading.value ?
        h('div', { class: 'p-12 text-center' }, ['加载中...']) :
        viewModel.applications.value.length === 0 ?
          h('div', { class: 'p-12 text-center' }, [
            h(LucideIcons.CheckCircle, { size: 48, class: 'mx-auto text-emerald-300 mb-4' }),
            h('p', { class: 'text-slate-500 text-lg' }, '暂无待审批申请')
          ]) :
          h('table', { class: 'w-full text-left' }, [
            h('thead', { class: 'bg-slate-50' }, [
              h('tr', [
                '申请ID', '申请人', '品牌', '申请理由', '申请时间', '状态', '操作'
              ].map(th => h('th', { class: 'px-6 py-4 text-xs font-black text-slate-400 uppercase tracking-widest' }, th)))
            ]),
            h('tbody', viewModel.applications.value.map(app => h('tr', { class: 'border-b border-slate-50 last:border-0 hover:bg-slate-50/40' }, [
              h('td', { class: 'px-6 py-4 text-sm text-slate-400 font-mono' }, String(app.id)),
              h('td', { class: 'px-6 py-4 text-sm font-bold text-slate-900' }, app.username || `用户${app.userId}`),
              h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, app.brandName || `品牌${app.brandId}`),
              h('td', { class: 'px-6 py-4 text-sm text-slate-600 max-w-xs truncate' }, app.reason || '-'),
              h('td', { class: 'px-6 py-4 text-sm text-slate-500' }, app.createdAt),
              h('td', { class: 'px-6 py-4' }, [
                h('span', {
                  class: `px-2 py-1 rounded-lg text-xs font-bold ${
                    app.status === 'pending' ? 'bg-amber-100 text-amber-800' :
                    app.status === 'approved' ? 'bg-emerald-100 text-emerald-800' :
                    'bg-rose-100 text-rose-800'
                  }`
                }, app.status === 'pending' ? '待审核' : app.status === 'approved' ? '已通过' : '已拒绝')
              ]),
              h('td', { class: 'px-6 py-4' }, [
                h('div', { class: 'flex gap-2' }, [
                  app.status === 'pending' && h('button', {
                    onClick: () => viewModel.openApproveModal(app),
                    class: 'px-3 py-1 text-xs bg-indigo-50 text-indigo-600 rounded-lg hover:bg-indigo-100 flex items-center gap-1'
                  }, [
                    h(LucideIcons.Check, { size: 12 }),
                    '审批'
                  ])
                ].filter(Boolean))
              ])
            ])))
          ])
    ]),

    // 审批弹窗
    viewModel.showApproveModal.value && h('div', {
      class: 'fixed inset-0 bg-black/50 flex items-center justify-center z-50',
      onClick: (e: any) => e.target === e.currentTarget && (viewModel.showApproveModal.value = false)
    }, [
      h('div', { class: 'bg-white rounded-2xl p-6 w-full max-w-md' }, [
        h('h3', { class: 'text-lg font-bold text-slate-900 mb-4' }, '审批分销商申请'),
        viewModel.currentApplication.value && h('div', { class: 'space-y-4 mb-6' }, [
          h('div', [
            h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '申请人'),
            h('p', { class: 'text-slate-900' }, viewModel.currentApplication.value.username || `用户${viewModel.currentApplication.value.userId}`)
          ]),
          h('div', [
            h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '申请理由'),
            h('p', { class: 'text-slate-600 bg-slate-50 p-3 rounded-lg' }, viewModel.currentApplication.value.reason || '无')
          ]),
          h('div', [
            h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '审批操作'),
            h('div', { class: 'flex gap-4' }, [
              h('label', { class: 'flex items-center gap-2 cursor-pointer' }, [
                h('input', {
                  type: 'radio',
                  checked: viewModel.approveForm.value.action === 'approve',
                  onChange: () => viewModel.approveForm.value.action = 'approve',
                  class: 'w-4 h-4 text-indigo-600'
                }),
                h('span', { class: 'text-sm text-slate-700' }, '批准')
              ]),
              h('label', { class: 'flex items-center gap-2 cursor-pointer' }, [
                h('input', {
                  type: 'radio',
                  checked: viewModel.approveForm.value.action === 'reject',
                  onChange: () => viewModel.approveForm.value.action = 'reject',
                  class: 'w-4 h-4 text-indigo-600'
                }),
                h('span', { class: 'text-sm text-slate-700' }, '拒绝')
              ])
            ])
          ]),
          viewModel.approveForm.value.action === 'approve' && h('div', [
            h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '设置级别'),
            h('select', {
              value: viewModel.approveForm.value.level,
              onChange: (e: any) => viewModel.approveForm.value.level = parseInt(e.target.value),
              class: 'w-full px-3 py-2 border border-slate-200 rounded-lg'
            }, [
              h('option', { value: 1 }, '一级分销商'),
              h('option', { value: 2 }, '二级分销商'),
              h('option', { value: 3 }, '三级分销商')
            ])
          ]),
          h('div', [
            h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '备注说明'),
            h('textarea', {
              value: viewModel.approveForm.value.reason,
              onInput: (e: any) => viewModel.approveForm.value.reason = e.target.value,
              placeholder: '请填写审批备注（选填）',
              class: 'w-full px-3 py-2 border border-slate-200 rounded-lg',
              rows: 3
            })
          ])
        ]),
        h('div', { class: 'flex gap-3 justify-end' }, [
          h('button', {
            onClick: () => viewModel.showApproveModal.value = false,
            class: 'px-4 py-2 text-slate-700 bg-slate-100 rounded-lg hover:bg-slate-200'
          }, '取消'),
          h('button', {
            onClick: viewModel.submitApproval,
            disabled: viewModel.processing.value,
            class: 'px-4 py-2 text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 disabled:opacity-50'
          }, viewModel.processing.value ? '处理中...' : '确认')
        ])
      ])
    ])
  ]);
};
