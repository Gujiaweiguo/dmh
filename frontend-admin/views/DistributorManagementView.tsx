import { computed, defineComponent, h, onMounted, onUnmounted, ref, watch } from 'vue';
import { Clock, RefreshCw, UserCheck, UserX, Users } from 'lucide-vue-next';
import { distributorApi } from '../services/distributorApi';

export const DISTRIBUTOR_REQUEST_TIMEOUT_MS = 2000;
export const DISTRIBUTOR_SEARCH_DEBOUNCE_MS = 300;

export interface DistributorRecord {
  id: number;
  userId?: number;
  username?: string;
  brandId?: number;
  brandName?: string;
  level: number;
  parentName?: string;
  totalEarnings?: number;
  subordinatesCount?: number;
  status: string;
  createdAt?: string;
  approvedAt?: string;
}

interface DistributorStats {
  total: number;
  active: number;
  suspended: number;
  pending: number;
}

interface StatusInfo {
  style: string;
  label: string;
}

interface DistributorViewModel {
  distributors: { value: DistributorRecord[] };
  loading: { value: boolean };
  errorMessage: { value: string };
  searchKeyword: { value: string };
  statusFilter: { value: string };
  levelFilter: { value: number };
  page: { value: number };
  pageSize: { value: number };
  total: { value: number };
  stats: { value: DistributorStats };
  filteredDistributors: { value: DistributorRecord[] };
  loadDistributors: () => Promise<void>;
  goPrev: () => Promise<void>;
  goNext: () => Promise<void>;
  getStatusInfo: (status: string) => StatusInfo;
}

const createTimeoutError = () => {
  const error = new Error('REQUEST_TIMEOUT');
  error.name = 'TimeoutError';
  return error;
};

const isTimeoutError = (error: unknown) =>
  error instanceof Error && (error.name === 'TimeoutError' || error.message === 'REQUEST_TIMEOUT');

export const withTimeout = <T,>(promise: Promise<T>, timeoutMs = DISTRIBUTOR_REQUEST_TIMEOUT_MS): Promise<T> =>
  new Promise((resolve, reject) => {
    const timer = setTimeout(() => reject(createTimeoutError()), timeoutMs);
    promise
      .then((value) => {
        clearTimeout(timer);
        resolve(value);
      })
      .catch((error) => {
        clearTimeout(timer);
        reject(error);
      });
  });

export const computeDistributorStats = (items: DistributorRecord[]): DistributorStats => ({
  total: items.length,
  active: items.filter((item) => item.status === 'active').length,
  suspended: items.filter((item) => item.status === 'suspended').length,
  pending: items.filter((item) => item.status === 'pending').length,
});

export const filterDistributors = (
  items: DistributorRecord[],
  status: string,
  level: number,
  keyword: string,
): DistributorRecord[] => {
  let result = items;

  if (status) {
    result = result.filter((item) => item.status === status);
  }

  if (level > 0) {
    result = result.filter((item) => item.level === level);
  }

  const normalizedKeyword = keyword.trim().toLowerCase();
  if (!normalizedKeyword) {
    return result;
  }

  return result.filter((item) => {
    const username = item.username?.toLowerCase() ?? '';
    const brandName = item.brandName?.toLowerCase() ?? '';
    return username.includes(normalizedKeyword) || brandName.includes(normalizedKeyword);
  });
};

const statusStyleMap: Record<string, string> = {
  active: 'bg-emerald-100 text-emerald-800',
  suspended: 'bg-rose-100 text-rose-800',
  pending: 'bg-amber-100 text-amber-800',
};

const statusLabelMap: Record<string, string> = {
  active: '正常',
  suspended: '已暂停',
  pending: '待审核',
};

const parseSelectNumber = (event: Event) => {
  const target = event.target as HTMLSelectElement | null;
  const value = Number(target?.value ?? 0);
  return Number.isNaN(value) ? 0 : value;
};

export const DistributorManagementView = defineComponent({
  props: {
    brandId: Number,
    readOnly: Boolean,
    isPlatformAdmin: Boolean,
  },
  setup() {
    const distributors = ref<DistributorRecord[]>([]);
    const loading = ref(false);
    const errorMessage = ref('');
    const searchKeyword = ref('');
    const debouncedSearchKeyword = ref('');
    const statusFilter = ref('');
    const levelFilter = ref(0);
    const page = ref(1);
    const pageSize = ref(50);
    const total = ref(0);

    let debounceTimer: ReturnType<typeof setTimeout> | undefined;

    watch(searchKeyword, (value) => {
      if (debounceTimer) {
        clearTimeout(debounceTimer);
      }
      debounceTimer = setTimeout(() => {
        debouncedSearchKeyword.value = value;
      }, DISTRIBUTOR_SEARCH_DEBOUNCE_MS);
    });

    onUnmounted(() => {
      if (debounceTimer) {
        clearTimeout(debounceTimer);
      }
    });

    const stats = computed(() => computeDistributorStats(distributors.value));

    const filteredDistributors = computed(() =>
      filterDistributors(
        distributors.value,
        statusFilter.value,
        levelFilter.value,
        debouncedSearchKeyword.value,
      ),
    );

    const loadDistributors = async () => {
      loading.value = true;
      errorMessage.value = '';
      try {
        const response = await withTimeout(
          distributorApi.getGlobalDistributors(
            undefined,
            statusFilter.value || undefined,
            page.value,
            pageSize.value,
          ),
        );

        const fetchedList = Array.isArray(response?.distributors)
          ? (response.distributors as DistributorRecord[])
          : [];
        distributors.value = fetchedList;
        total.value = typeof response?.total === 'number' ? response.total : fetchedList.length;
      } catch (error) {
        distributors.value = [];
        total.value = 0;
        errorMessage.value = isTimeoutError(error) ? '请求超时，请重试' : '加载失败，请稍后重试';
      } finally {
        loading.value = false;
      }
    };

    const goPrev = async () => {
      if (page.value <= 1 || loading.value) {
        return;
      }
      page.value -= 1;
      await loadDistributors();
    };

    const goNext = async () => {
      if (loading.value) {
        return;
      }
      const maxPage = Math.max(1, Math.ceil((total.value || 0) / (pageSize.value || 1)));
      if (page.value >= maxPage) {
        return;
      }
      page.value += 1;
      await loadDistributors();
    };

    const getStatusInfo = (status: string): StatusInfo => ({
      style: statusStyleMap[status] || 'bg-slate-100 text-slate-800',
      label: statusLabelMap[status] || status,
    });

    onMounted(() => {
      void loadDistributors();
    });

    return () =>
      renderDistributorManagementView({
        distributors,
        loading,
        errorMessage,
        searchKeyword,
        statusFilter,
        levelFilter,
        page,
        pageSize,
        total,
        stats,
        filteredDistributors,
        loadDistributors,
        goPrev,
        goNext,
        getStatusInfo,
      });
  },
});

const renderStatCards = (viewModel: DistributorViewModel) =>
  h('div', { class: 'grid grid-cols-1 md:grid-cols-4 gap-6 mb-6' }, [
    h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm' }, [
      h('div', { class: 'w-12 h-12 bg-indigo-600 text-white rounded-2xl flex items-center justify-center mb-4' },
        h(Users, { size: 24 }),
      ),
      h('p', { class: 'text-[10px] font-black text-slate-400 uppercase tracking-widest' }, '总数'),
      h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, String(viewModel.stats.value.total)),
    ]),
    h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm' }, [
      h('div', { class: 'w-12 h-12 bg-emerald-600 text-white rounded-2xl flex items-center justify-center mb-4' },
        h(UserCheck, { size: 24 }),
      ),
      h('p', { class: 'text-[10px] font-black text-slate-400 uppercase tracking-widest' }, '正常'),
      h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, String(viewModel.stats.value.active)),
    ]),
    h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm' }, [
      h('div', { class: 'w-12 h-12 bg-rose-600 text-white rounded-2xl flex items-center justify-center mb-4' },
        h(UserX, { size: 24 }),
      ),
      h('p', { class: 'text-[10px] font-black text-slate-400 uppercase tracking-widest' }, '暂停'),
      h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, String(viewModel.stats.value.suspended)),
    ]),
    h('div', { class: 'bg-white p-6 rounded-3xl border border-slate-100 shadow-sm' }, [
      h('div', { class: 'w-12 h-12 bg-amber-600 text-white rounded-2xl flex items-center justify-center mb-4' },
        h(Clock, { size: 24 }),
      ),
      h('p', { class: 'text-[10px] font-black text-slate-400 uppercase tracking-widest' }, '待审核'),
      h('p', { class: 'text-3xl font-black text-slate-900 mt-2' }, String(viewModel.stats.value.pending)),
    ]),
  ]);

const renderDistributorManagementView = (viewModel: DistributorViewModel) => {
  const tableHeaders = ['ID', '分销用户', '品牌', '级别', '上级', '累计收益', '下级数', '状态', '加入时间', '审批时间'];

  return h('div', { class: 'space-y-6' }, [
    renderStatCards(viewModel),
    h('div', { class: 'flex justify-between items-center' }, [
      h('div', [
        h('h2', { class: 'text-2xl font-black text-slate-900' }, '分销监控'),
        h('p', { class: 'text-slate-400 text-sm mt-1' }, '查看拥有分销权限的普通用户分销情况（只读）'),
      ]),
      h(
        'button',
        {
          onClick: () => {
            void viewModel.loadDistributors();
          },
          class:
            'bg-indigo-600 text-white px-6 py-3 rounded-2xl font-bold hover:bg-indigo-700 transition-colors flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed',
          disabled: viewModel.loading.value,
        },
        [h(RefreshCw, { size: 18 }), '刷新'],
      ),
    ]),
    h('div', { class: 'flex gap-4 items-center' }, [
      h('input', {
        value: viewModel.searchKeyword.value,
        onInput: (event: Event) => {
          const target = event.target as HTMLInputElement | null;
          viewModel.searchKeyword.value = target?.value ?? '';
        },
        placeholder: '搜索用户名或品牌...（输入防抖）',
        class: 'flex-1 max-w-xs px-4 py-2 border border-slate-200 rounded-xl',
      }),
      h(
        'select',
        {
          value: viewModel.statusFilter.value,
          onChange: (event: Event) => {
            const target = event.target as HTMLSelectElement | null;
            viewModel.statusFilter.value = target?.value ?? '';
          },
          class: 'px-4 py-2 border border-slate-200 rounded-xl',
        },
        [
          h('option', { value: '' }, '全部状态'),
          h('option', { value: 'active' }, '正常'),
          h('option', { value: 'suspended' }, '已暂停'),
          h('option', { value: 'pending' }, '待审核'),
        ],
      ),
      h(
        'select',
        {
          value: viewModel.levelFilter.value,
          onChange: (event: Event) => {
            viewModel.levelFilter.value = parseSelectNumber(event);
          },
          class: 'px-4 py-2 border border-slate-200 rounded-xl',
        },
        [
          h('option', { value: 0 }, '全部级别'),
          h('option', { value: 1 }, '一级'),
          h('option', { value: 2 }, '二级'),
          h('option', { value: 3 }, '三级'),
        ],
      ),
      h(
        'button',
        {
          onClick: () => {
            viewModel.page.value = 1;
            void viewModel.loadDistributors();
          },
          class: 'px-6 py-2 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed',
          disabled: viewModel.loading.value,
        },
        '搜索',
      ),
    ]),
    h('div', { class: 'bg-white rounded-3xl border border-slate-100 overflow-hidden shadow-sm' }, [
      viewModel.loading.value
        ? h('div', { class: 'p-12 text-center text-slate-500' }, '加载中...')
        : viewModel.errorMessage.value
          ? h('div', { class: 'p-12 text-center space-y-3' }, [
              h('p', { class: 'text-rose-600 font-medium' }, viewModel.errorMessage.value),
              h(
                'button',
                {
                  class: 'px-5 py-2 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700',
                  onClick: () => {
                    void viewModel.loadDistributors();
                  },
                },
                '重试',
              ),
            ])
          : viewModel.filteredDistributors.value.length === 0
            ? h('div', { class: 'p-12 text-center' }, [
                h(Users, { size: 48, class: 'mx-auto text-slate-300 mb-4' }),
                h('p', { class: 'text-slate-500 text-lg' }, '暂无数据'),
              ])
            : h('table', { class: 'w-full text-left' }, [
                h('thead', { class: 'bg-slate-50' }, [
                  h(
                    'tr',
                    tableHeaders.map((header) =>
                      h(
                        'th',
                        { class: 'px-6 py-4 text-xs font-black text-slate-400 uppercase tracking-widest' },
                        header,
                      ),
                    ),
                  ),
                ]),
                h(
                  'tbody',
                  viewModel.filteredDistributors.value.map((item) =>
                    h('tr', { class: 'border-b border-slate-50 last:border-0 hover:bg-slate-50/40' }, [
                      h('td', { class: 'px-6 py-4 text-sm text-slate-400 font-mono' }, String(item.id)),
                      h('td', { class: 'px-6 py-4' }, [
                        h('div', { class: 'text-sm font-bold text-slate-900' }, item.username || `用户${item.userId}`),
                      ]),
                      h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, item.brandName || `品牌${item.brandId}`),
                      h('td', { class: 'px-6 py-4 text-sm' }, [
                        h(
                          'span',
                          {
                            class: `px-2 py-1 rounded-lg text-xs font-bold ${
                              item.level === 1
                                ? 'bg-purple-100 text-purple-800'
                                : item.level === 2
                                  ? 'bg-blue-100 text-blue-800'
                                  : 'bg-cyan-100 text-cyan-800'
                            }`,
                          },
                          `${item.level}级`,
                        ),
                      ]),
                      h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, item.parentName || '-'),
                      h(
                        'td',
                        { class: 'px-6 py-4 text-sm font-medium text-emerald-600' },
                        `¥${(item.totalEarnings ?? 0).toFixed(2)}`,
                      ),
                      h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, String(item.subordinatesCount ?? 0)),
                      h('td', { class: 'px-6 py-4' }, [
                        (() => {
                          const statusInfo = viewModel.getStatusInfo(item.status);
                          return h(
                            'span',
                            { class: `px-2 py-1 rounded-lg text-xs font-bold ${statusInfo.style}` },
                            statusInfo.label,
                          );
                        })(),
                      ]),
                      h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, item.createdAt || '-'),
                      h('td', { class: 'px-6 py-4 text-sm text-slate-600' }, item.approvedAt || '-'),
                    ]),
                  ),
                ),
              ]),
    ]),
    h('div', { class: 'flex items-center justify-between' }, [
      h('div', { class: 'text-sm text-slate-500' }, `第 ${viewModel.page.value} 页 · 每页 ${viewModel.pageSize.value} 条`),
      h('div', { class: 'flex gap-2' }, [
        h(
          'button',
          {
            onClick: () => {
              void viewModel.goPrev();
            },
            class: 'px-4 py-2 text-slate-700 bg-slate-100 rounded-lg hover:bg-slate-200 disabled:opacity-50',
            disabled: viewModel.page.value <= 1 || viewModel.loading.value,
          },
          '上一页',
        ),
        h(
          'button',
          {
            onClick: () => {
              void viewModel.goNext();
            },
            class: 'px-4 py-2 text-slate-700 bg-slate-100 rounded-lg hover:bg-slate-200 disabled:opacity-50',
            disabled:
              viewModel.total.value <= viewModel.page.value * viewModel.pageSize.value ||
              viewModel.loading.value,
          },
          '下一页',
        ),
      ]),
    ]),
  ]);
};
