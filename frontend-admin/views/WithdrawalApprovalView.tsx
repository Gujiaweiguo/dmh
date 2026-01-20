// 提现审批视图（平台管理员）
import { ref, onMounted } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import { distributorApi } from '../services/distributorApi';

export const WithdrawalApprovalView = () => {
  const withdrawals = ref<any[]>([]);
  const loading = ref(false);
  const refreshing = ref(false);
  const page = ref(1);
  const pageSize = ref(20);
  const total = ref(0);
  const selectedBrandId = ref<number | null>(null);
  const statusFilter = ref('pending');

  const showDetailModal = ref(false);
  const currentWithdrawal = ref<any>(null);
  const approvalForm = ref({
    action: 'approve',
    notes: ''
  });
  const processing = ref(false);

  // 加载提现记录
  const loadWithdrawals = async (isRefresh = false) => {
    if (loading.value || (finished.value && !isRefresh)) {
      return;
    }

    loading.value = true;

    try {
      const response: any = await distributorApi.getWithdrawals(
        selectedBrandId.value || 0,
        statusFilter.value,
        page.value,
        pageSize.value
      );
      if (response.code === 200) {
        if (isRefresh) {
          withdrawals.value = response.data.withdrawals || [];
        } else {
          withdrawals.value.push(...(response.data.withdrawals || []));
        }
        total.value = response.data.total || 0;
        finished.value = withdrawals.value.length >= total.value;
      }
    } catch (error) {
      console.error('加载提现记录失败:', error);
      alert('加载失败，请重试');
    } finally {
      loading.value = false;
      refreshing.value = false;
    }
  };

  // 刷新
  const onRefresh = async () => {
    page.value = 1;
    finished.value = false;
    withdrawals.value = [];
    await loadWithdrawals(true);
  };

  // 加载更多
  const onLoad = async () => {
    if (finished.value) {
      return;
    }
    page.value++;
    await loadWithdrawals();
  };

  // 打开详情弹窗
  const openDetailModal = (withdrawal: any) => {
    currentWithdrawal.value = withdrawal;
    approvalForm.value = {
      action: 'approve',
      notes: ''
    };
    showDetailModal.value = true;
  };

  // 提交审批
  const submitApproval = async () => {
    try {
      processing.value = true;
      let response: any;

      if (approvalForm.value.action === 'approve') {
        response = await distributorApi.approveWithdrawal(
          currentWithdrawal.value.id,
          approvalForm.value.notes
        );
      } else {
        response = await distributorApi.rejectWithdrawal(
          currentWithdrawal.value.id,
          approvalForm.value.notes
        );
      }

      if (response.code === 200) {
        alert(approvalForm.value.action === 'approve' ? '已批准提现' : '已拒绝提现');
        showDetailModal.value = false;
        
        // 重新加载
        page.value = 1;
        finished.value = false;
        withdrawals.value = [];
        await loadWithdrawals();
      } else {
        alert(response.message || '操作失败');
      }
    } catch (error: any) {
      alert(error.response?.data?.message || '操作失败');
    } finally {
      processing.value = false;
    }
  };

  const getStatusType = (status: string) => {
    const types: Record<string, string> = {
      pending: 'warning',
      approved: 'primary',
      rejected: 'danger',
      processing: 'info',
      completed: 'success',
      failed: 'danger'
    };
    return types[status] || 'default';
  };

  const getStatusText = (status: string) => {
    const texts: Record<string, string> = {
      pending: '待审核',
      approved: '已批准',
      rejected: '已拒绝',
      processing: '处理中',
      completed: '已完成',
      failed: '失败'
    };
    return texts[status] || status;
  };

  const getPayTypeText = (payType: string) => {
    const texts: Record<string, string> = {
      wechat: '微信',
      alipay: '支付宝',
      bank: '银行卡'
    };
    return texts[payType] || payType;
  };

  const formatDate = (date: string) => {
    if (!date) return '';
    const d = new Date(date);
    return `${d.getFullYear()}-${(d.getMonth() + 1).toString().padStart(2, '0')}-${d.getDate().toString().padStart(2, '0')} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`;
  };

  const finished = ref(false);

  onMounted(() => {
    loadWithdrawals();
  });

  return {
    withdrawals,
    loading,
    refreshing,
    page,
    pageSize,
    total,
    selectedBrandId,
    statusFilter,
    showDetailModal,
    currentWithdrawal,
    approvalForm,
    processing,
    finished,
    onRefresh,
    onLoad,
    loadWithdrawals,
    openDetailModal,
    submitApproval,
    getStatusType,
    getStatusText,
    getPayTypeText,
    formatDate
  };
};

// 渲染提现审批视图
export const renderWithdrawalApprovalView = (viewModel: ReturnType<typeof WithdrawalApprovalView>) => {
  const { h } = (window as any).Vue || { h: () => null };

  return h('div', { class: 'space-y-6' }, [
    // 头部
    h('div', { class: 'flex justify-between items-center' }, [
      h('div', [
        h('h2', { class: 'text-2xl font-black text-slate-900' }, '提现审批'),
        h('p', { class: 'text-slate-400 text-sm mt-1' }, '平台管理员审批分销商提现申请'),
      ]),
      h('div', { class: 'flex gap-2' }, [
        h('select', {
          class: 'border rounded-lg px-3 py-2',
          value: selectedBrandId.value || '',
          onInput: (e: any) => {
            viewModel.selectedBrandId.value = e.target.value === '' ? null : Number(e.target.value);
            page.value = 1;
            finished.value = false;
            withdrawals.value = [];
            loadWithdrawals();
          }
        }, [
          h('option', { value: '' }, '全部品牌'),
          h('option', { value: 1 }, '品牌1'),
          h('option', { value: 2 }, '品牌2'),
          h('option', { value: 3 }, '品牌3'),
        ]),
        h('select', {
          class: 'border rounded-lg px-3 py-2',
          value: statusFilter.value,
          onInput: (e: any) => {
            statusFilter.value = e.target.value;
            page.value = 1;
            finished.value = false;
            withdrawals.value = [];
            loadWithdrawals();
          }
        }, [
          h('option', { value: 'pending' }, '待审核'),
          h('option', { value: 'approved' }, '已批准'),
          h('option', { value: 'rejected' }, '已拒绝'),
          h('option', { value: 'completed' }, '已完成'),
        ]),
      ]),
    ]),

    // 提现列表
    h('div', { class: 'bg-white rounded-lg border shadow-sm' }, [
      h('div', { class: 'p-6 border-b' }, [
        h('table', { class: 'w-full' }, [
          h('thead'),
          h('tr', [
            h('th', { class: 'text-left px-4 py-3' }, '金额'),
            h('th', { class: 'text-left px-4 py-3' }, '提现方式'),
            h('th', { class: 'text-left px-4 py-3' }, '提现账号'),
            h('th', { class: 'text-left px-4 py-3' }, '真实姓名'),
            h('th', { class: 'text-left px-4 py-3' }, '申请时间'),
            h('th', { class: 'text-left px-4 py-3' }, '状态'),
            h('th', { class: 'text-left px-4 py-3' }, '操作'),
          ]),
          h('tbody', { class: viewModel.loading.value ? 'animate-pulse' : '' }, [
            viewModel.withdrawals.value.map((withdrawal: any) => h('tr', { class: 'border-t' }, [
              h('td', { class: 'px-4 py-3' }, `¥${withdrawal.amount.toFixed(2)}`),
              h('td', { class: 'px-4 py-3' }, getPayTypeText(withdrawal.payType)),
              h('td', { class: 'px-4 py-3' }, withdrawal.payAccount),
              h('td', { class: 'px-4 py-3' }, withdrawal.payRealName),
              h('td', { class: 'px-4 py-3 text-slate-500' }, formatDate(withdrawal.createdAt)),
              h('td', { class: 'px-4 py-3' }, [
                h('span', { class: `px-2 py-1 rounded-full text-xs font-medium ${getStatusType(withdrawal.status)}` }, getStatusText(withdrawal.status)),
              ]),
              h('td', { class: 'px-4 py-3' }, [
                h('button', {
                  class: 'text-blue-600 hover:text-blue-800 mr-2',
                  onClick: () => openDetailModal(withdrawal)
                }, '详情'),
                withdrawal.status === 'pending' && h('button', {
                  class: `px-3 py-1 rounded ${viewModel.processing.value ? 'opacity-50 cursor-not-allowed' : ''}`,
                  disabled: viewModel.processing.value,
                  onClick: () => openDetailModal(withdrawal)
                }, '审批'),
              ]),
            ]),
          ])),
          viewModel.withdrawals.value.length === 0 && h('tr', [
            h('td', { colSpan: 8, class: 'px-4 py-8 text-center text-slate-400' }, '暂无提现记录'),
          ]),
        ]),
      ]),
    ]),

    // 加载更多
    h('div', { class: 'mt-4 flex justify-center' }, [
      h('button', {
        class: `px-4 py-2 bg-white border text-slate-700 rounded-lg hover:bg-slate-50 ${finished.value ? 'opacity-50 cursor-not-allowed' : ''}`,
        disabled: finished.value || viewModel.loading.value,
        onClick: () => onLoad()
      }, finished.value ? '没有更多了' : '加载更多'),
    ]),

    // 详情弹窗
    showDetailModal.value && h('div', { class: 'fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50' }, [
      h('div', { class: 'bg-white rounded-lg shadow-xl p-6 w-full max-w-md' }, [
        h('h3', { class: 'text-lg font-semibold mb-4' }, '提现详情'),
        h('div', { class: 'space-y-3' }, [
          h('div', { class: 'flex justify-between' }, [
            h('span', { class: 'text-slate-600' }, '提现金额'),
            h('span', { class: 'font-semibold' }, `¥${currentWithdrawal.value?.amount?.toFixed(2) || '0.00'}`),
          ]),
          h('div', { class: 'flex justify-between' }, [
            h('span', { class: 'text-slate-600' }, '提现方式'),
            h('span', {}, getPayTypeText(currentWithdrawal.value?.payType || '')),
          ]),
          h('div', { class: 'flex justify-between' }, [
            h('span', { class: 'text-slate-600' }, '提现账号'),
            h('span', {}, currentWithdrawal.value?.payAccount || ''),
          ]),
          h('div', { class: 'flex justify-between' }, [
            h('span', { class: 'text-slate-600' }, '真实姓名'),
            h('span', {}, currentWithdrawal.value?.payRealName || ''),
          ]),
          h('div', { class: 'flex justify-between' }, [
            h('span', { class: 'text-slate-600' }, '申请时间'),
            h('span', { class: 'text-slate-500' }, formatDate(currentWithdrawal.value?.createdAt || '')),
          ]),
          h('div', { class: 'flex justify-between' }, [
            h('span', { class: 'text-slate-600' }, '当前状态'),
            h('span', { class: `px-2 py-1 rounded-full text-xs font-medium ${getStatusType(currentWithdrawal.value?.status || '')}` }, getStatusText(currentWithdrawal.value?.status || '')),
          ]),
        ]),
        currentWithdrawal.value?.rejectedReason && h('div', { class: 'text-sm text-red-600 mt-2' }, [
          h('span', { class: 'font-medium' }, '拒绝原因：'),
          h('p', {}, currentWithdrawal.value.rejectedReason),
        ]),
        
        // 审批操作（仅pending状态）
        currentWithdrawal.value?.status === 'pending' && h('div', { class: 'flex justify-end gap-2 mt-6 pt-6 border-t' }, [
          h('button', {
            class: 'px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 mr-2',
            onClick: () => approvalForm.value = { action: 'reject', notes: '' }
          }, '拒绝'),
          h('button', {
            class: 'px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700',
            onClick: () => approvalForm.value = { action: 'approve', notes: '' }
          }, '批准'),
        ]),

        // 审批表单
        (approvalForm.value.action === 'reject' || approvalForm.value.action === 'approve') && h('div', { class: 'mt-4 space-y-3' }, [
          approvalForm.value.action === 'reject' && h('div', [
            h('label', { class: 'block text-sm font-medium text-slate-700 mb-2' }, '拒绝原因（必填）'),
            h('textarea', {
              class: 'w-full border rounded-lg p-3 focus:outline-none focus:ring-2 focus:ring-blue-500',
              rows: 3,
              placeholder: '请输入拒绝原因',
              value: approvalForm.value.notes,
              onInput: (e: any) => approvalForm.value.notes = e.target.value
            }),
          ]),
          h('div', { class: 'flex justify-end gap-2 mt-4' }, [
            h('button', {
              class: 'px-4 py-2 border text-slate-700 rounded-lg hover:bg-slate-50 mr-2',
              onClick: () => showDetailModal.value = false
            }, '取消'),
            h('button', {
              class: `px-4 py-2 text-white rounded-lg ${viewModel.processing.value ? 'opacity-50 cursor-not-allowed' : ''}`,
              disabled: viewModel.processing.value || (approvalForm.value.action === 'reject' && !approvalForm.value.notes),
              onClick: () => submitApproval()
            }, approvalForm.value.action === 'reject' ? '确认拒绝' : '确认批准'),
          ]),
        ]),
        
        // 已审批的信息
        currentWithdrawal.value?.status !== 'pending' && currentWithdrawal.value?.status !== 'processing' && h('div', { class: 'mt-4 pt-4 border-t text-center' }, [
          currentWithdrawal.value?.status === 'approved' && h('div', { class: 'text-green-600' }, '已批准，等待系统处理打款'),
          currentWithdrawal.value?.status === 'completed' && h('div', { class: 'text-green-600' }, `打款完成'),
          currentWithdrawal.value?.status === 'rejected' && h('div', { class: 'text-red-600' }, '已拒绝'),
        ]),
        
        // 关闭按钮
        h('div', { class: 'flex justify-center mt-4' }, [
          h('button', {
            class: 'px-6 py-2 bg-slate-100 text-slate-700 rounded-lg hover:bg-slate-200',
            onClick: () => showDetailModal.value = false
          }, '关闭'),
        ]),
      ]),

      h('button', {
        class: 'absolute top-4 right-4 text-gray-400 hover:text-gray-600',
        onClick: () => showDetailModal.value = false
      }, '✕'),
    ]),

    // 遮罩层点击关闭
    h('div', {
      class: 'fixed inset-0',
      onClick: () => showDetailModal.value = false
    }),
  ]);
};
